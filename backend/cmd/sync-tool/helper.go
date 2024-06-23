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
	baseURL   string
	token     string // jwt token
	batchSize int
}

func NewSyncHelper(baseURL, token string, batchSize int) SyncHelper {
	return SyncHelper{
		baseURL:   baseURL,
		token:     token,
		batchSize: batchSize,
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

// format to store blog transform errors
type TransformError struct {
	Ref                BlogInfo `json:"ref"`
	NoneMatchingTags   []string `json:"noneMatchingTags"`   // slug
	NoneMatchingTopics []string `json:"noneMatchingTopics"` // slug
}

func NewTransformError(ref BlogInfo, tags, topics []string) TransformError {
	return TransformError{
		Ref:                ref,
		NoneMatchingTags:   tags,
		NoneMatchingTopics: topics,
	}
}

// this struct should only be used once
type BlogTransformHelper struct {
	tagMap            map[string]int
	topicMap          map[string]int
	accumulatedErrors []TransformError
	errorFilePath     string
}

func NewBlogTransformHelper(tags []entities.Tag, topics []entities.Topic, sourcePath string) BlogTransformHelper {
	// prepare for lookup
	topicMap := map[string]int{}
	for _, topic := range topics {
		topicMap[topic.Slug] = topic.ID
	}
	tagMap := map[string]int{}
	for _, tag := range tags {
		tagMap[tag.Slug] = tag.ID
	}

	return BlogTransformHelper{
		tagMap:        tagMap,
		topicMap:      topicMap,
		errorFilePath: path.Join(sourcePath, "blog-transform-error.json"),
	}
}

// transform the blog the a format that the server api accepts.
func (b *BlogTransformHelper) Transform(blogs BlogGroup[BlogInfo]) (BlogGroup[entities.InBlog], error) {
	slog.Debug("Transform")

	result := BlogGroup[entities.InBlog]{}

	defer func() {
		slog.Debug(
			"blog transform result",
			"transformed blogs", len(result.create)+len(result.update),
			"errors", len(b.accumulatedErrors),
		)
	}()

	result.create = b.transform(blogs.create)
	result.update = b.transform(blogs.update)
	result.delete = blogs.delete

	if len(b.accumulatedErrors) > 0 {
		if err := b.saveError(); err != nil {
			return BlogGroup[entities.InBlog]{}, fmt.Errorf("Transform: save error failed: %w", err)
		}
		return BlogGroup[entities.InBlog]{}, BlogReferenceError
	}

	return result, nil
}

// save the error record to a file
func (b BlogTransformHelper) saveError() error {
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
func (b *BlogTransformHelper) transform(blogs []BlogInfo) []entities.InBlog {
	slog.Debug("transform")

	result := []entities.InBlog{}
	for _, blog := range blogs {
		slog.Debug("transforming blog", "blog", blog.Filename)
		currErr := NewTransformError(blog, []string{}, []string{})

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

		// we don't need to do a full transformation after we hit an error,
		// but we will still loop through all the blogs to get a complete error report.
		if len(currErr.NoneMatchingTags) > 0 ||
			len(currErr.NoneMatchingTopics) > 0 {
			b.accumulatedErrors = append(b.accumulatedErrors, currErr)
			continue
		} else if len(b.accumulatedErrors) > 0 {
			continue
		}

		// we still won't load the entire content
		newBlog := entities.NewBlog(
			blog.Frontmatter.Title,
			"",
			blog.Frontmatter.Description,
			blog.Frontmatter.Pined,
			blog.Frontmatter.Visible,
		)
		newInBlog := entities.NewInBlog(
			*newBlog,
			tagIDs,
			topicIDs,
		)
		result = append(result, *newInBlog)
	}
	return result
}
