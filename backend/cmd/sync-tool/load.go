package main

import (
	"blog/entities"
	"crypto/md5"
	"fmt"
	"log/slog"
	"os"
	"strings"

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

	return data, nil
}

type BlogFrontmatter struct {
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Pined       bool     `yaml:"pined"`
	Visible     bool     `yaml:"visible"`
	Tags        []string `yaml:"tags"`   // slug
	Topics      []string `yaml:"topics"` // slug
}

type BlogInfo struct {
	Frontmatter BlogFrontmatter
	Content_md5 string
}

func NewBlogInfo(frontmatter BlogFrontmatter, content string) BlogInfo {
	content_md5 := fmt.Sprintf("%x", md5.Sum([]byte(content)))
	return BlogInfo{
		Frontmatter: frontmatter,
		Content_md5: content_md5,
	}
}

// load all blogs in blogDir, parse their mdx frontmatter
func loadBlogs(blogDir string) ([]BlogInfo, error) {
	slog.Debug("loadBlogs")

	// get all the files
	files, err := os.ReadDir(blogDir)
	if err != nil {
		return []BlogInfo{}, fmt.Errorf("loadBlogs: read dir failed: %w", err)
	}
	slog.Debug("got blogs", "blog count", len(files))

	// load all of them
	// split by '---' and use  yaml to unmartial the data
	result := []BlogInfo{}

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

		blogsInfo := NewBlogInfo(parsedHeader, content)
		result = append(result, blogsInfo)
	}

	return result, nil
}
