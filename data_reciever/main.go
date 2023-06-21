package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gorilla/websocket"
	"github.com/mhg14/toll-calculator/types"
)

var kafkaTopic = "obudata"

type DataReciever struct {
	msgch chan types.OBUData
	conn  *websocket.Conn
	prod  *kafka.Producer
}

func NewDataReciever() (*DataReciever, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return nil, err
	}

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	return &DataReciever{
		msgch: make(chan types.OBUData, 128),
		prod:  p,
	}, nil
}

func (dr *DataReciever) produceData(data types.OBUData) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = dr.prod.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &kafkaTopic,
			Partition: kafka.PartitionAny},
		Value: b,
	}, nil)
	return err
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

		fmt.Printf("recieved OBU data from [%d] :: <Lat: %.2f, Long: %.2f>\n", data.OBUID, data.Lat, data.Long)
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
