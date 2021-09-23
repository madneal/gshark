package cmd

import (
	"github.com/madneal/gshark/core"
	"github.com/madneal/gshark/search"
	"github.com/urfave/cli/v2"
)

var Web = cli.Command{
	Name:        "web",
	Usage:       "Startup a web Service",
	Description: "Startup a web Service",
	Action:      core.RunWindowsServer,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "host, H",
			Value: "0.0.0.0",
			Usage: "web listen address",
		},
		&cli.IntFlag{
			Name:  "post, p",
			Value: 8000,
			Usage: "web listen port",
		},
	},
}

var Scan = cli.Command{
	Name:        "scan",
	Usage:       "Start to scan leak info",
	Description: "start to scan leak info",
	Action:      search.Scan,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "mode, m",
			Value: "github",
			Usage: "scan mode: github, searchcode, gitlab, all",
		},
		&cli.IntFlag{
			Name:  "time, t",
			Value: 900,
			Usage: "scan interval(second)",
		},
	},
}
