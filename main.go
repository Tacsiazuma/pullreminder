package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:  "pullreminder",
		Usage: "Shows pull requests which require your attention",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "owner",
				Required: true,
				Usage:    "owner of the repository",
			},
			&cli.StringFlag{
				Name:     "name",
				Required: true,
				Usage:    "name of the repository",
			},
			&cli.StringFlag{
				Name:     "username",
				Required: true,
				Usage:    "your github username",
			},
		},
		Action: func(ctx *cli.Context) error {
			service := New(NewGithubProvider(os.Getenv("GITHUB_TOKEN")))
			err := service.AddRepository(&Repository{Name: ctx.String("name"), Owner: ctx.String("owner"), Provider: "github"})
			if err != nil {
				return err
			}
			err = service.AddCredentials("github", os.Getenv("GITHUB_TOKEN"))
			if err != nil {
				return err
			}
			prs, err := service.NeedsAttention(context.Background())
			if err != nil {
				return err
			}
			fmt.Printf("You have %d PRs to check\n\n", len(prs))
			for _, pr := range prs {
				fmt.Printf("#%d %s\n", pr.Number, pr.URL)
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
