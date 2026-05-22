# Java RAG Service

`services/java-rag` is the Java implementation of the shared RAG-study backend contract. It uses Spring Boot style controllers and services, with explicit RAG modules for learning and comparison.

## Capabilities

* Document ingestion with metadata.
* Text preprocessing and overlapping word chunking.
* Deterministic local embeddings for offline development.
* In-memory cosine vector search.
* Query rewriting with user memory facts and previous chat turns.
* Chat sessions with citations and memory updates.
* Shared API routes from `docs/architecture/api-contract.md`.

## Run

Prerequisites:

* Java 17 or newer.
* Maven 3.9 or newer.

```bash
cd services/java-rag
mvn spring-boot:run
```

The service listens on `http://localhost:8082`.

## Verify

```bash
mvn test
```

Maven is not installed in the current workspace, so this implementation pass could not execute the Java build locally.
