package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

func main() {
	var url string
	var metaFile string
	var blogsDir string

	// use custom client to set timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	commonFlags := []cli.Flag{
		&cli.StringFlag{
			Name:        "url",
			Value:       "https://localhost:8080/api/v1",
			Usage:       "Base `URL` for api",
			Destination: &url,
		},
		&cli.StringFlag{
			Name:        "meta",
			Value:       "./meta.yaml",
			Usage:       "`PATH` to yaml file containing topic and tags",
			Destination: &metaFile,
		},
		&cli.StringFlag{
			Name:        "blogs",
			Value:       "./blogs",
			Usage:       "`PATH` to directory containing blogs",
			Destination: &blogsDir,
		},
	}

	app := &cli.App{
		Name:  "Coding notes sync tool",
		Usage: "Used for syncing files to blog page",
		Commands: []*cli.Command{
			{
				Name:  "sync",
				Usage: "Sync everything",
				Action: func(cCtx *cli.Context) error {
					// Sync everything
					return syncAll(url, metaFile, blogsDir, client)
				},
				Flags: commonFlags,
				Subcommands: []*cli.Command{
					{
						Name:  "meta",
						Usage: "Sync topic and tags",
						Action: func(cCtx *cli.Context) error {
							// Sync topic and tags
							// update create delete
							fmt.Println("Sync topic and tags")
							return nil
						},
						Flags: commonFlags,
					},
					{
						Name:  "blogs",
						Usage: "Sync blogs",
						Action: func(cCtx *cli.Context) error {
							// Sync blogs
							// update create delete
							fmt.Println("Sync blogs")
							return nil
						},
						Flags: commonFlags,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
