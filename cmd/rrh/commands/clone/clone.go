package clone

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize/english"
	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
	"github.com/tamada/rrh/cmd/rrh/commands/common"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clone <REMOTE_REPOs...>",
		Short: "run \"git clone\" and register it to a group",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return common.PerformRrhCommand(c, args, performClone)
		},
	}
	flags := cmd.Flags()
	flags.StringSliceVarP(&cloneOpts.groups, "groups", "g", []string{}, "specify the groups of the cloned repositories")
	flags.StringVarP(&cloneOpts.directory, "directory", "d", ".", "specify the destination directory")
	return cmd
}

var cloneOpts = &cloneOptions{}

type cloneOptions struct {
	groups    []string
	directory string
}

func updateGroups(config *rrh.Config) {
	if len(cloneOpts.groups) == 0 {
		cloneOpts.groups = []string{config.GetValue(rrh.DefaultGroupName)}
	}
}

func printResult(c *cobra.Command, count int) error {
	groupsString := fmt.Sprintf("%s: %s", english.PluralWord(len(cloneOpts.groups), "groups", ""), strings.Join(cloneOpts.groups, ", "))
	switch count {
	case 1:
		c.Printf("a repository cloned into %s and registered to %s\n", cloneOpts.directory, groupsString)
	default:
		c.Printf("%d repositories cloned into %s and registered to %s\n", count, cloneOpts.directory, groupsString)
	}
	return nil
}

func registerPath(db *rrh.Database, dest string, repoID string) (*rrh.Repository, error) {
	var path, err = filepath.Abs(dest)
	if err != nil {
		return nil, err
	}
	var remotes, err2 = rrh.FindRemotes(path)
	if err2 != nil {
		return nil, err2
	}
	var repo, err3 = db.CreateRepository(repoID, path, "", remotes)
	if err3 != nil {
		return nil, err3
	}
	return repo, nil
}

func toDir(db *rrh.Database, URL string, dest string, repoID string) (*rrh.Repository, error) {
	// clone.printIfVerbose(fmt.Sprintf("git clone %s %s (%s)", URL, dest, repoID))
	fmt.Printf("git clone %s %s (%s)\n", URL, dest, repoID)
	var cmd = exec.Command("git", "clone", URL, dest)
	var err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("%s: clone error (%s)", URL, err.Error())
	}
	return registerPath(db, dest, repoID)
}

func isExistDir(path string) bool {
	abs, err := filepath.Abs(path)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	stat, err := os.Stat(abs)
	return !os.IsNotExist(err) && stat.IsDir()
}

func doClone(args []string, db *rrh.Database) (int, error) {
	if len(args) == 1 {
		return 1, doCloneARepository(db, args[0])
	}
	el := common.NewErrorList()
	count := 0
	for _, argument := range args {
		err := doCloneEachRepository(db, argument)
		el = el.Append(err)
		if err == nil {
			count = count + 1
		}
	}
	return count, el.NilOrThis()
}

func relateTo(db *rrh.Database, groupIDs []string, repoID string) error {
	el := common.NewErrorList()
	for _, group := range groupIDs {
		_, err := db.AutoCreateGroup(group, "", false)
		el = el.Append(err)
		if err == nil {
			db.Relate(group, repoID)
		}
	}
	return el.NilOrThis()
}

func doCloneEachRepository(db *rrh.Database, URL string) error {
	id := findIDFromURL(URL)
	path := filepath.Clean(filepath.Join(cloneOpts.directory, id))
	_, err := toDir(db, URL, path, id)
	if err != nil {
		return err
	}
	return relateTo(db, cloneOpts.groups, id)
}

func doCloneARepository(db *rrh.Database, URL string) error {
	var id, path string

	if isExistDir(cloneOpts.directory) {
		id = findIDFromURL(URL)
		path = filepath.Join(cloneOpts.directory, id)
	} else {
		_, newid := filepath.Split(cloneOpts.directory)
		path = cloneOpts.directory
		id = newid
	}
	_, err := toDir(db, URL, path, id)
	if err != nil {
		return err
	}
	return relateTo(db, cloneOpts.groups, id)
}

func findIDFromURL(URL string) string {
	var _, dir = path.Split(URL)
	if strings.HasSuffix(dir, ".git") {
		return strings.TrimSuffix(dir, ".git")
	}
	return dir
}

func performClone(c *cobra.Command, args []string, db *rrh.Database) error {
	updateGroups(db.Config)
	count, err := doClone(args, db)
	if err != nil {
		return err
	}
	db.StoreAndClose()
	return printResult(c, count)
}
