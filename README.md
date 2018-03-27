# x-patrol

This project is based on [x-patrol](https://github.com/MiSecurity/x-patrol). It is utilize to scan sensitive information in github or local repos.

## Run

You should build the `main.go` file firstly with the command `go build main.go`.
```
USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   20180131

AUTHOR:
   netxfly <x@xsec.io>

COMMANDS:
     web      Startup a web Service
     scan     start to scan github leak info
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug, -d             Debug Mode
   --host value, -H value  web listen address (default: "0.0.0.0")
   --port value, -p value  web listen port (default: 8000)
   --mode value, -m value  scan mode: github, local, all (default: "github")
   --time value, -t value  scan interval(second) (default: 900)
   --help, -h              show help
   --version, -v           print the version
```

### Start the service

`main web`

*Note: If you wish your service can be accessed from outer network, the service should bind to 0.0.0.0 instead of 127.0.0.1. The host can be configed in app.ini file.*
