package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func CheckPathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func EnsureDirectories(paths ...string) error {
	for _, path := range paths {
		if !CheckPathExists(path) {
			err := os.MkdirAll(path, os.FileMode(0755))
			if err != nil {
				return err
			}
		}
	}
	return nil
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
