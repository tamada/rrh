package prune

import (
	"os"

	"github.com/tamadalab/rrh/common"
)

func (prune *PruneCommand) removeNotExistRepository(db *common.Database) int {
	var removeRepos = []string{}
	for _, repo := range db.Repositories {
		var path = common.ToAbsolutePath(repo.Path, db.Config)
		var _, err = os.Stat(path)
		if os.IsNotExist(err) {
			removeRepos = append(removeRepos, repo.ID)
		}
	}

	var count = 0
	for _, repo := range removeRepos {
		var err = db.DeleteRepository(repo)
		if err == nil {
			count++
		}
	}
	return count
}

func (prune *PruneCommand) countReposInGroups(db *common.Database) map[string]int {
	var repoFlags = map[string]int{}
	for _, repo := range db.Repositories {
		repoFlags[repo.ID] = 0
	}
	for _, group := range db.Groups {
		for _, item := range group.Items {
			repoFlags[item] = repoFlags[item] + 1
		}
	}
	return repoFlags
}

func (prune *PruneCommand) pruneGroup(db *common.Database) int {
	var newGroups = []common.Group{}
	var count = 0
	for _, group := range db.Groups {
		if len(group.Items) != 0 {
			newGroups = append(newGroups, group)
		} else {
			count++
		}
	}
	db.Groups = newGroups
	return count
}
