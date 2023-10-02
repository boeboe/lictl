package linkedin

import (
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"testing"
)

func TestGetCompanyPage(t *testing.T) {
	// Define the test matrix
	tests := []struct {
		fileName              string
		expectedFollowerCount int
		expectedFoundedOn     bool
		expectedHeadline      bool
		expectedSpecialties   bool
	}{
		{"company-0.html", 78, false, false, false},
		{"company-1.html", 64494, true, true, true},
		{"company-2.html", 4699, true, true, true},
		{"company-3.html", 3848, true, false, false},
		{"company-4.html", 6938, true, true, true},
		{"company-5.html", 15789, true, true, true},
	}

	// Directory containing test HTML files
	_, filename, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(filename)
	testDir := filepath.Join(basepath, "../..", "testdata", "company")

	// Start a local HTTP server to serve the test files
	server, addr := startLocalHTTPServer(testDir)
	defer server.Close()

	// Iterate over the test matrix
	for _, tt := range tests {
		t.Run(tt.fileName, func(t *testing.T) {
			companyUrl := fmt.Sprintf("http://%s/%s", addr, tt.fileName)
			req, err := http.NewRequest("GET", companyUrl, nil)
			if err != nil {
				t.Fatalf("Error creating HTTP request: %v", err)
			}
			company, err := GetCompanyPage(req, false)
			if err != nil {
				t.Fatalf("Error in SearchJobsPerPage for file %s: %s", tt.fileName, err)
			}

			// Mandatory company information
			if company.FollowerCount != tt.expectedFollowerCount {
				t.Errorf("Expected company.FollowerCount set %d for file %s, but got %d", tt.expectedFollowerCount, tt.fileName, company.FollowerCount)
			}
			if len(company.Headquarters) == 0 {
				t.Errorf("Expected company.Headquarters not to be empty")
			}
			if len(company.Industry) == 0 {
				t.Errorf("Expected company.Industry not to be empty")
			}
			if len(company.Name) == 0 {
				t.Errorf("Expected company.Name not to be empty")
			}
			if len(company.Size) == 0 {
				t.Errorf("Expected company.Size not to be empty")
			}
			if len(company.Type) == 0 {
				t.Errorf("Expected company.Type not to be empty")
			}
			if len(company.Website) == 0 {
				t.Errorf("Expected company.Website not to be empty")
			}

			// Optional company information
			foundedOnSet := len(company.FoundedOn) > 0
			if foundedOnSet != tt.expectedFoundedOn {
				t.Errorf("Expected company.FoundedOn set %t for file %s, but got %t", tt.expectedFoundedOn, tt.fileName, foundedOnSet)
			}
			headlineSet := len(company.Headline) > 0
			if headlineSet != tt.expectedHeadline {
				t.Errorf("Expected company.Headline set %t for file %s, but got %t", tt.expectedHeadline, tt.fileName, headlineSet)
			}
			specialtiesSet := len(company.Specialties) > 0
			if specialtiesSet != tt.expectedSpecialties {
				t.Errorf("Expected company.Specialties set %t for file %s, but got %t", tt.expectedSpecialties, tt.fileName, specialtiesSet)
			}
		})
	}
}
