package service

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

const (
	sessionExpiry   = 5 * time.Minute
	cleanupInterval = 1 * time.Minute
)

type CLISession struct {
	Token     string
	CreatedAt time.Time
	Completed bool
}

type CLISessionService struct {
	mu       sync.RWMutex
	sessions map[string]*CLISession
}

var (
	cliSessionInstance *CLISessionService
	cliSessionOnce     sync.Once
)

func GetCLISessionService() *CLISessionService {
	cliSessionOnce.Do(func() {
		cliSessionInstance = &CLISessionService{
			sessions: make(map[string]*CLISession),
		}
		go cliSessionInstance.cleanupExpired()
	})
	return cliSessionInstance
}

func (s *CLISessionService) CreateSession() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	sessionID := hex.EncodeToString(bytes)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[sessionID] = &CLISession{
		CreatedAt: time.Now(),
		Completed: false,
	}

	return sessionID, nil
}

func (s *CLISessionService) CompleteSession(sessionID, token string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return false
	}

	if time.Since(session.CreatedAt) > sessionExpiry {
		delete(s.sessions, sessionID)
		return false
	}

	session.Token = token
	session.Completed = true
	return true
}

func (s *CLISessionService) PollSession(sessionID string) (token string, pending bool, valid bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return "", false, false
	}

	if time.Since(session.CreatedAt) > sessionExpiry {
		return "", false, false
	}

	if session.Completed {
		return session.Token, false, true
	}

	return "", true, true
}

func (s *CLISessionService) DeleteSession(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, sessionID)
}

func (s *CLISessionService) cleanupExpired() {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for id, session := range s.sessions {
			if now.Sub(session.CreatedAt) > sessionExpiry {
				delete(s.sessions, id)
			}
		}
		s.mu.Unlock()
	}
}
