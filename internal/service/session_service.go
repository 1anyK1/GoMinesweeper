package service

import (
	"fmt"

	"minesweeper/internal/domain"
	"minesweeper/internal/storage/memory"
)

type SessionService struct {
	repo     *memory.SessionRepo
	gameSize int
	mines    int
}

func NewSessionService(repo *memory.SessionRepo, gameSize, mines int) *SessionService {
	return &SessionService{
		repo:     repo,
		gameSize: gameSize,
		mines:    mines,
	}
}

func (s *SessionService) CreateSession(id string) (*domain.Session, error) {
	if id == "" {
		id = "default"
	}

	game, err := domain.NewGame(s.gameSize, s.mines)
	if err != nil {
		return nil, err
	}

	session := domain.NewSession(id, game)

	if err := s.repo.Save(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *SessionService) JoinSession(id, playerName string) error {
	session, err := s.repo.Get(id)
	if err != nil {
		return err
	}

	player := domain.NewPlayer(playerName)

	if err := session.AddPlayer(player); err != nil {
		return err
	}

	return s.repo.Save(session)
}

func (s *SessionService) ListSessions() string {
	sessions := s.repo.List()

	if len(sessions) == 0 {
		return "sessions: empty\n"
	}

	result := "sessions:\n"
	for _, session := range sessions {
		result += fmt.Sprintf("- %s players=%d\n", session.ID, len(session.Players))
	}

	return result
}
