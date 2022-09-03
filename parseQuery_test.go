package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseQuery(t *testing.T) {
	var tests = []struct {
		query    string
		expected map[string]string
	}{
		{"Phone=+1 4345 43598", map[string]string{"Phone": "+1 4345 43598"}},
		{"Company=Apple&Phone=08001076285", map[string]string{"Company": "Apple", "Phone": "08001076285"}},
		{"Website=google.com&Code=435", map[string]string{"Website": "google.com", "Code": "435"}},
		{"Name=Alphabet&Code=24&Website=abc.com", map[string]string{"Name": "Alphabet", "Code": "24", "Website": "abc.com"}},
	}

	for i, tt := range tests {
		testName := fmt.Sprintf("Parse Query %d", i)
		t.Run(testName, func(t *testing.T) {
			mapOutput := parseQuery(tt.query)
			if !reflect.DeepEqual(mapOutput, tt.expected) {
				t.Errorf("got %v, wanted %v", mapOutput, tt.expected)
			}
		})
	}
}
