# dnshe-go

`dnshe-go` 是一个面向 DNSHE 的轻量 IPv6 DDNS 服务。它直接获取本机公网 IPv6，并调用 DNSHE API 更新 `AAAA` 记录。

项目目标是用一个服务完成 DNSHE IPv6 DDNS 更新，不再依赖额外 callback 服务。

## 与 DDNS-GO 的关系

`dnshe-go` 不是 `ddns-go` 的 fork，也不依赖 `ddns-go`。

这个项目最初来自 `ddns-go + dnshe-ddns-go-callback` 的使用场景：`ddns-go` 负责检测 IP 变化，callback 服务负责把更新请求转换成 DNSHE API 调用。这个方案能工作，但部署时需要同时维护两个服务。

`dnshe-go` 把这段流程合并到一个独立程序里：

- 内置 IPv6 获取。
- 内置 DNSHE API 调用。
- 内置 Web UI。
- 内置自动同步循环。

因此部署时只需要运行 `dnshe-go` 一个服务。

## 相关链接

- 项目仓库：<https://github.com/qrst1ks/dnshe-go>
- Releases：<https://github.com/qrst1ks/dnshe-go/releases>
- Docker 镜像：`ghcr.io/qrst1ks/dnshe-go:latest`
- ddns-go：<https://github.com/jeessy2/ddns-go>
- 原 callback 项目：<https://github.com/qrst1ks/dnshe-ddns-go-callback>
- DNSHE：<https://www.dnshe.com>

## 项目功能

- 仅支持 IPv6 `AAAA` 记录。
- 支持通过 URL、有效物理网卡或命令获取公网 IPv6。
- 支持多个域名批量同步。
- DNSHE 更新流程为：查询子域名、查询记录、比对当前值、使用 `record_id` 更新。
- 提供轻量 Web UI。
- 支持手动同步和自动同步。
- 提供独立日志页和最近同步结果。
- 支持配置持久化。
- 支持 `DNSHE_API_KEY`、`DNSHE_API_SECRET`、`DNSHE_API_BASE_URL` 环境变量覆盖配置文件。
- 默认同步间隔为 300 秒，默认 TTL 为 600。

Web UI 不包含登录系统。打开页面后填写 DNSHE API、域名和 IP 获取方式即可。

## 部署方法

这个项目只保留两种部署方式：

1. 下载 Release 二进制文件启动。
2. 使用 Docker 镜像部署。

### 二进制启动

从 GitHub Releases 下载对应系统的压缩包：

<https://github.com/qrst1ks/dnshe-go/releases>

解压后直接双击 `dnshe-go` 启动服务，并自动打开 Web UI。

双击启动时，配置文件会保存在可执行文件同目录下：

```text
data/config.json
```

也可以从命令行启动同一个二进制文件：

```bash
./dnshe-go -l 127.0.0.1:9876 -c data/config.json
```

打开 Web UI：

```text
http://127.0.0.1:9876
```

### Docker 部署

直接拉取镜像启动：

```bash
docker run -d \
  --name dnshe-go \
  --restart unless-stopped \
  -p 9876:9876 \
  -v "$(pwd)/data:/data" \
  ghcr.io/qrst1ks/dnshe-go:latest
```

配置文件保存在宿主机：

```text
./data/config.json
```

容器内配置路径为：

```text
/data/config.json
```

也可以使用 Docker Compose：

```yaml
services:
  dnshe-go:
    image: ghcr.io/qrst1ks/dnshe-go:latest
    container_name: dnshe-go
    restart: unless-stopped
    ports:
      - "9876:9876"
    volumes:
      - ./data:/data
```

启动：

```bash
docker compose up -d
```

镜像支持：

- `linux/amd64`
- `linux/arm64`

如果无法拉取 GHCR 镜像，请确认 GitHub Packages 中的 package visibility 已设置为 Public。

## 发布

源码推送到 GitHub 仓库，二进制文件不提交到 Git。

创建版本 tag 后会自动发布：

```bash
git tag -a v1.0.0 -m "v1.0.0"
git push origin v1.0.0
```

Release 会包含：

- `linux/amd64`
- `linux/arm64`
- `darwin/amd64`
- `darwin/arm64`
- `windows/amd64`
- `checksums.txt`

## 目录结构

```text
.
├── .github/workflows/    # Docker 镜像和 Release 自动发布
├── Dockerfile            # Docker 镜像构建
├── docker-compose.yml    # Docker Compose 部署
├── Makefile              # 构建二进制
├── README.md
├── go.mod
├── main.go
└── internal/             # 服务源码
```

## 许可证

本项目使用 MIT License。详见 [LICENSE](LICENSE)。
