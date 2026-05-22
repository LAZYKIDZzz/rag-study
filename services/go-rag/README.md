# Go RAG Service

`services/go-rag` is the Go implementation of the shared RAG-study backend contract. It uses only the Go standard library and in-memory stores so the RAG stages are easy to inspect.

## Capabilities

* Document import with metadata.
* Whitespace preprocessing and overlapping word chunking.
* Deterministic hash embeddings for offline development.
* In-memory cosine vector search.
* Query rewriting with user memory facts.
* Chat sessions with citations and memory updates.
* Shared error shape from `docs/architecture/api-contract.md`.

## Run

Prerequisite: Go 1.22 or newer.

```bash
cd services/go-rag
go run ./cmd/server
```

The service listens on `http://localhost:8080` by default. Set `PORT` to override it:

```bash
$env:PORT = "8081"
go run ./cmd/server
```

## Verify

```bash
cd services/go-rag
go test ./...
```

Go is not installed in the current workspace, so these commands could not be executed during this implementation pass.

## Example Flow

```bash
curl -X POST http://localhost:8080/documents `
  -H "Content-Type: application/json" `
  -d "{\"title\":\"RAG notes\",\"content\":\"Chunking splits documents. Embeddings create vectors.\",\"metadata\":{\"source\":\"manual\"}}"

curl -X POST http://localhost:8080/documents/doc_000001/index

curl -X POST http://localhost:8080/retrieval/search `
  -H "Content-Type: application/json" `
  -d "{\"query\":\"What are embeddings?\",\"top_k\":3,\"filters\":{}}"
```
