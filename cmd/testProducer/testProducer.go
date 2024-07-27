package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"golang.org/x/exp/slog"
)

func main() {
	url := "http://localhost:3000/publish"
	topics := []string{"topic_1"}

	for i := 0; i < len(topics); i++ {
		topic := topics[rand.Intn(len(topics))]
		payload := []byte(fmt.Sprintf("payload_%d", i))
		resp, err := http.Post(url+"/"+topic, "application/octet-stream", bytes.NewReader(payload))
		slog.Info("Producer", "resp", resp)
		if err != nil {
			log.Fatal(err)
		}

		if resp.StatusCode != http.StatusOK {
			log.Fatal("status code is not 200")
		}
	}
}
