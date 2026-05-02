package domain

import "errors"

type Session struct {
	ID      string
	Game    *Game
	Players []Player
}

func NewSession(id string, game *Game) *Session {
	return &Session{
		ID:      id,
		Game:    game,
		Players: make([]Player, 0),
	}
}

func (s *Session) AddPlayer(player Player) error {
	if player.Name == "" {
		return errors.New("player name is empty")
	}

	for _, p := range s.Players {
		if p.Name == player.Name {
			return nil
		}
	}

	s.Players = append(s.Players, player)
	return nil
}
