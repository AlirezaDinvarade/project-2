package main

import (
	"log"
	"tolling/aggregator/client"
)

const (
	kafkaTopic         = "obuData"
	// httpAggregatorEndpoint = "http://localhost:3000/aggregate"
	grpcAggregatorEndpoint = "localhost:3001"
)

func main() {
	var svc CalculatorServicer
	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)

	// httpClient := client.NewHTTPClient(httpAggregatorEndpoint)
	grpcClient, err := client.NewGRPCClient(grpcAggregatorEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc, grpcClient)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
