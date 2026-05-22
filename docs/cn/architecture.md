# 整体架构

## 项目定位

File Browser 是一个**单二进制文件的自托管 Web 文件管理器**。它将 Go 后端 API 与 Vue 3 单页面应用（SPA）打包进同一个可执行文件，无需任何额外运行时依赖，直接指向一个目录即可通过浏览器管理文件。

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端语言 | Go 1.25 |
| HTTP 路由 | `gorilla/mux` |
| CLI 框架 | `spf13/cobra` + `spf13/viper` |
| 虚拟文件系统 | `spf13/afero` |
| 元数据数据库 | BoltDB（通过 `storm/v3` ORM） |
| 认证令牌 | JWT HS256（`golang-jwt/jwt/v5`） |
| 前端框架 | Vue 3 + Pinia + Vue Router 5 + TypeScript |
| 构建工具 | Vite 8 |
| 国际化 | vue-i18n |
| 可续传上传 | TUS 协议（`tus-js-client`） |
| 图像处理 | `disintegration/imaging` + 信号量 worker 池 |
| 预览缓存 | `diskcache`（可选），TUS 状态可选 Redis |
| 构建流程 | `Taskfile.yml`（先构建前端，再 `go build` 内嵌） |

## 目录结构

```
filebrowser/
├── main.go              # 入口：调用 cmd.Execute()
├── cmd/                 # CLI 命令（server / config / users / rules / cmds）
├── http/                # HTTP 路由、处理器、中间件
├── auth/                # 认证后端（JSON / Proxy / Hook / NoAuth）
├── users/               # 用户模型、权限、密码
├── settings/            # 应用与服务器配置模型
├── storage/             # 存储层 Facade
│   └── bolt/            # BoltDB 实现（唯一后端）
├── files/               # 文件元数据、MIME、目录列表
├── fileutils/           # 文件复制、目录工具
├── rules/               # 路径访问规则（allow/deny）
├── runner/              # Shell 钩子（before_*/after_* 事件）
├── search/              # 文件名/类型搜索
├── share/               # 分享链接模型与存储
├── img/                 # 缩略图/预览图处理
├── diskcache/           # 磁盘预览缓存
├── errors/              # 公共哨兵错误
├── version/             # 构建时版本信息
├── frontend/            # Vue SPA（构建产物内嵌进 Go 二进制）
│   ├── src/
│   └── assets.go        # //go:embed dist/*
└── docker/              # Docker 镜像、健康检查、默认配置
```

## 启动流程

```
main.go
  └─ cmd.Execute()
       └─ rootCmd.RunE  (cmd/root.go)
            ├─ initViper()          ← 读取配置文件 / 环境变量 / CLI flags
            ├─ withViperAndStore()  ← 打开 BoltDB，构建 Storage facade
            ├─ quickSetup()         ← 首次运行时初始化默认用户与配置
            ├─ getServerSettings()  ← 合并运行时设置
            └─ http.NewHandler()    ← 构建路由，启动 HTTP 服务（含优雅关闭）
```

## 核心设计原则

### 1. 处理器返回状态码模式

所有业务处理器签名为 `func(w, r, *data) (int, error)`，由 `handle()` 统一映射到 HTTP 响应，业务逻辑无需直接写响应体。

### 2. 每用户虚拟沙箱文件系统

每个用户通过 `afero.NewBasePathFs` 将 OS 文件系统限定在其 Scope 目录下，防止路径穿越：

```go
u.Fs = afero.NewBasePathFs(afero.NewOsFs(), scope)
```

### 3. 高阶函数中间件

认证、权限、分享校验通过高阶函数组合（`withUser`、`withAdmin`、`withHashFile` 等），而非中间件链，使每条路由的权限语义一目了然。

### 4. 运行时配置 vs 应用配置分离

- **运行时**（Viper / CLI / 环境变量）：监听地址、端口、TLS、缓存目录、根路径等服务器参数。
- **应用配置**（BoltDB）：用户、规则、钩子命令、品牌、分享链接——可通过 Web UI 或 `filebrowser config` 命令管理。

### 5. 单二进制内嵌前端

生产构建时前端 `dist/` 通过 `//go:embed` 打包进 Go 二进制，部署只需单个可执行文件。
