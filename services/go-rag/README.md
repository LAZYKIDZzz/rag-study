# Go RAG 服务

`services/go-rag` 是本项目的 Go 对照实现，使用标准库把 RAG 的每一步显式展开，便于学习和比较。

## 这个目录做什么

- 实现共享 API 契约
- 展示 Go 的显式路由和服务编排
- 作为 Python / Java 的实现对照

## 关键入口

- HTTP 入口：[`internal/rag/http.go`](internal/rag/http.go)
- 服务编排：[`internal/rag/service.go`](internal/rag/service.go)
- 类型定义：[`internal/rag/types.go`](internal/rag/types.go)
- 测试：[`internal/rag/pipeline_test.go`](internal/rag/pipeline_test.go)

## 运行

```powershell
cd services/go-rag
go run ./cmd/server
```

默认端口：`8080`

可用 `PORT` 覆盖端口：

```powershell
$env:PORT = "8081"
go run ./cmd/server
```

## 验证

```powershell
cd services/go-rag
go test ./...
```

## 示例流程

1. `POST /documents`
2. `POST /documents/{id}/index`
3. `POST /retrieval/search`

## 相关文档

- [共享 API 契约](../../docs/architecture/api-contract.md)
- [三后端实现对照](../../docs/architecture/backend-comparison.md)
- [代码阅读地图](../../docs/architecture/code-reading-map.md)

