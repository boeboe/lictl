package linkedin

import (
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"runtime"
	"testing"
)

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
	testDir := filepath.Join(basepath, "../..", "testdata", "job")

	// Start a local HTTP server to serve the test files
	server, addr := startLocalHTTPServer(testDir)
	defer server.Close()

	// Iterate over the test matrix
	for _, tt := range tests {
		t.Run(tt.fileName, func(t *testing.T) {
			jobSearchURL := fmt.Sprintf("http://%s/%s", addr, tt.fileName)
			jobs, err := SearchJobsPerPage(jobSearchURL, false)
			if err != nil {
				t.Fatalf("Error in SearchJobsPerPage for file %s: %s", tt.fileName, err)
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
