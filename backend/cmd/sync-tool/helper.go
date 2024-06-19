package main

import (
	"blog/entities"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

var LimitReaderSize int64 = 10 * 1024 * 1024 // 10MB

func drainAndClose(body io.ReadCloser) error {
	reader := io.LimitReader(body, LimitReaderSize)
	_, drainErr := io.Copy(io.Discard, reader)
	if drainErr != nil {
		return fmt.Errorf("drainAndClose: drain failed: %w", drainErr)
	}
	return body.Close()
}

type SyncHelper struct {
	baseURL string
	token   string
}

func NewSyncHelper(baseURL, token string) SyncHelper {
	return SyncHelper{
		baseURL: baseURL,
		token:   token,
	}
}

// ================ tags  ================
func (s SyncHelper) GetAllTags() (oTag []entities.Tag, oErr error) {
	res, err := httpClient.Get(fmt.Sprintf("%s/tags", s.baseURL))
	if err != nil {
		return []entities.Tag{}, fmt.Errorf("GetAllTags: get all tags failed: %w", err)
	}

	// cleanup
	defer func() {
		oErr = errors.Join(oErr, drainAndClose(res.Body))
	}()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return []entities.Tag{}, fmt.Errorf("GetAllTags: read body failed: %w", err)
	}

	if res.StatusCode >= 400 {
		return []entities.Tag{}, errors.New(string(resBody))
	}

	data := entities.RetSuccess[[]entities.Tag]{}
	if err := json.Unmarshal(resBody, &data); err != nil {
		return []entities.Tag{}, fmt.Errorf("GetAllTags: decode body failed: %w", err)
	}

	return data.Msg, nil
}
func (s SyncHelper) CreateTags(tags []entities.Tag) error {
	return nil
}
func (s SyncHelper) UpdateTags(tags []entities.Tag) error {
	return nil
}
func (s SyncHelper) DeleteTags(tags []entities.Tag) error {
	return nil
}

// ================ topics ================
func (s SyncHelper) GetAllTopics() (oTag []entities.Topic, oErr error) {
	res, err := httpClient.Get(fmt.Sprintf("%s/topics", s.baseURL))
	if err != nil {
		return []entities.Topic{}, fmt.Errorf("GetAllTopics: get all tags failed: %w", err)
	}

	// cleanup
	defer func() {
		oErr = errors.Join(oErr, drainAndClose(res.Body))
	}()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return []entities.Topic{}, fmt.Errorf("GetAllTopics: read body failed: %w", err)
	}

	if res.StatusCode >= 400 {
		return []entities.Topic{}, errors.New(string(resBody))
	}

	data := entities.RetSuccess[[]entities.Topic]{}
	if err := json.Unmarshal(resBody, &data); err != nil {
		return []entities.Topic{}, fmt.Errorf("GetAllTopics: decode body failed: %w", err)
	}

	return data.Msg, nil
}
func (s SyncHelper) CreateTopics(topics []entities.Topic) error {
	return nil
}
func (s SyncHelper) UpdateTopics(topics []entities.Topic) error {
	return nil
}
func (s SyncHelper) DeleteTopics(topics []entities.Topic) error {
	return nil
}

// ================ blogs ================
func (s SyncHelper) GetAllBlogs() ([]entities.OutBlogSimple, error) {
	return []entities.OutBlogSimple{}, nil
}
func (s SyncHelper) CreateBlogs(blogs []entities.InBlog) error {
	// load blog content by batch
	return nil
}
func (s SyncHelper) UpdateBlogs(blogs []entities.InBlog) error {
	// load blog content by batch
	return nil
}
func (s SyncHelper) DeleteBlogs(blogs []entities.InBlog) error {
	return nil
}

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
	// log info
	return Groups[entities.Tag]{}, nil
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

type MetaFileContent struct {
	topics []entities.InTopic
	tags   []entities.InTag
}

func loadMetaFile(metaFile string) (MetaFileContent, error) {
	byteData, err := os.ReadFile(metaFile)
	if err != nil {
		return MetaFileContent{}, fmt.Errorf("loadMetaFile: read file failed: %w", err)
	}

	data := MetaFileContent{}
	if err := yaml.Unmarshal(byteData, &data); err != nil {
		return MetaFileContent{}, fmt.Errorf("loadMetaFile: yaml unmarshal failed: %w", err)
	}

	return data, nil
}

type BlogFrontmatter struct {
	title       string
	description string
	pined       bool
	visible     bool
	tags        []string // slug
	topics      []string // slug
}

type BlogInfo struct {
	frontmatter BlogFrontmatter
	content_md5 string
}

func NewBlogInfo(frontmatter BlogFrontmatter, content string) BlogInfo {
	content_md5 := fmt.Sprintf("%x", md5.Sum([]byte(content)))
	return BlogInfo{
		frontmatter: frontmatter,
		content_md5: content_md5,
	}
}

// load all blogs in blogDir, parse their mdx frontmatter
func loadBlogs(blogDir string) ([]BlogInfo, error) {
	// get all the files

	// load all of them
	// split by '---' and use  yaml to unmartial the data

	return []BlogInfo{}, nil
}
