package status

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/tamada/rrh/common"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

/*
StatusResult shows the result of the `rrh status` command.
*/
type StatusResult struct {
	GroupName      string
	RepositoryName string
	BranchName     string
	LastModified   *time.Time
	Description    string
}

type repo struct {
	gname string
	rname string
}

func (status *StatusCommand) lastCommitOnLocalBranch(name repo, r *git.Repository, ref *plumbing.Reference) (*StatusResult, error) {
	iter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, err
	}
	commit, err := iter.Next()
	if err != nil {
		return nil, err
	}
	var signature = commit.Author
	return &StatusResult{name.gname, name.rname, ref.Name().String(), &signature.When, ""}, nil
}

func (status *StatusCommand) openRepository(db *common.Database, repoID string) (*git.Repository, error) {
	var repo = db.FindRepository(repoID)
	if repo == nil {
		return nil, fmt.Errorf("%s: repository not found", repoID)
	}
	var r, err = git.PlainOpen(repo.Path)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func isTarget(options *statusOptions, ref *plumbing.Reference) bool {
	var refName = ref.Name()
	return refName.String() == "HEAD" || options.isRemoteTarget(refName) || options.isBranchTarget(refName)
}

func (status *StatusCommand) findLocalBranches(name repo, r *git.Repository) ([]StatusResult, error) {
	var results = []StatusResult{}
	var iter, err2 = r.References()
	if err2 != nil {
		return results, err2
	}

	iter.ForEach(func(ref *plumbing.Reference) error {
		if isTarget(status.Options, ref) {
			var result, err = status.lastCommitOnLocalBranch(name, r, ref)
			if err != nil {
				return err
			}
			if result.BranchName == "HEAD" {
				var others = []StatusResult{*result}
				results = append(others, results...)
			} else {
				results = append(results, *result)
			}
		}
		return nil
	})
	return results, nil
}

func (status *StatusCommand) findTime(db *common.Database, path string, repoID string) time.Time {
	var repo = db.FindRepository(repoID)
	var target = filepath.Join(repo.Path, path)

	var file, err2 = os.Open(target)
	defer file.Close()
	if err2 != nil {
		fmt.Println(err2.Error())
		return time.Unix(0, 0)
	}
	var fi, err3 = file.Stat()
	if err3 != nil {
		fmt.Println(err3.Error())
		return time.Unix(0, 0)
	}
	return fi.ModTime()
}

func (status *StatusCommand) flagChecker(db *common.Database, rname string, key string, lastModified *time.Time) *time.Time {
	var time = status.findTime(db, key, rname)
	if lastModified == nil || time.After(*lastModified) {
		return &time
	}
	return lastModified
}

func (status *StatusCommand) generateMessage(staging bool, changesNotAdded bool) string {
	if staging && changesNotAdded {
		return "Changes in staging"
	} else if !staging && changesNotAdded {
		return "Changes in workspace"
	}
	return "No changes"
}

func checkUpdateFlag(status git.StatusCode) bool {
	return status != git.Unmodified && status != git.Untracked
}

func (status *StatusCommand) findWorktree(name repo, r *git.Repository, db *common.Database) (*StatusResult, error) {
	var worktree, err = r.Worktree()
	if err != nil {
		return nil, err
	}
	var s, err2 = worktree.Status()
	if err2 != nil {
		return nil, err2
	}
	var lastModified *time.Time
	var staging, changesNotAdded = false, false
	for key, value := range s {
		staging = staging || checkUpdateFlag(value.Staging)
		changesNotAdded = changesNotAdded || checkUpdateFlag(value.Worktree)
		lastModified = status.flagChecker(db, name.rname, key, lastModified)
	}
	return &StatusResult{name.gname, name.rname, "WORKTREE", lastModified, status.generateMessage(staging, changesNotAdded)}, nil
}

func (status *StatusCommand) executeStatusOnRepository(db *common.Database, name repo) ([]StatusResult, error) {
	var r, err = status.openRepository(db, name.rname)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", name.rname, err.Error())
	}
	var results = []StatusResult{}
	var worktree, err2 = status.findWorktree(name, r, db)
	if err2 != nil {
		return nil, err2
	}
	var localBranches, err3 = status.findLocalBranches(name, r)
	if err3 != nil {
		return nil, err3
	}
	results = append(results, *worktree)
	results = append(results, localBranches...)

	return results, nil
}

func (status *StatusCommand) executeStatus(db *common.Database, name string) ([]StatusResult, []error) {
	if db.HasRepository(name) {
		var results, err = status.executeStatusOnRepository(db, repo{"unknown-group", name})
		if err != nil {
			return results, []error{err}
		}
		return results, nil
	} else if db.HasGroup(name) {
		return status.executeStatusOnGroup(db, name)
	}
	return nil, []error{fmt.Errorf("%s: group and repository not found", name)}
}

func (status *StatusCommand) executeStatusOnGroup(db *common.Database, groupName string) ([]StatusResult, []error) {
	var group = db.FindGroup(groupName)
	if group == nil {
		return nil, []error{fmt.Errorf("%s: group not found", groupName)}
	}
	var errors = []error{}
	var results = []StatusResult{}
	for _, repoID := range db.FindRelationsOfGroup(groupName) {
		var sr, err = status.executeStatusOnRepository(db, repo{groupName, repoID})
		if err != nil {
			errors = append(errors, err)
		} else {
			results = append(results, sr...)
		}
	}
	return results, errors
}
