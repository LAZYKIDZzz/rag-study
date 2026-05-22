from app.rag.models import Chunk, Document


class ParagraphChunker:
    def __init__(self, max_words: int = 120, overlap_words: int = 20) -> None:
        if overlap_words >= max_words:
            raise ValueError("overlap_words must be smaller than max_words")
        self.max_words = max_words
        self.overlap_words = overlap_words

    def chunk(self, document: Document, clean_text: str) -> list[Chunk]:
        words = clean_text.split()
        if not words:
            return []

        chunks: list[Chunk] = []
        start = 0
        ordinal = 0
        while start < len(words):
            end = min(start + self.max_words, len(words))
            chunk_words = words[start:end]
            chunks.append(
                Chunk(
                    document_id=document.id,
                    title=document.title,
                    ordinal=ordinal,
                    content=" ".join(chunk_words),
                    metadata={
                        **document.metadata,
                        "ordinal": ordinal,
                        "word_start": start,
                        "word_end": end,
                    },
                )
            )
            if end == len(words):
                break
            start = end - self.overlap_words
            ordinal += 1
        return chunks
