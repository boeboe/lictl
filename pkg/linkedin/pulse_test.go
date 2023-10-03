package linkedin

import (
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
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
				Author:               "John Doe",
				AuthorLinkedInUrl:    "https://linkedin.com/in/johndoe",
				AuthorTitle:          "Software Engineer",
				CommentCount:         10,
				AuthorFollowingCount: 5000,
				LikesCount:           100,
				PublishDate:          "2023-09-28",
				PulseLink:            "https://linkedin.com/pulse/12345",
				Title:                "Golang in 2023",
			},
			expected: "John Doe|https://linkedin.com/in/johndoe|Software Engineer|10|5000|100|2023-09-28|https://linkedin.com/pulse/12345|Golang in 2023",
		},
		{
			name:     "empty pulse",
			pulse:    Pulse{},
			expected: "||||||||",
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
	expected := "author|authorLinkedInUrl|authorTitle|commmentCount|authorFollowingCount|likesCount|publishDate|pulseLink|title"
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
				Author:               "John Doe",
				AuthorLinkedInUrl:    "https://linkedin.com/in/johndoe",
				AuthorTitle:          "Software Engineer",
				CommentCount:         10,
				AuthorFollowingCount: 5000,
				LikesCount:           100,
				PublishDate:          "2023-09-28",
				PulseLink:            "https://linkedin.com/pulse/12345",
				Title:                "Golang in 2023",
			},
			expected: `{
  "author": "John Doe",
  "authorLinkedInUrl": "https://linkedin.com/in/johndoe",
  "authorTitle": "Software Engineer",
  "commmentCount": 10,
  "authorFollowingCount": 5000,
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
  "author": "",
  "authorLinkedInUrl": "",
  "authorTitle": "",
  "commmentCount": 0,
  "authorFollowingCount": 0,
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

func TestGetPulseFromRequest(t *testing.T) {

	// Define the test matrix
	tests := []struct {
		fileName                     string
		expectedAuthor               string
		expectedAuthorLinkedInUrl    string
		expectedAuthorTitle          string
		expectedCommentCount         int
		expectedAuthorFollowingCount int
		expectedLikesCount           int
		expectedPublishDate          string
		expectedPulseLink            string
		expectedTitle                string
	}{
		{
			"pulse-0.html",
			"Darron B.",
			"https://www.linkedin.com/in/darronreginaldbrown",
			"Senior React Native Software Engineer | I Help Start-Up Non-Technical Founders Build Their Mobile Application | Podcast Host",
			3,
			242,
			4,
			"Published Sep 11, 2023",
			"https://www.linkedin.com/pulse/bill-gates-visionary-founder-who-redefined-start-ups-brown-msis",
			"Bill Gates: The Visionary Founder Who Redefined Start-Ups",
		},
		{
			"pulse-1.html",
			"Santiago Iniguez",
			"https://es.linkedin.com/in/siniguez",
			"President, IE University:  Reinventing Higher Education 欧阳圣德",
			8,
			56288,
			78,
			"Published Sep 30, 2023",
			"https://www.linkedin.com/pulse/so-you-think-youre-entrepreneur-try-elon-musk-test-santiago-iniguez",
			"So you think you’re an entrepreneur? Try the Elon Musk test",
		},
		{
			"pulse-2.html",
			"Jenesis Jones",
			"https://www.linkedin.com/in/jenesisjones",
			"Cybersecurity Professional | Security Analyst | Vulnerability Management Auditor | SIEM Analyst | Splunk",
			0,
			0,
			2,
			"Published Feb 10, 2023",
			"https://www.linkedin.com/pulse/beginners-guide-google-cloud-platform-itsbenefits-jenesis-jones",
			"A Beginner's Guide to Google Cloud Platform & its\u00a0Benefits",
		},
		{
			"pulse-3.html",
			"Manisha S.",
			"https://in.linkedin.com/in/manisha23",
			"Freelance",
			0,
			0,
			2,
			"Published Mar 12, 2023",
			"https://www.linkedin.com/pulse/google-cloud-platformgcp-manisha-sharma",
			"Google Cloud Platform(GCP)",
		},
		{
			"pulse-4.html",
			"Job Lefrandt",
			"https://nl.linkedin.com/in/job-lefrandt",
			"IT sparring partner - Efficiëntie en innovatie staan centraal in onze IT-services. Laten we samenwerken aan uw succes! - 06 38 29 54 62 - j.lefrandt@eshgro.nl",
			0,
			0,
			7,
			"Published Jul 10, 2023",
			"https://nl.linkedin.com/pulse/de-kracht-van-azure-een-nieuw-tijdperk-mogelijkheden-job-lefrandt",
			"De Kracht van Azure: Een Nieuw Tijdperk van Mogelijkheden",
		},
		{
			"pulse-5.html",
			"Paul Erlings",
			"https://nl.linkedin.com/in/paul-erlings-bb903517",
			"Principal Consultant at Rubicon Cloud Advisor",
			0,
			0,
			6,
			"Published Sep 29, 2023",
			"https://nl.linkedin.com/pulse/de-magie-van-windows-laps-azure-ad-ontmaskerd-paul-erlings",
			"De Magie van Windows LAPS Azure AD Ontmaskerd",
		},
	}

	// Directory containing test HTML files
	_, filename, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(filename)
	testDir := filepath.Join(basepath, "../..", "testdata", "pulse")

	// Start a local HTTP server to serve the test files
	server, addr := startLocalHTTPServer(testDir)
	defer server.Close()

	// Iterate over the test matrix
	for _, tt := range tests {
		t.Run(tt.fileName, func(t *testing.T) {
			pulseUrl := fmt.Sprintf("http://%s/%s", addr, tt.fileName)
			req, err := http.NewRequest("GET", pulseUrl, nil)
			if err != nil {
				t.Fatalf("Error creating HTTP request: %v", err)
			}
			pulse, err := getPulseFromRequest(req, false)
			if err != nil {
				t.Fatalf("Error in getPulseFromRequest for file %s: %s", tt.fileName, err)
			}

			if pulse.Author != tt.expectedAuthor {
				t.Errorf("Expected pulse.Author set %q for file %s, but got %q", tt.expectedAuthor, tt.fileName, pulse.Author)
			}
			if pulse.AuthorLinkedInUrl != tt.expectedAuthorLinkedInUrl {
				t.Errorf("Expected pulse.AuthorLinkedInUrl set %q for file %s, but got %q", tt.expectedAuthorLinkedInUrl, tt.fileName, pulse.AuthorLinkedInUrl)
			}
			if pulse.AuthorTitle != tt.expectedAuthorTitle {
				t.Errorf("Expected pulse.AuthorTitle set %q for file %s, but got %q", tt.expectedAuthorTitle, tt.fileName, pulse.AuthorTitle)
			}
			if pulse.CommentCount != tt.expectedCommentCount {
				t.Errorf("Expected pulse.CommentCount set %d for file %s, but got %d", tt.expectedCommentCount, tt.fileName, pulse.CommentCount)
			}
			if pulse.AuthorFollowingCount != tt.expectedAuthorFollowingCount {
				t.Errorf("Expected pulse.AuthorFollowingCount set %d for file %s, but got %d", tt.expectedAuthorFollowingCount, tt.fileName, pulse.AuthorFollowingCount)
			}
			if pulse.LikesCount != tt.expectedLikesCount {
				t.Errorf("Expected pulse.LikesCount set %d for file %s, but got %d", tt.expectedLikesCount, tt.fileName, pulse.LikesCount)
			}
			if pulse.PublishDate != tt.expectedPublishDate {
				t.Errorf("Expected pulse.PublishDate set %q for file %s, but got %q", tt.expectedPublishDate, tt.fileName, pulse.PublishDate)
			}
			if pulse.PulseLink != tt.expectedPulseLink {
				t.Errorf("Expected pulse.PulseLink set %q for file %s, but got %q", tt.expectedPulseLink, tt.fileName, pulse.PulseLink)
			}
			if pulse.Title != tt.expectedTitle {
				t.Errorf("Expected pulse.Title set %q for file %s, but got %q", tt.expectedTitle, tt.fileName, pulse.Title)
			}
		})
	}
}
