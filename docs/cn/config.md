# 配置系统

## 配置优先级

配置来源按以下优先级由高到低：

```
1. CLI 命令行参数      (--port 8080)
2. 环境变量           (FB_PORT=8080)
3. 配置文件           (.filebrowser.json / .filebrowser.yml)
4. 数据库存储         (应用配置，用户/规则/钩子等)
5. 内置默认值
```

> **说明：** 服务器运行时参数（端口、TLS、缓存目录等）使用 Viper 处理（前 3 项）；应用配置（用户、认证方法、规则、钩子、品牌等）存储在 BoltDB 中，通过 Web UI 或 `filebrowser config` 命令管理。

---

## 配置文件位置

Viper 按以下顺序查找配置文件（文件名 `.filebrowser`，支持 `.json`/`.yml`/`.yaml`）：

1. 当前工作目录
2. `$HOME`
3. `/etc/filebrowser/`

---

## 服务器运行时参数

| CLI 参数 | 环境变量 | 默认值 | 说明 |
|----------|----------|--------|------|
| `--port` / `-p` | `FB_PORT` | `8080` | 监听端口 |
| `--address` / `-a` | `FB_ADDRESS` | `""` (所有接口) | 监听地址 |
| `--root` / `-r` | `FB_ROOT` | `.` | 文件根目录 |
| `--database` / `-d` | `FB_DATABASE` | `./filebrowser.db` | 数据库路径 |
| `--config` / `-c` | `FB_CONFIG` | — | 配置文件路径 |
| `--baseurl` / `-b` | `FB_BASEURL` | `""` | 基础 URL 前缀（反代子路径） |
| `--log` / `-l` | `FB_LOG` | `stdout` | 日志输出（`stdout`/`stderr`/文件路径） |
| `--cert` | `FB_CERT` | — | TLS 证书路径 |
| `--key` | `FB_KEY` | — | TLS 密钥路径 |
| `--cacheDir` | `FB_CACHE_DIR` | — | 预览缓存目录（不设则禁用缓存） |
| `--noauth` | `FB_NOAUTH` | `false` | 禁用认证（危险，慎用） |
| `--createUserDir` | `FB_CREATE_USER_DIR` | `false` | 创建用户时自动建立目录 |
| `--redisCacheUrl` | `FB_REDIS_CACHE_URL` | — | Redis 地址，用于多实例 TUS 状态共享 |

---

## 应用配置（数据库存储）

通过 `filebrowser config set` 或 Web UI 管理员界面修改：

### 全局设置（`settings.Settings`）

| 字段 | 说明 |
|------|------|
| `signup` | 是否允许用户自行注册 |
| `authMethod` | 认证方式（`json`/`proxy`/`hook`/`noauth`） |
| `branding` | 品牌定制（名称、目录、禁用外部链接等） |
| `commands` | 事件钩子命令（`before_save`/`after_save` 等） |
| `shell` | Shell 解释器（如 `["/bin/sh", "-c"]`） |
| `rules` | 全局路径访问规则列表 |
| `tus` | TUS 上传配置（块大小、重试次数） |
| `defaults` | 新用户的默认权限与设置 |

### 服务器设置（`settings.Server`）

| 字段 | 说明 |
|------|------|
| `root` | 文件根目录 |
| `baseURL` | 基础 URL |
| `socket` | Unix socket 路径（替代 TCP） |
| `tlsCert` / `tlsKey` | TLS 证书 |
| `port` / `address` | 监听配置 |
| `log` | 日志输出 |
| `enableThumbnails` | 启用图片缩略图生成 |
| `resizePreview` | 压缩预览图以节省带宽 |
| `enableExec` | 允许用户执行 Shell 命令（默认关闭） |
| `typeDetectionByHeader` | 通过文件内容检测类型而非仅靠扩展名 |
| `authHook` | 钩子认证命令 |
| `tokenExpirationTime` | JWT 有效期 |

---

## CLI 配置命令

```bash
# 初始化配置（写入数据库默认值）
filebrowser config init

# 查看当前配置
filebrowser config cat

# 修改配置（示例）
filebrowser config set --signup=true
filebrowser config set --auth.method=proxy --auth.header=X-Forwarded-User
filebrowser config set --branding.name="My Files"
filebrowser config set --shell="/bin/sh -c"

# 导出配置为 JSON
filebrowser config export config.json

# 从 JSON 导入配置
filebrowser config import config.json
```

---

## 用户管理 CLI

```bash
# 列出所有用户
filebrowser users ls

# 添加用户
filebrowser users add username password --scope=/data/username

# 更新用户
filebrowser users update username --password=newpass

# 删除用户
filebrowser users rm username

# 导出/导入用户
filebrowser users export users.json
filebrowser users import users.json
```

---

## 规则管理 CLI

```bash
# 添加全局规则（拒绝访问 .git 目录）
filebrowser rules add --regex --deny --pattern "^\.git"

# 查看规则
filebrowser rules ls
```

---

## 首次运行自动初始化

若数据库文件不存在，`quickSetup()` 会自动执行：

1. 创建默认 `Settings`（signup=false，authMethod=json）
2. 创建默认 `Server` 配置
3. 创建 `admin` 用户（密码 `admin`）
4. **首次登录后请立即修改默认密码！**

---

## 环境变量命名规则

环境变量以 `FB_` 前缀，CLI 参数的连字符转为下划线：

```
--port          → FB_PORT
--base-url      → FB_BASE_URL
--auth.method   → FB_AUTH_METHOD
```

---

## Docker 部署配置示例

```yaml
# compose.yaml
services:
  filebrowser:
    image: filebrowser/filebrowser:latest
    ports:
      - "8080:8080"
    volumes:
      - /data/files:/srv          # 文件根目录
      - /data/fb.db:/database.db  # 数据库
      - /data/fb.json:/.filebrowser.json  # 配置文件
    environment:
      FB_ROOT: /srv
      FB_DATABASE: /database.db
```

最简 `.filebrowser.json`：

```json
{
  "port": 8080,
  "baseURL": "",
  "address": "",
  "log": "stdout",
  "database": "/database.db",
  "root": "/srv"
}
```
