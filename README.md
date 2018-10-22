<p align="center">
   <img alt="GgShark logo" src="https://s1.ax1x.com/2018/10/17/idhZvj.png" />
   <h3 align="center">GShark</h3>
   <p align="center">Scan for sensitive information in Github easily and effectively.</p>
</p>

# GShark [![Go Report Card](https://goreportcard.com/badge/github.com/neal1991/gshark)](https://goreportcard.com/report/github.com/neal1991/gshark)

This project is based on [x-patrol](https://github.com/MiSecurity/x-patrol). It can be utilized to detect the sensitive information based on keywords with custome rules.

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
