package app

import (
	"log/slog"

	"minesweeper/internal/config"
	"minesweeper/internal/service"
	"minesweeper/internal/transport/tcp"
)

type Server struct {
	cfg    config.Config
	logger *slog.Logger
}

func NewServer(cfg config.Config, logger *slog.Logger) *Server {
	return &Server{
		cfg:    cfg,
		logger: logger,
	}
}

func (s *Server) Run() error {
	gameService := service.NewGameService(s.cfg.GameSize, s.cfg.Mines)
	tcpServer := tcp.NewServer(s.cfg.TCPAddr, gameService, s.logger)

	return tcpServer.Run()
}
