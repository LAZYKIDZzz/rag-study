from __future__ import annotations

from app.rag.models import ChatSession, Chunk, Document, Message


class DocumentRepository:
    def __init__(self) -> None:
        self._documents: dict[str, Document] = {}
        self._chunks_by_document: dict[str, list[Chunk]] = {}

    def add(self, document: Document) -> Document:
        self._documents[document.id] = document
        return document

    def get(self, document_id: str) -> Document:
        document = self._documents.get(document_id)
        if document is None:
            raise ValueError(f"Document {document_id} was not found")
        return document

    def list(self) -> list[Document]:
        return sorted(self._documents.values(), key=lambda document: document.created_at, reverse=True)

    def replace_chunks(self, document_id: str, chunks: list[Chunk]) -> None:
        document = self.get(document_id)
        self._chunks_by_document[document_id] = chunks
        document.indexed = True

    def chunks_for_document(self, document_id: str) -> list[Chunk]:
        return self._chunks_by_document.get(document_id, [])

    def chunk_count(self, document_id: str) -> int:
        return len(self._chunks_by_document.get(document_id, []))


class ChatRepository:
    def __init__(self) -> None:
        self._sessions: dict[str, ChatSession] = {}

    def create(self, user_id: str, title: str | None = None) -> ChatSession:
        session = ChatSession(user_id=user_id, title=title)
        self._sessions[session.id] = session
        return session

    def get(self, session_id: str) -> ChatSession:
        session = self._sessions.get(session_id)
        if session is None:
            raise ValueError(f"Chat session {session_id} was not found")
        return session

    def add_message(self, session_id: str, message: Message) -> ChatSession:
        session = self.get(session_id)
        session.messages.append(message)
        return session
