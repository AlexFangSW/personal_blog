package main

import (
	"blog/api"
	"blog/db/models"
	"blog/structs"
	"blog/util"
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
		return fmt.Errorf("Run: load config error: %w", err)
	}
	config := structs.NewConfig()
	json.Unmarshal(rawConfig, config)

	// init logger
	util.InitLogger(config.Logger.Level)

	// db connection
	db, err := sql.Open("sqlite3", config.DB.DSNURL)
	if err != nil {
		return fmt.Errorf("Run: failed to open db connection: %w", err)
	}

	models := models.NewModels(db, config.DB)

	// setup server
	server := api.NewServer(*config, *models)
	if err := server.Start(); err != nil {
		return fmt.Errorf("Run: server error: %w", err)
	}

	return nil
}

// TODO: graceful shutdown
// TODO: activate sqlite foreign key
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
