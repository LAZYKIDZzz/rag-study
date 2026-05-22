from app.rag.chunking import ParagraphChunker
from app.rag.embeddings import LocalHashEmbeddingProvider
from app.rag.generation import ExtractiveAnswerGenerator
from app.rag.memory import MemoryStore
from app.rag.models import Document, Message, SearchResult
from app.rag.preprocessing import TextPreprocessor
from app.rag.query import QueryProcessor
from app.rag.repositories import ChatRepository, DocumentRepository
from app.rag.vector_store import InMemoryVectorStore
from app.schemas import (
    CapabilitiesResponse,
    ChatMessageRecord,
    ChatMessageRequest,
    ChatMessageResponse,
    ChatSessionCreateRequest,
    ChatSessionResponse,
    DocumentCreateRequest,
    DocumentListResponse,
    DocumentResponse,
    IndexResponse,
    MemoryResponse,
    RetrievalRequest,
    RetrievalResponse,
    RetrievedChunk,
)


class RagService:
    def __init__(self) -> None:
        self.documents = DocumentRepository()
        self.chats = ChatRepository()
        self.memory = MemoryStore()
        self.preprocessor = TextPreprocessor()
        self.chunker = ParagraphChunker()
        self.embeddings = LocalHashEmbeddingProvider()
        self.vector_store = InMemoryVectorStore()
        self.query_processor = QueryProcessor()
        self.generator = ExtractiveAnswerGenerator()

    def capabilities(self) -> CapabilitiesResponse:
        return CapabilitiesResponse(
            service="python-rag",
            features=[
                "document-ingestion",
                "preprocessing",
                "chunking",
                "local-embeddings",
                "in-memory-vector-search",
                "query-rewriting",
                "chat-memory",
                "extractive-answer-generation",
            ],
            embedding_provider="local-hash-embedding-128d",
            vector_store="in-memory-cosine",
            memory_strategy="in-memory",
        )

    def create_document(self, request: DocumentCreateRequest) -> DocumentResponse:
        document = Document(
            title=request.title.strip(),
            content=request.content,
            metadata=request.metadata,
        )
        self.documents.add(document)
        return self._document_response(document)

    def list_documents(self) -> DocumentListResponse:
        return DocumentListResponse(
            documents=[self._document_response(document) for document in self.documents.list()]
        )

    def index_document(self, document_id: str) -> IndexResponse:
        document = self.documents.get(document_id)
        clean_text = self.preprocessor.clean(document.content)
        chunks = self.chunker.chunk(document, clean_text)
        embeddings = [self.embeddings.embed(chunk.content) for chunk in chunks]
        self.vector_store.delete_by_document(document.id)
        self.vector_store.upsert_many(chunks, embeddings)
        self.documents.replace_chunks(document.id, chunks)
        return IndexResponse(
            document_id=document.id,
            chunk_count=len(chunks),
            vector_count=self.vector_store.count_for_document(document.id),
        )

    def search(self, request: RetrievalRequest) -> RetrievalResponse:
        rewritten_query = self.query_processor.rewrite(
            question=request.query,
            session=None,
            memory=self.memory.get("anonymous"),
        )
        results = self._retrieve(rewritten_query, request.top_k, request.filters)
        return RetrievalResponse(
            query=request.query,
            rewritten_query=rewritten_query,
            retrieved_chunks=[self._retrieved_chunk(result) for result in results],
        )

    def create_session(self, request: ChatSessionCreateRequest) -> ChatSessionResponse:
        session = self.chats.create(user_id=request.user_id, title=request.title)
        return self._session_response(session.id)

    def get_session(self, session_id: str) -> ChatSessionResponse:
        return self._session_response(session_id)

    def send_message(self, session_id: str, request: ChatMessageRequest) -> ChatMessageResponse:
        session = self.chats.get(session_id)
        user_memory = self.memory.get(request.user_id)
        rewritten_query = self.query_processor.rewrite(request.message, session, user_memory)
        results = self._retrieve(rewritten_query, request.top_k, {})
        answer = self.generator.generate(request.message, results, user_memory)

        self.chats.add_message(session_id, Message(role="user", content=request.message))
        self.chats.add_message(session_id, Message(role="assistant", content=answer, citations=results))
        updated_memory = self.memory.update_summary(request.user_id, request.message)

        return ChatMessageResponse(
            session_id=session_id,
            answer=answer,
            rewritten_query=rewritten_query,
            citations=[self._retrieved_chunk(result) for result in results],
            retrieved_chunks=[self._retrieved_chunk(result) for result in results],
            memory_updates=[updated_memory.recent_summary],
        )

    def get_memory(self, user_id: str) -> MemoryResponse:
        return self._memory_response(user_id)

    def add_memory_fact(self, user_id: str, fact: str) -> MemoryResponse:
        self.memory.add_fact(user_id, fact)
        return self._memory_response(user_id)

    def _retrieve(
        self,
        query: str,
        top_k: int,
        filters: dict[str, object],
    ) -> list[SearchResult]:
        query_embedding = self.embeddings.embed(query)
        return self.vector_store.search(query_embedding, top_k=top_k, filters=filters)

    def _document_response(self, document: Document) -> DocumentResponse:
        return DocumentResponse(
            id=document.id,
            title=document.title,
            metadata=document.metadata,
            indexed=document.indexed,
            chunk_count=self.documents.chunk_count(document.id),
            created_at=document.created_at,
        )

    def _retrieved_chunk(self, result: SearchResult) -> RetrievedChunk:
        return RetrievedChunk(
            chunk_id=result.chunk.id,
            document_id=result.chunk.document_id,
            document_title=result.chunk.title,
            ordinal=result.chunk.ordinal,
            content=result.chunk.content,
            score=round(result.score, 6),
            metadata=result.chunk.metadata,
        )

    def _session_response(self, session_id: str) -> ChatSessionResponse:
        session = self.chats.get(session_id)
        return ChatSessionResponse(
            id=session.id,
            user_id=session.user_id,
            title=session.title,
            messages=[
                ChatMessageRecord(
                    role=message.role,
                    content=message.content,
                    created_at=message.created_at,
                    citations=[self._retrieved_chunk(result) for result in message.citations],
                )
                for message in session.messages
            ],
            created_at=session.created_at,
        )

    def _memory_response(self, user_id: str) -> MemoryResponse:
        memory = self.memory.get(user_id)
        return MemoryResponse(
            user_id=memory.user_id,
            facts=memory.facts,
            recent_summary=memory.recent_summary,
        )
