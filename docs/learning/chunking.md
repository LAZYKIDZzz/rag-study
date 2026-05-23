# 分块（Chunking）

Chunking 的目标是把长文档变成“可检索、可引用、可拼接”的片段。

## 为什么分块

如果直接把整篇文档做一次检索：

- 粒度过粗，相关段落可能被淹没。
- Prompt 会快速超长。
- 无法精确给出引用片段。

## 当前实现策略

本项目三后端都采用“固定词数 + 重叠窗口”策略：

- `max_words = 120`
- `overlap_words = 20`

其核心意图：

- 120 词保证单 chunk 可读且信息密度适中。
- 20 词重叠减少“语义刚好落在边界处”导致的召回损失。

## Python 实现细节

代码：[`services/python-rag/app/rag/chunking.py`](../../services/python-rag/app/rag/chunking.py)

该实现会把以下信息放入 `chunk.metadata`：

- `ordinal`
- `word_start`
- `word_end`

这些字段有助于：

- 前端展示 chunk 顺序
- 排查分块效果
- 后续做高亮或回溯

## Java / Go 对照

- Java：[`services/java-rag/src/main/java/study/rag/core/ParagraphChunker.java`](../../services/java-rag/src/main/java/study/rag/core/ParagraphChunker.java)
- Go：[`services/go-rag/internal/rag/chunker.go`](../../services/go-rag/internal/rag/chunker.go)

## 如何判断 chunk 参数是否合适

观察三类现象：

1. `召回不准`：可能 chunk 太大。
2. `回答断裂`：可能 chunk 太小或 overlap 太小。
3. `回答过长且重复`：可能 top_k 太高，或 chunk 过碎。

## 关联文档

- [RAG 全流程（结合代码）](rag-pipeline.md)
- [向量与检索](embeddings-and-vectors.md)
- [共享 API 契约](../architecture/api-contract.md)
