package main

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type WSMessage struct {
	Action string   `json:"action"`
	Topics []string `json:"topic"`
}

func main() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:4000", nil)
	if err != nil {
		log.Fatal(err)
	}

	msg := WSMessage{
		Action: "subscribe",
		Topics: []string{"TestConsumer"},
	}

	if err := conn.WriteJSON(msg); err != nil {
		log.Fatal(err)
	}

	for {
		msg := WSMessage{}
		if err := conn.ReadJSON(msg); err != nil {
			log.Fatal(err)
		}
		fmt.Println("msg", msg)
	}
}
