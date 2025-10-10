package main

import (
	"log"
	"tolling/aggregator/client"
)

const (
	kafkaTopic         = "obuData"
	aggregatorEndpoint = "http://localhost:3000/aggregate"
)

func main() {
	var svc CalculatorServicer
	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc, client.NewClient(aggregatorEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
