# 查询改写与记忆（Query Rewriting and Memory）

在多轮对话中，用户问题往往是省略句。查询改写和记忆模块用于补全检索上下文。

## 为什么需要查询改写

用户会问：

- “那它和上一个有什么区别？”
- “再展开一下第二点。”

这类问题若直接检索，召回通常较差。改写模块会把历史上下文与用户 facts 合并，生成更完整的检索查询。

## 本项目当前策略

当前是“可解释、可测试”的确定性改写，不依赖 LLM：

- 提取用户最近 facts（最多 3 条）
- 提取上一次用户提问（Python/Java）
- 拼接当前问题形成 `rewritten_query`

## Python 实现

- QueryProcessor：[`services/python-rag/app/rag/query.py`](../../services/python-rag/app/rag/query.py)
- MemoryStore：[`services/python-rag/app/rag/memory.py`](../../services/python-rag/app/rag/memory.py)

行为要点：

- `add_fact` 会去重。
- `update_summary` 会把最近问题拼进 `recent_summary`。
- `send_message` 时会回传 `memory_updates`。

## Java / Go 对照

- Java Query：[`services/java-rag/src/main/java/study/rag/core/QueryProcessor.java`](../../services/java-rag/src/main/java/study/rag/core/QueryProcessor.java)
- Java Memory：[`services/java-rag/src/main/java/study/rag/core/MemoryStore.java`](../../services/java-rag/src/main/java/study/rag/core/MemoryStore.java)
- Go Query：[`services/go-rag/internal/rag/query.go`](../../services/go-rag/internal/rag/query.go)
- Go Memory：[`services/go-rag/internal/rag/memory.go`](../../services/go-rag/internal/rag/memory.go)

说明：Go 当前改写逻辑更简化，主要拼接 facts。

## 记忆的三层语义

- 短期记忆：会话消息（chat session messages）。
- 长期记忆：用户 facts（跨会话）。
- 摘要记忆：`recent_summary`（对近期对话的压缩表示）。

## 关联 API

- `POST /memory/users/{user_id}/facts`
- `GET /memory/users/{user_id}`
- `POST /chat/sessions/{session_id}/messages`

详见：[共享 API 契约](../architecture/api-contract.md)
