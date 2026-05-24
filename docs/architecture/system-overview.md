# 系统总览

RAG-study 是一个学习优先的知识库问答平台。它同时关注两件事：

- 能跑：前端可以调用后端完成文档写入、索引、检索、聊天、记忆观察。
- 能学：每个 RAG 阶段都在代码和文档中显式可见，方便从零理解流程。

## 总体架构

```text
                              ┌────────────────────┐
                              │ apps/web React UI  │
                              └─────────┬──────────┘
                                        │ Shared HTTP API
               ┌────────────────────────┼────────────────────────┐
               │                        │                        │
┌──────────────▼──────────────┐ ┌───────▼──────────────┐ ┌───────▼──────────────┐
│ services/python-rag FastAPI │ │ services/java-rag    │ │ services/go-rag      │
│ 学习基线与优先阅读入口       │ │ Spring Boot 对照实现  │ │ net/http 对照实现     │
└──────────────┬──────────────┘ └───────┬──────────────┘ └───────┬──────────────┘
               │                        │                        │
               └────────────────────────┼────────────────────────┘
                                        ▼
                         RAG pipeline concepts and contracts
```

## RAG 模块边界

每个后端都围绕同一条链路组织：

```text
文档写入
  -> 文本清洗
  -> 分块
  -> embedding
  -> 向量存储

用户提问
  -> 查询改写
  -> 检索召回
  -> 回答生成 + 引用
  -> 记忆更新
```

这种设计的价值：

- 初学者能把文档中的概念映射到代码文件。
- 前端能用同一 API 切换不同后端。
- 后续替换真实 embedding、向量库、LLM、rerank 时，不需要重写整个系统。

## 当前学习版实现

| 能力 | 当前实现 | 目的 |
| --- | --- | --- |
| 文档存储 | 内存仓储 | 降低本地启动门槛 |
| 文本清洗 | 空白归一化等轻量规则 | 展示预处理位置 |
| 分块 | 固定词数 + overlap | 展示 chunk 粒度和边界问题 |
| embedding | 确定性哈希向量 | 离线可运行、可测试 |
| 向量检索 | 内存向量库 + 余弦相似度 | 展示 top-k 召回 |
| 查询改写 | facts + 历史问题拼接 | 展示多轮上下文进入检索 |
| 回答生成 | 抽取式生成器 | 展示证据和引用，不依赖 LLM |
| 记忆 | facts + recent summary | 展示用户记忆和会话记忆 |

## 生产版演进方向

```text
文档解析：PDF / Word / HTML / 表格 / 网页
分块策略：标题感知 / 语义分块 / 层级分块
检索策略：向量 + 关键词 + metadata filter + rerank
生成策略：LLM grounded generation + 引用约束 + 答案自检
记忆策略：短期会话 + 用户 facts + 摘要 + 实体关系
评测策略：召回率 / 引用准确率 / 答案忠实度 / 延迟
```

更具体的进阶能力见：[进阶 RAG 技术地图](../learning/advanced-rag-patterns.md)

## 代码入口

- 前端：[`apps/web/src/App.tsx`](../../apps/web/src/App.tsx)
- Python API：[`services/python-rag/app/main.py`](../../services/python-rag/app/main.py)
- Python 编排：[`services/python-rag/app/rag/service.py`](../../services/python-rag/app/rag/service.py)
- Java API：[`services/java-rag/src/main/java/study/rag/api/RagController.java`](../../services/java-rag/src/main/java/study/rag/api/RagController.java)
- Java 编排：[`services/java-rag/src/main/java/study/rag/core/RagService.java`](../../services/java-rag/src/main/java/study/rag/core/RagService.java)
- Go API：[`services/go-rag/internal/rag/http.go`](../../services/go-rag/internal/rag/http.go)
- Go 编排：[`services/go-rag/internal/rag/service.go`](../../services/go-rag/internal/rag/service.go)

## 进一步阅读

- [RAG 全流程图解](../learning/rag-pipeline.md)
- [共享 API 契约](api-contract.md)
- [代码阅读地图](code-reading-map.md)
- [三后端实现对照](backend-comparison.md)
