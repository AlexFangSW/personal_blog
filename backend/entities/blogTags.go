package entities

type BlogTag struct {
	BlogID int `json:"blog_id"`
	TagID  int `json:"tag_id"`
}

func NewBlogTag(blogID, tagID int) *BlogTag {
	blogTag := &BlogTag{
		BlogID: blogID,
		TagID:  tagID,
	}
	return blogTag
}
