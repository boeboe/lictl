package linkedin

import (
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCsvContent(t *testing.T) {
	tests := []struct {
		name     string
		company  Company
		expected string
	}{
		{
			name: "happy path",
			company: Company{
				FollowerCount: 1000,
				FoundedOn:     "2000-01-01",
				Headline:      "Tech Company",
				Headquarters:  "San Francisco",
				Industry:      "Technology",
				Name:          "TechCorp",
				Size:          "100-500",
				Specialties:   "Software|Hardware",
				Type:          "Private",
				Website:       "https://techcorp.com",
			},
			expected: "1000|2000-01-01|Tech Company|San Francisco|Technology|TechCorp|100-500|Software Hardware|Private|https://techcorp.com",
		},
		{
			name:     "empty company",
			company:  Company{},
			expected: "|||||||||",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.company.CsvContent()
			if got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}

func TestCsvHeader(t *testing.T) {
	c := Company{}
	expected := "followerCount|foundedOn|headline|headquarters|industry|name|size|specialties|type|website"
	got := c.CsvHeader()
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func TestJson(t *testing.T) {
	tests := []struct {
		name     string
		company  Company
		expected string
	}{
		{
			name: "happy path",
			company: Company{
				FollowerCount: 1000,
				FoundedOn:     "2000-01-01",
				Headline:      "Tech Company",
				Headquarters:  "San Francisco",
				Industry:      "Technology",
				Name:          "TechCorp",
				Size:          "100-500",
				Specialties:   "Software|Hardware",
				Type:          "Private",
				Website:       "https://techcorp.com",
			},
			expected: `{
  "followerCount": 1000,
  "foundedOn": "2000-01-01",
  "headline": "Tech Company",
  "headquarters": "San Francisco",
  "industry": "Technology",
  "name": "TechCorp",
  "size": "100-500",
  "specialties": "Software|Hardware",
  "type": "Private",
  "website": "https://techcorp.com"
}`,
		},
		{
			name:    "empty company",
			company: Company{},
			expected: `{
  "followerCount": 0,
  "foundedOn": "",
  "headline": "",
  "headquarters": "",
  "industry": "",
  "name": "",
  "size": "",
  "specialties": "",
  "type": "",
  "website": ""
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.company.Json()
			if got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}

func TestGetCompanyFromRequest(t *testing.T) {
	// Define the test matrix
	tests := []struct {
		fileName              string
		expectedFollowerCount int
		expectedFoundedOn     string
		expectedHeadline      string
		expectedHeadquarters  string
		expectedIndustry      string
		expectedName          string
		expectedSize          string
		expectedSpecialties   string
		expectedType          string
		expectedWebsite       string
	}{
		{
			"company-0.html",
			78,
			"",
			"",
			"Jeannette, Pennsylvania",
			"Hospitals and Health Care",
			"Arriba Careers",
			"201-500 employees",
			"",
			"Public Company",
			"https://arribacareers.com/",
		},
		{
			"company-1.html",
			64494,
			"1985",
			"The leading source for biopharma news and jobs.\nConnecting industry pioneers with talented professionals.",
			"West Des Moines, Iowa",
			"Internet News",
			"BioSpace",
			"11-50 employees",
			"biotech jobs, pharma jobs, biotech news, pharma news, and life sciences news",
			"Privately Held",
			"http://www.biospace.com/",
		},
		{"company-2.html",
			4699,
			"1858",
			"Premier beverage production & packaging facility headquartered in La Crosse WI with three additional sites across the US",
			"La Crosse, WI",
			"Beverage Manufacturing",
			"City Brewing Company",
			"1,001-5,000 employees",
			"Contract Manufacturing",
			"Privately Held",
			"http://www.citybrewery.com",
		},
		{"company-3.html",
			3848,
			"2002",
			"",
			"San Luis Obispo, CA",
			"Education Administration Programs",
			"Grade Potential Tutoring",
			"11-50 employees",
			"",
			"Privately Held",
			"http://www.gradepotentialtutoring.com/",
		},
		{"company-4.html",
			6938,
			"2022",
			"K12 educational services company with a vision to end generational poverty and eliminate racial achievement gaps",
			"Blairsville, Pennsylvania",
			"Education",
			"Instructional Empowerment",
			"51-200 employees",
			"School Improvement, Core Instruction, Instructional Leadership, Rigor, Grade Level Proficiency, Upward Mobility, Agency, Autonomy, Self-Regulation, Critical Thinking, EdTech, Evaluation, Classroom Walks, Empowerment, Research, Systems, Team Building, Self-Efficacy, Standards-Based, and 21st Century Skills",
			"Self-Owned",
			"https://www.instructionalempowerment.com/",
		},
		{"company-5.html",
			15789,
			"2018",
			"Going beyond to advance treatments for patients with acid-related disorders",
			"Florham Park, New Jersey",
			"Biotechnology Research",
			"Phathom Pharmaceuticals",
			"51-200 employees",
			"Biotechnology, Pharmaceuticals, R&D, Commercialization, Gastrointestinal , GI, Innovation, and Drug development",
			"Public Company",
			"http://phathompharma.com",
		},
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
			company, err := getCompanyFromRequest(req, false)
			if err != nil {
				t.Fatalf("Error in SearchJobsPerPage for file %s: %s", tt.fileName, err)
			}

			if company.FollowerCount != tt.expectedFollowerCount {
				t.Errorf("Expected company.FollowerCount set %d for file %s, but got %d", tt.expectedFollowerCount, tt.fileName, company.FollowerCount)
			}
			if company.FoundedOn != tt.expectedFoundedOn {
				t.Errorf("Expected company.FoundedOn set %q for file %s, but got %q", tt.expectedFoundedOn, tt.fileName, company.FoundedOn)
			}
			if company.Headline != tt.expectedHeadline {
				t.Errorf("Expected company.Headline set %q for file %s, but got %q", tt.expectedHeadline, tt.fileName, company.Headline)
			}
			if company.Headquarters != tt.expectedHeadquarters {
				t.Errorf("Expected company.Headquarters set %q for file %s, but got %q", tt.expectedHeadquarters, tt.fileName, company.Headquarters)
			}
			if company.Industry != tt.expectedIndustry {
				t.Errorf("Expected company.Industry set %q for file %s, but got %q", tt.expectedIndustry, tt.fileName, company.Industry)
			}
			if company.Name != tt.expectedName {
				t.Errorf("Expected company.Name set %q for file %s, but got %q", tt.expectedName, tt.fileName, company.Name)
			}
			if company.Size != tt.expectedSize {
				t.Errorf("Expected company.Size set %q for file %s, but got %q", tt.expectedSize, tt.fileName, company.Size)
			}
			if company.Specialties != tt.expectedSpecialties {
				t.Errorf("Expected company.Specialties set %q for file %s, but got %q", tt.expectedSpecialties, tt.fileName, company.Specialties)
			}
			if company.Type != tt.expectedType {
				t.Errorf("Expected company.Type set %q for file %s, but got %q", tt.expectedType, tt.fileName, company.Type)
			}
			if company.Website != tt.expectedWebsite {
				t.Errorf("Expected company.Website set %q for file %s, but got %q", tt.expectedWebsite, tt.fileName, company.Website)
			}

		})
	}
}
