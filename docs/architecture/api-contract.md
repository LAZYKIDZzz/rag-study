# 共享 API 契约

三后端（Python / Java / Go）应提供等价接口与核心字段，前端可只改 base URL 即切换后端。

## 设计目标

1. 前端对后端实现无感知。
2. 学习者能在同一调用流程下比较不同语言。
3. 错误结构统一，便于联调与测试。

## 1. 健康检查

### `GET /health`

```json
{
  "status": "ok",
  "service": "python-rag",
  "version": "0.1.0"
}
```

## 2. 能力说明

### `GET /capabilities`

```json
{
  "service": "python-rag",
  "features": ["document-ingestion", "chunking", "chat-memory"],
  "embedding_provider": "local-hash-embedding-128d",
  "vector_store": "in-memory-cosine",
  "memory_strategy": "in-memory"
}
```

## 3. 文档管理

### `POST /documents`

请求：

```json
{
  "title": "RAG notes",
  "content": "Raw document text",
  "metadata": {
    "source": "manual"
  }
}
```

返回（示例）：

```json
{
  "id": "doc_xxx",
  "title": "RAG notes",
  "metadata": {"source": "manual"},
  "indexed": false,
  "chunk_count": 0,
  "created_at": "2026-05-23T12:00:00Z"
}
```

### `GET /documents`

```json
{
  "documents": [
    {
      "id": "doc_xxx",
      "title": "RAG notes",
      "metadata": {"source": "manual"},
      "indexed": true,
      "chunk_count": 3,
      "created_at": "2026-05-23T12:00:00Z"
    }
  ]
}
```

### `POST /documents/{document_id}/index`

触发：清洗 -> 分块 -> 向量化 -> 向量入库。

返回（示例）：

```json
{
  "document_id": "doc_xxx",
  "chunk_count": 3,
  "vector_count": 3
}
```

## 4. 检索

### `POST /retrieval/search`

请求：

```json
{
  "query": "What is chunking?",
  "top_k": 5,
  "filters": {}
}
```

返回（示例）：

```json
{
  "query": "What is chunking?",
  "rewritten_query": "Current question: What is chunking?",
  "retrieved_chunks": [
    {
      "chunk_id": "chk_xxx",
      "document_id": "doc_xxx",
      "document_title": "RAG notes",
      "ordinal": 0,
      "content": "Chunking splits documents into searchable passages.",
      "score": 0.82,
      "metadata": {"source": "manual"}
    }
  ]
}
```

## 5. 对话

### `POST /chat/sessions`

创建会话。支持 body 传 `user_id`。

### `GET /chat/sessions/{session_id}`

获取会话消息。

### `POST /chat/sessions/{session_id}/messages`

请求：

```json
{
  "user_id": "demo-user",
  "message": "Explain embeddings",
  "top_k": 5
}
```

返回（示例）：

```json
{
  "session_id": "ses_xxx",
  "answer": "Based on the retrieved knowledge...",
  "rewritten_query": "Current question: Explain embeddings",
  "citations": [],
  "retrieved_chunks": [],
  "memory_updates": ["Last asked: Explain embeddings"]
}
```

## 6. 记忆

### `GET /memory/users/{user_id}`

获取用户 facts 与 recent summary。

### `POST /memory/users/{user_id}/facts`

新增 fact：

```json
{
  "fact": "The user is studying RAG workflows."
}
```

## 7. 错误结构

```json
{
  "error": {
    "code": "validation_error",
    "message": "Request validation failed",
    "details": {}
  }
}
```

## 实现对照

- Python 路由：[`services/python-rag/app/main.py`](../../services/python-rag/app/main.py)
- Java 控制器：[`services/java-rag/src/main/java/study/rag/api/RagController.java`](../../services/java-rag/src/main/java/study/rag/api/RagController.java)
- Go HTTP：[`services/go-rag/internal/rag/http.go`](../../services/go-rag/internal/rag/http.go)

## 关联文档

- [系统总览](system-overview.md)
- [三后端实现对照](backend-comparison.md)
- [本地开发与联调](../runbooks/local-development.md)
