package group

import (
	"fmt"
	"strings"

	"github.com/tamada/rrh/common"
)

/*
GroupResult represents a group to show as the result.
*/
type GroupResult struct {
	Name        string
	Description string
	Repos       []string
}

func (group *groupListCommand) listGroups(db *common.Database, listOptions *listOptions) []GroupResult {
	var results = []GroupResult{}
	for _, group := range db.Groups {
		var result = GroupResult{group.Name, group.Description, []string{}}
		for _, relation := range db.Relations {
			if relation.GroupName == group.Name {
				result.Repos = append(result.Repos, relation.RepositoryID)
			}
		}
		results = append(results, result)
	}
	return results
}

func trueOrFalse(flag string) bool {
	var flagString = strings.ToLower(flag)
	if flagString == "true" {
		return true
	}
	return false
}

func (group *groupAddCommand) addGroups(db *common.Database, options *addOptions) error {
	for _, groupName := range options.args {
		var flag = trueOrFalse(options.omit)
		var _, err = db.CreateGroup(groupName, options.desc, flag)
		if err != nil {
			return err
		}
	}
	return nil
}

func (grc *groupRemoveCommand) removeGroupsImpl(db *common.Database, groupName string) error {
	if grc.Options.force {
		db.ForceDeleteGroup(groupName)
		grc.printIfVerbose(fmt.Sprintf("%s: group removed", groupName))
	} else if db.ContainsCount(groupName) == 0 {
		db.DeleteGroup(groupName)
		grc.printIfVerbose(fmt.Sprintf("%s: group removed", groupName))
	} else {
		return fmt.Errorf("%s: cannot remove group. the group has relations", groupName)
	}
	return nil
}

func (grc *groupRemoveCommand) removeGroups(db *common.Database) error {
	for _, groupName := range grc.Options.args {
		if !db.HasGroup(groupName) || !grc.Inquiry(groupName) {
			return nil
		}
		if err := grc.removeGroupsImpl(db, groupName); err != nil {
			return err
		}
	}
	return nil
}

func createNewGroup(opt *updateOptions, prevGroup *common.Group) common.Group {
	var newGroup = common.Group{Name: opt.newName, Description: opt.desc, OmitList: strings.ToLower(opt.omitList) == "true"}
	if opt.desc == "" {
		newGroup.Description = prevGroup.Description
	}
	if opt.newName == "" {
		newGroup.Name = prevGroup.Name
	}
	if opt.omitList == "" {
		newGroup.OmitList = prevGroup.OmitList
	}
	return newGroup
}

func (group *groupUpdateCommand) updateGroup(db *common.Database, opt *updateOptions) error {
	if !db.HasGroup(opt.target) {
		return fmt.Errorf("%s: group not found", opt.target)
	}
	var newGroup = createNewGroup(opt, db.FindGroup(opt.target))
	if !db.UpdateGroup(opt.target, newGroup) {
		return fmt.Errorf("%s: failed to update to {%s, %s, %s}", opt.target, opt.newName, opt.desc, opt.omitList)
	}
	return nil
}
