package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"github.com/ryanuber/columnize"
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
	Flags:   append(flagUser, flagCommon...),
	Action:  doUser,
}

var flagUser = []cli.Flag{
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
}

func doUser(c *cli.Context) error {
	// Validate args and flags
	if len(c.Args()) < 1 {
		return cli.NewExitError("is not input query", 1)
	}

	// Building search query
	query := BuildQuery(c, queryFlagsUser)

	// Building search options
	opts := BuildSearchOptions(c)

	// Do search
	client := github.NewClient(nil)
	ctx := context.Background()
	result, _, err := client.Search.Users(ctx, query, opts)

	// Draw result
	if err == nil {
		if c.Bool("oneline") {
			DrawUserOneline(result.Users)
		} else {
			DrawUserDefault(result.Users)
		}
	}

	return err
}

func DrawUserOneline(users []github.User) {
	var datas []string
	for _, user := range users {
		data := user.GetLogin()
		datas = append(datas, data)
	}
	result := columnize.SimpleFormat(datas)
	fmt.Println(result)
}

func DrawUserDefault(users []github.User) {
	client := github.NewClient(nil)
	ctx := context.Background()

	for _, user := range users {
		DrawMainContent("user", user.GetLogin())

		userDetail, _, _ := client.Users.Get(ctx, user.GetLogin())
		DrawSubContentOneline("Name", userDetail.GetName())
		DrawSubContentOneline("Location", userDetail.GetLocation())
		DrawSubContentOneline("Email", userDetail.GetEmail())
		DrawSubContentMultiLine(userDetail.GetBio())
		fmt.Fprintln(os.Stdout)
	}
}
