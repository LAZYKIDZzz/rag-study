# Product Scope

## Purpose

`RAG-study` teaches retrieval-augmented generation by building a complete knowledge base system and comparing how Java, Python, and Go backends implement the same workflow.

## MVP Goals

* Run a complete RAG flow through the Python backend.
* Expose the same conceptual API contract for Python, Java, and Go services.
* Provide a React workbench for document ingestion, retrieval inspection, chat, and memory inspection.
* Document RAG concepts alongside code and architecture.

## Core User Workflows

1. Add a source document.
2. Index the document into chunks and embeddings.
3. Ask a question.
4. Inspect retrieved chunks and scores.
5. Continue a chat session with memory.
6. Compare how another backend implements the same flow.

## Out of Scope for the First Milestone

* Production authentication and authorization.
* Cloud deployment.
* Multi-tenant isolation.
* Automated RAG evaluation.
* Advanced graph-based retrieval.

## Architecture Links

* [System overview](../architecture/system-overview.md)
* [API contract](../architecture/api-contract.md)
* [Data model](../architecture/data-model.md)

## Learning Links

* [RAG pipeline](../learning/rag-pipeline.md)
* [Embeddings and vectors](../learning/embeddings-and-vectors.md)
* [Query rewriting and memory](../learning/query-rewriting-and-memory.md)
