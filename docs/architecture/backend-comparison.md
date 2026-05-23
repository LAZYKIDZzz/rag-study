# 三后端实现对照（Python / Java / Go）

本页用于回答一个核心问题：`同一套 RAG 能力，在三种语言里如何落地？`

## 对照结论（先看这个）

- Python：开发效率高，适合作为教学与原型基线。
- Java：工程分层清晰，类型约束强，适合企业化扩展。
- Go：依赖少、控制流直接，适合理解底层实现细节。

## 路由层对照

- Python FastAPI：[`services/python-rag/app/main.py`](../../services/python-rag/app/main.py)
- Java Spring Controller：[`services/java-rag/src/main/java/study/rag/api/RagController.java`](../../services/java-rag/src/main/java/study/rag/api/RagController.java)
- Go net/http：[`services/go-rag/internal/rag/http.go`](../../services/go-rag/internal/rag/http.go)

差异点：

- Python/Java 有框架级参数校验能力。
- Go 当前自行处理 JSON 解析和参数错误。

## 核心编排层对照

- Python：[`services/python-rag/app/rag/service.py`](../../services/python-rag/app/rag/service.py)
- Java：[`services/java-rag/src/main/java/study/rag/core/RagService.java`](../../services/java-rag/src/main/java/study/rag/core/RagService.java)
- Go：[`services/go-rag/internal/rag/service.go`](../../services/go-rag/internal/rag/service.go)

共同点：

- 都按“文档 -> 分块 -> 向量 -> 检索 -> 回答 -> 记忆”组织。
- 都支持 `retrieved_chunks`、`rewritten_query`、`memory_updates`。

## 分块策略对照

- Python：[`services/python-rag/app/rag/chunking.py`](../../services/python-rag/app/rag/chunking.py)
- Java：[`services/java-rag/src/main/java/study/rag/core/ParagraphChunker.java`](../../services/java-rag/src/main/java/study/rag/core/ParagraphChunker.java)
- Go：[`services/go-rag/internal/rag/chunker.go`](../../services/go-rag/internal/rag/chunker.go)

共同策略：固定窗口 + overlap。

## 查询改写对照

- Python：支持用户 facts + 上一轮用户问题。
- Java：与 Python 基本一致。
- Go：当前更简化，主要拼接用户 facts。

对应代码：

- Python：[`services/python-rag/app/rag/query.py`](../../services/python-rag/app/rag/query.py)
- Java：[`services/java-rag/src/main/java/study/rag/core/QueryProcessor.java`](../../services/java-rag/src/main/java/study/rag/core/QueryProcessor.java)
- Go：[`services/go-rag/internal/rag/query.go`](../../services/go-rag/internal/rag/query.go)

## 错误处理对照

- Python：FastAPI 异常处理器统一 `error` 结构。
- Java：全局异常处理器统一输出。
- Go：`writeServiceError` + `writeError` 明确映射状态码和错误码。

## 选择建议

- 教学入门、快速迭代：优先 Python。
- 强工程治理、复杂业务落地：优先 Java。
- 高可控、低依赖、部署简单：优先 Go。

## 关联文档

- [共享 API 契约](api-contract.md)
- [代码阅读地图](code-reading-map.md)
- [RAG 全流程（结合代码）](../learning/rag-pipeline.md)
