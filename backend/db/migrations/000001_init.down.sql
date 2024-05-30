DROP TABLE IF EXISTS blogs;

DROP TABLE IF EXISTS topics;

DROP TABLE IF EXISTS tags;

DROP TABLE IF EXISTS blog_tags;
DROP INDEX IF EXISTS blog_tags_blog;
DROP INDEX IF EXISTS blog_tags_tag;

DROP TABLE IF EXISTS blog_topics;
DROP INDEX IF EXISTS blog_topics_blog;
DROP INDEX IF EXISTS blog_topics_topic;
