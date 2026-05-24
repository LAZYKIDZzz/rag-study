# GitHub Pages 文档部署

本文记录如何把 `docs/` 下的静态文档服务部署到 GitHub Pages。

当前文档站是纯静态结构：

- `docs/index.html` 是浏览器入口。
- `docs/**/*.md` 是被入口页动态加载的 Markdown 文档。
- 不需要 `npm build`、后端服务或额外构建产物。

## 推荐方案：从 `/docs` 目录发布

适合当前仓库，因为 GitHub Pages 原生支持从指定分支的 `/docs` 目录发布静态站点。

### 1. 提交并推送文档

确保这些文件已经提交到 GitHub：

- `docs/index.html`
- `docs/README.md`
- `docs/learning/`
- `docs/architecture/`
- `docs/requirements/`
- `docs/runbooks/`
- `docs/decisions/`

如果部署分支是 `main`，需要把文档合并或提交到 `main`；如果暂时使用 `dev` 分支，也可以在 Pages 设置中选择 `dev`。

### 2. 配置 GitHub Pages

在 GitHub 仓库页面进入：

```text
Settings -> Pages
```

然后配置：

```text
Source: Deploy from a branch
Branch: main 或 dev
Folder: /docs
```

保存后，GitHub 会自动发布 `docs/index.html`。

### 3. 访问地址

仓库地址是：

```text
git@github.com:LAZYKIDZzz/rag-study.git
```

默认 GitHub Pages 地址通常是：

```text
https://lazykidzzz.github.io/rag-study/
```

如果仓库 Pages 设置页显示了不同地址，以设置页中的地址为准。

## 验证部署

部署完成后检查：

1. 打开 `https://lazykidzzz.github.io/rag-study/`，确认能看到 RAG-study 文档阅读器。
2. 点击左侧导航中的几篇 Markdown 文档，确认正文可以加载。
3. 点击“打开原始 Markdown”，确认能打开当前文档对应的 `.md` 文件。
4. 刷新页面后确认当前 hash 路由仍能恢复，例如：

```text
https://lazykidzzz.github.io/rag-study/#learning/rag-basics.md
```

## 本地预览

不要直接双击打开 `index.html` 作为最终验证，因为浏览器通常会限制从本地文件读取 Markdown。

在项目根目录运行：

```powershell
python -m http.server 9000
```

然后访问：

```text
http://localhost:9000/docs/
```

## 可选方案：GitHub Actions 发布

如果后续文档站需要构建步骤，例如引入 Vite、Docusaurus、MkDocs 或静态资源生成流程，再改用 GitHub Actions 更合适。

可以新增 `.github/workflows/pages.yml`，把 `docs/` 上传为 Pages artifact。当前不需要这样做，因为 `/docs` 分支发布已经能覆盖需求。

## 常见问题

### 页面打开了，但 Markdown 加载失败

优先检查：

- Pages 是否选择了 `Folder: /docs`。
- 访问地址是否是仓库 Pages 根地址，而不是 `/docs/` 子路径。
- Markdown 文件是否已经提交并推送到部署分支。
- 浏览器控制台是否出现 404。

### 应该选择 `main` 还是 `dev`

推荐长期使用 `main` 发布稳定文档。

如果项目当前主要工作都在 `dev`，可以先选择 `dev` 快速预览；等文档稳定后，再合并到 `main` 并把 Pages 分支切回 `main`。

### 是否需要配置自定义域名

当前不需要。默认地址 `https://lazykidzzz.github.io/rag-study/` 已经可以公开访问。

如果后续需要自定义域名，再在 `Settings -> Pages -> Custom domain` 中配置，并按 GitHub 提示添加 DNS 记录。

## 官方参考

- GitHub Pages 发布源配置：<https://docs.github.com/en/pages/getting-started-with-github-pages/configuring-a-publishing-source-for-your-github-pages-site>
- GitHub Pages 自定义工作流：<https://docs.github.com/en/pages/getting-started-with-github-pages/using-custom-workflows-with-github-pages>
