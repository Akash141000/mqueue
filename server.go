package main

import (
	"fmt"
	"log"
	"sync"

	"golang.org/x/exp/slog"
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
	consumerch chan string
	quitch     chan struct{}
	consumers  []Consumer
	peers      map[string][]Peer
	mu         sync.RWMutex
}

func NewServer(cfg *Config) (*Server, error) {
	slog.Info("Http server", "Producer", cfg.HTTPListenAddr)
	slog.Info("WebSocker server", "Consumer", cfg.WSListenAddr)

	producerch := make(chan Message)
	consumerch := make(chan string)

	s := &Server{
		Config:     cfg,
		topics:     make(map[string]Storer),
		quitch:     make(chan struct{}),
		producers:  []Producer{NewHTTPProducer(cfg.HTTPListenAddr, producerch)},
		producerch: producerch,
		consumerch: consumerch,
		consumers:  []Consumer{},
		peers:      make(map[string][]Peer),
	}

	s.consumers = append(s.consumers, NewWSConsumer(cfg.WSListenAddr, s))
	return s, nil
}

func (s *Server) Start() {
	slog.Info("Server", "start", " consumers")
	for _, consumer := range s.consumers {
		go func(c Consumer) {
			if err := c.Start(); err != nil {
				fmt.Println(err)
			}
		}(consumer)
	}

	slog.Info("Server", "start", " producers")
	for _, producer := range s.producers {
		go func(p Producer) {
			if err := p.Start(); err != nil {
				fmt.Println(err)
			}
		}(producer)
	}

	wg := sync.WaitGroup{}
	//even if one routine stop other has to stop, so adding only one
	wg.Add(1)

	go s.startPublising(&wg)
	go s.startConsuming(&wg)

	wg.Wait()
}

func (s *Server) startPublising(wg *sync.WaitGroup) {
	slog.Info("Server", "producing", "started")
	for {
		select {
		case <-s.quitch:
			wg.Done()
			return
		case msg := <-s.producerch:
			slog.Info("Producer", "topic="+string(msg.Topic), "data="+string(msg.Data))
			offset, err := s.publish(msg)
			s.consumerch <- msg.Topic
			if err != nil {
				log.Fatal("unable to publish", err)
			}
			slog.Info("Producer", "produced offset", offset)
		}
	}
}

func (s *Server) startConsuming(wg *sync.WaitGroup) {
	slog.Info("Server", "consuming", "started")
	for {
		select {
		case <-s.quitch:
			wg.Done()
			return
		case topic := <-s.consumerch:
			slog.Info("Consumer", "start consuming for topic", topic)
			store := s.topics[topic]
			for _, peer := range s.peers[topic] {
				for i := 0; i < store.Len(); i++ {
					msg, err := store.Pull(i)
					if err != nil {
						slog.Info("Error", "send msg to peer", err)
						continue
					}
					peer.Send(msg)
				}
			}
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
		slog.Info("Topic", "new", topic)
	}
	return s.topics[topic]
}

func (s *Server) AddPeerToTopics(p Peer, topics ...string) {
	slog.Info("Peer", "add-to-topics", topics)
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, topic := range topics {
		if _, ok := s.peers[topic]; !ok {
			s.peers[topic] = make([]Peer, 0)
		}
		s.peers[topic] = append(s.peers[topic], p)
	}
}
