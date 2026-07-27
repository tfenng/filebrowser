<p align="center">
  <img src="https://raw.githubusercontent.com/filebrowser/filebrowser/master/branding/banner.png" width="550" alt="File Browser"/>
</p>

<p align="center">
  <a href="README.md">English</a> · 简体中文
</p>

[![Build](https://github.com/filebrowser/filebrowser/actions/workflows/ci.yaml/badge.svg)](https://github.com/filebrowser/filebrowser/actions/workflows/ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/filebrowser/filebrowser/v2)](https://goreportcard.com/report/github.com/filebrowser/filebrowser/v2)
[![Version](https://img.shields.io/github/release/filebrowser/filebrowser.svg)](https://github.com/filebrowser/filebrowser/releases/latest)

# File Browser

File Browser 是一个自托管的 Web 文件管理器。将它部署到服务器并指定一个目录，即可通过浏览器上传、下载、预览、编辑和整理文件。项目由一个 Go 后端和 Vue 3 前端组成，发布时可打包为单个可执行文件，适合搭建轻量的个人文件服务。

## 主要功能

- 文件与目录的创建、上传、下载、复制、移动、重命名和删除
- 图片、音视频、PDF、文本及代码文件预览
- 多用户、目录范围、访问规则和细粒度权限
- 公开分享链接，可设置密码和过期时间
- 文件搜索、分片续传、缩略图和在线代码编辑
- 自定义品牌、主题、多语言及可选的命令执行钩子

## 快速开始

### Docker

```bash
docker run -d \
  --name filebrowser \
  -p 8080:80 \
  -v /path/to/files:/srv \
  -v /path/to/database:/database \
  -v /path/to/config:/config \
  filebrowser/filebrowser:latest
```

打开 `http://localhost:8080`。首次启动时，管理员用户名默认为 `admin`，随机生成的密码会打印在容器日志中：

```bash
docker logs filebrowser
```

请在首次登录后立即修改密码。生产环境中应限制配置和数据库文件的访问权限，并使用反向代理启用 HTTPS。

### 本地运行

从 [Releases](https://github.com/filebrowser/filebrowser/releases/latest) 下载对应平台的程序，然后运行：

```bash
./filebrowser --root /path/to/files
```

默认监听 `127.0.0.1:8080`。完整的安装和配置说明见 [filebrowser.org](https://filebrowser.org)。

## 构建与开发

开发环境需要 Go 1.26、Node.js 24+、pnpm 10+，并建议安装 [Task](https://taskfile.dev/)。

```bash
# 构建前端和后端，输出 ./filebrowser
task build

# 仅运行 Go 后端
go mod download
go run .

# 启动前端开发服务器
cd frontend
pnpm install
pnpm run dev
```

首次使用新的 `filebrowser.db` 启动时，程序会创建 `admin` 用户，并在终端输出一次随机密码：

```text
User 'admin' initialized with randomly generated password: ...
```

请保存该密码并在首次登录后修改。若启动日志已丢失，请先停止正在运行的服务（前台运行时按 `Ctrl+C`），再重置密码并重新启动：

```bash
./filebrowser users update admin --password 'your-new-strong-password'
./filebrowser --root /home/tony/data
```

管理命令和服务进程不能同时写入同一个 BoltDB 数据库；若未先停服，重置命令可能报告 `timeout`。

## 测试与代码检查

提交变更前，请运行与修改范围相符的检查：

```bash
go test --race ./...

cd frontend
pnpm run lint
pnpm run test
pnpm run typecheck
```

后端测试与源码放在同一包内，文件名为 `*_test.go`；前端使用 Vitest，测试位于相邻的 `__tests__/` 目录中。

## 架构文档

- [整体架构](docs/cn/architecture.md)
- [认证与授权](docs/cn/auth.md)
- [存储层](docs/cn/storage.md)
- [前端架构](docs/cn/frontend.md)
- [HTTP API](docs/cn/api.md)
- [配置系统](docs/cn/config.md)
- [核心功能实现](docs/cn/features.md)

网站文档位于 `www/docs/`。运行 `task docs` 可构建文档，运行 `task docs:serve` 可在本地预览。

## 项目状态

File Browser 已实现其作为单文件 Web 文件管理器的核心目标，目前处于**仅维护模式**：

- Issue 主要用于跟踪缺陷，其他话题请使用 [Discussions](https://github.com/filebrowser/filebrowser/discussions)。
- 维护重点是问题分类、安全修复和缺陷修复。
- 暂无新功能规划，新功能 Pull Request 不保证会被审查。

更多背景请阅读维护者的[项目状态说明](https://hacdias.com/2026/03/11/filebrowser/)。

## 贡献

欢迎参与维护。提交 Pull Request 前请阅读[贡献指南](CONTRIBUTING.md)与[行为准则](CODE-OF-CONDUCT.md)。PR 应面向 `master`，标题遵循 Conventional Commits，例如 `fix(auth): 修复代理登录`。翻译内容请通过 [Transifex](https://app.transifex.com/file-browser/file-browser/) 提交。

安全问题请按照[安全策略](SECURITY.md)私下报告，不要创建公开 Issue。

## 许可证

[Apache License 2.0](LICENSE) © File Browser Contributors
