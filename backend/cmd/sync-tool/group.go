// In this project, the local copy of blogs, tags and topics are the source of truth.
// The server, frontend and database are all just a way to organize the data and show our content.
// The only thing we need to persist is a folder (such as the 'dummyData' folder),
// which stores markdown and yaml files that can easily be persisted by pushing to GitHub
package main

import (
	"blog/entities"
	"log/slog"
)

type GroupTypes interface {
	entities.Tag | entities.Topic | BlogInfo | entities.InBlog
}

type Groups[T GroupTypes] struct {
	create []T // exists locally but not on server
	update []T // exists locally and on server, but the name or description is different
	delete []T // doesn't exist in local, but exists on server
	noop   []T // exactly the same
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

		if !ok {
			newTag := entities.NewTag(localTag.Name, localTag.Description)
			result.create = append(result.create, *newTag)
			continue
		}

		if localTag.Description == remoteTag.Description {
			result.noop = append(result.noop, remoteTag)
			delete(tagMap, localTag.Name)
			continue

		} else {
			newTag := entities.NewTag(localTag.Name, localTag.Description)
			newTag.ID = remoteTag.ID
			result.update = append(result.noop, *newTag)
			delete(tagMap, localTag.Name)
			continue
		}
	}

	// the remaining remote data should be deleted
	for _, remoteTag := range tagMap {
		result.delete = append(result.delete, remoteTag)
	}

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

	}

	return Groups[entities.Topic]{}, nil
}

// this InBlog will have their id filled
func groupBlogs(localBlogs []BlogInfo, blogs []entities.OutBlogSimple) (Groups[BlogInfo], error) {
	return Groups[BlogInfo]{}, nil
}

// map topics and tag slugs to their id
// if there is no match, simply remove it
func transformBlogs(tags Groups[entities.Tag], topics Groups[entities.Topic], blogs Groups[BlogInfo]) (Groups[entities.InBlog], error) {
	return Groups[entities.InBlog]{}, nil
}
