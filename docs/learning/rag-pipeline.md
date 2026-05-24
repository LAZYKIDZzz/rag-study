# RAG 全流程图解

RAG 不是一个单点算法，而是一条数据流水线。理解它最好的方式是分成两条线：

- 索引线：文档如何进入知识库。
- 问答线：用户问题如何变成答案。

## 总览图

```text
┌─────────────────────────────── 索引线：把文档做成可搜索知识库 ───────────────────────────────┐
│                                                                                              │
│  原始文档 -> 文本解析/清洗 -> 分块 -> 元数据补充 -> embedding -> 向量库/关键词索引              │
│                                                                                              │
└──────────────────────────────────────────────────────────────────────────────────────────────┘
                                                   │
                                                   ▼
┌─────────────────────────────── 问答线：把问题变成有证据的答案 ───────────────────────────────┐
│                                                                                              │
│  用户问题 -> 意图识别 -> 查询改写/扩展 -> 检索召回 -> rerank -> 构造上下文 -> 生成回答 -> 记忆更新 │
│                                                                                              │
└──────────────────────────────────────────────────────────────────────────────────────────────┘
```

## 阶段 1：文档写入

接口：`POST /documents`

输入是标题、正文和元数据。元数据可以包含来源、业务线、权限标签、更新时间、作者等。它的价值不只在展示，还会用于后续过滤检索，例如“只查客服 SOP”或“只查 2026 年文档”。

当前代码：

- Python API：[`services/python-rag/app/main.py`](../../services/python-rag/app/main.py)
- Python 仓储：[`services/python-rag/app/rag/repositories.py`](../../services/python-rag/app/rag/repositories.py)
- 数据模型：[`services/python-rag/app/rag/models.py`](../../services/python-rag/app/rag/models.py)

为什么需要这一步：

- 知识库问答必须知道答案来自哪份文档。
- 后续引用、权限控制、增量更新都依赖文档实体。

## 阶段 2：文本解析与清洗

真实系统里，文档可能来自 PDF、Word、网页、Markdown、Excel、工单系统。解析会把不同格式统一成文本，清洗会处理空白、换行、页眉页脚、乱码、重复段落等噪音。

当前项目只处理已经传入的纯文本，并做轻量清洗：

- Python：[`services/python-rag/app/rag/preprocessing.py`](../../services/python-rag/app/rag/preprocessing.py)
- Java：[`services/java-rag/src/main/java/study/rag/core/TextPreprocessor.java`](../../services/java-rag/src/main/java/study/rag/core/TextPreprocessor.java)

为什么需要这一步：

- 噪音文本会污染 chunk。
- chunk 被污染后，embedding 和检索都会受影响。
- 解析质量往往决定 RAG 上限，尤其是 PDF、表格、代码文档。

## 阶段 3：分块

分块把长文档切成适合检索的片段。检索粒度通常不应该是整篇文档，因为整篇文档太大、主题太杂，也不利于引用。

当前策略：固定词数窗口 + overlap。

```text
chunk 1: word 0   - word 120
chunk 2: word 100 - word 220
chunk 3: word 200 - word 320
```

当前代码：

- Python：[`services/python-rag/app/rag/chunking.py`](../../services/python-rag/app/rag/chunking.py)
- Java：[`services/java-rag/src/main/java/study/rag/core/ParagraphChunker.java`](../../services/java-rag/src/main/java/study/rag/core/ParagraphChunker.java)
- Go：[`services/go-rag/internal/rag/chunker.go`](../../services/go-rag/internal/rag/chunker.go)

为什么需要 overlap：

- 一段关键语义可能刚好跨越边界。
- overlap 可以降低边界切断导致的召回损失。

详见：[分块 Chunking](chunking.md)

## 阶段 4：embedding

embedding 把文本转成定长数字向量。向量的意义是：语义相近的文本在向量空间更接近。

```text
"上传失败排查" -> [0.12, -0.33, 0.07, ...]
"文件上传报错如何处理" -> [0.10, -0.30, 0.08, ...]
```

当前项目使用确定性哈希 embedding：

- 离线可运行。
- 便于测试。
- 语义能力弱，不代表生产效果。

当前代码：

- Python：[`services/python-rag/app/rag/embeddings.py`](../../services/python-rag/app/rag/embeddings.py)
- Java：[`services/java-rag/src/main/java/study/rag/core/LocalHashEmbeddingProvider.java`](../../services/java-rag/src/main/java/study/rag/core/LocalHashEmbeddingProvider.java)
- Go：[`services/go-rag/internal/rag/embedding.go`](../../services/go-rag/internal/rag/embedding.go)

生产系统通常会替换为真实 embedding 模型，例如 OpenAI-compatible embedding、本地 embedding 模型，或云厂商模型。

详见：[向量与检索](embeddings-and-vectors.md)

## 阶段 5：索引存储

索引存储至少要保存：

- chunk 文本。
- chunk 所属 document。
- chunk 顺序。
- chunk metadata。
- chunk embedding。

当前策略：内存向量库 + 余弦相似度。

当前代码：

- Python：[`services/python-rag/app/rag/vector_store.py`](../../services/python-rag/app/rag/vector_store.py)
- Java：[`services/java-rag/src/main/java/study/rag/core/InMemoryVectorStore.java`](../../services/java-rag/src/main/java/study/rag/core/InMemoryVectorStore.java)
- Go：[`services/go-rag/internal/rag/vector_store.go`](../../services/go-rag/internal/rag/vector_store.go)

生产演进：

- PostgreSQL + pgvector：适合把业务数据、元数据、向量放在一个数据库里。
- Qdrant / Weaviate / Milvus：适合专用向量检索和更复杂的过滤能力。
- Elasticsearch / OpenSearch：适合关键词检索和混合检索。

## 阶段 6：用户问题理解

用户问题经常不是完整检索语句：

- “它和上一个有什么区别？”
- “第二点展开讲讲。”
- “按我们的场景怎么做？”

问答线通常会先做：

- 意图识别：用户是在查知识、要求总结、要求对比，还是继续追问。
- 查询改写：把省略句补完整。
- 查询扩展：补充同义词、关键词、业务术语。
- 多查询生成：从多个角度检索，避免单一查询漏召回。

当前代码：

- Python：[`services/python-rag/app/rag/query.py`](../../services/python-rag/app/rag/query.py)
- Java：[`services/java-rag/src/main/java/study/rag/core/QueryProcessor.java`](../../services/java-rag/src/main/java/study/rag/core/QueryProcessor.java)
- Go：[`services/go-rag/internal/rag/query.go`](../../services/go-rag/internal/rag/query.go)

详见：[查询改写与记忆](query-rewriting-and-memory.md)

## 阶段 7：检索召回

检索召回的目标不是立刻给最终答案，而是先找出“可能有用”的候选片段。

常见召回方式：

| 方式 | 擅长 | 局限 |
| --- | --- | --- |
| 向量检索 | 语义相近、表达不同 | 对精确关键词、编号、专有名词可能不稳定 |
| 关键词检索 | 精确词、错误码、产品名、接口名 | 同义表达召回差 |
| 混合检索 | 结合语义和关键词 | 需要融合分数和调参 |
| 图检索 | 实体关系、跨文档推理 | 建图成本高 |

当前实现是向量检索。后续里程碑会加入混合检索和 rerank。

## 阶段 8：rerank

召回阶段通常宁可多拿一些候选，rerank 再精排。rerank 模型会逐条判断“这个 chunk 对当前问题是否真的有帮助”。

```text
召回 top 30 -> rerank -> 取 top 5 放进上下文
```

为什么需要 rerank：

- 向量相似不等于可回答。
- chunk 可能语义接近但缺少关键事实。
- rerank 能减少无关上下文占用 prompt。

当前项目暂未实现 rerank，文档中把它列为 M4 进阶能力。

## 阶段 9：构造上下文

构造上下文是把检索结果整理成模型可读的 prompt。关键问题是：

- 放几个 chunk。
- 按什么顺序放。
- 是否合并同一文档相邻 chunk。
- 是否附带标题、来源、时间、metadata。
- 是否去重。

一个简化上下文可能长这样：

```text
用户问题：
上传失败时应该先排查什么？

可用资料：
[1] 文档：客服 SOP，片段 3
上传失败通常先检查文件大小、格式、网络超时和鉴权状态...

[2] 文档：错误码说明，片段 8
ERR_UPLOAD_TOKEN_EXPIRED 表示上传凭证过期，需要重新获取...
```

## 阶段 10：回答生成与引用

当前项目使用抽取式回答生成器，不调用 LLM。它会基于检索结果组织一个可解释答案，并返回引用。

当前代码：

- Python：[`services/python-rag/app/rag/generation.py`](../../services/python-rag/app/rag/generation.py)
- Java：[`services/java-rag/src/main/java/study/rag/core/ExtractiveAnswerGenerator.java`](../../services/java-rag/src/main/java/study/rag/core/ExtractiveAnswerGenerator.java)
- Go：[`services/go-rag/internal/rag/chat.go`](../../services/go-rag/internal/rag/chat.go)

生产系统中，生成阶段通常使用 LLM，并要求：

- 只基于给定资料回答。
- 不知道就说不知道。
- 给出引用。
- 不把引用和推断混在一起。

## 阶段 11：记忆更新

记忆不是把所有聊天记录无限塞进 prompt，而是把对后续有用的信息结构化保存。

常见记忆：

- 会话短期记忆：最近几轮对话。
- 用户事实记忆：用户偏好、角色、业务背景。
- 摘要记忆：把长对话压缩成摘要。
- 实体记忆：用户、项目、产品、指标之间的关系。

当前代码：

- Python：[`services/python-rag/app/rag/memory.py`](../../services/python-rag/app/rag/memory.py)
- Java：[`services/java-rag/src/main/java/study/rag/core/MemoryStore.java`](../../services/java-rag/src/main/java/study/rag/core/MemoryStore.java)
- Go：[`services/go-rag/internal/rag/memory.go`](../../services/go-rag/internal/rag/memory.go)

## 主编排入口

三后端都把这些阶段收敛在 service 层：

- Python：[`services/python-rag/app/rag/service.py`](../../services/python-rag/app/rag/service.py)
- Java：[`services/java-rag/src/main/java/study/rag/core/RagService.java`](../../services/java-rag/src/main/java/study/rag/core/RagService.java)
- Go：[`services/go-rag/internal/rag/service.go`](../../services/go-rag/internal/rag/service.go)

## 当前能力与后续能力对照

| 能力 | 当前项目 | 生产演进 |
| --- | --- | --- |
| 文档解析 | 纯文本 | PDF/Word/HTML/表格解析 |
| 分块 | 固定词数 + overlap | 语义分块、标题感知分块、层级分块 |
| embedding | 哈希向量 | 真实语义 embedding |
| 检索 | 内存向量 top-k | pgvector/Qdrant + filter + hybrid search |
| 查询改写 | facts + recent question 拼接 | LLM 改写、多查询、HyDE、关键词联想 |
| rerank | 暂无 | cross-encoder / rerank API |
| 生成 | 抽取式 | LLM grounded generation |
| 记忆 | facts + recent summary | 摘要、实体、权限隔离、过期策略 |
| 评测 | 基础测试 | 召回率、引用准确率、答案忠实度回归 |

## 关联文档

- [RAG 从 0 到 1](rag-basics.md)
- [查询改写与记忆](query-rewriting-and-memory.md)
- [进阶 RAG 技术地图](advanced-rag-patterns.md)
- [共享 API 契约](../architecture/api-contract.md)
- [代码阅读地图](../architecture/code-reading-map.md)
