from dataclasses import dataclass, field
from datetime import UTC, datetime
from typing import Any, Literal
from uuid import uuid4


def utc_now() -> str:
    return datetime.now(tz=UTC).isoformat()


def new_id(prefix: str) -> str:
    return f"{prefix}_{uuid4().hex[:12]}"


@dataclass
class Document:
    title: str
    content: str
    metadata: dict[str, Any]
    id: str = field(default_factory=lambda: new_id("doc"))
    indexed: bool = False
    created_at: str = field(default_factory=utc_now)


@dataclass
class Chunk:
    document_id: str
    title: str
    ordinal: int
    content: str
    metadata: dict[str, Any]
    id: str = field(default_factory=lambda: new_id("chk"))


@dataclass
class VectorRecord:
    chunk: Chunk
    embedding: list[float]


@dataclass
class SearchResult:
    chunk: Chunk
    score: float


@dataclass
class Message:
    role: Literal["user", "assistant"]
    content: str
    citations: list[SearchResult] = field(default_factory=list)
    created_at: str = field(default_factory=utc_now)


@dataclass
class ChatSession:
    user_id: str
    title: str | None = None
    id: str = field(default_factory=lambda: new_id("ses"))
    messages: list[Message] = field(default_factory=list)
    created_at: str = field(default_factory=utc_now)


@dataclass
class UserMemory:
    user_id: str
    facts: list[str] = field(default_factory=list)
    recent_summary: str = ""
