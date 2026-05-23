# 里程碑（Milestones）

## M1：Python 端到端基线（已完成）

- 建立仓库结构与文档地图。
- 完成 Python FastAPI RAG 后端。
- 前端工作台支持文档写入、检索、聊天、记忆查看。
- API 契约沉淀并形成跨语言对照基础。

## M2：持久化与向量数据库

- 将文档、chunk、session、memory 持久化到 PostgreSQL。
- 接入 `pgvector` 支持持久化向量检索。
- 保留内存实现用于测试与离线演示。

## M3：Java / Go 完整对齐与深度对比

- 补齐 Java / Go 与 Python 的能力一致性。
- 强化异常、测试、配置管理的横向对比。
- 前端提供更清晰的后端切换与差异观测。

## M4：高级 RAG 能力

- 更强查询改写策略（含 LLM 改写）。
- 混合检索与 rerank。
- 记忆摘要策略优化。
- 评测集与回归机制。

## 关联文档

- [产品范围](product-scope.md)
- [数据模型](../architecture/data-model.md)
- [三后端实现对照](../architecture/backend-comparison.md)
