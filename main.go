package main

import (
	"log"

	"golang.org/x/exp/slog"
)

func StoreFunc() StoreProducerFunc {
	return func() Storer { return NewMemoryStore() }
}

func main() {
	slog.Info("Server", "Mqueue server starting...")
	cfg := &Config{
		HTTPListenAddr:    ":3000",
		StoreProducerFunc: StoreFunc(),
		WSListenAddr:      ":4000",
	}
	s, err := NewServer(cfg)
	if err != nil {
		log.Fatal(err)
	}
	s.Start()
}
