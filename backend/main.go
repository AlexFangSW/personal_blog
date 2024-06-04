package main

import (
	"blog/api"
	"blog/config"
	"blog/db/models/sqlite"
	"blog/util"
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func run() error {
	log.SetFlags(log.Ltime | log.Llongfile | log.LstdFlags)

	// flags
	configPath := flag.String("config", "./config.json", "Config filepath")
	flag.Parse()
	slog.Info("load config", "path:", *configPath)

	// load config file
	rawConfig, err := os.ReadFile(*configPath)
	if err != nil {
		return fmt.Errorf("run: load config failed: %w", err)
	}
	config := config.NewConfig()
	json.Unmarshal(rawConfig, config)

	// init logger
	util.InitLogger(config.Logger.Level)

	// db connection
	db, err := sql.Open("sqlite3", config.DB.DSNURL)
	if err != nil {
		return fmt.Errorf("run: open db connection failed: %w", err)
	}

	// TODO: refector models to <db>model.<blogs | tags | topics>
	// let server use interfaces
	model := sqlite.New(db, config.DB)
	ctx := context.Background()
	if err := model.Prepare(ctx, config.DB.Timeout); err != nil {
		return fmt.Errorf("run: model prepare failed: %w", err)
	}

	// setup server
	server := api.NewServer(*config, *model)
	if err := server.Start(); err != nil {
		return fmt.Errorf("run: server start failed: %w", err)
	}

	return nil
}

// TODO: graceful shutdown
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
