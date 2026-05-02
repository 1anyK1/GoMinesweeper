package service

import (
	"errors"
	"sync"

	"minesweeper/internal/domain"
)

type GameService struct {
	mu       sync.Mutex
	game     *domain.Game
	gameSize int
	mines    int
}

func NewGameService(gameSize, mines int) *GameService {
	return &GameService{
		gameSize: gameSize,
		mines:    mines,
	}
}

func (s *GameService) Create() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	game, err := domain.NewGame(s.gameSize, s.mines)
	if err != nil {
		return "", err
	}

	s.game = game

	return game.Render(), nil
}

func (s *GameService) State() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.game == nil {
		return "", errors.New("no active game")
	}

	return s.game.Render(), nil
}

func (s *GameService) Open(x, y int) (string, domain.GameStatus, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.game == nil {
		return "", domain.StatusPlaying, errors.New("no active game")
	}

	err := s.game.Open(x, y)

	return s.game.Render(), s.game.Status, err
}

func (s *GameService) ToggleFlag(x, y int) (string, domain.GameStatus, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.game == nil {
		return "", domain.StatusPlaying, errors.New("no active game")
	}

	err := s.game.ToggleFlag(x, y)

	return s.game.Render(), s.game.Status, err
}

func (s *GameService) Reset() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.game == nil {
		return "", errors.New("no active game")
	}

	s.game.Reset()

	return s.game.Render(), nil
}
