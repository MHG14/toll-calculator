package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/mhg14/toll-calculator/types"
)

type DataReciever struct {
	msgch chan types.OBUData
	conn  *websocket.Conn
	prod  DataProducer
}

func NewDataReciever() (*DataReciever, error) {
	var (
		p   DataProducer
		err error
		kafkaTopic = "obudata"
	)
	p, err = NewKafkaProducer(kafkaTopic)
	if err != nil {
		return nil, err
	}
	p = NewLogMiddleware(p)
	return &DataReciever{
		msgch: make(chan types.OBUData, 128),
		prod:  p,
	}, nil
}

func (dr *DataReciever) produceData(data types.OBUData) error {
	return dr.prod.ProduceData(data)
}

func (dr *DataReciever) handleWS(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	dr.conn = conn
	go dr.wsRecieveLoop()
}

func (dr *DataReciever) wsRecieveLoop() {
	fmt.Println("New OBU connected")
	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println(err)
			continue
		}

		if err := dr.produceData(data); err != nil {
			fmt.Println("Kafka produce error:", err)
		}

		// fmt.Printf("recieved OBU data from [%d] :: <Lat: %.2f, Long: %.2f>\n", data.OBUID, data.Lat, data.Long)
		// dr.msgch <- data
	}
}

func main() {
	reciever, err := NewDataReciever()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/ws", reciever.handleWS)
	http.ListenAndServe(":30000", nil)
}
