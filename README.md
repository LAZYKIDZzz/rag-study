# RAG-study

`RAG-study` 是一个学习型 RAG 项目，目标是让你通过文档和代码同时理解：

- RAG 是什么
- RAG 在工程里怎么落地
- Python / Java / Go 三种后端怎么做出同一套能力

## 目录入口

- [文档中心](docs/README.md)
- [前端工作台](apps/web/README.md)
- [Python 后端](services/python-rag/README.md)
- [Java 后端](services/java-rag/README.md)
- [Go 后端](services/go-rag/README.md)

## 建议阅读顺序

1. [文档中心](docs/README.md)
2. [RAG 从 0 到 1](docs/learning/rag-basics.md)
3. [代码阅读地图](docs/architecture/code-reading-map.md)
4. [共享 API 契约](docs/architecture/api-contract.md)

## 快速启动

### Python 后端

```powershell
cd services/python-rag
python -m venv .venv
.\.venv\Scripts\Activate.ps1
pip install -r requirements.txt
uvicorn app.main:app --reload --port 8000
```

### 前端工作台

```powershell
cd apps/web
npm install
npm run dev
```

### 本地基础设施

```powershell
docker compose -f infra/docker-compose.yml up -d
```

## 仓库结构

```text
apps/web/            前端工作台
services/python-rag/ Python 参考实现
services/java-rag/   Java 对照实现
services/go-rag/     Go 对照实现
docs/                中文文档中心
infra/               本地基础设施
```

## 进一步阅读

- [产品范围](docs/requirements/product-scope.md)
- [系统总览](docs/architecture/system-overview.md)
- [三后端实现对照](docs/architecture/backend-comparison.md)
- [本地开发与联调](docs/runbooks/local-development.md)

