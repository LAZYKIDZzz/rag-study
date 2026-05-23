# Java RAG 服务

`services/java-rag` 是本项目的 Java 对照实现，用来展示同一套 RAG 能力如何在 Spring 风格工程里落地。

## 这个目录做什么

- 实现共享 API 契约
- 对照 Python 的同类模块
- 展示 Java 里的分层、DTO 和服务编排

## 关键入口

- 控制器：[`src/main/java/study/rag/api/RagController.java`](src/main/java/study/rag/api/RagController.java)
- 服务层：[`src/main/java/study/rag/core/RagService.java`](src/main/java/study/rag/core/RagService.java)
- 测试：[`src/test/java/study/rag/core/RagServiceTest.java`](src/test/java/study/rag/core/RagServiceTest.java)

## 运行

```powershell
cd services/java-rag
mvn spring-boot:run
```

默认端口：`8082`

## 验证

```powershell
mvn test
```

## 相关文档

- [共享 API 契约](../../docs/architecture/api-contract.md)
- [三后端实现对照](../../docs/architecture/backend-comparison.md)
- [代码阅读地图](../../docs/architecture/code-reading-map.md)

