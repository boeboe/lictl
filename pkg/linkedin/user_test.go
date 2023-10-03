package linkedin

import (
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
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
				ConnectionCount: "500+",
				FollowerCount:   "1K",
				UserTitle:       "Software|Developer",
				Location:        "San Francisco, CA",
				Name:            "John Doe",
				UserLink:        "https://linkedin.com/in/johndoe",
			},
			expected: "500+|1K|Software Developer|San Francisco, CA|John Doe|https://linkedin.com/in/johndoe",
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
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestUserCsvHeader(t *testing.T) {
	u := User{}
	expected := "connectionCount|followerCount|userTitle|location|name|userLink"
	got := u.CsvHeader()
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
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
				ConnectionCount: "500+",
				FollowerCount:   "1K",
				UserTitle:       "Software Developer",
				Location:        "San Francisco, CA",
				Name:            "John Doe",
				UserLink:        "https://linkedin.com/in/johndoe",
			},
			expected: `{
  "connectionCount": "500+",
  "followerCount": "1K",
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
  "connectionCount": "",
  "followerCount": "",
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

func TestGetUserFromRequest(t *testing.T) {
	// Define the test matrix
	tests := []struct {
		fileName                string
		expectedConnectionCount string
		expectedFollowerCount   string
		expectedUserTitle       string
		expectedLocation        string
		expectedName            string
		expectedUserLink        string
	}{
		{
			"user-0.html",
			"1K followers",
			"500+ connections",
			"Google for Education lead Belgium & Luxembourg",
			"Belgium",
			"Louise Van Lint",
			"https://be.linkedin.com/in/louisevanlint",
		},
		{
			"user-1.html",
			"2K followers",
			"500+ connections",
			"Analytical Lead - CPG & Technology at Google / Board Member at iN²POWER",
			"Antwerp, Flemish Region, Belgium",
			"Sebastiaan Monsieurs",
			"https://be.linkedin.com/in/sebastiaanmonsieurs",
		},
		{
			"user-2.html",
			"3K followers",
			"500+ connections",
			"BUSINESS GROUP LEAD MICROSOFT AZURE at Microsoft",
			"Rijswijk, South Holland, Netherlands",
			"Jeffrey Vermeulen",
			"https://nl.linkedin.com/in/jeffvermeulen",
		},
		{
			"user-3.html",
			"2K followers",
			"500+ connections",
			"Azure Cloud Evangelist @ Cegeka | Driving Digital Transformation",
			"Geel, Flemish Region, Belgium",
			"Ivo Haagen",
			"https://be.linkedin.com/in/ivohaagen",
		},
		{
			"user-4.html",
			"1K followers",
			"500+ connections",
			"Lawyer, Certified Data Privacy Officer (DPO) & EU GDPR Representative",
			"Brussels Metropolitan Area",
			"Gauthier Broze, LL.M, CIPP/E, GCP (Good Clinical Practice)",
			"https://be.linkedin.com/in/gauthier-broze-ll-m-cipp-e-gcp-good-clinical-practice-9612ab5",
		},
		{
			"user-5.html",
			"409 followers",
			"405 connections",
			"Consultant - Process Improvement & Change Management at GCP Consulting",
			"Waremme, Walloon Region, Belgium",
			"Julien Hernaut",
			"https://be.linkedin.com/in/julienhernaut",
		},
	}

	// Directory containing test HTML files
	_, filename, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(filename)
	testDir := filepath.Join(basepath, "../..", "testdata", "user")

	// Start a local HTTP server to serve the test files
	server, addr := startLocalHTTPServer(testDir)
	defer server.Close()

	// Iterate over the test matrix
	for _, tt := range tests {
		t.Run(tt.fileName, func(t *testing.T) {
			userUrl := fmt.Sprintf("http://%s/%s", addr, tt.fileName)
			req, err := http.NewRequest("GET", userUrl, nil)
			if err != nil {
				t.Fatalf("Error creating HTTP request: %v", err)
			}
			user, err := getUserFromRequest(req, false)
			if err != nil {
				t.Fatalf("Error in SearchJobsPerPage for file %s: %s", tt.fileName, err)
			}

			if user.ConnectionCount != tt.expectedConnectionCount {
				t.Errorf("Expected user.ConnectionCount set %q for file %s, but got %q", tt.expectedConnectionCount, tt.fileName, user.ConnectionCount)
			}
			if user.FollowerCount != tt.expectedFollowerCount {
				t.Errorf("Expected user.FollowerCount set %q for file %s, but got %q", tt.expectedFollowerCount, tt.fileName, user.FollowerCount)
			}
			if user.UserTitle != tt.expectedUserTitle {
				t.Errorf("Expected user.UserTitle set %q for file %s, but got %q", tt.expectedUserTitle, tt.fileName, user.UserTitle)
			}
			if user.Location != tt.expectedLocation {
				t.Errorf("Expected user.Location set %q for file %s, but got %q", tt.expectedLocation, tt.fileName, user.Location)
			}
			if user.Name != tt.expectedName {
				t.Errorf("Expected user.Name set %q for file %s, but got %q", tt.expectedName, tt.fileName, user.Name)
			}
			if user.UserLink != tt.expectedUserLink {
				t.Errorf("Expected user.UserLink set %q for file %s, but got %q", tt.expectedUserLink, tt.fileName, user.UserLink)
			}
		})
	}
}
