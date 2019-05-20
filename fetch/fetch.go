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

func foundIn(relation common.Relation, list []common.Relation) bool {
	for _, r := range list {
		if r.GroupName == relation.GroupName &&
			r.RepositoryID == relation.RepositoryID {
			return true
		}
	}
	return false
}

func eliminateDuplicate(relations []common.Relation) []common.Relation {
	var result = []common.Relation{}
	for _, relation := range relations {
		if !foundIn(relation, result) {
			result = append(result, relation)
		}
	}
	return result
}

func toRelations(groupName string, repoNames []string) []common.Relation {
	var result = []common.Relation{}
	for _, repo := range repoNames {
		result = append(result, common.Relation{RepositoryID: repo, GroupName: groupName})
	}
	return result
}

/*
FindTargets returns the instances of Relation objects with given groupNames.
*/
func (fetch *Command) FindTargets(db *common.Database, groupNames []string) []common.Relation {
	var result = []common.Relation{}
	for _, groupName := range groupNames {
		var repos = db.FindRelationsOfGroup(groupName)
		var relations = toRelations(groupName, repos)
		result = append(result, relations...)
	}
	return eliminateDuplicate(result)
}

/*
DoFetch exec fetch operation of git.
Currently, fetch is conducted by the system call.
Ideally, fetch is performed by using go-git.
*/
func (fetch *Command) DoFetch(repo *common.Repository, relation *common.Relation, progress *Progress) error {
	var cmd = exec.Command("git", "fetch", fetch.options.remote)
	cmd.Dir = repo.Path
	progress.increment()
	fmt.Printf("%s fetching %s....", progress, relation)
	var output, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("%s,%s", relation, err.Error())
	}
	fmt.Printf("done\n%s", output)
	return nil
}

/*
FetchRepository execute `git fetch` on the given repository.
*/
func (fetch *Command) FetchRepository(db *common.Database, relation *common.Relation, progress *Progress) error {
	var repository = db.FindRepository(relation.RepositoryID)
	if repository == nil {
		return fmt.Errorf("%s: repository not found", relation)
	}
	return fetch.DoFetch(repository, relation, progress)
}
