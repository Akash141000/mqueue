# Makefile for mqueue

# Variables
OUTPUT_PATH = bin/mqueue

build:
	go build -o bin/mqueue

run: build
	./bin/mqueue

testConsumer:
	go run cmd/testConsumer/testConsumer.go

testProducer:
	go run cmd/testProducer/testProducer.go

