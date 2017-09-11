package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	// "github.com/ryanuber/columnize"
	"github.com/urfave/cli"
)

var queryFlagsUser = []string{
	"type",
	"in",
	"repos",
	"location",
	"language",
	"created",
	"followers",
}

var commandUser = cli.Command{
	Name:    "user",
	Aliases: []string{"u"},
	Usage:   "Search users",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "type",
			Usage: "With this qualifier you can restrict the search to just personal accounts (user) or just organization accounts (org).",
		},
		cli.StringFlag{
			Name:  "in",
			Usage: "Qualifies which fields are searched. With this qualifier you can restrict the search to just the username (login), public email (email), full name (fullname), or any combination of these.",
		},
		cli.StringFlag{
			Name:  "repos",
			Usage: "Filters users based on the number of repositories they have.",
		},
		cli.StringFlag{
			Name:  "location",
			Usage: "Filter users by the location indicated in their profile.",
		},
		cli.StringFlag{
			Name:  "language",
			Usage: "Search for users that have repositories that match a certain language.",
		},
		cli.StringFlag{
			Name:  "created",
			Usage: "Filter users based on when they joined.",
		},
		cli.StringFlag{
			Name:  "followers",
			Usage: "Filter users based on the number of followers they have.",
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
	},
	Action: doUser,
}

func doUser(c *cli.Context) error {
	// Validate args and flags
	if len(c.Args()) < 1 {
		return cli.NewExitError("is not input query", 1)
	}

	ctx := context.Background()
	client := github.NewClient(nil)

	// Building search query
	query := BuildQuery(c, queryFlagsUser)

	// Building search options
	opts := BuildSearchOptions(c)

	// Do search
	client := github.NewClient(nil)
	ctx := context.Background()
	result, _, err := client.Search.Users(ctx, query, opts)

	if err == nil {
		var datas []string
		for _, user := range result.Users {
			data := user.GetLogin()
			datas = append(datas, data)
			fmt.Println(data)
		}
		// result := columnize.SimpleFormat(datas)
		// fmt.Println(result)
	}

	return err
}
