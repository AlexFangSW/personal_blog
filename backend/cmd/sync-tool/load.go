package main

import (
	"blog/entities"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/gosimple/slug"
	"gopkg.in/yaml.v3"
)

type MetaFileContent struct {
	Topics []entities.InTopic `yaml:"topics"`
	Tags   []entities.InTag   `yaml:"tags"`
}

func loadMetaFile(metaFile string) (MetaFileContent, error) {
	slog.Debug("loadMetaFile")

	byteData, err := os.ReadFile(metaFile)
	if err != nil {
		return MetaFileContent{}, fmt.Errorf("loadMetaFile: read file failed: %w", err)
	}

	data := MetaFileContent{}
	if err := yaml.Unmarshal(byteData, &data); err != nil {
		return MetaFileContent{}, fmt.Errorf("loadMetaFile: yaml unmarshal failed: %w", err)
	}

	slog.Debug("load meta file", "content", data)
	return data, nil
}

type BlogFrontmatter struct {
	ID          int    // this will be loaded separately
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Pined       bool   `yaml:"pined"`
	Visible     bool   `yaml:"visible"`

	// will be transformed into slugs
	Tags   []string `yaml:"tags"`
	Topics []string `yaml:"topics"`
}

func (b *BlogFrontmatter) slugify() {
	// tags
	newTags := make([]string, 0, len(b.Tags))
	for _, tag := range b.Tags {
		newTags = append(newTags, slug.Make(tag))
	}
	b.Tags = newTags

	// topics
	newTopics := make([]string, 0, len(b.Tags))
	for _, topic := range b.Topics {
		newTopics = append(newTopics, slug.Make(topic))
	}
	b.Topics = newTopics
}

type BlogInfo struct {
	Frontmatter BlogFrontmatter
	Content_md5 string
	Filename    string
}

func NewBlogInfo(frontmatter BlogFrontmatter, content string, filename string) BlogInfo {
	content_md5 := fmt.Sprintf("%x", md5.Sum([]byte(content)))
	frontmatter.slugify()
	return BlogInfo{
		Frontmatter: frontmatter,
		Content_md5: content_md5,
		Filename:    filename,
	}
}

// load all blogs in blogDir, parse their mdx frontmatter
func loadBlogs(blogDir string, idMap map[string]int) ([]BlogInfo, error) {
	slog.Debug("loadBlogs")

	// get all the files
	files, err := os.ReadDir(blogDir)
	if err != nil {
		return []BlogInfo{}, fmt.Errorf("loadBlogs: read dir failed: %w", err)
	}
	slog.Debug("got blogs", "blog count", len(files))

	result := []BlogInfo{}

	// load all of them
	// split by '---' and use  yaml to unmartial the data
	for _, file := range files {
		filepath := fmt.Sprintf("%s/%s", blogDir, file.Name())
		rawData, err := os.ReadFile(filepath)
		if err != nil {
			return []BlogInfo{}, fmt.Errorf("loadBlogs: read file %q failed: %w", filepath, err)
		}

		full := string(rawData)
		splited := strings.Split(full, "---")
		header := splited[1]
		content := splited[2]

		parsedHeader := BlogFrontmatter{}
		if err := yaml.Unmarshal([]byte(header), &parsedHeader); err != nil {
			return []BlogInfo{}, fmt.Errorf("loadBlogs: parse header from %q failed: %w", filepath, err)
		}
		id, ok := idMap[file.Name()]
		if ok {
			slog.Debug("got id for blog", "filename", file.Name(), "id", id)
			parsedHeader.ID = id
		} else {
			slog.Debug("new blog", "filename", file.Name())
		}

		blogsInfo := NewBlogInfo(parsedHeader, content, file.Name())
		result = append(result, blogsInfo)
	}

	slog.Debug("local blogs loaded", "blog count", len(result))
	return result, nil
}

// load ids.json
func loadIDMap(idFile string) (map[string]int, error) {
	idMap := map[string]int{}
	file, err := os.ReadFile(idFile)
	if err != nil {
		slog.Warn("load ids.json failed, might be first sync ?", "err", err)
	} else {
		if err := json.Unmarshal(file, &idMap); err != nil {
			return map[string]int{}, fmt.Errorf("loadBlogs: load unmarshal ids.json failed: %w", err)
		}
	}
	slog.Debug("ids.json", "content", idMap)
	return idMap, nil
}
