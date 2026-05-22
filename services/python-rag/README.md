# Python RAG Service

FastAPI reference implementation for the shared RAG-study API contract.

## Run

```powershell
cd services/python-rag
python -m venv .venv
.\.venv\Scripts\Activate.ps1
pip install -r requirements.txt
uvicorn app.main:app --reload --port 8000
```

The service uses in-memory stores and deterministic local embeddings, so it does
not require PostgreSQL, pgvector, or model provider keys for the first milestone.

## Verify

```powershell
python -m compileall app
```

OpenAPI docs are available at `http://localhost:8000/docs` while the service is
running.

