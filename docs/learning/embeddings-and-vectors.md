# 向量与检索（Embeddings and Vectors）

本页说明“向量化为什么能检索文本”，并对应到项目实现。

## 核心概念

- Embedding：把文本映射到一个定长向量。
- 相似度：通常使用余弦相似度衡量“方向是否接近”。
- Top-k 检索：返回最相似的前 k 个 chunk。

## 本项目当前策略

为保证离线可运行，三后端都实现了“确定性哈希 embedding”。

它的特点：

- 不依赖外部模型 API。
- 结果可复现，便于测试。
- 语义能力弱于生产模型，只用于学习与流程验证。

## Python 实现

- Embedding：[`services/python-rag/app/rag/embeddings.py`](../../services/python-rag/app/rag/embeddings.py)
- Vector Store：[`services/python-rag/app/rag/vector_store.py`](../../services/python-rag/app/rag/vector_store.py)

你会看到：

- Token 经过哈希后写入固定维度桶。
- 向量归一化后用于余弦相似度。
- 查询时按分数倒序取 `top_k`。

## Java / Go 对照

- Java embedding：[`services/java-rag/src/main/java/study/rag/core/LocalHashEmbeddingProvider.java`](../../services/java-rag/src/main/java/study/rag/core/LocalHashEmbeddingProvider.java)
- Java vector store：[`services/java-rag/src/main/java/study/rag/core/InMemoryVectorStore.java`](../../services/java-rag/src/main/java/study/rag/core/InMemoryVectorStore.java)
- Go embedding：[`services/go-rag/internal/rag/embedding.go`](../../services/go-rag/internal/rag/embedding.go)
- Go vector store：[`services/go-rag/internal/rag/vector_store.go`](../../services/go-rag/internal/rag/vector_store.go)

## 面向生产的演进方向

- 接入真实语义 embedding 模型。
- 将向量持久化到 `pgvector` 或专用向量数据库。
- 增加混合检索（关键词 + 向量）与 rerank。

关联里程碑：[milestones.md](../requirements/milestones.md)

## 关联文档

- [RAG 全流程（结合代码）](rag-pipeline.md)
- [数据模型](../architecture/data-model.md)
- [系统总览](../architecture/system-overview.md)
