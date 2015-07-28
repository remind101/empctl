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
}

var Commands = []cli.Command{}

func AddCommand(c cli.Command) {
	Commands = append(Commands, c)
}

func main() {
	app := cli.NewApp()
	app.Name = "empctl"
	app.Usage = "Usage"
	app.Commands = Commands

	app.Run(os.Args)
}
