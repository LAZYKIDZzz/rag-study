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

- 查询理解：LLM 改写、多查询生成、HyDE、关键词联想。
- 检索增强：关键词检索、向量检索、metadata filter、混合检索、结果融合。
- 精排增强：rerank、去重、相邻 chunk 合并。
- 记忆增强：短期会话、用户 facts、摘要记忆、实体关系记忆、过期与删除策略。
- 生成增强：LLM grounded generation、引用约束、答案自检。
- 评测集与回归机制：召回率、引用准确率、答案忠实度、延迟。

## M5：研究型 RAG 能力（可选）

- RAPTOR / 层级摘要索引，用于长文档总结和跨章节问题。
- GraphRAG，用于实体关系明显、跨文档推理较多的知识库。
- Self-RAG 式检索决策与答案批判，用于更严格的质量控制。
- Agentic retrieval，用于复杂任务分解和多工具检索。

## 关联文档

- [产品范围](product-scope.md)
- [数据模型](../architecture/data-model.md)
- [三后端实现对照](../architecture/backend-comparison.md)
- [进阶 RAG 技术地图](../learning/advanced-rag-patterns.md)
