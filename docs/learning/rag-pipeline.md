# RAG Pipeline

Retrieval-augmented generation connects a language model to external knowledge. This project keeps every stage visible so each backend can be compared.

## Stages

1. Preprocess source documents.
2. Split text into chunks.
3. Generate embeddings for each chunk.
4. Store vectors and metadata.
5. Rewrite or enrich a user query.
6. Embed the query.
7. Retrieve relevant chunks.
8. Build an answer prompt with citations.
9. Generate an answer.
10. Update chat and user memory.

## Code Links

* Python pipeline: [services/python-rag/app/rag/service.py](../../services/python-rag/app/rag/service.py)
* Chunking: [services/python-rag/app/rag/chunking.py](../../services/python-rag/app/rag/chunking.py)
* Embeddings: [services/python-rag/app/rag/embeddings.py](../../services/python-rag/app/rag/embeddings.py)
* Retrieval: [services/python-rag/app/rag/vector_store.py](../../services/python-rag/app/rag/vector_store.py)

## Requirement Links

* [Product scope](../requirements/product-scope.md)
* [API contract](../architecture/api-contract.md)
