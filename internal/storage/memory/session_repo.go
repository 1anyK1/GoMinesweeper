package memory

import (
	"errors"
	"sync"

	"minesweeper/internal/domain"
)

type SessionRepo struct {
	mu       sync.RWMutex
	sessions map[string]*domain.Session
}

func NewSessionRepo() *SessionRepo {
	return &SessionRepo{
		sessions: make(map[string]*domain.Session),
	}
}

func (r *SessionRepo) Save(session *domain.Session) error {
	if session == nil {
		return errors.New("session is nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.sessions[session.ID] = session
	return nil
}

func (r *SessionRepo) Get(id string) (*domain.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	session, ok := r.sessions[id]
	if !ok {
		return nil, errors.New("session not found")
	}

	return session, nil
}

func (r *SessionRepo) List() []*domain.Session {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.Session, 0, len(r.sessions))
	for _, session := range r.sessions {
		result = append(result, session)
	}

	return result
}
