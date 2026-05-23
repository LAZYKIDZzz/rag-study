# 系统总览

本项目是“学习优先”的 RAG 实验平台：

- 前端：React Workbench，用于文档写入、检索追踪、聊天与记忆观察。
- 后端：Python / Java / Go 三实现，遵循同一 API 契约。
- 目标：同一 RAG 流程下，对比不同语言生态的工程组织方式。

## 总体架构

```text
apps/web (React)
  -> Shared API Contract
    -> services/python-rag (FastAPI)
    -> services/java-rag (Spring Boot)
    -> services/go-rag (net/http)

Each backend contains:
  preprocessing -> chunking -> embeddings -> vector search -> answer generation -> memory
```

## 端到端数据流

```text
文档写入
  -> 文本清洗
  -> 分块
  -> 向量化
  -> 向量存储

用户提问
  -> 查询改写
  -> 检索 top-k
  -> 回答生成 + 引用
  -> 记忆更新
```

## 模块边界设计原则

1. `RAG 各阶段显式可见`：便于教学与排障。
2. `三后端同契约`：便于公平对比。
3. `默认离线可运行`：优先使用内存存储与本地 embedding。
4. `可替换`：embedding、vector store、generator 都可逐步替换。

## 代码入口

- 前端：[`apps/web/src/App.tsx`](../../apps/web/src/App.tsx)
- Python API：[`services/python-rag/app/main.py`](../../services/python-rag/app/main.py)
- Python 编排：[`services/python-rag/app/rag/service.py`](../../services/python-rag/app/rag/service.py)
- Java API：[`services/java-rag/src/main/java/study/rag/api/RagController.java`](../../services/java-rag/src/main/java/study/rag/api/RagController.java)
- Go API：[`services/go-rag/internal/rag/http.go`](../../services/go-rag/internal/rag/http.go)

## 进一步阅读

- [代码阅读地图](code-reading-map.md)
- [共享 API 契约](api-contract.md)
- [RAG 全流程（结合代码）](../learning/rag-pipeline.md)
