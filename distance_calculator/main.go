package main

import (
	"log"
)

const kafkaTopic = "obuData"

func main() {
	var svc CalculatorServicer
	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
