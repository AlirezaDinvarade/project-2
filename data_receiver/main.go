package main

import (
	"fmt"
	"log"
	"net/http"
	types "tolling/Types"

	"github.com/gorilla/websocket"
)


func main() {
	receiver, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/ws", receiver.HandleWS)
	http.ListenAndServe(":3000", nil)
}

type DataReceiver struct {
	messageChan chan types.OBUData
	conn        *websocket.Conn
	prod        DataProducer
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
	return dr.prod.ProduceData(data)
}

func NewDataReceiver() (*DataReceiver, error) {
	var p DataProducer
	p, err := NewKafkaProducer("obuData")
	if err != nil {
		return nil, err
	}
	p = NewLogMiddleware(p)
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
