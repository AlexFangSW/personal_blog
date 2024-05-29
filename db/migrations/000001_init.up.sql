CREATE TABLE IF NOT EXISTS blogs(
  id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,

  -- ISO 8061
  created_at TEXT NOT NULL DEFAULT (strftime('%FT%T+00:00')), 
  updated_at TEXT NOT NULL DEFAULT (strftime('%FT%T+00:00')), 
  deleted_at TEXT DEFAULT NULL, 

  title TEXT NOT NULL UNIQUE,
  content TEXT DEFAULT "",
  description TEXT DEFAULT "",
  slug TEXT NOT NULL UNIQUE,
  pined BOOLEAN DEFAULT 0,
  visible BOOLEAN DEFAULT 0
);
-- SQlite automatically creates an index for UNIQUE columns
-- CREATE UNIQUE INDEX IF NOT EXISTS blog_slug ON blogs (slug);

CREATE TABLE IF NOT EXISTS topics(
  id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,

  -- ISO 8061
  created_at TEXT NOT NULL DEFAULT (strftime('%FT%T+00:00')), 
  updated_at TEXT NOT NULL DEFAULT (strftime('%FT%T+00:00')), 

  name TEXT NOT NULL UNIQUE,
  description TEXT DEFAULT "",
  slug TEXT NOT NULL UNIQUE
);
-- CREATE UNIQUE INDEX IF NOT EXISTS topic_slug ON topics (slug);

CREATE TABLE IF NOT EXISTS tags(
  id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,

  -- ISO 8061
  created_at TEXT NOT NULL DEFAULT (strftime('%FT%T+00:00')), 
  updated_at TEXT NOT NULL DEFAULT (strftime('%FT%T+00:00')), 

  name TEXT NOT NULL UNIQUE,
  description TEXT DEFAULT "",
  slug TEXT NOT NULL UNIQUE
);
-- CREATE UNIQUE INDEX IF NOT EXISTS tag_slug ON tags (slug);

CREATE TABLE IF NOT EXISTS blog_tags(
  blog_id INTEGER NOT NULL,
  tag_id INTEGER NOT NULL,
  FOREIGN KEY(blog_id) REFERENCES blogs(id),
  FOREIGN KEY(tag_id) REFERENCES tags(id)
);
CREATE INDEX IF NOT EXISTS blog_tags_blog ON blog_tags (blog_id);
CREATE INDEX IF NOT EXISTS blog_tags_tag ON blog_tags (tag_id);

CREATE TABLE IF NOT EXISTS blog_topics(
  blog_id INTEGER NOT NULL,
  topic_id INTEGER NOT NULL,
  FOREIGN KEY(blog_id) REFERENCES blogs(id),
  FOREIGN KEY(topic_id) REFERENCES topics(id)
);
CREATE INDEX IF NOT EXISTS blog_topics_blog ON blog_topics (blog_id);
CREATE INDEX IF NOT EXISTS blog_topics_topic ON blog_topics (topic_id);

-- Remamber to set foreign_keys ON to enforce foreign key constraint
-- SQlite disables it by default... 
-- PRAGMA foreign_keys = ON;
