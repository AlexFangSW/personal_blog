package main

import (
	"blog/api/handlers"
	"blog/config"
	"blog/db/models/sqlite"
	"blog/entities"
	"blog/repositories"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func getUserRepo(configPath string) (*repositories.Users, *sql.DB, error) {
	// load config file
	rawConfig, err := os.ReadFile(configPath)
	if err != nil {
		return &repositories.Users{}, &sql.DB{}, fmt.Errorf("getUserRepo: load config failed: %w", err)
	}
	config := config.NewConfig()
	json.Unmarshal(rawConfig, config)

	// db connection
	db, err := sql.Open("sqlite3", config.DB.DSNURL)
	if err != nil {
		return &repositories.Users{}, &sql.DB{}, fmt.Errorf("getUserRepo: open db connection failed: %w", err)
	}
	db.SetMaxOpenConns(config.DB.Connections)

	// db prepare
	model := sqlite.New(db, config.DB)
	ctx := context.Background()
	if err := model.Prepare(ctx); err != nil {
		return &repositories.Users{}, &sql.DB{}, fmt.Errorf("getUserRepo: model prepare failed: %w", err)
	}

	// models
	usersModel := sqlite.NewUsers()

	// repositories
	usersRepoModels := repositories.NewUsersRepoModels(
		usersModel,
	)
	return repositories.NewUsers(db, config.DB, *usersRepoModels), db, nil
}

func genBasicAuth(username, password string) {
	cred := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
	fmt.Println("Use this in header to get jwt token on login")
	fmt.Printf("Authorization: Basic %s\n", cred)
}

func run() error {
	// flags
	configPath := flag.String("config", "./config.json", "Config filepath")
	createUser := flag.Bool("create", false, "Create user")
	updateUser := flag.Bool("update", false, "Update user")
	deleteUser := flag.Bool("delete", false, "Delete user")
	flag.Parse()
	slog.Info("load config", "path:", *configPath)

	username := flag.Arg(0)
	password := flag.Arg(1)

	fmt.Printf("username: %q\n", username)
	fmt.Printf("password: %q\n", password)

	userRepo, db, err := getUserRepo(*configPath)
	defer db.Close()
	if err != nil {
		return fmt.Errorf("run: get user repo failed: %w", err)
	}

	// process
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if *createUser {
		if username == "" || password == "" {
			return fmt.Errorf("run: must provide username and password for create")
		}
		hashedPassword, err := handlers.HashPassword(password)
		if err != nil {
			return fmt.Errorf("run: hash password failed: %w", err)
		}

		inUser := entities.NewInUser(username, hashedPassword)
		newUser, err := userRepo.Create(ctxTimeout, *inUser)
		if err != nil {
			return fmt.Errorf("run: create user failed: %w", err)
		}
		fmt.Println("User created !!", newUser)
		genBasicAuth(username, password)

	} else if *updateUser {
		if username == "" || password == "" {
			return fmt.Errorf("run: must provide username and password for create")
		}
		hashedPassword, err := handlers.HashPassword(password)
		if err != nil {
			return fmt.Errorf("run: hash password failed: %w", err)
		}

		inUser := entities.NewInUser(username, hashedPassword)
		newUser, err := userRepo.Update(ctxTimeout, *inUser)
		if err != nil {
			return fmt.Errorf("run: update user failed: %w", err)
		}
		fmt.Println("User updated !!", newUser)
		genBasicAuth(username, password)

	} else if *deleteUser {
		if err := userRepo.Delete(ctxTimeout); err != nil {
			return fmt.Errorf("run: update user failed: %w", err)
		}
		fmt.Println("User deleted !!")
	}

	return nil
}

// For User CRUD
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
