package rag

import (
	"strings"
	"sync"
)

type memoryStore struct {
	mu      sync.RWMutex
	records map[string]UserMemory
}

func newMemoryStore() *memoryStore {
	return &memoryStore{
		records: make(map[string]UserMemory),
	}
}

func (s *memoryStore) get(userID string) UserMemory {
	userID = normalizeUserID(userID)

	s.mu.RLock()
	record, ok := s.records[userID]
	s.mu.RUnlock()
	if !ok {
		return UserMemory{UserID: userID}
	}

	record.Facts = append([]string(nil), record.Facts...)
	return record
}

func (s *memoryStore) addFact(userID string, fact string) (UserMemory, bool) {
	userID = normalizeUserID(userID)
	fact = strings.TrimSpace(fact)
	if fact == "" {
		return s.get(userID), false
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	record := s.records[userID]
	if record.UserID == "" {
		record.UserID = userID
	}
	for _, existing := range record.Facts {
		if strings.EqualFold(existing, fact) {
			record.Facts = append([]string(nil), record.Facts...)
			return record, false
		}
	}

	record.Facts = append(record.Facts, fact)
	if len(record.Facts) > 20 {
		record.Facts = record.Facts[len(record.Facts)-20:]
	}
	record.RecentSummary = summarizeFacts(record.Facts)
	s.records[userID] = record
	record.Facts = append([]string(nil), record.Facts...)
	return record, true
}

func (s *memoryStore) updateFromMessage(userID string, message string) []string {
	message = strings.TrimSpace(message)
	if message == "" {
		return nil
	}

	lower := strings.ToLower(message)
	prefixes := []string{"remember ", "my preference is ", "i prefer ", "i am studying "}
	for _, prefix := range prefixes {
		if strings.HasPrefix(lower, prefix) {
			fact := strings.TrimSpace(message[len(prefix):])
			_, added := s.addFact(userID, fact)
			if added {
				return []string{fact}
			}
			return nil
		}
	}
	return nil
}

func summarizeFacts(facts []string) string {
	if len(facts) == 0 {
		return ""
	}
	if len(facts) <= 3 {
		return strings.Join(facts, "; ")
	}
	return strings.Join(facts[len(facts)-3:], "; ")
}

func normalizeUserID(userID string) string {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return "demo-user"
	}
	return userID
}
