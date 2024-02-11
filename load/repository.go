package load

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"source.cyberpi.de/go/teminel/utils"
)

type ArchiveSource struct {
	Host        string
	Versions    []string
	Archive     string
	UseBaseName bool
}

type GitSource struct {
	Archive   *ArchiveSource
	Protocols []string
}

func (source *GitSource) EnsureBareRepository(name string, path string, cache string) error {
	err := source.EnsureRepository(name, cache)
	if err != nil {
		return err
	}
	barePath := source.Archive.formatWorkingPath(name+".git", path)
	if utils.VerifyPath(barePath) {
		os.RemoveAll(barePath)
	}
	options := &git.CloneOptions{
		URL:          source.Archive.formatWorkingPath(name, cache),
		SingleBranch: true,
		Depth:        1,
		Tags:         git.NoTags,
	}
	_, err = git.PlainClone(barePath, true, options)
	return err
}

func (source *GitSource) EnsureRepository(name string, path string) error {
	workingPath := source.Archive.formatWorkingPath(name, path)
	fmt.Println("Ensuring repository:", name, "from host:", source.Archive.Host, "on path:", workingPath)
	if utils.VerifyPath(workingPath) {
		fmt.Println("Updating repository:", name)
		source.Archive.updateRepository(name, path)
	} else {
		fmt.Println("Cloning repository:", name)
		options := &git.CloneOptions{
			SingleBranch: true,
			Depth:        1,
			Tags:         git.NoTags,
		}
		err := fmt.Errorf("No protocol was set")
		for _, protocol := range source.Protocols {
			fmt.Println("Trying to clone repo with:", protocol)
			options.URL = fmt.Sprintf(selectUrlTemplate(protocol), source.Archive.Host, name)
			_, err = git.PlainClone(workingPath, false, options)
			if err == nil {
				fmt.Println("Repo clone successful using:", protocol)
				break
			}
		}
		if err != nil {
			fmt.Println("Unable to clone repo using git:", err)
			return source.Archive.ensureTarballRepository(name, path)
		}
	}
	return nil
}

func selectUrlTemplate(protocol string) string {
	switch protocol {
	case "ssh":
		return protocol + "://git@%v/%v.git"
	}
	return protocol + "://%v/%v.git"
}

func (source *ArchiveSource) updateRepository(name string, path string) error {
	workingPath := source.formatWorkingPath(name, path)
	repository, err := git.PlainOpen(workingPath)
	remotes, err := repository.Remotes()
	if err != nil {
		return nil
	}
	if len(remotes) != 0 {
		worktree, err := repository.Worktree()
		if err != nil {
			return err
		}
		options := &git.PullOptions{
			RemoteName: "origin",
		}
		err = worktree.Pull(options)
		if err == git.NoErrAlreadyUpToDate {

		} else if err != nil {
			return nil
		}
	} else {
		return source.ensureTarballRepository(name, path)
	}
	return nil
}

func (source *ArchiveSource) ensureTarballRepository(name string, path string) error {
	for _, version := range source.Versions {
		fmt.Println("Ensuring tarball repository:", name, "version:", version)
		url := fmt.Sprintf("https://%v/%v/%v/%v.tar.gz", source.Host, name, source.Archive, version)
		workingPath := source.formatWorkingPath(name, path)
		err := LoadTarball(url, workingPath)
		if err != nil {
			fmt.Println("Tarball failed to load:", err)
			continue
		}
		repository, err := openOrInit(workingPath, version)
		if err != nil {
			fmt.Println("Failed to open or init repository:", err)
			return err
		}
		err = commit(repository)
		if err != nil {
			return err
		}
		break
	}
	return nil
}

func (source *ArchiveSource) formatWorkingPath(name string, path string) string {
	if source.UseBaseName {
		return filepath.Join(path, filepath.Base(name))
	} else {
		return filepath.Join(path, name)
	}
}

func openOrInit(path string, branch string) (*git.Repository, error) {
	if utils.VerifyPath(path + "/.git") {
		return git.PlainOpen(path)
	} else {
		options := &git.PlainInitOptions{
			InitOptions: git.InitOptions{
				DefaultBranch: plumbing.NewBranchReferenceName(branch),
			},
			Bare:         false,
			ObjectFormat: config.DefaultObjectFormat,
		}
		return git.PlainInitWithOptions(path, options)
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
