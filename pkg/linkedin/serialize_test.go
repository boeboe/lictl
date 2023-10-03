package linkedin

import (
	"fmt"
	"testing"
)

type TestStruct struct {
	Name  string `json:"name" csv:"name"`
	Age   int    `json:"age" csv:"age"`
	Alive bool   `json:"-"    csv:"-"`
}

func (t *TestStruct) CsvContent() string {
	if t == nil {
		return ""
	}
	return CsvContent(t)
}

func (t *TestStruct) CsvHeader() string {
	if t == nil {
		return ""
	}
	return CsvHeader(t)
}

func (t *TestStruct) Json() string {
	if t == nil {
		return ""
	}
	return Json(t)
}

type TestStructs []*TestStruct

func (tls TestStructs) Len() int {
	return len(tls)
}

func (tls TestStructs) Get(i int) Serializable {
	return Serializable(tls[i])
}

func TestConvertToJSON(t *testing.T) {
	tls := TestStructs{
		&TestStruct{Name: "John", Age: 30, Alive: true},
		&TestStruct{Name: "Jane", Age: 25, Alive: false},
	}
	expected := `[
  {
    "name": "John",
    "age": 30
  },
  {
    "name": "Jane",
    "age": 25
  }
]`
	result := ConvertToJSON(tls)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestConvertToCSV(t *testing.T) {
	tls := TestStructs{
		&TestStruct{Name: "John", Age: 30, Alive: true},
		&TestStruct{Name: "Jane", Age: 25, Alive: false},
	}
	sep := string(CSVSeparator)
	expected := fmt.Sprintf("name%sage\nJohn%s30\nJane%s25\n", sep, sep, sep)
	result := ConvertToCSV(tls)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}
