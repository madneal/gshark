中文 | [英文](README.md)
<p align="center">
   <img alt="GgShark logo" src="https://s1.ax1x.com/2018/10/17/idhZvj.png" />
   <h3 align="center">GShark</h3>
   <p align="center">快速有效地扫描敏感信息</p>
</p>

# GShark [![Go Report Card](https://goreportcard.com/badge/github.com/madneal/gshark)](https://goreportcard.com/report/github.com/madneal/gshark)   

项目基于 Go 以及 Vue 搭建的敏感信息检测管理系统。关于的完整介绍请参考[这里](https://mp.weixin.qq.com/s/Yoo1DdC2lCtqOMAreF9K0w)。对于项目的详细介绍，可以参考资源里面的视频或者文章链接。目前，本项目针对的都是公开环境而不是本地环境。

# 特性

* 支持多个搜索平台，包括 Github，Gitlab，Searchcode
* 灵活的菜单以及 API 权限管理
* 灵活的规则以及过滤规则设置
* 支持 gobuster 作为子域名爆破的支持
* 方便易用

# 快速开始

![GShark](https://user-images.githubusercontent.com/12164075/114326875-58e1da80-9b69-11eb-82a5-b2e3751a2304.png)

## 部署

### 前端项目部署

建议通过 nginx 部署前端项目。 将 `dist` 文件夹放在 `/var/www/html`下，修改 `nginx.conf` 来反向代理后端服务。在[bilibili](https://www.bilibili.com/video/BV1Py4y1s7ap/) 和 [youtube](https://youtu.be/bFrKm5t4M54) 可以参考部署视频教程。 Windows 下的部署请参考[这里](https://www.bilibili.com/video/BV1CA411L7ux/)。

```
location /api/ {
proxy_set_header Host $http_host;
proxy_set_header  X-Real-IP $remote_addr;
proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
proxy_set_header X-Forwarded-Proto $scheme;
rewrite ^/api/(.*)$ /$1 break;
proxy_pass http://127.0.0.1:8888;
}
```

部署工作非常简单。 从 [releases](https://github.com/madneal/gshark/releases) 中找到对应的版本 zip 文件。 解压后得将 `dist` 中的文件复制到  `/var/www/html` 文件夹中。

### 后端部署

后端项目无须部署，直接在文件夹内启动运行即可。web 服务和 scan 服务分别是两个独立的服务，需要独立运行。

```
./gshark web
```

### 扫描服务

```
./gshark scan
```

## 开发

### 服务端

``` 
git clone https://github.com/madneal/gshark.git

cd server

go mod tidy

mv config-temp.yaml config.yaml

go build

./gshark web
```

启动扫描服务：

```
./gshark scan
```



### Web 端

```
cd ../web

npm install

npm run serve
```

## 运行

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

### 添加 Token

执行扫描任务需要在 Github 的 Github token。你可以在 [tokens](https://github.com/settings/tokens) 中生成令牌，只需要 public_repo 的读权限即可。对于 Gitlab 扫描，请记得添加令牌。

[![iR2TMt.md.png](https://s1.ax1x.com/2018/10/31/iR2TMt.md.png)](https://imgchr.com/i/iR2TMt)

## FAQ

1. 默认登录的用户名密码（**及时修改**）：

gshark/gshark

2. 数据库初始化失败

确保数据库 mysql 版本高于 5.6。如果是第二次初始化的时候记得移除第一次初始化产生的实例。

3. `go get ./... connection error`

[使用 GOPROXY](https://madneal.com/post/gproxy/:

```
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GO111MODULE=on
```
4. 部署前端项目后，页面空白

尝试清除 LocalStorage

## Reference

* [gin-vue-admin](https://github.com/flipped-aurora/gin-vue-admin)

## Wechat

如果您想加入微信群，可以添加我的微信 `mmadneal`，并留言 `gshark`。

## License

[Apache License 2.0](https://github.com/madneal/gshark/blob/master/LICENSE)

## 404StarLink 2.0 - Galaxy

![](https://github.com/knownsec/404StarLink-Project/raw/master/logo.png)

GShark 是 404Team [星链计划2.0](https://github.com/knownsec/404StarLink2.0-Galaxy)中的一环，如果对 GShark 有任何疑问又或是想要找小伙伴交流，可以参考星链计划的加群方式。

- [https://github.com/knownsec/404StarLink2.0-Galaxy#community](https://github.com/knownsec/404StarLink2.0-Galaxy#community)
