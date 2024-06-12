package repositories_test

import (
	"blog/config"
	"blog/db"
	"blog/db/models/sqlite"
	"blog/entities"
	"blog/repositories"
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

func upDB(dbConn *sql.DB, migrations fs.FS, dialect string, path string) error {
	goose.SetLogger(goose.NopLogger())
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect(dialect); err != nil {
		return fmt.Errorf("upDB: set dialect failed: %w", err)
	}

	if err := goose.Up(dbConn, path); err != nil {
		return fmt.Errorf("upDB: up failed: %w", err)
	}

	return nil
}

func downDB(dbConn *sql.DB, migrations fs.FS, dialect string, path string) error {
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect(dialect); err != nil {
		return fmt.Errorf("downDB: set dialect failed: %w", err)
	}

	if err := goose.Down(dbConn, path); err != nil {
		return fmt.Errorf("downDB: down failed: %w", err)
	}

	return nil
}

func TestTagsCreateSqlite(t *testing.T) {
	// connect
	dbConn, err := sql.Open("sqlite3", "file:test.db?mode=memory&_foreign_keys=on")
	if err != nil {
		t.Fatalf("TestTagsCreateSqlite: open db connection failed: %s", err)
	}
	defer dbConn.Close()

	// migrate db
	downDB(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite")
	if err := upDB(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite"); err != nil {
		t.Fatalf("TestTagsCreateSqlite: migrate up failed: %s", err)
	}

	// setup repo
	tagsModel := sqlite.NewTags()
	blogTagsModel := sqlite.NewBlogTags()
	tagsRepoModels := repositories.NewTagsRepoModels(blogTagsModel, tagsModel)
	tagsRepo := repositories.NewTags(dbConn, config.NewConfig().DB, *tagsRepoModels)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// test create
	tag1 := entities.NewTag("name 1", "desc 1")
	newTag1, err := tagsRepo.Create(ctxTimeout, *tag1)
	if err != nil {
		t.Fatalf("TestTagsCreateSqlite: create failed: %s", err)
	}
	if !cmp.Equal(tag1, newTag1, cmpopts.IgnoreFields(entities.Tag{}, "ID", "Created_at", "Updated_at")) {
		t.Fatalf("TestTagsCreateSqlite: create cmp failed")
	}

	// test create failed
	tag2 := tag1
	_, err2 := tagsRepo.Create(ctxTimeout, *tag2)
	if err2 == nil {
		t.Fatalf("TestTagsCreateSqlite: create should have failed")
	}

}

func TestTagsUpdateSqlite(t *testing.T) {
	// connect
	dbConn, err := sql.Open("sqlite3", "file:test.db?mode=memory&_foreign_keys=on")
	if err != nil {
		t.Fatalf("TestTagsUpdateSqlite: open db connection failed: %s", err)
	}
	defer dbConn.Close()

	// migrate db
	downDB(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite")
	if err := upDB(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite"); err != nil {
		t.Fatalf("TestTagsUpdateSqlite: migrate up failed: %s", err)
	}

	// setup repo
	tagsModel := sqlite.NewTags()
	blogTagsModel := sqlite.NewBlogTags()
	tagsRepoModels := repositories.NewTagsRepoModels(blogTagsModel, tagsModel)
	tagsRepo := repositories.NewTags(dbConn, config.NewConfig().DB, *tagsRepoModels)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// fill in rows, just assume they will succeed
	tag1 := entities.NewTag("name 1", "desc 1")
	tagsRepo.Create(ctxTimeout, *tag1)

	tag2 := entities.NewTag("name 2", "desc 2")
	tagsRepo.Create(ctxTimeout, *tag2)

	// test update
	tag3 := entities.NewTag("updated name 3", "updated desc 3")
	newTag3, err := tagsRepo.Update(ctxTimeout, *tag3, 1)
	if err != nil {
		t.Fatalf("TestTagsUpdateSqlite: update failed: %s", err)
	}
	if !cmp.Equal(tag3, newTag3, cmpopts.IgnoreFields(entities.Tag{}, "ID", "Created_at", "Updated_at")) {
		t.Fatalf("TestTagsUpdateSqlite: update cmp failed")
	}

	// test update failed becuse of unique constraint
	tag4 := tag3
	_, err4 := tagsRepo.Update(ctxTimeout, *tag4, 2)
	if err4 == nil {
		t.Fatalf("TestTagsUpdateSqlite: update sould have failed failed")
	}
}

func TestTagsDeleteSqlite(t *testing.T) {
	// connect
	dbConn, err := sql.Open("sqlite3", "file:test.db?mode=memory&_foreign_keys=on")
	if err != nil {
		t.Fatalf("TestTagsDeleteSqlite: open db connection failed: %s", err)
	}
	defer dbConn.Close()

	// migrate db
	downDB(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite")
	if err := upDB(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite"); err != nil {
		t.Fatalf("TestTagsDeleteSqlite: migrate up failed: %s", err)
	}

	// setup repo
	tagsModel := sqlite.NewTags()
	blogTagsModel := sqlite.NewBlogTags()
	tagsRepoModels := repositories.NewTagsRepoModels(blogTagsModel, tagsModel)
	tagsRepo := repositories.NewTags(dbConn, config.NewConfig().DB, *tagsRepoModels)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// fill in rows, just assume they will succeed
	tag1 := entities.NewTag("name 1", "desc 1")
	tagsRepo.Create(ctxTimeout, *tag1)

	// test delete
	affectedRows, err := tagsRepo.Delete(ctxTimeout, 1)
	if err != nil {
		t.Fatalf("TestTagsDeleteSqlite: delete failed: %s", err)
	}
	if affectedRows == 0 {
		t.Fatalf("TestTagsDeleteSqlite: affected rows should not be zero")
	}

	// test delete no affected rows
	affectedRows2, err2 := tagsRepo.Delete(ctxTimeout, 1)
	if err2 != nil {
		t.Fatalf("TestTagsDeleteSqlite: delete failed: %s", err)
	}
	if affectedRows2 != 0 {
		t.Fatalf("TestTagsDeleteSqlite: affected rows should BE zero")
	}

	// TODO: should fail with foreign key constraint
}
