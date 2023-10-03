package linkedin

import (
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
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
				ActivityURN:       "urn:li:activity:12345",
				Author:            "John Doe",
				AuthorLinkedInUrl: "https://linkedin.com/in/johndoe",
				AuthorTitle:       "Software|Engineer",
				CommentCount:      10,
				Freshness:         "1 hour ago",
				LikesCount:        100,
				PostLink:          "https://linkedin.com/post/12345",
				PublishDate:       "2023-09-28",
				ShareURN:          "urn:li:share:67890",
			},
			expected: "urn:li:activity:12345|John Doe|https://linkedin.com/in/johndoe|Software Engineer|10||1 hour ago|100|https://linkedin.com/post/12345|2023-09-28|urn:li:share:67890",
		},
		{
			name:     "empty post",
			post:     Post{},
			expected: "||||||||||",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.post.CsvContent()
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestPostCsvHeader(t *testing.T) {
	p := Post{}
	expected := "activityURN|author|authorLinkedInUrl|authorTitle|commmentCount|companyFollowerCount|freshness|likesCount|postLink|publishDate|shareURN"
	got := p.CsvHeader()
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
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
				ActivityURN:          "urn:li:activity:12345",
				Author:               "John Doe",
				AuthorLinkedInUrl:    "https://linkedin.com/in/johndoe",
				AuthorTitle:          "Software Engineer",
				CommentCount:         10,
				CompanyFollowerCount: 5,
				Freshness:            "1 hour ago",
				LikesCount:           100,
				PostLink:             "https://linkedin.com/post/12345",
				PublishDate:          "2023-09-28",
				ShareURN:             "urn:li:share:67890",
			},
			expected: `{
  "activityURN": "urn:li:activity:12345",
  "author": "John Doe",
  "authorLinkedInUrl": "https://linkedin.com/in/johndoe",
  "authorTitle": "Software Engineer",
  "commmentCount": 10,
  "companyFollowerCount": 5,
  "freshness": "1 hour ago",
  "likesCount": 100,
  "postLink": "https://linkedin.com/post/12345",
  "publishDate": "2023-09-28",
  "shareURN": "urn:li:share:67890"
}`,
		},
		{
			name: "empty post",
			post: Post{},
			expected: `{
  "activityURN": "",
  "author": "",
  "authorLinkedInUrl": "",
  "authorTitle": "",
  "commmentCount": 0,
  "companyFollowerCount": 0,
  "freshness": "",
  "likesCount": 0,
  "postLink": "",
  "publishDate": "",
  "shareURN": ""
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

func TestGetPostFromRequest(t *testing.T) {
	// Define the test matrix
	tests := []struct {
		fileName                      string
		expectedActivityURN           string
		expectedAuthor                string
		expectedAuthorLinkedInUrl     string
		expectedAuthorTitle           string
		expectedCommentCount          int
		expectedCompanyFollowerCountt int
		expectedFreshness             string
		expectedLikesCount            int
		expectedPostLink              string
		expectedShareURN              string
	}{
		{
			"post-0.html",
			"urn:li:activity:6983753436961865728",
			"Julien Blanchez",
			"https://be.linkedin.com/in/julienblanchez",
			"Digital Sovereignty Solution Lead at Google",
			0,
			0,
			"12mo",
			4,
			"https://www.linkedin.com/posts/julienblanchez_hacking-google-series-trailer-activity-6983753436961865728-RoKg",
			"urn:li:share:6983753435862999040",
		},
		{
			"post-1.html",
			"urn:li:activity:7105127933211467777",
			"Sasja Nothacker",
			"https://de.linkedin.com/in/sasja-nothacker-a051648a",
			"Empower customers to embrace digital transformation",
			0,
			0,
			"3w",
			11,
			"https://www.linkedin.com/posts/sasja-nothacker-a051648a_all-161-things-we-announced-at-google-cloud-activity-7105127933211467777-POG4",
			"urn:li:share:7105127931164676096",
		},
		{
			"post-2.html",
			"urn:li:activity:7051984245006753792",
			"Ignasi Barrera",
			"https://es.linkedin.com/in/ignasibarrera",
			"Founding Engineer at Tetrate | ASF Member | Building Zero Trust (ZTA) with Envoy and Istio",
			0,
			0,
			"5mo",
			13,
			"https://www.linkedin.com/posts/ignasibarrera_introducing-tetrate-service-express-activity-7051984245006753792-0CTS",
			"urn:li:share:7051984244465647616",
		},
		{
			"post-3.html",
			"urn:li:activity:7039669241490415616",
			"solo.io",
			"https://www.linkedin.com/company/solo.io",
			"",
			0,
			10757,
			"6mo",
			222,
			"https://www.linkedin.com/posts/solo%2Eio_compare-top-api-gateways-activity-7039669241490415616-DNov",
			"urn:li:share:7039669240731238401",
		},
		{
			"post-4.html",
			"urn:li:activity:6979083456680984576",
			"Amazon Web Services (AWS)",
			"https://www.linkedin.com/company/amazon-web-services",
			"",
			0,
			9026507,
			"1y",
			127,
			"https://www.linkedin.com/posts/amazon-web-services_learn-more-about-data-protection-at-aws-activity-6979083456680984576-k1b2",
			"urn:li:ugcPost:6979083432681197568",
		},
		{
			"post-5.html",
			"urn:li:activity:7087150407113719808",
			"Max Uritsky",
			"https://www.linkedin.com/in/max-uritsky-170a441",
			"",
			0,
			0,
			"2mo",
			254,
			"https://www.linkedin.com/posts/max-uritsky-170a441_introducing-microsoft-azure-boost-preview-activity-7087150407113719808-gMLa",
			"urn:li:share:7087150406597840896",
		},
	}

	// Directory containing test HTML files
	_, filename, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(filename)
	testDir := filepath.Join(basepath, "../..", "testdata", "post")

	// Start a local HTTP server to serve the test files
	server, addr := startLocalHTTPServer(testDir)
	defer server.Close()

	// Iterate over the test matrix
	for _, tt := range tests {
		t.Run(tt.fileName, func(t *testing.T) {
			postUrl := fmt.Sprintf("http://%s/%s", addr, tt.fileName)
			req, err := http.NewRequest("GET", postUrl, nil)
			if err != nil {
				t.Fatalf("Error creating HTTP request: %v", err)
			}
			post, err := getPostFromRequest(req, false)
			if err != nil {
				t.Fatalf("Error in getPostFromRequest for file %s: %s", tt.fileName, err)
			}

			if post.ActivityURN != tt.expectedActivityURN {
				t.Errorf("Expected post.ActivityURN set %q for file %s, but got %q", tt.expectedActivityURN, tt.fileName, post.ActivityURN)
			}
			if post.Author != tt.expectedAuthor {
				t.Errorf("Expected post.Author set %q for file %s, but got %q", tt.expectedAuthor, tt.fileName, post.Author)
			}
			if post.AuthorLinkedInUrl != tt.expectedAuthorLinkedInUrl {
				t.Errorf("Expected post.AuthorLinkedInUrl set %q for file %s, but got %q", tt.expectedAuthorLinkedInUrl, tt.fileName, post.AuthorLinkedInUrl)
			}
			if post.AuthorTitle != tt.expectedAuthorTitle {
				t.Errorf("Expected post.AuthorTitle set %q for file %s, but got %q", tt.expectedAuthorTitle, tt.fileName, post.AuthorTitle)
			}
			if post.CommentCount != tt.expectedCommentCount {
				t.Errorf("Expected post.CommentCount set %d for file %s, but got %d", tt.expectedCommentCount, tt.fileName, post.CommentCount)
			}
			if post.Freshness != tt.expectedFreshness {
				t.Errorf("Expected post.Freshness set %q for file %s, but got %q", tt.expectedFreshness, tt.fileName, post.Freshness)
			}
			if post.LikesCount != tt.expectedLikesCount {
				t.Errorf("Expected post.LikesCount set %d for file %s, but got %d", tt.expectedLikesCount, tt.fileName, post.LikesCount)
			}
			if post.PostLink != tt.expectedPostLink {
				t.Errorf("Expected post.PostLink set %q for file %s, but got %q", tt.expectedPostLink, tt.fileName, post.PostLink)
			}
			if post.ShareURN != tt.expectedShareURN {
				t.Errorf("Expected post.ShareURN set %q for file %s, but got %q", tt.expectedShareURN, tt.fileName, post.ShareURN)
			}
		})
	}
}
