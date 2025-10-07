package main

import (
	"log"
)

const kafkaTopic = "obuData"

func main() {
	svc := NewCalculatorService()
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
