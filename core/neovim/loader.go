package neovim

import (
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"

	"github.com/go-git/go-git/v5"
	"github.com/sosedoff/gitkit"
)

const loaderDir = "/tmp/teminel/loader"

type Loader struct {
	server *gitkit.Server
	loaded map[string]bool
}

func (cache *Loader) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		fmt.Println("Host:", request.Host, "URL:", request.URL)
		matcher := regexp.MustCompile(`(^.+?)\.git.*$`)
		matches := matcher.FindStringSubmatch(request.URL.String())
		fmt.Println("Matches:", matches)
		_, ok := cache.loaded[matches[1]]
		if !ok {
			options := &git.CloneOptions{
				URL:          fmt.Sprintf("https://%v%v.git", "github.com", matches[1]),
				SingleBranch: true,
				Depth:        1,
				Tags:         git.NoTags,
			}
			_, err := git.PlainClone(filepath.Join(
				mirrorDir,
				fmt.Sprintf("%v.git", matches[1]),
			), true, options)
			if err != nil {
				fmt.Println(err)
			}
			cache.loaded[matches[1]] = true
		}
	}
	cache.server.ServeHTTP(writer, request)
}
