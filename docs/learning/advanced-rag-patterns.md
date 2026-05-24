# 进阶 RAG 技术地图

这篇回答一个常见问题：基础 RAG 跑通后，为什么还需要查询改写、关键词联想、混合检索、rerank、记忆关联、GraphRAG 这些方案？

结论先说：它们不是同一层的替代品，而是分别解决 RAG 链路中的不同失败点。

## 先看技术位置

```text
用户问题
  -> 查询理解
      -> 查询改写
      -> 关键词联想
      -> 多查询生成
      -> HyDE
  -> 召回
      -> 向量检索
      -> 关键词检索
      -> 混合检索
      -> 图检索
  -> 精排
      -> rerank
      -> 去重
      -> 相邻 chunk 合并
  -> 生成
      -> grounded generation
      -> 引用约束
      -> Self-RAG 式自检
  -> 记忆
      -> 会话记忆
      -> 用户 facts
      -> 摘要记忆
      -> 实体关系记忆
```

## 1. 查询改写

查询改写把“不适合检索的用户表达”改成“适合检索的查询”。

例子：

```text
历史：我们刚讨论了 RAG 的分块策略。
用户：那 overlap 设多少合适？

改写后：RAG 文档分块 chunk overlap 应该设置多少，过大或过小分别有什么影响？
```

它解决的问题：

- 多轮对话里的代词和省略。
- 口语表达太短，检索信号不足。
- 用户问题混合了多个意图，需要拆分。

当前项目实现：

- Python / Java：结合用户 facts 和上一轮问题。
- Go：主要拼接用户 facts，后续可补齐上下文改写。

代码入口：

- [`services/python-rag/app/rag/query.py`](../../services/python-rag/app/rag/query.py)
- [`services/java-rag/src/main/java/study/rag/core/QueryProcessor.java`](../../services/java-rag/src/main/java/study/rag/core/QueryProcessor.java)
- [`services/go-rag/internal/rag/query.go`](../../services/go-rag/internal/rag/query.go)

## 2. 关键词联想

关键词联想是给查询补充同义词、业务词、缩写、错误码、产品名。

例子：

```text
用户：上传失败怎么排查？

联想关键词：
upload failed, 上传报错, 文件上传, ERR_UPLOAD, 超时, 鉴权, token 过期, 文件大小限制
```

它解决的问题：

- 用户说“上传失败”，文档写的是“upload error”。
- 用户说“登录不了”，文档写的是“鉴权失败”。
- 用户说“卡住”，文档写的是具体错误码。

常见实现：

- 维护业务词典。
- 从文档标题、标签、错误码中抽取关键词。
- 用 LLM 生成同义改写。
- 使用 BM25/关键词检索补充向量检索。

在知识库问答里的作用：提高召回率，尤其适合客服、运维、接口文档、产品手册。

## 3. 多查询检索

多查询检索会从同一个问题生成多个检索角度。

```text
原问题：RAG 为什么检索结果不准？

查询 A：RAG 检索召回不准的原因
查询 B：chunk size overlap 对召回的影响
查询 C：embedding 模型质量和向量库相似度问题
查询 D：混合检索 rerank 如何改善 RAG
```

它解决的问题：

- 单个查询表达太窄。
- 用户问题背后可能有多个原因。
- 向量检索对 phrasing 敏感。

常见方案：

- Multi-Query Retriever：生成多条改写查询并合并结果。
- RAG-Fusion：多查询检索后用 Reciprocal Rank Fusion 之类方法融合排序。

代价：

- 检索次数增加。
- 需要去重和结果融合。
- 不适合极低延迟场景直接无脑开启。

## 4. HyDE

HyDE 的思路是：先让模型根据问题写一个“假设性答案”，再用这个假设答案去检索。

```text
用户问题 -> 生成假设答案 -> 对假设答案做 embedding -> 检索真实文档
```

它解决的问题：

- 原始问题太短，embedding 信号弱。
- 文档和问题表达差异大。
- 用户问的是抽象概念，文档里是具体描述。

风险：

- 假设答案可能引入错误方向。
- 必须把 HyDE 结果当“检索辅助”，不能当真实答案。

## 5. 混合检索

混合检索结合向量检索和关键词检索。

```text
向量检索：找到语义相近内容
关键词检索：找到精确词、编号、错误码、函数名
融合排序：把两路结果合并
```

适合场景：

- 技术文档：接口名、类名、错误码很重要。
- 客服知识库：用户表达和文档表达差异大。
- 产品手册：专有名词、型号、版本号不能丢。

常见技术：

- Dense vector：语义检索。
- Sparse vector / BM25：关键词检索。
- RRF：融合多路排序。
- Metadata filter：按业务线、权限、时间过滤。

当前项目尚未实现混合检索，建议作为 M4 的关键能力之一。

## 6. rerank

rerank 在召回之后工作。召回负责“多找一些可能相关的”，rerank 负责“把真正适合回答的问题片段排到前面”。

```text
召回 top 30 -> rerank -> 取 top 5 -> 生成回答
```

它解决的问题：

- 向量相似但信息不完整。
- chunk 相关但不能回答用户问题。
- 多路召回后分数不可比。

常见实现：

- cross-encoder reranker。
- 云服务 rerank API。
- LLM 打分，但成本和延迟更高。

在知识库问答里的作用：提升上下文质量，减少模型看到无关资料。

## 7. 记忆关联

记忆关联不是“把所有历史消息都塞给模型”，而是判断哪些历史信息会影响当前问题。

例子：

```text
用户事实：用户负责 B 端 SaaS 客服系统。
当前问题：我们的知识库应该怎么分块？

关联后查询：B 端 SaaS 客服知识库 RAG 文档分块策略，SOP、FAQ、工单记录如何切分
```

常见记忆类型：

| 类型 | 内容 | 用途 |
| --- | --- | --- |
| 短期会话 | 最近几轮消息 | 处理追问、省略 |
| 用户 facts | 角色、偏好、业务背景 | 个性化回答和检索 |
| 摘要记忆 | 对长对话的压缩 | 控制上下文长度 |
| 实体记忆 | 用户、项目、产品、指标关系 | 复杂业务关联 |

注意点：

- 记忆要有权限边界。
- 记忆要能更新、过期、删除。
- 敏感信息不能随意进入检索上下文。

## 8. Self-RAG

Self-RAG 的核心思想是让模型学习在生成过程中判断是否需要检索，并对生成内容进行反思和批判。它更接近“生成阶段的质量控制”。

它解决的问题：

- 不是每个问题都需要检索。
- 检索结果可能不支持当前回答。
- 模型需要判断答案是否被证据支撑。

在本项目中的定位：

- 当前不实现 Self-RAG。
- 可以作为后续“答案质量自检”和“是否需要检索”的研究方向。

## 9. RAPTOR

RAPTOR 会把文档片段递归聚类并生成摘要，形成树状索引。检索时既能查细粒度 chunk，也能查高层摘要。

它解决的问题：

- 单个 chunk 太碎，难回答全局总结类问题。
- 用户问“这份文档整体讲了什么”时，普通 chunk 检索不够好。
- 长文档需要层级理解。

适合场景：

- 长报告。
- 研究文档。
- 会议纪要。
- 多章节知识库。

## 10. GraphRAG

GraphRAG 会从文档中抽取实体和关系，构建图结构，再结合图上的社区摘要或关系检索回答问题。

它解决的问题：

- 答案依赖多个实体之间的关系。
- 信息分散在多个文档里。
- 问题需要跨文档归纳，而不是单段事实查找。

适合场景：

- 企业知识网络。
- 风控、投研、情报分析。
- 有大量实体、组织、项目、事件关系的资料库。

代价：

- 抽取和建图成本高。
- 更新链路复杂。
- 需要额外评估图谱质量。

## 能力选择表

| 现象 | 优先考虑 |
| --- | --- |
| 用户追问“它”“刚才那个”检索不到 | 查询改写 + 短期记忆 |
| 术语、错误码、接口名检索不到 | 关键词联想 + 混合检索 |
| 查询表达换一下结果差很多 | 多查询检索 / RAG-Fusion |
| 问题太短或太抽象 | HyDE |
| 召回结果有相关但不适合回答的片段 | rerank |
| 长文档总结效果差 | RAPTOR / 层级摘要 |
| 多文档实体关系推理差 | GraphRAG |
| 答案缺少证据约束 | 引用约束 + Self-RAG 式自检 |

## 本项目推荐演进顺序

1. 先把当前基础 RAG 链路跑稳。
2. 接入真实 embedding 和持久化向量库。
3. 加入关键词检索，形成混合检索。
4. 增加 rerank。
5. 增强 query rewriting：LLM 改写、多查询、关键词联想。
6. 完善 memory：短期会话、facts、摘要、过期策略。
7. 增加评测集，验证每个增强是否真的提升效果。
8. 最后再研究 GraphRAG、RAPTOR、Self-RAG 等更复杂方案。

## 参考资料

- RAG 原始论文：[Retrieval-Augmented Generation for Knowledge-Intensive NLP Tasks](https://arxiv.org/abs/2005.11401)
- Self-RAG 论文：[Self-RAG: Learning to Retrieve, Generate, and Critique through Self-Reflection](https://arxiv.org/abs/2310.11511)
- RAPTOR 论文：[RAPTOR: Recursive Abstractive Processing for Tree-Organized Retrieval](https://arxiv.org/abs/2401.18059)
- Microsoft GraphRAG 文档：[GraphRAG Overview](https://microsoft.github.io/graphrag/index/overview/)
- LangChain 文档：[MultiQueryRetriever](https://api.python.langchain.com/en/latest/langchain/retrievers/langchain.retrievers.multi_query.MultiQueryRetriever.html)
- LlamaIndex 文档：[HyDE Query Transform](https://docs.llamaindex.ai/en/stable/examples/query_transformations/HyDEQueryTransformDemo/)
- Qdrant 文档：[Hybrid Queries](https://qdrant.tech/documentation/concepts/hybrid-queries/)
- Cohere 文档：[Rerank](https://docs.cohere.com/docs/reranking-with-cohere)
