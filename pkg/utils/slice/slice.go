package slice

import "strings"

func RemoveDuplicates(strSlice []string) []string {
	allKeys := make(map[string]bool)
	var list []string
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func ContainsElement(list []string, element string) bool {
	for _, listElement := range list {
		if strings.TrimSpace(strings.ToLower(listElement)) == strings.TrimSpace(strings.ToLower(element)) {
			return true
		}
	}
	return false
}
