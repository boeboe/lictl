package cmd

import (
	"fmt"
	"os"
	"time"
)

func writeOutput(content, directory, prefix, extension string) (string, error) {
	// If directory is empty, use the current working directory
	if directory == "" {
		var err error
		directory, err = os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current working directory: %v", err)
		}
	}

	// Check if directory exists
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		// Create directory if it doesn't exist
		if err := os.MkdirAll(directory, 0755); err != nil {
			return "", fmt.Errorf("failed to create directory: %v", err)
		}
	} else if err != nil {
		return "", fmt.Errorf("failed to check directory: %v", err)
	}

	// Check if directory is writable
	if dir, err := os.OpenFile(directory, os.O_WRONLY, 0755); err != nil {
		return "", fmt.Errorf("directory is not writable: %v", err)
	} else {
		dir.Close()
	}

	// Create the unique file
	file, err := createUniqueFile(directory, prefix, extension)
	if err != nil {
		return "", fmt.Errorf("failed to create unique file: %v", err)
	}
	defer file.Close()

	// Write content to file
	if _, err := file.WriteString(content); err != nil {
		return "", fmt.Errorf("failed to write to file: %v", err)
	}

	return file.Name(), nil
}

func createUniqueFile(directory, prefix, extension string) (*os.File, error) {
	currentTime := time.Now()
	timestamp := fmt.Sprintf("%s-%d", currentTime.Format("2006-01-02T15-04-05"), currentTime.Nanosecond())
	filename := fmt.Sprintf("%s/%s_%s.%s", directory, prefix, timestamp, extension)
	return os.Create(filename)
}
