package main

import (
	"log"

	"github.com/gorilla/websocket"
	"golang.org/x/exp/slog"
)

type WSMessage struct {
	Action string   `json:"action"`
	Topics []string `json:"topics"`
}

func main() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:4000", nil)
	if err != nil {
		log.Fatal(err)
	}

	msg := WSMessage{
		Action: "subscribe",
		Topics: []string{"topic_1"},
	}

	if err := conn.WriteJSON(msg); err != nil {
		slog.Info("Consumer", "write error", err)
	}

	for {
		// var msg WSResponse
		var msg []byte
		_, msg, err := conn.ReadMessage()
		if err != nil {
			slog.Info("Consumer", "read error", err)
		}
		slog.Info("Consumer", "response", msg)
	}
}
