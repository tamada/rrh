package group

import (
	"fmt"

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

func (grc *groupRemoveCommand) removeGroupsImpl(db *common.Database, groupName string) error {
	var group = db.FindGroup(groupName)
	if grc.Options.force {
		db.ForceDeleteGroup(groupName)
		grc.printIfVerbose(fmt.Sprintf("%s: group removed", groupName))
	} else if len(group.Items) == 0 {
		db.DeleteGroup(groupName)
		grc.printIfVerbose(fmt.Sprintf("%s: group removed", groupName))
	} else {
		return fmt.Errorf("%s: cannot remove group. the group has relations", groupName)
	}
	return nil
}

func (grc *groupRemoveCommand) removeGroups(db *common.Database) error {
	for _, groupName := range grc.Options.args {
		if db.HasGroup(groupName) && grc.Inquiry(groupName) {
			if err := grc.removeGroupsImpl(db, groupName); err != nil {
				return err
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
