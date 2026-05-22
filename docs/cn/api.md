# HTTP API 参考

## 路由总览

所有 API 路由在 `http/http.go` 的 `NewHandler` 中注册。基础路径由 `server.BaseURL` 配置（默认 `/`）。

---

## 公开接口（无需认证）

| 方法 | 路径 | 处理器 | 说明 |
|------|------|--------|------|
| GET | `/health` | `healthHandler` | 健康检查，返回 `200 OK` |
| POST | `/api/login` | `loginHandler` | 用户登录，返回 JWT |
| POST | `/api/signup` | `signupHandler` | 用户注册（需开启 `settings.Signup`） |
| GET | `/api/public/share/{hash}` | `publicShareHandler` | 获取分享链接元信息 |
| GET | `/api/public/dl/{hash}` | `publicDlHandler` | 通过分享链接下载文件 |
| GET | `/static/*` | 静态文件 | 内嵌前端资源 |
| GET | `/*` | SPA 回退 | 返回 `index.html`（含注入的 Bootstrap 数据） |

---

## 需要认证的接口

### 认证

| 方法 | 路径 | 说明 |
|------|------|------|
| PUT | `/api/renew` | 续签 JWT 令牌 |

### 文件操作（`http/resource.go`）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/resources{path}` | 获取文件信息或目录列表 |
| POST | `/api/resources{path}` | 创建目录或上传文件（覆盖模式） |
| PUT | `/api/resources{path}` | 更新文件内容（编辑器保存） |
| PATCH | `/api/resources{path}` | 复制/移动/重命名（`action` 参数：`copy`/`rename`） |
| DELETE | `/api/resources{path}` | 删除文件或目录 |

**GET 参数：**
- `checksum=md5|sha1|sha256|sha512` — 计算并返回文件校验和
- `listing=true/false` — 是否展开目录内容（默认 true）
- `sort=name|size|modified` — 目录排序字段
- `order=asc|desc` — 排序方向

**PATCH 参数：**
- `action=copy` — 复制到目标路径（`destination` 参数）
- `action=rename` — 移动/重命名到目标路径

### 原始下载（`http/raw.go`）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/raw{path}` | 直接下载文件（无 `Content-Disposition`） |
| GET | `/api/raw{path}?algo=zip` | 打包目录为 ZIP 下载 |

支持的压缩格式（通过 `algo` 参数）：`zip`、`tar`、`targz`、`tarbz2`、`tarxz`、`tarzst`

### 预览（`http/preview.go`）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/preview/{size}/{path}` | 获取图片缩略图/预览 |

`{size}` 可选值：`thumb`（150px）、`small`（360px）、`medium`（720px）、`big`（1080px）

缓存键为 `SHA256(path + size + modtime)`，若启用 `--cacheDir` 则复用磁盘缓存。

### 可续传上传（TUS，`http/tus_handlers.go`）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/tus` | 初始化上传，返回上传 URL |
| HEAD | `/api/tus/{id}` | 查询上传进度（offset） |
| PATCH | `/api/tus/{id}` | 上传数据块 |
| DELETE | `/api/tus/{id}` | 取消上传 |

遵循 [TUS Protocol 1.0.0](https://tus.io/protocols/resumable-upload.html)。

### 搜索（`http/search.go`）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/search?query=...` | 按文件名/类型搜索 |

查询语法：
- 普通字符串：按文件名包含匹配
- `type:image` / `type:video` / `type:audio` / `type:doc` — 按 MIME 大类过滤
- 多词以空格分隔，AND 逻辑

### 分享（`http/share.go`）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/share{path}` | 获取当前路径的分享链接 |
| POST | `/api/share{path}` | 创建分享链接（可设密码和过期时间） |
| DELETE | `/api/share/{hash}` | 删除分享链接 |
| GET | `/api/shares` | 列出我的所有分享链接 |

POST 请求体：
```json
{
  "expires": "2026-12-31T00:00:00Z",
  "unit": "hours",
  "password": "可选密码"
}
```

### 磁盘使用（`http/resource.go`）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/usage` | 获取根目录所在磁盘的用量信息 |

### 字幕（`http/subtitle.go`）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/subtitle{path}` | 获取视频文件的同名字幕文件 |

### 命令执行（`http/commands.go`）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET (WebSocket) | `/api/command` | 执行用户白名单内的 Shell 命令（需 `enableExec` 且 `Perm.Execute`） |

---

## 管理员接口

### 用户管理（`http/users.go`）

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| GET | `/api/users` | 管理员 | 获取用户列表 |
| POST | `/api/users` | 管理员 | 创建用户 |
| GET | `/api/users/{id}` | 管理员或自身 | 获取用户信息 |
| PUT | `/api/users/{id}` | 管理员或自身 | 更新用户 |
| DELETE | `/api/users/{id}` | 管理员 | 删除用户 |

### 设置（`http/settings.go`）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/settings` | 获取全局设置 |
| PUT | `/api/settings` | 更新全局设置 |

---

## 错误响应格式

业务错误以标准 HTTP 状态码返回，响应体为纯文本错误消息：

```
403 Forbidden
404 Not Found
409 Conflict        ← 文件已存在
500 Internal Server Error
```

---

## 响应头

所有响应携带 CSP 安全头，`X-Renew-Token` 在需要续签令牌时由服务器主动下发。
