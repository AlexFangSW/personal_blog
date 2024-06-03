package main

import (
	"blog/api"
	"blog/db/models"
	"blog/structs"
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
	config := structs.NewConfig()
	json.Unmarshal(rawConfig, config)

	// init logger
	util.InitLogger(config.Logger.Level)

	// db connection
	db, err := sql.Open("sqlite3", config.DB.DSNURL)
	if err != nil {
		return fmt.Errorf("run: open db connection failed: %w", err)
	}

	newModels := models.NewModels(db, config.DB)
	ctx := context.Background()
	if err := newModels.PrepareSqlite(ctx, config.DB.Timeout); err != nil {
		return fmt.Errorf("run: prepare sqlite failed: %w", err)
	}

	// setup server
	server := api.NewServer(*config, *newModels)
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
