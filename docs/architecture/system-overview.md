# System Overview

## High-Level Shape

```text
React Workbench
  -> Shared HTTP API Contract
    -> Python FastAPI RAG service
    -> Java RAG service
    -> Go RAG service
      -> Document store
      -> Chunker
      -> Embedding provider
      -> Vector repository
      -> Query processor
      -> Memory store
      -> Answer generator
```

The first complete implementation is [services/python-rag](../../services/python-rag). Java and Go services follow the same capability contract so their framework choices can be compared.

## Data Flow

```text
Source document
  -> preprocessing
  -> chunking
  -> embedding
  -> vector storage
  -> retrieval
  -> answer prompt construction
  -> model response
  -> memory update
  -> frontend trace display
```

This flow is explained conceptually in [RAG pipeline](../learning/rag-pipeline.md).

## Design Principles

* Keep API shapes comparable across backend languages.
* Keep RAG stages explicit instead of hiding everything behind one framework call.
* Default to local/offline behavior where possible so the project can be studied without API keys.
* Use provider interfaces for embeddings and generation so real model providers can be added later.

## Code Links

* Python service: [services/python-rag](../../services/python-rag)
* Java service: [services/java-rag](../../services/java-rag)
* Go service: [services/go-rag](../../services/go-rag)
* Frontend: [apps/web](../../apps/web)
