package main

import (
	"fmt"
	"github.com/madneal/gshark/cli"
	"os"
)

func main() {
	if err := cli.NewCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
