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
	)

	// use custom client to set timeout
	commonFlags := []cli.Flag{
		&cli.StringFlag{
			Name:        "url",
			Value:       "https://localhost:8080/api/v1",
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
				Usage:                  "Sync everything",
				UseShortOptionHandling: true,
				Action: func(cCtx *cli.Context) error {
					return syncAll(ctxCancel, url, sourcePath)
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
