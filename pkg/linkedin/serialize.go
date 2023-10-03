package linkedin

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

const CSVSeparator = '|'

type Serializable interface {
	CsvContent() string
	CsvHeader() string
	Json() string
}

func CsvContent(s Serializable) string {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	var csvContent []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		csvTag := field.Tag.Get("csv")
		if csvTag != "" && csvTag != "-" {
			var value string
			switch v.Field(i).Kind() {
			case reflect.String:
				value = v.Field(i).String()
				value = strings.ReplaceAll(value, string(CSVSeparator), " ")
			case reflect.Bool:
				value = fmt.Sprintf("%v", v.Field(i).Bool())
			case reflect.Int:
				if v.Field(i).Int() != 0 {
					value = fmt.Sprintf("%d", v.Field(i).Int())
				}
			}
			csvContent = append(csvContent, value)
		}
	}
	return strings.Join(csvContent, string(CSVSeparator))
}

func CsvHeader(s Serializable) string {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	var csvHeader []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		csvTag := field.Tag.Get("csv")

		if csvTag != "" && csvTag != "-" {
			csvHeader = append(csvHeader, csvTag)
		}
	}
	return strings.Join(csvHeader, string(CSVSeparator))
}

func Json(s Serializable) string {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		typeName := reflect.TypeOf(s).Elem().Name()
		return fmt.Sprintf("error converting %s to JSON: %v", typeName, err)
	}
	return string(data)
}

type ListSerializable interface {
	Len() int
	Get(i int) Serializable
}

func ConvertToJSON(ls ListSerializable) string {
	if ls.Len() == 0 {
		return "[]"
	}

	var validItems []Serializable
	for i := 0; i < ls.Len(); i++ {
		item := ls.Get(i)
		if item != nil {
			validItems = append(validItems, item)
		}
	}

	data, err := json.MarshalIndent(validItems, "", "  ")
	if err != nil {
		return fmt.Sprintf("error converting slice to JSON: %v", err)
	}
	return string(data)
}

func ConvertToCSV(ls ListSerializable) string {
	if ls.Len() == 0 {
		return ""
	}

	var csvOutput string
	for i := 0; i < ls.Len(); i++ {
		item := ls.Get(i)
		if item != nil {
			csvOutput = item.CsvHeader() + "\n"
			break
		}
	}

	for i := 0; i < ls.Len(); i++ {
		item := ls.Get(i)
		if item != nil {
			csvOutput += item.CsvContent() + "\n"
		}
	}
	return csvOutput
}
