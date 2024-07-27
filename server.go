package main

import (
	"fmt"
	"sync"
)

type Message struct {
	Topic string
	Data  []byte
}

type Config struct {
	HTTPListenAddr    string
	WSListenAddr      string
	StoreProducerFunc StoreProducerFunc
}

type Server struct {
	*Config
	topics     map[string]Storer
	producers  []Producer
	producerch chan Message
	quitch     chan struct{}
	consumers  []Consumer
	peers      map[Peer]bool
	mu         sync.RWMutex
}

func NewServer(cfg *Config) (*Server, error) {
	producerch := make(chan Message)

	s := &Server{
		Config:     cfg,
		topics:     make(map[string]Storer),
		quitch:     make(chan struct{}),
		producers:  []Producer{NewHTTPProducer(cfg.HTTPListenAddr, producerch)},
		producerch: producerch,
		consumers:  []Consumer{},
		peers:      make(map[Peer]bool),
	}
	s.consumers = append(s.consumers, NewWSConsumer(cfg.WSListenAddr, s))
	return s, nil
}

func (s *Server) Start() {
	for _, consumer := range s.consumers {
		go func(c Consumer) {
			if err := c.Start(); err != nil {
				fmt.Println(err)
			}
		}(consumer)
	}

	for _, producer := range s.producers {
		go func(p Producer) {
			if err := p.Start(); err != nil {
				fmt.Println(err)
			}
		}(producer)

	}
	s.loop()
}

func (s *Server) loop() {
	for {
		select {
		case <-s.quitch:
			return
		case msg := <-s.producerch:
			fmt.Println("produced ->", msg)
			offset, err := s.publish(msg)
			if err != nil {
				fmt.Println("unable to publish", err)
			}
			fmt.Println("Produced offset", offset)
		}
	}
}

func (s *Server) publish(msg Message) (int, error) {
	store := s.getStoreForTopic(msg.Topic)
	return store.Push(msg.Data)
}

func (s *Server) getStoreForTopic(topic string) Storer {
	if _, ok := s.topics[topic]; !ok {
		s.topics[topic] = s.StoreProducerFunc()
		fmt.Println("create new topic", topic)
	}
	return s.topics[topic]
}

func (s *Server) AddPeer(p Peer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.peers[p] = true
}

func (s *Server) AddPeerToTopics(p Peer, topics ...string) {
	fmt.Println("adding peer to topics", topics)
	for _, topic := range topics {
		store := s.getStoreForTopic(topic)
		size := store.Len()
		for i := 0; i < size; i++ {
			b, _ := store.Pull(i)
			for p := range s.peers {
				p.Send(b)
			}
		}
	}
	fmt.Println("adding peer to topics", topics, "peers", p)
}
