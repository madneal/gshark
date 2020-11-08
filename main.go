package main

import (
	"github.com/madneal/gshark/cmd"
	"github.com/madneal/gshark/logger"
	"github.com/urfave/cli"
	"os"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	app := cli.NewApp()
	app.Name = "GShark"
	app.Author = "madneal"
	app.Email = "bing.ecnu@gmail.com"
	app.Version = "20201109"
	app.Usage = "Scan for sensitive information easily and effectively."
	app.Commands = []cli.Command{cmd.Web, cmd.Scan}
	app.Flags = append(app.Flags, cmd.Web.Flags...)
	app.Flags = append(app.Flags, cmd.Scan.Flags...)
	err := app.Run(os.Args)
	if err != nil {
		logger.Log.Error(err)
	}
}
