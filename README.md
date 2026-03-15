# gproxy

HTTP 代理转发服务 - 通过指定的 HTTP 代理转发请求到目标 HTTPS 地址。

## 功能

- 支持通过 HTTP 代理转发请求到 HTTPS 目标地址
- 支持命令行参数配置
- 轻量级、跨平台

## 安装

### 从源码编译

```bash
git clone https://github.com/graydovee/gproxy.git
cd gproxy
go build -o gproxy
```

### 使用 Docker

```bash
docker pull ghcr.io/graydovee/gproxy:latest
```

## 使用方法

### 命令行参数

| 参数 | 短参数 | 说明 | 默认值 | 示例 |
|------|--------|------|--------|------|
| `--listen` | `-l` | 监听地址和端口 | `:8080` | `:8080` |
| `--proxy` | `-p` | HTTP 代理地址 | (必填) | `http://127.0.0.1:7890` |
| `--target` | `-t` | 目标 HTTPS 地址 | (必填) | `https://api.openai.com` |

### 运行

```bash
# 使用短参数
./gproxy -l :8080 -p http://127.0.0.1:7890 -t https://api.openai.com

# 使用长参数
./gproxy --listen :8080 --proxy http://127.0.0.1:7890 --target https://api.openai.com

# 查看帮助
./gproxy --help
```

### Docker 运行

```bash
docker run -d \
  -p 8080:8080 \
  ghcr.io/graydovee/gproxy:latest \
  -l :8080 \
  -p http://host.docker.internal:7890 \
  -t https://api.openai.com
```

### 测试

```bash
curl http://localhost:8080/v1/models
```

## 使用场景

- 通过 HTTP 代理访问外部 HTTPS API
- 网络受限环境下的 API 代理
- 开发调试时的请求转发

## Docker 镜像

镜像支持以下平台：
- `linux/amd64`
- `linux/arm64`

```bash
# 拉取最新版本
docker pull ghcr.io/graydovee/gproxy:latest

# 拉取指定版本
docker pull ghcr.io/graydovee/gproxy:v1.0.0
```

## License

MIT
