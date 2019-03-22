package fetch

import (
	"fmt"
	"os/exec"

	"github.com/tamada/rrh/common"
)

/*
DoFetch exec fetch operation of git.
Currently, fetch is conducted by the system call.
Ideally, fetch is performed by using go-git.
*/
func (fetch *FetchCommand) DoFetch(repo *common.Repository, group string, options *FetchOptions, config *common.Config) error {
	var cmd = exec.Command("git", "fetch", options.remote)
	cmd.Dir = common.ToAbsolutePath(repo.Path, config)
	fmt.Printf("fetching %s,%s....", group, repo.ID)
	var output, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("%s,%s,%s", group, repo.ID, err.Error())
	}
	fmt.Printf("done\n%s", output)
	return nil
}

func (fetch *FetchCommand) fetchRepository(db *common.Database, groupName string, repoID string, options *FetchOptions) error {
	var repository = db.FindRepository(repoID)
	if repository == nil {
		return fmt.Errorf("%s,%s: repository not found", groupName, repoID)
	}
	return fetch.DoFetch(repository, groupName, options, db.Config)
}

func (fetch *FetchCommand) FetchGroup(db *common.Database, groupName string, options *FetchOptions) []error {
	var list = []error{}
	var group = db.FindGroup(groupName)
	if group == nil {
		return []error{fmt.Errorf("%s: group not found", groupName)}
	}
	for _, relation := range db.Relations {
		var err = fetch.executeFetch(db, groupName, relation, options)
		if err == nil {
			continue
		}
		if db.Config.GetValue(common.RrhOnError) == common.FailImmediately {
			return []error{err}
		} else {
			list = append(list, err)
		}
	}
	return list
}

func (fetch *FetchCommand) executeFetch(db *common.Database, groupName string, relation common.Relation, options *FetchOptions) error {
	if relation.GroupName == groupName {
		var err = fetch.fetchRepository(db, groupName, relation.RepositoryID, options)
		if err != nil {
			return err
		}
	}
	return nil
}
