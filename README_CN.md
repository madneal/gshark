<p align="center">
   <img alt="GShark logo" src="https://s1.ax1x.com/2018/10/17/idhZvj.png" />
   <h3 align="center">GShark</h3>
   <p align="center">轻松有效地扫描敏感信息。</p>
</p>

<div align="center">
   <strong>🇨🇳 中文版</strong> | <a href="README.md">🇺🇸 English</a>
</div>

# GShark [![Release](https://github.com/madneal/gshark/actions/workflows/release.yml/badge.svg)](https://github.com/madneal/gshark/actions/workflows/release.yml)

GShark 是一个敏感信息检测和管理平台。后端基于 Go 和 Gin 构建，当前前端基于 Vue 3、Vite、Vue Router 4、Vuex 4 和 Element Plus 构建。完整介绍请参考[文章](https://mp.weixin.qq.com/mp/appmsgalbum?__biz=MzI3MjA3MTY3Mw==&action=getalbum&album_id=2376148333116850178#wechat_redirect)和[视频](https://mp.weixin.qq.com/mp/appmsgalbum?__biz=MzI3MjA3MTY3Mw==&action=getalbum&album_id=1834365721464651778#wechat_redirect)。GShark 扫描配置平台可访问的仓库，不扫描本地源码目录。

关于 GShark 的使用，请参考 [wiki](https://github.com/madneal/gshark/wiki)。

# 主要特性

* 🌐 多平台支持：GitHub、GitLab、Sourcegraph、Postman 等
* 🔍 灵活的规则管理：自定义扫描规则和过滤，支持白名单/黑名单
* 🔑 细粒度访问控制：可配置的菜单和 API 权限
* 🔄 子域名发现：集成 gobuster 进行子域名枚举
* 🚀 Docker 部署：容器化部署，易于设置
* 📊 Vue 3 管理界面：基于 Vite 的 Web 界面，用于任务和结果管理
* 🔁 更健壮的扫描：GitHub 触发限流后会自动重试，GitLab 全局搜索不可用时会自动降级为逐项目爬取

# 快速开始

初始化后的默认登录账号（未自定义时）：

```text
gshark / gshark
```

非本地部署完成后请立即修改默认密码。

可通过脚本参数自定义管理员账号（无需打开初始化页面）：

```bash
./scripts/quick-docker.sh --admin-user myadmin --admin-password 'S3cret!'
# 或
./gshark init --host 127.0.0.1 --user root --password madneal --db gshark \
  --admin-user myadmin --admin-password 'S3cret!'
```

## Docker 部署

```bash
# 克隆仓库
git clone https://github.com/madneal/gshark.git
cd gshark

# 构建镜像、初始化 MySQL，并启动 server/web
./scripts/quick-docker.sh

# 初始化时设置自定义管理员
./scripts/quick-docker.sh --admin-user myadmin --admin-password 'S3cret!'

# 初始化完成后在同一命令中启动扫描器
./scripts/quick-docker.sh --with-scan
```

> [!IMPORTANT]
> Docker quick 脚本会先启动 MySQL、初始化数据库，只有使用 `--with-scan` 时才会在初始化完成后启动扫描器。如果手动使用 Docker Compose 启动扫描器，请先等待数据库初始化完成。

> [!TIP]
> 不使用 `--with-scan` 时，请登录 `http://localhost:8080` 配置 Token 和规则，再执行 `docker compose up -d scan`；如果使用了 `--skip-init`，则需先在网页中完成数据库初始化。

扫描器在 `docker-compose.yaml` 中默认配置了资源保护：内存上限 512 MB、CPU 上限 1 核，并通过 `GOMEMLIMIT=384MiB` 让 Go 更积极地回收内存。搜索结果会按分页写入数据库，避免大规模搜索结果全部驻留内存。如果 scanner 因超过内存上限被 OOM kill，Compose 会自动重启；可以先使用 `docker stats gshark-scanner` 观察实际占用，再决定是否继续下调限制。

### Docker 运维

```bash
docker compose ps
docker compose up -d scan
docker compose logs -f server scan
docker compose restart scan
docker compose stop scan
```

## Release 包部署

此方式需要在 macOS 或 Linux 上准备 MySQL、Nginx、`curl`、`jq` 和 `unzip`。

```bash
git clone https://github.com/madneal/gshark.git
cd gshark

# 下载最新发布包、配置 Nginx、初始化数据库并启动后端
./scripts/quick-release.sh

# 也可以部署已经下载的发布包
./scripts/quick-release.sh --file ./gshark_linux_amd64.zip
```

## 手动部署

### 环境要求

* Nginx
* MySQL **8.0+**
* Go **1.25+**，用于构建后端
* Node.js **20+** 和 npm，用于构建前端

建议使用 Nginx 部署前端。构建 Vite 项目后，将生成的 `web/dist` 文件放置在 `/var/www/html` 中，并配置 Nginx 将 `/api/` 反向代理到后端服务。详细的部署教程可以观看 [bilibili](https://www.bilibili.com/video/BV1Py4y1s7ap/) 或 [youtube](https://youtu.be/bFrKm5t4M54) 上的视频。Windows 部署请参考[此链接](https://www.bilibili.com/video/BV1CA411L7ux/)。

### Nginx

使用 `nginx -t` 定位当前生效的 `nginx.conf`，然后添加下面的 server 配置。请根据实际安装位置调整前端目录。

```nginx
worker_processes  1;

events {
    worker_connections  1024;
}

http {
    include       mime.types;
    default_type  application/octet-stream;
    sendfile        on;
    keepalive_timeout  65;
    server {
        listen       8080;
        server_name  localhost;

        location / {
            root   /var/www/html;
            index  index.html index.htm;
            try_files $uri $uri/ /index.html;
        }
        location /api/ {
            proxy_set_header Host $http_host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            rewrite ^/api/(.*)$ /$1 break;
            proxy_pass http://127.0.0.1:8888;
        }
        error_page   500 502 503 504  /50x.html;
        location = /50x.html {
            root   /var/www/html;
        }
    }
}
```

从 [releases](https://github.com/madneal/gshark/releases) 下载对应平台的发布包，然后将 `dist` 中的全部内容复制到 Nginx 前端目录：

```bash
unzip gshark*.zip
cd gshark*
sudo mkdir -p /var/www/html
sudo cp -R dist/. /var/www/html/
```

Homebrew macOS 常用的前端目录是 `$(brew --prefix)/var/www`；Nginx 的 `root` 和文件复制目标必须保持一致。

修改配置后先校验，再重启 Nginx：

```bash
sudo nginx -t
# Homebrew macOS
brew services restart nginx
# 使用 systemd 的 Linux
sudo systemctl restart nginx
```

### 服务器服务

```shell
cp config-temp.yaml config.yaml
# 启动服务前编辑 config.yaml，配置 MySQL 连接。
./gshark serve
```

后端默认监听 `8888`；使用 Nginx 时，应通过前端端口（例如 `8080`）访问网页。

如果您之前没有初始化数据库，您将首先被重定向到数据库初始化页面。

<img width="936" alt="image" src="https://github.com/user-attachments/assets/dfa7e53e-dc4a-4697-831f-a4f4f3810c3c">

### 扫描服务

```shell
./gshark scan
```

启动扫描服务前，请先配置所需平台的 Token 和规则。

### 增量部署

升级前应备份数据库并阅读对应版本的 Release Notes。模型字段新增会由 GORM 在启动时自动迁移；仅当发布说明明确要求时，才执行 [sql.md](https://github.com/madneal/gshark/blob/master/sql.md) 中对应版本的语句，不要在每次升级时执行整个文件。

## 开发

### 服务器

```shell
git clone https://github.com/madneal/gshark.git
cd gshark/server
go mod download
cp config-temp.yaml config.yaml
# 编辑 config.yaml，配置 MySQL 连接。
go build -o gshark .
```

运行 Web 服务器：

```shell
./gshark serve
```

配置 Token 和规则后，在另一个终端运行扫描任务：

```shell
./gshark scan
```

开发时如不需要生成二进制文件，也可以使用 `go run`：

```shell
go run main.go serve
go run main.go scan
```

> [!NOTE]
> 在 macOS ARM 上，服务器状态页面的 CPU 百分比采集依赖 cgo。如果需要显示 CPU 使用率，请在运行或构建后端时启用 `CGO_ENABLED=1`：
>
> ```shell
> CGO_ENABLED=1 go run main.go serve
> ```

### Web 前端

```bash
cd ../web
npm install
npm run serve
```

## 使用方法

### 添加 Token

#### GitHub

创建一个[细粒度个人访问令牌](https://github.com/settings/personal-access-tokens/new)，仅授予规则所需的仓库和权限。建议设置较短的有效期，不要复用管理员令牌。具体请参考 GitHub 的[个人访问令牌指南](https://docs.github.com/zh/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens)。GitLab 搜索需要单独配置 GitLab Token。

[![iR2TMt.md.png](https://s1.ax1x.com/2018/10/31/iR2TMt.md.png)](https://imgchr.com/i/iR2TMt)

### 规则配置

GitHub 和 GitLab 规则使用对应平台的搜索语法。规则内容应填写在该平台上实际使用的搜索表达式。同一套 GitHub Token 会覆盖三个搜索面：

* `github`：仓库代码搜索（`in:file`）
* `github_issue`：Issue 和 PR（默认 `in:title,body,comments`，也可自行加 `is:issue` / `is:pr`）
* `gist`：公开 Gist 搜索，命中后再通过官方 Gist API 拉取文件内容

同一条规则可以勾选多种类型，例如 `github,github_issue,gist`。

您可以下载规则导入模板 CSV 文件，然后批量导入规则。

<img width="572" alt="image" src="https://user-images.githubusercontent.com/12164075/212504597-3e1ad5bd-bacf-433e-83e8-08de7eee6509.png">

### 过滤器配置

过滤器目前针对 GitHub 相关扫描。`keyword` 适用于 `github`、`github_issue` 和 `gist`；`extension` 适用于 `github` 代码搜索和 `gist`。两类都可以配置为黑名单或白名单。

更多信息，您可以参考这个[视频](https://www.bilibili.com/video/BV1aG4y1c72N/?vd_source=ef4657ebf0549af8755f75118b6e81bb)。

## 扫描运营

1. 初始化数据库并登录管理页面。
2. 配置有效的平台 Token，并启用需要的规则；当平台侧搜索较宽泛时，可以增加本地匹配正则，对候选证据做更严格的校验。
3. Docker 部署执行 `docker compose up -d scan`，手动部署执行 `./gshark scan`，启动扫描器。
4. 查看扫描日志页面，或执行 `docker compose logs -f scan`，确认各平台扫描正常完成。
5. 在结果页面确认真实密钥，忽略示例值和误报，并按需导出结果。

## 配置

手动部署时，将 `config-temp.yaml` 复制为 `config.yaml`，再根据实际环境配置数据库和其他设置。

### GitLab 基础 URL

<img width="363" alt="image" src="https://user-images.githubusercontent.com/12164075/203898719-1ce66395-083d-4226-937f-b6eed859addc.png">

当 GitLab 全局搜索不可用时，GShark 会降级为爬取近期活跃的公开项目。`config.yaml` 中的 `search.gitlab-discover-pages` 和 `search.gitlab-batch-size` 用于控制每个扫描周期发现的项目分页数和搜索的项目数量（默认分别为 5 和 50）。

### Sourcegraph 全局代码搜索

Sourcegraph provider 使用 Stream API 对 Sourcegraph 索引中的全部仓库执行规则搜索，不需要逐个配置仓库。默认使用 `https://sourcegraph.com` 公共实例；如果使用自建实例，可在 `config.yaml` 中设置 `search.sourcegraph-url`，并通过 `search.sourcegraph-token` 或 `SOURCEGRAPH_TOKEN` 提供访问令牌。公共实例默认排除的 fork 和 archived 仓库会被 GShark 显式包含，以尽量扩大覆盖范围。Sourcegraph 的结果仍受其索引范围、搜索超时和结果限制影响，扫描日志会记录上游返回的限制信息。

数据库中已有的 `searchcode` 类型规则也会由该 provider 继续执行，不需要迁移规则数据。

### 入库前 AI 过滤

GShark 支持在搜索结果写入 `search_result` 之前，将每条新结果发送到一个或多个兼容 OpenAI Chat Completions 的接口进行判断。请按优先级配置 `system.ai_providers`；旧的 `ai_server`、`ai_token`、`model` 字段仍兼容为单个 Provider。网络、HTTP 或鉴权失败时切换到下一个 Provider；模型明确返回 `real: false` 时不再切换。模型必须返回 `{"real":true|false,"confidence":0.0,"reason":"..."}` 格式的 JSON；只有 `real: true` 的结果才会入库，响应格式错误、超时或接口异常都会安全拒绝入库。每次请求使用内置限制：超时 30 秒、最多发送 6000 个字符。该功能默认关闭。

示例：

```yaml
system:
  ai_analysis_enabled: true
  ai_providers:
    - name: primary
      server: https://api.openai.com/v1/chat/completions
      token: your-openai-token
      model: gpt-4o-mini
    - name: backup
      server: https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions
      token: your-dashscope-token
      model: qwen-plus
```

系统配置页提供“测试 AI 配置”按钮，会使用合成的占位符内容验证接口连通性、鉴权、模型可用性和响应格式，不会写入搜索结果。每次请求使用内置限制：超时 30 秒、最多发送 6000 个字符。

### 本地上下文匹配

规则可以在规则页面配置可选的 `matchPattern`。`content` 仍作为平台侧的搜索表达式，用于发现候选结果；“本地匹配正则”使用兼容 Go/RE2 的正则表达式，在结果入库前校验返回的代码片段。`matchPattern` 为空时保持原有行为。

在规则页面新增或编辑规则时，可以按以下方式填写：

| 字段 | 示例 | 用途 |
| --- | --- | --- |
| 规则内容 | `ghp_` | 平台侧的宽泛搜索表达式 |
| 本地匹配正则 | `ghp_(?P<key>[A-Za-z0-9_]{20,})` | 对候选结果进行最终的本地校验，并提取 `key` 分组 |
| Key 验证 API | `github` | 入库前调用对应平台 API 验证提取出的 key |

CSV 导入模板也支持这些字段，“本地匹配正则”和 “Key 验证 API”分别位于最后两列。

例如，GitHub token 规则可以配置为：

```text
content: ghp_
matchPattern: ghp_(?P<key>[A-Za-z0-9_]{20,})
validationType: github
```

支持的 Key 验证 API 为 `github`、`gitlab`、`sourcegraph` 和 `postman`，分别调用各平台的当前用户接口。验证返回 2xx 才会入库；401 会作为无效 key 忽略，网络错误、限流和其他非确定性错误会暂不入库并在下一轮重试。未配置验证 API 时保持原有的仅搜索和本地匹配行为。

## 常见问题

1. GShark 扫描的是本地代码还是公开平台代码？

GShark 扫描配置平台可访问的仓库，不扫描本地源码目录。GitHub 扫描基于 GitHub Search API；GitLab 扫描依赖 GitLab 搜索能力；Sourcegraph 扫描其索引中的仓库。私有仓库的覆盖范围取决于平台 API、Sourcegraph 实例和 Token 权限。

2. 推荐怎么部署？

新用户优先使用 quick 脚本：

```shell
./scripts/quick-docker.sh
./scripts/quick-docker.sh --with-scan
./scripts/quick-release.sh
```

手动部署适合需要自定义 Nginx、MySQL、后端配置的场景。

3. 部署环境有什么要求？

MySQL 需要 8.0+。手动构建时需要 Go 1.25+、Node.js 20+、npm 和 Nginx。Docker 部署建议直接使用项目提供的 compose 和 quick 脚本，避免老教程里的配置差异。

4. 初始化后默认账号是什么？

默认账号密码是 `gshark / gshark`。生产环境部署完成后应立即修改密码。

5. Docker 部署后 scanner 为什么没启动或没结果？

scanner 依赖数据库初始化。`./scripts/quick-docker.sh --with-scan` 会等待初始化完成后再启动 scanner；手动使用 Compose 时，如果 scanner 启动过早可能会退出，初始化完成后执行 `docker compose up -d scan` 即可。排查时优先查看 scanner 和 server 日志。

6. GShark 的核心运行链路是什么？

基本链路是：配置数据库 -> 初始化系统 -> 登录后台 -> 添加 Token -> 添加规则 -> 启动 scan 服务 -> 拉取并过滤搜索结果 -> 人工确认或忽略 -> 导出结果。

7. 配置 token 和规则后为什么没有扫描结果？

常见原因包括：scan 服务没启动、scanner 连不上数据库、token 无效、规则没有命中、GitHub/GitLab API 网络不通、DNS 配置错误、触发平台 rate limit。应先看后端和 scanner 日志。

8. 扫描任务是手动触发还是自动循环？

新版中扫描服务会循环执行。只要 scan 服务在运行，并且存在有效 token 和规则，就会周期性扫描。旧版任务管理相关问题不适用于新版 FAQ。

9. GitHub 规则应该怎么写？

GitHub 规则可以直接使用 GitHub 搜索语法，例如：

```text
password in:file
access_token org:example
secret repo:owner/repo
api_key extension:yaml
```

规则不是只能写普通关键词，可以带 `repo:`、`org:`、`user:`、`in:file` 等限定符。

10. 一条规则可以写多个关键词吗？

单条规则建议只写一个搜索表达式。多个规则应使用批量导入能力，不要把多个无关关键词塞到一条规则里。

11. 如何减少 `.json`、`.csv`、日志文件等噪声结果？

使用 GitHub 过滤器，通过 `extension` 和 `keyword` 缩小初始搜索范围。对于宽泛规则，可增加本地 `matchPattern` 正则，要求候选内容包含更强的证据后再入库。

12. GitHub rate limit 怎么处理？

GitHub 搜索限制无法可靠绕过，也不建议通过多账号规避，存在封号风险。scanner 现在会在触发限流后自动重试该页，而不是直接丢弃，但仍建议减少规则噪声、缩小搜索范围、接受扫描延迟。

13. 自建 GitLab 能接入吗？

可以配置 GitLab Base URL。GShark 会优先尝试 GitLab 的全局代码搜索（`scope=blobs`），该功能依赖 Advanced Search/Elasticsearch —— 自建实例如果开启了该功能默认可用，GitLab.com 则只有开通 Advanced Search 的账号可用。当服务端返回不支持全局搜索时，GShark 会自动降级为爬取并搜索近期活跃的公开项目，因此没有 Advanced Search 的账号依然能拿到结果，只是覆盖范围比真正的全局搜索要窄。

14. 搜索结果可以导出吗？

可以。新版已有搜索结果导出能力，适合离线分析、归档和后续处置。

15. 遇到问题应该先提供哪些信息？

建议提供版本号、部署方式、操作系统、MySQL 版本、是否 Docker、server 日志、scanner 日志、浏览器控制台错误、相关配置截图或脱敏后的 token/rule 配置。这样比单独贴页面截图更容易定位。

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
