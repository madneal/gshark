package main

import (
	"github.com/urfave/cli"
	"gshark/cmd"
	"os"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	app := cli.NewApp()
	app.Name = "GShark"
	app.Author = "neal"
	app.Email = "bing.ecnu@gmail.com"
	app.Version = "20180131"
	app.Usage = "Scan for sensitive information in Github easily and effectively."
	app.Commands = []cli.Command{cmd.Web, cmd.Scan}
	app.Flags = append(app.Flags, cmd.Web.Flags...)
	app.Flags = append(app.Flags, cmd.Scan.Flags...)
	app.Run(os.Args)
}
