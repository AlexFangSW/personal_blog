package entities

import "github.com/gosimple/slug"

// xxx_at are all in ISO 8601.
type Tag struct {
	ID          int    `json:"id"`
	Created_at  string `json:"created_at"`
	Updated_at  string `json:"updated_at"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Slug        string `json:"slug"`
}

func (t *Tag) GenSlug() {
	t.Slug = slug.Make(t.Name)
}

func NewTag(name, description string) *Tag {
	tag := &Tag{
		Name:        name,
		Description: description,
	}
	tag.GenSlug()
	return tag
}

type InTag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
