package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/urfave/cli"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	newApp().Run(os.Args)
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = "ghs"
	app.Usage = "Call GitHub search API."
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

var commandRepository = cli.Command{
	Name:    "repository",
	Aliases: []string{"repo, r"},
	Usage:   "Search repositorys",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "query",
			Usage: "The search keywords, as well as any qualifiers.",
		},
		cli.StringFlag{
			Name:  "sort",
			Usage: "The sort field. One of stars, forks, or updated. Default: results are sorted by best match.",
		},
		cli.StringFlag{
			Name:  "order",
			Usage: "The sort order if sort parameter is provided. One of asc or desc. Default: desc",
		},
	},
	Action: doRepository,
}

var commandCommit = cli.Command{
	Name:    "commit",
	Aliases: []string{"m"},
	Usage:   "Search commits",
	Action:  doCommit,
}

var commandCode = cli.Command{
	Name:    "code",
	Aliases: []string{"c"},
	Usage:   "Search codes",
	Action:  doCode,
}

var commandIssue = cli.Command{
	Name:    "issue",
	Aliases: []string{"i"},
	Usage:   "Search issues",
	Action:  doIssue,
}

var commandUser = cli.Command{
	Name:    "user",
	Aliases: []string{"u"},
	Usage:   "Search users",
	Action:  doUser,
}

func doRepository(c *cli.Context) error {
	ctx := context.Background()
	client := github.NewClient(nil)

	// Setting search option
	opts := &github.SearchOptions{
		Sort:        c.String("sort"),
		Order:       c.String("order"),
		TextMatch:   false,
		ListOptions: github.ListOptions{PerPage: 10, Page: 1},
	}

	// Do repository search
	result, _, err := client.Search.Repositories(ctx, c.String("query"), opts)

	// Draw result
	for _, repo := range result.Repositories {
		fmt.Println(*repo.FullName)
	}

	return err
}

func doCommit(c *cli.Context) error {
	return nil
}

func doIssue(c *cli.Context) error {
	return nil
}

func doUser(c *cli.Context) error {
	return nil
}

func doCode(c *cli.Context) error {
	return nil
}

func GetRequest(url string) string {
	res, _ := http.Get(url)
	defer res.Body.Close()
	byteArray, _ := ioutil.ReadAll(res.Body)
	return string(byteArray)
}
