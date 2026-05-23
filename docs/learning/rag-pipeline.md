# RAG 全流程（结合代码）

本页把“RAG 概念流程”与“项目代码路径”一一对应，便于边学边看实现。

## 流程总览

```text
文档写入 -> 文本清洗 -> 分块 -> 向量化 -> 向量存储
                                      |
用户提问 -> 查询改写 -> 查询向量化 -> 相似度检索 -> 回答生成 -> 记忆更新
```

## 阶段 1：文档写入

接口：`POST /documents`

职责：保存原始文档与元数据。

关键代码：

- Python API：[`services/python-rag/app/main.py`](../../services/python-rag/app/main.py)
- Python 仓储：[`services/python-rag/app/rag/repositories.py`](../../services/python-rag/app/rag/repositories.py)

## 阶段 2：文本清洗

职责：统一空白、减少后续分块噪音。

关键代码：

- Python：[`services/python-rag/app/rag/preprocessing.py`](../../services/python-rag/app/rag/preprocessing.py)

## 阶段 3：分块（Chunking）

职责：把文档切成可检索的片段。

当前策略：`max_words=120`，`overlap_words=20`。

关键代码：

- Python：[`services/python-rag/app/rag/chunking.py`](../../services/python-rag/app/rag/chunking.py)
- Java：[`services/java-rag/src/main/java/study/rag/core/ParagraphChunker.java`](../../services/java-rag/src/main/java/study/rag/core/ParagraphChunker.java)
- Go：[`services/go-rag/internal/rag/chunker.go`](../../services/go-rag/internal/rag/chunker.go)

扩展阅读：[分块（Chunking）](chunking.md)

## 阶段 4：向量化（Embedding）

职责：把 chunk 文本转成向量。

当前策略：离线可运行的确定性哈希向量（学习用途，不是生产语义模型）。

关键代码：

- Python：[`services/python-rag/app/rag/embeddings.py`](../../services/python-rag/app/rag/embeddings.py)
- Java：[`services/java-rag/src/main/java/study/rag/core/LocalHashEmbeddingProvider.java`](../../services/java-rag/src/main/java/study/rag/core/LocalHashEmbeddingProvider.java)
- Go：[`services/go-rag/internal/rag/embedding.go`](../../services/go-rag/internal/rag/embedding.go)

扩展阅读：[向量与检索](embeddings-and-vectors.md)

## 阶段 5：向量存储与检索

职责：

- 写入 chunk + embedding
- 查询时计算相似度并返回 top-k

当前策略：内存向量库 + 余弦相似度。

关键代码：

- Python：[`services/python-rag/app/rag/vector_store.py`](../../services/python-rag/app/rag/vector_store.py)
- Java：[`services/java-rag/src/main/java/study/rag/core/InMemoryVectorStore.java`](../../services/java-rag/src/main/java/study/rag/core/InMemoryVectorStore.java)
- Go：[`services/go-rag/internal/rag/vector_store.go`](../../services/go-rag/internal/rag/vector_store.go)

## 阶段 6：查询改写

职责：补充用户 facts、历史问题，使检索查询更完整。

关键代码：

- Python：[`services/python-rag/app/rag/query.py`](../../services/python-rag/app/rag/query.py)
- Java：[`services/java-rag/src/main/java/study/rag/core/QueryProcessor.java`](../../services/java-rag/src/main/java/study/rag/core/QueryProcessor.java)
- Go：[`services/go-rag/internal/rag/query.go`](../../services/go-rag/internal/rag/query.go)

扩展阅读：[查询改写与记忆](query-rewriting-and-memory.md)

## 阶段 7：回答生成

职责：根据检索结果生成回答，并返回引用。

当前策略：抽取式回答生成器（离线可跑）。

关键代码：

- Python：[`services/python-rag/app/rag/generation.py`](../../services/python-rag/app/rag/generation.py)
- Java：[`services/java-rag/src/main/java/study/rag/core/ExtractiveAnswerGenerator.java`](../../services/java-rag/src/main/java/study/rag/core/ExtractiveAnswerGenerator.java)
- Go：[`services/go-rag/internal/rag/chat.go`](../../services/go-rag/internal/rag/chat.go)

## 阶段 8：记忆更新

职责：更新用户长期 facts 与近期摘要。

关键代码：

- Python：[`services/python-rag/app/rag/memory.py`](../../services/python-rag/app/rag/memory.py)
- Java：[`services/java-rag/src/main/java/study/rag/core/MemoryStore.java`](../../services/java-rag/src/main/java/study/rag/core/MemoryStore.java)
- Go：[`services/go-rag/internal/rag/memory.go`](../../services/go-rag/internal/rag/memory.go)

## 流程编排入口

- Python 主编排：[`services/python-rag/app/rag/service.py`](../../services/python-rag/app/rag/service.py)
- Java 主编排：[`services/java-rag/src/main/java/study/rag/core/RagService.java`](../../services/java-rag/src/main/java/study/rag/core/RagService.java)
- Go 主编排：[`services/go-rag/internal/rag/service.go`](../../services/go-rag/internal/rag/service.go)

## 关联文档

- [共享 API 契约](../architecture/api-contract.md)
- [代码阅读地图](../architecture/code-reading-map.md)
- [系统总览](../architecture/system-overview.md)
