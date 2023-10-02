package linkedin

import (
	"encoding/json"
	"fmt"
	"time"
)

// Post represents the structure of a LinkedIn post.
type Post struct {
	Author       int      `json:"author"`
	AuthorTitle  string   `json:"authorTitle"`
	CommentCount int      `json:"commmentCount"`
	HashTags     []string `json:"hashTags"`
	LikesCount   int      `json:"likesCount"`
	PostLink     string   `json:"postLink"`
	PublishDate  string   `json:"publishDate"`
}

func (p Post) Dump() string {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Sprintf("error dumping post: %v", err)
	}
	return string(data)
}

func SearchPostsOnline(keywords []string, interval time.Duration, debug bool) ([]Post, error) {
	return make([]Post, 0), nil
}
