APP_NAME=minesweeper
SERVER_BIN=bin/server
CLIENT_BIN=bin/client
DOCKER_IMAGE=minesweeper-go

.PHONY: run-server run-client build test fmt vet clean docker-build docker-run

run-server:
	go run ./cmd/server

run-client:
	go run ./cmd/client

build:
	mkdir -p bin
	go build -o $(SERVER_BIN) ./cmd/server
	go build -o $(CLIENT_BIN) ./cmd/client

clean:
	rm -rf bin

docker-build:
	docker build -t $(DOCKER_IMAGE) .

docker-run:
	docker run --rm -p 8080:8080 \
		-e TCP_ADDR=:8080 \
		-e GAME_SIZE=8 \
		-e MINES=10 \
		$(DOCKER_IMAGE)