package entities

import "github.com/gosimple/slug"

// xxx_at are all in ISO 8601.
type Blog struct {
	ID          int    `json:"id"`
	Created_at  string `json:"created_at"`
	Updated_at  string `json:"updated_at"`
	Deleted_at  string `json:"deleted_at"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Description string `json:"description"`
	Slug        string `json:"slug"`
	Pined       bool   `json:"pined"`
	Visible     bool   `json:"visible"`
}

func (b *Blog) GenSlug() {
	b.Slug = slug.Make(b.Title)
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
