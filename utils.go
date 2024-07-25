package main

func compareOccurs[T comparable](slice1 []T, slice2 []T) (same bool) {
	mapOccurs := func(slice []T) map[T]int {
		occurs := make(map[T]int)
		for _, item := range slice {
			occurs[item]++
		}
		return occurs
	}
	map1 := mapOccurs(slice1)
	map2 := mapOccurs(slice2)

	if len(map1) != len(map2) {
		return false
	}

	for key, count1 := range map1 {
		if count2, exists := map2[key]; !exists || count1 != count2 {
			return false
		}
	}

	return true
}
