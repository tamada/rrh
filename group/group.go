package group

import (
	"fmt"

	"github.com/tamadalab/rrh/common"
)

/*
GroupResult represents a group to show as the result.
*/
type GroupResult struct {
	Name        string
	Description string
	Repos       []string
}

func (group *groupListCommand) listGroups(db *common.Database, listOptions *listOptions) ([]GroupResult, error) {
	var results = []GroupResult{}
	for _, group := range db.Groups {
		results = append(results, GroupResult{group.Name, group.Description, group.Items})
	}
	return results, nil
}

func (group *groupAddCommand) addGroups(db *common.Database, options *addOptions) error {
	for _, group := range options.args {
		var _, err = db.CreateGroup(group, options.desc)
		if err != nil {
			return err
		}
	}
	return nil
}

func (group *groupRemoveCommand) removeGroups(db *common.Database, options *removeOptions) error {
	for _, groupName := range options.args {
		if db.HasGroup(groupName) && options.Inquiry(groupName) {
			var group = db.FindGroup(groupName)
			if len(group.Items) == 0 || options.force {
				db.DeleteGroup(groupName)
				options.printIfVerbose(fmt.Sprintf("%s: group removed", groupName))
			} else {
				return fmt.Errorf("%s: cannot remove group. the group has relations", groupName)
			}
		}
	}
	return nil
}

func (group *groupUpdateCommand) updateGroup(db *common.Database, opt *updateOptions) error {
	if !db.HasGroup(opt.target) {
		return fmt.Errorf("%s: group not found", opt.target)
	}
	if !db.UpdateGroup(opt.target, opt.newName, opt.desc) {
		return fmt.Errorf("%s: failed to update to {%s, %s}", opt.target, opt.newName, opt.desc)
	}
	return nil
}
