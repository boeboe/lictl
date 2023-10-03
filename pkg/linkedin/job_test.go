package linkedin

import (
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"runtime"
	"testing"
)

func TestJobCsvContent(t *testing.T) {
	tests := []struct {
		name     string
		job      Job
		expected string
	}{
		{
			name: "happy path",
			job: Job{
				CompanyLinkedInURL: "https://linkedin.com/company/techcorp",
				CompanyName:        "Techcorp",
				DatePosted:         "2023-01-01",
				JobLink:            "https://linkedin.com/jobs/view/123456",
				JobTitle:           "Software|Engineer",
				JobURN:             "urn:li:job:123456",
				Location:           "San Francisco, CA",
			},
			expected: "https://linkedin.com/company/techcorp|Techcorp|2023-01-01|https://linkedin.com/jobs/view/123456|Software Engineer|urn:li:job:123456|San Francisco, CA",
		},
		{
			name:     "empty job",
			job:      Job{},
			expected: "||||||",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.job.CsvContent()
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestJobCsvHeader(t *testing.T) {
	j := Job{}
	expected := "companyLinkedInURL|companyName|datePosted|jobLink|jobTitle|jobURN|location"
	got := j.CsvHeader()
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestJobJson(t *testing.T) {
	tests := []struct {
		name     string
		job      Job
		expected string
	}{
		{
			name: "happy path",
			job: Job{
				JobTitle:           "Software Engineer",
				CompanyName:        "Techcorp",
				CompanyLinkedInURL: "https://linkedin.com/company/techcorp",
				Location:           "San Francisco, CA",
				DatePosted:         "2023-01-01",
				JobLink:            "https://linkedin.com/jobs/view/123456",
				JobURN:             "urn:li:job:123456",
			},
			expected: `{
  "companyLinkedInURL": "https://linkedin.com/company/techcorp",
  "companyName": "Techcorp",
  "datePosted": "2023-01-01",
  "jobLink": "https://linkedin.com/jobs/view/123456",
  "jobTitle": "Software Engineer",
  "jobURN": "urn:li:job:123456",
  "location": "San Francisco, CA"
}`,
		},
		{
			name: "empty job",
			job:  Job{},
			expected: `{
  "companyLinkedInURL": "",
  "companyName": "",
  "datePosted": "",
  "jobLink": "",
  "jobTitle": "",
  "jobURN": "",
  "location": ""
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.job.Json()
			if got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}

func TestSearchJobsPerPage(t *testing.T) {
	// Define the test matrix
	tests := []struct {
		fileName             string
		expectedJobsCount    int
		expectedEmptyCompany int
	}{
		{"jobs-0.html", 25, 0},
		{"jobs-25.html", 25, 0},
		{"jobs-50.html", 25, 0},
		{"jobs-75.html", 25, 0},
		{"jobs-100.html", 25, 0},
		{"jobs-final.html", 16, 0},
	}

	// Directory containing test HTML files
	_, filename, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(filename)
	testDir := filepath.Join(basepath, "../..", "testdata", "job-search")

	// Start a local HTTP server to serve the test files
	server, addr := startLocalHTTPServer(testDir)
	defer server.Close()

	// Iterate over the test matrix
	for _, tt := range tests {
		t.Run(tt.fileName, func(t *testing.T) {
			jobSearchURL := fmt.Sprintf("http://%s/%s", addr, tt.fileName)
			jobs, err := GetJobsFromSearchUrl(jobSearchURL, false)
			if err != nil {
				t.Fatalf("Error in GetJobsFromSearchUrl for file %s: %s", tt.fileName, err)
			}

			if len(jobs) != tt.expectedJobsCount {
				t.Errorf("Expected %d jobs for file %s, but got %d", tt.expectedJobsCount, tt.fileName, len(jobs))
			}

			emptyCompanyCount := 0
			for _, job := range jobs {
				if job.CompanyLinkedInURL == "" {
					emptyCompanyCount++
				}
			}

			if emptyCompanyCount != tt.expectedEmptyCompany {
				t.Errorf("Expected %d jobs with empty CompanyLinkedInURL for file %s, but got %d", tt.expectedEmptyCompany, tt.fileName, emptyCompanyCount)
			}
		})
	}
}

func startLocalHTTPServer(dir string) (*http.Server, string) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	server := &http.Server{Handler: http.FileServer(http.Dir(dir))}

	go func() {
		if err := server.Serve(listener); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return server, listener.Addr().String()
}
