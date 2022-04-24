[ZH](README-ZH.md) | EN

<p align="center">
   <img alt="GgShark logo" src="https://s1.ax1x.com/2018/10/17/idhZvj.png" />
   <h3 align="center">GShark</h3>
   <p align="center">Scan for sensitive information easily and effectively.</p>
</p>

# GShark [![Go Report Card](https://goreportcard.com/badge/github.com/madneal/gshark)](https://goreportcard.com/report/github.com/madneal/gshark)   

The project is based on Go and Vue to build a management system for sensitive information detection. For the full introduction, please refer [here](https://mp.weixin.qq.com/s/Yoo1DdC2lCtqOMAreF9K0w).


# Features

* Support multi platforms, including Gitlab, Github, Searchcode
* Flexible menu and API permission setting
* Flexible rules and filter rules
* Utilize gobuster to brute force subdomain
* Easily used management system

# Quick start

![GShark](https://user-images.githubusercontent.com/12164075/114326875-58e1da80-9b69-11eb-82a5-b2e3751a2304.png)

## Deployment

It's suggested to deploy the frontend project by nginx. Place the `dist` folder under `/var/www/html`, modify the `nginx.conf` to reverse proxy the backend service. For the detailed deployment videos, refer to [bilibili](https://www.bilibili.com/video/BV1Py4y1s7ap/) or [youtube](https://youtu.be/bFrKm5t4M54). For the deployment in windows, refer to [here](https://www.bilibili.com/video/BV1CA411L7ux/).

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

The deployment work is very easy. Find the corresponding version zip file from [releases](https://github.com/madneal/gshark/releases). Unzip and copy the files inside `dist` to `/var/www/html` folder of nginx. Start the nginx and the frontend is deploy successfully.

### Web service

```
./gshark web
```

### Scan service

```
./gshark scan
```

## Development

### Server side

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



### Web side

```
cd ../web

npm install

npm run serve
```

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

### Add Token

To execute `./gshark scan`, you need to add a Github token for crawl information in github. You can generate a token in [tokens](https://github.com/settings/tokens). Most access scopes are enough. For Gitlab search, remember to add token too.

[![iR2TMt.md.png](https://s1.ax1x.com/2018/10/31/iR2TMt.md.png)](https://imgchr.com/i/iR2TMt)

## FAQ

1. Default username and password to login

gshark/gshark

2. Database initial failed

make sure the version of mysql is over 5.6. And remove the databse before initial the second time.

3. `go get ./... connection error`

It's suggested to enable goproxy(refer this [article](https://madneal.com/post/gproxy/) for golang upgrade):

```
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GO111MODULE=on
```
4. When deployed the web to `nginx`, the page was empty

try to clear the LocalStorage

## Reference

* [gin-vue-admin](https://github.com/flipped-aurora/gin-vue-admin)

## Wechat

If you would like to join wechat group, you can add my wechat `mmadneal` with the message `gshark`.

## License

[Apache License 2.0](https://github.com/madneal/gshark/blob/master/LICENSE)

## 404StarLink 2.0 - Galaxy

![](https://github.com/knownsec/404StarLink-Project/raw/master/logo.png)

GShark 是 404Team [星链计划2.0](https://github.com/knownsec/404StarLink2.0-Galaxy)中的一环，如果对 GShark 有任何疑问又或是想要找小伙伴交流，可以参考星链计划的加群方式。

- [https://github.com/knownsec/404StarLink2.0-Galaxy#community](https://github.com/knownsec/404StarLink2.0-Galaxy#community)
