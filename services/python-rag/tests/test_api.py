from fastapi.testclient import TestClient

from app.main import app


client = TestClient(app)


def test_full_rag_flow() -> None:
    document_response = client.post(
        "/documents",
        json={
            "title": "RAG basics",
            "content": (
                "Retrieval augmented generation retrieves relevant chunks before answering. "
                "Chunking splits documents into smaller units. "
                "Embeddings convert text into vectors for similarity search."
            ),
            "metadata": {"source": "test"},
        },
    )
    assert document_response.status_code == 200
    document = document_response.json()
    assert document["indexed"] is False

    index_response = client.post(f"/documents/{document['id']}/index")
    assert index_response.status_code == 200
    assert index_response.json()["chunk_count"] >= 1

    session_response = client.post("/chat/sessions", json={"user_id": "demo-user"})
    assert session_response.status_code == 200
    session = session_response.json()

    memory_response = client.post(
        "/memory/users/demo-user/facts",
        json={"fact": "The user is studying RAG workflows."},
    )
    assert memory_response.status_code == 200
    assert memory_response.json()["facts"]

    chat_response = client.post(
        f"/chat/sessions/{session['id']}/messages",
        json={"user_id": "demo-user", "message": "What does chunking do?", "top_k": 3},
    )
    assert chat_response.status_code == 200
    payload = chat_response.json()
    assert "chunking" in payload["answer"].lower()
    assert payload["citations"]
    assert payload["retrieved_chunks"]
    assert payload["memory_updates"]
    assert "Current question" in payload["rewritten_query"]


def test_validation_error_shape() -> None:
    response = client.post("/documents", json={"title": "", "content": ""})
    assert response.status_code == 422
    assert response.json()["error"]["code"] == "validation_error"
