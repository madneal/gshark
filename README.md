<p align="center">
   <img alt="GgShark logo" src="https://s1.ax1x.com/2018/10/17/idhZvj.png" />
   <h3 align="center">GShark</h3>
   <p align="center">Scan for sensitive information in Github easily and effectively.</p>
</p>

# GShark [![Go Report Card](https://goreportcard.com/badge/github.com/madneal/gshark)](https://goreportcard.com/report/github.com/madneal/gshark)   

The project is based on golang with AdminLTE to build a management system to manage the Github search results. Github API is utilized to crawl the related results according to key words and some rules. It proves to be a proper way to detect the information related to your company.:rocket::rocket::rocket: For a detailed introduction, you can refer [here](https://mp.weixin.qq.com/s?__biz=MzI3MjA3MTY3Mw==&mid=2247483770&idx=1&sn=9f02c2803e1c946e8c23b16ff3eba757&chksm=eb396fecdc4ee6fa2f378e846f354f45acf6e6f540cfd54190e9353df47c7707e3a2aadf714f&token=1263666156&lang=zh_CN#rd).

![ezgif com-optimize](https://user-images.githubusercontent.com/12164075/47776907-72db2a00-dd2e-11e8-9862-db4aa5c458ff.gif)

## Features

* Support multi platform, including Gitlab, Github, Searchcode
* Support search keyword in Huawei app store
* flexible rules
* utilize gobuster to brute force subdomain

## Quick start

```
git clone https://github.com/madneal/gshark

go get ./...

go build main.go

# check the config
mv app-template.ini app.ini 

# start web service
./main web 

# start crawler
./main scan
```

## Config

The configuration can be set according to app-template.ini. You should rename it to app.ini to config the project.

```
HTTP_HOST = 127.0.0.1
HTTP_PORT = 8000
MAX_INDEXERS = 2
DEBUG_MODE = true
REPO_PATH = repos
MAX_Concurrency_REPOS = 5

; server酱配置口令
SCKEY =
; gobuster file path
gobuster_path =
; gobuster subdomain wordlist file path
subdomain_wordlist_file =

[database]
;support sqlite3, mysql, postgres
DB_TYPE = sqlite
HOST = 127.0.0.1
PORT = 3306
NAME = misec
USER = root
PASSWD = 
SSL_MODE = disable
;the path to store the database file of sqlite3
PATH = 
```

## Before Running

* Make sure you have installed dependencies, suggest to use go mod
* Make sure the app.ini in config folder, you can rename app-template.ini to app.ini
* Make sure that you have config and set database correctly
* Make sure that you have config [policy](https://github.com/madneal/gshark/blob/master/conf/policy.csv) properly
* Make sure that you have config corresponding tokens for github or gitlab

## Run

You should build the `main.go` file firstly with the command `go build main.go`.
```
USAGE:
   main [global options] command [command options] [arguments...]

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

### Initial Running

If it's the first time to run, there are some [initial works](https://github.com/madneal/gshark/blob/0ea3365f88e012df3fef1079df04a4f4b266319d/models/models.go#L31) will be finished automatically.

* [Init Rules](https://github.com/madneal/gshark/blob/0ea3365f88e012df3fef1079df04a4f4b266319d/models/models.go#L98)
* [Init admin](https://github.com/madneal/gshark/blob/0ea3365f88e012df3fef1079df04a4f4b266319d/models/models.go#L117)
* [Init database](https://github.com/madneal/gshark/blob/0ea3365f88e012df3fef1079df04a4f4b266319d/models/models.go#L47)

### Add Token

To execute `main scan`, you need to add a Github token for crawl information in github. You can generate a token in [tokens](https://github.com/settings/tokens). Most access scopes are enough. For Gitlab search, remember to add token too.

[![iR2TMt.md.png](https://s1.ax1x.com/2018/10/31/iR2TMt.md.png)](https://imgchr.com/i/iR2TMt)

### Add notification

Now support notification by `server 酱`. Set the config of `SCKEY` in `app.ini` file.

## FAQ

1. Access web service 403 forbidden

Access to http://127.0.0.1/admin/login

2. Default username and password

gshark/gshark

3. `# github.com/mattn/go-sqlite3
exec: "gcc": executable file not found in %PATH%`

https://github.com/mattn/go-sqlite3/issues/435#issuecomment-314247676

4. `go get ./... connection error`

It's suggested to enable goproxy(refer this [article](https://madneal.com/post/gproxy/) for golang upgrade):

```
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GO111MODULE=on
```

## Reference

* [x-patrol](https://github.com/MiSecurity/x-patrol)
* [authz](https://github.com/go-macaron/authz)
* [macaron](https://github.com/go-macaron/macaron)

## Wechat

If you would like to join wechat group, you can add my wechat `mmadneal` with the message `gshark`.

## License

[Apache License 2.0](https://github.com/madneal/gshark/blob/master/LICENSE)


