package load

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"source.cyberpi.de/go/teminel/utils"
)

func LoadTarball(url string, path string) error {
	fmt.Println("Loading tarball at:", url)
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	if response.StatusCode >= 300 || response.StatusCode < 200 {
		return fmt.Errorf("Error on tarball retrieval: %v", response.Status)
	}
	defer response.Body.Close()
	gzipReader, err := gzip.NewReader(response.Body)
	if err != nil {
		return err
	}
	defer gzipReader.Close()
	if err := Untar(gzipReader, path); err != nil {
		return err
	}
	return nil
}

func Untar(archive io.Reader, path string) error {
	if err := utils.EnsureDirectory(path); err != nil {
		return err
	}
	reader := tar.NewReader(archive)
	root, err := SeekRoot(reader)
	if err != nil {
		return err
	}
	for {
		header, err := reader.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		itemPath := filepath.Join(path, header.Name[len(root):])
		switch header.Typeflag {
		case tar.TypeDir:
			if err := UntarDir(itemPath, header); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := UntarFile(itemPath, reader, header); err != nil {
				return err
			}
		}
	}
}

func SeekRoot(reader *tar.Reader) (string, error) {
	for {
		header, err := reader.Next()
		if err == io.EOF {
			return "", fmt.Errorf("Tarball does not contain root dir")
		}
		if err != nil {
			return "", err
		}
		switch header.Typeflag {
		case tar.TypeDir:
			return header.Name, nil
		}
	}
}

func UntarDir(path string, header *tar.Header) error {
	return os.MkdirAll(path, os.FileMode(header.Mode))
}

func UntarFile(path string, reader *tar.Reader, header *tar.Header) error {
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		return err
	}
	if err := os.Chmod(path, os.FileMode(header.Mode)); err != nil {
		return err
	}
	if _, err := io.Copy(file, reader); err != nil {
		return err
	}
	return nil
}
