import math

from app.rag.models import Chunk, SearchResult, VectorRecord


class InMemoryVectorStore:
    def __init__(self) -> None:
        self._records: dict[str, VectorRecord] = {}

    def upsert_many(self, chunks: list[Chunk], embeddings: list[list[float]]) -> None:
        for chunk, embedding in zip(chunks, embeddings, strict=True):
            self._records[chunk.id] = VectorRecord(chunk=chunk, embedding=embedding)

    def delete_by_document(self, document_id: str) -> None:
        chunk_ids = [
            chunk_id
            for chunk_id, record in self._records.items()
            if record.chunk.document_id == document_id
        ]
        for chunk_id in chunk_ids:
            del self._records[chunk_id]

    def search(
        self,
        query_embedding: list[float],
        top_k: int,
        filters: dict[str, object] | None = None,
    ) -> list[SearchResult]:
        filters = filters or {}
        scored: list[SearchResult] = []
        for record in self._records.values():
            if not self._matches_filters(record.chunk, filters):
                continue
            score = cosine_similarity(query_embedding, record.embedding)
            scored.append(SearchResult(chunk=record.chunk, score=score))

        scored.sort(key=lambda result: result.score, reverse=True)
        return scored[:top_k]

    def count_for_document(self, document_id: str) -> int:
        return sum(1 for record in self._records.values() if record.chunk.document_id == document_id)

    def _matches_filters(self, chunk: Chunk, filters: dict[str, object]) -> bool:
        for key, expected in filters.items():
            if chunk.metadata.get(key) != expected:
                return False
        return True


def cosine_similarity(left: list[float], right: list[float]) -> float:
    if len(left) != len(right):
        return 0.0

    dot = sum(a * b for a, b in zip(left, right, strict=True))
    left_length = math.sqrt(sum(a * a for a in left))
    right_length = math.sqrt(sum(b * b for b in right))
    if left_length == 0 or right_length == 0:
        return 0.0
    return dot / (left_length * right_length)
