package main

import (
	"fmt"
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
			Name:  "q",
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
	Action: func(c *cli.Context) error {
		SearchRepository(c.String("q"), c.String("sort"), c.String("order"))
		return nil
	},
}

var commandCommit = cli.Command{
	Name:    "commit",
	Aliases: []string{"m"},
	Usage:   "Search commits",
	Action: func(c *cli.Context) error {
		SearchCommit()
		return nil
	},
}

var commandCode = cli.Command{
	Name:    "code",
	Aliases: []string{"c"},
	Usage:   "Search codes",
	Action: func(c *cli.Context) error {
		SearchCode()
		return nil
	},
}

var commandIssue = cli.Command{
	Name:    "issue",
	Aliases: []string{"i"},
	Usage:   "Search issues",
	Action: func(c *cli.Context) error {
		SearchIssue()
		return nil
	},
}

var commandUser = cli.Command{
	Name:    "user",
	Aliases: []string{"u"},
	Usage:   "Search users",
	Action: func(c *cli.Context) error {
		SearchUser()
		return nil
	},
}

func SearchRepository(q string, sort string, order string) {
	fmt.Println(GetRequest("https://api.github.com/search/repositories?q=tetris+language:assembly&sort=stars&order=desc"))
}

func SearchCommit() {
	fmt.Println(GetRequest("https://api.github.com/search/repositories?q=tetris+language:assembly&sort=stars&order=desc"))
}

func SearchCode() {
	fmt.Println(GetRequest("https://api.github.com/search/repositories?q=tetris+language:assembly&sort=stars&order=desc"))
}

func SearchIssue() {
	fmt.Println(GetRequest("https://api.github.com/search/repositories?q=tetris+language:assembly&sort=stars&order=desc"))
}

func SearchUser() {
	fmt.Println(GetRequest("https://api.github.com/search/repositories?q=tetris+language:assembly&sort=stars&order=desc"))
}

func GetRequest(url string) string {
	res, _ := http.Get(url)
	defer res.Body.Close()
	byteArray, _ := ioutil.ReadAll(res.Body)
	return string(byteArray)
}
