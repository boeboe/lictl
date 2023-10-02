package linkedin

import (
	"testing"
)

func TestPostCsvContent(t *testing.T) {
	tests := []struct {
		name     string
		post     Post
		expected string
	}{
		{
			name: "happy path",
			post: Post{
				Author:       12345,
				AuthorTitle:  "Software|Engineer",
				CommentCount: 10,
				HashTags:     "#tech #golang",
				LikesCount:   100,
				PostLink:     "https://linkedin.com/post/12345",
				PublishDate:  "2023-09-28",
			},
			expected: "12345|Software Engineer|10|#tech #golang|100|https://linkedin.com/post/12345|2023-09-28",
		},
		{
			name:     "empty post",
			post:     Post{},
			expected: "||||||",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.post.CsvContent()
			if got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}

func TestPostCsvHeader(t *testing.T) {
	p := Post{}
	expected := "author|authorTitle|commmentCount|hashTags|likesCount|postLink|publishDate"
	got := p.CsvHeader()
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func TestPostJson(t *testing.T) {
	tests := []struct {
		name     string
		post     Post
		expected string
	}{
		{
			name: "happy path",
			post: Post{
				Author:       12345,
				AuthorTitle:  "Software Engineer",
				CommentCount: 10,
				HashTags:     "#tech #golang",
				LikesCount:   100,
				PostLink:     "https://linkedin.com/post/12345",
				PublishDate:  "2023-09-28",
			},
			expected: `{
  "author": 12345,
  "authorTitle": "Software Engineer",
  "commmentCount": 10,
  "hashTags": "#tech #golang",
  "likesCount": 100,
  "postLink": "https://linkedin.com/post/12345",
  "publishDate": "2023-09-28"
}`,
		},
		{
			name: "empty post",
			post: Post{},
			expected: `{
  "author": 0,
  "authorTitle": "",
  "commmentCount": 0,
  "hashTags": "",
  "likesCount": 0,
  "postLink": "",
  "publishDate": ""
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.post.Json()
			if got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}
