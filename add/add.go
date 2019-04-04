package add

import (
	"fmt"
	"path/filepath"

	"github.com/tamada/rrh/common"
	git "gopkg.in/src-d/go-git.v4"
)

func checkDuplication(db *common.Database, repoID string, path string) error {
	var repo = db.FindRepository(repoID)
	if repo != nil && repo.Path != path {
		return fmt.Errorf("%s: duplicate repository id", repoID)
	}
	return nil
}

func (add *AddCommand) addRepositoryToGroup(db *common.Database, groupName string, path string) []error {
	var absPath, _ = filepath.Abs(path)
	var id = filepath.Base(absPath)
	if err1 := common.IsExistAndGitRepository(absPath, path); err1 != nil {
		return []error{err1}
	}
	if err1 := checkDuplication(db, id, absPath); err1 != nil {
		return []error{err1}
	}
	var remotes, err2 = FindRemotes(absPath)
	if err2 != nil {
		return []error{err2}
	}
	db.CreateRepository(id, absPath, remotes)

	var err = db.Relate(groupName, id)
	if err != nil {
		return []error{fmt.Errorf("%s: cannot create relation to group %s", id, groupName)}
	}
	return []error{}
}

/*
AddRepositoriesToGroup registers the given repositories to the specified group.
*/
func (add *AddCommand) AddRepositoriesToGroup(db *common.Database, args []string, groupName string) []error {
	var _, err = db.AutoCreateGroup(groupName, "", false)
	if err != nil {
		return []error{err}
	}
	var errorlist = []error{}
	for _, item := range args {
		var list = add.addRepositoryToGroup(db, groupName, item)
		errorlist = append(errorlist, list...)
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
