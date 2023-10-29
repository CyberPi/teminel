package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func CheckFileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func ParsePath(path string) (string, error) {
	result := path
	if strings.Contains(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return home, err
		}
		result = filepath.Join(home, strings.Trim(path, "~ \t"))
	}
	result = filepath.Clean(result)
	return result, nil
}
