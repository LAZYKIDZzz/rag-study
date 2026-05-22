# Data Model

The first implementation uses in-memory repositories so the full RAG flow can run without infrastructure. PostgreSQL + pgvector is the intended persistent model.

## Entities

### Document

* `id`
* `title`
* `content`
* `metadata`
* `created_at`
* `indexed`

### Chunk

* `id`
* `document_id`
* `ordinal`
* `content`
* `metadata`
* `embedding`

### Chat Session

* `id`
* `user_id`
* `created_at`
* `messages`

### Message

* `role`
* `content`
* `created_at`
* `citations`

### User Memory

* `user_id`
* `facts`
* `recent_summary`

## PostgreSQL Direction

Documents, chunks, chat sessions, messages, and user memory should live in regular tables. Chunk embeddings should be stored with `pgvector`, allowing a single database to support metadata filtering and vector search.
