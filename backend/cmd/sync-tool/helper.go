package main

import (
	"blog/entities"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
)

var (
	BlogReferenceError       = errors.New("blog reference error")
	LimitReaderSize    int64 = 10 * 1024 * 1024 // 10MB
)

func drainAndClose(body io.ReadCloser) error {
	reader := io.LimitReader(body, LimitReaderSize)
	_, drainErr := io.Copy(io.Discard, reader)
	if drainErr != nil {
		return fmt.Errorf("drainAndClose: drain failed: %w", drainErr)
	}
	return body.Close()
}

type SyncHelper struct {
	baseURL    string
	token      string // jwt token
	batchSize  int
	sourcePath string
}

func NewSyncHelper(baseURL, token string, batchSize int, sourcePath string) SyncHelper {
	return SyncHelper{
		baseURL:    baseURL,
		token:      token,
		batchSize:  batchSize,
		sourcePath: sourcePath,
	}
}

func batch[T any](inpt []T, size int, out chan<- []T) {
	inptLength := len(inpt)
	for i := 0; i < inptLength; i += size {
		upperBound := min(i+size, inptLength)
		out <- inpt[i:upperBound:upperBound]
	}
	close(out)
	return
}

type MaperError struct {
	Ref                BlogInfo `json:"ref"`
	NoneMatchingTags   []string `json:"noneMatchingTags"`   // slug
	NoneMatchingTopics []string `json:"noneMatchingTopics"` // slug
}

func NewMaperError(ref BlogInfo, tags, topics []string) MaperError {
	return MaperError{
		Ref:                ref,
		NoneMatchingTags:   tags,
		NoneMatchingTopics: topics,
	}
}

// this struct should only be used once
// used for mapping tag and topic slugs to their ids
type BlogMaper struct {
	tagMap            map[string]int
	topicMap          map[string]int
	accumulatedErrors []MaperError
	errorFilePath     string
}

func NewBlogMaper(tags []entities.Tag, topics []entities.Topic, sourcePath string) BlogMaper {
	// prepare for lookup
	topicMap := map[string]int{}
	for _, topic := range topics {
		topicMap[topic.Slug] = topic.ID
	}
	tagMap := map[string]int{}
	for _, tag := range tags {
		tagMap[tag.Slug] = tag.ID
	}

	return BlogMaper{
		tagMap:        tagMap,
		topicMap:      topicMap,
		errorFilePath: path.Join(sourcePath, "blog-map-error.json"),
	}
}

// map topic and tag slugs to their ids
func (b *BlogMaper) MapIDs(blogs BlogGroup[BlogInfo]) (BlogGroup[BlogInfo], error) {
	slog.Info("MapIDs")

	result := BlogGroup[BlogInfo]{}

	defer func() {
		slog.Info(
			"map id result",
			"processed blogs", len(result.create)+len(result.update),
			"errors", len(b.accumulatedErrors),
		)
	}()

	result.create = b.mapIds(blogs.create)
	result.update = b.mapIds(blogs.update)
	result.delete = blogs.delete

	if len(b.accumulatedErrors) > 0 {
		if err := b.saveError(); err != nil {
			return BlogGroup[BlogInfo]{}, fmt.Errorf("MapIDs: save error failed: %w", err)
		}
		return BlogGroup[BlogInfo]{}, BlogReferenceError
	}

	return result, nil
}

// save the error record to a file
func (b BlogMaper) saveError() error {
	data, err := json.Marshal(b.accumulatedErrors)
	if err != nil {
		return fmt.Errorf("saveError: encode accumulated error failed: %w", err)
	}
	slog.Info("saving blog transform error to file", "file", b.errorFilePath)
	if err := os.WriteFile(b.errorFilePath, data, 0644); err != nil {
		return fmt.Errorf("saveError: write file failed: %w", err)
	}
	return nil
}

// map topics and tag names to their id
// if there is no match, record it and return an error
func (b *BlogMaper) mapIds(blogs []BlogInfo) []BlogInfo {
	slog.Debug("mapIds")

	result := []BlogInfo{}
	for _, blog := range blogs {
		slog.Debug("map ids", "blog", blog.Filename)
		currErr := NewMaperError(blog, []string{}, []string{})

		// check if the blog contains any topic or tags that doesn't exist
		tagIDs := []int{}
		for _, tag := range blog.Frontmatter.Tags {
			id, ok := b.tagMap[tag]
			if !ok {
				slog.Error("blog refereced a none existent tag", "tag", tag, "filename", blog.Filename)
				currErr.NoneMatchingTags = append(currErr.NoneMatchingTags, tag)
			}
			tagIDs = append(tagIDs, id)
		}

		topicIDs := []int{}
		for _, topic := range blog.Frontmatter.Topics {
			id, ok := b.topicMap[topic]
			if !ok {
				slog.Error("blog refereced a none existent topic", "topic", topic, "filename", blog.Filename)
				currErr.NoneMatchingTopics = append(currErr.NoneMatchingTopics, topic)
			}
			topicIDs = append(topicIDs, id)
		}

		// we don't need to finish this after we hit an error,
		// but we will still loop through all the blogs to get a complete error report.
		if len(currErr.NoneMatchingTags) > 0 ||
			len(currErr.NoneMatchingTopics) > 0 {
			b.accumulatedErrors = append(b.accumulatedErrors, currErr)
			continue
		} else if len(b.accumulatedErrors) > 0 {
			continue
		}

		blog.Frontmatter.TagIDs = tagIDs
		blog.Frontmatter.TopicIDs = topicIDs

		result = append(result, blog)
	}
	return result
}
