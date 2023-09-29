package linkedin

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestSearchJobsPerPage(t *testing.T) {
	// Directory containing test HTML files
	_, filename, _, _ := runtime.Caller(0) // Get the current file's directory
	basepath := filepath.Dir(filename)     // Get the directory of the current file
	testDir := filepath.Join(basepath, "../..", "testdata", "job")

	// Start a local HTTP server to serve the test files
	server, addr := startLocalHTTPServer(testDir)
	defer server.Close()

	// Iterate over all files in the test directory
	err := filepath.Walk(testDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			t.Errorf("Failed to access file: %s", err)
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Construct the file URL for the local server
		fileURL := fmt.Sprintf("http://%s/%s", addr, filepath.Base(path))

		// Call the function with the file URL
		jobs, err := SearchJobsPerPage(fileURL)
		if err != nil {
			t.Errorf("Error in SearchJobsPerPage for file %s: %s", path, err)
		}

		// Add any assertions you want, for example:
		if len(jobs) == 0 {
			t.Errorf("No jobs found for file %s", path)
		}

		return nil
	})

	if err != nil {
		t.Errorf("Error walking the path %s: %v", testDir, err)
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
