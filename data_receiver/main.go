package main

import (
	"fmt"
	"log"
	"net/http"
	types "tolling/Types"

	"github.com/gorilla/websocket"
)

func main() {
	reveiver := NewDataReceiver()
	http.HandleFunc("/ws", reveiver.HandleWS)
	http.ListenAndServe(":3000", nil)
	fmt.Println("data receiver work fine ")
}

type DataReceiver struct {
	messageChan chan types.OBUData
	conn        *websocket.Conn
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

func NewDataReceiver() *DataReceiver {
	return &DataReceiver{
		messageChan: make(chan types.OBUData, 1028),
	}
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
		dr.messageChan <- data
	}
}
