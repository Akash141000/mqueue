package main

import (
	"fmt"
	"log"
)

func StoreFunc() StoreProducerFunc {
	return func() Storer { return NewMemoryStore() }
}

func main() {
	fmt.Println("Starting mqueue")
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
