package main

import (
	"github.com/gorilla/websocket"
	"golang.org/x/exp/slog"
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

	go p.readWrite()

	return p
}

func (p *WSPeer) readWrite() {
	var msg WSMessage
	for {
		if err := p.conn.ReadJSON(&msg); err != nil {
			slog.Info("Peer", "ws peer read error", err)
			break
		}

		if err := p.handleMessage(msg); err != nil {
			slog.Info("Peer", "ws peer handle msg error", err)
			continue
		}

	}
}

func (p *WSPeer) handleMessage(msg WSMessage) (err error) {
	slog.Info("Peer", "handle message", msg)
	if msg.Action == "subscribe" {
		p.server.AddPeerToTopics(p, msg.Topics...)
	}
	return nil
}

func (p *WSPeer) Send(b []byte) (err error) {
	return p.conn.WriteMessage(websocket.BinaryMessage, b)
}
