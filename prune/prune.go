package prune

import (
	"os"

	"github.com/tamada/rrh/common"
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
