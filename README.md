# RAG-study

`RAG-study` is a learning-first RAG knowledge base project. It is designed to be both runnable software and a comparison lab for how Java, Python, and Go backends support the same retrieval-augmented generation workflow.

## What This Project Builds

* A React workbench for uploading documents, inspecting retrieval, chatting with a knowledge base, and comparing backend implementations.
* A complete Python RAG backend built with FastAPI.
* Java and Go backend implementations that follow the same API contract and RAG module boundaries.
* Local infrastructure for PostgreSQL with pgvector, with room to add Qdrant as a dedicated vector database comparison.
* Documentation for requirements, architecture, and RAG learning notes, all cross-linked.

## Repository Layout

```text
apps/
  web/                  React RAG workbench
services/
  python-rag/           FastAPI implementation and first complete backend
  java-rag/             Java/Spring-style implementation
  go-rag/               Go HTTP implementation
docs/
  requirements/         Product goals, scope, milestones
  architecture/         System design, API contract, data model
  learning/             RAG concept notes linked to code
  decisions/            ADR-style technical decisions
  runbooks/             Local development and troubleshooting
infra/
  docker-compose.yml    Local databases and middleware
```

## Quick Start

Python backend:

```bash
cd services/python-rag
python -m venv .venv
.\.venv\Scripts\Activate.ps1
pip install -r requirements.txt
uvicorn app.main:app --reload --port 8000
```

Frontend:

```bash
cd apps/web
npm install
npm run dev
```

Local infrastructure:

```bash
docker compose -f infra/docker-compose.yml up -d
```

## Reading Order

Start with [docs/README.md](docs/README.md), then read:

1. [Product scope](docs/requirements/product-scope.md)
2. [System overview](docs/architecture/system-overview.md)
3. [RAG pipeline learning guide](docs/learning/rag-pipeline.md)
4. [API contract](docs/architecture/api-contract.md)

## Current Implementation Strategy

The first milestone uses Python as the complete backend so the full RAG pipeline can be inspected quickly. Java and Go follow the same API shape and module boundaries so differences in framework support are easier to compare.
