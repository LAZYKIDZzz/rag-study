# Quality Guidelines

> Code quality standards for frontend development.

---

## Overview

The frontend is an operational RAG workbench. It should be responsive, readable,
and contract-focused. Prefer clear state and visible retrieval traces over
decorative UI.

---

## Forbidden Patterns

* Do not build a landing page as the primary screen.
* Do not hide retrieval traces, scores, rewritten queries, or memory state.
* Do not introduce backend-specific UI branches unless the shared contract
  cannot express the behavior.
* Do not use decorative gradients/orbs or nested cards that reduce scanability.

---

## Required Patterns

* Read the API base URL from `VITE_RAG_API_BASE_URL`, defaulting to
  `http://localhost:8000`.
* Keep loading/error status visible.
* Support document add, index, retrieval search, chat message, and memory fact
  workflows.
* Keep text wrapping resilient for long IDs, chunk text, and error messages.

---

## Testing Requirements

Run `npm run build` in `apps/web` before committing frontend changes.

---

## Code Review Checklist

* UI calls the shared API contract.
* TypeScript build passes.
* The workbench is usable on desktop and mobile widths.
* Retrieval traces and memory are visible without opening browser devtools.

---

## Scenario: RAG Workbench API Boundary

### 1. Scope / Trigger

Trigger: any frontend change that reads or writes RAG API payloads. The
frontend is the visible contract consumer for Python, Java, and Go backends.

### 2. Signatures

The workbench must call:

* `POST /documents`
* `GET /documents`
* `POST /documents/{document_id}/index`
* `POST /retrieval/search`
* `POST /chat/sessions`
* `POST /chat/sessions/{session_id}/messages`
* `GET /chat/sessions/{session_id}`
* `GET /memory/users/{user_id}`
* `POST /memory/users/{user_id}/facts`

### 3. Contracts

* API base URL comes from `VITE_RAG_API_BASE_URL`; default is
  `http://localhost:8000`.
* Document list responses are read from `documents`.
* Retrieval traces are read from `retrieved_chunks`.
* Chat responses display `answer`, `rewritten_query`, `citations`,
  `retrieved_chunks`, and `memory_updates`.

### 4. Validation & Error Matrix

* Non-2xx response with shared error shape -> show `error.message`.
* Network failure -> show a visible status error.
* Invalid metadata JSON -> keep the operation in error state and do not submit
  malformed payload.

### 5. Good/Base/Bad Cases

* Good: user adds a document, indexes it, asks a question, sees answer,
  rewritten query, retrieved chunks, and memory.
* Base: no documents shows an empty state that points the user to document
  intake.
* Bad: backend-specific UI branch is added for Java/Go response names instead
  of fixing the shared contract.

### 6. Tests Required

* `npm run build` must pass.
* When adding a test runner later, add mocked API tests for document intake,
  retrieval, chat, and memory fact save.

### 7. Wrong vs Correct

Wrong:

```ts
setRetrievedChunks(payload.results);
```

Correct:

```ts
setRetrievedChunks(payload.retrieved_chunks);
```
