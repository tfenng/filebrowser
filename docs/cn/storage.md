# 存储层

## 概述

存储层分为两部分：**元数据存储**（BoltDB）和**文件内容存储**（操作系统文件系统）。两者都通过接口抽象，理论上可替换后端实现。

---

## 元数据存储（BoltDB + Storm）

### 数据库文件

默认路径：`./filebrowser.db`，通过 `--database` / `-d` 或环境变量 `FB_DATABASE` 指定。

底层使用 [BoltDB](https://github.com/etcd-io/bbolt)（嵌入式 KV 存储）配合 [Storm v3](https://github.com/asdine/storm) ORM，Storm 提供结构体到 Bucket 的自动映射。

### Storage Facade

```go
// storage/storage.go
type Storage struct {
    Users    *users.Storage
    Share    *share.Storage
    Auth     *auth.Storage
    Settings *settings.Storage
}
```

`storage/bolt/bolt.go` 中的 `NewStorage(db *storm.DB)` 将四个 Storm 后端注入各 Storage 对象：

| Store | 后端文件 | 存储内容 |
|-------|----------|----------|
| `Users` | `storage/bolt/users.go` | 用户列表（usersBackend） |
| `Share` | `storage/bolt/share.go` | 分享链接（shareBackend） |
| `Auth` | `storage/bolt/auth.go` | 认证方法配置（authBackend） |
| `Settings` | `storage/bolt/config.go` | 应用与服务器配置（settingsBackend） |

### 存储后端接口

每个领域都定义了 `StorageBackend` 接口，便于未来替换实现：

```go
// 示例：users/storage.go
type StorageBackend interface {
    GetAll() ([]*User, error)
    Get(id uint) (*User, error)
    GetByUsername(username string) (*User, error)
    Save(u *User) error
    Update(u *User, fields ...string) error
    Delete(id uint) error
}
```

数据库版本通过 `save(db, "version", 2)` 存储，可用于迁移。

---

## 文件系统存储

### 根目录

服务器全局根目录通过 `settings.Server.Root` 配置（`--root` / `FB_ROOT`），所有用户文件都在此目录之下。

### 用户沙箱

每个用户拥有独立的 Scope（相对于 Root 的子路径），通过 `afero.NewBasePathFs` 实现文件系统隔离：

```go
// users/users.go
u.Fs = afero.NewBasePathFs(afero.NewOsFs(), scope)
```

所有文件操作通过 `user.Fs` 进行，无法访问 Scope 之外的路径。钩子命令需要真实路径时使用 `u.FullPath(path)` 获取绝对路径。

### 文件元数据（`files` 包）

`files.NewFileInfo` 负责构建文件元信息：

```go
// files/file.go
type FileInfo struct {
    Path      string
    Name      string
    Size      int64
    Extension string
    ModTime   time.Time
    Mode      fs.FileMode
    IsDir     bool
    IsSymlink bool
    Type      string       // MIME 大类：image/video/audio/text/blob
    Subtitles []string
    Content   string       // 文本文件内容（编辑器使用）
    Listing   *Listing     // 目录时填充
    Checksums map[string]string
}
```

MIME 检测优先使用文件扩展名映射（`files/mime.go`），回退到 `http.DetectContentType`。

---

## 可选缓存层

### 预览/缩略图缓存

接口：`fbhttp.FileCache`（`http/preview.go`）

```go
type FileCache interface {
    Store(ctx context.Context, key string, value []byte) error
    Get(ctx context.Context, key string) ([]byte, bool, error)
    Delete(ctx context.Context, key string) error
}
```

- **启用：** `--cacheDir <path>` → `diskcache.New(path)`，缩略图以 SHA256(path+params) 为键存储为文件
- **禁用：** `diskcache.NewNoOp()`（默认）

### TUS 上传状态缓存

接口：`fbhttp.UploadCache`

- **内存缓存（默认）：** `newMemoryUploadCache()`，带 TTL 自动清理
- **Redis 缓存：** `--redisCacheUrl redis://...` → `newRedisUploadCache()`，适用于多实例部署

---

## 数据持久化流程示例

**用户登录时：**

```
loginHandler
  → store.Auth.Get(method)          ← 从 BoltDB 读取认证配置
  → auther.Auth(r)                  ← 验证凭据
  → store.Users.GetByUsername()     ← 从 BoltDB 加载用户
  → printToken()                    ← 签发 JWT（用户数据内嵌 Payload）
```

**文件列表请求时：**

```
resourceGetHandler
  → data.Check(path)                ← 从内存规则检查路径权限
  → files.NewFileInfo(d.user.Fs)    ← 通过用户沙箱 FS 读取文件系统
  → renderJSON(w, fi)               ← 序列化返回
```
