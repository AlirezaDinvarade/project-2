package main

import (
	"log"
	"math"
	"math/rand"
	"time"
	types "tolling/Types"

	"github.com/gorilla/websocket"
)

const weEndpoint = "ws://127.0.0.1:3000/ws"

var sendInterval = time.Second * 5

func genCoords() float64 {
	n := float64(rand.Intn(100) + 1)
	f := rand.Float64()
	return n + f
}

func gerLocatin() (float64, float64) {
	return genCoords(), genCoords()
}

func main() {
	obuIDS := generateOBUIDS(20)
	conn, _, err := websocket.DefaultDialer.Dial(weEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		for i := range obuIDS {
			lat, lon := gerLocatin()
			data := types.OBUData{
				OBUID: obuIDS[i],
				Lat:   lat,
				Lon:   lon,
			}
			if err := conn.WriteJSON(data); err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(sendInterval)
	}
}

func generateOBUIDS(n int) []int {
	ids := make([]int, n)
	for i := range n {
		ids[i] = rand.Intn(math.MaxInt)
	}
	return ids
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
