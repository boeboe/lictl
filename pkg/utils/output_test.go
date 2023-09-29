package utils

import (
	"bufio"
	"encoding/csv"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

type TestStruct struct {
	FieldOne   string `json:"fieldone" csv:"fieldone"`
	FieldTwo   string `json:"fieldtwo" csv:"fieldtwo"`
	FieldThree string `json:"fieldthree" csv:"fieldthree"`
}

func TestOutputFunctions(t *testing.T) {
	// Setup: Create a temporary directory for the tests
	_, filename, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(filename)
	testDir := filepath.Join(basepath, "../..", "testdata", "output")
	dir, err := os.MkdirTemp(testDir, "test_output")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Register cleanup function to run after all tests/subtests
	t.Cleanup(func() {
		if t.Failed() {
			t.Logf("Tests failed. Temporary directory retained at: %s", dir)
		} else {
			os.RemoveAll(dir)
		}
	})

	t.Run("TestDumpToJSON", func(t *testing.T) {
		testDumpToJSON(t, dir)
	})

	t.Run("TestDumpToCSV", func(t *testing.T) {
		testDumpToCSV(t, dir)
	})
}

func testDumpToJSON(t *testing.T, dir string) {
	tests := []struct {
		name     string
		input    []TestStruct
		expected string
	}{
		{
			name: "Happy Path",
			input: []TestStruct{
				{
					FieldOne:   "Test1",
					FieldTwo:   "Test2",
					FieldThree: "Test3",
				},
			},
			expected: `[
  {
    "fieldone": "Test1",
    "fieldtwo": "Test2",
    "fieldthree": "Test3"
  }
]`,
		},
		{
			name:     "Empty Struct",
			input:    []TestStruct{},
			expected: "[]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename, err := DumpToJSON(tt.input, dir)
			if err != nil {
				t.Fatalf("Failed to dump to JSON: %v", err)
			}

			files, _ := os.ReadDir(dir)
			if len(files) == 0 {
				t.Fatal("No files found in the directory")
			}

			content, _ := os.ReadFile(filename)
			if strings.TrimSpace(string(content)) != strings.TrimSpace(tt.expected) {
				t.Fatalf("Expected %s but got %s", tt.expected, string(content))
			}
		})
	}
}

func testDumpToCSV(t *testing.T, dir string) {
	tests := []struct {
		name     string
		input    []TestStruct
		expected string
	}{
		{
			name: "Happy Path",
			input: []TestStruct{
				{
					FieldOne:   "Test1",
					FieldTwo:   "Test2|WithPipe",
					FieldThree: "Test3",
				},
			},
			expected: "fieldone|fieldtwo|fieldthree\nTest1|Test2 WithPipe|Test3",
		},
		{
			name:     "Empty Struct",
			input:    []TestStruct{},
			expected: "fieldone|fieldtwo|fieldthree",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename, err := DumpToCSV(tt.input, dir)
			if err != nil {
				t.Fatalf("Failed to dump to CSV: %v", err)
			}

			files, _ := os.ReadDir(dir)
			if len(files) == 0 {
				t.Fatal("No files found in the directory")
			}

			file, _ := os.Open(filename)
			defer file.Close()

			reader := csv.NewReader(bufio.NewReader(file))
			reader.Comma = CSVSeparator
			records, _ := reader.ReadAll()

			var result []string
			for _, record := range records {
				result = append(result, strings.Join(record, "|"))
			}

			if strings.Join(result, "\n") != tt.expected {
				t.Fatalf("Expected %s but got %s", tt.expected, strings.Join(result, "\n"))
			}
		})
	}
}
