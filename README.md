# dnshe-go

`dnshe-go` 是一个轻量的 IPv6 DDNS 服务，支持 DNSHE 和 Cloudflare，用于自动更新 IPv6 `AAAA` 记录。

它直接获取本机公网 IPv6，并调用对应服务商的 API 更新域名记录。部署后只需要运行一个服务。

当前版本：`v2.0.0`，更新说明见 [GitHub Releases](https://github.com/qrst1ks/dnshe-go/releases)。

## 适用场景

- 你使用 DNSHE 或 Cloudflare 管理域名解析。
- 你只需要更新 IPv6 `AAAA` 记录。
- 你的公网 IPv6 会变化，需要自动同步到域名服务商。
- 你希望用一个轻量 Web UI 完成配置和查看日志。

## 与 ddns-go 的关系

`dnshe-go` 基于 `ddns-go` 的 DDNS 工作方式实现，面向 IPv6 `AAAA` 记录更新场景做了简化和整合。

它保留了 DDNS 服务中常用的 IP 检测、周期同步、手动同步和 Web UI 配置思路，同时直接内置 DNSHE 和 Cloudflare 的更新流程。

因此，如果你之前熟悉 `ddns-go`，可以把 `dnshe-go` 理解为一个面向 DNSHE / Cloudflare 和 IPv6 `AAAA` 记录的轻量实现。

## 功能

- 支持 DNSHE 和 Cloudflare，两个服务商可独立启用、同时同步。
- 仅支持 IPv6 `AAAA` 记录。
- 支持通过 URL、有效物理网卡或命令获取 IPv6。
- 每个服务商支持多个域名批量同步。
- Cloudflare 记录不存在时自动创建 `AAAA` 记录（`proxied=false`）。
- 支持自动同步和手动同步。
- 提供轻量 Web UI。
- 提供日志页和最近同步结果。
- 配置自动保存。
- 默认同步间隔为 300 秒。
- 默认 TTL 为 600。

> **安全提示**：Web UI 没有登录鉴权，并且可以查看完整 API 密钥。仅推荐在本机使用，不建议开放到公网或不可信局域网。如需远程访问，请使用 SSH 隧道，或带 TLS 和身份认证的反向代理；因开放网络访问产生的风险由使用者自行承担。

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

容器默认只监听 `127.0.0.1:9999`，仅允许本机访问。如确需从局域网访问，请自行覆盖容器启动参数中的监听地址，并自行承担相应风险。

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

两个服务商默认均未启用。首次打开 Web UI 后，按需启用服务商并填写：

- DNSHE：API Key、API Secret、需要同步的域名。
- Cloudflare：API Token、需要同步的域名。
- IPv6 获取方式。
- 同步间隔和 TTL。

两个服务商互不影响，可以只启用一个，也可以同时启用。各自的域名列表独立维护。

已保存的密钥默认只显示首尾掩码，点击输入框右侧的眼睛按钮可从本地服务读取并显示完整内容，再次点击隐藏。

Cloudflare API Token 建议只授予以下权限：

- `Zone:Read`
- `DNS:Edit`

Token 可以在 Cloudflare 控制台的 My Profile → API Tokens 页面创建。Cloudflare 上不存在的 `AAAA` 记录会自动创建（`proxied=false`）；更新已有记录时会保留原有的 `proxied` 设置。

IPv6 获取方式支持：

- URL：从外部 IPv6 查询服务获取。
- 网卡：从本机有效物理网卡读取。
- 命令：执行自定义命令，命令输出中需要包含 IPv6。

保存配置后，可以点击“立即同步”进行手动测试。

## 环境变量

除配置文件外，也可以通过环境变量提供凭据：

| 变量 | 说明 |
| --- | --- |
| `DNSHE_API_KEY` | DNSHE API Key |
| `DNSHE_API_SECRET` | DNSHE API Secret |
| `DNSHE_API_BASE_URL` | DNSHE API 地址，一般无需修改 |
| `CLOUDFLARE_API_TOKEN` | Cloudflare API Token |
| `CLOUDFLARE_API_BASE_URL` | Cloudflare API 地址，一般无需修改 |

环境变量只覆盖凭据和 API 地址，是否启用服务商仍以 Web UI 或配置文件中的 `enable` 为准。

## 常见问题

### Docker 启动后访问不到 Web UI

先确认容器是否正在运行：

```bash
docker ps -a --filter name=dnshe-go
```

默认的 Docker 配置使用 host 网络模式，`docker ps` 不会显示 `9999:9999` 端口映射，这是正常现象。直接访问：

```text
http://127.0.0.1:9999
```

如果你更新过镜像或修改过 compose，建议重新创建容器：

```bash
docker pull ghcr.io/qrst1ks/dnshe-go:latest
docker compose up -d --force-recreate dnshe-go
```

### 日志显示 DNSHE API key/secret not configured

这是启用了 DNSHE 但尚未填写凭据的正常状态。打开 Web UI 填写 DNSHE API Key 和 Secret 后保存即可，或者关闭 DNSHE 开关。

### 日志显示 Cloudflare API token not configured

这是启用了 Cloudflare 但尚未填写 Token 的正常状态。打开 Web UI 填写 Cloudflare API Token 后保存即可。

### Cloudflare 提示 zone not found

确认域名已经添加到当前 API Token 可访问的 Cloudflare 账号中，并且 Token 具备 `Zone:Read` 和 `DNS:Edit` 权限。

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
- Cloudflare API Tokens：<https://dash.cloudflare.com/profile/api-tokens>

## 许可证

本项目使用 MIT License。详见 [LICENSE](LICENSE)。
