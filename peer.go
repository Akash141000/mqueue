package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Peer interface {
	Send([]byte) error
}

type WSPeer struct {
	conn   *websocket.Conn
	server *Server
}

func NewWSPeer(conn *websocket.Conn, server *Server) *WSPeer {
	p := &WSPeer{
		conn:   conn,
		server: server,
	}

	go p.readLoop()

	return p
}

func (p *WSPeer) readLoop() {
	var msg WSMessage
	for {
		if err := p.conn.ReadJSON(&msg); err != nil {
			fmt.Println("ws peer read error", err)
			continue
		}

		if err := p.handleMessage(msg); err != nil {
			fmt.Println("ws peer handle msg error", err)
			continue
		}

	}
}

func (p *WSPeer) handleMessage(msg WSMessage) (err error) {
	fmt.Printf("handle message =>  %+v\n", msg)
	if msg.Action == "subscribe" {
		p.server.AddPeerToTopics(p, msg.Topics...)
	}
	return nil
}

func (p *WSPeer) Send(b []byte) (err error) {
	return p.conn.WriteMessage(websocket.BinaryMessage, b)
}
