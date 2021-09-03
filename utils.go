package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

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

func CopyFile(source string, destination string) (int64, error) {
	sourceFileStat, err := os.Stat(source)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", source)
	}

	sourceFile, err := os.Open(source)
	if err != nil {
		return 0, err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(destination)
	if err != nil {
		return 0, err
	}

	defer destinationFile.Close()

	return io.Copy(destinationFile, sourceFile)
}
