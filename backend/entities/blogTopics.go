package entities

type BlogTopic struct {
	BlogID  int `json:"blog_id"`
	TopicID int `json:"tag_id"`
}

func NewBlogTopic(blogID, topicID int) *BlogTopic {
	blogTopic := &BlogTopic{
		BlogID:  blogID,
		TopicID: topicID,
	}
	return blogTopic
}
