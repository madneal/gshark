<p align="center">
  <img alt="GShark logo" src="https://s1.ax1x.com/2018/10/17/idhZvj.png" />
  <h3 align="center">GShark</h3>
  <p align="center">轻松有效地扫描公开代码平台中的敏感信息。</p>
</p>

<div align="center">
  <strong>🇨🇳 中文版</strong> | <a href="README.md">🇺🇸 English</a>
</div>

# GShark [![Go Report Card](https://goreportcard.com/badge/github.com/madneal/gshark)](https://goreportcard.com/report/github.com/madneal/gshark) [![Release](https://github.com/madneal/gshark/actions/workflows/release.yml/badge.svg)](https://github.com/madneal/gshark/actions/workflows/release.yml)

GShark 是一个基于 Go/Gin 后端和 Vue 3 管理界面的敏感信息监测平台，支持 GitHub、GitLab、Sourcegraph、Postman 和 Gobuster 子域名发现。

GShark 通过公开 API 和公开索引进行搜索。私有仓库是否可见取决于平台 API 和 Token 权限。完整介绍可参考[文章](https://mp.weixin.qq.com/mp/appmsgalbum?__biz=MzI3MjA3MTY3Mw==&action=getalbum&album_id=2376148333116850178#wechat_redirect)、[视频](https://mp.weixin.qq.com/mp/appmsgalbum?__biz=MzI3MjA3MTY3Mw==&action=getalbum&album_id=1834365721464651778#wechat_redirect)和 [Wiki](https://github.com/madneal/gshark/wiki)。

## 主要特性

- GitHub 代码、Issue/PR 和公开 Gist 搜索
- GitLab 全局搜索，不支持时自动降级为项目爬取
- Sourcegraph 全局代码搜索（已有的 `searchcode` 规则也会继续支持）
- Postman 集合搜索和子域名发现
- 自定义规则、关键词/扩展名过滤、结果审核和导出
- 可选的入库前 AI 过滤，支持多个兼容 OpenAI 的 Provider
- Docker 部署和 scanner 资源限制

## 快速开始：Docker

环境要求：Docker Engine 和 Docker Compose。

```bash
git clone https://github.com/madneal/gshark.git
cd gshark

# 启动 MySQL、后端和前端；默认不会启动 scanner
./scripts/quick-docker.sh \
  --admin-user myadmin \
  --admin-password 'change-this-password'
```

打开 <http://localhost:8080>。登录后添加 Token 和规则，再启动扫描：

```bash
docker compose up -d scan
```

未自定义账号时，默认账号密码为 `gshark / gshark`。非本地部署请立即修改。quick 脚本会先初始化数据库，再启动 scanner；数据库初始化完成前不要手动启动 `scan`。

常用 Docker 命令：

```bash
docker compose ps
docker compose logs -f server
docker compose logs -f scan
docker compose restart server
docker compose stop scan       # 暂停扫描，不影响前端和后端
docker compose down            # 保留 ./mysql 数据目录
```

scanner 默认限制为 512 MB 内存、1 个 CPU，并设置 `GOMEMLIMIT=384MiB`。只有在查看 `docker stats gshark-scanner` 后，才建议调整 [docker-compose.yaml](docker-compose.yaml) 中的 `scan` 限制。除非确认要删除数据库，否则不要使用 `docker compose down -v`。

> **安全提示：** Docker 示例文件包含本地开发用凭据。生产环境使用前，请同时修改 `docker-compose.yaml` 和 `server/config.docker.yaml` 中的 MySQL 和 JWT 密钥。

## 第一次扫描

1. 打开 `http://localhost:8080` 并登录后台。
2. 在对应平台添加 Token。Token 只应保存在 GShark 中，不要提交到代码仓库或粘贴到 Issue。
3. 一条规则写一个搜索表达式，例如：

   ```text
   password in:file
   access_token org:example
   secret repo:owner/repository
   api_key extension:yaml
   ```

   GitHub 规则类型包括 `github`（代码）、`github_issue`（Issue/PR）和 `gist`（公开 Gist）。一条规则可以同时选择多个类型，例如 `github,github_issue,gist`。
4. 如果结果噪声较多，添加关键词或扩展名过滤器。目前过滤器用于 GitHub 代码和 Gist。
5. 如需 AI 入库前过滤，先配置 Provider，在系统页面点击“测试 AI 配置”，确认成功后再开启。该功能默认关闭。
6. 使用 `docker compose up -d scan` 启动 scanner；手动部署使用 `./gshark scan`。只要 Token 和规则有效，scanner 就会周期性运行。
7. 在结果页面确认真实结果，忽略示例值和误报，并按需导出。

## 平台配置

- **GitHub：** 在 [GitHub Token 设置](https://github.com/settings/tokens) 创建 Token。Token 必须具备规则所需的搜索权限。GitHub 搜索仍有 rate limit，GShark 会重试限流页面，但无法消除上游限制。
- **GitLab：** 配置 Token；自建 GitLab 还需要配置 GitLab Base URL。全局代码搜索不可用时，GShark 会搜索近期活跃的公开项目。
- **Sourcegraph：** 配置 `search.sourcegraph-url`（默认 `https://sourcegraph.com`）以及 `search.sourcegraph-token` 或 `SOURCEGRAPH_TOKEN`。结果受实例索引、超时和上游返回限制影响。
- **Postman/子域名发现：** 配置对应规则和字典后再启动 scanner。

规则支持通过 CSV 模板批量导入。数据库中已有的 `searchcode` 类型规则会由 Sourcegraph provider 继续处理，不需要迁移规则数据。

## 可选：入库前 AI 过滤

AI 过滤发生在结果写入 `search_result` 之前。只有模型严格返回 `real: true` 的结果才会入库；格式错误、超时、HTTP 错误和网络错误都会安全拒绝入库。内置限制为请求超时 30 秒、最多发送 6000 个字符。

该功能默认关闭。可在 `server/config.yaml`（Docker 使用 `server/config.docker.yaml`）中配置：

```yaml
system:
  ai_analysis_enabled: true
  ai_providers:
    - name: primary
      server: https://api.openai.com/v1/chat/completions
      token: your-token
      model: gpt-4o-mini
    - name: backup
      server: https://example.com/v1/chat/completions
      token: backup-token
      model: backup-model
```

当网络、鉴权或 HTTP 请求失败时按顺序切换 Provider。旧的 `ai_server`、`ai_token` 和 `model` 字段仍兼容为单 Provider 配置。系统页面的测试功能只发送合成占位证据，不会写入扫描结果。

## Release 包部署

Release 脚本支持 macOS 和 Linux，需要 `curl`、`jq`、`unzip`、`nginx` 和 `sudo`：

```bash
./scripts/quick-release.sh \
  --admin-user myadmin \
  --admin-password 'change-this-password'
```

脚本会下载匹配系统的 Release zip，将前端部署到检测到的 Nginx 根目录，把 `/api/` 代理到后端 `8888` 端口，并通过 `8080` 提供网页。使用本地包时增加 `--file ./gshark_linux_amd64.zip`；数据库已初始化时可增加 `--skip-init`。

## 手动部署

环境要求：MySQL 8.0+、Go 1.25+、Node.js 20+、npm 和 Nginx。

```bash
cd server
cp config-temp.yaml config.yaml
go build -o gshark
./gshark serve       # 后端，默认 8888 端口
./gshark scan        # scanner，数据库初始化后再运行
```

在 `web/` 中执行 `npm install && npm run build`，将 `web/dist` 交给 Nginx，并把 `/api/` 反向代理到 `127.0.0.1:8888`。浏览器访问 Nginx 端口（示例为 8080），不要直接访问后端端口。

已有部署升级时，请先备份 MySQL，再执行 [sql.md](sql.md) 中对应版本的 SQL，重启后端，确认迁移完成后再启动 scanner。不能假设只替换二进制就会完成所有历史迁移。

## 配置说明

主配置文件是 `server/config.yaml`。请从 [server/config-temp.yaml](server/config-temp.yaml) 开始；Docker 使用 [server/config.docker.yaml](server/config.docker.yaml)。可以通过 `-c` 或 `GSHARK_CONFIG` 指定配置文件：

```bash
./gshark -c /path/to/config.yaml serve
GSHARK_CONFIG=/path/to/config.yaml ./gshark scan
```

## 常见问题

- **GitHub 返回 401：** 用 GitHub API 单独验证 Token，检查过期时间、权限以及是否绑定到正确的平台规则。
- **没有扫描结果：** 确认 `scan` 正在运行、规则语法正确、Token 有效且平台可访问，然后查看 `docker compose logs server scan`。
- **scanner 没启动：** 等待 MySQL health check 和数据库初始化完成，再执行 `docker compose up -d scan`。
- **页面慢或内存高：** 先查看 `docker stats`、scanner 日志和 MySQL 状态，再调整资源限制。
- **GitHub 限流：** 减少规则噪声和搜索范围，接受扫描延迟；scanner 会重试限流页面。
- **升级异常：** 确认版本、备份数据库、执行对应的 [sql.md](sql.md) 迁移，再按顺序重启后端和 scanner。

## 开发

后端：

```bash
cd server
go mod tidy
cp config-temp.yaml config.yaml
go run main.go serve
```

前端：

```bash
cd web
npm install
npm run serve
```

可先运行 AI 相关的聚焦测试：`go test ./service -run 'Test(SearchResultContent|ParseSearchResultAnalysis|AnalyzeSearchResult|TestAIConfig)'`。完整测试请只在隔离的测试数据库中运行。macOS ARM 如果需要在服务器信息页面显示 CPU 百分比，请设置 `CGO_ENABLED=1`。

## 资源

### 文章

* [多平台的敏感信息监测工具-GShark](https://mp.weixin.qq.com/s?__biz=MzI3MjA3MTY3Mw==&mid=2247484283&idx=1&sn=3232df7d321c0f62ce61b7e6368204ad&chksm=eb396deddc4ee4fb0c825a378c085223b87fc45f05648d46e7bdc24a03fb83ad6c7ade414df7#rd)
* [GShark-监测你的 Github 敏感信息泄露](https://mp.weixin.qq.com/s?__biz=MzI3MjA3MTY3Mw==&mid=2247483770&idx=1&sn=9f02c2803e1c946e8c23b16ff3eba757&chksm=eb396fecdc4ee6fa2f378e846f354f45acf6e6f540cfd54190e9353df47c7707e3a2aadf714f&token=1578822041&lang=zh_CN#rd)

### 视频

* [GShark v1.5.0 版本及 Docker 使用指南](https://www.bilibili.com/video/BV1oUe3eBEMz/)
* [GShark v1.3.0 版本支持 Docker](https://www.bilibili.com/video/BV1BH4y1C7Ga/)
* [GShark 支持多种规则类型以及规则配置建议](https://www.bilibili.com/video/BV1uY4y177SX)
* [批量导入规则](https://mp.weixin.qq.com/s?__biz=MzI3MjA3MTY3Mw==&mid=2247484546&idx=1&sn=818915279c5199457340ade89d6cbd54&chksm=eb396a14dc4ee302039bcb1474380a6049dba84370345b7813049aa8feb49a98f89d47ec5d5b#rd)
* [GShark部署](https://mp.weixin.qq.com/s?__biz=MzI3MjA3MTY3Mw==&mid=2247484487&idx=1&sn=78f942ccf6861f433fc7f4a60564441c&chksm=eb396ad1dc4ee3c7505362da243433e54a2b558c96fbbb50f8b6cea87d1f9bc920b249b72705#rd)
* [windows 部署](https://mp.weixin.qq.com/s?__biz=MzI3MjA3MTY3Mw==&mid=2247484289&idx=1&sn=2b0f1c38b88c924ad514fb64b559b784&chksm=eb396d17dc4ee4018573dde6c3bfce83903c86034403539eaf1b87b89c4a4dd44f957a308818#rd)
* [GShark v1.0.2 版本发布](https://www.bilibili.com/video/BV1Zx4y1G7FX/)
* [GShark v1.1.0 更新内容介绍](https://www.bilibili.com/video/BV1aG4y1c72N/)

## 许可证

[Apache License 2.0](https://github.com/madneal/gshark/blob/master/LICENSE)

## 404StarLink 2.0 - Galaxy

![](https://github.com/knownsec/404StarLink-Project/raw/master/logo.png)

GShark 是 404Team [星链计划2.0](https://github.com/knownsec/404StarLink2.0-Galaxy)中的一环，如果对 GShark 有任何疑问又或是想要找小伙伴交流，可以参考星链计划的加群方式。

- [https://github.com/knownsec/404StarLink2.0-Galaxy#community](https://github.com/knownsec/404StarLink2.0-Galaxy#community)
