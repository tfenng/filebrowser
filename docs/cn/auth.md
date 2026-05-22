# 认证与授权模块

## 概述

认证系统由 `auth/` 包实现，通过 `auth.Auther` 接口支持四种可插拔的认证后端，认证通过后颁发 JWT 令牌用于后续请求鉴权。

## 认证接口

```go
// auth/auth.go
type Auther interface {
    Auth(r *http.Request, userStore *users.Storage, settings *settings.Settings, tokenKey []byte) (*users.User, error)
    LoginPage() bool
}
```

## 四种认证后端

### 1. JSON 认证（默认）

**文件：** `auth/json.go`  
**常量：** `MethodJSONAuth`

用户通过 POST `/api/login` 提交用户名/密码 JSON，后端以 bcrypt 校验密码。

特性：
- 支持 Google reCAPTCHA v2/v3（配置 `recaptchaKey`/`recaptchaHost`）
- 登录失败时执行虚假 bcrypt 比较，防止时序攻击
- 支持用户注册（当 `settings.Signup = true`）

### 2. 代理认证

**文件：** `auth/proxy.go`  
**常量：** `MethodProxyAuth`

从指定 HTTP 请求头（如 `X-Forwarded-User`）读取用户名，适合前置反向代理（Nginx、Traefik、Authentik 等）已完成认证的场景。

特性：
- 若数据库中不存在该用户，可自动创建（使用默认权限）
- 不显示登录页（`LoginPage()` 返回 `false`）

### 3. 钩子认证

**文件：** `auth/hook.go`  
**常量：** `MethodHookAuth`

将用户名/密码通过环境变量传递给外部命令，由外部命令输出 `auth`、`pass` 或 `block` 决定认证结果。

环境变量：
- `USERNAME` — 提交的用户名
- `PASSWORD` — 提交的密码

输出协议：
- `auth` — 认证通过，允许访问
- `pass` — 交由数据库密码校验
- `block` — 明确拒绝

### 4. 无认证模式

**文件：** `auth/none.go`  
**常量：** `MethodNoAuth`

始终以用户 ID 1（管理员）身份访问，用于单用户完全信任环境（如本地开发）。

---

## JWT 会话流程

```
POST /api/login
  └─ loginHandler (http/auth.go)
       ├─ d.store.Auth.Get()          ← 取出配置的 Auther
       ├─ auther.Auth(r, ...)          ← 验证凭据，返回 User
       └─ printToken(w, r, d, user)   ← 签发 JWT

每次 API 请求
  └─ withUser (http/auth.go)
       ├─ extractor.ExtractToken()    ← 从 X-Auth 头或 auth Cookie 提取令牌
       ├─ jwt.ParseWithClaims()       ← 验证签名与有效期
       ├─ d.store.Users.Get()         ← 从数据库重新加载用户（确保最新权限）
       └─ 若令牌临近过期或用户已更新 → 设置 X-Renew-Token 响应头

PUT /api/renew
  └─ renewHandler                    ← 重新签发令牌
```

JWT Payload 中内嵌了用户权限、locale、允许命令等信息（`userInfo` 结构体），减少数据库查询。

---

## 授权分层

| 层级 | 实现位置 | 说明 |
|------|----------|------|
| JWT 校验 | `withUser` | 所有 `/api/*` 路由（公开分享除外） |
| 管理员权限 | `withAdmin` | 用户管理、全局设置接口 |
| 自身或管理员 | `withSelfOrAdmin` | 用户个人资料 GET/PUT/DELETE |
| 用户权限位 | `users.Permissions` | 在各处理器中检查 `Create/Delete/Share/Download/Rename/Execute` |
| 路径规则 | `data.Check(path)` | 全局规则 + 用户规则按最长匹配合并，allow/deny |
| 分享权限 | `withPermShare` | 需同时具备 `Perm.Share` 和 `Perm.Download` |
| 命令执行权限 | `Perm.Execute` + 白名单 | `EnableExec` 服务器开关 + 用户命令白名单 |

---

## 路径访问规则

规则在 `rules/` 包中实现，`data.Check(path)` 将全局规则与用户规则合并后按如下逻辑匹配：

- 支持精确路径和正则表达式两种模式
- 多条规则按**最后匹配优先**（即更具体的规则覆盖通用规则）
- 默认隐藏点文件（`.` 开头），可通过 `HideDotfiles` 配置关闭

---

## 安全注意事项

- 默认关闭 Shell 命令执行（`disableExec = true`）
- 注册新用户不获得管理员或执行权限
- 密码使用 bcrypt 存储
- 响应头包含 CSP（Content Security Policy）
- 分享链接支持可选密码保护和过期时间
