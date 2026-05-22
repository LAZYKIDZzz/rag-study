# RAG-study Multi-Backend Knowledge Base

## Goal

Build `RAG-study` as a study-oriented, full-stack RAG knowledge base system. The project must be useful as runnable software and as a learning repository: one React frontend, multiple backend implementations in at least Java, Python, and Go, and a documentation system that explains requirements, code architecture, and RAG concepts with clear cross-links.

## What I Already Know

* Project name: `RAG-study`.
* Main goals:
  * Build a complete RAG knowledge base system with frontend and backend.
  * Learn the core RAG workflow: preprocessing, chunking, embeddings, vectors, retrieval, query rewriting, chat/user memory, and answer generation.
  * Compare how different backend languages/frameworks support RAG.
* Frontend: React. It only needs enough depth to understand and operate the RAG workflows.
* Backends: at least Java, Python, and Go. Each backend should implement the same core RAG capabilities so they can be compared.
* Documentation must be substantial and separated into:
  * requirement/product analysis documentation,
  * code/architecture documentation,
  * RAG learning documentation.
* Documentation categories must link to each other clearly.
* Git commits must be made as `LAZYKIDZzz`; no push is required.
* Code should prioritize readability, maintainability, extensibility, and normal coding standards.

## Requirements

* Provide a monorepo-style project structure with:
  * `apps/web` for the React frontend,
  * `services/java-rag` for the Java backend,
  * `services/python-rag` for the Python backend,
  * `services/go-rag` for the Go backend,
  * shared documentation under `docs/`,
  * local infrastructure configuration for databases and middleware.
* Provide a consistent RAG capability contract across all backends:
  * document upload/import,
  * data preprocessing,
  * chunking,
  * embedding generation,
  * vector storage,
  * retrieval,
  * query processing/query rewriting,
  * answer generation,
  * conversation and user memory management,
  * observability-friendly logs and metadata.
* Prefer shared API concepts and data schemas across backend implementations so the frontend can switch target backend with minimal differences.
* Document each RAG stage with:
  * conceptual explanation,
  * implementation notes,
  * language-specific comparison,
  * links to related code and requirements.
* Provide a development environment that can run locally with Docker Compose where practical.
* Use multiple focused git commits, with each commit describing one coherent change.

## Acceptance Criteria

* [ ] Repository has a documented monorepo layout for frontend, Java backend, Python backend, Go backend, docs, and infrastructure.
* [ ] Root README explains project purpose, quick start, backend comparison goal, and documentation map.
* [ ] `docs/requirements/` contains product and scope documents.
* [ ] `docs/architecture/` contains code/system architecture documents.
* [ ] `docs/learning/` contains RAG concept learning notes.
* [ ] The documentation categories link to each other through an index or map.
* [ ] React frontend can call at least one backend for the basic RAG flow.
* [ ] Java, Python, and Go backends each expose comparable RAG API endpoints.
* [ ] Each backend includes data preprocessing, chunking, embedding abstraction, vector storage abstraction, retrieval, query processing, and memory management.
* [ ] Local infrastructure supports a relational store and a vector search option.
* [ ] Commits are authored as `LAZYKIDZzz` and are grouped by purpose.

## Definition of Done

* Tests added or updated where implementation risk justifies it.
* Lint/typecheck/build commands documented and passing for implemented packages.
* Documentation updated alongside behavior or architecture changes.
* Local startup path documented.
* Implementation choices and trade-offs recorded in architecture docs.
* Commit history contains focused commits and no push is performed.

## Research References

* [`research/stack-selection.md`](research/stack-selection.md) - framework, vector database, middleware, frontend, and documentation stack findings.

## Technical Approach

The initial direction is to build a learning-first monorepo:

* Use one shared API contract and comparable endpoint names across Java, Python, and Go.
* Use PostgreSQL plus `pgvector` as the default relational/vector baseline because it keeps metadata, memory, and vectors in one familiar local database.
* Optionally add Qdrant as a second vector backend for comparing dedicated vector database behavior, filtering, and hybrid search.
* Use provider abstractions for LLM and embedding calls so OpenAI-compatible APIs, local Ollama models, or other providers can be swapped.
* Implement the simplest useful RAG path first, then layer in query rewriting, memory, hybrid search, reranking, and evaluation.
* MVP sequencing: use "one complete backend first" to validate the full RAG workflow before implementing the other backend languages.
* First complete backend: Python, using FastAPI and a deliberately readable RAG pipeline.
* Keep Java, Python, and Go service directories in the initial structure so the comparison target is visible from the beginning, even if only one backend is feature-complete in the first milestone.
* Product and technical decisions that do not require user preference should follow the recommended route in research notes rather than asking additional questions.

## Decision (ADR-lite)

**Context**: The project has two modes of success: runnable software and explicit learning/comparison across language ecosystems.

**Decision**: Start with a monorepo and shared RAG contract. Implement Python as the first complete backend, then add comparable Java and Go services against the proven contract rather than building three unrelated demos.

**Consequences**:

* Higher upfront design work, but much better backend comparison.
* More documentation is required to keep concepts, requirements, and code connected.
* Shared contracts reduce frontend duplication and make language differences easier to inspect.
* Cross-language parity is intentionally delayed until the first full RAG workflow is validated.
* Python-first makes the initial RAG implementation easier to inspect because the ecosystem is mature and examples are abundant.

## Open Questions

* None currently. Continue with recommended defaults unless a blocking product decision appears.

## Out of Scope

* Production-grade authentication and authorization in the first milestone.
* Cloud deployment in the first milestone.
* Multi-tenant security hardening in the first milestone.
* Advanced RAG features such as GraphRAG, agentic retrieval, and automated evaluation until the baseline RAG flow exists.

## Technical Notes

* Current repository is effectively empty except Trellis project scaffolding and `AGENTS.md`.
* Existing Trellis specs are placeholders for backend and frontend guidelines; they should be filled as conventions emerge.
* This feature spans many layers, so the cross-layer data-flow guide applies: Source -> Transform -> Store -> Retrieve -> Transform -> Display.
* Git currently has uncommitted scaffold files from project initialization. These should be handled carefully in a commit plan rather than silently bundled with later implementation work.
