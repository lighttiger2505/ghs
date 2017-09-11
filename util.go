package main

import (
	"strings"

	"github.com/urfave/cli"
)

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
