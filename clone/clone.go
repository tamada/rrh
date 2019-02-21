package clone

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/tamadalab/rrh/add"
	"github.com/tamadalab/rrh/common"
)

func (clone *CloneCommand) toDir(db *common.Database, url string, dest string, repoID string) (*common.Repository, error) {
	fmt.Printf("git clone %s %s (%s)\n", url, dest, repoID)
	var cmd = exec.Command("git", "clone", url, dest)
	var err = cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	remotes, err := add.FindRemotes(common.ToAbsolutePath(dest, db.Config))
	if err != nil {
		return nil, err
	}

	repo, err := db.CreateRepository(repoID, dest, remotes)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func (clone *CloneCommand) isExistDir(path string) bool {
	abs, err := filepath.Abs(path)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	stat, err := os.Stat(abs)
	return !os.IsNotExist(err) && stat.IsDir()
}

func (clone *CloneCommand) DoClone(db *common.Database, options *cloneOptions) (int, []error) {
	if len(options.args) == 1 {
		return 1, []error{clone.DoCloneARepository(db, options, options.args[0])}
	}
	var errorlist = []error{}
	var count = 0
	for _, url := range options.args {
		var id = findID(url)
		var path = filepath.Join(options.dest, id)
		var _, err = clone.toDir(db, url, path, id)
		if err != nil {
			if db.Config.GetValue(common.RrhOnError) == common.FailImmediately {
				return count, []error{err}
			}
			errorlist = append(errorlist, err)
		} else {
			db.Relate(options.group, id)
			count++
		}
	}
	return count, errorlist
}

func (clone *CloneCommand) DoCloneARepository(db *common.Database, options *cloneOptions, URL string) error {
	var id, path string

	if clone.isExistDir(options.dest) {
		id = findID(URL)
		path = filepath.Join(options.dest, id)
	} else {
		var _, newid = filepath.Split(options.dest)
		path = options.dest
		id = newid
	}
	var _, err = clone.toDir(db, URL, path, id)
	if err != nil {
		return err
	}
	db.Relate(options.group, id)
	return nil
}

func findID(URL string) string {
	var _, dir = path.Split(URL)
	if strings.HasSuffix(dir, ".git") {
		return strings.TrimSuffix(dir, ".git")
	}
	return dir
}
