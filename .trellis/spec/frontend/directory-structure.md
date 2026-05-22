# Directory Structure

> How frontend code is organized in this project.

---

## Overview

The frontend lives under `apps/web` and is a React + TypeScript + Vite RAG
workbench. It should be a usable tool for document intake, indexing, retrieval
inspection, chat, and memory inspection rather than a marketing page.

---

## Directory Layout

```text
apps/web/
  src/
    App.tsx          # Workbench UI and API calls
    main.tsx         # React entrypoint
    styles.css       # Workbench styling
    vite-env.d.ts    # Vite type declarations
```

---

## Module Organization

Keep small workbench features in `App.tsx` until they become meaningfully
large. Extract components only when it improves readability or reuse. Keep API
types close to the code that consumes the shared backend contract.

---

## Naming Conventions

* React components use PascalCase.
* TypeScript types use PascalCase.
* API JSON fields should mirror backend snake_case fields at the boundary.
* CSS class names should describe UI roles, not implementation details.

---

## Examples

* `apps/web/src/App.tsx`
* `apps/web/src/styles.css`
