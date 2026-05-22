package rag

import (
	"errors"
	"strings"
	"sync"
	"time"
)

var ErrSessionNotFound = errors.New("chat session not found")

type chatStore struct {
	mu       sync.RWMutex
	ids      idGenerator
	sessions map[string]ChatSession
}

func newChatStore() *chatStore {
	return &chatStore{
		sessions: make(map[string]ChatSession),
	}
}

func (s *chatStore) create(userID string) ChatSession {
	session := ChatSession{
		ID:        s.ids.next("session"),
		UserID:    normalizeUserID(userID),
		CreatedAt: time.Now().UTC(),
		Messages:  []ChatMessage{},
	}

	s.mu.Lock()
	s.sessions[session.ID] = session
	s.mu.Unlock()

	return session
}

func (s *chatStore) get(sessionID string) (ChatSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.sessions[sessionID]
	if !ok {
		return ChatSession{}, ErrSessionNotFound
	}
	session.Messages = append([]ChatMessage(nil), session.Messages...)
	return session, nil
}

func (s *chatStore) appendMessage(sessionID string, message ChatMessage) (ChatSession, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[sessionID]
	if !ok {
		return ChatSession{}, ErrSessionNotFound
	}

	session.Messages = append(session.Messages, message)
	s.sessions[sessionID] = session
	session.Messages = append([]ChatMessage(nil), session.Messages...)
	return session, nil
}

type answerGenerator struct{}

func newAnswerGenerator() answerGenerator {
	return answerGenerator{}
}

func (answerGenerator) answer(message string, chunks []RetrievedChunk, memory UserMemory) string {
	var builder strings.Builder
	if len(chunks) == 0 {
		builder.WriteString("I could not find indexed knowledge that matches the question yet.")
	} else {
		builder.WriteString("Based on the retrieved knowledge: ")
		for i, chunk := range chunks {
			if i > 0 {
				builder.WriteString(" ")
			}
			builder.WriteString(trimForAnswer(chunk.Content, 220))
		}
	}

	if memory.RecentSummary != "" {
		builder.WriteString(" User memory considered: ")
		builder.WriteString(memory.RecentSummary)
		builder.WriteString(".")
	}

	return builder.String()
}
