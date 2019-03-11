package add

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tamada/rrh/common"
	git "gopkg.in/src-d/go-git.v4"
)

func (add *AddCommand) isExistAndGitRepository(absPath string, path string) error {
	var fmode, err = os.Stat(absPath)
	if err != nil {
		return err
	}
	if !fmode.IsDir() {
		return fmt.Errorf("%s: not directory", path)
	}
	fmode, err = os.Stat(filepath.Join(absPath, ".git"))
	// If the repository of path is submodule, `.git` will be a file to indicate the `.git` directory.
	if os.IsNotExist(err) {
		return fmt.Errorf("%s: not git repository", path)
	}
	return nil
}

func (add *AddCommand) createGroupIfNeeded(db *common.Database, groupName string) error {
	if !db.HasGroup(groupName) {
		if db.Config.GetValue(common.RrhAutoCreateGroup) == "true" {
			var _, err = db.CreateGroup(groupName, "")
			return err
		}
	}
	if db.HasGroup(groupName) {
		return nil
	}
	return fmt.Errorf("%s: group not found", groupName)
}

func checkDuplication(db *common.Database, repoID string, path string) error {
	var repo = db.FindRepository(repoID)
	if repo != nil && repo.Path != path {
		return fmt.Errorf("%s: duplicate repository id", repoID)
	}
	return nil
}

func (add *AddCommand) addRepositoryToGroup(db *common.Database, groupName string, path string, list []error) []error {
	var absPath, _ = filepath.Abs(path)
	var id = filepath.Base(absPath)
	if err1 := add.isExistAndGitRepository(absPath, path); err1 != nil {
		return append(list, err1)
	}
	var repoPath = common.NormalizePath(absPath)
	if err1 := checkDuplication(db, id, absPath); err1 != nil {
		return append(list, err1)
	} else {
		var remotes, err2 = FindRemotes(absPath)
		if err2 != nil {
			return append(list, err2)
		}
		db.CreateRepository(id, repoPath, remotes)
	}
	var err = db.Relate(groupName, id)
	if err != nil {
		return append(list, fmt.Errorf("%s: cannot create relation to group %s", id, groupName))
	}
	return list
}

func (add *AddCommand) AddRepositoriesToGroup(db *common.Database, args []string, groupName string) []error {
	var err = add.createGroupIfNeeded(db, groupName)
	if err != nil {
		return []error{err}
	}

	var errorlist = []error{}
	for _, item := range args {
		errorlist = add.addRepositoryToGroup(db, groupName, item, errorlist)
	}
	return errorlist
}

/*
FindRemotes function returns the remote of the given git repository.
*/
func FindRemotes(path string) ([]common.Remote, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}
	remotes, err := r.Remotes()
	if err != nil {
		return nil, err
	}
	var crs = []common.Remote{}
	for _, remote := range remotes {
		var config = remote.Config()
		crs = append(crs, common.Remote{Name: config.Name, URL: config.URLs[0]})
	}
	return crs, nil
}
