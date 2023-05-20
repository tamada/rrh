package add

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
	"github.com/tamada/rrh/cmd/rrh/commands/utils"
	"github.com/tamada/rrh/common"
)

type addOptions struct {
	groups     []string
	repoId     string
	dryRunFlag bool
}

var addOpts = &addOptions{groups: []string{}}

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <REPOs...>",
		Args:  validateArgs,
		Short: "add the given repositories to the rrh database",
		RunE: func(c *cobra.Command, args []string) error {
			return utils.PerformRrhCommand(c, args, perform)
		},
	}
	flags := cmd.Flags()
	defaultGroup := findDefaultGroupName()
	flags.StringSliceVarP(&addOpts.groups, "group", "g", []string{defaultGroup}, "group for the repositories")
	flags.StringVarP(&addOpts.repoId, "repository-id", "r", "", "specifies the repository id. Specifying this option fails on multiple arguments")
	flags.BoolVarP(&addOpts.dryRunFlag, "dry-run", "D", false, "dry-run mode")
	return cmd
}

func findDefaultGroupName() string {
	defaultGroupName := os.Getenv("RRH_DEFAULT_GROUP")
	if defaultGroupName == "" {
		defaultGroupName = "no-group"
	}
	return defaultGroupName
}

func validateArgs(c *cobra.Command, args []string) error {
	repo, _ := c.Flags().GetString("repository-id")
	if repo != "" && len(args) > 1 {
		return errors.New("too many arguments in specifying repository-id")
	}
	if len(args) == 0 {
		return errors.New("too few arguments")
	}
	return nil
	// return validateGitDirs(args)
}

func validateGitDirs(args []string) error {
	messages := common.NewErrorList()
	for _, arg := range args {
		absPath, _ := filepath.Abs(arg)
		err := rrh.IsExistAndGitRepository(absPath, arg)
		messages = messages.Append(err)
	}
	return messages.NilOrThis()
}

func perform(c *cobra.Command, args []string, db *rrh.Database) error {
	el := createGroups(db, addOpts.groups)
	if el.IsErr() {
		return el
	}
	for _, targetPath := range args {
		err := addRepositoryToGroup(db, targetPath, addOpts.groups)
		el = el.Append(err)
	}
	if !addOpts.dryRunFlag {
		db.StoreAndClose()
	}
	if el.IsErr() {
		return el
	}
	return nil
}

func createGroups(db *rrh.Database, groups []string) common.ErrorList {
	el := common.NewErrorList()
	if len(groups) == 0 && len(db.Groups) == 0 {
		el = el.Append(createDefaultGroupNameIfNeeded(db))
	}
	for _, groupName := range groups {
		_, err := db.AutoCreateGroup(groupName, "", false)
		el = el.Append(err)
	}
	return el
}

func createDefaultGroupNameIfNeeded(db *rrh.Database) error {
	defaultGroupName := findDefaultGroupName()
	if db.HasGroup(defaultGroupName) {
		return nil
	}
	_, err := db.CreateGroup(defaultGroupName, "default group", false)
	return err
}

func findIDFromPath(repoIdFromOpts string, absPath string) string {
	if repoIdFromOpts == "" {
		return filepath.Base(absPath)
	}
	return repoIdFromOpts
}

func isDuplicateRepositoryId(db *rrh.Database, repoId, path string) error {
	var repo = db.FindRepository(repoId)
	if repo != nil && repo.Path != path {
		return fmt.Errorf("%s: duplicate repository id", repoId)
	}
	return nil
}

func addRepositoryToGroup(db *rrh.Database, path string, groupNames []string) error {
	absPath, _ := filepath.Abs(path)
	if err := rrh.IsExistAndGitRepository(absPath, path); err != nil {
		return err
	}
	id := findIDFromPath(addOpts.repoId, absPath)
	if err1 := isDuplicateRepositoryId(db, id, absPath); err1 != nil {
		return err1
	}
	remotes, err2 := rrh.FindRemotes(absPath)
	if err2 != nil {
		return err2
	}
	db.CreateRepository(id, absPath, "", remotes)

	for _, groupName := range groupNames {
		err := db.Relate(groupName, id)
		if err != nil {
			return err
		}
	}
	return nil
}
