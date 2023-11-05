package load

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"source.cyberpi.de/go/teminel/utils"
)

const loaderDir = "/tmp/teminel/loader"

func CloneBare(hostUrl string, path string, repository string) error {
	bareRepository := fmt.Sprintf("%v.git", repository)
	barePath := filepath.Join(path, bareRepository)
	if !utils.VerifyPath(barePath) {
		fmt.Println("Loading repository:", repository, "Using git clone from:", hostUrl)
		options := &git.CloneOptions{
			URL:          hostUrl + bareRepository,
			SingleBranch: true,
			Depth:        1,
			Tags:         git.NoTags,
		}
		_, err := git.PlainClone(barePath, true, options)
		if err != nil {
			fmt.Println("Try to get archive download due to error:", err)
			url := fmt.Sprintf("https://github.com/%v/archive/refs/heads/main.tar.gz", repository)
			repositoryPath := filepath.Join(loaderDir, repository)
			err = LoadTarball(url, repositoryPath)
			if err != nil {
				return err
			}
			_, err = git.PlainInit(repositoryPath, false)
			if err != nil {
				return err
			}
			//options.URL = repositoryPath
			//_, err := git.PlainClone(barePath, true, options)
			return err
		}
	}
	return nil
}

func LoadTarball(url string, path string) error {
	fmt.Println("Loading tarball at:", url)
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	gzipReader, err := gzip.NewReader(response.Body)
	if err != nil {
		return err
	}
	defer gzipReader.Close()
	fmt.Println("Untaring tarball")
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
	for {
		header, err := reader.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		itemPath := filepath.Join(path, header.Name)
		fmt.Println("Found item:", itemPath, "with type flag:", header.Typeflag, "and mode:", header.Mode)
		switch header.Typeflag {
		case tar.TypeDir:
			fmt.Println("Creating folder:", itemPath)
			if err := os.MkdirAll(itemPath, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			fmt.Println("Creating file:", itemPath)
			file, err := os.Create(itemPath)
			if err != nil {
				return err
			}
			if err := os.Chmod(itemPath, os.FileMode(header.Mode)); err != nil {
				return err
			}
			if _, err := io.Copy(file, reader); err != nil {
				return err
			}
			file.Close()
		}
	}
}
