package main

func startIndexStrOccurencies(str string, substr string) []int {
	if len(substr) == 0 {
		return []int{}
	}

	indexes := make([]int, 0)

	runeStr := []rune(str)
	runeSubstr := []rune(substr)
	strLen := len(runeStr)
	substrLen := len(runeSubstr)

	offset := 0
	for offset <= strLen-substrLen {
		if string(runeStr[offset:offset+substrLen]) == string(runeSubstr) {
			indexes = append(indexes, offset)
			offset += substrLen
		} else {
			offset++
		}
	}

	return indexes
}
