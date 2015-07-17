package main

import (
	"os"

	"github.com/codegangsta/cli"
)

var commonFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "app, a",
		Usage: "The name of the app",
	},
	cli.StringFlag{
		Name:   "aws.id",
		EnvVar: "AWS_ACCESS_KEY_ID",
		Usage:  "AWS Access Key ID",
	},
	cli.StringFlag{
		Name:   "aws.key",
		EnvVar: "AWS_SECRET_ACCESS_KEY",
		Usage:  "AWS Secret Access Key",
	},
}

var Commands = []cli.Command{
	{
		Name:   "elb-status",
		Usage:  "Check the ELB status for an app",
		Flags:  commonFlags,
		Action: runELBStatus,
	},
}

func main() {
	app := cli.NewApp()
	app.Name = "empctl"
	app.Usage = "Usage"
	app.Commands = Commands

	app.Run(os.Args)
}

func runELBStatus(c *cli.Context) {

}
