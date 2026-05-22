# 核心功能实现

## 文件浏览与操作

### 文件列表（`GET /api/resources{path}`）

处理器：`http/resource.go` `resourceGetHandler`

流程：
1. `data.Check(path)` — 校验路径规则（全局 + 用户规则）
2. `files.NewFileInfo(d.user.Fs, path, ...)` — 从用户沙箱 FS 读取元信息
3. 若为目录：填充 `Listing`（支持排序、分页）
4. 若请求 `checksum` 参数：异步计算 MD5/SHA 系列校验和
5. 返回序列化的 `FileInfo` JSON

排序支持：`name`（文件名）、`size`（大小）、`modified`（修改时间），正序/倒序。

### 文件上传

**覆盖上传（`POST /api/resources{path}`）**
- 直接写入 `d.user.Fs`
- 触发 `runner.RunHook("upload", ...)`

**TUS 可续传上传（`/api/tus`）**
- 遵循 TUS 1.0 协议，支持断点续传
- 上传状态存入 `UploadCache`（内存或 Redis）
- 上传完成后移动临时文件到目标路径
- 支持通过 `--redisCacheUrl` 配置 Redis，实现多实例间状态共享

### 文件复制/移动（`PATCH /api/resources{path}`）

通过 `action` 参数区分操作类型：

```
PATCH /api/resources/old/path?action=rename&destination=/new/path
PATCH /api/resources/src/path?action=copy&destination=/dst/path
```

底层使用 `fileutils` 包的递归复制和 `os.Rename`（同卷移动）。

---

## 文件预览

### 图片预览（`GET /api/preview/{size}/{path}`）

处理器：`http/preview.go`

流程：
1. 检查磁盘缓存（若启用）
2. 通过 `img.Service.Resize()` 处理图片
   - 信号量限制并发 worker 数量（`marusama/semaphore`）
   - 支持 EXIF 自动旋转
   - 输出 JPEG（质量 75）
3. 结果写入磁盘缓存
4. 返回图片字节

`img.Service` 支持的格式：JPEG、PNG、GIF、BMP、TIFF、WebP

### 视频/音频预览
- 前端直接使用 `<video>`/`<audio>` HTML5 标签
- 通过 `/api/raw{path}` 流式传输文件内容
- 支持 HTTP Range 请求（断点续播）

### 文本编辑器
- 使用 Ace Editor，支持语法高亮
- `GET /api/resources{path}` 返回文件内容（`FileInfo.Content`）
- 编辑后 `PUT /api/resources{path}` 保存
- 触发 `after_save` 钩子

### 字幕
- `GET /api/subtitle{path}` 查找视频同名的 `.vtt`/`.srt` 文件
- 前端将字幕轨道注入 `<video>` 元素

---

## 搜索

处理器：`http/search.go` → `search.Search()`

实现（`search/search.go`）：
1. 从用户根路径递归 `Walk`
2. 每个文件路径经 `data.Check()` 过滤（遵守访问规则）
3. 按查询解析结果（`search/conditions.go`）匹配：
   - 纯字符串：文件名包含匹配（不区分大小写）
   - `type:image`、`type:video` 等：按 MIME 大类过滤
4. 流式返回结果（`text/event-stream` 或 JSON 数组）

---

## 分享链接

### 创建分享

处理器：`http/share.go` `sharePostHandler`

1. 生成随机哈希作为链接标识
2. 记录分享者用户 ID、路径、创建时间
3. 可选设置：
   - 过期时间（绝对时间戳）
   - 访问密码（bcrypt 存储）
4. 存入 BoltDB `share` Bucket

### 公开访问

处理器：`http/public.go` `withHashFile`

1. 从 URL 取哈希，查询 `store.Share`
2. 检查是否过期
3. 若有密码，验证请求中的 `token` 参数（前端先 POST 密码换取 token）
4. **以分享者身份构造虚拟 FS**（使用分享者的 Scope），访问目标文件
5. 分享链接可访问文件或整个目录（支持 ZIP 打包下载）

---

## 事件钩子（Shell Hooks）

### 触发时机

在以下操作前后可配置钩子命令（通过 `settings.Commands` 映射）：

| 事件 | 时机 |
|------|------|
| `before_save` | 文件保存前 |
| `after_save` | 文件保存后 |
| `before_copy` | 复制操作前 |
| `after_copy` | 复制操作后 |
| `before_rename` | 重命名/移动前 |
| `after_rename` | 重命名/移动后 |
| `before_upload` | 上传前 |
| `after_upload` | 上传后 |
| `before_delete` | 删除前 |
| `after_delete` | 删除后 |

### 实现（`runner/runner.go`）

钩子命令在子进程中执行，传入以下环境变量：

```
FILE           操作的文件路径（绝对路径）
SCOPE          用户根目录
USERNAME       触发操作的用户名
DESTINATION    目标路径（移动/复制操作）
```

命令支持通过 `settings.Shell` 配置解释器（如 `["/bin/sh", "-c"]`），也可直接作为 `exec` 调用。

---

## 命令执行（WebSocket）

处理器：`http/commands.go` `commandsHandler`

条件：
1. 服务器开启 `enableExec: true`
2. 用户有 `Perm.Execute` 权限
3. 命令在用户的 `Commands` 白名单中

实现：
- WebSocket 连接后，读取命令名称（不带参数）
- `runner.ParseCommand()` 解析命令（支持 Shell 包装）
- 命令在用户的 Scope 目录下执行
- stdout/stderr 实时流式发送到客户端

---

## 品牌定制

通过 `settings.Branding` 配置：

| 字段 | 说明 |
|------|------|
| `name` | 页面标题和 Logo 文字 |
| `filesPath` | 本地品牌文件目录（放置 `logo.png`、`banner.png`、`custom.css` 等） |
| `theme` | 主题（`light`/`dark`/`auto`） |
| `color` | 主题色（CSS 颜色值） |
| `disableExternal` | 隐藏外部链接（如官网、GitHub） |
| `disableUsedPercentage` | 隐藏磁盘用量百分比 |

自定义 CSS 文件会自动加载，可覆盖任何样式。
