package main

import (
	"errors"
	"sync"
)

type MemorySessionStore struct {
	mu       sync.Mutex
	sessions map[string]Session
}

func NewMemorySessionStore() *MemorySessionStore {
	return &MemorySessionStore{
		sessions: make(map[string]Session),
	}
}

func (s *MemorySessionStore) Add(session Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.sessions[session.ID]; exists {
		return errors.New("session already exists")
	}

	s.sessions[session.ID] = session
	return nil
}

func (s *MemorySessionStore) Get(id string) (Session, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.sessions[id]

	return session, exists
}

func (s *MemorySessionStore) ListOwned(
	principal Principal,
) []Session {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]Session, 0)

	for _, session := range s.sessions {
		if principal.CanAccessSession(session) {
			result = append(result, session)
		}
	}

	return result
}
