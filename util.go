package main

import (
	"strings"

	"github.com/google/go-github/github"
	"github.com/urfave/cli"
)

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
