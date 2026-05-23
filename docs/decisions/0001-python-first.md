# ADR 0001：先完成 Python 后端

## 状态

已采纳（Accepted）

## 背景

项目目标不是“只跑起来一个后端”，而是：

- 完成一套可学习的 RAG 闭环。
- 同时对比 Python、Java、Go 的实现差异。

如果三后端同时推进，早期会把需求不确定性放大为三份重复返工。

## 决策

先把 Python 作为首个完整参考实现（FastAPI + 显式 RAG 模块），
再让 Java、Go 按同一 API 契约对齐。

## 决策理由

- Python 在 RAG 生态里迭代快，验证成本低。
- 先形成“正确的参考行为”更利于跨语言对齐。
- 可把“业务问题”和“语言框架差异”分开处理。

## 影响

正向影响：

- 更快拿到可演示、可教学的端到端链路。
- 文档与测试可以先围绕参考实现稳定下来。

代价：

- Java/Go 在早期阶段会短暂落后于 Python。

## 关联文档

- [产品范围](../requirements/product-scope.md)
- [三后端实现对照](../architecture/backend-comparison.md)
- [共享 API 契约](../architecture/api-contract.md)
