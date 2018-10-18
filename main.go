package main

import (
	"github.com/urfave/cli"
	"os"
	"runtime"
	"x-patrol/cmd"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
func main() {
	app := cli.NewApp()
	app.Name = "Github leaked patrol"
	app.Author = "netxfly"
	app.Email = "x@xsec.io"
	app.Version = "20180131"
	app.Usage = "Github leaked patrol, support search github and local repos"
	app.Commands = []cli.Command{cmd.Web, cmd.Scan}
	app.Flags = append(app.Flags, cmd.Web.Flags...)
	app.Flags = append(app.Flags, cmd.Scan.Flags...)
	app.Run(os.Args)
}
