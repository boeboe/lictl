package utils

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"
)

const CSVSeparator = '|'

// getStructInfo returns the name of the struct and its fields for CSV.
func getStructInfoForCSV(data interface{}) (string, []string) {
	t := reflect.TypeOf(data).Elem()
	structName := strings.ToLower(t.Name())

	var fields []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		csvTag := field.Tag.Get("csv")

		if csvTag != "" && csvTag != "-" {
			fields = append(fields, csvTag)
		} else {
			fields = append(fields, strings.ToLower(field.Name))
		}
	}

	return structName, fields
}

// DumpToJSON writes a slice of structs to a JSON file.
func DumpToJSON(data interface{}, dir string) (string, error) {
	structName := strings.ToLower(reflect.TypeOf(data).Elem().Name())

	file, err := createUniqueFile(structName, "json", dir)
	if err != nil {
		return "", err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	err = encoder.Encode(data)

	return file.Name(), err
}

// DumpToCSV writes a slice of structs to a CSV file.
func DumpToCSV(data interface{}, dir string) (string, error) {
	structName, fields := getStructInfoForCSV(data)

	file, err := createUniqueFile(structName, "csv", dir)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Comma = CSVSeparator

	// Write the headers first
	if err := writer.Write(fields); err != nil {
		return file.Name(), fmt.Errorf("failed to write headers: %w", err)
	}

	s := reflect.ValueOf(data)
	// Only write the data if the slice has elements
	if s.Len() > 0 {
		for i := 0; i < s.Len(); i++ {
			var record []string
			for j := 0; j < len(fields); j++ {
				value := s.Index(i).Field(j).String()
				record = append(record, strings.ReplaceAll(value, string(CSVSeparator), " "))
			}
			if err := writer.Write(record); err != nil {
				return file.Name(), fmt.Errorf("failed to write recored: %w", err)
			}
		}
	}

	return file.Name(), nil
}

type Dumper interface {
	Dump() string
}

func DumpFallback(slice interface{}) {
	fmt.Println("Falling back to printing users:")

	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		fmt.Println("Error: provided item is not a slice")
		return
	}

	for i := 0; i < s.Len(); i++ {
		item := s.Index(i).Interface()
		if dumper, ok := item.(Dumper); ok {
			fmt.Println(dumper.Dump())
		} else {
			fmt.Println("Error: item does not implement Dumper interface")
		}
	}
}

func createUniqueFile(structName, extension, dir string) (*os.File, error) {
	currentTime := time.Now()
	timestamp := fmt.Sprintf("%s-%d", currentTime.Format("2006-01-02T15-04-05"), currentTime.Nanosecond())
	filename := fmt.Sprintf("%s/%s_%s.%s", dir, structName, timestamp, extension)
	return os.Create(filename)
}
