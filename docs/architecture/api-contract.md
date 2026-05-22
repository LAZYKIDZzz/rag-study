# Shared API Contract

All backend implementations should expose equivalent endpoints and response shapes. The frontend uses this contract to switch backend targets.

## Health

`GET /health`

```json
{
  "status": "ok",
  "service": "python-rag",
  "version": "0.1.0"
}
```

## Capabilities

`GET /capabilities`

Returns supported features, embedding provider, vector store, and memory strategy.

```json
{
  "service": "python-rag",
  "features": ["document-ingestion", "chunking", "chat-memory"],
  "embedding_provider": "local-hash-embedding-128d",
  "vector_store": "in-memory-cosine",
  "memory_strategy": "in-memory"
}
```

## Documents

`POST /documents`

```json
{
  "title": "RAG notes",
  "content": "Raw document text",
  "metadata": {
    "source": "manual"
  }
}
```

`GET /documents`

Returns all known documents and indexing status.

```json
{
  "documents": [
    {
      "id": "doc_123",
      "title": "RAG notes",
      "metadata": {
        "source": "manual"
      },
      "indexed": true,
      "chunk_count": 3,
      "created_at": "2026-05-22T00:00:00Z"
    }
  ]
}
```

`POST /documents/{document_id}/index`

Chunks and embeds a document.

## Retrieval

`POST /retrieval/search`

```json
{
  "query": "What is chunking?",
  "top_k": 5,
  "filters": {}
}
```

Returns retrieved chunks, similarity scores, and source metadata.

```json
{
  "query": "What is chunking?",
  "rewritten_query": "Current question: What is chunking?",
  "retrieved_chunks": [
    {
      "chunk_id": "chk_123",
      "document_id": "doc_123",
      "document_title": "RAG notes",
      "ordinal": 0,
      "content": "Chunking splits documents into searchable passages.",
      "score": 0.82,
      "metadata": {
        "source": "manual"
      }
    }
  ]
}
```

## Chat

`POST /chat/sessions`

Creates a chat session.

`POST /chat/sessions/{session_id}/messages`

```json
{
  "user_id": "demo-user",
  "message": "Explain embeddings",
  "top_k": 5
}
```

Returns an answer, citations, rewritten query, retrieved chunks, and memory updates.

```json
{
  "session_id": "ses_123",
  "answer": "Based on the retrieved knowledge...",
  "rewritten_query": "Current question: Explain embeddings",
  "citations": [],
  "retrieved_chunks": [],
  "memory_updates": ["Last asked: Explain embeddings"]
}
```

`GET /chat/sessions/{session_id}`

Returns messages for a session.

## Memory

`GET /memory/users/{user_id}`

Returns stored facts and recent chat summary for a user.

`POST /memory/users/{user_id}/facts`

Adds a user memory fact.

## Error Shape

```json
{
  "error": {
    "code": "document_not_found",
    "message": "Document was not found",
    "details": {}
  }
}
```
