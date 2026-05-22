# 前端架构

## 技术栈

| 技术 | 用途 |
|------|------|
| Vue 3 | UI 框架（Composition API） |
| TypeScript | 类型安全 |
| Vite 8 | 构建工具 |
| Pinia | 状态管理 |
| Vue Router 5 | 客户端路由 |
| vue-i18n | 国际化 |
| Ace Editor | 文本文件编辑器 |
| tus-js-client | TUS 可续传上传 |

构建产物 `frontend/dist/` 通过 Go `//go:embed` 打包进二进制（`frontend/assets.go`）。

---

## 目录结构

```
frontend/src/
├── main.ts              # 应用入口：挂载 Pinia / Router / i18n
├── router/
│   └── index.ts         # 路由定义与导航守卫
├── stores/              # Pinia 状态模块
│   ├── auth.ts          # 用户与 JWT
│   ├── file.ts          # 文件列表与选中状态
│   ├── upload.ts        # 上传队列
│   ├── clipboard.ts     # 剪切板（复制/移动）
│   ├── layout.ts        # 视图布局（列表/网格）
│   └── router.ts        # 路由辅助状态
├── views/               # 页面级组件
│   ├── Files.vue        # 主文件浏览视图
│   ├── FileListing.vue  # 目录内容展示
│   ├── Editor.vue       # 文本编辑器（Ace）
│   ├── Preview.vue      # 文件预览
│   ├── Share.vue        # 公开分享页
│   └── settings/        # 设置页（个人、全局、用户管理、分享管理）
├── api/                 # API 请求封装
│   ├── files.ts         # 文件 CRUD
│   ├── users.ts         # 用户管理
│   ├── share.ts         # 分享链接
│   └── settings.ts      # 设置
├── utils/
│   ├── auth.ts          # 登录/续签/Token 解析
│   └── ...
└── components/          # 可复用 UI 组件
```

---

## 路由设计

```
/login                   → Login.vue（无需认证）
/files/:path*            → Files.vue（需要认证）
/share/:hash/:path*      → Share.vue（公开，无需认证）
/settings/profile        → 个人设置
/settings/shares         → 分享管理
/settings/global         → 全局配置（管理员）
/settings/users          → 用户管理（管理员）
/settings/users/:id      → 用户编辑（管理员）
*                        → 重定向到 /files/
```

**导航守卫（`beforeResolve`）：**
1. 调用 `initAuth()` 初始化认证状态（首次加载时从 Cookie/localStorage 恢复）
2. 未认证用户重定向到 `/login`
3. 非管理员访问管理员路由时重定向
4. 检测 JWT 临近过期时自动调用 `renew()` 续签

---

## 认证客户端（`utils/auth.ts`）

```
login(username, password, reCaptcha?)
  → POST /api/login
  → 解析 JWT，存入 Cookie 和 localStorage
  → 更新 useAuthStore

parseToken(token)
  → 从 JWT Payload 还原用户信息（权限、locale、命令白名单等）

renew()
  → PUT /api/renew
  → 更新 Token

validateLogin()
  → 检测空闲超时，自动登出
```

Token 的存储与传递：
- 优先存入 Cookie（`auth`），便于同域请求自动携带
- 同时写入 `localStorage` 作为备份
- API 请求通过 `X-Auth` 请求头显式传递（`api/` 模块统一封装）

---

## 状态管理（Pinia）

### `auth` store
- `user`: 当前用户对象（含权限位）
- `jwt`: 当前 JWT 字符串

### `file` store
- `req`: 当前目录的 `FileInfo`（含 `Listing`）
- `selected`: 选中文件路径数组
- `reload()`: 触发重新请求当前目录

### `upload` store
- 维护上传队列，对接 `tus-js-client`
- 上传完成后触发 `file.reload()`

### `clipboard` store
- 存储剪切板操作类型（copy/cut）和来源路径列表
- 粘贴时调用 `PATCH /api/resources`（action: copy/rename）

---

## 服务端注入数据（Bootstrap Data）

服务器在返回 `index.html` 时（`http/static.go` `handleWithStaticData`）将以下数据注入 HTML 模板：

```json
{
  "baseURL": "/",
  "staticURL": "/static",
  "signup": false,
  "version": "2.x.x",
  "noAuth": false,
  "authMethod": "json",
  "loginPage": true,
  "recaptcha": { "host": "", "key": "" },
  "tusSettings": { "chunkSize": 10485760, "retryCount": 5 },
  "branding": { "name": "File Browser", "disableExternal": false, ... }
}
```

前端从全局变量 `window.FileBrowser` 读取，无需额外 API 请求即可初始化页面状态。

---

## 构建集成

```bash
# 开发
cd frontend && pnpm run dev

# 生产构建（由 Taskfile 编排）
task build
# 等价于:
# 1. cd frontend && pnpm run build   → 输出到 frontend/dist/
# 2. go build -o filebrowser .      → 内嵌 dist/ 进二进制
```

开发前端时，后端直接用 `go run .` 启动；浏览器访问 Vite dev server，
由 `frontend/vite.config.ts` 将 `/api` 请求代理到默认后端 `127.0.0.1:8080`。
