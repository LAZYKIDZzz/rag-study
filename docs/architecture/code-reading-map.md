# 代码阅读地图

本页给出“从哪看、按什么顺序看、每个文件看什么”的最短路径。

## 0. 建议阅读顺序

1. API 入口（看接口面）
2. Service 编排（看流程面）
3. RAG 子模块（看算法面）
4. 前端调用（看联调面）

## 1. Python（推荐第一站）

### 1.1 API 入口

- [`services/python-rag/app/main.py`](../../services/python-rag/app/main.py)

看点：

- 路由定义是否和 API 契约一致。
- 错误处理与校验错误结构。

### 1.2 流程编排

- [`services/python-rag/app/rag/service.py`](../../services/python-rag/app/rag/service.py)

看点：

- `index_document`：文档索引主链路。
- `search`：纯检索链路。
- `send_message`：聊天 + 检索 + 记忆更新链路。

### 1.3 RAG 模块

- 预处理：[`services/python-rag/app/rag/preprocessing.py`](../../services/python-rag/app/rag/preprocessing.py)
- 分块：[`services/python-rag/app/rag/chunking.py`](../../services/python-rag/app/rag/chunking.py)
- 向量化：[`services/python-rag/app/rag/embeddings.py`](../../services/python-rag/app/rag/embeddings.py)
- 向量检索：[`services/python-rag/app/rag/vector_store.py`](../../services/python-rag/app/rag/vector_store.py)
- 查询改写：[`services/python-rag/app/rag/query.py`](../../services/python-rag/app/rag/query.py)
- 回答生成：[`services/python-rag/app/rag/generation.py`](../../services/python-rag/app/rag/generation.py)
- 记忆：[`services/python-rag/app/rag/memory.py`](../../services/python-rag/app/rag/memory.py)

### 1.4 测试

- [`services/python-rag/tests/test_api.py`](../../services/python-rag/tests/test_api.py)

看点：

- 端到端 happy path。
- 校验失败时错误结构。

## 2. Java（与 Python 对照）

- API：[`services/java-rag/src/main/java/study/rag/api/RagController.java`](../../services/java-rag/src/main/java/study/rag/api/RagController.java)
- 编排：[`services/java-rag/src/main/java/study/rag/core/RagService.java`](../../services/java-rag/src/main/java/study/rag/core/RagService.java)
- 分块：[`services/java-rag/src/main/java/study/rag/core/ParagraphChunker.java`](../../services/java-rag/src/main/java/study/rag/core/ParagraphChunker.java)
- 改写：[`services/java-rag/src/main/java/study/rag/core/QueryProcessor.java`](../../services/java-rag/src/main/java/study/rag/core/QueryProcessor.java)

## 3. Go（与 Python 对照）

- API：[`services/go-rag/internal/rag/http.go`](../../services/go-rag/internal/rag/http.go)
- 编排：[`services/go-rag/internal/rag/service.go`](../../services/go-rag/internal/rag/service.go)
- 类型：[`services/go-rag/internal/rag/types.go`](../../services/go-rag/internal/rag/types.go)
- 分块：[`services/go-rag/internal/rag/chunker.go`](../../services/go-rag/internal/rag/chunker.go)
- 改写：[`services/go-rag/internal/rag/query.go`](../../services/go-rag/internal/rag/query.go)

## 4. 前端调用入口

- [`apps/web/src/App.tsx`](../../apps/web/src/App.tsx)

看点：

- `api()` 统一请求。
- 文档写入、索引、检索、聊天、记忆对应哪些按钮动作。
- 如何消费 `retrieved_chunks` 和 `memory_updates`。

## 关联文档

- [系统总览](system-overview.md)
- [RAG 全流程（结合代码）](../learning/rag-pipeline.md)
- [共享 API 契约](api-contract.md)
