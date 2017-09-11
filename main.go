package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	err := newApp().Run(os.Args)
	var exitCode = 0
	if err != nil {
		fmt.Println(os.Stderr, err.Error())
		exitCode = 1
	}
	os.Exit(exitCode)
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = "ghs"
	app.HelpName = "ghs"
	app.Usage = "Call GitHub search API."
	app.UsageText = "ghs [global options] command [command options] [query]"
	app.Version = "0.0.1"
	app.Author = "lighttiger2505"
	app.Email = "lighttiger2505@gmail.com"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "~/.ghsrc",
			Usage: "Load configuration from `FILE`",
		},
	}
	app.Commands = []cli.Command{
		commandRepository,
		commandCommit,
		commandCode,
		commandIssue,
		commandUser,
	}
	return app
}

var commandIssue = cli.Command{
	Name:    "issue",
	Aliases: []string{"i"},
	Usage:   "Search issues",
	Action:  doIssue,
}

func doIssue(c *cli.Context) error {
	return nil
}
