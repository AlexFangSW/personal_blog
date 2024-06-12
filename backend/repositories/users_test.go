package repositories_test

import (
	"blog/config"
	"blog/db"
	"blog/db/models/sqlite"
	"blog/entities"
	"blog/repositories"
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func TestUsersCreateSqlite(t *testing.T) {
	// connect
	dbConn, err := sql.Open("sqlite3", "file:test.db?mode=memory&_foreign_keys=on")
	if err != nil {
		t.Fatalf("TestUsersCreateSqlite: open db connection failed: %s", err)
	}
	defer dbConn.Close()

	// migrate db
	downDB(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite")
	if err := upDB(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite"); err != nil {
		t.Fatalf("TestUsersCreateSqlite: migrate up failed: %s", err)
	}

	// setup repo
	usersModel := sqlite.NewUsers()
	usersRepoModels := repositories.NewUsersRepoModels(usersModel)
	usersRepo := repositories.NewUsers(dbConn, config.NewConfig().DB, *usersRepoModels)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// test create
	user1 := entities.NewInUser("username1", "password1")
	newUser1, err := usersRepo.Create(ctxTimeout, *user1)
	if err != nil {
		t.Fatalf("TestUsersCreateSqlite: create failed: %s", err)
	}
	if !compareInUser(*user1, *newUser1) {
		t.Fatalf("TestUsersCreateSqlite: create cmp failed")
	}

	// test create failed, only one user should exist
	user2 := entities.NewInUser("username2", "password2")
	_, err2 := usersRepo.Create(ctxTimeout, *user2)
	if err2 == nil {
		t.Fatalf("TestUsersCreateSqlite: create should have failed")
	}
}

func TestUsersUpdateSqlite(t *testing.T) {
	// connect
	dbConn, err := sql.Open("sqlite3", "file:test.db?mode=memory&_foreign_keys=on")
	if err != nil {
		t.Fatalf("TestTopicsUpdateSqlite: open db connection failed: %s", err)
	}
	defer dbConn.Close()

	// migrate db
	downDB(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite")
	if err := upDB(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite"); err != nil {
		t.Fatalf("TestTopicsUpdateSqlite: migrate up failed: %s", err)
	}

	// setup repo
	usersModel := sqlite.NewUsers()
	usersRepoModels := repositories.NewUsersRepoModels(usersModel)
	usersRepo := repositories.NewUsers(dbConn, config.NewConfig().DB, *usersRepoModels)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// fill in rows, just assume they will succeed
	user1 := entities.NewInUser("username1", "password1")
	usersRepo.Create(ctxTimeout, *user1)

	// test update
	user2 := entities.NewInUser("username2", "password2")
	newUser2, err := usersRepo.Update(ctxTimeout, *user2)
	if err != nil {
		t.Fatalf("TestUsersUpdateSqlite: update failed: %s", err)
	}

	if !compareInUser(*user2, *newUser2) {
		t.Fatalf("TestUsersUpdateSqlite: update cmp failed")
	}

	// jwt should be cleared after an user update
	newjwt := "aabbb.fafsa.fdsaf"
	if err := usersRepo.UpdateJWT(ctxTimeout, newjwt); err != nil {
		t.Fatalf("TestUsersClearJWTSqlite: get failed: %s", err)
	}
	user3 := entities.NewInUser("username3", "password3")
	newUser3, err := usersRepo.Update(ctxTimeout, *user3)
	if err != nil {
		t.Fatalf("TestUsersUpdateSqlite: update failed: %s", err)
	}
	if newUser3.JWT != "" {
		t.Fatalf("TestUsersUpdateSqlite: jwt should be cleared")
	}
}

func TestUsersDeleteSqlite(t *testing.T) {
	// connect
	dbConn, err := sql.Open("sqlite3", "file:test.db?mode=memory&_foreign_keys=on")
	if err != nil {
		t.Fatalf("TestUsersDeleteSqlite: open db connection failed: %s", err)
	}
	defer dbConn.Close()

	// migrate db
	downDB(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite")
	if err := upDB(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite"); err != nil {
		t.Fatalf("TestUsersDeleteSqlite: migrate up failed: %s", err)
	}

	// setup repo
	usersModel := sqlite.NewUsers()
	usersRepoModels := repositories.NewUsersRepoModels(usersModel)
	usersRepo := repositories.NewUsers(dbConn, config.NewConfig().DB, *usersRepoModels)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// fill in rows, just assume they will succeed
	user1 := entities.NewInUser("username1", "password1")
	usersRepo.Create(ctxTimeout, *user1)

	// test delete
	if err := usersRepo.Delete(ctxTimeout); err != nil {
		t.Fatalf("TestUsersDeleteSqlite: delete failed: %s", err)
	}

	// is should be find deleting multiple times
	if err := usersRepo.Delete(ctxTimeout); err != nil {
		t.Fatalf("TestUsersDeleteSqlite: delete failed: %s", err)
	}
}

// This only compares username and password
func compareInUser(inUser entities.InUser, user entities.User) bool {
	if inUser.Name != user.Name {
		return false
	}
	if inUser.Password != user.Password {
		return false
	}

	return true
}

func TestUsersGetSqlite(t *testing.T) {
	// connect
	dbConn, err := sql.Open("sqlite3", "file:test.db?mode=memory&_foreign_keys=on")
	if err != nil {
		t.Fatalf("TestUsersDeleteSqlite: open db connection failed: %s", err)
	}
	defer dbConn.Close()

	// migrate db
	downDB(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite")
	if err := upDB(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite"); err != nil {
		t.Fatalf("TestUsersDeleteSqlite: migrate up failed: %s", err)
	}

	// setup repo
	usersModel := sqlite.NewUsers()
	usersRepoModels := repositories.NewUsersRepoModels(usersModel)
	usersRepo := repositories.NewUsers(dbConn, config.NewConfig().DB, *usersRepoModels)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// test get fail
	_, err2 := usersRepo.Get(ctxTimeout)
	if err2 == nil {
		t.Fatalf("TestUsersGetSqlite: this should have failed: %s", err2)
	}

	// fill in rows, just assume they will succeed
	user1 := entities.NewInUser("username1", "password1")
	usersRepo.Create(ctxTimeout, *user1)

	// test get
	user2, err2 := usersRepo.Get(ctxTimeout)
	if err2 != nil {
		t.Fatalf("TestUsersGetSqlite: get failed: %s", err2)
	}
	if !compareInUser(*user1, *user2) {
		t.Fatalf("TestUsersGetSqlite: get cmp failed")
	}
}

func TestUsersUpdateJWTSqlite(t *testing.T) {
	// connect
	dbConn, err := sql.Open("sqlite3", "file:test.db?mode=memory&_foreign_keys=on")
	if err != nil {
		t.Fatalf("TestUsersUpdateJWTSqlite: open db connection failed: %s", err)
	}
	defer dbConn.Close()

	// migrate db
	downDB(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite")
	if err := upDB(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite"); err != nil {
		t.Fatalf("TestUsersUpdateJWTSqlite: migrate up failed: %s", err)
	}

	// setup repo
	usersModel := sqlite.NewUsers()
	usersRepoModels := repositories.NewUsersRepoModels(usersModel)
	usersRepo := repositories.NewUsers(dbConn, config.NewConfig().DB, *usersRepoModels)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// fill in rows, just assume they will succeed
	user1 := entities.NewInUser("username1", "password1")
	usersRepo.Create(ctxTimeout, *user1)

	// test update jwt
	newjwt := "aabbb.fafsa.fdsaf"
	if err := usersRepo.UpdateJWT(ctxTimeout, newjwt); err != nil {
		t.Fatalf("TestUsersUpdateJWTSqlite: get failed: %s", err)
	}

	// check
	user2, _ := usersRepo.Get(ctxTimeout)
	if user2.JWT != newjwt {
		t.Fatalf("TestUsersUpdateJWTSqlite: jwt was not updated")
	}
}

func TestUsersClearJWTSqlite(t *testing.T) {
	// connect
	dbConn, err := sql.Open("sqlite3", "file:test.db?mode=memory&_foreign_keys=on")
	if err != nil {
		t.Fatalf("TestUsersClearJWTSqlite: open db connection failed: %s", err)
	}
	defer dbConn.Close()

	// migrate db
	downDB(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite")
	if err := upDB(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite"); err != nil {
		t.Fatalf("TestUsersClearJWTSqlite: migrate up failed: %s", err)
	}

	// setup repo
	usersModel := sqlite.NewUsers()
	usersRepoModels := repositories.NewUsersRepoModels(usersModel)
	usersRepo := repositories.NewUsers(dbConn, config.NewConfig().DB, *usersRepoModels)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// fill in rows, just assume they will succeed
	user1 := entities.NewInUser("username1", "password1")
	usersRepo.Create(ctxTimeout, *user1)

	// test update jwt
	newjwt := "aabbb.fafsa.fdsaf"
	if err := usersRepo.UpdateJWT(ctxTimeout, newjwt); err != nil {
		t.Fatalf("TestUsersClearJWTSqlite: get failed: %s", err)
	}

	// check
	user2, _ := usersRepo.Get(ctxTimeout)
	if user2.JWT != newjwt {
		t.Fatalf("TestUsersClearJWTSqlite: jwt was not updated")
	}

	// clear jwt
	if err := usersRepo.ClearJWT(ctxTimeout); err != nil {
		t.Fatalf("TestUsersClearJWTSqlite: get failed: %s", err)
	}

	// check
	user3, _ := usersRepo.Get(ctxTimeout)
	if user3.JWT != "" {
		t.Fatalf("TestUsersClearJWTSqlite: jwt should be empty")
	}
}
