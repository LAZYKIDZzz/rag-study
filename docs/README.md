# RAG-study Documentation

This documentation is split by purpose and cross-linked so product intent, code design, and RAG learning notes stay connected.

## Documentation Map

| Area | Purpose | Start Here |
| --- | --- | --- |
| Requirements | Product goals, milestones, and acceptance criteria | [requirements/product-scope.md](requirements/product-scope.md) |
| Architecture | Code structure, API contracts, data model, backend comparison | [architecture/system-overview.md](architecture/system-overview.md) |
| Learning | RAG concepts explained through this codebase | [learning/rag-pipeline.md](learning/rag-pipeline.md) |
| Decisions | ADR-style records for technical choices | [decisions/0001-python-first.md](decisions/0001-python-first.md) |
| Runbooks | Local setup and troubleshooting | [runbooks/local-development.md](runbooks/local-development.md) |

## Cross-Link Rule

* Requirement docs link to the architecture that satisfies them.
* Architecture docs link to concrete code paths and learning notes.
* Learning docs link back to the implementation and API behavior that demonstrates the concept.
