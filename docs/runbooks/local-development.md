# Local Development

## Prerequisites

* Python 3.11 or newer.
* Node.js 20 or newer.
* Docker for PostgreSQL + pgvector.
* Java 17 for the Java service.
* Go 1.22 or newer for the Go service.

## Python Backend

```bash
cd services/python-rag
python -m venv .venv
.\.venv\Scripts\Activate.ps1
pip install -r requirements.txt
uvicorn app.main:app --reload --port 8000
```

## Frontend

```bash
cd apps/web
npm install
npm run dev
```

Set `VITE_RAG_API_BASE_URL=http://localhost:8000` if the backend is not on the default URL.

## Infrastructure

```bash
docker compose -f infra/docker-compose.yml up -d
```

The first Python implementation can run without Docker because it uses in-memory stores. Docker is needed for persistent pgvector work in later milestones.
