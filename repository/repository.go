package repository

import (
	"fmt"

	"github.com/tamada/rrh/common"
)

func findAll(db *common.Database, args []string) ([]common.Repository, []error) {
	if len(args) > 0 {
		return findResults(db, args)
	}
	return db.Repositories, []error{}
}

func findResult(db *common.Database, arg string) (*common.Repository, error) {
	var repo = db.FindRepository(arg)
	if repo == nil {
		return nil, fmt.Errorf("%s: repository not found", arg)
	}
	return repo, nil
}

func findResults(db *common.Database, args []string) ([]common.Repository, []error) {
	var results = []common.Repository{}
	var errs = []error{}
	for _, arg := range args {
		var repo, err = findResult(db, arg)
		if err != nil {
			errs = append(errs, err)
			if db.Config.GetValue(common.RrhOnError) == common.FailImmediately {
				return []common.Repository{}, errs
			}
		} else {
			results = append(results, *repo)
		}
	}
	return results, errs
}

func (update *updateCommand) perform(db *common.Database, targetRepoID string) error {
	var repo = db.FindRepository(targetRepoID)
	if repo == nil {
		return fmt.Errorf("%s: repository not found", targetRepoID)
	}
	var newRepo = buildNewRepo(update.options, repo)
	if !db.UpdateRepository(targetRepoID, newRepo) {
		return fmt.Errorf("%s: repository update failed", targetRepoID)
	}
	return nil
}

func buildNewRepo(options *updateOptions, repo *common.Repository) common.Repository {
	var newRepo = common.Repository{ID: repo.ID, Path: repo.Path, Description: repo.Description}
	if options.description != "" {
		newRepo.Description = options.description
	}
	if options.newID != "" {
		newRepo.ID = options.newID
	}
	if options.newPath != "" {
		newRepo.Path = options.newPath
	}
	return newRepo
}
