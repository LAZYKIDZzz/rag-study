from typing import Any, Literal

from pydantic import BaseModel, Field


class HealthResponse(BaseModel):
    status: str
    service: str
    version: str


class CapabilitiesResponse(BaseModel):
    service: str
    features: list[str]
    embedding_provider: str
    vector_store: str
    memory_strategy: str


class DocumentCreateRequest(BaseModel):
    title: str = Field(min_length=1)
    content: str = Field(min_length=1)
    metadata: dict[str, Any] = Field(default_factory=dict)


class DocumentResponse(BaseModel):
    id: str
    title: str
    metadata: dict[str, Any]
    indexed: bool
    chunk_count: int
    created_at: str


class DocumentListResponse(BaseModel):
    documents: list[DocumentResponse]


class IndexResponse(BaseModel):
    document_id: str
    chunk_count: int
    vector_count: int


class RetrievalRequest(BaseModel):
    query: str = Field(min_length=1)
    top_k: int = Field(default=5, ge=1, le=20)
    filters: dict[str, Any] = Field(default_factory=dict)


class RetrievedChunk(BaseModel):
    chunk_id: str
    document_id: str
    document_title: str
    ordinal: int
    content: str
    score: float
    metadata: dict[str, Any]


class RetrievalResponse(BaseModel):
    query: str
    rewritten_query: str
    retrieved_chunks: list[RetrievedChunk]


class ChatSessionCreateRequest(BaseModel):
    user_id: str = "demo-user"
    title: str | None = None


class ChatSessionResponse(BaseModel):
    id: str
    user_id: str
    title: str | None = None
    messages: list["ChatMessageRecord"]
    created_at: str


class ChatMessageRecord(BaseModel):
    role: Literal["user", "assistant"]
    content: str
    created_at: str
    citations: list[RetrievedChunk] = Field(default_factory=list)


class ChatMessageRequest(BaseModel):
    user_id: str = "demo-user"
    message: str = Field(min_length=1)
    top_k: int = Field(default=5, ge=1, le=20)


class ChatMessageResponse(BaseModel):
    session_id: str
    answer: str
    rewritten_query: str
    citations: list[RetrievedChunk]
    retrieved_chunks: list[RetrievedChunk]
    memory_updates: list[str]


class AddFactRequest(BaseModel):
    fact: str = Field(min_length=1)


class MemoryResponse(BaseModel):
    user_id: str
    facts: list[str]
    recent_summary: str


ChatSessionResponse.model_rebuild()
ChatMessageResponse.model_rebuild()
