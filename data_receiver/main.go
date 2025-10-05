package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	types "tolling/Types"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gorilla/websocket"
)

var kafkaTopic = "obuData"

func main() {
	receiver, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/ws", receiver.HandleWS)
	http.ListenAndServe(":3000", nil)
	fmt.Println("data receiver work fine ")
}

type DataReceiver struct {
	messageChan chan types.OBUData
	conn        *websocket.Conn
	prod        *kafka.Producer
}

func (dr *DataReceiver) HandleWS(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
	}
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn

	go dr.weReceiveLoop()
}

func (dr *DataReceiver) ProduceData(data types.OBUData) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = dr.prod.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kafkaTopic, Partition: kafka.PartitionAny},
		Value:          b,
	}, nil)
	return err
}

func NewDataReceiver() (*DataReceiver, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return nil, err
	}

	// start another go routin to check if we have delivered data
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

	return &DataReceiver{
		messageChan: make(chan types.OBUData, 1028),
		prod:        p,
	}, nil
}

func (dr *DataReceiver) weReceiveLoop() {
	fmt.Println("New OBU connected client connected !")
	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("read error: ", err)
			continue
		}
		fmt.Printf("received OBU data from %d :: lat %.2f and lon %.2f \n", data.OBUID, data.Lat, data.Lon)
		if err := dr.ProduceData(data); err != nil {
			fmt.Println("Kafka produce error ", err)
		}
	}
}
