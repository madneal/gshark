package cmd

import (
	"github.com/madneal/gshark/search"
	"github.com/madneal/gshark/web"
	"github.com/urfave/cli"
)

var Web = cli.Command{
	Name:        "web",
	Usage:       "Startup a web Service",
	Description: "Startup a web Service",
	Action:      web.RunWeb,
	Flags: []cli.Flag{
		boolFlag("debug, d", "Debug Mode"),
		stringFlag("host, H", "0.0.0.0", "web listen address"),
		intFlag("port, p", 8000, "web listen port"),
	},
}

var Scan = cli.Command{
	Name:        "scan",
	Usage:       "Start to scan leak info",
	Description: "start to scan leak info",
	Action:      search.Scan,
	Flags: []cli.Flag{
		stringFlag("mode, m", "github", "scan mode: github, searchcode, gitlab, all"),
		intFlag("time, t", 900, "scan interval(second)"),
	},
}

func stringFlag(name, value, usage string) cli.StringFlag {
	return cli.StringFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}

func boolFlag(name, usage string) cli.BoolFlag {
	return cli.BoolFlag{
		Name:  name,
		Usage: usage,
	}
}

func intFlag(name string, value int, usage string) cli.IntFlag {
	return cli.IntFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}
