package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"time"
)

func main() {

	c := make(chan int)
	dialer := websocket.DefaultDialer
	for i := range 10 {
		time.Sleep(1 * time.Millisecond)

		go func(i int) {
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

		}(i)
	}

	<-c
}
