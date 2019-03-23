package list

import (
	"fmt"

	"github.com/tamada/rrh/common"
)

/*
Repo represents the result for showing of repositories.
*/
type Repo struct {
	Name    string
	Path    string
	Remotes []common.Remote
}

/*
ListResult represents the result for showing.
*/
type ListResult struct {
	GroupName   string
	Description string
	Repos       []Repo
}

func (list *ListCommand) findList(db *common.Database, groupName string) (*ListResult, error) {
	var repos = []Repo{}
	var group = db.FindGroup(groupName)
	if group == nil {
		return nil, fmt.Errorf("%s: group not found", groupName)
	}
	for _, relation := range db.Relations {
		if relation.GroupName == groupName {
			var repo = db.FindRepository(relation.RepositoryID)
			if repo == nil {
				return nil, fmt.Errorf("%s: repository not found", relation.RepositoryID)
			}
			repos = append(repos, Repo{repo.ID, repo.Path, repo.Remotes})
		}
	}

	return &ListResult{group.Name, group.Description, repos}, nil
}

func (list *ListCommand) findAllGroupNames(db *common.Database) []string {
	var names = []string{}
	for _, group := range db.Groups {
		names = append(names, group.Name)
	}
	return names
}

/*
FindResults returns the result list of list command.
*/
func (list *ListCommand) FindResults(db *common.Database) ([]ListResult, error) {
	var groups = list.Options.args
	if len(groups) == 0 {
		groups = list.findAllGroupNames(db)
	}
	var results = []ListResult{}
	for _, group := range groups {
		var list, err = list.findList(db, group)
		if err != nil {
			return nil, err
		}
		results = append(results, *list)
	}
	return results, nil
}
