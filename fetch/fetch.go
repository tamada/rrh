package fetch

import (
	"fmt"
	"os/exec"

	"github.com/tamada/rrh/common"
)

/*
Progress represents a fetching progress.
*/
type Progress struct {
	current int
	total   int
}

func (progress *Progress) String() string {
	return fmt.Sprintf("%3d/%3d", progress.current, progress.total)
}

func (progress *Progress) increment() {
	progress.current++
}

/*
DoFetch exec fetch operation of git.
Currently, fetch is conducted by the system call.
Ideally, fetch is performed by using go-git.
*/
func (fetch *Command) DoFetch(repo *common.Repository, group string, progress *Progress) error {
	var cmd = exec.Command("git", "fetch", fetch.options.remote)
	cmd.Dir = repo.Path
	progress.increment()
	fmt.Printf("%s fetching %s/%s....", progress, group, repo.ID)
	var output, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("%s,%s,%s", group, repo.ID, err.Error())
	}
	fmt.Printf("done\n%s", output)
	return nil
}

func (fetch *Command) fetchRepository(db *common.Database, groupName string, repoID string, progress *Progress) error {
	var repository = db.FindRepository(repoID)
	if repository == nil {
		return fmt.Errorf("%s,%s: repository not found", groupName, repoID)
	}
	return fetch.DoFetch(repository, groupName, progress)
}

/*
FetchGroup performs `git fetch` command in the repositories belonging in the specified group.
*/
func (fetch *Command) FetchGroup(db *common.Database, groupName string, progress *Progress) []error {
	var list = []error{}
	var group = db.FindGroup(groupName)
	if group == nil {
		return []error{fmt.Errorf("%s: group not found", groupName)}
	}
	for _, relation := range db.Relations {
		var err = fetch.executeFetch(db, groupName, relation, progress)
		if err == nil {
			continue
		}
		if db.Config.GetValue(common.RrhOnError) == common.FailImmediately {
			return []error{err}
		}
		list = append(list, err)
	}
	return list
}

func (fetch *Command) executeFetch(db *common.Database, groupName string, relation common.Relation, progress *Progress) error {
	if relation.GroupName == groupName {
		var err = fetch.fetchRepository(db, groupName, relation.RepositoryID, progress)
		if err != nil {
			return err
		}
	}
	return nil
}
