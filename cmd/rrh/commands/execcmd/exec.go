package execcmd

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
	"github.com/tamada/rrh/cmd/rrh/commands/utils"
	"github.com/tamada/rrh/common"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exec <COMMAND>",
		Short: "execute the specified command on the specified repositories",
		Args:  validateArguments,
		RunE: func(c *cobra.Command, args []string) error {
			return utils.PerformRrhCommand(c, args, performExec)
		},
	}
	flags := cmd.Flags()
	flags.StringSliceVarP(&execOpts.repositories, "repositories", "r", []string{}, "specify the target repositories")
	flags.StringSliceVarP(&execOpts.groups, "groups", "g", []string{}, "specify the target group")
	flags.BoolVarP(&execOpts.withoutHeader, "no-header", "H", false, "print without header")
	flags.BoolVarP(&execOpts.dryRunFlag, "dry-run", "D", false, "dry-run mode")
	return cmd
}

func validateArguments(c *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("some command should be specified")
	}
	if len(execOpts.repositories) == 0 && len(execOpts.groups) == 0 {
		return errors.New("either of repositories and groups options should be specified")
	}
	return nil
}

var execOpts = &execOptions{}

type execOptions struct {
	repositories  []string
	groups        []string
	withoutHeader bool
	dryRunFlag    bool
}

func findRelatedRepositories(groups []string, db *rrh.Database) []string {
	repos := []string{}
	for _, group := range groups {
		r := db.FindRelationsOfGroup(group)
		repos = append(repos, r...)
	}
	return eliminateDuplication(repos)
}

func eliminateDuplication(froms []string) []string {
	repos := []string{}
	for _, from := range froms {
		if !rrh.FindIn(from, repos) {
			repos = append(repos, from)
		}
	}
	return repos
}

func validateRepos(repos []string, db *rrh.Database) ([]string, error) {
	el := common.NewErrorList()
	results := []string{}
	for _, r := range repos {
		repo := db.FindRepository(r)
		if repo == nil {
			el.Append(fmt.Errorf("%s: repository not found", r))
		} else {
			results = append(results, r)
		}
	}
	return results, el.NilOrThis()
}

func findTargetRepositories(db *rrh.Database) ([]string, error) {
	repos := []string{}
	if len(execOpts.groups) > 0 {
		repos2 := findRelatedRepositories(execOpts.groups, db)
		repos = append(repos, repos2...)
	}
	if len(execOpts.repositories) > 0 {
		repos = append(repos, execOpts.repositories...)
	}
	return validateRepos(repos, db)
}

func execute(c *cobra.Command, repo *rrh.Repository, args []string) error {
	if !execOpts.withoutHeader {
		c.Printf("----- %s (%s) -----\n", repo.ID, repo.Path)
	}
	if execOpts.dryRunFlag {
		c.Printf("%s: %v\n", repo.Path, strings.Join(args, " "))
		return nil
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = repo.Path
	cmd.Stdout = c.OutOrStdout()
	cmd.Stdin = c.InOrStdin()
	cmd.Stderr = c.ErrOrStderr()
	return cmd.Run()
}

func performExec(c *cobra.Command, args []string, db *rrh.Database) error {
	repositories, err := findTargetRepositories(db)
	if err != nil {
		return err
	}
	el := common.NewErrorList()
	for _, repository := range repositories {
		repo := db.FindRepository(repository)
		el.Append(execute(c, repo, args))
	}
	return el.NilOrThis()
}
