package entities

import (
	"crypto/md5"
	"fmt"

	"github.com/gosimple/slug"
)

// xxx_at are all in ISO 8601.
type Blog struct {
	ID          int    `json:"id"`
	Created_at  string `json:"created_at"`
	Updated_at  string `json:"updated_at"`
	Deleted_at  string `json:"deleted_at"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	ContentMD5  string `json:"contentMD5"`
	Description string `json:"description"`
	Slug        string `json:"slug"`
	Pined       bool   `json:"pined"`
	Visible     bool   `json:"visible"`
}

func (b *Blog) GenSlug() {
	b.Slug = slug.Make(b.Title)
}
func (b *Blog) GenMD5() {
	b.ContentMD5 = fmt.Sprintf("%x", md5.Sum([]byte(b.Content)))
}

func NewBlog(title, content, description string, pined, visible bool) *Blog {
	blog := &Blog{
		Title:       title,
		Content:     content,
		Description: description,
		Pined:       pined,
		Visible:     visible,
	}
	blog.GenSlug()
	blog.GenMD5()
	return blog
}

func NewBlogWithID(id int, title, content, description string, pined, visible bool) *Blog {
	blog := &Blog{
		ID:          id,
		Title:       title,
		Content:     content,
		Description: description,
		Pined:       pined,
		Visible:     visible,
	}
	blog.GenSlug()
	blog.GenMD5()
	return blog
}

type InBlog struct {
	Blog
	Tags   []int `json:"tags"`
	Topics []int `json:"topics"`
}

func NewInBlog(blog Blog, tags, topics []int) *InBlog {
	return &InBlog{
		Blog:   blog,
		Tags:   tags,
		Topics: topics,
	}
}

func NewInBlogWithID(blog Blog, tags, topics []int) *InBlog {
	return &InBlog{
		Blog:   blog,
		Tags:   tags,
		Topics: topics,
	}
}

type OutBlog struct {
	Blog
	Tags   []Tag   `json:"tags"`
	Topics []Topic `json:"topics"`
}

func NewOutBlog(blog Blog, tags []Tag, topics []Topic) *OutBlog {
	return &OutBlog{
		Blog:   blog,
		Tags:   tags,
		Topics: topics,
	}
}

// tags an topics as slugs
type OutBlogSimple struct {
	Blog
	Tags   []string `json:"tags"`
	Topics []string `json:"topics"`
}

func NewOutBlogSimple(blog Blog, tags []string, topics []string) OutBlogSimple {
	return OutBlogSimple{
		Blog:   blog,
		Tags:   tags,
		Topics: topics,
	}
}

type ReqInBlog struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	Description string `json:"description"`
	Pined       bool   `json:"pined"`
	Visible     bool   `json:"visible"`
	Tags        []int  `json:"tags"`
	Topics      []int  `json:"topics"`
}
