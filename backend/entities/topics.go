package entities

import "github.com/gosimple/slug"

// xxx_at are all in ISO 8601.
type Topic struct {
	ID          int    `json:"id"`
	Created_at  string `json:"created_at"`
	Updated_at  string `json:"updated_at"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Slug        string `json:"slug"`
}

func (t *Topic) GenSlug() {
	t.Slug = slug.Make(t.Name)
}

func NewTopic(name, description string) *Topic {
	topic := &Topic{
		Name:        name,
		Description: description,
	}
	topic.GenSlug()
	return topic
}

func NewTopicWithID(id int, name, description string) *Topic {
	topic := &Topic{
		ID:          id,
		Name:        name,
		Description: description,
	}
	topic.GenSlug()
	return topic
}

type InTopic struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
}
