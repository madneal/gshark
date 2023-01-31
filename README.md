[ZH](README-ZH.md) | EN

<p align="center">
   <img alt="GgShark logo" src="https://s1.ax1x.com/2018/10/17/idhZvj.png" />
   <h3 align="center">GShark</h3>
   <p align="center">Scan for sensitive information easily and effectively.</p>
</p>

# GShark [![Go Report Card](https://goreportcard.com/badge/github.com/madneal/gshark)](https://goreportcard.com/report/github.com/madneal/gshark)  [![Release](https://github.com/madneal/gshark/actions/workflows/release.yml/badge.svg)](https://github.com/madneal/gshark/actions/workflows/release.yml)

The project is based on Go and Vue to build a management system for sensitive information detection. For the full introduction, please refer to [articles](https://mp.weixin.qq.com/mp/appmsgalbum?__biz=MzI3MjA3MTY3Mw==&action=getalbum&album_id=2376148333116850178#wechat_redirect) and [videos](https://mp.weixin.qq.com/mp/appmsgalbum?__biz=MzI3MjA3MTY3Mw==&action=getalbum&album_id=1834365721464651778#wechat_redirect). For now, all the scans are only targeted to the public environments, not local environments.

# Features

* Support multi platforms, including Gitlab, Github, Searchcode, Postman
* Flexible menu and API permission setting
* Flexible rules and filter rules
* Utilize gobuster to brute force subdomain
* Easily used management system

# Quick start

![GShark](https://user-images.githubusercontent.com/12164075/114326875-58e1da80-9b69-11eb-82a5-b2e3751a2304.png)

## Deployment

### Requirements

* Nginx
* MySQL(version above **8.0**)

It's suggested to deploy the frontend project by nginx. Place the `dist` folder under `/var/www/html`, modify the `nginx.conf` to reverse proxy the backend service. For the detailed deployment videos, refer to [bilibili](https://www.bilibili.com/video/BV1Py4y1s7ap/) or [youtube](https://youtu.be/bFrKm5t4M54). For the deployment in windows, refer to [here](https://www.bilibili.com/video/BV1CA411L7ux/).

### Nginx

Modify the `nginx.conf`:

```
location /api/ {
   proxy_set_header Host $http_host;
   proxy_set_header X-Real-IP $remote_addr;
   proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
   proxy_set_header X-Forwarded-Proto $scheme;
   rewrite ^/api/(.*)$ /$1 break;
   proxy_pass http://127.0.0.1:8888;
}
```

The deployment work is very easy. Find the corresponding version zip file from [releases](https://github.com/madneal/gshark/releases). 

Unzip and copy the files inside `dist` to `/var/www/html` folder of Nginx. 

```
unzip gshark*.zip
```

Start the Nginx and the Front-End is deployed successfully.

### Incremental Deployment

For the incremental deployment, [sql.md](https://github.com/madneal/gshark/blob/master/sql.md) should be executed for the corresponding database operations.

### Web service

For the first time, you need to rename `config-temp.yaml` to `config.yaml`.

```
./gshark web
```

### Scan service

```
./gshark scan
```

For the scan service, it's necessary to config the corresponding rules. For example, Github or Gitlab rules.

<<<<<<< HEAD
## Development

### Server

``` 
git clone https://github.com/madneal/gshark.git

cd server

go mod tidy

mv config-temp.yaml config.yaml

go build

./gshark web
```

If you want to set up the scan service, please run:

```
./gshark scan
```



### Web 

```
cd ../web

npm install

npm run serve
```

=======
>>>>>>> 18b49ff51aa3e2394c3f907af876d4620d17cd80
## Run

```
USAGE:
   gshark [global options] command [command options] [arguments...]

COMMANDS:
     web      Startup a web Service
     scan     Start to scan github leak info
     help, h  Show a list of commands or help for one command

GLOBAL OPTIONS:
   --debug, -d             Debug Mode
   --host value, -H value  web listen address (default: "0.0.0.0")
   --port value, -p value  web listen port (default: 8000)
   --time value, -t value  scan interval(second) (default: 900)
   --help, -h              show help
   --version, -v           print the version
```

## Usage
### Add Token

To execute `./gshark scan`, you need to add a GitHub token for crawl information in GitHub. You can generate a token in [tokens](https://github.com/settings/tokens). Most access scopes are enough. For GitLab search, remember to add token too.

[![iR2TMt.md.png](https://s1.ax1x.com/2018/10/31/iR2TMt.md.png)](https://imgchr.com/i/iR2TMt)

### Rule Configuration

For the Github or Gitlab rule, the rule will be matched by the syntax in the corresponding platforms. Directly, you config what you search at GitHub. You can download the rule import template csv file, then batch import rules.

<img width="572" alt="image" src="https://user-images.githubusercontent.com/12164075/212504597-3e1ad5bd-bacf-433e-83e8-08de7eee6509.png">


### Filter Configuration

Filter is only addressed to GitHub search now. There are three classes of filter, include `extension`, `keyword`, `sec_keyword`. For `extension` and `keyword`, they can used for blacklist or whitelist.

For more information, you can refer to this [video](https://www.bilibili.com/video/BV1aG4y1c72N/?vd_source=ef4657ebf0549af8755f75118b6e81bb).

## Configuration

You are supposed to rename `config-temp.yaml` to `config.yaml` and config the database information and other information according to your environment.

### GitLab Base Url

<img width="363" alt="image" src="https://user-images.githubusercontent.com/12164075/203898719-1ce66395-083d-4226-937f-b6eed859addc.png">


## FAQ

1. Default username and password to login

gshark/gshark

2. Database initial failed

make sure the version of mysql is over 5.6. And remove the database before initial the second time.

3. `go get ./... connection error`

It's suggested to enable GOPROXY(refer this [article](https://madneal.com/post/gproxy/) for golang upgrade):

```
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GO111MODULE=on
```
4. When deployed the web to `nginx`, the page was empty

try to clear the LocalStorage

## Resources 

### Articles

* [多平台的敏感信息监测工具-GShark](https://mp.weixin.qq.com/s?__biz=MzI3MjA3MTY3Mw==&mid=2247484283&idx=1&sn=3232df7d321c0f62ce61b7e6368204ad&chksm=eb396deddc4ee4fb0c825a378c085223b87fc45f05648d46e7bdc24a03fb83ad6c7ade414df7#rd)
* [GShark-监测你的 Github 敏感信息泄露](https://mp.weixin.qq.com/s?__biz=MzI3MjA3MTY3Mw==&mid=2247483770&idx=1&sn=9f02c2803e1c946e8c23b16ff3eba757&chksm=eb396fecdc4ee6fa2f378e846f354f45acf6e6f540cfd54190e9353df47c7707e3a2aadf714f&token=1578822041&lang=zh_CN#rd)


### Videos

* [GShark 支持多种规则类型以及规则配置建议](https://www.bilibili.com/video/BV1uY4y177SX) 
* [批量导入规则](https://mp.weixin.qq.com/s?__biz=MzI3MjA3MTY3Mw==&mid=2247484546&idx=1&sn=818915279c5199457340ade89d6cbd54&chksm=eb396a14dc4ee302039bcb1474380a6049dba84370345b7813049aa8feb49a98f89d47ec5d5b#rd)
* [GShark部署](https://mp.weixin.qq.com/s?__biz=MzI3MjA3MTY3Mw==&mid=2247484487&idx=1&sn=78f942ccf6861f433fc7f4a60564441c&chksm=eb396ad1dc4ee3c7505362da243433e54a2b558c96fbbb50f8b6cea87d1f9bc920b249b72705#rd)
* [windows 部署](https://mp.weixin.qq.com/s?__biz=MzI3MjA3MTY3Mw==&mid=2247484289&idx=1&sn=2b0f1c38b88c924ad514fb64b559b784&chksm=eb396d17dc4ee4018573dde6c3bfce83903c86034403539eaf1b87b89c4a4dd44f957a308818#rd)
* [GShark v1.0.2 版本发布](https://www.bilibili.com/video/BV1Zx4y1G7FX/)
* [GShark v1.1.0 更新内容介绍](https://www.bilibili.com/video/BV1aG4y1c72N/)


## Reference

* [gin-vue-admin](https://github.com/flipped-aurora/gin-vue-admin)

## Wechat

If you would like to join WeChat group, you can add my WeChat `mmadneal` with the message `gshark`. But I would strongly suggest you to raise an issue directly instead of just asking in WeChat.

## License

[Apache License 2.0](https://github.com/madneal/gshark/blob/master/LICENSE)

## 404StarLink 2.0 - Galaxy

![](https://github.com/knownsec/404StarLink-Project/raw/master/logo.png)

GShark 是 404Team [星链计划2.0](https://github.com/knownsec/404StarLink2.0-Galaxy)中的一环，如果对 GShark 有任何疑问又或是想要找小伙伴交流，可以参考星链计划的加群方式。

- [https://github.com/knownsec/404StarLink2.0-Galaxy#community](https://github.com/knownsec/404StarLink2.0-Galaxy#community)
