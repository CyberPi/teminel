package load

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	workingPath := source.Archive.formatWorkingPath(name, cache)
	if utils.VerifyPath(barePath) {
		os.RemoveAll(barePath)
	}
	return utils.CopyDirectory(filepath.Join(workingPath, ".git"), barePath)
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
			SingleBranch: false,
			Tags:         git.AllTags,
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
		url := fmt.Sprintf("https://%v/%v/%v/%v.tar.gz", source.Host, name, source.Archive, version)
		fmt.Print("Ensuring tarball repository:", name)
		if utils.RequestShouldSucceed(url) {
			fmt.Println(" version:", version, "found")
			workingPath := source.formatWorkingPath(name, path)

			clearGitRepository(workingPath)

			repository, err := openOrInit(workingPath, version)
			if err != nil {
				fmt.Println("Error on opening repo:", err)
				return err
			}

			worktree, err := repository.Worktree()
			if err != nil {
				fmt.Println("Error on getting worktree:", err)
				return err
			}
			ref := plumbing.NewBranchReferenceName(version)
			fmt.Println("Ref:", ref)
			err = worktree.Checkout(&git.CheckoutOptions{
				Branch: ref,
				Force:  true,
				Create: !checkBranchExists(repository, version),
			})
			if err != nil {
				fmt.Println("Error on branch creation:", err)
				return err
			}

			err = LoadTarball(url, workingPath)
			if err != nil {
				fmt.Println("Tarball failed to load:", err)
				return err
			}

			if err != nil {
				fmt.Println("Failed to open or init repository:", err)
				return err
			}
			err = commit(repository)
			if err != nil {
				return err
			}
		} else {
			fmt.Println(" version:", version, "not found")
		}
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
		repository, err := git.PlainOpen(path)
		if err != nil {
			return nil, err
		}
		return repository, nil
	} else {
		options := &git.PlainInitOptions{
			InitOptions: git.InitOptions{
				DefaultBranch: plumbing.NewBranchReferenceName(branch),
			},
			Bare:         false,
			ObjectFormat: config.DefaultObjectFormat,
		}
		repository, err := git.PlainInitWithOptions(path, options)
		if err != nil {
			return nil, err
		}
		worktree, err := repository.Worktree()
		if err != nil {
			return nil, err
		}
		filename := "README.md"
		_, _ = worktree.Filesystem.Create(filename)
		commit(repository)
		return repository, nil
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

func clearGitRepository(repositoryPath string) error {
	return filepath.Walk(repositoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == repositoryPath || strings.Contains(path, ".git") {
			return nil
		}
		println("Removing:", path, "type:", info.Mode())
		return os.RemoveAll(path)
	})
}

func checkBranchExists(repo *git.Repository, branchName string) bool {
	branchRefName := plumbing.NewBranchReferenceName(branchName)

	refs, err := repo.References()
	if err != nil {
		return false
	}

	exists := false
	err = refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name() == branchRefName {
			exists = true
		}
		return nil
	})

	return exists
}
