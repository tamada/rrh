package move

import (
	"fmt"

	"github.com/tamada/rrh/common"
)

func createGroupIfNeeded(db *common.Database, groupName string, description string) error {
	if !db.HasGroup(groupName) {
		if db.Config.GetValue(common.RrhAutoCreateGroup) != "true" {
			return fmt.Errorf("%s: group not found", groupName)
		}
		db.CreateGroup(groupName, description)
	}
	return nil
}

func (mv *MoveCommand) moveRepositoryToRepository(db *common.Database, from target, to target) error {
	if from.repositoryName != to.repositoryName {
		return fmt.Errorf("repository name did not match: %s, %s", from.original, to.original)
	}
	if err := createGroupIfNeeded(db, to.groupName, ""); err != nil {
		return err
	}
	if from.targetType == GroupAndRepoType {
		db.Unrelate(from.groupName, from.repositoryName)
		mv.Options.printIfNeeded(fmt.Sprintf("unrelate group %s and repository %s", from.groupName, from.repositoryName))
	}
	db.Relate(to.groupName, to.repositoryName)
	mv.Options.printIfNeeded(fmt.Sprintf("relate group %s and repository %s", to.groupName, to.repositoryName))
	return nil
}

func (mv *MoveCommand) moveRepositoryToGroup(db *common.Database, from target, to target) error {
	if to.targetType == GroupType || to.targetType == GroupOrRepoType {
		if err := createGroupIfNeeded(db, to.original, ""); err != nil {
			return err
		}
	}
	if from.targetType == GroupAndRepoType {
		db.Unrelate(from.groupName, from.repositoryName)
	}
	db.Relate(to.original, from.repositoryName)
	return nil
}

func isFailImmediately(config *common.Config) bool {
	return config.GetValue(common.RrhOnError) == common.FailImmediately
}

func (mv *MoveCommand) moveRepositoriesToGroup(db *common.Database, froms []target, to target) []error {
	var list = []error{}
	for _, from := range froms {
		var err = mv.moveRepositoryToGroup(db, from, to)
		if err != nil {
			if isFailImmediately(db.Config) {
				return []error{err}
			}
			list = append(list, err)
		}
	}
	return list
}

func (mv *MoveCommand) moveGroupsToGroup(db *common.Database, froms []target, to target) []error {
	var list = []error{}
	for _, from := range froms {
		var errs = mv.moveGroupToGroup(db, from, to)
		if len(errs) != 0 {
			if isFailImmediately(db.Config) {
				return errs
			}
			list = append(list, errs...)
		}
	}
	return list
}

func (mv *MoveCommand) moveGroupToGroup(db *common.Database, from target, to target) []error {
	if err := createGroupIfNeeded(db, to.groupName, ""); err != nil {
		return []error{err}
	}
	var repos = db.FindRelationsOfGroup(from.groupName)
	for _, repo := range repos {
		db.Unrelate(from.groupName, repo)
		mv.Options.printIfNeeded(fmt.Sprintf("unrelate group %s and repository %s", from.groupName, repo))
		db.Relate(to.groupName, repo)
		mv.Options.printIfNeeded(fmt.Sprintf("relate group %s and repository %s", to.groupName, repo))
	}
	return []error{}
}
