package rag

import "testing"

func TestIndexSearchAndChatFlow(t *testing.T) {
	service := NewService()

	doc, err := service.CreateDocument(CreateDocumentRequest{
		Title:   "RAG notes",
		Content: "Chunking splits documents into smaller passages. Embeddings turn text into vectors for retrieval.",
		Metadata: map[string]string{
			"source": "test",
		},
	})
	if err != nil {
		t.Fatalf("CreateDocument returned error: %v", err)
	}

	indexed, err := service.IndexDocument(doc.ID)
	if err != nil {
		t.Fatalf("IndexDocument returned error: %v", err)
	}
	if indexed["chunk_count"].(int) == 0 {
		t.Fatal("expected at least one chunk")
	}

	search, err := service.Search(SearchRequest{Query: "What are embeddings?", TopK: 2})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(search.RetrievedChunks) == 0 {
		t.Fatal("expected retrieved chunks")
	}

	session := service.CreateChatSession(CreateChatSessionRequest{UserID: "student"})
	response, err := service.AddChatMessage(session.ID, ChatMessageRequest{
		UserID:  "student",
		Message: "Explain embeddings",
		TopK:    2,
	})
	if err != nil {
		t.Fatalf("AddChatMessage returned error: %v", err)
	}
	if response.Answer == "" {
		t.Fatal("expected answer")
	}
	if len(response.Citations) == 0 {
		t.Fatal("expected citations")
	}
}

func TestMemoryFactDeduplicates(t *testing.T) {
	store := newMemoryStore()

	_, added := store.addFact("u1", "prefers short answers")
	if !added {
		t.Fatal("expected first fact to be added")
	}

	memory, added := store.addFact("u1", "Prefers short answers")
	if added {
		t.Fatal("expected duplicate fact to be ignored")
	}
	if len(memory.Facts) != 1 {
		t.Fatalf("expected one fact, got %d", len(memory.Facts))
	}
}
