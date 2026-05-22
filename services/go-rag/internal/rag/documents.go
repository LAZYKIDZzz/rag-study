package rag

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

var ErrDocumentNotFound = errors.New("document not found")

type documentStore struct {
	mu        sync.RWMutex
	ids       idGenerator
	documents map[string]Document
	chunkCounts map[string]int
}

func newDocumentStore() *documentStore {
	return &documentStore{
		documents:  make(map[string]Document),
		chunkCounts: make(map[string]int),
	}
}

func (s *documentStore) create(title string, content string, metadata map[string]string) (Document, error) {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)
	if title == "" {
		return Document{}, newValidationError("title is required")
	}
	if content == "" {
		return Document{}, newValidationError("content is required")
	}

	doc := Document{
		ID:        s.ids.next("doc"),
		Title:     title,
		Content:   content,
		Metadata:  cloneStringMap(metadata),
		CreatedAt: time.Now().UTC(),
		Indexed:   false,
	}

	s.mu.Lock()
	s.documents[doc.ID] = doc
	s.mu.Unlock()

	return doc, nil
}

func (s *documentStore) list() []Document {
	s.mu.RLock()
	defer s.mu.RUnlock()

	documents := make([]Document, 0, len(s.documents))
	for _, doc := range s.documents {
		doc.Metadata = cloneStringMap(doc.Metadata)
		doc.ChunkCount = s.chunkCounts[doc.ID]
		documents = append(documents, doc)
	}

	sort.Slice(documents, func(i int, j int) bool {
		return documents[i].CreatedAt.Before(documents[j].CreatedAt)
	})
	return documents
}

func (s *documentStore) get(id string) (Document, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	doc, ok := s.documents[id]
	if !ok {
		return Document{}, ErrDocumentNotFound
	}
	doc.Metadata = cloneStringMap(doc.Metadata)
	doc.ChunkCount = s.chunkCounts[doc.ID]
	return doc, nil
}

func (s *documentStore) markIndexed(id string, chunkCount int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	doc, ok := s.documents[id]
	if !ok {
		return ErrDocumentNotFound
	}
	doc.Indexed = true
	doc.ChunkCount = chunkCount
	s.documents[id] = doc
	s.chunkCounts[id] = chunkCount
	return nil
}
