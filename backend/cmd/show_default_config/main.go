package main

import (
	"blog/config"
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	cfg := config.NewConfig()
	bytes, err := json.Marshal(cfg)
	if err != nil {
		log.Fatal("json marshal config failed", err)
	}
	fmt.Println(string(bytes))
}
