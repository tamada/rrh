package add

import (
	"errors"
	"fmt"
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
	flags.StringSliceVarP(&addOpts.groups, "group", "g", []string{"no-group"}, "group for the repositories")
	flags.StringVarP(&addOpts.repoId, "repository-id", "r", "", "specifies the repository id. Specifying this option fails on multiple arguments")
	flags.BoolVarP(&addOpts.dryRunFlag, "dry-run", "D", false, "dry-run mode")
	return cmd
}

func validateArgs(c *cobra.Command, args []string) error {
	repo, _ := c.Flags().GetString("repository-id")
	if repo != "" && len(args) > 1 {
		return errors.New("too many arguments in specifying repository-id")
	}
	if len(args) == 0 {
		return errors.New("too few arguments")
	}
	return validateGitDirs(args)
}

func validateGitDirs(args []string) error {
	messages := common.NewErrorList()
	for _, arg := range args {
		absPath, _ := filepath.Abs(arg)
		err := rrh.IsExistAndGitRepository(absPath, arg)
		messages.Append(err)
	}
	if messages.IsNil() {
		return nil
	}
	return messages
}

func perform(c *cobra.Command, args []string, db *rrh.Database) error {
	el := createGroups(db, addOpts.groups)
	if el.IsErr() {
		return el
	}
	for _, targetPath := range args {
		err := addRepositoryToGroup(db, targetPath, addOpts.groups)
		el.Append(err)
	}
	if el.IsErr() {
		return el
	}
	if !addOpts.dryRunFlag {
		db.StoreAndClose()
	}
	return nil
}

func createGroups(db *rrh.Database, groups []string) common.ErrorList {
	el := common.NewErrorList()
	for _, groupName := range groups {
		_, err := db.AutoCreateGroup(groupName, "", false)
		el.Append(err)
	}
	return el
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
	var absPath, _ = filepath.Abs(path)
	var id = findIDFromPath(addOpts.repoId, absPath)
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
