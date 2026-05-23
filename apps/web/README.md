# 前端工作台

`apps/web` 是 RAG-study 的 React 工作台，用于把后端能力可视化。

## 这个目录做什么

- 写入文档
- 触发索引
- 执行检索
- 发起聊天
- 查看用户记忆与引用片段

## 关键入口

- 应用入口：[`src/App.tsx`](src/App.tsx)
- 挂载入口：[`src/main.tsx`](src/main.tsx)
- 页面样式：[`src/styles.css`](src/styles.css)

## 运行

```powershell
cd apps/web
npm install
npm run dev
```

默认后端地址：`http://localhost:8000`

切换后端示例：

```powershell
$env:VITE_RAG_API_BASE_URL="http://localhost:8082"
npm run dev
```

## 构建

```powershell
npm run build
```

## 相关文档

- [共享 API 契约](../../docs/architecture/api-contract.md)
- [代码阅读地图](../../docs/architecture/code-reading-map.md)
- [本地开发与联调](../../docs/runbooks/local-development.md)

