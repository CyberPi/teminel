package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func VerifyPath(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func EnsureDirectory(path string) error {
	if !VerifyPath(path) {
		return os.Mkdir(path, os.FileMode(0755))
	}
	return nil
}

func EnsureDirectories(paths ...string) error {
	for _, path := range paths {
		err := EnsureDirectory(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func ResolvePath(path string) (string, error) {
	result := path
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		result = filepath.Join(home, strings.Trim(path, "~"))
	}
	result = filepath.Clean(result)
	return result, nil
}
