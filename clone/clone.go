package clone

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/tamada/rrh/add"
	"github.com/tamada/rrh/common"
)

func (clone *CloneCommand) toDir(db *common.Database, url string, dest string, repoID string) (*common.Repository, error) {
	clone.printIfVerbose(fmt.Sprintf("git clone %s %s (%s)", url, dest, repoID))
	var cmd = exec.Command("git", "clone", url, dest)
	var err = cmd.Run()
	if err != nil {
		fmt.Printf("clone error: %s\n", err.Error())
		return nil, err
	}

	path, err := filepath.Abs(dest)
	if err != nil {
		return nil, err
	}
	remotes, err := add.FindRemotes(path)
	if err != nil {
		return nil, err
	}
	repo, err := db.CreateRepository(repoID, path, remotes)
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

func (clone *CloneCommand) DoClone(db *common.Database, arguments []string) (int, []error) {
	if len(arguments) == 1 {
		var err = clone.DoCloneARepository(db, arguments[0])
		if err != nil {
			return 0, []error{err}
		}
		return 1, []error{}
	}
	var errorlist = []error{}
	var count = 0
	for _, url := range arguments {
		var increment, err = clone.DoCloneRepositories(db, url)
		if err != nil {
			errorlist = append(errorlist, err)
			if db.Config.GetValue(common.RrhOnError) == common.FailImmediately {
				return count, errorlist
			}
		}
		count += increment
	}
	return count, errorlist
}

func (clone *CloneCommand) relateTo(db *common.Database, groupID string, repoID string) error {
	if !db.HasGroup(groupID) {
		if db.Config.GetValue(common.RrhAutoCreateGroup) == "true" {
			db.CreateGroup(groupID, "")
		} else {
			return fmt.Errorf("%s: group not found", groupID)
		}
	}
	db.Relate(groupID, repoID)
	return nil
}

func (clone *CloneCommand) DoCloneRepositories(db *common.Database, url string) (int, error) {
	var count int
	var id = findID(url)
	var path = filepath.Join(clone.Options.dest, id)
	var _, err = clone.toDir(db, url, path, id)
	if err == nil {
		if err := clone.relateTo(db, clone.Options.group, id); err != nil {
			return count, err
		}
		count++
	}
	return count, err
}

func (clone *CloneCommand) DoCloneARepository(db *common.Database, URL string) error {
	var id, path string

	if clone.isExistDir(clone.Options.dest) {
		id = findID(URL)
		path = filepath.Join(clone.Options.dest, id)
	} else {
		var _, newid = filepath.Split(clone.Options.dest)
		path = clone.Options.dest
		id = newid
	}
	var _, err = clone.toDir(db, URL, path, id)
	if err != nil {
		return err
	}
	return clone.relateTo(db, clone.Options.group, id)
}

func findID(URL string) string {
	var _, dir = path.Split(URL)
	if strings.HasSuffix(dir, ".git") {
		return strings.TrimSuffix(dir, ".git")
	}
	return dir
}
