# 向量与检索

向量检索回答的是一个核心问题：用户问题和知识库片段虽然用词不同，系统如何判断它们语义相关？

## embedding 是什么

embedding 会把文本转成一组数字：

```text
"如何处理上传失败" -> [0.12, -0.03, 0.44, ...]
"文件上传报错排查" -> [0.10, -0.05, 0.40, ...]
```

这些数字不是给人读的，而是给检索算法比较的。语义越接近，向量方向通常越接近。

## 向量检索在流程中的位置

```text
索引阶段：
chunk 文本 -> embedding -> 存入向量库

问答阶段：
用户查询 -> 查询改写 -> query embedding -> 向量库 top-k -> retrieved_chunks
```

## 相似度是什么

当前项目使用余弦相似度。你可以把它理解为比较两个向量“方向是否接近”。

```text
score 越高，表示 query 和 chunk 越相似。
```

注意：相似不等于正确。一个 chunk 可能和问题语义接近，但不包含答案所需事实。因此生产系统常在向量召回后加入 rerank。

## 当前项目实现

为保证离线可运行，三后端都实现了确定性哈希 embedding。

特点：

- 不依赖外部模型 API。
- 每次运行结果可复现。
- 适合测试流程。
- 语义能力弱，不代表真实 RAG 效果。

代码入口：

- Python embedding：[`services/python-rag/app/rag/embeddings.py`](../../services/python-rag/app/rag/embeddings.py)
- Python vector store：[`services/python-rag/app/rag/vector_store.py`](../../services/python-rag/app/rag/vector_store.py)
- Java embedding：[`services/java-rag/src/main/java/study/rag/core/LocalHashEmbeddingProvider.java`](../../services/java-rag/src/main/java/study/rag/core/LocalHashEmbeddingProvider.java)
- Java vector store：[`services/java-rag/src/main/java/study/rag/core/InMemoryVectorStore.java`](../../services/java-rag/src/main/java/study/rag/core/InMemoryVectorStore.java)
- Go embedding：[`services/go-rag/internal/rag/embedding.go`](../../services/go-rag/internal/rag/embedding.go)
- Go vector store：[`services/go-rag/internal/rag/vector_store.go`](../../services/go-rag/internal/rag/vector_store.go)

## 向量检索的强项和弱项

| 类型 | 表现 |
| --- | --- |
| 同义表达 | 通常较好，例如“上传失败”和“文件上传报错” |
| 概念问题 | 通常较好，例如“为什么要分块” |
| 错误码/编号 | 可能不稳定，例如 `ERR_4012` |
| 人名/接口名/产品名 | 取决于模型和语料 |
| 权限过滤 | 不能只靠向量，必须配合 metadata filter |

这就是为什么生产 RAG 常用混合检索，而不是只用向量检索。

## 关键词检索和混合检索

关键词检索擅长精确匹配：

- 错误码。
- 接口名。
- 产品型号。
- 人名、项目名。
- 日志关键字。

混合检索会把向量检索和关键词检索结合起来：

```text
用户问题
  -> 查询改写
  -> 向量检索 top-k
  -> 关键词/BM25 检索 top-k
  -> 融合排序
  -> rerank
  -> 生成回答
```

常见融合方式：

- 分数归一化后加权。
- Reciprocal Rank Fusion。
- 先合并去重，再交给 reranker。

## rerank 为什么重要

向量库负责快速召回，rerank 负责精排。两者的目标不同。

| 阶段 | 目标 | 输入输出 |
| --- | --- | --- |
| 召回 | 多找一些可能相关的 | query -> top 20/50 chunks |
| rerank | 判断哪些最能回答问题 | query + chunks -> top 3/5 chunks |

在知识库问答中，rerank 往往能明显改善“检索到了但上下文不好”的问题。

## 生产存储选择

| 方案 | 适合场景 |
| --- | --- |
| pgvector | 想把业务表、metadata、向量放在 PostgreSQL 中统一管理 |
| Qdrant | 需要专用向量库、过滤、混合查询、向量检索能力 |
| Weaviate | 需要向量检索、混合检索、schema 管理 |
| Milvus | 更大规模向量检索 |
| Elasticsearch/OpenSearch | 关键词检索强，适合和向量能力组合 |

当前项目里程碑 M2 优先考虑 PostgreSQL + pgvector，后续可以加入 Qdrant 对照。

## 评估检索效果

不要只看“有结果”。至少要看：

- 召回率：答案所需证据是否被召回。
- 精确率：召回结果里无关片段是否太多。
- MRR/NDCG：正确片段是否排在前面。
- 引用准确率：答案引用是否真的支撑回答。
- 延迟：top-k、混合检索、rerank 是否拖慢交互。

## 关联文档

- [RAG 全流程图解](rag-pipeline.md)
- [分块 Chunking](chunking.md)
- [查询改写与记忆](query-rewriting-and-memory.md)
- [进阶 RAG 技术地图](advanced-rag-patterns.md)
- [数据模型](../architecture/data-model.md)
