# 本地开发与联调

本页目标：让你在本地快速跑通“前端 + 任一后端”，并能验证完整 RAG 链路。

## 1. 环境要求

- Python 3.11+
- Node.js 20+
- Java 17+（运行 Java 后端时）
- Go 1.22+（运行 Go 后端时）
- Docker（后续持久化阶段使用）

## 2. 启动 Python 后端（推荐起点）

```powershell
cd services/python-rag
python -m venv .venv
.\.venv\Scripts\Activate.ps1
pip install -r requirements.txt
uvicorn app.main:app --reload --port 8000
```

访问：`http://localhost:8000/docs`

## 3. 启动前端

```powershell
cd apps/web
npm install
npm run dev
```

默认请求地址：`http://localhost:8000`

如需切换：

```powershell
$env:VITE_RAG_API_BASE_URL="http://localhost:8082"
npm run dev
```

## 4. 启动 Java 后端

```powershell
cd services/java-rag
mvn spring-boot:run
```

默认端口：`8082`

## 5. 启动 Go 后端

```powershell
cd services/go-rag
go run ./cmd/server
```

默认端口：`8080`

## 6. 最小验证流程（建议按顺序）

1. `POST /documents` 新增文档。
2. `POST /documents/{id}/index` 建索引。
3. `POST /retrieval/search` 看检索片段。
4. `POST /chat/sessions` 创建会话。
5. `POST /chat/sessions/{id}/messages` 看回答与引用。
6. `GET /memory/users/{user_id}` 看记忆更新。

可直接在前端 Workbench 完成以上步骤。

## 7. 启动基础设施（可选）

```powershell
docker compose -f infra/docker-compose.yml up -d
```

说明：当前里程碑以内存实现为主，不依赖数据库也可运行。

## 8. 常见问题

### 8.1 前端报跨域

确认后端已启用 CORS（当前三后端默认允许 `*`）。

### 8.2 检索结果为空

通常是文档未执行 index；先调用 `/documents/{id}/index`。

### 8.3 字段不一致导致前端展示异常

对照 [共享 API 契约](../architecture/api-contract.md) 检查字段名。

## 关联文档

- [共享 API 契约](../architecture/api-contract.md)
- [代码阅读地图](../architecture/code-reading-map.md)
- [RAG 全流程（结合代码）](../learning/rag-pipeline.md)
