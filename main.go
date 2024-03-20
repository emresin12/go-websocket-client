package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"sync"
	"time"
)

var n_conn = 100

type LatencyData struct {
	sync.Mutex
	latencies   []time.Duration
	clientCount int
}

func connectAndListen(i int, dialer *websocket.Dialer) {
	url := "ws://localhost:3000"

	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket server:", err)
	}
	defer conn.Close()
	fmt.Println("connection i=" + strconv.Itoa(i) + " connected")

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			fmt.Println("hata ", strconv.Itoa(i))
			break
		}

	}

}

func latency_test(i int, dialer *websocket.Dialer, latencyData *LatencyData, shutdown chan int) {
	url := "ws://localhost:3000"

	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket server:", err)
	}

	fmt.Println("connection i=" + strconv.Itoa(i) + " connected")
	message := make([]byte, 1024)
	latencies := make([]time.Duration, 0)
	for m := range 10 {

		startTime := time.Now()
		err2 := conn.WriteMessage(websocket.TextMessage, message)
		if err2 != nil {
			log.Println("error while sending message on client ", i, " try: ", m, " err: ", err2)
		}

		_, _, err2 = conn.ReadMessage()
		if err2 != nil {
			log.Println("error while reading message on client ", i, " try: ", m, " err: ", err2)
		}
		endTime := time.Now()

		latencies = append(latencies, endTime.Sub(startTime))
	}
	terminate := false
	latencyData.Lock()
	latencyData.latencies = append(latencyData.latencies, latencies...)
	latencyData.clientCount++
	if latencyData.clientCount == n_conn {
		terminate = true
	}
	latencyData.Unlock()
	closeErr := conn.Close()
	if closeErr != nil {
		log.Println("error while closing connection ", i, " err: ", closeErr)
	}

	if terminate {
		close(shutdown)
	}

}

func calculateLatencyStats(latencyData *LatencyData) {
	sumLatencies := time.Duration(0)
	minLatency := latencyData.latencies[0]
	maxLatency := latencyData.latencies[0]
	for _, latency := range latencyData.latencies {
		sumLatencies += latency
		if latency < minLatency {
			minLatency = latency
		}
		if latency > maxLatency {
			maxLatency = latency
		}
	}
	averageLatency := float64(sumLatencies.Microseconds()) / float64(len(latencyData.latencies))
	fmt.Printf("average latency: %f ms \n", averageLatency/float64(1000))
	fmt.Printf("min latency: %f ms \n", float64(minLatency.Microseconds())/float64(1000))
	fmt.Printf("max latency: %f ms \n", float64(maxLatency.Microseconds())/float64(1000))
}

func main() {

	c := make(chan int)
	dialer := websocket.DefaultDialer
	//latencyData := LatencyData{}
	//shutdown := make(chan int)
	for i := range n_conn {
		time.Sleep(100 * time.Microsecond)

		go connectAndListen(i, dialer)
		//go latency_test(i, dialer, &latencyData, shutdown)

	}

	//<-shutdown

	//calculateLatencyStats(&latencyData)

	<-c
}
