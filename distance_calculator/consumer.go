package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/mhg14/toll-calculator/aggregator/client"
	"github.com/mhg14/toll-calculator/types"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	consumer    *kafka.Consumer
	isRunning   bool
	calcService CalculatorServicer
	aggClient   client.Client
}

func NewKafkaConsumer(topic string, svc CalculatorServicer, aggClient client.Client) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}

	c.SubscribeTopics([]string{topic}, nil)

	return &KafkaConsumer{
		consumer:    c,
		calcService: svc,
		aggClient:   aggClient,
	}, nil
}

func (c *KafkaConsumer) Start() {
	c.isRunning = true
	logrus.Info("kafka transport started")
	c.readMessageLoop()
}

func (c *KafkaConsumer) readMessageLoop() {
	for c.isRunning {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("kafka consume error %s", err)
			continue
		}
		var data types.OBUData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("JSON serialization error %s", err)
			continue
		}
		distance, err := c.calcService.CalculateDistance(data)
		if err != nil {
			logrus.Errorf("distance calculation service error %s", err)
			continue
		}
		req := &types.AggregateRequest{
			OBUID: int32(data.OBUID),
			Value: distance,
			Unix:  time.Now().UnixNano(),
		}
		if err := c.aggClient.Aggregate(context.Background(), req); err != nil {
			logrus.Error("aggregate error:", err)
			continue
		}
	}
}
