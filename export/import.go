package export

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/tamada/rrh/common"
)

func (command *ImportCommand) readNewDB(path string, config *common.Config) (*common.Database, error) {
	var db = common.Database{Timestamp: common.Now(), Repositories: []common.Repository{}, Groups: []common.Group{}, Relations: []common.Relation{}, Config: config}
	var bytes, err = ioutil.ReadFile(path)
	if err != nil {
		return &db, nil
	}
	var homeReplacedString = replaceHome(bytes)

	if err := json.Unmarshal([]byte(homeReplacedString), &db); err != nil {
		return nil, err
	}
	return &db, nil
}

func (command *ImportCommand) copyDB(from *common.Database, to *common.Database) []error {
	var errs = []error{}
	var errs1 = command.copyGroups(from, to)
	var errs2 = command.copyRepositories(from, to)
	var errs3 = command.copyRelations(from, to)
	errs = append(errs, errs1...)
	errs = append(errs, errs2...)
	return append(errs, errs3...)
}

func (command *ImportCommand) copyGroups(from *common.Database, to *common.Database) []error {
	var list = []error{}
	for _, group := range from.Groups {
		if to.HasGroup(group.Name) {
			var successFlag = to.UpdateGroup(group.Name, group)
			if !successFlag {
				list = append(list, fmt.Errorf("%s: update failed", group.Name))
			}
		} else {
			var _, err = to.CreateGroup(group.Name, group.Description, group.OmitList)
			if err != nil {
				list = append(list, err)
			}
			command.options.printIfNeeded(fmt.Sprintf("%s: create group", group.Name))
		}
	}
	return list
}

func isFailImmediately(config *common.Config) bool {
	return config.GetValue(common.RrhOnError) == common.FailImmediately
}

func findOrigin(remotes []common.Remote) common.Remote {
	for _, remote := range remotes {
		if remote.Name == "origin" {
			return remote
		}
	}
	return remotes[0]
}

func (command *ImportCommand) doClone(repository common.Repository, remote common.Remote) error {
	var cmd = exec.Command("git", "clone", remote.URL, repository.Path)
	var err = cmd.Run()
	if err != nil {
		return fmt.Errorf("%s: clone error (%s)", remote.URL, err.Error())
	}
	return nil
}

func (command *ImportCommand) cloneRepository(repository common.Repository) error {
	if len(repository.Remotes) == 0 {
		return fmt.Errorf("%s: could not clone, did not have remotes", repository.ID)
	}
	var remote = findOrigin(repository.Remotes)
	return command.doClone(repository, remote)
}

func (command *ImportCommand) copyRepositories(from *common.Database, to *common.Database) []error {
	var list = []error{}
	for _, repository := range from.Repositories {
		if to.HasRepository(repository.ID) {
			continue
		}
		var _, err = os.Stat(repository.Path)
		if err != nil {
			if command.options.autoClone {
				command.cloneRepository(repository)
			} else {
				list = append(list, fmt.Errorf("%s: repository path did not exist at %s", repository.ID, repository.Path))
				if isFailImmediately(from.Config) {
					return list
				}
				continue
			}
		}
		if err := common.IsExistAndGitRepository(repository.Path, repository.ID); err != nil {
			list = append(list, err)
			if isFailImmediately(from.Config) {
				return list
			}
		} else {
			to.CreateRepository(repository.ID, repository.Path, repository.Remotes)
			command.options.printIfNeeded(fmt.Sprintf("%s: create repository", repository.ID))
		}
	}
	return list
}

func (command *ImportCommand) copyRelations(from *common.Database, to *common.Database) []error {
	var list = []error{}
	for _, rel := range from.Relations {
		if to.HasGroup(rel.GroupName) && to.HasRepository(rel.RepositoryID) {
			to.Relate(rel.GroupName, rel.RepositoryID)
			command.options.printIfNeeded(fmt.Sprintf("%s, %s: create relation", rel.GroupName, rel.RepositoryID))
		} else {
			list = append(list, fmt.Errorf("group %s and repository %s: could not relate", rel.GroupName, rel.RepositoryID))
			if isFailImmediately(to.Config) {
				return list
			}
		}
	}
	return list
}

func replaceHome(bytes []byte) string {
	var home, err = homedir.Dir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Warning: could not get home directory")
	}
	var absPath, _ = filepath.Abs(home)
	return strings.Replace(string(bytes), "${HOME}", absPath, -1)
}
