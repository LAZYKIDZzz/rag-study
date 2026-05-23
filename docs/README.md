# RAG-study 文档总览

这里是项目文档的统一入口。目标不是把 RAG 术语堆在一起，而是让没有 RAG 经验的互联网从业者，能按顺序理解“知识库问答系统为什么这样设计、每一步用什么技术、每个技术在回答问题时承担什么职责”。

## 推荐阅读方式

- 想在浏览器中读：打开 [index.html](index.html)。它会把 Markdown 文档加载成 HTML 阅读页。
- 想直接看源文档：继续阅读本目录下的 `.md` 文件。
- 如果浏览器从本地文件打开时无法加载 Markdown，在项目根目录执行：

```powershell
python -m http.server 9000
```

然后访问 `http://localhost:9000/docs/`。

## 零基础学习路径

1. [RAG 从 0 到 1](learning/rag-basics.md)：先建立基本概念，知道 RAG 解决什么问题。
2. [RAG 全流程图解](learning/rag-pipeline.md)：按“文档进入系统”和“用户提问”两条线理解完整流程。
3. [分块 Chunking](learning/chunking.md)：理解为什么不能直接把整篇文档塞给模型。
4. [向量与检索](learning/embeddings-and-vectors.md)：理解 embedding、向量库、相似度、混合检索、rerank。
5. [查询改写与记忆](learning/query-rewriting-and-memory.md)：理解多轮对话、用户记忆、查询扩展怎么提高召回。
6. [进阶 RAG 技术地图](learning/advanced-rag-patterns.md)：理解 HyDE、Multi-Query、RAG-Fusion、Self-RAG、GraphRAG 等方案适合解决什么问题。

## 工程阅读路径

1. [系统总览](architecture/system-overview.md)：看前端、三后端、RAG 模块的整体关系。
2. [共享 API 契约](architecture/api-contract.md)：看前端如何用同一套接口切换 Python / Java / Go。
3. [代码阅读地图](architecture/code-reading-map.md)：按文件顺序进入实现。
4. [三后端实现对照](architecture/backend-comparison.md)：比较三种语言落地同一 RAG 流程的差异。
5. [本地开发与联调](runbooks/local-development.md)：跑通前端和任一后端。

## 文档分区

| 目录 | 作用 |
| --- | --- |
| `learning/` | 面向 RAG 初学者的概念、流程、技术方案 |
| `architecture/` | 系统结构、API、数据模型、代码入口 |
| `requirements/` | 产品目标、范围、里程碑 |
| `runbooks/` | 本地运行、联调、排障 |
| `decisions/` | 关键技术决策记录 |

## 当前实现与目标实现的关系

当前代码优先保证“离线可运行、流程可观察、适合学习”：

- embedding 使用确定性哈希向量，不依赖外部模型。
- vector store 使用内存实现，便于测试和演示。
- answer generator 是抽取式生成器，便于看清引用来源。
- query rewriting 和 memory 已有基础实现，但还不是生产级策略。

后续目标是逐步演进到生产常见配置：

- 真实 embedding 模型。
- PostgreSQL + pgvector 或专用向量数据库。
- 关键词 + 向量混合检索。
- rerank、查询扩展、查询改写、记忆摘要、评测回归。

## 技术参考

- RAG 原始论文：[Retrieval-Augmented Generation for Knowledge-Intensive NLP Tasks](https://arxiv.org/abs/2005.11401)
- Self-RAG 论文：[Self-RAG: Learning to Retrieve, Generate, and Critique through Self-Reflection](https://arxiv.org/abs/2310.11511)
- Microsoft GraphRAG 文档：[GraphRAG Overview](https://microsoft.github.io/graphrag/index/overview/)
- LangChain MultiQueryRetriever 文档：[MultiQueryRetriever](https://api.python.langchain.com/en/latest/langchain/retrievers/langchain.retrievers.multi_query.MultiQueryRetriever.html)
