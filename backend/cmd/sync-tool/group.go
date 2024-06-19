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

	// log info
	tagMap := map[string]entities.Tag{}
	for _, tag := range tags {
		tagMap[tag.Name] = tag
	}

	result := Groups[entities.Tag]{}
	for _, localTag := range localTags {
		remoteTag, ok := tagMap[localTag.Name]

		if !ok {
			newTag := entities.NewTag(localTag.Name, localTag.Description)
			newTag.ID = remoteTag.ID
			result.create = append(result.create, *newTag)
			delete(tagMap, localTag.Name)
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

	for _, remoteTag := range tagMap {
		result.delete = append(result.delete, remoteTag)
	}

	return result, nil
}

func groupTopics(localTopics []entities.InTopic, topics []entities.Topic) (Groups[entities.Topic], error) {
	// log info
	return Groups[entities.Topic]{}, nil
}

// this InBlog will have their id filled
func groupBlogs(localBlogs []BlogInfo, blogs []entities.OutBlogSimple) (Groups[BlogInfo], error) {
	// log info
	return Groups[BlogInfo]{}, nil
}

// map topics and tag slugs to their id
// if there is no match, simply remove it
func transformBlogs(tags Groups[entities.Tag], topics Groups[entities.Topic], blogs Groups[BlogInfo]) (Groups[entities.InBlog], error) {
	return Groups[entities.InBlog]{}, nil
}
