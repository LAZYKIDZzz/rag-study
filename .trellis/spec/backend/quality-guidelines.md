# Quality Guidelines

> Code quality standards for backend development.

---

## Overview

Backend changes must preserve the shared API contract and keep RAG stages
inspectable. The first milestone supports offline execution with deterministic
local embeddings and in-memory stores; production providers can be added behind
the same interfaces later.

---

## Forbidden Patterns

* Do not introduce language-specific endpoint shapes unless the shared contract
  is updated at the same time.
* Do not require external LLM or embedding credentials for the default local
  demo path.
* Do not mix unrelated RAG stages into one large controller method.
* Do not log document contents, secrets, or user memory facts at info level.

---

## Required Patterns

* Return the shared error shape:
  `{"error":{"code":"...","message":"...","details":{}}}`.
* Use snake_case JSON fields for cross-language API responses.
* Keep document IDs, chunk IDs, scores, source metadata, and rewritten query in
  retrieval/chat responses so the frontend can show traces.
* Keep provider abstractions for embeddings, vector search, memory, and answer
  generation.

---

## Testing Requirements

* Python backend: run `python -m pytest services/python-rag/tests -q`.
* Frontend contract consumers: run `npm run build` in `apps/web`.
* Java backend: run `mvn test` when Maven is available.
* Go backend: run `go test ./...` when Go is available.
* At minimum, each backend should have a full flow test: create document,
  index, search, create chat session, send message, verify citations/memory.

---

## Code Review Checklist

* API response shapes match `docs/architecture/api-contract.md`.
* RAG stages remain independently readable and replaceable.
* Retrieval traces include `chunk_id`, `document_id`, `document_title`,
  `ordinal`, `score`, and metadata.
* Missing local tools are reported explicitly rather than hidden.

---

## Scenario: Shared RAG API Contract

### 1. Scope / Trigger

Trigger: any backend or frontend change that adds or changes RAG HTTP request
or response shapes. This is a cross-layer contract because React, Python, Java,
and Go all consume or produce the same payloads.

### 2. Signatures

Required backend endpoints:

* `GET /health`
* `GET /capabilities`
* `POST /documents`
* `GET /documents`
* `POST /documents/{document_id}/index`
* `POST /retrieval/search`
* `POST /chat/sessions`
* `GET /chat/sessions/{session_id}`
* `POST /chat/sessions/{session_id}/messages`
* `GET /memory/users/{user_id}`
* `POST /memory/users/{user_id}/facts`

### 3. Contracts

Required response fields:

* Capabilities: `service`, `features`, `embedding_provider`,
  `vector_store`, `memory_strategy`.
* Document list: `{"documents": [...]}`.
* Retrieval/chat chunks: `chunk_id`, `document_id`, `document_title`,
  `ordinal`, `content`, `score`, `metadata`.
* Chat response: `session_id`, `answer`, `rewritten_query`, `citations`,
  `retrieved_chunks`, `memory_updates`.
* Errors: `{"error":{"code":"...","message":"...","details":{}}}`.

### 4. Validation & Error Matrix

* Empty document title/content -> `validation_error`.
* Empty retrieval query -> `validation_error`.
* Empty chat message -> `validation_error`.
* Missing document ID -> `document_not_found` or `not_found`.
* Missing chat session ID -> `session_not_found` or `not_found`.

### 5. Good/Base/Bad Cases

* Good: document is created, indexed, queried, cited, and memory summary updates.
* Base: no indexed chunks returns an answer explaining that no matching content
  is available yet.
* Bad: backend returns `results` instead of `retrieved_chunks`; frontend traces
  disappear or need backend-specific branches.

### 6. Tests Required

* Python: full API flow test in `services/python-rag/tests/test_api.py`.
* Java: service-level full RAG flow test in `services/java-rag/src/test`.
* Go: service-level full RAG flow test in `services/go-rag/internal/rag`.
* Frontend: `npm run build` must pass after API type changes.

### 7. Wrong vs Correct

Wrong:

```json
{
  "results": [{"id": "chk_1", "title": "Doc", "score": 0.8}]
}
```

Correct:

```json
{
  "retrieved_chunks": [
    {
      "chunk_id": "chk_1",
      "document_id": "doc_1",
      "document_title": "Doc",
      "ordinal": 0,
      "content": "Chunk text",
      "score": 0.8,
      "metadata": {}
    }
  ]
}
```
