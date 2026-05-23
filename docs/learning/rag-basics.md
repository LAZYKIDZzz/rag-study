# RAG 从 0 到 1（面向零基础）

本页假设你完全不了解 RAG。读完后你应该能回答三件事：

1. RAG 到底解决什么问题。
2. RAG 在工程里分成哪些阶段。
3. 本项目代码里每个阶段在哪。

## 1. 为什么需要 RAG

大模型本身存在三个常见问题：

- 训练数据有时间边界，可能不知道你的私有文档。
- 回答可能“看起来合理但实际错误”（幻觉）。
- 很难说明“答案来自哪里”。

RAG（Retrieval-Augmented Generation，检索增强生成）通过“先检索、再回答”降低这些问题。

## 2. 一句话理解 RAG

RAG = `把你的文档做成可搜索知识库` + `回答前先从知识库取证据`。

## 3. 本项目中的 RAG 主流程

1. 写入文档（`/documents`）
2. 文档索引（`/documents/{id}/index`）
   - 清洗文本
   - 切分 chunk
   - 生成向量
   - 存入向量库
3. 用户提问（`/retrieval/search` 或 `/chat/sessions/{id}/messages`）
   - 查询改写
   - 查询向量化
   - 相似度检索
4. 生成回答（带引用片段）
5. 更新用户记忆（facts + recent summary）

详见：[RAG 全流程（结合代码）](rag-pipeline.md)

## 4. 三个最容易混淆的概念

### 4.1 Chunk 是什么

chunk 是“可检索的文本片段”，不是整篇文档。检索时比整文更精准。

详见：[分块（Chunking）](chunking.md)

### 4.2 Embedding 是什么

embedding 是把文本映射为向量。语义越接近，向量越接近。

详见：[向量与检索](embeddings-and-vectors.md)

### 4.3 Query Rewriting 是什么

把口语化或上下文依赖问题改写为更适合检索的查询。

详见：[查询改写与记忆](query-rewriting-and-memory.md)

## 5. 先看哪套代码

建议先看 Python 实现，因为链路最直观、最完整：

- 入口：[`services/python-rag/app/main.py`](../../services/python-rag/app/main.py)
- 编排：[`services/python-rag/app/rag/service.py`](../../services/python-rag/app/rag/service.py)

然后再看 Java/Go 对照：

- Java：[`services/java-rag/src/main/java/study/rag/core/RagService.java`](../../services/java-rag/src/main/java/study/rag/core/RagService.java)
- Go：[`services/go-rag/internal/rag/service.go`](../../services/go-rag/internal/rag/service.go)

## 下一步

- 想继续概念：读 [RAG 全流程（结合代码）](rag-pipeline.md)
- 想动手跑：读 [本地开发与联调](../runbooks/local-development.md)
