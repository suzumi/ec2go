package main

import (
	"os"

	"github.com/suzumi/ec2go"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "ec2go"
	app.Usage = "show ec2 instance list tool"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "show instance list",
			Action:  ec2go.ListAction,
			Subcommands: []cli.Command{
				{
					Name:   "all",
					Usage:  "show all state",
					Action: ec2go.ListAllSubAction,
				},
			},
		},
		{
			Name:    "ssh",
			Aliases: []string{"s"},
			Usage:   "filter instance",
			Action:  ec2go.SshAction,
		},
	}

	app.Run(os.Args)
}
