package main

import (
	"os"

	"github.com/codegangsta/cli"
)

var commonFlags = []cli.Flag{
	{
		Name:      "app",
		ShortName: "a",
		Usage:     "The name of the app",
	},
}

var Commands = []cli.Command{
	{
		Name:   "elb-status",
		Usage:  "Check the ELB status for an app",
		Flags:  commonFlags,
		Action: runServer,
	},
}

func main() {
	app := cli.NewApp()
	app.Name = "empctl"
	app.Usage = "Usage"
	app.Commands = Commands

	app.Run(os.Args)
}
