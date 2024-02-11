package utils

import (
	"fmt"
	"io"
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
		return os.MkdirAll(path, os.FileMode(0755))
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

func CopyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func CopyDirectory(src string, dst string) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	sfi, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !sfi.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, sfi.Mode())
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDirectory(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
