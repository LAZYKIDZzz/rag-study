import { FormEvent, useEffect, useMemo, useState } from "react";

const API_BASE_URL =
  import.meta.env.VITE_RAG_API_BASE_URL?.replace(/\/$/, "") ?? "http://localhost:8000";

type DocumentRecord = {
  id: string;
  title: string;
  metadata: Record<string, unknown>;
  indexed: boolean;
  chunk_count?: number;
  created_at: string;
};

type RetrievedChunk = {
  chunk_id?: string;
  id?: string;
  document_id: string;
  document_title?: string;
  title?: string;
  ordinal?: number;
  content: string;
  score: number;
  metadata: Record<string, unknown>;
};

type ChatSession = {
  id: string;
  user_id: string;
  title?: string | null;
  messages: Array<{
    role: "user" | "assistant";
    content: string;
    created_at: string;
    citations?: RetrievedChunk[];
  }>;
};

type Memory = {
  user_id: string;
  facts: string[];
  recent_summary: string;
};

type ChatResponse = {
  answer: string;
  rewritten_query: string;
  citations: RetrievedChunk[];
  retrieved_chunks: RetrievedChunk[];
  memory_updates: string[];
};

type Status = { kind: "idle" | "busy" | "error" | "ok"; message: string };

async function api<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...init?.headers
    }
  });
  const payload = await response.json().catch(() => ({}));
  if (!response.ok) {
    const message = payload?.error?.message ?? `Request failed with ${response.status}`;
    throw new Error(message);
  }
  return payload as T;
}

function chunksFromSearch(payload: { retrieved_chunks?: RetrievedChunk[]; results?: RetrievedChunk[] }) {
  return payload.retrieved_chunks ?? payload.results ?? [];
}

export function App() {
  const [documents, setDocuments] = useState<DocumentRecord[]>([]);
  const [session, setSession] = useState<ChatSession | null>(null);
  const [memory, setMemory] = useState<Memory | null>(null);
  const [retrievedChunks, setRetrievedChunks] = useState<RetrievedChunk[]>([]);
  const [answer, setAnswer] = useState<ChatResponse | null>(null);
  const [status, setStatus] = useState<Status>({ kind: "idle", message: "Ready" });
  const [userId, setUserId] = useState("demo-user");
  const [title, setTitle] = useState("RAG notes");
  const [content, setContent] = useState(
    "Chunking splits documents into searchable passages. Embeddings map text into vectors. Retrieval selects relevant chunks before answer generation."
  );
  const [metadata, setMetadata] = useState('{"source":"manual"}');
  const [question, setQuestion] = useState("What does chunking do?");
  const [fact, setFact] = useState("The user is studying RAG workflows.");

  const indexedCount = useMemo(() => documents.filter((document) => document.indexed).length, [documents]);

  useEffect(() => {
    void refreshDocuments();
    void refreshMemory(userId);
  }, []);

  async function refreshDocuments() {
    const payload = await api<{ documents?: DocumentRecord[] } | DocumentRecord[]>("/documents");
    setDocuments(Array.isArray(payload) ? payload : payload.documents ?? []);
  }

  async function refreshMemory(nextUserId = userId) {
    const payload = await api<Memory>(`/memory/users/${encodeURIComponent(nextUserId)}`);
    setMemory(payload);
  }

  async function runTask<T>(message: string, task: () => Promise<T>) {
    setStatus({ kind: "busy", message });
    try {
      const result = await task();
      setStatus({ kind: "ok", message: "Synced" });
      return result;
    } catch (error) {
      setStatus({ kind: "error", message: error instanceof Error ? error.message : "Unknown error" });
      throw error;
    }
  }

  async function addDocument(event: FormEvent) {
    event.preventDefault();
    await runTask("Adding document", async () => {
      const parsedMetadata = metadata.trim() ? JSON.parse(metadata) : {};
      const document = await api<DocumentRecord>("/documents", {
        method: "POST",
        body: JSON.stringify({ title, content, metadata: parsedMetadata })
      });
      setDocuments((current) => [document, ...current]);
    });
  }

  async function indexDocument(documentId: string) {
    await runTask("Indexing document", async () => {
      await api(`/documents/${documentId}/index`, { method: "POST" });
      await refreshDocuments();
    });
  }

  async function createSession() {
    await runTask("Creating chat session", async () => {
      const created = await api<ChatSession>("/chat/sessions", {
        method: "POST",
        body: JSON.stringify({ user_id: userId, title: "Workbench session" })
      });
      setSession(created);
      await refreshMemory(userId);
    });
  }

  async function addFact() {
    await runTask("Saving memory", async () => {
      const payload = await api<Memory>(`/memory/users/${encodeURIComponent(userId)}/facts`, {
        method: "POST",
        body: JSON.stringify({ fact })
      });
      setMemory(payload);
    });
  }

  async function searchOnly() {
    await runTask("Retrieving chunks", async () => {
      const payload = await api<{ retrieved_chunks?: RetrievedChunk[]; results?: RetrievedChunk[] }>(
        "/retrieval/search",
        {
          method: "POST",
          body: JSON.stringify({ query: question, top_k: 5, filters: {} })
        }
      );
      setRetrievedChunks(chunksFromSearch(payload));
    });
  }

  async function sendMessage() {
    await runTask("Generating answer", async () => {
      let activeSession = session;
      if (!activeSession) {
        activeSession = await api<ChatSession>("/chat/sessions", {
          method: "POST",
          body: JSON.stringify({ user_id: userId, title: "Workbench session" })
        });
        setSession(activeSession);
      }

      const payload = await api<ChatResponse>(`/chat/sessions/${activeSession.id}/messages`, {
        method: "POST",
        body: JSON.stringify({ user_id: userId, message: question, top_k: 5 })
      });
      setAnswer(payload);
      setRetrievedChunks(payload.retrieved_chunks ?? payload.citations ?? []);
      await refreshMemory(userId);
      const updated = await api<ChatSession>(`/chat/sessions/${activeSession.id}`);
      setSession(updated);
    });
  }

  return (
    <main className="shell">
      <header className="topbar">
        <div>
          <p className="eyebrow">RAG-study</p>
          <h1>Knowledge Workbench</h1>
        </div>
        <div className={`status ${status.kind}`}>
          <span>{status.message}</span>
          <small>{API_BASE_URL}</small>
        </div>
      </header>

      <section className="summary" aria-label="Backend summary">
        <div>
          <strong>{documents.length}</strong>
          <span>documents</span>
        </div>
        <div>
          <strong>{indexedCount}</strong>
          <span>indexed</span>
        </div>
        <div>
          <strong>{retrievedChunks.length}</strong>
          <span>chunks in trace</span>
        </div>
        <label>
          User
          <input value={userId} onChange={(event) => setUserId(event.target.value)} onBlur={() => refreshMemory(userId)} />
        </label>
      </section>

      <div className="workspace">
        <section className="panel compose">
          <div className="panelHeading">
            <h2>Document Intake</h2>
            <button onClick={() => refreshDocuments()} type="button">Refresh</button>
          </div>
          <form onSubmit={addDocument}>
            <label>
              Title
              <input value={title} onChange={(event) => setTitle(event.target.value)} />
            </label>
            <label>
              Content
              <textarea value={content} onChange={(event) => setContent(event.target.value)} rows={7} />
            </label>
            <label>
              Metadata JSON
              <input value={metadata} onChange={(event) => setMetadata(event.target.value)} />
            </label>
            <button className="primary" type="submit">Add document</button>
          </form>
        </section>

        <section className="panel documents">
          <div className="panelHeading">
            <h2>Documents</h2>
            <span>{documents.length} loaded</span>
          </div>
          <div className="list">
            {documents.map((document) => (
              <article className="item" key={document.id}>
                <div>
                  <h3>{document.title}</h3>
                  <p>{document.id}</p>
                </div>
                <div className="itemMeta">
                  <span>{document.indexed ? "Indexed" : "Raw"}</span>
                  <span>{document.chunk_count ?? 0} chunks</span>
                </div>
                <button onClick={() => indexDocument(document.id)} type="button">Index</button>
              </article>
            ))}
            {documents.length === 0 && <p className="empty">Add a document to start the retrieval loop.</p>}
          </div>
        </section>

        <section className="panel chat">
          <div className="panelHeading">
            <h2>Chat And Retrieval</h2>
            <button onClick={createSession} type="button">{session ? "New session" : "Create session"}</button>
          </div>
          <label>
            Question
            <textarea value={question} onChange={(event) => setQuestion(event.target.value)} rows={3} />
          </label>
          <div className="actions">
            <button onClick={searchOnly} type="button">Retrieve</button>
            <button className="primary" onClick={sendMessage} type="button">Ask</button>
          </div>
          {answer && (
            <div className="answer">
              <h3>Answer</h3>
              <p>{answer.answer}</p>
              <small>Rewritten query: {answer.rewritten_query}</small>
            </div>
          )}
          <div className="trace">
            {retrievedChunks.map((chunk) => (
              <article className="chunk" key={chunk.chunk_id ?? chunk.id}>
                <header>
                  <strong>{chunk.document_title ?? chunk.title ?? chunk.document_id}</strong>
                  <span>{chunk.score.toFixed(4)}</span>
                </header>
                <p>{chunk.content}</p>
              </article>
            ))}
          </div>
        </section>

        <section className="panel memory">
          <div className="panelHeading">
            <h2>User Memory</h2>
            <button onClick={() => refreshMemory(userId)} type="button">Reload</button>
          </div>
          <label>
            Fact
            <input value={fact} onChange={(event) => setFact(event.target.value)} />
          </label>
          <button onClick={addFact} type="button">Save fact</button>
          <div className="memoryBlock">
            <h3>Facts</h3>
            {(memory?.facts.length ?? 0) > 0 ? (
              <ul>{memory?.facts.map((item) => <li key={item}>{item}</li>)}</ul>
            ) : (
              <p className="empty">No saved facts for this user.</p>
            )}
            <h3>Recent summary</h3>
            <p>{memory?.recent_summary || "No chat summary yet."}</p>
          </div>
        </section>
      </div>
    </main>
  );
}

