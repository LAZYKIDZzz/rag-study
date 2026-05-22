# Stack Selection Research

## Scope

This research supports the first planning pass for `RAG-study`: framework choices, vector storage choices, middleware, frontend tooling, and documentation structure for a multi-backend RAG learning project.

## Primary Sources Checked

* Spring AI RAG and advisors documentation: https://docs.spring.io/spring-ai/reference/api/retrieval-augmented-generation.html
* Spring AI chat memory documentation: https://docs.spring.io/spring-ai/reference/api/chat-memory.html
* Spring AI advisors documentation: https://docs.spring.io/spring-ai/reference/api/advisors.html
* LangChain4j RAG documentation: https://docs.langchain4j.dev/tutorials/rag
* LangChain4j embedding stores documentation: https://docs.langchain4j.dev/tutorials/embedding-stores/
* LangChain Python retrieval documentation: https://docs.langchain.com/oss/python/langchain/retrieval
* LangChain Python vector stores documentation: https://docs.langchain.com/oss/python/integrations/vectorstores/
* LangChain Python text splitters documentation: https://docs.langchain.com/oss/python/integrations/splitters/index
* LlamaIndex ingestion pipeline documentation: https://docs.llamaindex.ai/en/stable/module_guides/loading/ingestion_pipeline/
* LlamaIndex memory documentation: https://docs.llamaindex.ai/en/stable/module_guides/deploying/agents/memory/
* FastAPI first steps documentation: https://fastapi.tiangolo.com/tutorial/first-steps/
* Eino overview and core modules: https://www.cloudwego.io/docs/eino/overview/ and https://www.cloudwego.io/docs/eino/core_modules/
* Eino embedding and retriever guides: https://www.cloudwego.io/docs/eino/core_modules/components/embedding_guide/ and https://www.cloudwego.io/docs/eino/core_modules/components/retriever_guide/
* Qdrant official interfaces: https://qdrant.tech/documentation/interfaces/
* Qdrant documentation: https://qdrant.tech/documentation/
* pgvector project documentation: https://github.com/pgvector/pgvector
* pgvector Go support: https://github.com/pgvector/pgvector-go
* Vite guide: https://vite.dev/guide/
* TanStack Query documentation: https://tanstack.com/query/v4/docs
* React Router documentation: https://reactrouter.com/
* shadcn/ui Vite installation documentation: https://ui.shadcn.com/docs/installation/vite

## Recommended Baseline Stack

### Frontend

Recommended:

* React + TypeScript + Vite
* TanStack Query for server-state fetching/caching
* React Router for local navigation
* shadcn/ui + Tailwind CSS for a practical UI component baseline

Why:

* Vite provides a standard React TypeScript template and fast local development.
* TanStack Query keeps API state explicit, which is useful when switching among Java/Python/Go backends.
* React Router is enough for a study app with pages such as documents, chat, retrieval trace, memory, and backend comparison.
* shadcn/ui gives editable components rather than a black-box component dependency, which fits a learning project.

### Java Backend

Recommended:

* Spring Boot + Spring AI as the primary Java implementation.
* Keep LangChain4j as an optional comparison track or later alternate service.

Why:

* Spring AI has first-class RAG primitives through advisors and vector-store integration.
* Spring AI has chat memory abstractions and repositories, including JDBC-backed memory.
* Spring Boot is widely used in Java backends, so this implementation demonstrates how enterprise Java can host a RAG system.
* LangChain4j is also strong for Java RAG and vector store integrations, but using both immediately may dilute the first milestone.

### Python Backend

Recommended:

* FastAPI for HTTP service structure.
* LangChain + LangGraph for the first implementation path.
* LlamaIndex as a comparison/reference for ingestion pipelines and memory concepts.

Why:

* FastAPI is a straightforward API framework with async support and automatic OpenAPI output.
* LangChain has well-documented retrieval, vector store, and text splitter abstractions.
* LangGraph is a good future fit for query rewriting, memory, and multi-step RAG orchestration.
* LlamaIndex is particularly strong for ingestion pipeline concepts, transformations, caching, and document-centric workflows.

### Go Backend

Recommended:

* Go HTTP service using Gin or Chi for the API layer.
* Eino for AI component abstractions where useful.
* Direct pgvector or Qdrant client integration for vector storage.

Why:

* Eino is a Go-native LLM application framework with components such as `ChatModel`, `Embedding`, and `Retriever`.
* Go RAG frameworks are less standardized than Python, so the Go implementation should expose more of the underlying pipeline for learning.
* Direct database/client usage in Go helps demonstrate embeddings, vector insert/query, filtering, and memory storage without excessive framework magic.

### Databases and Middleware

Recommended default:

* PostgreSQL + pgvector for relational data, metadata, memory, and baseline vector search.
* Redis for optional cache/rate-limit/job state later.
* Object/file storage as local filesystem first, optionally MinIO later for document originals.
* Docker Compose for local development.

Recommended comparison vector store:

* Qdrant, added once the pgvector baseline is working.

Why:

* pgvector keeps vectors with normal relational data and supports approximate nearest neighbor indexes such as HNSW/IVFFlat.
* Qdrant is a dedicated vector search engine with official clients, including Go, and is useful for comparing payload filtering, hybrid search, and vector-database operations.
* Starting with both as required infrastructure may slow the first milestone; a default plus optional comparison path is cleaner.

## Shared RAG Capability Contract

Each backend should implement the same conceptual modules:

* `DocumentIngestion`: accept files/text, normalize metadata, persist source records.
* `Preprocessor`: extract text, clean text, detect format, preserve provenance.
* `Chunker`: split content into chunks with stable IDs and overlap metadata.
* `EmbeddingProvider`: generate embeddings through an abstract provider interface.
* `VectorRepository`: upsert/search/delete vectorized chunks.
* `Retriever`: combine query embedding, filters, top-k search, and optional reranking.
* `QueryProcessor`: rewrite or expand user questions and attach conversation context.
* `MemoryStore`: store recent chat messages, summaries, and longer-term user facts.
* `AnswerGenerator`: build prompt context, call chat model, return answer with citations.
* `TraceRecorder`: record retrieved chunks, scores, prompt metadata, and model choices.

## Suggested API Shape

The frontend should be able to call equivalent endpoints for any backend:

* `GET /health`
* `GET /capabilities`
* `POST /documents`
* `GET /documents`
* `POST /documents/{documentId}/index`
* `POST /chat/sessions`
* `POST /chat/sessions/{sessionId}/messages`
* `GET /chat/sessions/{sessionId}`
* `POST /retrieval/search`
* `GET /memory/users/{userId}`
* `POST /memory/users/{userId}/facts`

## Documentation Architecture

Recommended root docs layout:

* `docs/README.md` - documentation map and reading order.
* `docs/requirements/` - product goals, milestones, user stories, acceptance criteria.
* `docs/architecture/` - system architecture, API contracts, data model, deployment, backend comparisons.
* `docs/learning/` - RAG concepts: embeddings, vectors, chunking, retrieval, query rewriting, memory, evaluation.
* `docs/decisions/` - ADR-style records for stack and architecture decisions.
* `docs/runbooks/` - local startup, troubleshooting, model/provider configuration.

Cross-link rule:

* Requirement docs link to the architecture docs that satisfy them.
* Architecture docs link to code paths and learning notes.
* Learning docs link to one or more concrete implementations in Java/Python/Go.

## Feasible MVP Approaches

### Approach A: Contract-first, one complete backend first

Build shared docs, API contract, local infrastructure, frontend shell, and one complete backend first. Add the other two backends after the baseline is proven.

Pros:

* Fastest path to a working end-to-end RAG system.
* Reduces the chance of duplicating a flawed design three times.
* Lets the project learn from the first implementation before cloning the contract.

Cons:

* Cross-language comparison is delayed.

### Approach B: Three backend skeletons first

Create React app, shared contract, and Java/Python/Go services with matching endpoint skeletons, then fill RAG features in parallel.

Pros:

* Establishes language comparison structure immediately.
* Makes shared API design visible early.

Cons:

* More scaffolding before users can run a meaningful RAG flow.
* Higher risk of shallow implementations.

### Approach C: Documentation-first study repository

Create the complete docs map, architecture decisions, learning notes, and diagrams before building runnable services.

Pros:

* Strongest learning structure.
* Makes later implementation easier to reason about.

Cons:

* Delays executable feedback and may over-design before code reveals constraints.

## Recommendation

Use Approach A with enough skeleton for all three services to reserve the comparison shape:

1. Create docs, infrastructure, shared API contract, frontend shell, and one complete backend.
2. Keep placeholder service directories and README files for the other two backends.
3. Implement the second and third backends against the already-tested contract.

This balances the user's three goals: complete system, RAG learning, and backend-language comparison.
