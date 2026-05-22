# Query Rewriting and Memory

Query rewriting improves retrieval by turning conversational messages into standalone search queries. Memory keeps useful context across turns.

## Baseline Strategy

The first implementation uses deterministic rewriting:

* include recent user facts,
* include the previous user question when available,
* keep the current user message as the main query.

This keeps the behavior inspectable before adding LLM-based rewriting.

## Memory Types

* Short-term memory: recent chat messages in a session.
* Long-term memory: user facts that persist across sessions.
* Summary memory: compact representation of older conversation context.

## Code Links

* Query processing: [services/python-rag/app/rag/query.py](../../services/python-rag/app/rag/query.py)
* Memory store: [services/python-rag/app/rag/memory.py](../../services/python-rag/app/rag/memory.py)
