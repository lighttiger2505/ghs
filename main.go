package main

import (
	"context"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/google/go-github/github"
	"github.com/olekukonko/tablewriter"
	// "github.com/ryanuber/columnize"
	"github.com/urfave/cli"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	newApp().Run(os.Args)
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

var commandRepository = cli.Command{
	Name:    "repository",
	Aliases: []string{"repo", "r"},
	Usage:   "Search repositorys",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "in",
			Usage: "The in qualifier limits what fields are searched. With this qualifier you can restrict the search to just the repository name, description, README, or any combination of these. Without the qualifier, only the name and description are searched.",
		},
		cli.StringFlag{
			Name:  "size",
			Usage: "The size qualifier finds repositories that match a certain size (in kilobytes), using greater than, less than, and range qualifiers.",
		},
		cli.StringFlag{
			Name:  "mirror",
			Usage: "You can search repositories based on whether or not they're a mirror and are hosted elsewhere.",
		},
		cli.StringFlag{
			Name:  "forks",
			Usage: " The forks qualifier specifies the number of forks a repository should have, using greater than, less than, and range qualifiers.",
		},
		cli.StringFlag{
			Name:  "created",
			Usage: "You can filter repositories based on time of creation or time of last update. For repository creation, you can use the created qualifier; to find out when a repository was last updated, you'll want to use the pushed qualifier. The pushed qualifier will return a list of repositories, sorted by the most recent commit made on any branch in the repository.",
		},
		cli.StringFlag{
			Name:  "pushed",
			Usage: "You can filter repositories based on time of creation or time of last update. For repository creation, you can use the created qualifier; to find out when a repository was last updated, you'll want to use the pushed qualifier. The pushed qualifier will return a list of repositories, sorted by the most recent commit made on any branch in the repository.",
		},
		cli.StringFlag{
			Name:  "user",
			Usage: "To grab a list of a user's repositories.",
		},
		cli.StringFlag{
			Name:  "org",
			Usage: "To grab a list of a organization's repositories.",
		},
		cli.StringFlag{
			Name:  "topic",
			Usage: "You can find all of the repositories that are classified with a particular topic.",
		},
		cli.StringFlag{
			Name:  "topics",
			Usage: "You can find repositories by the number of applied topics, using the topics qualifier along with greater than, less than, and range qualifiers.",
		},
		cli.StringFlag{
			Name:  "language",
			Usage: "You can also search repositories based on what language they're written in.",
		},
		cli.StringFlag{
			Name:  "stars",
			Usage: "You can search repositories based on the number of stars a repository has, using greater than, less than, and range qualifiers",
		},
		cli.StringFlag{
			Name:  "sort, s",
			Usage: "The sort field. One of stars, forks, or updated. Default: results are sorted by best match.",
		},
		cli.StringFlag{
			Name:  "order, o",
			Usage: "The sort order if sort parameter is provided. One of asc or desc. Default: desc",
		},
		cli.BoolFlag{
			Name:  "only",
			Usage: "Draw repository full name only",
		},
		cli.BoolFlag{
			Name:  "oneline",
			Usage: "Draw repository online",
		},
		cli.BoolFlag{
			Name:  "table",
			Usage: "Draw table",
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

	// Validate args and flags
	if len(c.Args()) < 1 {
		return cli.NewExitError("is not input query", 1)
	}

	ctx := context.Background()
	client := github.NewClient(nil)

	// Building search query
	query := BuildQuery(c)

	// Setting search option
	opts := &github.SearchOptions{
		Sort:        c.String("sort"),
		Order:       c.String("order"),
		TextMatch:   false,
		ListOptions: github.ListOptions{PerPage: 100, Page: 1},
	}

	// Draw spinner
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Start()

	// Do repository search
	result, resp, err := client.Search.Repositories(ctx, query, opts)

	fmt.Println(resp.Request.URL)

	// Clear spinner
	s.Stop()

	if err == nil {
		// Draw result content
		DrawTable(result.Repositories)
	}

	return err
}

func BuildQuery(c *cli.Context) string {
	var query []string
	ignoreFlags := []string{
		"sort",
		"order",
		"only",
		"oneline",
		"table",
	}

	for _, flagName := range c.FlagNames() {
		ignore := false
		for _, ignoreFlag := range ignoreFlags {
			if ignoreFlag == flagName {
				ignore = true
			}
		}
		if ignore == true {
			continue
		}

		if c.String(flagName) != "" {
			query = append(query, flagName+":"+c.String(flagName))
		}
	}
	query = append(query, c.Args()[0])
	return strings.Join(query, " ")
}

func DrawOnly([]github.Repository) {
}

func DrawOneline([]github.Repository) {
}

func DrawTable(repos []github.Repository) {
	fmt.Println(len(repos))

	var datas [][]string
	for _, repo := range repos {
		var desc string
		if repo.Description != nil {
			desc = *repo.Description
		}
		var lang string
		if repo.Language != nil {
			lang = *repo.Language
		}

		data := []string{
			*repo.FullName,
			fmt.Sprint(*repo.StargazersCount),
			lang,
			desc,
		}
		datas = append(datas, data)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Star", "Language", "Description"})

	for _, data := range datas {
		table.Append(data)
	}
	table.Render()
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
