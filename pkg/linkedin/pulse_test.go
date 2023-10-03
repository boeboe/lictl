package linkedin

import (
	"testing"
)

func TestPulseCsvContent(t *testing.T) {
	tests := []struct {
		name     string
		pulse    Pulse
		expected string
	}{
		{
			name: "happy path",
			pulse: Pulse{
				Author:       12345,
				AuthorTitle:  "Software Engineer",
				CommentCount: 10,
				HashTags:     "#tech #golang",
				LikesCount:   100,
				PublishDate:  "2023-09-28",
				PulseLink:    "https://linkedin.com/pulse/12345",
				Title:        "Golang in 2023",
			},
			expected: "12345|Software Engineer|10|#tech #golang|100|2023-09-28|https://linkedin.com/pulse/12345|Golang in 2023",
		},
		{
			name:     "empty pulse",
			pulse:    Pulse{},
			expected: "|||||||",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pulse.CsvContent()
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestPulseCsvHeader(t *testing.T) {
	p := Pulse{}
	expected := "author|authorTitle|commmentCount|hashTags|likesCount|publishDate|pulseLink|title"
	got := p.CsvHeader()
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestPulseJson(t *testing.T) {
	tests := []struct {
		name     string
		pulse    Pulse
		expected string
	}{
		{
			name: "happy path",
			pulse: Pulse{
				Author:       12345,
				AuthorTitle:  "Software Engineer",
				CommentCount: 10,
				HashTags:     "#tech #golang",
				LikesCount:   100,
				PublishDate:  "2023-09-28",
				PulseLink:    "https://linkedin.com/pulse/12345",
				Title:        "Golang in 2023",
			},
			expected: `{
  "author": 12345,
  "authorTitle": "Software Engineer",
  "commmentCount": 10,
  "hashTags": "#tech #golang",
  "likesCount": 100,
  "publishDate": "2023-09-28",
  "pulseLink": "https://linkedin.com/pulse/12345",
  "title": "Golang in 2023"
}`,
		},
		{
			name:  "empty pulse",
			pulse: Pulse{},
			expected: `{
  "author": 0,
  "authorTitle": "",
  "commmentCount": 0,
  "hashTags": "",
  "likesCount": 0,
  "publishDate": "",
  "pulseLink": "",
  "title": ""
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pulse.Json()
			if got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}
