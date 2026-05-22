package rag

import (
	"errors"
	"strings"
	"time"
)

type Service struct {
	documents  *documentStore
	chunker    chunker
	embedder   embeddingProvider
	vectors    *vectorStore
	query      queryProcessor
	memory     *memoryStore
	chats      *chatStore
	generator  answerGenerator
}

func NewService() *Service {
	return &Service{
		documents: newDocumentStore(),
		chunker:   newChunker(),
		embedder:  newEmbeddingProvider(),
		vectors:   newVectorStore(),
		query:     newQueryProcessor(),
		memory:    newMemoryStore(),
		chats:     newChatStore(),
		generator: newAnswerGenerator(),
	}
}

func (s *Service) Health() map[string]string {
	return map[string]string{
		"status":  "ok",
		"service": ServiceName,
		"version": ServiceVersion,
	}
}

func (s *Service) Capabilities() map[string]any {
	return map[string]any{
		"service":            ServiceName,
		"version":            ServiceVersion,
		"embedding_provider": s.embedder.Name(),
		"vector_store":       "in-memory-cosine",
		"memory_strategy":    "in-memory-user-facts",
		"features": []string{
			"document_import",
			"preprocessing",
			"chunking",
			"deterministic_embedding",
			"vector_search",
			"query_rewriting",
			"chat_memory",
			"citation_metadata",
		},
	}
}

func (s *Service) CreateDocument(req CreateDocumentRequest) (Document, error) {
	return s.documents.create(req.Title, req.Content, req.Metadata)
}

func (s *Service) ListDocuments() []Document {
	return s.documents.list()
}

func (s *Service) IndexDocument(documentID string) (map[string]any, error) {
	doc, err := s.documents.get(documentID)
	if err != nil {
		return nil, err
	}

	chunks := s.chunker.split(doc)
	for i := range chunks {
		chunks[i].Embedding = s.embedder.Embed(chunks[i].Content)
	}

	s.vectors.replaceDocumentChunks(documentID, chunks)
	if err := s.documents.markIndexed(documentID, len(chunks)); err != nil {
		return nil, err
	}

	return map[string]any{
		"document_id":  documentID,
		"chunk_count":  len(chunks),
		"vector_count": len(chunks),
	}, nil
}

func (s *Service) Search(req SearchRequest) (SearchResponse, error) {
	query := strings.TrimSpace(req.Query)
	if query == "" {
		return SearchResponse{}, newValidationError("query is required")
	}

	rewritten := s.query.rewrite(query, UserMemory{})
	chunks := s.vectors.search(s.embedder.Embed(rewritten), req.TopK, req.Filters)
	return SearchResponse{
		Query:           query,
		RewrittenQuery:  rewritten,
		RetrievedChunks: chunks,
	}, nil
}

func (s *Service) CreateChatSession(req CreateChatSessionRequest) ChatSession {
	return s.chats.create(req.UserID)
}

func (s *Service) GetChatSession(sessionID string) (ChatSession, error) {
	return s.chats.get(sessionID)
}

func (s *Service) AddChatMessage(sessionID string, req ChatMessageRequest) (ChatMessageResponse, error) {
	message := strings.TrimSpace(req.Message)
	if message == "" {
		return ChatMessageResponse{}, newValidationError("message is required")
	}

	session, err := s.chats.get(sessionID)
	if err != nil {
		return ChatMessageResponse{}, err
	}

	userID := normalizeUserID(req.UserID)
	if req.UserID == "" && session.UserID != "" {
		userID = session.UserID
	}

	memoryBefore := s.memory.get(userID)
	rewritten := s.query.rewrite(message, memoryBefore)
	retrieved := s.vectors.search(s.embedder.Embed(rewritten), req.TopK, nil)
	answer := s.generator.answer(message, retrieved, memoryBefore)
	memoryUpdates := s.memory.updateFromMessage(userID, message)

	now := time.Now().UTC()
	if _, err := s.chats.appendMessage(sessionID, ChatMessage{
		Role:      "user",
		Content:   message,
		CreatedAt: now,
	}); err != nil {
		return ChatMessageResponse{}, err
	}
	if _, err := s.chats.appendMessage(sessionID, ChatMessage{
		Role:      "assistant",
		Content:   answer,
		CreatedAt: now,
		Citations: retrieved,
	}); err != nil {
		return ChatMessageResponse{}, err
	}

	return ChatMessageResponse{
		SessionID:       sessionID,
		Answer:          answer,
		Citations:       retrieved,
		RewrittenQuery:  rewritten,
		RetrievedChunks: retrieved,
		MemoryUpdates:   memoryUpdates,
	}, nil
}

func (s *Service) GetMemory(userID string) UserMemory {
	return s.memory.get(userID)
}

func (s *Service) AddFact(userID string, req AddFactRequest) (UserMemory, error) {
	if strings.TrimSpace(req.Fact) == "" {
		return UserMemory{}, newValidationError("fact is required")
	}
	record, _ := s.memory.addFact(userID, req.Fact)
	return record, nil
}

func errorCode(err error) string {
	var validation validationError
	switch {
	case errors.Is(err, ErrDocumentNotFound):
		return "document_not_found"
	case errors.Is(err, ErrSessionNotFound):
		return "session_not_found"
	case errors.As(err, &validation):
		return "validation_error"
	default:
		return "internal_error"
	}
}
