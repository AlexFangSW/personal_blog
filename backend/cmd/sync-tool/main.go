package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"
)

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

func main() {
	var (
		url        string
		sourcePath string
		verbose    int
		batchSize  int
		username   string
		password   string
	)

	// use custom client to set timeout
	commonFlags := []cli.Flag{
		&cli.StringFlag{
			Name:        "url",
			Value:       "http://localhost:8080/api/v1",
			Usage:       "Base `URL` for api",
			Destination: &url,
		},
		&cli.StringFlag{
			Name:        "source",
			Value:       "./",
			Usage:       "`PATH` to folder containing meta.yaml, blogs/ and ids.json",
			Destination: &sourcePath,
		},
		&cli.BoolFlag{
			Name:  "v",
			Usage: "verbose, shows debug log",
			Count: &verbose,
		},
		&cli.IntFlag{
			Name:        "bs",
			Value:       5,
			Usage:       "max `SIZE` of concurrent requests",
			Destination: &batchSize,
		},
		&cli.StringFlag{
			Name:        "username",
			Value:       "",
			Usage:       "blog username (optional)",
			Destination: &username,
			EnvVars:     []string{"BLOG_USERNAME"},
		},
		&cli.StringFlag{
			Name:        "password",
			Value:       "",
			Usage:       "blog password (optional)",
			Destination: &password,
			EnvVars:     []string{"BLOG_PASSWORD"},
		},
	}

	ctxCancel, cancel := context.WithCancel(context.Background())
	defer cancel()
	notifyDone := make(chan bool, 1)

	app := &cli.App{
		Name:  "Coding notes sync tool",
		Usage: "Used for syncing files to blog page",
		Commands: []*cli.Command{
			{
				Name:                   "sync",
				Usage:                  `Sync everything. USERNAME and PASSWORD can be passed in as enviroment variables or input interactively.`,
				UseShortOptionHandling: true,
				Action: func(cCtx *cli.Context) error {
					return syncAll(
						ctxCancel,
						username,
						password,
						url,
						sourcePath,
						batchSize,
					)
				},
				Flags: commonFlags,
				Before: func(ctx *cli.Context) error {
					log.SetFlags(log.Llongfile | log.Ltime)
					if verbose > 0 {
						slog.SetLogLoggerLevel(slog.LevelDebug)
					}
					return nil
				},
			},
		},
	}

	go func() {
		if err := app.Run(os.Args); err != nil {
			log.Fatal(err)
		}
		notifyDone <- true
	}()

	notifyClose := make(chan os.Signal, 1)
	signal.Notify(notifyClose, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-notifyClose:
		cancel()
		slog.Warn("Stoping...")
		<-notifyDone
		slog.Warn("Stoped")
	case <-notifyDone:
	}
}
