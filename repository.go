package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"github.com/olekukonko/tablewriter"
	"github.com/ryanuber/columnize"
	"github.com/urfave/cli"
)

var queryFlagsRepository = []string{
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

var commandRepository = cli.Command{
	Name:      "repository",
	UsageText: "ghs repository [command options] [query]",
	Aliases:   []string{"repo", "r"},
	Usage:     "Search repositorys",
	Flags:     append(flagRepository, flagCommon...),
	Action:    doRepository,
}

var flagRepository = []cli.Flag{
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
}

func doRepository(c *cli.Context) error {
	// Validate args and flags
	if len(c.Args()) < 1 {
		return cli.NewExitError("is not input query", 1)
	}

	// Building search query
	query := BuildQuery(c, queryFlagsRepository)

	// Building search options
	opts := BuildSearchOptions(c)

	// Do search
	client := github.NewClient(nil)
	ctx := context.Background()
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
