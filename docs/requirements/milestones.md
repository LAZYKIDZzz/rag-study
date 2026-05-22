# Milestones

## Milestone 1: Python End-to-End Baseline

* Create repository structure and documentation map.
* Implement Python FastAPI RAG backend with in-memory vector search and deterministic local embeddings.
* Build React workbench that can ingest text, ask questions, and inspect retrieval traces.
* Add infrastructure files for PostgreSQL + pgvector.

## Milestone 2: Persistent Storage

* Persist documents, chunks, sessions, and memory to PostgreSQL.
* Add pgvector-backed search.
* Keep the in-memory store for tests and offline demos.

## Milestone 3: Java and Go Parity

* Implement Java and Go services against the shared API contract.
* Document framework differences and code-level trade-offs.
* Add backend switching in the frontend.

## Milestone 4: Advanced RAG

* Query rewriting strategies.
* Hybrid retrieval.
* Reranking.
* User memory summarization.
* Evaluation datasets and regression checks.
