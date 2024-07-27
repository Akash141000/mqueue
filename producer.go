package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/exp/slog"
)

type Producer interface {
	Start() error
}

type HTTPProducer struct {
	listenAddr string
	producerch chan<- Message
}

func NewHTTPProducer(listenAddr string, producerch chan Message) *HTTPProducer {
	return &HTTPProducer{
		listenAddr: listenAddr,
		producerch: producerch,
	}
}

func (p *HTTPProducer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("Producer", "url", r.URL)
	var (
		path  = strings.Trim(r.URL.Path, "/")
		parts = strings.Split(path, "/")
	)

	if r.Method == "POST" {
		if len(parts) != 2 {
			fmt.Println("Invalid action")
			return
		}
		topic := parts[1]

		payload, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal("error reading payload", err)
		}

		p.producerch <- Message{
			Data:  []byte(payload),
			Topic: topic,
		}
	}
}

func (p *HTTPProducer) Start() error {
	return http.ListenAndServe(p.listenAddr, p)
}
