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
