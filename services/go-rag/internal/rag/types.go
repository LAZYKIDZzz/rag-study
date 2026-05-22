package rag

import "time"

const (
	ServiceName    = "go-rag"
	ServiceVersion = "0.1.0"
)

type Document struct {
	ID        string            `json:"id"`
	Title     string            `json:"title"`
	Content   string            `json:"content,omitempty"`
	Metadata  map[string]string `json:"metadata"`
	CreatedAt time.Time         `json:"created_at"`
	Indexed   bool              `json:"indexed"`
	ChunkCount int               `json:"chunk_count"`
}

type Chunk struct {
	ID         string            `json:"id"`
	DocumentID string            `json:"document_id"`
	Ordinal    int               `json:"ordinal"`
	Content    string            `json:"content"`
	Metadata   map[string]string `json:"metadata"`
	Embedding  []float64         `json:"-"`
}

type RetrievedChunk struct {
	ChunkID       string            `json:"chunk_id"`
	DocumentID    string            `json:"document_id"`
	DocumentTitle string            `json:"document_title"`
	Ordinal       int               `json:"ordinal"`
	Content       string            `json:"content"`
	Metadata      map[string]string `json:"metadata"`
	Score         float64           `json:"score"`
}

type ChatSession struct {
	ID        string        `json:"id"`
	UserID    string        `json:"user_id,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	Messages  []ChatMessage `json:"messages"`
}

type ChatMessage struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
	Citations []RetrievedChunk `json:"citations,omitempty"`
}

type UserMemory struct {
	UserID        string   `json:"user_id"`
	Facts         []string `json:"facts"`
	RecentSummary string   `json:"recent_summary"`
}

type CreateDocumentRequest struct {
	Title    string            `json:"title"`
	Content  string            `json:"content"`
	Metadata map[string]string `json:"metadata"`
}

type SearchRequest struct {
	Query   string            `json:"query"`
	TopK    int               `json:"top_k"`
	Filters map[string]string `json:"filters"`
}

type SearchResponse struct {
	Query           string           `json:"query"`
	RewrittenQuery  string           `json:"rewritten_query"`
	RetrievedChunks []RetrievedChunk `json:"retrieved_chunks"`
}

type CreateChatSessionRequest struct {
	UserID string `json:"user_id"`
}

type ChatMessageRequest struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
	TopK    int    `json:"top_k"`
}

type ChatMessageResponse struct {
	SessionID       string           `json:"session_id"`
	Answer          string           `json:"answer"`
	Citations       []RetrievedChunk `json:"citations"`
	RewrittenQuery  string           `json:"rewritten_query"`
	RetrievedChunks []RetrievedChunk `json:"retrieved_chunks"`
	MemoryUpdates   []string         `json:"memory_updates"`
}

type AddFactRequest struct {
	Fact string `json:"fact"`
}

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details"`
}
