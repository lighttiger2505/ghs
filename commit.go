package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/github"
	"github.com/ryanuber/columnize"
	"github.com/urfave/cli"
)

var commandCommit = cli.Command{
	Name:      "commit",
	UsageText: "ghs commit [command options] [query]",
	Aliases:   []string{"m"},
	Usage:     "Search commits",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "author",
			Usage: "Matches commits authored by a user (based on email settings).",
		},
		cli.StringFlag{
			Name:  "commiter",
			Usage: "Matches commits committed by a user (based on email settings).",
		},
		cli.StringFlag{
			Name:  "author-name",
			Usage: "Matches commits by author name.",
		},
		cli.StringFlag{
			Name:  "committer-name",
			Usage: "Matches commits by committer name.",
		},
		cli.StringFlag{
			Name:  "author-email",
			Usage: "Matches commits by author email.",
		},
		cli.StringFlag{
			Name:  "committer-email",
			Usage: "Matches commits by committer email.",
		},
		cli.StringFlag{
			Name:  "author-date",
			Usage: "Matches commits by author date range.",
		},
		cli.StringFlag{
			Name:  "ommitter-date",
			Usage: "Matches commits by committer date range.",
		},
		cli.StringFlag{
			Name:  "merge",
			Usage: "true filters to merge commits, false filters out merge commits.",
		},
		cli.StringFlag{
			Name:  "hash",
			Usage: "Matches commits by hash.",
		},
		cli.StringFlag{
			Name:  "parent",
			Usage: "Matches commits that have a particular parent.",
		},
		cli.StringFlag{
			Name:  "tree",
			Usage: "Matches commits by tree hash.",
		},
		cli.StringFlag{
			Name:  "is",
			Usage: "Matches public or private repositories.",
		},
		cli.StringFlag{
			Name:  "user",
			Usage: "Limits searches to a specific user.",
		},
		cli.StringFlag{
			Name:  "org",
			Usage: "Limits searches to a specific organization.",
		},
		cli.StringFlag{
			Name:  "repo",
			Usage: "Limits searches to a specific repository.",
		},
		cli.StringFlag{
			Name:  "sort, s",
			Usage: "The sort field. Can be author-date or committer-date. Default: results are sorted by best match.",
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
	},
	Action: doCommit,
}

func doCommit(c *cli.Context) error {
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
	result, _, err := client.Search.Commits(ctx, query, opts)

	if err == nil {
		// Draw result content
		fmt.Println("検索できたよー")
		var datas []string
		for _, commit := range result.Commits {
			data := commit.Commit.GetSHA() + "|" + commit.Commit.GetMessage()
			datas = append(datas, data)
		}
		result := columnize.SimpleFormat(datas)
		fmt.Println(result)
	}

	return err
}

func BuildCommitQuery(c *cli.Context) string {
	var query []string
	queryFlags := []string{
		"author",
		"commiter",
		"author-name",
		"committer-name",
		"author-email",
		"committer-email",
		"author-date",
		"ommitter-date",
		"merge",
		"hash",
		"parent",
		"tree",
		"is",
		"user",
		"org",
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
