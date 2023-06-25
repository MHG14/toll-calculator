package main

import (
	"fmt"
	"log"

	"github.com/mhg14/toll-calculator/aggregator/client"
)

const (
	kafkaTopic         = "obudata"
	aggregatorEndpoint = "http://127.0.0.1:5000"
)

func main() {
	var (
		err error
		svc CalculatorServicer
	)
	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)
	httpClient := client.NewHTTPClient(aggregatorEndpoint)
	// grpcClient, err := client.NewGRPCClient(aggregatorEndpoint)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc, httpClient)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
	fmt.Println("this is working just fine !!")
}
