# Embeddings and Vectors

An embedding maps text to a numeric vector. Similar text should produce vectors that are close under a similarity metric such as cosine similarity.

## Why This Project Starts with Local Embeddings

The first backend uses deterministic local embeddings so the RAG flow works without model API keys. This is not a production semantic embedding model; it is a learning and testing baseline.

Later milestones can add:

* OpenAI-compatible embeddings.
* Ollama/local model embeddings.
* pgvector persistence.
* Qdrant comparison.

## Code Links

* Local embedding provider: [services/python-rag/app/rag/embeddings.py](../../services/python-rag/app/rag/embeddings.py)
* Vector search: [services/python-rag/app/rag/vector_store.py](../../services/python-rag/app/rag/vector_store.py)
