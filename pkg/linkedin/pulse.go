package linkedin

import (
	"encoding/json"
	"fmt"
	"time"
)

// Pulse represents the structure of a LinkedIn pulse.
type Pulse struct {
	Author       int      `json:"author"`
	AuthorTitle  string   `json:"authorTitle"`
	CommentCount int      `json:"commmentCount"`
	HashTags     []string `json:"hashTags"`
	LikesCount   int      `json:"likesCount"`
	PublishDate  string   `json:"publishDate"`
	PulseLink    string   `json:"pulseLink"`
	Title        string   `json:"title"`
}

func (p Pulse) Dump() string {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Sprintf("error dumping pulse: %v", err)
	}
	return string(data)
}

func SearchPulsesOnline(keywords []string, interval time.Duration, debug bool) ([]Pulse, error) {
	return make([]Pulse, 0), nil
}
