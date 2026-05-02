package main

import (
	"log"
	"net"

	"minesweeper/internal/transport/tcp"
)

const serverAddr = "localhost:8080"

func main() {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := tcp.NewClient(conn)
	client.Run()
}
