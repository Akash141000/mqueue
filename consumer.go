package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Consumer interface {
	Start() error
}

type WSConsumer struct {
	listenAddr string
	server     *Server
}

func NewWSConsumer(listenAddr string, server *Server) *WSConsumer {
	return &WSConsumer{
		listenAddr: listenAddr,
		server:     server,
	}
}

func (ws *WSConsumer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("error upgrading the connection", err)
	}
	NewWSPeer(conn, ws.server)
	// ws.server.AddPeer(p)
}

type WSMessage struct {
	Action string   `json:"action"`
	Topics []string `json:"topics"`
}

func (ws *WSConsumer) Start() error {
	return http.ListenAndServe(ws.listenAddr, ws)
}
