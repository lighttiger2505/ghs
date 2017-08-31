package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"github.com/olekukonko/tablewriter"
	"github.com/ryanuber/columnize"
	"github.com/urfave/cli"
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
	Name:      "repository",
	UsageText: "ghs repository [command options] [query]",
	Aliases:   []string{"repo", "r"},
	Usage:     "Search repositorys",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "created",
			Usage: "Filters repositories based on date of creation.",
		},
		cli.StringFlag{
			Name:  "pushed",
			Usage: "Filters repositories based on when they were last updated.",
		},
		cli.StringFlag{
			Name:  "forks",
			Usage: "Filters repositories based on the number of forks.",
		},
		cli.StringFlag{
			Name:  "in",
			Usage: "Qualifies which fields are searched. With this qualifier you can restrict the search to just the repository name, description, readme, or any combination of these.",
		},
		cli.StringFlag{
			Name:  "language",
			Usage: "Searches repositories based on the language they're written in.",
		},
		cli.StringFlag{
			Name:  "repo",
			Usage: "Limits searches to a specific repository.",
		},
		cli.StringFlag{
			Name:  "user",
			Usage: "Limits searches to a specific user.",
		},
		cli.StringFlag{
			Name:  "size",
			Usage: "Finds repositories that match a certain size (in kilobytes).",
		},
		cli.StringFlag{
			Name:  "stars",
			Usage: "Searches repositories based on the number of stars.",
		},
		cli.StringFlag{
			Name:  "topics",
			Usage: "Filters repositories based on the specified topic.",
		},
		cli.StringFlag{
			Name:  "sort, s",
			Usage: "The sort field. One of stars, forks, or updated. Default: results are sorted by best match.",
		},
		cli.StringFlag{
			Name:  "order, o",
			Value: "desc",
			Usage: "The sort order if sort parameter is provided. One of asc or desc.",
		},
		cli.IntFlag{
			Name:  "num, n",
			Value: 10,
			Usage: "Display num of search results.",
		},
		cli.IntFlag{
			Name:  "page, p",
			Value: 1,
			Usage: "Display page of search results.",
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
	Name:      "code",
	UsageText: "ghs code [command options] [query]",
	Aliases:   []string{"c"},
	Usage:     "Search codes",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "in",
			Usage: "Qualifies which fields are searched. With this qualifier you can restrict the search to the file contents (file), the file path (path), or both.",
		},
		cli.StringFlag{
			Name:  "language",
			Usage: "Searches code based on the language it's written in.",
		},
		cli.StringFlag{
			Name:  "fork",
			Usage: "Specifies that code from forked repositories should be searched (true). Repository forks will not be searchable unless the fork has more stars than the parent repository.",
		},
		cli.StringFlag{
			Name:  "size",
			Usage: "Finds files that match a certain size (in bytes).",
		},
		cli.StringFlag{
			Name:  "path",
			Usage: "Specifies the path prefix that the resulting file must be under.",
		},
		cli.StringFlag{
			Name:  "filename",
			Usage: "Matches files by a substring of the filename.",
		},
		cli.StringFlag{
			Name:  "extension",
			Usage: "Matches files with a certain extension after a dot.",
		},
		cli.StringFlag{
			Name:  "user",
			Usage: "Limits searches to a specific user.",
		},
		cli.StringFlag{
			Name:  "repo",
			Usage: "Limits searches to a specific repository.",
		},
		cli.StringFlag{
			Name:  "sort, s",
			Usage: "The sort field. One of stars, forks, or updated. Default: results are sorted by best match.",
		},
		cli.StringFlag{
			Name:  "order, o",
			Value: "desc",
			Usage: "The sort order if sort parameter is provided. One of asc or desc.",
		},
		cli.IntFlag{
			Name:  "num, n",
			Value: 10,
			Usage: "Display num of search results.",
		},
		cli.IntFlag{
			Name:  "page, p",
			Value: 1,
			Usage: "Display page of search results.",
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
	Action: doCode,
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
	query := BuildRepositoryQuery(c)

	// Setting search option
	opts := &github.SearchOptions{
		Sort:        c.String("sort"),
		Order:       c.String("order"),
		TextMatch:   false,
		ListOptions: github.ListOptions{PerPage: c.Int("num"), Page: c.Int("page")},
	}

	// Do repository search
	result, _, err := client.Search.Repositories(ctx, query, opts)

	if err == nil {
		// Draw result content
		if c.Bool("only") {
			DrawOnly(result.Repositories)
		} else if c.Bool("oneline") {
			DrawOneline(result.Repositories)
		} else {
			DrawTable(result.Repositories)
		}
	}

	return err
}

func BuildRepositoryQuery(c *cli.Context) string {
	var query []string
	queryFlags := []string{
		"created",
		"pushed",
		"forks",
		"in",
		"language",
		"repo",
		"user",
		"size",
		"stars",
		"topics",
	}

	for _, flagName := range c.FlagNames() {
		isQuery := false
		for _, queryFlag := range queryFlags {
			if queryFlag == flagName {
				isQuery = true
				break
			}
		}
		if !isQuery {
			continue
		}

		if c.String(flagName) != "" {
			query = append(query, flagName+":"+c.String(flagName))
		}
	}
	query = append(query, c.Args()[0])
	return strings.Join(query, " ")
}

func DrawOnly(repos []github.Repository) {
	for _, repo := range repos {
		fmt.Println(repo.GetFullName())
	}
}

func DrawOneline(repos []github.Repository) {

	var datas []string
	for _, repo := range repos {
		data := repo.GetFullName() + "|" + repo.GetDescription()
		datas = append(datas, data)
	}
	result := columnize.SimpleFormat(datas)
	fmt.Println(result)
}

func DrawTable(repos []github.Repository) {
	var datas [][]string
	for _, repo := range repos {
		data := []string{
			repo.GetFullName(),
			fmt.Sprint(repo.GetStargazersCount()),
			repo.GetLanguage(),
			repo.GetDescription(),
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
	// Validate args and flags
	if len(c.Args()) < 1 {
		return cli.NewExitError("is not input query", 1)
	}

	ctx := context.Background()
	client := github.NewClient(nil)

	// Building search query
	query := BuildCodeQuery(c)

	// Setting search option
	opts := &github.SearchOptions{
		Sort:        c.String("sort"),
		Order:       c.String("order"),
		TextMatch:   false,
		ListOptions: github.ListOptions{PerPage: c.Int("num"), Page: c.Int("page")},
	}

	// Do repository search
	result, _, err := client.Search.Code(ctx, query, opts)

	if err == nil {
		var datas []string
		for _, code := range result.CodeResults {
			data := code.GetName() + "|" + code.GetPath() + "|" + code.GetHTMLURL() + "|" + code.GetSHA()
			datas = append(datas, data)
		}
		result := columnize.SimpleFormat(datas)
		fmt.Println(result)
	}

	return err
}

func BuildCodeQuery(c *cli.Context) string {
	var query []string
	queryFlags := []string{
		"in",
		"language",
		"fork",
		"size",
		"path",
		"filename",
		"extension",
		"user",
		"repo",
	}

	for _, flagName := range c.FlagNames() {
		isQuery := false
		for _, queryFlag := range queryFlags {
			if queryFlag == flagName {
				isQuery = true
				break
			}
		}
		if !isQuery {
			continue
		}

		if c.String(flagName) != "" {
			query = append(query, flagName+":"+c.String(flagName))
		}
	}
	query = append(query, c.Args()[0])
	return strings.Join(query, " ")
}
