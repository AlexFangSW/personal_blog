package main

import (
	"blog/api"
	"blog/structs"
	"blog/util"
	"encoding/json"
	"flag"
	"log"
	"log/slog"
	"os"
)

func main() {
	log.SetFlags(log.Ltime | log.Llongfile | log.LstdFlags)

	// flags
	configPath := flag.String("config", "./config.json", "Config filepath")
	flag.Parse()
	slog.Info("load config", "path:", *configPath)

	// load config file
	rawConfig, err := os.ReadFile(*configPath)
	if err != nil {
		log.Fatal("load config error: ", err)
	}
	config := structs.NewConfig()
	json.Unmarshal(rawConfig, config)

	// init logger
	util.InitLogger(config.Logger.Level)

	// setup server
	server := api.NewServer(*config)
	log.Fatal(server.Start())
}
