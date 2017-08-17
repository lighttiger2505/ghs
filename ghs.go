package main

import (
	"github.com/urfave/cli"
	"os"
)

func main() {
	var language string
	app := cli.NewApp()

	app.Name = "ghs"
	app.Usage = "Call GitHub search API."

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "~/.ghsrc",
			Usage: "Load configuration from `FILE`",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "repository",
			Aliases: []string{"repo, r"},
			Usage:   "Search repositorys",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:    "commit",
			Aliases: []string{"m"},
			Usage:   "Search commits",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:    "code",
			Aliases: []string{"c"},
			Usage:   "Search codes",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:    "issue",
			Aliases: []string{"i"},
			Usage:   "Search issues",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:    "user",
			Aliases: []string{"u"},
			Usage:   "Search users",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	}

	app.Run(os.Args)
}
