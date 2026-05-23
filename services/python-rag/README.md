# Python RAG 服务

`services/python-rag` 是本项目的 Python 参考实现，也是最先完成的完整 RAG 基线。

## 这个目录做什么

- 提供共享 API 契约的 Python 实现
- 展示完整 RAG 流程编排
- 作为 Java / Go 的行为参考

## 关键入口

- API 入口：[`app/main.py`](app/main.py)
- 编排层：[`app/rag/service.py`](app/rag/service.py)
- 数据模型：[`app/rag/models.py`](app/rag/models.py)
- 测试：[`tests/test_api.py`](tests/test_api.py)

## 运行

```powershell
cd services/python-rag
python -m venv .venv
.\.venv\Scripts\Activate.ps1
pip install -r requirements.txt
uvicorn app.main:app --reload --port 8000
```

## 验证

```powershell
python -m compileall app
pytest
```

## 相关文档

- [共享 API 契约](../../docs/architecture/api-contract.md)
- [RAG 全流程（结合代码）](../../docs/learning/rag-pipeline.md)
- [代码阅读地图](../../docs/architecture/code-reading-map.md)

