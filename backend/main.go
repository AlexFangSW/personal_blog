package main

import (
	"blog/api"
	"blog/api/handlers"
	"blog/config"
	"blog/db/models/sqlite"
	"blog/repositories"
	"blog/util"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	db.SetMaxOpenConns(config.DB.Connections)
	defer db.Close()

	// db prepare
	model := sqlite.New(db, config.DB)
	ctx := context.Background()
	if err := model.Prepare(ctx); err != nil {
		return fmt.Errorf("run: model prepare failed: %w", err)
	}

	// models
	blogsModel := sqlite.NewBlogs()
	blogTagsModel := sqlite.NewBlogTags()
	blogTopicsModel := sqlite.NewBlogTopics()
	TagsModel := sqlite.NewTags()
	TopicsModel := sqlite.NewTopics()

	// repositories
	blogsRepoModels := repositories.NewBlogsRepoModels(
		blogsModel,
		blogTagsModel,
		blogTopicsModel,
		TagsModel,
		TopicsModel,
	)
	blogsRepo := repositories.NewBlogs(db, config.DB, *blogsRepoModels)

	tagsRepoModels := repositories.NewTagsRepoModels(
		blogTagsModel,
		TagsModel,
	)
	tagsRepo := repositories.NewTags(db, config.DB, *tagsRepoModels)

	topicsRepoModels := repositories.NewTopicsRepoModels(
		blogTopicsModel,
		TopicsModel,
	)
	topicsRepo := repositories.NewTopics(db, config.DB, *topicsRepoModels)

	// handlers
	blogsHandler := handlers.NewBlogs(blogsRepo)
	tagsHandler := handlers.NewTags(tagsRepo)
	topicsHandler := handlers.NewTopics(topicsRepo)

	// setup server
	server := api.NewServer(
		config.Server,
		*blogsHandler,
		*tagsHandler,
		*topicsHandler,
	)

	// start server
	go func() {
		if err := server.Start(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("run: server failed:", err)
		}
		slog.Info("run: server gracfully stoped")
	}()

	// graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownTimeout, shutdownCancel := context.WithTimeout(
		ctx,
		time.Duration(config.Server.ShutdownTimeout)*time.Second,
	)
	defer shutdownCancel()

	return server.Stop(shutdownTimeout)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
