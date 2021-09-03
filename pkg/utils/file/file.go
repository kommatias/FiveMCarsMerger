package file

import (
	"fmt"
	"io"
	"os"
)

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
