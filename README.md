# dnshe-go

`dnshe-go` 是一个面向 DNSHE 的轻量 IPv6 DDNS 服务。它直接获取本机公网 IPv6，并调用 DNSHE API 更新 `AAAA` 记录

## 当前能力

- 支持 IPv6 `AAAA` 记录。
- 支持通过 URL、有效物理网卡或命令获取公网 IPv6。
- 支持多个域名批量同步。
- DNSHE 更新流程为：查询子域名、查询记录、比对当前值、使用 `record_id` 更新。
- 提供轻量 Web UI、手动同步、独立日志页和最近同步结果。
- 支持 `DNSHE_API_KEY`、`DNSHE_API_SECRET`、`DNSHE_API_BASE_URL` 环境变量覆盖配置文件。
- 默认同步间隔为 300 秒，默认 TTL 为 600。

## 本地运行

```bash
go run . -l 127.0.0.1:9876 -c data/config.json
```

打开 `http://127.0.0.1:9876` 填写 DNSHE API、域名和 IP 获取方式。

Web UI 不包含登录系统。`API Secret` 使用前端遮罩显示，不触发浏览器登录表单。

## 构建二进制

```bash
make build
```

生成的启动文件位于项目根目录的 `dnshe-go`。直接双击这个文件会启动服务，并自动打开 Web UI。

双击启动时，配置文件会保存在可执行文件同目录下的 `data/config.json`。

命令行启动也可以显式指定监听地址和配置文件：

```bash
./dnshe-go -l 127.0.0.1:9876 -c data/config.json
```

## 单次同步

```bash
DNSHE_API_KEY=xxx DNSHE_API_SECRET=yyy go run . -once -c data/config.json
```

## Docker

直接拉取镜像启动：

```bash
docker run -d \
  --name dnshe-go \
  --restart unless-stopped \
  -p 9876:9876 \
  -v "$(pwd)/data:/app/data" \
  ghcr.io/qrst1ks/dnshe-go:latest
```

也可以使用 compose：

```bash
docker compose up -d
```

配置文件保存在 `./data/config.json`，文件权限为 `0600`。

镜像发布到 `ghcr.io/qrst1ks/dnshe-go`，支持 `linux/amd64` 和 `linux/arm64`。推送 `main` 分支后会自动发布 `latest`，推送 `v*` tag 后会发布对应版本标签。

这个项目把这段流程内置在同步循环里：IP 变化后直接查询 DNSHE 并更新记录，所以部署时只需要运行一个服务。
