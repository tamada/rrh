package remove

import (
	"fmt"
	"strings"

	"github.com/tamada/rrh/common"
)

func (rm *RemoveCommand) executeRemoveGroup(db *common.Database, groupName string) error {
	var group = db.FindGroup(groupName)
	if group == nil {
		return fmt.Errorf("%s: group not found", groupName)
	}
	if rm.Options.inquiry && !common.IsInputYes(fmt.Sprintf("%s: Remove group? [yN]> ", groupName)) {
		rm.Options.printIfVerbose(fmt.Sprintf("%s: group do not removed", groupName))
		return nil
	}
	if !rm.Options.recursive && len(group.Items) > 0 {
		return fmt.Errorf("%s: cannot remove, it contains %d repository(es)", group.Name, len(group.Items))
	}

	for i, g := range db.Groups {
		if g.Name == group.Name {
			db.Groups[i].Items = []string{}
		}
	}
	var err = db.DeleteGroup(groupName)
	if err == nil {
		rm.Options.printIfVerbose(fmt.Sprintf("%s: group removed", group.Name))
		return nil
	}
	return err
}

func (rm *RemoveCommand) executeRemoveRepository(db *common.Database, repoID string) error {
	if !db.HasRepository(repoID) {
		return fmt.Errorf("%s: repository not found", repoID)
	}
	if rm.Options.inquiry && !common.IsInputYes(fmt.Sprintf("%s: Remove repository? [yN]> ", repoID)) {
		rm.Options.printIfVerbose(fmt.Sprintf("%s: repository do not removed", repoID))
		return nil
	}
	if err := db.DeleteRepository(repoID); err != nil {
		return err
	}
	rm.Options.printIfVerbose(fmt.Sprintf("%s: repository removed", repoID))
	return nil
}

func (rm *RemoveCommand) executeRemoveFromGroup(db *common.Database, groupName string, repoID string) error {
	db.Unrelate(groupName, repoID)
	rm.Options.printIfVerbose(fmt.Sprintf("%s: removed from group %s", repoID, groupName))
	return nil
}

func (rm *RemoveCommand) executeRemove(db *common.Database, target string) error {
	var data = strings.Split(target, "/")
	if len(data) == 2 {
		return rm.executeRemoveFromGroup(db, data[0], data[1])
	}
	var repoFlag = db.HasRepository(target)
	var groupFlag = db.HasGroup(target)
	if repoFlag && groupFlag {
		return fmt.Errorf("%s: exists in repositories and groups", target)
	}
	if repoFlag {
		return rm.executeRemoveRepository(db, target)
	}
	if groupFlag {
		return rm.executeRemoveGroup(db, target)
	}
	return fmt.Errorf("%s: not found in repositories and groups", target)
}
