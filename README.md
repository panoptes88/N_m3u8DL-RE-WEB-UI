# N_m3u8DL-RE Web UI

基于 Go + Gin + Vue3 + Ant Design Vue 的 Web 界面，用于管理 N_m3u8DL-RE 下载任务。

## 功能特性

- 用户认证（单用户登录、修改密码）
- 创建和管理下载任务，实时查看下载进度与日志
- 下载方案（按域名保存常用配置，新建任务自动匹配）
- 自动选择最佳视频+音频并合并为 mp4、并行下载、跳过分片完整性检测
- 解密支持（mp4decrypt）、自定义请求头、代理、自定义参数
- 文件浏览、下载与管理
- URL / 输出文件名点击复制
- Docker 容器化部署

## 推荐配合使用

**[m3u8dl-browser-extension](https://github.com/panoptes88/m3u8dl-browser-extension)** — Chrome 视频嗅探插件

自动嗅探网页中的 m3u8、mp4、mpd 等媒体资源，支持一键发送到本项目进行远程下载，省去手动复制链接的步骤。

> 扩展本身只负责「嗅探 + 转发链接」，不下载视频；实际下载由本项目服务器完成。

## 界面预览

| 首页 | 下载任务 | 文件管理 |
|:---:|:---:|:---:|
| ![首页](https://i.loi.li/j/0/2026/03/69a3d70ba4fa3.webp) | ![下载任务](https://i.loi.li/j/0/2026/06/6a2054ffbc66e.png) | ![文件管理](https://i.loi.li/j/0/2026/03/69a3d80e450e2.webp) |

## 快速开始

### Docker 快速体验

```bash
mkdir /data/m3u8dl -p
cd /data/m3u8dl
docker run -d --name m3u8dl -p 8080:8080 -e ALLOW_INSECURE=true -e ADMIN_PASSWORD=admin123 -v ./db:/app/db -v ./downloads:/app/downloads ghcr.io/panoptes88/n_m3u8dl-re-web-ui:latest
```

> **提示**：国内网络可将 `ghcr.io` 替换为 `ghcr.1ms.run`，避免镜像拉取失败。

访问 http://localhost:8080

### 使用 Docker Compose（推荐）

```bash
mkdir /data/m3u8dl -p
cd /data/m3u8dl
curl https://raw.githubusercontent.com/panoptes88/N_m3u8DL-RE-WEB-UI/refs/heads/main/docker-compose-example.yml -o docker-compose.yml
docker compose up -d
```

访问 http://localhost:8080

默认登录信息：
- 用户名: admin
- 密码: admin123

## 目录结构

```
N_m3u8DL-RE_WEB_UI/
├── cmd/
│   └── server/           # Go 主程序入口
├── internal/
│   ├── config/           # 配置
│   ├── handler/          # HTTP 处理器
│   ├── middleware/       # 中间件
│   ├── model/            # 数据模型
│   └── service/          # 业务逻辑
├── web/                  # Vue 前端
│   ├── src/
│   │   ├── views/        # 页面组件
│   │   ├── stores/       # Pinia 状态管理
│   │   └── api/          # API 调用
│   └── package.json
├── db/                   # 数据库目录（运行时生成）
├── downloads/            # 下载文件目录（运行时生成）
├── bin/                  # 二进制工具目录（仅源码部署需手动准备；Docker 部署自动下载）
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## 配置项

通过环境变量配置：

| | 变量 | 默认值 | 说明 |
|-:|------|--------|------|
| 🔴 | PORT | 8080 | 服务端口 |
| 🔴 | DOWNLOAD_DIR | ./downloads | 下载文件目录 |
| 🔴 | BIN_DIR | ./bin | 工具目录 |
| 🔴 | DB_PATH | ./db/data.db | 数据库文件路径 |
| 🔴 | LANG | C.UTF-8 | 容器字符集，影响终端内中文文件名显示（镜像已内置） |
| 🟢 | ADMIN_PASSWORD | admin123 | 管理员密码 |
| 🟢 | TZ | Asia/Shanghai | 时区设置 |
| 🟢 | ALLOW_INSECURE | false | 是否允许非 HTTPS 环境，如需 HTTP 访问需设为 true |
| 🟢 | ALLOW_ORIGINS | http://localhost:8080,http://127.0.0.1:8080 | 允许的跨域来源，多个用逗号分隔 |
| 🟢 | DOWNLOAD_TIMEOUT | 0 | 下载超时时间（秒），0 表示不限制 |

> **说明**：🔴 建议保持默认，🟢 推荐按需修改。

### 详细说明

#### ALLOW_INSECURE
控制 Cookie 的 Secure 属性：
- `false`：Cookie 设置 Secure 标志，仅 HTTPS 传输，适用于生产环境（默认）
- `true`：Cookie 不设置 Secure 标志，适用于 HTTP 环境

#### ALLOW_ORIGINS
配置允许跨域访问的来源地址，防止 CSRF 攻击。多个地址用逗号分隔。

#### DOWNLOAD_TIMEOUT
设置下载任务的最大超时时间：
- `0` 或负数：不限制超时
- 正整数：超时时间（秒）

#### TZ
容器内部时区，影响日志时间显示。推荐使用 Asia/Shanghai。

#### LANG
容器字符集。设为 `C.UTF-8` 可使容器终端（如 `docker exec` 进容器 `ls`）正确显示中文文件名，避免被转义成 `$'\xxx'`。镜像已内置默认值，通常无需修改。

## API 接口

所有接口在 `/api` 下。除认证接口外，均需登录后携带 Cookie（session）访问。

### 认证
- `POST /api/auth/login` - 登录
- `POST /api/auth/logout` - 登出
- `POST /api/auth/change-password` - 修改密码
- `GET /api/user` - 获取当前用户

### 任务管理
- `GET /api/tasks` - 获取任务列表
- `POST /api/tasks` - 创建任务
- `GET /api/tasks/:id` - 获取任务详情
- `DELETE /api/tasks/:id` - 删除任务
- `GET /api/tasks/:id/log` - 获取任务日志
- `POST /api/tasks/:id/save-as-profile` - 保存任务为下载方案

### 文件管理
- `GET /api/files` - 获取文件列表
- `GET /api/files/download` - 下载文件
- `DELETE /api/files/:name` - 删除文件

### 下载方案
- `GET /api/profiles` - 方案列表
- `POST /api/profiles` - 创建方案
- `GET /api/profiles/:id` - 获取方案
- `PUT /api/profiles/:id` - 更新方案
- `DELETE /api/profiles/:id` - 删除方案
- `GET /api/profiles/by-domain` - 按域名查询方案

### 创建任务示例

```bash
# 1. 登录，保存 Cookie
curl -c cookie.txt -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 2. 创建任务（携带 Cookie）
curl -b cookie.txt -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com/index.m3u8","output_name":"my-video","auto_select":true}'
```

主要请求参数：

| 字段 | 类型 | 说明 |
|------|------|------|
| url | string | m3u8 链接（必填） |
| output_name | string | 输出文件名 |
| auto_select | bool | 自动选最佳视频+音频并合并为 mp4，**建议传 true** |
| thread_count | int | 下载线程数，默认 32 |
| retry_count | int | 重试次数，默认 15 |
| headers | string | 自定义请求头，分号分隔 |
| base_url | string | Base URL |
| custom_proxy | string | 自定义代理 |
| custom_args | string | 附加给 N_m3u8DL-RE 的自定义参数 |
| concurrent_download | bool | 并行下载音视频 |
| skip_segments_check | bool | 跳过分片数量完整性检测 |
| key | string | 解密密钥 |
| decryption_engine | string | 解密引擎，默认 MP4DECRYPT |

> ⚠️ **重要**：API 创建任务时建议传 `"auto_select": true`。不传（默认 false）会导致多码率 m3u8 下载失败、且下载后视频音频不合并（无声音）。Web UI 已默认勾选，API 调用需自行带上。

## 二进制文件

bin/ 目录需要准备以下文件（**Docker 部署会自动下载，无需手动准备**；仅源码部署需自行下载放到 bin/）：

| 文件 | 说明 | 获取方式 |
|------|------|----------|
| N_m3u8DL-RE | m3u8 下载器 | [Release 页面](https://github.com/nilaoda/N_m3u8DL-RE/releases) |
| ffmpeg | 视频处理工具 | [BtbN Builds](https://github.com/BtbN/FFmpeg-Builds/releases) |
| mp4decrypt | MP4 解密工具 | [Bento4](https://www.bok.net/Bento4/binaries/) |

> 注意：bin/ 目录已添加到 .gitignore，不会被提交到版本控制。

## 从源码构建

如需自行构建镜像（定制或开发）：

```bash
git clone https://github.com/panoptes88/N_m3u8DL-RE-WEB-UI.git
cd N_m3u8DL-RE-WEB-UI
docker build -t m3u8dl .
docker run -d --name m3u8dl -p 8080:8080 -v ./db:/app/db -v ./downloads:/app/downloads m3u8dl
```
