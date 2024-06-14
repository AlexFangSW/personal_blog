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
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	_ "github.com/mattn/go-sqlite3"
)

// Prepres blogs, tags, and topics repo
func prepareRepos(dbConn *sql.DB) (repositories.Blogs, repositories.Tags, repositories.Topics) {
	blogsModel := sqlite.NewBlogs()
	blogTagsModel := sqlite.NewBlogTags()
	blogTopicsModel := sqlite.NewBlogTopics()
	tagsModel := sqlite.NewTags()
	topicsModel := sqlite.NewTopics()

	topicsRepoModels := repositories.NewTopicsRepoModels(blogTopicsModel, topicsModel)
	topicsRepo := repositories.NewTopics(dbConn, config.NewConfig().DB, *topicsRepoModels)

	tagsRepoModels := repositories.NewTagsRepoModels(blogTagsModel, tagsModel)
	tagsRepo := repositories.NewTags(dbConn, config.NewConfig().DB, *tagsRepoModels)

	blogsRepoModels := repositories.NewBlogsRepoModels(
		blogsModel,
		blogTagsModel,
		blogTopicsModel,
		tagsModel,
		topicsModel,
	)
	blogsRepo := repositories.NewBlogs(dbConn, config.NewConfig().DB, *blogsRepoModels)

	return *blogsRepo, *tagsRepo, *topicsRepo
}

func TestBlogsCreateSqlite(t *testing.T) {
	// connect
	dbConn, err := sql.Open("sqlite3", "file:test.db?mode=memory&_foreign_keys=on")
	if err != nil {
		t.Fatalf("TestBlogsCreateSqlite: open db connection failed: %s", err)
	}
	defer dbConn.Close()

	// migrate db
	if err := db.Up(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite"); err != nil {
		t.Fatalf("TestBlogsCreateSqlite migrate up failed: %s", err)
	}

	// setup repo
	blogsRepo, tagsRepo, topicsRepo := prepareRepos(dbConn)
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// prepare topic and tags
	topic1, _ := topicsRepo.Create(ctxTimeout, *entities.NewTopic("topic1", "topic1"))
	topic2, _ := topicsRepo.Create(ctxTimeout, *entities.NewTopic("topic2", "topic2"))
	tag1, _ := tagsRepo.Create(ctxTimeout, *entities.NewTag("tag1", "tag1"))
	tag2, _ := tagsRepo.Create(ctxTimeout, *entities.NewTag("tag2", "tag2"))

	// test create
	newBlog := entities.NewBlog(
		"title1",
		"content1",
		"description1",
		false,
		true,
	)
	newInBlog := entities.NewInBlog(
		*newBlog,
		[]int{1, 2},
		[]int{1, 2},
	)
	createResult1, err := blogsRepo.Create(ctxTimeout, *newInBlog)
	if err != nil {
		t.Fatalf("TestBlogsCreateSqlite: create failed: %s", err)
	}
	if !cmp.Equal(newBlog, &createResult1.Blog, cmpopts.IgnoreFields(entities.Blog{}, "ID", "Created_at", "Updated_at")) {
		t.Fatalf("TestBlogsCreateSqlite: create cmp blog failed")
	}
	if !cmp.Equal(topic1, &createResult1.Topics[0]) {
		t.Fatalf("TestBlogsCreateSqlite: create cmp topics 1 failed")
	}
	if !cmp.Equal(topic2, &createResult1.Topics[1]) {
		t.Fatalf("TestBlogsCreateSqlite: create cmp topics 2 failed")
	}
	if !cmp.Equal(tag1, &createResult1.Tags[0]) {
		t.Fatalf("TestBlogsCreateSqlite: create cmp tags 1 failed")
	}
	if !cmp.Equal(tag2, &createResult1.Tags[1]) {
		t.Fatalf("TestBlogsCreateSqlite: create cmp tags 2 failed")
	}

	// create blog fail with unique contraint
	newInBlog2 := newInBlog
	_, err2 := blogsRepo.Create(ctxTimeout, *newInBlog2)
	if err2 == nil {
		t.Fatalf("TestBlogsCreateSqlite: create should have failed")
	}

	// create blog fail with foreign key contraint (topic)
	newInBlog3 := newInBlog
	newInBlog3.Topics = []int{1, 2, 5}
	_, err3 := blogsRepo.Create(ctxTimeout, *newInBlog3)
	if err3 == nil {
		t.Fatalf("TestBlogsCreateSqlite: create should have failed")
	}

	// create blog fail with foreign key contraint (tag)
	newInBlog4 := newInBlog
	newInBlog4.Tags = []int{1, 2, 5}
	_, err4 := blogsRepo.Create(ctxTimeout, *newInBlog4)
	if err4 == nil {
		t.Fatalf("TestBlogsCreateSqlite: create should have failed")
	}
}

func TestBlogsUpdateSqlite(t *testing.T) {
	// connect
	dbConn, err := sql.Open("sqlite3", "file:test.db?mode=memory&_foreign_keys=on")
	if err != nil {
		t.Fatalf("TestBlogsUpdateSqlite: open db connection failed: %s", err)
	}
	defer dbConn.Close()

	// migrate db
	if err := db.Up(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite"); err != nil {
		t.Fatalf("TestBlogsUpdateSqlite migrate up failed: %s", err)
	}

	// setup repo
	blogsRepo, tagsRepo, topicsRepo := prepareRepos(dbConn)
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// lets just assume they work
	// prepare topic and tags
	topic1, _ := topicsRepo.Create(ctxTimeout, *entities.NewTopic("topic1", "topic1"))
	topicsRepo.Create(ctxTimeout, *entities.NewTopic("topic2", "topic2"))
	tag1, _ := tagsRepo.Create(ctxTimeout, *entities.NewTag("tag1", "tag1"))
	tagsRepo.Create(ctxTimeout, *entities.NewTag("tag2", "tag2"))

	// prepare blog
	newBlog := entities.NewBlog(
		"title1",
		"content1",
		"description1",
		false,
		true,
	)
	newInBlog := entities.NewInBlog(
		*newBlog,
		[]int{1, 2},
		[]int{1, 2},
	)
	blogsRepo.Create(ctxTimeout, *newInBlog)

	// update blog
	newBlog2 := entities.NewBlog(
		"title2",
		"content2",
		"description2",
		true,
		false,
	)
	newInBlog2 := entities.NewInBlog(
		*newBlog2,
		[]int{1},
		[]int{1},
	)
	updatedBlog, err := blogsRepo.Update(ctxTimeout, *newInBlog2, 1)
	if err != nil {
		t.Fatalf("TestBlogsUpdateSqlite: update failed: %s", err)
	}
	if !cmp.Equal(newBlog2, &updatedBlog.Blog, cmpopts.IgnoreFields(entities.Blog{}, "ID", "Created_at", "Updated_at")) {
		t.Fatalf("TestBlogsUpdateSqlite: update cmp blog failed")
	}
	if len(updatedBlog.Topics) != 1 {
		t.Fatalf("TestBlogsUpdateSqlite: update cmp topics len failed")
	}
	if len(updatedBlog.Tags) != 1 {
		t.Fatalf("TestBlogsUpdateSqlite: update cmp tags len failed")
	}
	if !cmp.Equal(topic1, &updatedBlog.Topics[0]) {
		t.Fatalf("TestBlogsUpdateSqlite: update cmp topics 1 failed")
	}
	if !cmp.Equal(tag1, &updatedBlog.Tags[0]) {
		t.Fatalf("TestBlogsUpdateSqlite: update cmp tags 1 failed")
	}
}

func TestBlogsGetSqlite(t *testing.T) {
	// connect
	dbConn, err := sql.Open("sqlite3", "file:test.db?mode=memory&_foreign_keys=on")
	if err != nil {
		t.Fatalf("TestBlogsGetSqlite: open db connection failed: %s", err)
	}
	defer dbConn.Close()

	// migrate db
	if err := db.Up(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite"); err != nil {
		t.Fatalf("TestBlogsGetSqlite: migrate up failed: %s", err)
	}

	// setup repo
	blogsRepo, tagsRepo, topicsRepo := prepareRepos(dbConn)
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// lets just assume they work
	// prepare topic and tags
	topic1, _ := topicsRepo.Create(ctxTimeout, *entities.NewTopic("topic1", "topic1"))
	tag1, _ := tagsRepo.Create(ctxTimeout, *entities.NewTag("tag1", "tag1"))

	// prepare blog
	visibleBlog := entities.NewBlog(
		"title1",
		"content1",
		"description1",
		false,
		true,
	)
	newInBlog1 := entities.NewInBlog(
		*visibleBlog,
		[]int{1},
		[]int{1},
	)
	blogsRepo.Create(ctxTimeout, *newInBlog1)

	notVisibleBlog := entities.NewBlog(
		"title2",
		"content2",
		"description2",
		false,
		false,
	)
	newInBlog2 := entities.NewInBlog(
		*notVisibleBlog,
		[]int{1},
		[]int{1},
	)
	blogsRepo.Create(ctxTimeout, *newInBlog2)

	// Get visible
	blog1, err := blogsRepo.Get(ctxTimeout, 1)
	if err != nil {
		t.Fatalf("TestBlogsGetSqlite: get failed: %s", err)
	}
	if !cmp.Equal(visibleBlog, &blog1.Blog, cmpopts.IgnoreFields(entities.Blog{}, "ID", "Created_at", "Updated_at")) {
		t.Fatalf("TestBlogsGetSqlite: get cmp blog failed")
	}
	if !cmp.Equal(topic1, &blog1.Topics[0]) {
		t.Fatalf("TestBlogsGetSqlite: get cmp topics 1 failed")
	}
	if !cmp.Equal(tag1, &blog1.Tags[0]) {
		t.Fatalf("TestBlogsGetSqlite: get cmp tags 1 failed")
	}

	// Get not visible ( should fail)
	_, err2 := blogsRepo.Get(ctxTimeout, 2)
	if err2 == nil {
		t.Fatalf("TestBlogsGetSqlite: get should not succeed failed: %s", err)
	}
}

func TestBlogsListSqlite(t *testing.T) {
	// connect
	dbConn, err := sql.Open("sqlite3", "file:test.db?mode=memory&_foreign_keys=on")
	if err != nil {
		t.Fatalf("TestBlogsListSqlite: open db connection failed: %s", err)
	}
	defer dbConn.Close()

	// migrate db
	if err := db.Up(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite"); err != nil {
		t.Fatalf("TestBlogsListSqlite: migrate up failed: %s", err)
	}

	// setup repo
	blogsRepo, tagsRepo, topicsRepo := prepareRepos(dbConn)
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// lets just assume they work
	// prepare topic and tags
	topic1, _ := topicsRepo.Create(ctxTimeout, *entities.NewTopic("topic1", "topic1"))
	tag1, _ := tagsRepo.Create(ctxTimeout, *entities.NewTag("tag1", "tag1"))

	// prepare blog
	visibleBlog := entities.NewBlog(
		"title1",
		"content1",
		"description1",
		false,
		true,
	)
	newInBlog1 := entities.NewInBlog(
		*visibleBlog,
		[]int{1},
		[]int{1},
	)
	blogsRepo.Create(ctxTimeout, *newInBlog1)

	notVisibleBlog := entities.NewBlog(
		"title2",
		"content2",
		"description2",
		false,
		false,
	)
	newInBlog2 := entities.NewInBlog(
		*notVisibleBlog,
		[]int{1},
		[]int{1},
	)
	blogsRepo.Create(ctxTimeout, *newInBlog2)

	// list visible
	blogs, err := blogsRepo.List(ctxTimeout)
	if err != nil {
		t.Fatalf("TestBlogsListSqlite: list failed: %s", err)
	}
	if len(blogs) != 1 {
		t.Fatalf("TestBlogsListSqlite: list should only return one")
	}
	// List will not return content (too large)
	if !cmp.Equal(visibleBlog, &blogs[0].Blog, cmpopts.IgnoreFields(entities.Blog{}, "ID", "Created_at", "Updated_at", "Content")) {
		t.Fatalf("TestBlogsListSqlite: list cmp blog failed")
	}
	if !cmp.Equal(topic1, &blogs[0].Topics[0]) {
		t.Fatalf("TestBlogsListSqlite: list cmp topics 1 failed")
	}
	if !cmp.Equal(tag1, &blogs[0].Tags[0]) {
		t.Fatalf("TestBlogsListSqlite: list cmp tags 1 failed")
	}
}

func TestBlogsListByTopicIDsSqlite(t *testing.T) {
	// connect
	dbConn, err := sql.Open("sqlite3", "file:test.db?mode=memory&_foreign_keys=on")
	if err != nil {
		t.Fatalf("TestBlogsListByTopicIDsSqlite: open db connection failed: %s", err)
	}
	defer dbConn.Close()

	// migrate db
	if err := db.Up(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite"); err != nil {
		t.Fatalf("TestBlogsListByTopicIDsSqlite: migrate up failed: %s", err)
	}

	// setup repo
	blogsRepo, tagsRepo, topicsRepo := prepareRepos(dbConn)
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// lets just assume they work
	// prepare topic and tags
	topic1, _ := topicsRepo.Create(ctxTimeout, *entities.NewTopic("topic1", "topic1"))
	tag1, _ := tagsRepo.Create(ctxTimeout, *entities.NewTag("tag1", "tag1"))

	// prepare blog
	visibleBlog := entities.NewBlog(
		"title1",
		"content1",
		"description1",
		false,
		true,
	)
	newInBlog1 := entities.NewInBlog(
		*visibleBlog,
		[]int{1},
		[]int{1},
	)
	blogsRepo.Create(ctxTimeout, *newInBlog1)

	notVisibleBlog := entities.NewBlog(
		"title2",
		"content2",
		"description2",
		false,
		false,
	)
	newInBlog2 := entities.NewInBlog(
		*notVisibleBlog,
		[]int{1},
		[]int{1},
	)
	blogsRepo.Create(ctxTimeout, *newInBlog2)

	// list
	blogs, err := blogsRepo.ListByTopicIDs(ctxTimeout, []int{1})
	if err != nil {
		t.Fatalf("TestBlogsListByTopicIDsSqlite: list failed: %s", err)
	}
	if len(blogs) != 1 {
		t.Fatalf("TestBlogsListByTopicIDsSqlite: list should only return one")
	}
	// List will not return content (too large)
	if !cmp.Equal(visibleBlog, &blogs[0].Blog, cmpopts.IgnoreFields(entities.Blog{}, "ID", "Created_at", "Updated_at", "Content")) {
		t.Fatalf("TestBlogsListByTopicIDsSqlite: list cmp blog failed")
	}
	if !cmp.Equal(topic1, &blogs[0].Topics[0]) {
		t.Fatalf("TestBlogsListByTopicIDsSqlite: list cmp topics 1 failed")
	}
	if !cmp.Equal(tag1, &blogs[0].Tags[0]) {
		t.Fatalf("TestBlogsListByTopicIDsSqlite: list cmp tags 1 failed")
	}

	// return empty slice
	blogs2, err := blogsRepo.ListByTopicIDs(ctxTimeout, []int{2})
	if err != nil {
		t.Fatalf("TestBlogsListByTopicIDsSqlite: list failed: %s", err)
	}
	if len(blogs2) != 0 {
		t.Fatalf("TestBlogsListByTopicIDsSqlite: should return a empty slice")
	}
}

func TestBlogsListByTopicAndTagIDsSqlite(t *testing.T) {
	// connect
	dbConn, err := sql.Open("sqlite3", "file:test.db?mode=memory&_foreign_keys=on")
	if err != nil {
		t.Fatalf("TestBlogsListByTopicAndTagIDsSqlite: open db connection failed: %s", err)
	}
	defer dbConn.Close()

	// migrate db
	if err := db.Up(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite"); err != nil {
		t.Fatalf("TestBlogsListByTopicAndTagIDsSqlite: migrate up failed: %s", err)
	}

	// setup repo
	blogsRepo, tagsRepo, topicsRepo := prepareRepos(dbConn)
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// lets just assume they work
	// prepare topic and tags
	topic1, _ := topicsRepo.Create(ctxTimeout, *entities.NewTopic("topic1", "topic1"))
	tag1, _ := tagsRepo.Create(ctxTimeout, *entities.NewTag("tag1", "tag1"))

	// prepare blog
	visibleBlog := entities.NewBlog(
		"title1",
		"content1",
		"description1",
		false,
		true,
	)
	newInBlog1 := entities.NewInBlog(
		*visibleBlog,
		[]int{1},
		[]int{1},
	)
	blogsRepo.Create(ctxTimeout, *newInBlog1)

	notVisibleBlog := entities.NewBlog(
		"title2",
		"content2",
		"description2",
		false,
		false,
	)
	newInBlog2 := entities.NewInBlog(
		*notVisibleBlog,
		[]int{1},
		[]int{1},
	)
	blogsRepo.Create(ctxTimeout, *newInBlog2)

	// list
	blogs, err := blogsRepo.ListByTopicAndTagIDs(ctxTimeout, []int{1}, []int{1})
	if err != nil {
		t.Fatalf("TestBlogsListByTopicAndTagIDsSqlite: list failed: %s", err)
	}
	if len(blogs) != 1 {
		t.Fatalf("TestBlogsListByTopicAndTagIDsSqlite: list should only return one")
	}
	// List will not return content (too large)
	if !cmp.Equal(visibleBlog, &blogs[0].Blog, cmpopts.IgnoreFields(entities.Blog{}, "ID", "Created_at", "Updated_at", "Content")) {
		t.Fatalf("TestBlogsListByTopicAndTagIDsSqlite: list cmp blog failed")
	}
	if !cmp.Equal(topic1, &blogs[0].Topics[0]) {
		t.Fatalf("TestBlogsListByTopicAndTagIDsSqlite: list cmp topics 1 failed")
	}
	if !cmp.Equal(tag1, &blogs[0].Tags[0]) {
		t.Fatalf("TestBlogsListByTopicAndTagIDsSqlite: list cmp tags 1 failed")
	}

	// return empty slice
	blogs2, err := blogsRepo.ListByTopicAndTagIDs(ctxTimeout, []int{2}, []int{1})
	if err != nil {
		t.Fatalf("TestBlogsListByTopicAndTagIDsSqlite: list failed: %s", err)
	}
	if len(blogs2) != 0 {
		t.Fatalf("TestBlogsListByTopicAndTagIDsSqlite: should return a empty slice")
	}
}

func TestBlogsAdminGetSqlite(t *testing.T) {
	// connect
	dbConn, err := sql.Open("sqlite3", "file:test.db?mode=memory&_foreign_keys=on")
	if err != nil {
		t.Fatalf("TestBlogsAdminGetSqlite: open db connection failed: %s", err)
	}
	defer dbConn.Close()

	// migrate db
	if err := db.Up(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite"); err != nil {
		t.Fatalf("TestBlogsAdminGetSqlite: migrate up failed: %s", err)
	}

	// setup repo
	blogsRepo, tagsRepo, topicsRepo := prepareRepos(dbConn)
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// lets just assume they work
	// prepare topic and tags
	topic1, _ := topicsRepo.Create(ctxTimeout, *entities.NewTopic("topic1", "topic1"))
	tag1, _ := tagsRepo.Create(ctxTimeout, *entities.NewTag("tag1", "tag1"))

	// prepare blog
	visibleBlog := entities.NewBlog(
		"title1",
		"content1",
		"description1",
		false,
		true,
	)
	newInBlog1 := entities.NewInBlog(
		*visibleBlog,
		[]int{1},
		[]int{1},
	)
	blogsRepo.Create(ctxTimeout, *newInBlog1)

	notVisibleBlog := entities.NewBlog(
		"title2",
		"content2",
		"description2",
		false,
		false,
	)
	newInBlog2 := entities.NewInBlog(
		*notVisibleBlog,
		[]int{1},
		[]int{1},
	)
	blogsRepo.Create(ctxTimeout, *newInBlog2)

	// Get visible
	blog1, err := blogsRepo.AdminGet(ctxTimeout, 1)
	if err != nil {
		t.Fatalf("TestBlogsAdminGetSqlite: get failed: %s", err)
	}
	if !cmp.Equal(visibleBlog, &blog1.Blog, cmpopts.IgnoreFields(entities.Blog{}, "ID", "Created_at", "Updated_at")) {
		t.Fatalf("TestBlogsAdminGetSqlite: get cmp blog failed")
	}
	if !cmp.Equal(topic1, &blog1.Topics[0]) {
		t.Fatalf("TestBlogsAdminGetSqlite: get cmp topics 1 failed")
	}
	if !cmp.Equal(tag1, &blog1.Tags[0]) {
		t.Fatalf("TestBlogsAdminGetSqlite: get cmp tags 1 failed")
	}

	// Get not visible ( should succeed)
	blog2, err := blogsRepo.AdminGet(ctxTimeout, 2)
	if err != nil {
		t.Fatalf("TestBlogsAdminGetSqlite: get failed: %s", err)
	}
	if !cmp.Equal(notVisibleBlog, &blog2.Blog, cmpopts.IgnoreFields(entities.Blog{}, "ID", "Created_at", "Updated_at")) {
		t.Fatalf("TestBlogsAdminGetSqlite: get cmp blog failed")
	}
	if !cmp.Equal(topic1, &blog2.Topics[0]) {
		t.Fatalf("TestBlogsAdminGetSqlite: get cmp topics 1 failed")
	}
	if !cmp.Equal(tag1, &blog2.Tags[0]) {
		t.Fatalf("TestBlogsAdminGetSqlite: get cmp tags 1 failed")
	}
}

func TestBlogsAdminListSqlite(t *testing.T) {
	// connect
	dbConn, err := sql.Open("sqlite3", "file:test.db?mode=memory&_foreign_keys=on")
	if err != nil {
		t.Fatalf("TestBlogsAdminListSqlite: open db connection failed: %s", err)
	}
	defer dbConn.Close()

	// migrate db
	if err := db.Up(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite"); err != nil {
		t.Fatalf("TestBlogsAdminListSqlite: migrate up failed: %s", err)
	}

	// setup repo
	blogsRepo, tagsRepo, topicsRepo := prepareRepos(dbConn)
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// lets just assume they work
	// prepare topic and tags
	topic1, _ := topicsRepo.Create(ctxTimeout, *entities.NewTopic("topic1", "topic1"))
	tag1, _ := tagsRepo.Create(ctxTimeout, *entities.NewTag("tag1", "tag1"))

	// prepare blog
	visibleBlog := entities.NewBlog(
		"title1",
		"content1",
		"description1",
		false,
		true,
	)
	newInBlog1 := entities.NewInBlog(
		*visibleBlog,
		[]int{1},
		[]int{1},
	)
	blogsRepo.Create(ctxTimeout, *newInBlog1)

	<-time.After(time.Second * 1)

	notVisibleBlog := entities.NewBlog(
		"title2",
		"content2",
		"description2",
		false,
		false,
	)
	newInBlog2 := entities.NewInBlog(
		*notVisibleBlog,
		[]int{1},
		[]int{1},
	)
	blogsRepo.Create(ctxTimeout, *newInBlog2)

	// list visible
	blogs, err := blogsRepo.AdminList(ctxTimeout)
	if err != nil {
		t.Fatalf("TestBlogsAdminListSqlite: list failed: %s", err)
	}
	if len(blogs) != 2 {
		t.Fatalf("TestBlogsAdminListSqlite: should return two")
	}

	// List will not return content (too large)
	err2 := compareListBlog(
		blogs,
		*visibleBlog,
		*notVisibleBlog,
		*topic1,
		*tag1,
	)
	if err2 != nil {
		t.Fatalf("TestBlogsAdminListSqlite: compare list blog failed: %s", err2)
	}
}

func TestBlogsAdminListByTopicIDsSqlite(t *testing.T) {
	// connect
	dbConn, err := sql.Open("sqlite3", "file:test.db?mode=memory&_foreign_keys=on")
	if err != nil {
		t.Fatalf("TestBlogsAdminListByTopicIDsSqlite: open db connection failed: %s", err)
	}
	defer dbConn.Close()

	// migrate db
	if err := db.Up(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite"); err != nil {
		t.Fatalf("TestBlogsAdminListByTopicIDsSqlite: migrate up failed: %s", err)
	}

	// setup repo
	blogsRepo, tagsRepo, topicsRepo := prepareRepos(dbConn)
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// lets just assume they work
	// prepare topic and tags
	topic1, _ := topicsRepo.Create(ctxTimeout, *entities.NewTopic("topic1", "topic1"))
	tag1, _ := tagsRepo.Create(ctxTimeout, *entities.NewTag("tag1", "tag1"))

	// prepare blog
	visibleBlog := entities.NewBlog(
		"title1",
		"content1",
		"description1",
		false,
		true,
	)
	newInBlog1 := entities.NewInBlog(
		*visibleBlog,
		[]int{1},
		[]int{1},
	)
	blogsRepo.Create(ctxTimeout, *newInBlog1)

	<-time.After(time.Second * 1)

	notVisibleBlog := entities.NewBlog(
		"title2",
		"content2",
		"description2",
		false,
		false,
	)
	newInBlog2 := entities.NewInBlog(
		*notVisibleBlog,
		[]int{1},
		[]int{1},
	)
	blogsRepo.Create(ctxTimeout, *newInBlog2)

	// list
	blogs, err := blogsRepo.AdminListByTopicIDs(ctxTimeout, []int{1})
	if err != nil {
		t.Fatalf("TestBlogsAdminListByTopicIDsSqlite: list failed: %s", err)
	}
	if len(blogs) != 2 {
		t.Fatalf("TestBlogsAdminListByTopicIDsSqlite: should return two")
	}

	// List will not return content (too large)
	err2 := compareListBlog(
		blogs,
		*visibleBlog,
		*notVisibleBlog,
		*topic1,
		*tag1,
	)
	if err2 != nil {
		t.Fatalf("TestBlogsAdminListSqlite: compare list blog failed: %s", err)
	}
}

func TestBlogsAdminListByTopicAndTagIDsSqlite(t *testing.T) {
	// connect
	dbConn, err := sql.Open("sqlite3", "file:test.db?mode=memory&_foreign_keys=on")
	if err != nil {
		t.Fatalf("TestBlogsAdminListByTopicAndTagIDsSqlite: open db connection failed: %s", err)
	}
	defer dbConn.Close()

	// migrate db
	if err := db.Up(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite"); err != nil {
		t.Fatalf("TestBlogsAdminListByTopicAndTagIDsSqlite: migrate up failed: %s", err)
	}

	// setup repo
	blogsRepo, tagsRepo, topicsRepo := prepareRepos(dbConn)
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// lets just assume they work
	// prepare topic and tags
	topic1, _ := topicsRepo.Create(ctxTimeout, *entities.NewTopic("topic1", "topic1"))
	tag1, _ := tagsRepo.Create(ctxTimeout, *entities.NewTag("tag1", "tag1"))

	// prepare blog
	visibleBlog := entities.NewBlog(
		"title1",
		"content1",
		"description1",
		false,
		true,
	)
	newInBlog1 := entities.NewInBlog(
		*visibleBlog,
		[]int{1},
		[]int{1},
	)
	blogsRepo.Create(ctxTimeout, *newInBlog1)

	<-time.After(time.Second * 1)

	notVisibleBlog := entities.NewBlog(
		"title2",
		"content2",
		"description2",
		false,
		false,
	)
	newInBlog2 := entities.NewInBlog(
		*notVisibleBlog,
		[]int{1},
		[]int{1},
	)
	blogsRepo.Create(ctxTimeout, *newInBlog2)

	// list
	blogs, err := blogsRepo.AdminListByTopicAndTagIDs(ctxTimeout, []int{1}, []int{1})
	if err != nil {
		t.Fatalf("TestBlogsAdminListByTopicAndTagIDsSqlite: list failed: %s", err)
	}
	if len(blogs) != 2 {
		t.Fatalf("TestBlogsAdminListByTopicAndTagIDsSqlite: should return two")
	}

	// List will not return content (too large)
	err2 := compareListBlog(
		blogs,
		*visibleBlog,
		*notVisibleBlog,
		*topic1,
		*tag1,
	)
	if err2 != nil {
		t.Fatalf("TestBlogsAdminListSqlite: compare list blog failed: %s", err)
	}
}

// This only compares slice with length of two
func compareListBlog(
	blogs []entities.OutBlog,
	visibleBlog,
	notVisibleBlog entities.Blog,
	topic1 entities.Topic,
	tag1 entities.Tag,
) error {
	// first blog
	if !cmp.Equal(notVisibleBlog, blogs[0].Blog, cmpopts.IgnoreFields(entities.Blog{}, "ID", "Created_at", "Updated_at", "Content")) {
		return fmt.Errorf("compareListBlog: list cmp blogs[0].Blog failed")
	}
	if !cmp.Equal(topic1, blogs[0].Topics[0]) {
		return fmt.Errorf("compareListBlog: list cmp blogs[0].Topics[0] failed")
	}
	if !cmp.Equal(tag1, blogs[0].Tags[0]) {
		return fmt.Errorf("compareListBlog: list cmp blogs[0].Tags[0] failed")
	}

	// second blog
	if !cmp.Equal(visibleBlog, blogs[1].Blog, cmpopts.IgnoreFields(entities.Blog{}, "ID", "Created_at", "Updated_at", "Content")) {
		return fmt.Errorf("compareListBlog: list cmp blogs[1].Blog failed")
	}
	if !cmp.Equal(topic1, blogs[1].Topics[0]) {
		return fmt.Errorf("compareListBlog: list cmp blogs[1].Topics[0] failed")
	}
	if !cmp.Equal(tag1, blogs[1].Tags[0]) {
		return fmt.Errorf("compareListBlog: list cmp blogs[1].Tags[0] failed")
	}

	return nil
}

func TestBlogsSoftDeleteSqlite(t *testing.T) {
	// connect
	dbConn, err := sql.Open("sqlite3", "file:test.db?mode=memory&_foreign_keys=on")
	if err != nil {
		t.Fatalf("TestBlogsSoftDeleteSqlite: open db connection failed: %s", err)
	}
	defer dbConn.Close()

	// migrate db
	if err := db.Up(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite"); err != nil {
		t.Fatalf("TestBlogsSoftDeleteSqlite: migrate up failed: %s", err)
	}

	// setup repo
	blogsRepo, tagsRepo, topicsRepo := prepareRepos(dbConn)
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// lets just assume they work
	// prepare topic and tags
	topicsRepo.Create(ctxTimeout, *entities.NewTopic("topic1", "topic1"))
	tagsRepo.Create(ctxTimeout, *entities.NewTag("tag1", "tag1"))

	// prepare blog
	visibleBlog := entities.NewBlog(
		"title1",
		"content1",
		"description1",
		false,
		true,
	)
	newInBlog1 := entities.NewInBlog(
		*visibleBlog,
		[]int{1},
		[]int{1},
	)
	blogsRepo.Create(ctxTimeout, *newInBlog1)

	affectedRows, err := blogsRepo.SoftDelete(ctxTimeout, 1)
	if err != nil {
		t.Fatalf("TestBlogsSoftDeleteSqlite: soft delete failed: %s", err)
	}
	if affectedRows != 1 {
		t.Fatalf("TestBlogsSoftDeleteSqlite: affectedRows should be 1")
	}

	// normal list and get shouldn't see it anymore
	_, err2 := blogsRepo.Get(ctxTimeout, 1)
	if err2 == nil {
		t.Fatalf("TestBlogsSoftDeleteSqlite: get shouldn't see this")
	}

	listResult, _ := blogsRepo.List(ctxTimeout)
	if len(listResult) != 0 {
		t.Fatalf("TestBlogsSoftDeleteSqlite: list shouldn't see this")
	}

	listByTopicResult, _ := blogsRepo.ListByTopicIDs(ctxTimeout, []int{1})
	if len(listByTopicResult) != 0 {
		t.Fatalf("TestBlogsSoftDeleteSqlite: list by topic shouldn't see this")
	}

	listByTopicAndTagIDsResult, _ := blogsRepo.ListByTopicAndTagIDs(ctxTimeout, []int{1}, []int{1})
	if len(listByTopicAndTagIDsResult) != 0 {
		t.Fatalf("TestBlogsSoftDeleteSqlite: list topic and tag ids shouldn't see this")
	}
}

func TestBlogsDeleteSqlite(t *testing.T) {
	// connect
	dbConn, err := sql.Open("sqlite3", "file:test.db?mode=memory&_foreign_keys=on")
	if err != nil {
		t.Fatalf("TestBlogsDeleteSqlite: open db connection failed: %s", err)
	}
	defer dbConn.Close()

	// migrate db
	if err := db.Up(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite"); err != nil {
		t.Fatalf("TestBlogsDeleteSqlite: migrate up failed: %s", err)
	}

	// setup repo
	blogsRepo, tagsRepo, topicsRepo := prepareRepos(dbConn)
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// lets just assume they work
	// prepare topic and tags
	topicsRepo.Create(ctxTimeout, *entities.NewTopic("topic1", "topic1"))
	tagsRepo.Create(ctxTimeout, *entities.NewTag("tag1", "tag1"))

	// prepare blog
	visibleBlog := entities.NewBlog(
		"title1",
		"content1",
		"description1",
		false,
		true,
	)
	newInBlog1 := entities.NewInBlog(
		*visibleBlog,
		[]int{1},
		[]int{1},
	)
	blogsRepo.Create(ctxTimeout, *newInBlog1)

	// delete (needs to be soft deleted first)
	affectedRows, err := blogsRepo.Delete(ctxTimeout, 1)
	if err != nil {
		t.Fatalf("TestBlogsDeleteSqlite: delete failed: %s", err)
	}
	if affectedRows == 1 {
		t.Fatalf("TestBlogsDeleteSqlite: affected rows should be zero")
	}

	// soft delete
	affectedRows2, err2 := blogsRepo.SoftDelete(ctxTimeout, 1)
	if err2 != nil {
		t.Fatalf("TestBlogsDeleteSqlite: soft delete failed: %s", err)
	}
	if affectedRows2 != 1 {
		t.Fatalf("TestBlogsDeleteSqlite: affectedRows should be 1")
	}

	// this time it will be deleted
	affectedRows3, err3 := blogsRepo.Delete(ctxTimeout, 1)
	if err3 != nil {
		t.Fatalf("TestBlogsDeleteSqlite: delete failed: %s", err)
	}
	if affectedRows3 != 1 {
		t.Fatalf("TestBlogsDeleteSqlite: affectedRows should be 1")
	}

	// try to get
	_, err5 := blogsRepo.AdminGet(ctxTimeout, 1)
	if err5 == nil {
		t.Fatalf("TestBlogsDeleteSqlite: should not be able to get this blog")
	}
}

func TestBlogsRestoreDeletedSqlite(t *testing.T) {
	// connect
	dbConn, err := sql.Open("sqlite3", "file:test.db?mode=memory&_foreign_keys=on")
	if err != nil {
		t.Fatalf("TestBlogsRestoreDeletedSqlite: open db connection failed: %s", err)
	}
	defer dbConn.Close()

	// migrate db
	if err := db.Up(dbConn, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite"); err != nil {
		t.Fatalf("TestBlogsRestoreDeletedSqlite: migrate up failed: %s", err)
	}

	// setup repo
	blogsRepo, tagsRepo, topicsRepo := prepareRepos(dbConn)
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// lets just assume they work
	// prepare topic and tags
	topicsRepo.Create(ctxTimeout, *entities.NewTopic("topic1", "topic1"))
	tagsRepo.Create(ctxTimeout, *entities.NewTag("tag1", "tag1"))

	// prepare blog
	visibleBlog := entities.NewBlog(
		"title1",
		"content1",
		"description1",
		false,
		true,
	)
	newInBlog1 := entities.NewInBlog(
		*visibleBlog,
		[]int{1},
		[]int{1},
	)
	blogsRepo.Create(ctxTimeout, *newInBlog1)

	// soft delete
	affectedRows, err := blogsRepo.SoftDelete(ctxTimeout, 1)
	if err != nil {
		t.Fatalf("TestBlogsRestoreDeletedSqlite: soft delete failed: %s", err)
	}
	if affectedRows != 1 {
		t.Fatalf("TestBlogsRestoreDeletedSqlite: affectedRows should be 1")
	}

	// restore soft delete
	blog, err := blogsRepo.RestoreDeleted(ctxTimeout, 1)
	if err != nil {
		t.Fatalf("TestBlogsRestoreDeletedSqlite: restore delete failed: %s", err)
	}
	if !cmp.Equal(visibleBlog, &blog.Blog, cmpopts.IgnoreFields(entities.Blog{}, "ID", "Created_at", "Updated_at", "Deleted_at")) {
		t.Fatalf("TestBlogsRestoreDeletedSqlite: restored blog cmp failed")
	}
}
