# dnshe-go

`dnshe-go` 是一个面向 DNSHE 的轻量 IPv6 DDNS 服务，用于自动更新 DNSHE 的 IPv6 `AAAA` 记录。

它直接获取本机公网 IPv6，并调用 DNSHE API 更新域名记录。部署后只需要运行一个服务。

## 适用场景

- 你使用 DNSHE 管理域名解析。
- 你只需要更新 IPv6 `AAAA` 记录。
- 你的公网 IPv6 会变化，需要自动同步到 DNSHE。
- 你希望用一个轻量 Web UI 完成配置和查看日志。

## 与 ddns-go 的关系

`dnshe-go` 基于 `ddns-go` 的 DDNS 工作方式实现，面向 DNSHE 的 IPv6 `AAAA` 记录更新场景做了简化和整合。

它保留了 DDNS 服务中常用的 IP 检测、周期同步、手动同步和 Web UI 配置思路，同时直接内置 DNSHE 更新流程。

因此，如果你之前熟悉 `ddns-go`，可以把 `dnshe-go` 理解为一个面向 DNSHE 和 IPv6 `AAAA` 记录的轻量实现。

## 功能

- 仅支持 IPv6 `AAAA` 记录。
- 支持通过 URL、有效物理网卡或命令获取 IPv6。
- 支持多个域名批量同步。
- 支持自动同步和手动同步。
- 提供轻量 Web UI。
- 提供日志页和最近同步结果。
- 配置自动保存。
- 默认同步间隔为 300 秒。
- 默认 TTL 为 600。

Web UI 不包含登录系统。建议只在可信内网环境访问，或自行通过反向代理增加访问控制。

## Docker 部署

推荐使用 Docker 部署。

### docker run

```bash
docker run -d \
  --name dnshe-go \
  --restart unless-stopped \
  --network host \
  -v "$(pwd)/data:/data" \
  ghcr.io/qrst1ks/dnshe-go:latest
```

启动后打开：

```text
http://127.0.0.1:9999
```

配置文件保存在宿主机：

```text
./data/config.json
```

容器内配置路径：

```text
/data/config.json
```

### Docker Compose

```yaml
services:
  dnshe-go:
    image: ghcr.io/qrst1ks/dnshe-go:latest
    container_name: dnshe-go
    restart: unless-stopped
    network_mode: host
    volumes:
      - ./data:/data
```

启动：

```bash
docker compose up -d
```

更新镜像：

```bash
docker pull ghcr.io/qrst1ks/dnshe-go:latest
docker compose up -d --force-recreate dnshe-go
```

镜像支持：

- `linux/amd64`
- `linux/arm64`

如果无法拉取 GHCR 镜像，请确认 GitHub Packages 中的 package visibility 已设置为 Public。

这里默认使用 host 网络模式，是为了让容器能够读取宿主机网卡信息。host 模式下不需要配置 `ports`，服务会直接监听宿主机的 `9999` 端口。

## 本机部署

不使用 Docker 时，可以从发布页下载对应系统的压缩包：

<https://github.com/qrst1ks/dnshe-go/releases>

解压后直接双击 `dnshe-go` 启动服务。程序会自动打开 Web UI。

双击启动时，配置文件保存在可执行文件同目录下：

```text
data/config.json
```

也可以从命令行启动：

```bash
./dnshe-go -l 127.0.0.1:9999 -c data/config.json
```

打开 Web UI：

```text
http://127.0.0.1:9999
```

Release 文件说明：

- `darwin_amd64`：Intel Mac
- `darwin_arm64`：Apple Silicon Mac
- `linux_amd64`：普通 x86_64 Linux
- `linux_arm64`：ARM64 Linux
- `windows_amd64`：64 位 Windows

## Web UI 配置

首次打开 Web UI 后填写：

- DNSHE API Key
- DNSHE API Secret
- 需要同步的域名
- IPv6 获取方式
- 同步间隔和 TTL

IPv6 获取方式支持：

- URL：从外部 IPv6 查询服务获取。
- 网卡：从本机有效物理网卡读取。
- 命令：执行自定义命令，命令输出中需要包含 IPv6。

保存配置后，可以点击“立即同步”进行手动测试。

## 常见问题

### Docker 启动后访问不到 Web UI

先确认容器是否正在运行：

```bash
docker ps -a --filter name=dnshe-go
```

README 默认使用 host 网络模式，`docker ps` 不会显示 `9999:9999` 端口映射，这是正常现象。直接访问：

```text
http://127.0.0.1:9999
```

如果你更新过镜像或修改过 compose，建议重新创建容器：

```bash
docker pull ghcr.io/qrst1ks/dnshe-go:latest
docker compose up -d --force-recreate dnshe-go
```

### 日志显示 DNSHE API key/secret not configured

这是首次启动的正常状态。打开 Web UI 填写 DNSHE API Key 和 Secret 后保存即可。

### Docker 配置文件在哪里

如果按 README 的 Docker 命令启动，配置文件在宿主机：

```text
./data/config.json
```

### 二进制配置文件在哪里

双击二进制启动时，配置文件在二进制同目录：

```text
data/config.json
```

## 相关链接

- 项目仓库：<https://github.com/qrst1ks/dnshe-go>
- Releases：<https://github.com/qrst1ks/dnshe-go/releases>
- Docker 镜像：`ghcr.io/qrst1ks/dnshe-go:latest`
- ddns-go：<https://github.com/jeessy2/ddns-go>
- DNSHE：<https://www.dnshe.com>

## 许可证

本项目使用 MIT License。详见 [LICENSE](LICENSE)。
