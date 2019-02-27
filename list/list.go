package list

import (
	"fmt"

	"github.com/tamada/rrh/common"
)

type Repo struct {
	Name    string
	Path    string
	Remotes []common.Remote
}

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
	for _, repoName := range group.Items {
		var repo = db.FindRepository(repoName)
		if repo == nil {
			return nil, fmt.Errorf("%s: repository not found", repoName)
		}
		repos = append(repos, Repo{repo.ID, repo.Path, repo.Remotes})
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
