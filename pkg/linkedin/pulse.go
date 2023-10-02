package linkedin

import "time"

// Pulse represents the structure of a LinkedIn pulse.
type Pulse struct {
	Author      int      `json:"author"`
	AuthorTitle string   `json:"authorTitle"`
	HashTags    []string `json:"hashTags"`
	LikesCount  int      `json:"likesCount"`
	PublishDate string   `json:"publishDate"`
	Title       string   `json:"title"`
}

func SearchPulsesOnline(keywords []string, interval time.Duration, debug bool) ([]Pulse, error) {
	return make([]Pulse, 0), nil
}
