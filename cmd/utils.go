package cmd

import (
	"os"
	"path/filepath"
)

func checkFilePathExist(path string) bool {
	chainConfigFilePath := path

	if !filepath.IsAbs(path) {
		chainConfigFilePath, _ = filepath.Abs(chainConfigFilePath)
	}

	if _, err := os.Stat(chainConfigFilePath); os.IsNotExist(err) {
		errorLogger.Println("File does not exist.")
		return false
	}

	return true
}
