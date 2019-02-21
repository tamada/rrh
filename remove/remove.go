package remove

import (
	"fmt"
	"strings"

	"github.com/tamadalab/rrh/common"
)

func (rm *RemoveCommand) executeRemoveGroup(db *common.Database, groupName string, options *removeOptions) error {
	if options.inquiry && !common.IsInputYes(fmt.Sprintf("%s: Remove group? [yN]> ", groupName)) {
		rm.printIfVerbose(fmt.Sprintf("%s: group do not removed", groupName), options)
		return nil
	}
	var group = db.FindGroup(groupName)
	if group == nil {
		return fmt.Errorf("%s: group not found", groupName)
	}
	if !options.recursive && len(group.Items) > 0 {
		return fmt.Errorf("%s: cannot remove, it contains %d repository(es)", group.Name, len(group.Items))
	}

	for i, g := range db.Groups {
		if g.Name == group.Name {
			db.Groups[i].Items = []string{}
		}
	}
	var err = db.DeleteGroup(groupName)
	if err == nil {
		rm.printIfVerbose(fmt.Sprintf("%s: group removed", group.Name), options)
		return nil
	}
	return err
}

func (rm *RemoveCommand) executeRemoveRepository(db *common.Database, repoID string, options *removeOptions) error {
	if options.inquiry && !common.IsInputYes(fmt.Sprintf("%s: Remove repository? [yN]> ", repoID)) {
		rm.printIfVerbose(fmt.Sprintf("%s: repository do not removed", repoID), options)
		return nil
	}
	if err := db.DeleteRepository(repoID); err != nil {
		return err
	}
	rm.printIfVerbose(fmt.Sprintf("%s: repository removed", repoID), options)
	return nil
}

func (rm *RemoveCommand) executeRemoveFromGroup(db *common.Database, groupName string, repoID string, options *removeOptions) error {
	var err = db.Unrelate(groupName, repoID)
	if err == nil {
		rm.printIfVerbose(fmt.Sprintf("%s: removed from group %s", repoID, groupName), options)
		return nil
	}
	return err
}

func (rm *RemoveCommand) executeRemove(db *common.Database, target string, options *removeOptions) error {
	var data = strings.Split(target, "/")
	if len(data) == 2 {
		return rm.executeRemoveFromGroup(db, data[0], data[1], options)
	}
	var repoFlag = db.HasRepository(target)
	var groupFlag = db.HasGroup(target)
	if repoFlag && groupFlag {
		return fmt.Errorf("%s: exists in repositories and groups", target)
	}
	if repoFlag {
		return rm.executeRemoveRepository(db, target, options)
	}
	if groupFlag {
		return rm.executeRemoveGroup(db, target, options)
	}
	return fmt.Errorf("%s: not found in repositories and groups", target)
}
