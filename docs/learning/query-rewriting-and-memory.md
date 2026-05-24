# 查询改写与记忆

查询改写和记忆解决的是同一个问题：用户真实想问的内容，往往不完整地出现在当前这句话里。

## 为什么直接检索会失败

多轮对话中，用户经常这样问：

```text
用户：RAG 的分块为什么要有 overlap？
助手：overlap 可以减少语义跨边界导致的召回损失。
用户：那它一般设多少？
```

如果直接拿“那它一般设多少？”去检索，系统不知道“它”指的是 chunk overlap。查询改写要做的事，就是把当前问题补成可检索的问题：

```text
RAG 文档分块中的 chunk overlap 一般设置多少，过大或过小有什么影响？
```

## 查询改写在流程中的位置

```text
用户问题
  -> 读取会话上下文
  -> 读取用户记忆 facts
  -> 生成 rewritten_query
  -> embedding / 关键词检索 / 混合检索
  -> 召回 chunk
```

查询改写发生在检索之前。它的目标不是生成最终答案，而是提高检索命中率。

## 常见查询改写策略

| 策略 | 做法 | 适合场景 |
| --- | --- | --- |
| 上下文补全 | 把上一轮问题、主题、实体补进当前问题 | 多轮追问 |
| 关键词联想 | 补充同义词、错误码、产品名、英文缩写 | 客服、运维、技术文档 |
| 查询拆分 | 把一个复杂问题拆成多个子问题 | 对比、归因、方案设计 |
| 多查询生成 | 从多个角度生成检索查询 | 单次召回不稳定 |
| HyDE | 先生成假设答案，再用假设答案检索 | 问题短、抽象、表达弱 |
| 过滤条件提取 | 从问题中提取时间、业务线、权限范围 | 企业知识库 |

详见：[进阶 RAG 技术地图](advanced-rag-patterns.md)

## 关键词联想示例

关键词联想不是简单扩写句子，而是补充文档中可能出现的检索信号。

```text
用户问题：支付失败怎么查？

联想：
- 支付失败
- payment failed
- 交易失败
- 支付回调
- 网关超时
- 余额不足
- 风控拦截
- error_code
```

如果知识库里有错误码文档、接口文档、工单记录，关键词联想能显著改善召回。

## 记忆是什么

记忆是系统对用户和会话历史的结构化保存。它不是“无限聊天记录”，而是对当前和未来问答有用的信息。

| 记忆类型 | 保存什么 | 在检索中怎么用 |
| --- | --- | --- |
| 短期会话 | 最近几轮问题和回答 | 处理“上一个”“它”“第二点” |
| 用户 facts | 用户角色、偏好、业务背景 | 让查询更贴近用户场景 |
| 摘要记忆 | 长会话压缩总结 | 控制上下文长度 |
| 实体记忆 | 人、产品、项目、指标关系 | 支持跨会话关联 |

## 记忆关联示例

```text
用户 facts：
- 用户负责 B 端 SaaS 客服知识库。
- 用户更关注可解释性和引用来源。

用户问题：
怎么优化检索？

关联后的查询：
B 端 SaaS 客服知识库 RAG 如何优化检索准确率，要求可解释、可引用，包含混合检索、rerank、chunk 参数。
```

记忆关联能把“泛泛的问题”变成“符合用户上下文的问题”。

## 当前项目实现

当前实现是可解释、可测试的确定性策略，不依赖 LLM。

### Python

- QueryProcessor：[`services/python-rag/app/rag/query.py`](../../services/python-rag/app/rag/query.py)
- MemoryStore：[`services/python-rag/app/rag/memory.py`](../../services/python-rag/app/rag/memory.py)

行为：

- 从用户 facts 中取最多 3 条。
- 结合上一轮用户问题。
- 拼接形成 `rewritten_query`。
- 发送消息后更新 `memory_updates`。

### Java

- QueryProcessor：[`services/java-rag/src/main/java/study/rag/core/QueryProcessor.java`](../../services/java-rag/src/main/java/study/rag/core/QueryProcessor.java)
- MemoryStore：[`services/java-rag/src/main/java/study/rag/core/MemoryStore.java`](../../services/java-rag/src/main/java/study/rag/core/MemoryStore.java)

Java 与 Python 的设计目标一致，更强调类型模型和 service 分层。

### Go

- Query：[`services/go-rag/internal/rag/query.go`](../../services/go-rag/internal/rag/query.go)
- Memory：[`services/go-rag/internal/rag/memory.go`](../../services/go-rag/internal/rag/memory.go)

Go 当前策略更简化，主要拼接 facts，后续可补齐上一轮问题、摘要记忆和关键词扩展。

## 生产系统应注意什么

1. 不要把所有历史消息都塞进 prompt。
2. 记忆需要权限隔离，不能把 A 用户 facts 用到 B 用户问题里。
3. 记忆要有删除、过期、纠错机制。
4. LLM 改写要保留原始问题，便于排障。
5. 改写查询、召回 chunk、最终答案都要记录日志。
6. 每种增强策略都要评测，不要只凭感觉上线。

## 相关 API

- `POST /memory/users/{user_id}/facts`
- `GET /memory/users/{user_id}`
- `POST /chat/sessions/{session_id}/messages`

详见：[共享 API 契约](../architecture/api-contract.md)

## 下一步

- 想看全链路：[RAG 全流程图解](rag-pipeline.md)
- 想看现代增强方案：[进阶 RAG 技术地图](advanced-rag-patterns.md)
- 想看代码入口：[代码阅读地图](../architecture/code-reading-map.md)
