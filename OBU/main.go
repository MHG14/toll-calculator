package main

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mhg14/toll-calculator/types"
)

const wsURL = "ws://127.0.0.1:30000/ws"

var sendInterval = time.Second

func genCoord() float64 {
	n := float64(rand.Intn(100) + 1)
	f := rand.Float64()
	return n + f
}

func genLatLong() (float64, float64) {
	return genCoord(), genCoord()
}

func genOBUIDs(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(math.MaxInt)
	}
	return ids
}

func main() {
	OBUIDs := genOBUIDs(20)
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		for i := 0; i < len(OBUIDs); i++ {
			lat, long := genLatLong()
			data := types.OBUData{
				OBUID: OBUIDs[i],
				Lat:   lat,
				Long:  long,
			}
			if err := conn.WriteJSON(data); err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(sendInterval)
	}
}

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}
