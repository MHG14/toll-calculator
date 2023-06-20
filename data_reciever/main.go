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
}

func NewDataReciever() *DataReciever {
	return &DataReciever{
		msgch: make(chan types.OBUData, 128),
	}
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

		fmt.Printf("recieved OBU data from [%d] :: <Lat: %.2f, Long: %.2f>\n", data.OBUID, data.Lat, data.Long)
		dr.msgch <- data
	}
}

func main() {
	fmt.Println("reciever service working properly")
	reciever := NewDataReciever()
	http.HandleFunc("/ws", reciever.handleWS)
	http.ListenAndServe(":30000", nil)
}
