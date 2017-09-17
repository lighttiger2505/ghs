package main

import (
	"strings"

	"github.com/google/go-github/github"
	"github.com/urfave/cli"
)

var flagCommon = []cli.Flag{
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
		Name:  "oneline",
		Usage: "Draw content online",
	},
}

// BuildQuery is format value of flag from cli into github search API query.
func BuildQuery(c *cli.Context, queryFlags []string) string {
	var buildFlags []string
	for _, flagName := range c.FlagNames() {
		for _, queryFlag := range queryFlags {
			if queryFlag == flagName {
				buildFlags = append(buildFlags, queryFlag)
			}
		}
	}

	var query []string
	for _, flagName := range buildFlags {
		flagValue := c.String(flagName)
		if flagValue != "" {
			query = append(query, flagName+":"+flagValue)
		}
	}
	query = append(query, c.Args()[0])
	return strings.Join(query, " ")
}

// BuildSearchOptions is format value of flag from cli into github search API options.
func BuildSearchOptions(c *cli.Context) *github.SearchOptions {
	opts := &github.SearchOptions{
		Sort:        c.String("sort"),
		Order:       c.String("order"),
		TextMatch:   false,
		ListOptions: github.ListOptions{PerPage: c.Int("num"), Page: c.Int("page")},
	}
	return opts
}
