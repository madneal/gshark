<p align="center">
   <img alt="GShark logo" src="https://s1.ax1x.com/2018/10/17/idhZvj.png" />
   <h3 align="center">GShark</h3>
   <p align="center">轻松有效地扫描敏感信息。</p>
</p>

<div align="center">
   <strong>🇨🇳 中文版</strong> | <a href="README.md">🇺🇸 English</a>
</div>

# GShark [![Go Report Card](https://goreportcard.com/badge/github.com/madneal/gshark)](https://goreportcard.com/report/github.com/madneal/gshark) [![Release](https://github.com/madneal/gshark/actions/workflows/release.yml/badge.svg)](https://github.com/madneal/gshark/actions/workflows/release.yml)

GShark 是一个敏感信息检测和管理平台。后端基于 Go 和 Gin 构建，当前前端基于 Vue 3、Vite、Vue Router 4、Vuex 4 和 Element Plus 构建。完整介绍请参考[文章](https://mp.weixin.qq.com/mp/appmsgalbum?__biz=MzI3MjA3MTY3Mw==&action=getalbum&album_id=2376148333116850178#wechat_redirect)和[视频](https://mp.weixin.qq.com/mp/appmsgalbum?__biz=MzI3MjA3MTY3Mw==&action=getalbum&album_id=1834365721464651778#wechat_redirect)。目前，所有扫描仅针对公共环境，不针对本地环境。

关于 GShark 的使用，请参考 [wiki](https://github.com/madneal/gshark/wiki)。

# 主要特性

* 🌐 多平台支持：GitHub、GitLab、Searchcode、Postman 等
* 🔍 灵活的规则管理：自定义扫描规则和过滤，支持白名单/黑名单
* 🔑 细粒度访问控制：可配置的菜单和 API 权限
* 🔄 子域名发现：集成 gobuster 进行子域名枚举
* 🚀 Docker 部署：容器化部署，易于设置
* 📊 Vue 3 管理界面：基于 Vite 的 Web 界面，用于任务和结果管理

# 快速开始

初始化后的默认登录账号：

```text
gshark / gshark
```

## Docker 部署

```
# 克隆仓库
git clone https://github.com/madneal/gshark

cd gshark

# 构建并启动容器
docker-compose build && docker-compose up 
```

> [!IMPORTANT]
> 在 MySQL 数据库初始化之前，扫描器容器会退出。需要在 MySQL 数据库初始化后重启扫描器。

## 本地部署

```bash
# 克隆仓库
git clone https://github.com/madneal/gshark.git
cd gshark

# 添加执行权限（只需要执行一次）
chmod +x deployment.sh

# 执行脚本构建并启动服务
./deployment.sh
```

## 手动部署

### 环境要求

* Nginx
* MySQL **8.0+**
* Go **1.25+**，用于构建后端
* Node.js **20+** 和 npm，用于构建前端

建议使用 Nginx 部署前端。构建 Vite 项目后，将生成的 `web/dist` 文件放置在 `/var/www/html` 中，并配置 Nginx 将 `/api/` 反向代理到后端服务。详细的部署教程可以观看 [bilibili](https://www.bilibili.com/video/BV1Py4y1s7ap/) 或 [youtube](https://youtu.be/bFrKm5t4M54) 上的视频。Windows 部署请参考[此链接](https://www.bilibili.com/video/BV1CA411L7ux/)。

### Nginx

可以使用 `nginx -t` 定位 `nginx.conf` 文件，然后修改 `nginx.conf`：

```
// 根据您的需要配置用户
user  www www;
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
            autoindex on;
            root   html;
            index  index.html index.htm;
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
            root   html;
        }
    }
    include servers/*;
}

```

部署工作很简单。从 [releases](https://github.com/madneal/gshark/releases) 找到对应版本的 zip 文件。

解压并将 `dist` 内的文件复制到 Nginx 的 `/var/www/html` 文件夹。

```
unzip gshark*.zip
cd gshark*
mv dist/* /var/www/html/
# Mac 系统
mv dist/* /usr/local/www/html/
```

启动 Nginx，前端部署成功。

> [!TIP]
> 如果您通过 Homebrew 安装了 Nginx，需要停止 Nginx：
> ```shell
> brew services stop nginx
> ```
> Ubuntu 启动 Nginx：
> ```shell
> systemctl start nginx
> ```

### 服务器服务

```shell
./gshark serve
```

初始时，将 `config-temp.yaml` 复制为 `config.yaml`，并根据环境修改配置。之后，您可以直接运行 `gshark` 二进制文件。然后，访问 `localhost:8080` 进行本地部署。

如果您之前没有初始化数据库，您将首先被重定向到数据库初始化页面。

<img width="936" alt="image" src="https://github.com/user-attachments/assets/dfa7e53e-dc4a-4697-831f-a4f4f3810c3c">

### 扫描服务

```shell
./gshark scan
```

对于扫描服务，需要配置相应的规则。例如，GitHub 或 Gitlab 规则。

### 增量部署

对于增量部署，应该执行 [sql.md](https://github.com/madneal/gshark/blob/master/sql.md) 进行相应的数据库操作。

## 开发

### 服务器

```shell
git clone https://github.com/madneal/gshark.git
cd gshark/server
go mod tidy
cp config-temp.yaml config.yaml
go build
```

运行 Web 服务器：

```shell
go build
./gshark serve 
```

或者

```shell
go run main.go serve
```

运行扫描任务：

```shell
go build
./gshark scan 
```

或者

```shell
go run main.go scan
```

> [!NOTE]
> 在 macOS ARM 上，服务器状态页面的 CPU 百分比采集依赖 cgo。如果需要显示 CPU 使用率，请在运行或构建后端时启用 `CGO_ENABLED=1`：
>
> ```shell
> CGO_ENABLED=1 go run main.go serve
> ```

### Web 前端

```
cd ../web

npm install

npm run serve
```

## 使用方法
### 添加 Token

#### GitHub

要执行 GitHub 的扫描任务，您需要添加 GitHub token 来爬取 GitHub 中的信息。您可以在 [tokens](https://github.com/settings/tokens) 中生成 token。大多数访问范围就足够了。对于 GitLab 搜索，记得也要添加 token。

[![iR2TMt.md.png](https://s1.ax1x.com/2018/10/31/iR2TMt.md.png)](https://imgchr.com/i/iR2TMt)

### 规则配置

对于 Github 或 Gitlab 规则，规则将按照相应平台的语法进行匹配。您可以直接配置在 GitHub 中搜索的内容。您可以下载规则导入模板 CSV 文件，然后批量导入规则。

<img width="572" alt="image" src="https://user-images.githubusercontent.com/12164075/212504597-3e1ad5bd-bacf-433e-83e8-08de7eee6509.png">

### 过滤器配置

过滤器目前仅针对 GitHub 搜索。有三类过滤器，包括 `extension`、`keyword`、`sec_keyword`。对于 `extension` 和 `keyword`，它们可以用于黑名单或白名单。

更多信息，您可以参考这个[视频](https://www.bilibili.com/video/BV1aG4y1c72N/?vd_source=ef4657ebf0549af8755f75118b6e81bb)。

## 配置

您应该将 `config-temp.yaml` 复制为 `config.yaml`，并根据您的环境配置数据库信息和其他信息。

### GitLab 基础 URL

<img width="363" alt="image" src="https://user-images.githubusercontent.com/12164075/203898719-1ce66395-083d-4226-937f-b6eed859addc.png">

## 常见问题

1. 默认登录用户名和密码

gshark/gshark

2. 数据库初始化失败

确保 MySQL 版本为 8.0 或更高。并在第二次初始化前删除数据库。

3. `go get ./... connection error`

建议启用 GOPROXY（参考这篇[文章](https://madneal.com/post/gproxy/)进行 golang 升级）：

```
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GO111MODULE=on
```

4. 将 Web 部署到 `nginx` 时，页面为空

尝试清除 LocalStorage，并确认 Nginx 的 `/api/` 反向代理指向后端服务。

5. macOS ARM 上服务器状态页面的 CPU 使用率显示为 0

运行或构建后端时启用 cgo：

```shell
CGO_ENABLED=1 go run main.go serve
```

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
