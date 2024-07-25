package main

import (
	"fmt"
	"net/http"
	"strings"
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
	fmt.Println(r.URL.Path)
	var (
		path  = strings.Trim(r.URL.Path, "/")
		parts = strings.Split(path, "/")
	)

	if r.Method == "GET" {
		fmt.Println("GET")
	}

	if r.Method == "POST" {
		if len(parts) != 2 {
			fmt.Println("Invalid action")
			return
		}
		topic := parts[1]
		fmt.Println("topic", topic)
		p.producerch <- Message{
			Data:  []byte("We don't know yet"),
			Topic: parts[1],
		}
	}
	fmt.Println("Parts", parts)
}

func (p *HTTPProducer) Start() error {
	fmt.Println("HTTP transport started", "port", p.listenAddr)
	return http.ListenAndServe(p.listenAddr, p)
}
