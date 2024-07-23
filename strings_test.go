package main

import (
	"reflect"
	"testing"
)

func TestStartIndexStrOccurencies(t *testing.T) {
	text := "złzłzłz"
	substr := "zł"

	result := startIndexStrOccurencies(text, substr)
	expected := []int{0, 2, 4}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v but got %v", expected, result)
	}
}
