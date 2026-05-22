# Chunking

Chunking splits documents into retrievable pieces. Chunk size and overlap directly affect answer quality.

## Baseline Strategy

The first implementation uses paragraph-aware fixed-size chunking:

* normalize whitespace,
* keep paragraphs where possible,
* split long paragraphs by word count,
* add overlap between chunks,
* retain document ID and ordinal metadata.

## Trade-Offs

Small chunks retrieve precise snippets but may lose context. Large chunks preserve context but can dilute retrieval scores and prompt space.

## Code Link

* [services/python-rag/app/rag/chunking.py](../../services/python-rag/app/rag/chunking.py)
