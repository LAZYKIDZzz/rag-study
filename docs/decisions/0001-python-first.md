# ADR 0001: Python First Backend

## Status

Accepted.

## Context

The project must produce a complete RAG system and also compare Java, Python, and Go backends. Implementing all three fully before validating the product flow would multiply design mistakes.

## Decision

Implement Python as the first complete backend using FastAPI and explicit RAG modules. Keep Java and Go service directories aligned to the same API contract so they can be completed against a proven reference.

## Consequences

* The first runnable path arrives faster.
* The RAG workflow can be learned in a mature ecosystem first.
* Java and Go parity follows after the contract is validated.
