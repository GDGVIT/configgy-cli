package utils

import (
	"fmt"
	"os"
)

func FilePathExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File does not exist
			return false
		} // Some other error occurred
		panic(fmt.Errorf("failed to check if file exists: %s", err))
	}
	return true
}
