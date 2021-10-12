package importcmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
	"github.com/tamada/rrh/cmd/rrh/commands/common"
)

type importOptions struct {
	autoClone bool
	overwrite bool
	dryRun    bool
}

var importOpts = &importOptions{}

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "import the given database",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return common.PerformRrhCommand(c, args, perform)
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&importOpts.autoClone, "auto-clone", "", false, "clone the repository, if paths do not exist")
	flags.BoolVarP(&importOpts.overwrite, "overwrite", "", false, "replace the local RRH database to the given database")
	flags.BoolVarP(&importOpts.dryRun, "dry-run", "D", false, "dry-run mode")
	return cmd
}

func perform(c *cobra.Command, args []string, db *rrh.Database) error {
	if importOpts.overwrite {
		eraseDatabase(db)
	}
	var db2, err = readNewDB(args[0], db.Config)
	if err != nil {
		return err
	}
	return copyDB(db2, db)
}

func eraseDatabase(db *rrh.Database) {
	db.Groups = []*rrh.Group{}
	db.Repositories = []*rrh.Repository{}
	db.Relations = []*rrh.Relation{}
}

func readNewDB(path string, config *rrh.Config) (*rrh.Database, error) {
	var db = rrh.Database{Timestamp: rrh.Now(), Repositories: []*rrh.Repository{}, Groups: []*rrh.Group{}, Relations: []*rrh.Relation{}, Config: config}
	var bytes, err = ioutil.ReadFile(path)
	if err != nil {
		return &db, nil
	}
	var homeReplacedString = replaceHome(bytes)

	if err := json.Unmarshal([]byte(homeReplacedString), &db); err != nil {
		return nil, err
	}
	return &db, nil
}

func copyDB(from *rrh.Database, to *rrh.Database) common.ErrorList {
	var errs = []error{}
	var errs1 = copyGroups(from, to)
	var errs2 = copyRepositories(from, to)
	var errs3 = copyRelations(from, to)
	errs = append(errs, errs1...)
	errs = append(errs, errs2...)
	return append(errs, errs3...)
}

func copyGroup(group *rrh.Group, to *rrh.Database) common.ErrorList {
	var list = []error{}
	if to.HasGroup(group.Name) {
		var successFlag = to.UpdateGroup(group.Name, group)
		if !successFlag {
			list = append(list, fmt.Errorf("%s: update failed", group.Name))
		}
	} else {
		var _, err = to.CreateGroup(group.Name, group.Description, group.OmitList)
		if err != nil {
			list = append(list, err)
		}
	}
	return list
}

func copyGroups(from *rrh.Database, to *rrh.Database) common.ErrorList {
	list := common.NewErrorList()
	for _, group := range from.Groups {
		var errs = copyGroup(group, to)
		list = list.Append(errs)
	}
	return list
}

func findOrigin(remotes []*rrh.Remote) *rrh.Remote {
	for _, remote := range remotes {
		if remote.Name == "origin" {
			return remote
		}
	}
	return remotes[0]
}

func doClone(repository *rrh.Repository, remote *rrh.Remote) error {
	var cmd = exec.Command("git", "clone", remote.URL, repository.Path)
	var err = cmd.Run()
	if err != nil {
		return fmt.Errorf("%s: clone error (%s)", remote.URL, err.Error())
	}
	return nil
}

func cloneRepository(repository *rrh.Repository) error {
	if len(repository.Remotes) == 0 {
		return fmt.Errorf("%s: could not clone, did not have remotes", repository.ID)
	}
	var remote = findOrigin(repository.Remotes)
	var err = doClone(repository, remote)
	return err
}

func cloneIfNeeded(repository *rrh.Repository) error {
	if !importOpts.autoClone {
		return fmt.Errorf("%s: repository path did not exist at %s", repository.ID, repository.Path)
	}
	cloneRepository(repository)
	return nil
}

func copyRepository(repository *rrh.Repository, to *rrh.Database) common.ErrorList {
	if to.HasRepository(repository.ID) {
		return []error{}
	}
	var _, err = os.Stat(repository.Path)
	if err != nil {
		var err1 = cloneIfNeeded(repository)
		if err1 != nil {
			return []error{err1}
		}
	}
	return copyRepositoryImpl(repository, to)
}

func copyRepositoryImpl(repository *rrh.Repository, to *rrh.Database) common.ErrorList {
	if err := rrh.IsExistAndGitRepository(repository.Path, repository.ID); err != nil {
		return []error{err}
	}
	to.CreateRepository(repository.ID, repository.Path, repository.Description, repository.Remotes)
	return []error{}
}

func copyRepositories(from *rrh.Database, to *rrh.Database) []error {
	var list = common.NewErrorList()
	for _, repository := range from.Repositories {
		var errs = copyRepository(repository, to)
		list = list.Append(errs)
	}
	return list
}

func copyRelation(rel *rrh.Relation, to *rrh.Database) common.ErrorList {
	var list = []error{}
	if to.HasGroup(rel.GroupName) && to.HasRepository(rel.RepositoryID) {
		to.Relate(rel.GroupName, rel.RepositoryID)
	} else {
		list = append(list, fmt.Errorf("group %s and repository %s: could not relate", rel.GroupName, rel.RepositoryID))
	}
	return list
}

func copyRelations(from *rrh.Database, to *rrh.Database) common.ErrorList {
	var list = common.NewErrorList()
	for _, rel := range from.Relations {
		var errs = copyRelation(rel, to)
		list = list.Append(errs)
	}
	return list
}

func replaceHome(bytes []byte) string {
	var home, err = os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Warning: could not get home directory")
	}
	var absPath, _ = filepath.Abs(home)
	return strings.Replace(string(bytes), "${HOME}", absPath, -1)
}
