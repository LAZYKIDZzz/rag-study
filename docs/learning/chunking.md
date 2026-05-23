# 分块 Chunking

分块是 RAG 中最容易被低估的一步。很多“检索不准”“答案断裂”“引用很奇怪”的问题，根因都在分块。

## 为什么不能直接检索整篇文档

如果把整篇文档作为一个检索单元，会遇到几个问题：

- 文档太长，放不进 prompt。
- 一篇文档可能有多个主题，向量会被平均，相关段落被稀释。
- 引用只能指向整篇文档，无法告诉用户具体依据。
- top-k 返回的是文档而不是证据片段，生成阶段仍然要二次寻找答案。

chunk 的目标是让每个片段都尽量满足：

- 语义完整。
- 大小适中。
- 可独立引用。
- 能和相邻片段拼接。

## 分块在流程中的位置

```text
原始文档 -> 文本清洗 -> 分块 -> chunk metadata -> embedding -> 向量库
```

分块发生在 embedding 之前。embedding 不是对整篇文档做，而是对每个 chunk 做。

## 当前项目策略

三后端都采用“固定词数 + 重叠窗口”策略：

- `max_words = 120`
- `overlap_words = 20`

示意：

```text
文档词序： 0 ........................................ 320

chunk 1:  0 ---------------------- 120
chunk 2:                    100 ---------------------- 220
chunk 3:                                      200 ---------------------- 320
```

为什么这样做：

- 120 词让单个 chunk 不至于太长，便于引用和调试。
- 20 词 overlap 降低边界切断风险。
- 固定策略最容易跨 Python / Java / Go 对照实现。

## 当前代码

- Python：[`services/python-rag/app/rag/chunking.py`](../../services/python-rag/app/rag/chunking.py)
- Java：[`services/java-rag/src/main/java/study/rag/core/ParagraphChunker.java`](../../services/java-rag/src/main/java/study/rag/core/ParagraphChunker.java)
- Go：[`services/go-rag/internal/rag/chunker.go`](../../services/go-rag/internal/rag/chunker.go)

Python chunk metadata 会记录：

- `ordinal`：chunk 在文档中的顺序。
- `word_start`：起始词位置。
- `word_end`：结束词位置。

这些字段可以用于前端展示、引用定位、排查边界问题。

## 常见分块策略

| 策略 | 做法 | 优点 | 局限 |
| --- | --- | --- | --- |
| 固定长度 | 按字符、词、token 切 | 简单稳定 | 可能切断语义 |
| 固定长度 + overlap | 相邻 chunk 重叠一部分 | 降低边界损失 | 会增加存储和重复 |
| 段落分块 | 按段落/标题切 | 更自然 | 段落长度不稳定 |
| 递归分块 | 先按标题、段落，再按句子兜底 | 兼顾结构和长度 | 实现复杂 |
| 语义分块 | 按语义变化切 | 片段更完整 | 需要模型或规则，成本更高 |
| 层级分块 | 章节摘要 + 小 chunk | 支持全局和局部问题 | 索引结构复杂 |

## 参数怎么影响效果

| 现象 | 可能原因 | 调整方向 |
| --- | --- | --- |
| 检索结果很泛 | chunk 太大 | 减小 chunk size |
| 答案缺上下文 | chunk 太小 | 增大 chunk size 或合并相邻 chunk |
| 关键句被切断 | overlap 太小 | 增大 overlap |
| 结果重复严重 | overlap 太大或 top-k 太高 | 减小 overlap 或做去重 |
| 长文档总结差 | 只有小 chunk | 加层级摘要或 RAPTOR 类索引 |

## 分块和引用的关系

知识库问答需要让用户知道“答案来自哪里”。引用通常指向 chunk。如果 chunk 太大，引用不精确；如果 chunk 太小，引用缺上下文。

一个好的 chunk 应该让用户点开引用后能马上判断：

- 这个片段确实支持答案。
- 片段本身可读。
- 如需更多上下文，可以追溯到同文档相邻片段。

## 生产演进建议

当前固定窗口适合学习。生产系统可以逐步加入：

1. 标题感知分块：把章节标题写进 chunk metadata 或 chunk 内容前缀。
2. 表格特殊处理：表头和行内容不能随意切散。
3. 代码文档特殊处理：按函数、类、接口切分。
4. 相邻 chunk 合并：检索到 chunk 后，把前后片段一起放入上下文。
5. 分块评测：统计“答案所需证据是否被完整召回”。

## 关联文档

- [RAG 全流程图解](rag-pipeline.md)
- [向量与检索](embeddings-and-vectors.md)
- [进阶 RAG 技术地图](advanced-rag-patterns.md)
- [共享 API 契约](../architecture/api-contract.md)
