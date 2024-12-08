package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/kirsle/configdir"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
)

func main() {
	configpath := configdir.LocalConfig("pullreminder")
	err := configdir.MakePath(configpath)
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("sqlite3", filepath.Join(configpath, "db.sqlite"))
	if err != nil {
		log.Fatal(err)
	}
	service := New(NewGithubProvider(os.Getenv("GITHUB_TOKEN")), NewSqliteStore(db))
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "add-repo",
				Usage: "Adds a repository",
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
						Name:     "provider",
						Required: true,
						Usage:    "provider of the repository",
					},
				},
				Action: func(ctx *cli.Context) error {
					err = service.AddRepository(&Repository{Name: ctx.String("name"), Owner: ctx.String("owner"), Provider: ctx.String("provider")})
					if err != nil {
						return err
					}
					return nil
				},
			}, {
				Name:  "add-creds",
				Usage: "Adds credentials",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "provider",
						Required: true,
						Usage:    "provider of the repository",
					},
				},
				Action: func(ctx *cli.Context) error {
					err = service.AddCredentials(ctx.String("provider"), os.Getenv("GITHUB_TOKEN"))
					if err != nil {
						return err
					}
					return nil
				},
			}, {
				Name:  "check",
				Usage: "Check for PRs",
				Action: func(ctx *cli.Context) error {
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
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
