-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS blogs(
  id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,

  -- ISO 8061
  created_at TEXT NOT NULL DEFAULT (strftime('%FT%T+00:00')), 
  updated_at TEXT NOT NULL DEFAULT (strftime('%FT%T+00:00')), 
  deleted_at TEXT DEFAULT "", 

  title TEXT NOT NULL UNIQUE,
  content TEXT DEFAULT "",
  description TEXT DEFAULT "",
  slug TEXT NOT NULL UNIQUE,
  pined BOOLEAN DEFAULT 0,
  visible BOOLEAN DEFAULT 0
);

-- Auto update 'update_ts' on updated row
-- Currently, sqlite trigger only supports FOR EACH ROW (fire on each effected row)
CREATE TRIGGER IF NOT EXISTS blogs_update_ts
AFTER UPDATE ON blogs
BEGIN 
  UPDATE blogs SET updated_at = (strftime('%FT%T+00:00')) WHERE id = NEW.id;
END;

CREATE TABLE IF NOT EXISTS topics(
  id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,

  -- ISO 8061
  created_at TEXT NOT NULL DEFAULT (strftime('%FT%T+00:00')), 
  updated_at TEXT NOT NULL DEFAULT (strftime('%FT%T+00:00')), 

  name TEXT NOT NULL UNIQUE,
  description TEXT DEFAULT "",
  slug TEXT NOT NULL UNIQUE
);

CREATE TRIGGER IF NOT EXISTS topics_update_ts
AFTER UPDATE ON topics
BEGIN 
  UPDATE topics SET updated_at = (strftime('%FT%T+00:00')) WHERE id = NEW.id;
END;

CREATE TABLE IF NOT EXISTS tags(
  id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,

  -- ISO 8061
  created_at TEXT NOT NULL DEFAULT (strftime('%FT%T+00:00')), 
  updated_at TEXT NOT NULL DEFAULT (strftime('%FT%T+00:00')), 

  name TEXT NOT NULL UNIQUE,
  description TEXT DEFAULT "",
  slug TEXT NOT NULL UNIQUE
);

CREATE TRIGGER IF NOT EXISTS tags_update_ts
AFTER UPDATE ON tags
BEGIN 
  UPDATE tags SET updated_at = (strftime('%FT%T+00:00')) WHERE id = NEW.id;
END;

CREATE TABLE IF NOT EXISTS blog_tags(
  blog_id INTEGER NOT NULL,
  tag_id INTEGER NOT NULL,
  FOREIGN KEY(blog_id) REFERENCES blogs(id),
  FOREIGN KEY(tag_id) REFERENCES tags(id),
  PRIMARY KEY(blog_id, tag_id)
);
CREATE INDEX IF NOT EXISTS blog_tags_blog ON blog_tags (blog_id);
CREATE INDEX IF NOT EXISTS blog_tags_tag ON blog_tags (tag_id);

CREATE TABLE IF NOT EXISTS blog_topics(
  blog_id INTEGER NOT NULL,
  topic_id INTEGER NOT NULL,
  FOREIGN KEY(blog_id) REFERENCES blogs(id),
  FOREIGN KEY(topic_id) REFERENCES topics(id),
  PRIMARY KEY(blog_id, topic_id)
);
CREATE INDEX IF NOT EXISTS blog_topics_blog ON blog_topics (blog_id);
CREATE INDEX IF NOT EXISTS blog_topics_topic ON blog_topics (topic_id);

-- only one user 
CREATE TABLE IF NOT EXISTS users(
  id INTEGER NOT NULL UNIQUE PRIMARY KEY CHECK (id = 0),
  name TEXT NOT NULL UNIQUE,
  -- encoded password
  pwd TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS blogs;
DROP TRIGGER IF EXISTS blogs_update_ts;

DROP TABLE IF EXISTS topics;
DROP TRIGGER IF EXISTS topics_update_ts; 

DROP TABLE IF EXISTS tags;
DROP TRIGGER IF EXISTS tags_update_ts;

DROP TABLE IF EXISTS blog_tags;
DROP INDEX IF EXISTS blog_tags_blog;
DROP INDEX IF EXISTS blog_tags_tag;

DROP TABLE IF EXISTS blog_topics;
DROP INDEX IF EXISTS blog_topics_blog;
DROP INDEX IF EXISTS blog_topics_topic;

DROP TABLE IF EXISTS users;
-- +goose StatementEnd
