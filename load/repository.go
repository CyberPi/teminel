package load

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"source.cyberpi.de/go/teminel/utils"
)

const bareDirectory = "/tmp/teminel/bare"
const workingDirectory = "/tmp/teminel/working"

func CloneBare(protocol string, host string, path string, repositoryName string) error {
	bareRepository := repositoryName + ".git"
	barePath := filepath.Join(path, bareRepository)
	workingPath := filepath.Join(workingDirectory, repositoryName)
	// Check if working copy needs an update
	if utils.VerifyPath(workingPath) {
		updateTarballRepo(host, workingDirectory, repositoryName)
	}

	if !utils.VerifyPath(barePath) {
		fmt.Println("Loading repository:", repositoryName, "Using git clone from:", host)
		options := &git.CloneOptions{
			URL:          host + bareRepository,
			SingleBranch: true,
			Depth:        1,
			Tags:         git.NoTags,
		}
		_, err := git.PlainClone(barePath, true, options)
		if err != nil {
			fmt.Println("Try to get archive download due to error:", err)
			repository, err := git.PlainInit(workingPath, false)
			if err != nil {
				return err
			}
			worktree, err := repository.Worktree()
			if err != nil {
				return err
			}
			err = worktree.AddGlob(".")
			if err != nil {
				return err
			}

			options.URL = workingPath
			_, err = git.PlainClone(barePath, true, options)
			return err
		}
	}
	return nil
}

func updateTarballRepo(host string, path string, repositoryName string) error {
	url := fmt.Sprintf("https://%v/%v/archive/refs/heads/main.tar.gz", host, repositoryName)
	workingPath := filepath.Join(path, repositoryName)
	err := LoadTarball(url, workingPath)
	if err != nil {
		return err
	}
	repository, err := openOrInit(workingPath)
	if err != nil {
		return err
	}
	err = commit(repository)
	if err != nil {
		return err
	}
	return nil
}

func openOrInit(path string) (*git.Repository, error) {
	if utils.VerifyPath(path + "/.git") {
		return git.PlainOpen(path)
	} else {
		return git.PlainInit(path, false)
	}
}

func commit(repository *git.Repository) error {
	worktree, err := repository.Worktree()
	if err != nil {
		return err
	}
	status, err := worktree.Status()
	if err != nil {
		return err
	}
	if !status.IsClean() {
		options := &git.AddOptions{All: true}
		worktree.AddWithOptions(options)
		_, err := worktree.Commit("New commit", &git.CommitOptions{
			Author: &object.Signature{
				Name:  "Auto Maton",
				Email: "teminel@cyberpi.de",
				When:  time.Now(),
			},
		})
		if err != nil {
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
		fmt.Println("Found item:", itemPath, "with type flag:", header.Typeflag, "and mode:", header.Mode)
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
	fmt.Println("Creating folder:", path)
	return os.MkdirAll(path, os.FileMode(header.Mode))
}

func UntarFile(path string, reader *tar.Reader, header *tar.Header) error {
	fmt.Println("Creating file:", path)
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
