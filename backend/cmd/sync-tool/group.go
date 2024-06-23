// In this project, the local copy of blogs, tags and topics are the source of truth.
package main

import (
	"blog/entities"
	"errors"
	"log/slog"
	"slices"
)

var (
	BlogTagReferenceError   = errors.New("referece error for blog.tags")
	BlogTopicReferenceError = errors.New("referece error for blog.topics")
)

type GroupTypes interface {
	entities.Tag | entities.Topic
}

type BlogGroupTypes interface {
	BlogInfo | entities.InBlog
}

type Groups[T GroupTypes] struct {
	create []T // exists locally but not on server
	update []T // exists locally and on server, but the name or description is different
	delete []T // doesn't exist in local, but exists on server
	noop   []T // exactly the same
}

type BlogGroup[T BlogGroupTypes] struct {
	create []T                      // exists locally but not on server
	update []T                      // exists locally and on server, but the name or description is different
	delete []entities.OutBlogSimple // doesn't exist in local, but exists on server
	noop   []T                      // exactly the same
}

func groupTags(localTags []entities.InTag, tags []entities.Tag) (Groups[entities.Tag], error) {
	slog.Debug("groupTags")

	// for faster lookup
	tagMap := map[string]entities.Tag{}
	for _, tag := range tags {
		tagMap[tag.Name] = tag
	}

	// go throgh local records, match them to remote data and categorize them
	result := Groups[entities.Tag]{}
	for _, localTag := range localTags {
		remoteTag, ok := tagMap[localTag.Name]

		// the record dosen't exist on remote, we should create it
		if !ok {
			newTag := entities.NewTag(localTag.Name, localTag.Description)
			result.create = append(result.create, *newTag)
			continue
		}

		// the record is identical, do nothing
		if localTag.Description == remoteTag.Description {
			result.noop = append(result.noop, remoteTag)
			delete(tagMap, localTag.Name)
			continue

		} else {
			// the content is differrnt, we should update it
			newTag := entities.NewTagWithID(remoteTag.ID, localTag.Name, localTag.Description)
			result.update = append(result.noop, *newTag)
			delete(tagMap, localTag.Name)
			continue
		}
	}

	// the remaining remote data should be deleted
	for _, remoteTag := range tagMap {
		result.delete = append(result.delete, remoteTag)
	}

	slog.Debug(
		"grouped tags",
		"create", len(result.create),
		"update", len(result.update),
		"delete", len(result.delete),
		"noop", len(result.noop),
	)
	return result, nil
}

func groupTopics(localTopics []entities.InTopic, topics []entities.Topic) (Groups[entities.Topic], error) {
	slog.Debug("groupTopics")

	topicMap := map[string]entities.Topic{}
	for _, topic := range topics {
		topicMap[topic.Name] = topic
	}

	result := Groups[entities.Topic]{}
	for _, localTopic := range localTopics {
		remoteTopic, ok := topicMap[localTopic.Name]

		// the record dosen't exist on remote, we should create it
		if !ok {
			newTopic := entities.NewTopic(localTopic.Name, localTopic.Description)
			result.create = append(result.create, *newTopic)
			continue
		}

		// the record is identical, do nothing
		if localTopic.Description == remoteTopic.Description {
			result.noop = append(result.noop, remoteTopic)
			delete(topicMap, localTopic.Name)
			continue

		} else {
			// the content is differrnt, we should update it
			newTopic := entities.NewTopicWithID(remoteTopic.ID, localTopic.Name, localTopic.Description)
			result.update = append(result.update, *newTopic)
			delete(topicMap, localTopic.Name)
			continue
		}
	}

	// the remaining remote data should be deleted
	for _, remoteTopic := range topicMap {
		result.delete = append(result.delete, remoteTopic)
	}

	slog.Debug(
		"grouped topics",
		"create", len(result.create),
		"update", len(result.update),
		"delete", len(result.delete),
		"noop", len(result.noop),
	)
	return result, nil
}

func groupBlogs(localBlogs []BlogInfo, blogs []entities.OutBlogSimple) (BlogGroup[BlogInfo], error) {
	slog.Debug("groupBlogs")

	blogMap := map[int]entities.OutBlogSimple{}
	for _, blog := range blogs {
		blogMap[blog.ID] = blog
	}

	result := BlogGroup[BlogInfo]{}
	for _, localBlog := range localBlogs {
		remoteBlog, ok := blogMap[localBlog.Frontmatter.ID]

		// the record dosen't exist on remote, we should create it.
		// whether the blog should use its current ID will be decided at a later stage.
		if !ok {
			result.create = append(result.create, localBlog)
			continue
		}

		// the record is identical, do nothing
		if blogEqual(localBlog, remoteBlog) {
			result.noop = append(result.noop, localBlog)
			delete(blogMap, localBlog.Frontmatter.ID)
			continue

		} else {
			// the content is differrnt, we should update it
			result.update = append(result.update, localBlog)
			delete(blogMap, localBlog.Frontmatter.ID)
			continue
		}
	}

	// the remaining remote data should be deleted
	for _, remoteBlog := range blogMap {
		result.delete = append(result.delete, remoteBlog)
	}

	slog.Debug(
		"grouped blogs",
		"create", len(result.create),
		"update", len(result.update),
		"delete", len(result.delete),
		"noop", len(result.noop),
	)
	return result, nil
}

func blogEqual(localBlog BlogInfo, remoteBlog entities.OutBlogSimple) bool {
	if localBlog.Frontmatter.Title != remoteBlog.Title {
		slog.Debug("Title not equal", "filename", localBlog.Filename)
		return false
	}
	if localBlog.Frontmatter.Description != remoteBlog.Description {
		slog.Debug("Description not equal", "filename", localBlog.Filename)
		return false
	}
	if localBlog.Frontmatter.Pined != remoteBlog.Pined {
		slog.Debug("Pined not equal", "filename", localBlog.Filename)
		return false
	}
	if localBlog.Frontmatter.Visible != remoteBlog.Visible {
		slog.Debug("Visible not equal", "filename", localBlog.Filename)
		return false
	}
	if !slices.Equal[[]string](localBlog.Frontmatter.Tags, remoteBlog.Tags) {
		slog.Debug("Tags not equal", "filename", localBlog.Filename)
		return false
	}
	if !slices.Equal[[]string](localBlog.Frontmatter.Topics, remoteBlog.Topics) {
		slog.Debug("Topics not equal", "filename", localBlog.Filename)
		return false
	}
	if localBlog.Content_md5 != remoteBlog.ContentMD5 {
		slog.Debug("Content_md5 not equal", "filename", localBlog.Filename)
		return false
	}
	return true
}

// map topics and tag names to their id
// if there is no match, record it and return an error
func transformBlogs(tags []entities.Tag, topics []entities.Topic, blogs BlogGroup[BlogInfo]) (BlogGroup[entities.InBlog], error) {
	slog.Debug("transformBlogs")

	// prepare for lookup
	topicMap := map[string]entities.Topic{}
	for _, topic := range topics {
		topicMap[topic.Slug] = topic
	}
	tagMap := map[string]entities.Tag{}
	for _, tag := range tags {
		tagMap[tag.Slug] = tag
	}

	// go though all the blogs, and transform them

	return BlogGroup[entities.InBlog]{}, nil
}
