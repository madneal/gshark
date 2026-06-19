<p align="center">
   <img alt="GShark logo" src="https://s1.ax1x.com/2018/10/17/idhZvj.png" />
   <h3 align="center">GShark</h3>
   <p align="center">Scan for sensitive information easily and effectively.</p>
</p>

<div align="center">
   <a href="README_CN.md">🇨🇳 中文版</a> | <strong>🇺🇸 English</strong>
</div>

# GShark [![Go Report Card](https://goreportcard.com/badge/github.com/madneal/gshark)](https://goreportcard.com/report/github.com/madneal/gshark) [![Release](https://github.com/madneal/gshark/actions/workflows/release.yml/badge.svg)](https://github.com/madneal/gshark/actions/workflows/release.yml)

GShark is a sensitive information detection and management platform. The backend is built with Go and Gin, and the current frontend is built with Vue 3, Vite, Vue Router 4, Vuex 4, and Element Plus. For the full introduction, please refer to [articles](https://mp.weixin.qq.com/mp/appmsgalbum?__biz=MzI3MjA3MTY3Mw==&action=getalbum&album_id=2376148333116850178#wechat_redirect) and [videos](https://mp.weixin.qq.com/mp/appmsgalbum?__biz=MzI3MjA3MTY3Mw==&action=getalbum&album_id=1834365721464651778#wechat_redirect). For now, all scans target public environments, not local environments.

For the usage of GShark, please refer to the [wiki](https://github.com/madneal/gshark/wiki).

# Key Features

* 🌐 Multi-Platform Support: GitHub, GitLab, Searchcode, Postman, and more
* 🔍 Flexible Rule Management: Custom scanning rules and filtering with whitelist/blacklist support
* 🔑 Fine-grained Access Control: Configurable menu and API permissions
* 🔄 Subdomain Discovery: Integrated gobuster for subdomain enumeration
* 🚀 Docker Deployment: Containerized deployment for easy setup
* 📊 Vue 3 Management Interface: Vite-powered web interface for task and result management

# Quick start

Default login after initialization:

```text
gshark / gshark
```

## Quick one-click deployment

Use one of the two quick deployment entries:

```bash
# Option 1: Docker quick. Build and start mysql/server/web in the background.
./scripts/quick-docker.sh

# Start the scan container too.
./scripts/quick-docker.sh --with-scan
```

```bash
# Option 2: Release quick. Download the matching release package,
# configure Nginx, and start the gshark backend in the background.
./scripts/quick-release.sh

# Or deploy from a local release zip.
./scripts/quick-release.sh --file ./gshark_linux_amd64.zip
```

## Docker Deployment

```
# Clone the repository
git clone https://github.com/madneal/gshark

cd gshark

# Build and start the containers
./scripts/quick-docker.sh
```

> [!IMPORTANT]
> Before the MySQL database initial, the scanner container will exit. Need to restart the scanner after the MySQL database initial.

## Local Deployment 

```bash  
# Clone the repository  
git clone https://github.com/madneal/gshark.git  
cd gshark  

# Run Release quick to download the release package, configure Nginx,
# and start the backend.
./scripts/quick-release.sh
```

## Manual Deployment

### Requirements

* Nginx
* MySQL **8.0+**
* Go **1.25+** for building the backend
* Node.js **20+** and npm for building the frontend

It is recommended to deploy the frontend with Nginx. Build the Vite project, place the generated `web/dist` files in `/var/www/html`, and configure Nginx to reverse proxy `/api/` to the backend service. For detailed deployment tutorials, you can watch videos on [bilibili](https://www.bilibili.com/video/BV1Py4y1s7ap/) or [youtube](https://youtu.be/bFrKm5t4M54). For deployment on Windows, refer to [this link](https://www.bilibili.com/video/BV1CA411L7ux/).

### Nginx

Can use `nginx -t` to locate the `nginx.conf` file, then modify the `nginx.conf`:

```
// config the user accoring to your need
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

The deployment work is straightforward. Find the corresponding version zip file from [releases](https://github.com/madneal/gshark/releases).

Unzip and copy the files inside `dist` to `/var/www/html` folder of Nginx. 

```
unzip gshark*.zip
cd gshark*
mv dist/* /var/www/html/
# for Mac
mv dist/* /usr/local/www/html/
```

Start the Nginx and the Front-End is deployed successfully.

> [!TIP]
> If you installed Nginx by Homebrew, you need to stop Nginx by:
> ```shell
> brew services stop nginx
> ```
> Start Nginx for Ubuntu:
> ```shell
> systemctl start nginx
> ```

### Server service

```shell
./gshark serve
```

Initially, copy `config-temp.yaml` to `config.yaml` and update it for your environment. After that, you can run the `gshark` binary file directly. Then, access `localhost:8080` for local deployment.

If you haven't initialized the database before, you will be redirected to the database initialization page first.

<img width="936" alt="image" src="https://github.com/user-attachments/assets/dfa7e53e-dc4a-4697-831f-a4f4f3810c3c">

### Scan service

```shell
./gshark scan
```

For the scan service, it's necessary to config the corresponding rules. For example, GitHub or Gitlab rules.

### Incremental Deployment

For the incremental deployment, [sql.md](https://github.com/madneal/gshark/blob/master/sql.md) should be executed for the corresponding database operations.

## Development

### Server

```shell
git clone https://github.com/madneal/gshark.git
cd gshark/server
go mod tidy
cp config-temp.yaml config.yaml
go build
```

Run the web server:

```shell
go build
./gshark serve 
```

Or

```shell
go run main.go serve
```

Run the scan task:

```shell
go build
./gshark scan 
```

Or

```shell
go run main.go scan
```

> [!NOTE]
> On macOS ARM, CPU percentage collection in the server-info page depends on cgo. Use `CGO_ENABLED=1` when running or building the backend if you need CPU usage percentages:
>
> ```shell
> CGO_ENABLED=1 go run main.go serve
> ```

### Web 

```
cd ../web

npm install

npm run serve
```

## Usage
### Add Token

#### GitHub

To execute the scan task for GitHub, you need to add a GitHub token for crawl information in GitHub. You can generate a token in [tokens](https://github.com/settings/tokens). Most access scopes are enough. For the GitLab search, remember to add a token too.

[![iR2TMt.md.png](https://s1.ax1x.com/2018/10/31/iR2TMt.md.png)](https://imgchr.com/i/iR2TMt)

### Rule Configuration

For the Github or Gitlab rule, the rule will be matched by the syntax in the corresponding platforms. Directly, you config what you search at GitHub. You can download the rule import template CSV file, then batch import rules.

<img width="572" alt="image" src="https://user-images.githubusercontent.com/12164075/212504597-3e1ad5bd-bacf-433e-83e8-08de7eee6509.png">


### Filter Configuration

Filter is only addressed to GitHub search now. There are three classes of filters, including `extension`, `keyword`, `sec_keyword`. For `extension` and `keyword`, they can used for blacklist or whitelist.

For more information, you can refer to this [video](https://www.bilibili.com/video/BV1aG4y1c72N/?vd_source=ef4657ebf0549af8755f75118b6e81bb).

## Configuration

You are supposed to rename `config-temp.yaml` to `config.yaml` and config the database information and other information according to your environment.

### GitLab Base Url

<img width="363" alt="image" src="https://user-images.githubusercontent.com/12164075/203898719-1ce66395-083d-4226-937f-b6eed859addc.png">


## FAQ

1. Does GShark scan local code or public platforms?

GShark is designed to scan public environments, not local source trees. GitHub scanning is based on the GitHub Search API, and GitLab scanning depends on GitLab search. Whether private repositories can be scanned depends on the platform API and token permissions.

2. What is the recommended deployment method?

New users should prefer the quick scripts:

```shell
./scripts/quick-docker.sh
./scripts/quick-docker.sh --with-scan
./scripts/quick-release.sh
```

Manual deployment is useful when you need custom Nginx, MySQL, or backend configuration.

3. What are the deployment requirements?

MySQL 8.0+ is required. Manual builds require Go 1.25+, Node.js 20+, npm, and Nginx. For Docker deployment, prefer the compose file and quick scripts provided by this repository to avoid configuration drift from older tutorials.

4. What is the default account after initialization?

The default account is `gshark / gshark`. Change the password immediately after deploying to a production environment.

5. Why did the scanner not start or produce results after Docker deployment?

The scanner depends on database initialization. Before MySQL is initialized, the scanner container may exit. Restart the scanner after database initialization. When troubleshooting, check the scanner/server container logs first instead of only checking the web page.

6. What is the core GShark workflow?

The basic workflow is: configure the database -> initialize the system -> log in -> add tokens -> add rules -> start the scan service -> fetch search results -> filter or run secondary filtering -> manually confirm or ignore findings -> export results.

7. Why are there no scan results after configuring tokens and rules?

Common causes include: the scan service is not running, the scanner cannot connect to the database, the token is invalid, no rule matched, the GitHub/GitLab API is unreachable, DNS is misconfigured, or the platform rate limit was triggered. Check backend and scanner logs first.

8. Are scans manually triggered or automatically repeated?

In the current version, the scan service runs in a loop. As long as the scan service is running and valid tokens and rules exist, scans will run periodically. Old task-management issues do not apply to the current FAQ.

9. How should GitHub rules be written?

GitHub rules can directly use GitHub search syntax, for example:

```text
password in:file
access_token org:example
secret repo:owner/repo
api_key extension:yaml
```

Rules are not limited to plain keywords. You can use qualifiers such as `repo:`, `org:`, `user:`, and `in:file`.

10. Can one rule contain multiple keywords?

One rule should normally contain one search expression. Use batch import for multiple rules instead of placing unrelated keywords into a single rule.

11. How can I reduce noisy results from `.json`, `.csv`, log files, and similar files?

Use filters. Filters are focused on GitHub search and support types such as `extension`, `keyword`, and `sec_keyword`. Extension filtering happens before results are stored. Secondary filtering uses secondary keywords to refine results; it is not the same feature as extension filtering.

12. How should GitHub rate limits be handled?

GitHub search limits cannot be reliably bypassed, and using multiple accounts to avoid them is not recommended because it may risk account bans. A better approach is to reduce noisy rules, narrow the search scope, accept scan delays, and check whether failed tasks are retried.

13. Can GShark connect to self-hosted GitLab?

Yes, by configuring the GitLab Base URL. However, the self-hosted GitLab instance must support code search/indexing. If global search is disabled on the server, GShark cannot bypass that platform limitation.

14. Can search results be exported?

Yes. Current versions include search result export, which is useful for offline analysis, archiving, and follow-up handling.

15. What information should I provide when reporting a problem?

Provide the version, deployment method, operating system, MySQL version, whether Docker is used, server logs, scanner logs, browser console errors, relevant screenshots, and redacted token/rule configuration. This is more useful than a page screenshot alone.

## Resources 

### Articles

* [多平台的敏感信息监测工具-GShark](https://mp.weixin.qq.com/s?__biz=MzI3MjA3MTY3Mw==&mid=2247484283&idx=1&sn=3232df7d321c0f62ce61b7e6368204ad&chksm=eb396deddc4ee4fb0c825a378c085223b87fc45f05648d46e7bdc24a03fb83ad6c7ade414df7#rd)
* [GShark-监测你的 Github 敏感信息泄露](https://mp.weixin.qq.com/s?__biz=MzI3MjA3MTY3Mw==&mid=2247483770&idx=1&sn=9f02c2803e1c946e8c23b16ff3eba757&chksm=eb396fecdc4ee6fa2f378e846f354f45acf6e6f540cfd54190e9353df47c7707e3a2aadf714f&token=1578822041&lang=zh_CN#rd)


### Videos

* [GShark v1.5.0 版本及 Docker 使用指南](https://www.bilibili.com/video/BV1oUe3eBEMz/)
* [GShark v1.3.0 版本支持 Docker](https://www.bilibili.com/video/BV1BH4y1C7Ga/)
* [GShark 支持多种规则类型以及规则配置建议](https://www.bilibili.com/video/BV1uY4y177SX) 
* [批量导入规则](https://mp.weixin.qq.com/s?__biz=MzI3MjA3MTY3Mw==&mid=2247484546&idx=1&sn=818915279c5199457340ade89d6cbd54&chksm=eb396a14dc4ee302039bcb1474380a6049dba84370345b7813049aa8feb49a98f89d47ec5d5b#rd)
* [GShark部署](https://mp.weixin.qq.com/s?__biz=MzI3MjA3MTY3Mw==&mid=2247484487&idx=1&sn=78f942ccf6861f433fc7f4a60564441c&chksm=eb396ad1dc4ee3c7505362da243433e54a2b558c96fbbb50f8b6cea87d1f9bc920b249b72705#rd)
* [windows 部署](https://mp.weixin.qq.com/s?__biz=MzI3MjA3MTY3Mw==&mid=2247484289&idx=1&sn=2b0f1c38b88c924ad514fb64b559b784&chksm=eb396d17dc4ee4018573dde6c3bfce83903c86034403539eaf1b87b89c4a4dd44f957a308818#rd)
* [GShark v1.0.2 版本发布](https://www.bilibili.com/video/BV1Zx4y1G7FX/)
* [GShark v1.1.0 更新内容介绍](https://www.bilibili.com/video/BV1aG4y1c72N/)


## License

[Apache License 2.0](https://github.com/madneal/gshark/blob/master/LICENSE)

## 404StarLink 2.0 - Galaxy

![](https://github.com/knownsec/404StarLink-Project/raw/master/logo.png)

GShark 是 404Team [星链计划2.0](https://github.com/knownsec/404StarLink2.0-Galaxy)中的一环，如果对 GShark 有任何疑问又或是想要找小伙伴交流，可以参考星链计划的加群方式。

- [https://github.com/knownsec/404StarLink2.0-Galaxy#community](https://github.com/knownsec/404StarLink2.0-Galaxy#community)
