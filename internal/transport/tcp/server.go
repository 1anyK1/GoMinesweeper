package tcp

import (
	"log/slog"
	"net"

	"minesweeper/internal/service"
)

type Server struct {
	addr    string
	handler *Handler
	logger  *slog.Logger
}

func NewServer(addr string, game *service.GameService, logger *slog.Logger) *Server {
	return &Server{
		addr:    addr,
		handler: NewHandler(game, logger),
		logger:  logger,
	}
}

func (s *Server) Run() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer ln.Close()

	s.logger.Info("tcp server started", "addr", s.addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			s.logger.Error("accept failed", "error", err)
			continue
		}

		go s.handler.Handle(conn)
	}
}
