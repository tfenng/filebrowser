<p align="center">
  <img src="https://raw.githubusercontent.com/filebrowser/filebrowser/master/branding/banner.png" width="550"/>
</p>

[![Build](https://github.com/filebrowser/filebrowser/actions/workflows/ci.yaml/badge.svg)](https://github.com/filebrowser/filebrowser/actions/workflows/ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/filebrowser/filebrowser/v2)](https://goreportcard.com/report/github.com/filebrowser/filebrowser/v2)
[![Version](https://img.shields.io/github/release/filebrowser/filebrowser.svg)](https://github.com/filebrowser/filebrowser/releases/latest)

**File Browser** 是一个自托管的 Web 文件管理器。将其安装在服务器上，指定一个目录，即可通过简洁的 Web 界面上传、下载、预览、编辑和管理文件。它是一个**搭建私有云**类型的软件——单个可执行文件，开箱即用。

> **项目状态：** 本项目处于 **仅维护模式**，不再计划新功能。开发重心在于修复缺陷和安全问题。

---

## 目录

- [快速开始](#快速开始)
- [功能概览](#功能概览)
- [功能架构](#功能架构)
  - [整体架构](docs/cn/architecture.md)
  - [认证与授权](docs/cn/auth.md)
  - [存储层](docs/cn/storage.md)
  - [前端架构](docs/cn/frontend.md)
  - [HTTP API 参考](docs/cn/api.md)
  - [配置系统](docs/cn/config.md)
  - [核心功能实现](docs/cn/features.md)
- [构建与开发](#构建与开发)
- [贡献](#贡献)
- [许可证](#许可证)

---

## 快速开始

### Docker（推荐）

```bash
docker run -d \
  -v /path/to/files:/srv \
  -v /path/to/filebrowser.db:/database.db \
  -p 8080:8080 \
  --name filebrowser \
  filebrowser/filebrowser:latest
```

### Docker Compose

```bash
# 使用项目自带的 compose.yaml
docker compose up -d
```

### 直接运行二进制

```bash
# 从 Releases 页面下载对应平台的二进制
./filebrowser --root /path/to/files --port 8080
```

首次运行会自动创建数据库并初始化默认管理员账号：
- **用户名：** `admin`
- **密码：** `admin`

> **请在首次登录后立即修改密码！**

### 配置示例

```bash
./filebrowser \
  --root /data/files \
  --database /data/filebrowser.db \
  --port 8080 \
  --address 0.0.0.0 \
  --baseurl /files
```

更多配置选项参见 [配置系统](docs/cn/config.md)。

---

## 功能概览

| 功能 | 说明 |
|------|------|
| **文件浏览** | 列表/网格视图，支持排序（名称/大小/修改时间） |
| **上传下载** | 普通上传 + TUS 可续传上传，支持 ZIP 批量下载 |
| **文件操作** | 新建目录、重命名、复制、移动、删除 |
| **在线编辑** | 内置 Ace Editor，支持代码语法高亮 |
| **文件预览** | 图片（含缩略图）、视频/音频播放、PDF、文本 |
| **多用户** | 独立账号，每用户可设不同根目录和权限 |
| **分享链接** | 生成公开分享链接，支持密码保护和过期时间 |
| **搜索** | 按文件名和文件类型搜索 |
| **磁盘用量** | 显示当前磁盘使用情况 |
| **字幕** | 视频播放时自动加载同名字幕文件 |
| **Shell 钩子** | 文件操作前后触发自定义 Shell 命令 |
| **命令执行** | （可选）允许用户通过 WebSocket 执行白名单命令 |
| **品牌定制** | 自定义名称、Logo、主题色、CSS |
| **多语言** | 内置 i18n 国际化支持 |

---

## 功能架构

### 技术栈概览

```
┌─────────────────────────────────────────────────────┐
│                    浏览器                            │
│         Vue 3 SPA (Pinia + Vue Router + TS)         │
└──────────────────────┬──────────────────────────────┘
                       │ HTTP / WebSocket
┌──────────────────────▼──────────────────────────────┐
│              Go HTTP 服务器（gorilla/mux）            │
│  ┌──────────┐ ┌──────────┐ ┌───────────┐            │
│  │ 认证中间件│ │ 文件处理器│ │ TUS 处理器│            │
│  └──────────┘ └──────────┘ └───────────┘            │
│  ┌──────────┐ ┌──────────┐ ┌───────────┐            │
│  │ 用户管理 │ │ 分享管理 │ │ 搜索引擎  │            │
│  └──────────┘ └──────────┘ └───────────┘            │
└──────┬──────────────────────────────────────────────┘
       │
┌──────▼──────┐    ┌─────────────────────────────────┐
│  BoltDB     │    │  操作系统文件系统（afero 虚拟FS） │
│  用户/配置  │    │  每用户 BasePathFs 沙箱隔离       │
│  规则/分享  │    └─────────────────────────────────┘
└─────────────┘
```

### 模块说明

| 模块 | 路径 | 职责 |
|------|------|------|
| **CLI / 启动** | `cmd/` | Cobra 命令树，服务器启动，配置管理命令 |
| **HTTP 路由层** | `http/` | 路由注册、处理器、认证中间件、TUS 处理 |
| **认证后端** | `auth/` | JSON/Proxy/Hook/NoAuth 四种认证实现 |
| **用户领域** | `users/` | 用户模型、权限位、密码 bcrypt |
| **应用配置** | `settings/` | 全局设置与服务器配置模型 |
| **存储 Facade** | `storage/` | 聚合四个存储后端（用户/分享/认证/设置） |
| **BoltDB 后端** | `storage/bolt/` | Storm ORM 实现，唯一元数据存储后端 |
| **文件元数据** | `files/` | FileInfo 构建、MIME 检测、目录列表、排序 |
| **文件工具** | `fileutils/` | 跨目录复制、递归删除等底层工具 |
| **访问规则** | `rules/` | 路径 allow/deny 规则，最长前缀匹配 |
| **Shell 钩子** | `runner/` | 事件钩子执行，命令解析，环境变量注入 |
| **搜索引擎** | `search/` | 文件系统 Walk + 查询条件匹配 |
| **分享链接** | `share/` | 分享模型，哈希生成，密码/过期管理 |
| **图像处理** | `img/` | 缩略图生成，EXIF 旋转，并发 worker 池 |
| **预览缓存** | `diskcache/` | 磁盘缓存，内容哈希为键 |
| **前端** | `frontend/` | Vue 3 SPA，构建产物内嵌进 Go 二进制 |

### 详细文档

- [整体架构](docs/cn/architecture.md) — 技术栈、目录结构、启动流程、核心设计原则
- [认证与授权](docs/cn/auth.md) — 四种认证后端、JWT 流程、授权分层、路径规则
- [存储层](docs/cn/storage.md) — BoltDB 元数据存储、文件系统沙箱、可选缓存层
- [前端架构](docs/cn/frontend.md) — Vue 3 组件结构、路由、Pinia 状态、认证客户端
- [HTTP API 参考](docs/cn/api.md) — 全部接口端点、参数、响应格式
- [配置系统](docs/cn/config.md) — 配置优先级、所有参数说明、CLI 命令、Docker 示例
- [核心功能实现](docs/cn/features.md) — 文件预览、TUS 上传、搜索、分享链接、Shell 钩子实现细节

---

## 构建与开发

### 前置依赖

- Go 1.25+
- Node.js 24+ 和 pnpm 10+（构建前端）
- [Task](https://taskfile.dev/)（构建工具，可选）

### 构建完整二进制

```bash
# 安装 Task
brew install go-task

# 构建（前端 + Go）
task build

# 输出：./filebrowser
```

### 仅构建前端

```bash
cd frontend
pnpm install
pnpm run build   # 输出到 frontend/dist/
```

### 本地开发

```bash
# 终端 1：前端开发服务器（热重载）
cd frontend && pnpm run dev

# 终端 2：后端 API
go run . --root /tmp/files
```

当前本机调试数据库的登录账号：

```text
用户名：admin
密码：adminuser987^
```

开发前端时请从 Vite dev server 的地址访问界面；`frontend/vite.config.ts`
会把 `/api` 请求代理到默认后端 `127.0.0.1:8080`。

### 运行测试

```bash
go test ./...
```

---

## 贡献

欢迎提交 Pull Request。贡献前请阅读 [贡献指南](CONTRIBUTING.md)。

**注意：** 项目处于维护模式，仅接受 bug 修复 PR，新功能 PR 不保证被 Review。

---

## 许可证

[Apache License 2.0](LICENSE) © File Browser Contributors
