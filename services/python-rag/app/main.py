from fastapi import FastAPI, Request
from fastapi.exceptions import RequestValidationError
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse

from app.rag.service import RagService
from app.schemas import (
    AddFactRequest,
    CapabilitiesResponse,
    ChatMessageRequest,
    ChatMessageResponse,
    ChatSessionCreateRequest,
    ChatSessionResponse,
    DocumentCreateRequest,
    DocumentListResponse,
    DocumentResponse,
    HealthResponse,
    IndexResponse,
    MemoryResponse,
    RetrievalRequest,
    RetrievalResponse,
)

app = FastAPI(
    title="RAG-study Python Backend",
    version="0.1.0",
    description="Reference Python implementation for the RAG-study API contract.",
)
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=False,
    allow_methods=["*"],
    allow_headers=["*"],
)

rag_service = RagService()


@app.exception_handler(ValueError)
async def value_error_handler(_: Request, exc: ValueError) -> JSONResponse:
    return JSONResponse(
        status_code=404,
        content={"error": {"code": "not_found", "message": str(exc), "details": {}}},
    )


@app.exception_handler(RequestValidationError)
async def validation_error_handler(_: Request, exc: RequestValidationError) -> JSONResponse:
    return JSONResponse(
        status_code=422,
        content={
            "error": {
                "code": "validation_error",
                "message": "Request validation failed",
                "details": {"errors": exc.errors()},
            }
        },
    )


@app.get("/health", response_model=HealthResponse)
def health() -> HealthResponse:
    return HealthResponse(status="ok", service="python-rag", version="0.1.0")


@app.get("/capabilities", response_model=CapabilitiesResponse)
def capabilities() -> CapabilitiesResponse:
    return rag_service.capabilities()


@app.post("/documents", response_model=DocumentResponse)
def create_document(request: DocumentCreateRequest) -> DocumentResponse:
    return rag_service.create_document(request)


@app.get("/documents", response_model=DocumentListResponse)
def list_documents() -> DocumentListResponse:
    return rag_service.list_documents()


@app.post("/documents/{document_id}/index", response_model=IndexResponse)
def index_document(document_id: str) -> IndexResponse:
    return rag_service.index_document(document_id)


@app.post("/retrieval/search", response_model=RetrievalResponse)
def search(request: RetrievalRequest) -> RetrievalResponse:
    return rag_service.search(request)


@app.post("/chat/sessions", response_model=ChatSessionResponse)
def create_session(request: ChatSessionCreateRequest) -> ChatSessionResponse:
    return rag_service.create_session(request)


@app.get("/chat/sessions/{session_id}", response_model=ChatSessionResponse)
def get_session(session_id: str) -> ChatSessionResponse:
    return rag_service.get_session(session_id)


@app.post("/chat/sessions/{session_id}/messages", response_model=ChatMessageResponse)
def send_chat_message(session_id: str, request: ChatMessageRequest) -> ChatMessageResponse:
    return rag_service.send_message(session_id, request)


@app.get("/memory/users/{user_id}", response_model=MemoryResponse)
def get_user_memory(user_id: str) -> MemoryResponse:
    return rag_service.get_memory(user_id)


@app.post("/memory/users/{user_id}/facts", response_model=MemoryResponse)
def add_memory_fact(user_id: str, request: AddFactRequest) -> MemoryResponse:
    return rag_service.add_memory_fact(user_id, request.fact)
