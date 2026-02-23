.PHONY: build run docker-up docker-down seed test

build:
	go build -o bin/cargomax-api ./cmd/server

run:
	go run ./cmd/server

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

seed:
	go run ./cmd/server --seed

test:
	go test ./...
