# 数据模型

当前实现以“内存模型”优先，保证离线学习可运行；后续演进到 PostgreSQL + pgvector。

## 核心实体

### Document

- `id`
- `title`
- `content`
- `metadata`
- `created_at`
- `indexed`

### Chunk

- `id`
- `document_id`
- `ordinal`
- `content`
- `metadata`
- `embedding`

### ChatSession

- `id`
- `user_id`
- `created_at`
- `messages`

### Message

- `role`
- `content`
- `created_at`
- `citations`

### UserMemory

- `user_id`
- `facts`
- `recent_summary`

## 代码映射

- Python 数据类：[`services/python-rag/app/rag/models.py`](../../services/python-rag/app/rag/models.py)
- Go 类型定义：[`services/go-rag/internal/rag/types.go`](../../services/go-rag/internal/rag/types.go)
- Java 模型：[`services/java-rag/src/main/java/study/rag/core/model`](../../services/java-rag/src/main/java/study/rag/core/model)

## 存储演进方向

### 当前（M1）

- Document、Chunk、Session、Memory 都在内存中。
- 便于快速实验，不适合生产持久化。

### 目标（M2）

- 业务数据：PostgreSQL 常规表。
- 向量数据：`pgvector` 列存储 embedding。
- 检索能力：元数据过滤 + 向量相似度查询。

详见：[里程碑](../requirements/milestones.md)

## 关联文档

- [共享 API 契约](api-contract.md)
- [向量与检索](../learning/embeddings-and-vectors.md)
