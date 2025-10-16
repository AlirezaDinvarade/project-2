package main

import (
	"encoding/json"
	"time"
	"tolling/aggregator/client"
	types "tolling/types"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	consumer    *kafka.Consumer
	isRunning   bool
	calcservice CalculatorServicer
	aggClient   *client.HTTPClient
}

func NewKafkaConsumer(topic string, svc CalculatorServicer, aggClient *client.HTTPClient) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err

	}
	err = c.SubscribeTopics([]string{topic}, nil)
	return &KafkaConsumer{
		consumer:    c,
		calcservice: svc,
		aggClient:   aggClient,
	}, nil
}

func (c *KafkaConsumer) Start() {
	logrus.Info("kafka transport started")
	c.isRunning = true
	c.readMessageLoop()
}

func (c *KafkaConsumer) readMessageLoop() {
	for c.isRunning {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("Kafka consume error %s", err)
			continue
		}
		var data types.OBUData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("serialization error  %s ", err)
			continue
		}
		distance, err := c.calcservice.CalculateDistance(data)
		if err != nil {
			logrus.Errorf("calculation error  %s ", err)
			continue
		}
		req := types.Distance{
			Value: distance,
			Unix:  time.Now().UnixNano(),
			OBUID: data.OBUID,
		}
		if err := c.aggClient.AggregateInvoice(req); err != nil {
			logrus.Errorf("Aggregate error:", err)
			continue
		}
	}
}
