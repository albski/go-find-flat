package main

import (
	"strings"
)

func startIndexStrOccurencies(str string, substr string) []int {
	if len(substr) == 0 {
		return []int{}
	}

	indexes := make([]int, 0)
	offset := 0

	for {
		index := strings.Index(str[offset:], substr)
		if index == -1 {
			break
		}
		indexes = append(indexes, offset+index)
		offset += index + len(substr)
	}

	return indexes
}
