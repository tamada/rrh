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

func findID(repoID string, absPath string) string {
	if repoID == "" {
		return filepath.Base(absPath)
	}
	return repoID
}

func (add *Command) addRepositoryToGroup(db *common.Database, repoID string, groupName string, path string) []error {
	var absPath, _ = filepath.Abs(path)
	var id = findID(repoID, absPath)
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

func (add *Command) validateArguments(args []string, repoID string) error {
	if repoID != "" && len(args) > 1 {
		return fmt.Errorf("specifying repository id do not accept multiple arguments: %v", args)
	}
	return nil
}

/*
AddRepositoriesToGroup registers the given repositories to the specified group.
*/
func (add *Command) AddRepositoriesToGroup(db *common.Database, opt *options) []error {
	var _, err = db.AutoCreateGroup(opt.group, "", false)
	if err != nil {
		return []error{err}
	}
	if err := add.validateArguments(opt.args, opt.repoID); err != nil {
		return []error{err}
	}
	var errorlist = []error{}
	for _, item := range opt.args {
		var list = add.addRepositoryToGroup(db, opt.repoID, opt.group, item)
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
