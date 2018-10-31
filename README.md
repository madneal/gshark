<p align="center">
   <img alt="GgShark logo" src="https://s1.ax1x.com/2018/10/17/idhZvj.png" />
   <h3 align="center">GShark</h3>
   <p align="center">Scan for sensitive information in Github easily and effectively.</p>
</p>

# GShark [![Go Report Card](https://goreportcard.com/badge/github.com/neal1991/gshark)](https://goreportcard.com/report/github.com/neal1991/gshark)

The project is based on golang with AdminLTE to build a management system to manage the Github search results. Github API is utilized to scawl the related results according to key words and some rules. It proves to be a proper way to detect the information related to your company.:rocket::rocket::rocket:

## Requirements

* go version 1.10+
* the project should be placed in GOPATH

## Config

The configuration can be set according to app-template.ini. You should rename it to app.ini to config the project.

```
HTTP_HOST = 127.0.0.1
HTTP_PORT = 8000
MAX_INDEXERS = 2
DEBUG_MODE = true
REPO_PATH = repos
MAX_Concurrency_REPOS = 5

[database]
;support sqlite3, mysql, postgres
DB_TYPE = sqlite
HOST = 127.0.0.1
PORT = 3306
NAME = misec
USER = root
PASSWD = 
SSL_MODE = disable
PATH = 
```

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

If it's the first time to run, there are some [initial works](https://github.com/neal1991/gshark/blob/0ea3365f88e012df3fef1079df04a4f4b266319d/models/models.go#L31) will be finished automatically.

* [Init Rules](https://github.com/neal1991/gshark/blob/0ea3365f88e012df3fef1079df04a4f4b266319d/models/models.go#L98)
* [Init admin](https://github.com/neal1991/gshark/blob/0ea3365f88e012df3fef1079df04a4f4b266319d/models/models.go#L117)
* [Init database](https://github.com/neal1991/gshark/blob/0ea3365f88e012df3fef1079df04a4f4b266319d/models/models.go#L47)

## Reference

* [x-patrol](https://github.com/MiSecurity/x-patrol)
* [authz](https://github.com/go-macaron/authz)


