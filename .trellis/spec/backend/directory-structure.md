# Directory Structure

> How backend code is organized in this project.

---

## Overview

Backend code lives under `services/<language>-rag`. Each backend implements
the shared RAG API contract documented in `docs/architecture/api-contract.md`.
Implementations may use different frameworks, but they must preserve the same
conceptual RAG modules so Java, Python, and Go remain comparable.

---

## Directory Layout

```text
services/
  python-rag/
    app/main.py              # FastAPI routes and exception mapping
    app/schemas.py           # HTTP DTOs
    app/rag/                 # RAG pipeline modules
    tests/                   # API and pipeline tests
  java-rag/
    src/main/java/.../api/   # Spring controllers and exception handling
    src/main/java/.../core/  # RAG services, repositories, DTOs, model
    src/test/                # Java tests
  go-rag/
    cmd/server/              # HTTP entrypoint
    internal/rag/            # RAG service, stores, HTTP handlers, tests
```

---

## Module Organization

Keep RAG stages explicit:

* document ingestion and repositories,
* preprocessing,
* chunking,
* embedding provider,
* vector store,
* query processing,
* memory store,
* answer generation,
* API layer and error mapping.

Do not hide the whole pipeline behind one framework call. The repository is a
study project, so readability and stage-by-stage comparison matter.

---

## Naming Conventions

* Keep endpoint names and JSON fields aligned with the shared API contract.
* Use snake_case JSON fields across languages.
* Python modules use `snake_case.py`.
* Java packages use `study.rag.<layer>` and DTO records for API shapes.
* Go keeps implementation under `internal/rag` with small focused files.

---

## Examples

* `services/python-rag/app/rag/service.py`
* `services/java-rag/src/main/java/study/rag/core/RagService.java`
* `services/go-rag/internal/rag/service.go`
