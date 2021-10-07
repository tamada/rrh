package repository

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
	"github.com/tamada/rrh/cmd/rrh/commands/common"
)

func newUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Args:  validateUpdate,
		Short: "update repository information",
		RunE: func(c *cobra.Command, args []string) error {
			return common.PerformRrhCommand(c, args, performUpdate)
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&updateOpts.dryRunMode, "dry-run", "D", false, "dry-run mode")
	flags.StringVarP(&updateOpts.newId, "id", "", "", "specify the new repository id")
	flags.StringVarP(&updateOpts.newPath, "path", "", "", "specify the new repository path")
	flags.StringSliceVarP(&updateOpts.appendGroups, "append-group", "", []string{}, "specify the appending groups")
	flags.StringSliceVarP(&updateOpts.groups, "group", "", []string{}, "replace the groups of the repository")
	flags.StringVarP(&updateOpts.newDescription, "description", "", "", "specify the new repository description")
	return cmd
}

var updateOpts = &updateOptions{}

type updateOptions struct {
	dryRunMode     bool
	newId          string
	newDescription string
	appendGroups   []string
	groups         []string
	newPath        string
}

func validateUpdate(c *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("arguments must be only one")
	}
	if updateOpts.newDescription == "" && updateOpts.newId == "" && updateOpts.newPath == "" {
		return fmt.Errorf("either option of id, description and path must be required")
	}
	if len(updateOpts.groups) != 0 && len(updateOpts.appendGroups) != 0 {
		return fmt.Errorf("only either option of group and append-group are available")
	}
	return nil
}

func performUpdate(c *cobra.Command, args []string, db *rrh.Database) error {
	repo := db.FindRepository(args[0])
	if repo == nil {
		return fmt.Errorf("%s: repository not found", args[0])
	}
	results, err := updateRepoInfo(repo, db)
	if err != nil {
		return err
	}
	showInfo(c, results)
	if !updateOpts.dryRunMode {
		return db.StoreAndClose()
	}
	return nil
}

func showInfo(c *cobra.Command, results [][]string) {
	table := tablewriter.NewWriter(c.OutOrStdout())
	table.SetBorder(false)
	table.SetNoWhiteSpace(true)
	table.SetCenterSeparator("")
	table.SetRowSeparator("")
	newLabel := "New"
	if updateOpts.dryRunMode {
		newLabel = "New (dry-run mode)"
	}

	table.SetHeader([]string{"", "Old", newLabel})
	table.SetTablePadding("   ")
	table.AppendBulk(results)
	table.Render()
}

func updateRepoInfo(repo *rrh.Repository, db *rrh.Database) ([][]string, error) {
	results := [][]string{}

	results = append(results, []string{"Repository ID", repo.ID, updateOpts.newId})
	if updateOpts.newId != "" {
		repo.ID = updateOpts.newId
	}
	results = append(results, []string{"Description", repo.Description, updateOpts.newDescription})
	if updateOpts.newDescription != "" {
		repo.Description = updateOpts.newDescription
	}
	if updateOpts.newPath != "" {
		abs, err := checkPath(updateOpts.newPath)
		if err != nil {
			return [][]string{}, err
		}
		results = append(results, []string{"Path", repo.Path, abs})
		repo.Path = abs
	} else {
		results = append(results, []string{"Path", repo.Path, ""})
	}
	oldGroups, newGroups, err := updateGroups(repo, db)
	if err != nil {
		return nil, err
	}
	results = append(results, []string{"Groups", strings.Join(oldGroups, ", "), strings.Join(newGroups, ", ")})
	return results, err
}

func removeAllGroups(repo *rrh.Repository, db *rrh.Database) error {
	groups := db.FindRelationsOfRepository(repo.ID)
	for _, group := range groups {
		db.Unrelate(group, repo.ID)
	}
	return nil
}

func updateGroups(repo *rrh.Repository, db *rrh.Database) ([]string, []string, error) {
	oldGroups := db.FindRelationsOfRepository(repo.ID)
	if len(updateOpts.appendGroups) == 0 && len(updateOpts.groups) == 0 {
		return oldGroups, []string{}, nil
	}
	groups := updateOpts.appendGroups
	if len(updateOpts.groups) > 0 {
		if err := removeAllGroups(repo, db); err != nil {
			return nil, nil, err
		}
		groups = updateOpts.groups
	}
	for _, group := range groups {
		db.Relate(group, repo.ID)
	}
	return oldGroups, db.FindRelationsOfRepository(repo.ID), nil
}

func checkPath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	if err := rrh.IsExistAndGitRepository(abs, path); err != nil {
		return "", err
	}
	return abs, nil
}
