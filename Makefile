# Makefile for mqueue

# Variables
OUTPUT_PATH = bin/mqueue

build:
	go build -o bin/mqueue

run: build
	./bin/mqueue

test:
	go run cmd/test/main.go

