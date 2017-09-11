package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"github.com/ryanuber/columnize"
	"github.com/urfave/cli"
)

var queryFlagsCode = []string{
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

func doCode(c *cli.Context) error {
	// Validate args and flags
	if len(c.Args()) < 1 {
		return cli.NewExitError("is not input query", 1)
	}

	ctx := context.Background()
	client := github.NewClient(nil)

	// Building search query
	query := BuildQuery(c, queryFlagsCode)

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
