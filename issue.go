package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/github"
	"github.com/ryanuber/columnize"
	"github.com/urfave/cli"
)

var commandIssue = cli.Command{
	Name:      "issue",
	UsageText: "ghs issue [command options] [query]",
	Aliases:   []string{"i"},
	Usage:     "Search issues",
	Action:    doIssue,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "type",
			Usage: "With this qualifier you can restrict the search to issues (issue) or pull request (pr) only.",
		},
		cli.StringFlag{
			Name:  "in",
			Usage: "Qualifies which fields are searched. With this qualifier you can restrict the search to just the title (title), body (body), comments (comments), or any combination of these.",
		},
		cli.StringFlag{
			Name:  "author",
			Usage: "Finds issues or pull requests created by a certain user.",
		},
		cli.StringFlag{
			Name:  "assignee",
			Usage: "Finds issues or pull requests that are assigned to a certain user.",
		},
		cli.StringFlag{
			Name:  "mentions",
			Usage: "Finds issues or pull requests that mention a certain user.",
		},
		cli.StringFlag{
			Name:  "commenter",
			Usage: "Finds issues or pull requests that a certain user commented on.",
		},
		cli.StringFlag{
			Name:  "involves",
			Usage: "Finds issues or pull requests that were either created by a certain user, assigned to that user, mention that user, or were commented on by that user.",
		},
		cli.StringFlag{
			Name:  "team",
			Usage: "For organizations you're a member of, finds issues or pull requests that @mention a team within the organization.",
		},
		cli.StringFlag{
			Name:  "state",
			Usage: "Filter issues or pull requests based on whether they're open or closed.",
		},
		cli.StringFlag{
			Name:  "labels",
			Usage: "Filters issues or pull requests based on their labels.",
		},
		cli.StringFlag{
			Name:  "no",
			Usage: "Filters items missing certain metadata, such as label, milestone, or assignee",
		},
		cli.StringFlag{
			Name:  "language",
			Usage: "Searches for issues or pull requests within repositories that match a certain language.",
		},
		cli.StringFlag{
			Name:  "is",
			Usage: "Searches for items within repositories that match a certain state, such as open, closed, or merged.",
		},
		cli.StringFlag{
			Name:  "created",
			Usage: "Filters issues or pull requests based on date of creation.",
		},
		cli.StringFlag{
			Name:  "updated",
			Usage: "Filters issues or pull requests based on when they were last updated.",
		},
		cli.StringFlag{
			Name:  "merged",
			Usage: "Filters pull requests based on the date when they were merged.",
		},
		cli.StringFlag{
			Name:  "status",
			Usage: "Filters pull requests based on the commit status.",
		},
		cli.StringFlag{
			Name:  "head",
			Usage: "Filters pull requests based on the branch that they came from.",
		},
		cli.StringFlag{
			Name:  "base",
			Usage: "Filters pull requests based on the branch that they are modifying.",
		},
		cli.StringFlag{
			Name:  "close",
			Usage: "Filters issues or pull requests based on the date when they were closed.",
		},
		cli.StringFlag{
			Name:  "comments",
			Usage: "Filters issues or pull requests based on the quantity of comments.",
		},
		cli.StringFlag{
			Name:  "user",
			Usage: "user Limits searches to a specific user or repository.",
		},
		cli.StringFlag{
			Name:  "repo",
			Usage: "repo Limits searches to a specific user or repository.",
		},
		cli.StringFlag{
			Name:  "project",
			Usage: "Limits searches to a specific project board in a repository or organization.",
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
}

func doIssue(c *cli.Context) error {
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

func BuildIssueQuery(c *cli.Context) string {
	var query []string
	queryFlags := []string{
		"type",
		"in",
		"author",
		"assignee",
		"mentions",
		"commenter",
		"involves",
		"team",
		"state",
		"no",
		"language",
		"is",
		"created",
		"updated",
		"merged",
		"status",
		"head",
		"base",
		"close",
		"comments",
		"user",
		"repo",
		"project",
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
