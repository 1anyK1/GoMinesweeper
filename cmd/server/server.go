package main

import (
	"log"

	"minesweeper/internal/app"
	"minesweeper/internal/config"
	"minesweeper/internal/logger"
)

func main() {
	cfg := config.Load()
	logg := logger.New()

	server := app.NewServer(cfg, logg)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
