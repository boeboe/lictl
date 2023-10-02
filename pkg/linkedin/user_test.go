package linkedin

import (
	"testing"
)

func TestUserCsvContent(t *testing.T) {
	tests := []struct {
		name     string
		user     User
		expected string
	}{
		{
			name: "happy path",
			user: User{
				ConnectionCount: 500,
				FollowerCount:   1000,
				UserTitle:       "Software|Developer",
				Location:        "San Francisco, CA",
				Name:            "John Doe",
				UserLink:        "https://linkedin.com/in/johndoe",
			},
			expected: "500|1000|Software Developer|San Francisco, CA|John Doe|https://linkedin.com/in/johndoe",
		},
		{
			name:     "empty user",
			user:     User{},
			expected: "|||||",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.user.CsvContent()
			if got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}

func TestUserCsvHeader(t *testing.T) {
	u := User{}
	expected := "connectionCount|followerCount|userTitle|location|name|userLink"
	got := u.CsvHeader()
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func TestUserJson(t *testing.T) {
	tests := []struct {
		name     string
		user     User
		expected string
	}{
		{
			name: "happy path",
			user: User{
				ConnectionCount: 500,
				FollowerCount:   1000,
				UserTitle:       "Software Developer",
				Location:        "San Francisco, CA",
				Name:            "John Doe",
				UserLink:        "https://linkedin.com/in/johndoe",
			},
			expected: `{
  "connectionCount": 500,
  "followerCount": 1000,
  "userTitle": "Software Developer",
  "location": "San Francisco, CA",
  "name": "John Doe",
  "userLink": "https://linkedin.com/in/johndoe"
}`,
		},
		{
			name: "empty user",
			user: User{},
			expected: `{
  "connectionCount": 0,
  "followerCount": 0,
  "userTitle": "",
  "location": "",
  "name": "",
  "userLink": ""
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.user.Json()
			if got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}
