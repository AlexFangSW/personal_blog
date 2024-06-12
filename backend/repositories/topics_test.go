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

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	_ "github.com/mattn/go-sqlite3"
)

func TestTopicsCreateSqlite(t *testing.T) {
	// connect
	dbConn, err := sql.Open("sqlite3", "file:test.db?mode=memory&_foreign_keys=on")
	if err != nil {
		t.Fatalf("TestTopicsCreateSqlite: open db connection failed: %s", err)
	}
	defer dbConn.Close()

	// migrate db
	downDB(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite")
	if err := upDB(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite"); err != nil {
		t.Fatalf("TestTopicsCreateSqlite: migrate up failed: %s", err)
	}

	// setup repo
	topicsModel := sqlite.NewTopics()
	blogTopicsModel := sqlite.NewBlogTopics()
	topicsRepoModels := repositories.NewTopicsRepoModels(blogTopicsModel, topicsModel)
	topicsRepo := repositories.NewTopics(dbConn, config.NewConfig().DB, *topicsRepoModels)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// test create
	entry1 := entities.NewTopic("name 1", "desc 1")
	newEntry1, err := topicsRepo.Create(ctxTimeout, *entry1)
	if err != nil {
		t.Fatalf("TestTopicsCreateSqlite: create failed: %s", err)
	}
	if !cmp.Equal(entry1, newEntry1, cmpopts.IgnoreFields(entities.Topic{}, "ID", "Created_at", "Updated_at")) {
		t.Fatalf("TestTopicsCreateSqlite: create cmp failed")
	}

	// test create failed
	entry2 := entry1
	_, err2 := topicsRepo.Create(ctxTimeout, *entry2)
	if err2 == nil {
		t.Fatalf("TestTopicsCreateSqlite: create should have failed")
	}

}

func TestTopicsUpdateSqlite(t *testing.T) {
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
	topicsModel := sqlite.NewTopics()
	blogTopicsModel := sqlite.NewBlogTopics()
	topicsRepoModels := repositories.NewTopicsRepoModels(blogTopicsModel, topicsModel)
	topicsRepo := repositories.NewTopics(dbConn, config.NewConfig().DB, *topicsRepoModels)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// fill in rows, just assume they will succeed
	topic1 := entities.NewTopic("name 1", "desc 1")
	topicsRepo.Create(ctxTimeout, *topic1)

	topic2 := entities.NewTopic("name 2", "desc 2")
	topicsRepo.Create(ctxTimeout, *topic2)

	// test update
	topic3 := entities.NewTopic("updated name 3", "updated desc 3")
	newTopic3, err := topicsRepo.Update(ctxTimeout, *topic3, 1)
	if err != nil {
		t.Fatalf("TestTopicsUpdateSqlite: update failed: %s", err)
	}
	if !cmp.Equal(topic3, newTopic3, cmpopts.IgnoreFields(entities.Topic{}, "ID", "Created_at", "Updated_at")) {
		t.Fatalf("TestTopicsUpdateSqlite: update cmp failed")
	}

	// test update failed becuse of unique constraint
	topic4 := topic3
	_, err4 := topicsRepo.Update(ctxTimeout, *topic4, 2)
	if err4 == nil {
		t.Fatalf("TestTopicsUpdateSqlite: update sould have failed failed")
	}
}

func TestTopicsDeleteSqlite(t *testing.T) {
	// connect
	dbConn, err := sql.Open("sqlite3", "file:test.db?mode=memory&_foreign_keys=on")
	if err != nil {
		t.Fatalf("TestTopicsDeleteSqlite: open db connection failed: %s", err)
	}
	defer dbConn.Close()

	// migrate db
	downDB(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite")
	if err := upDB(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite"); err != nil {
		t.Fatalf("TestTopicsDeleteSqlite: migrate up failed: %s", err)
	}

	// setup repo
	topicsModel := sqlite.NewTopics()
	blogTopicsModel := sqlite.NewBlogTopics()
	topicsRepoModels := repositories.NewTopicsRepoModels(blogTopicsModel, topicsModel)
	topicsRepo := repositories.NewTopics(dbConn, config.NewConfig().DB, *topicsRepoModels)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// fill in rows, just assume they will succeed
	topic1 := entities.NewTopic("name 1", "desc 1")
	topicsRepo.Create(ctxTimeout, *topic1)

	// test delete
	affectedRows, err := topicsRepo.Delete(ctxTimeout, 1)
	if err != nil {
		t.Fatalf("TestTopicsDeleteSqlite: delete failed: %s", err)
	}
	if affectedRows == 0 {
		t.Fatalf("TestTopicsDeleteSqlite: affected rows should not be zero")
	}

	// test delete no affected rows
	affectedRows2, err2 := topicsRepo.Delete(ctxTimeout, 1)
	if err2 != nil {
		t.Fatalf("TestTopicsDeleteSqlite: delete failed: %s", err)
	}
	if affectedRows2 != 0 {
		t.Fatalf("TestTopicsDeleteSqlite: affected rows should BE zero")
	}

	// TODO: should fail with foreign key constraint
}
