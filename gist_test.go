package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestFetchLatestGist(t *testing.T) {
	gistID := "280bc244d9616f2c91f6a361ae58e05b"
	result, err := fetchLatestGist(gistID)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(strings.Split(result, "\n"))
	t.Skip("okok")
}
