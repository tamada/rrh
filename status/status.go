package status

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/tamadalab/rrh/common"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type StatusResult struct {
	GroupName      string
	RepositoryName string
	BranchName     string
	LastModified   *time.Time
	Description    string
}

func (status *StatusCommand) lastCommitOnLocalBranch(gname string, rname string, r *git.Repository, ref *plumbing.Reference) (*StatusResult, error) {
	iter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, err
	}
	commit, err := iter.Next()
	if err != nil {
		return nil, err
	}
	var signature = commit.Author
	return &StatusResult{gname, rname, ref.Name().String(), &signature.When, ""}, nil
}

func (status *StatusCommand) openRepository(db *common.Database, repoID string) (*git.Repository, error) {
	var repo = db.FindRepository(repoID)
	if repo == nil {
		return nil, fmt.Errorf("%s: repository not found", repoID)
	}
	var r, err = git.PlainOpen(common.ToAbsolutePath(repo.Path, db.Config))
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (status *StatusCommand) findLocalBranches(gname string, rname string, r *git.Repository, options *statusOptions) ([]StatusResult, error) {
	var results = []StatusResult{}
	var iter, err2 = r.References()
	if err2 != nil {
		return results, err2
	}

	iter.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().String() == "HEAD" ||
			options.remote && ref.Name().IsRemote() ||
			options.branch && ref.Name().IsBranch() {
			var result, err = status.lastCommitOnLocalBranch(gname, rname, r, ref)
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

func (status *StatusCommand) findTime(path string, repoID string, db *common.Database) time.Time {
	var repo = db.FindRepository(repoID)
	var absPath = common.ToAbsolutePath(repo.Path, db.Config)
	var target = filepath.Join(absPath, path)

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

func (status *StatusCommand) findWorktree(gname string, rname string, r *git.Repository, db *common.Database) (*StatusResult, error) {
	var worktree, err = r.Worktree()
	if err != nil {
		return nil, err
	}
	var s, err2 = worktree.Status()
	if err2 != nil {
		return nil, err2
	}
	var lastModified *time.Time
	var staging = false
	var changesNotAdded = false
	for key, value := range s {
		if value.Staging != ' ' && value.Staging != '?' {
			staging = true
		}
		if value.Worktree != ' ' && value.Worktree != '?' {
			changesNotAdded = true
		}
		var time = status.findTime(key, rname, db)
		if lastModified == nil || time.After(*lastModified) {
			lastModified = &time
		}
		// fmt.Printf("%-20s(%c, %c)\t%s\n", key, value.Staging, value.Worktree, time)
	}
	var message = "No changes"
	if staging && changesNotAdded {
		message = "Changes in staging"
	} else if !staging && changesNotAdded {
		message = "Changes in workspace"
	}
	return &StatusResult{gname, rname, "WORKTREE", lastModified, message}, nil
}

func (status *StatusCommand) executeStatusOnRepository(db *common.Database, gname string, repoID string, options *statusOptions) ([]StatusResult, error) {
	var r, err = status.openRepository(db, repoID)
	if err != nil {
		return nil, err
	}
	var results = []StatusResult{}
	var worktree, err2 = status.findWorktree(gname, repoID, r, db)
	if err2 != nil {
		return nil, err2
	}
	var localBranches, err3 = status.findLocalBranches(gname, repoID, r, options)
	if err3 != nil {
		return nil, err3
	}
	results = append(results, *worktree)
	results = append(results, localBranches...)

	return results, nil
}

func (status *StatusCommand) executeStatus(db *common.Database, name string, options *statusOptions) ([]StatusResult, []error) {
	if db.HasRepository(name) {
		var results, err = status.executeStatusOnRepository(db, "unknown-group", name, options)
		if err != nil {
			return results, []error{err}
		}
		return results, nil
	} else if db.HasGroup(name) {
		return status.executeStatusOnGroup(db, name, options)
	}
	return nil, []error{fmt.Errorf("%s: group and repository not found", name)}
}

func (status *StatusCommand) executeStatusOnGroup(db *common.Database, groupName string, options *statusOptions) ([]StatusResult, []error) {
	var group = db.FindGroup(groupName)
	if group == nil {
		return nil, []error{fmt.Errorf("%s: group not found", groupName)}
	}
	var errors = []error{}
	var results = []StatusResult{}
	for _, repo := range group.Items {
		var sr, err = status.executeStatusOnRepository(db, groupName, repo, options)
		if err != nil {
			errors = append(errors, err)
		} else {
			results = append(results, sr...)
		}
	}

	return results, errors
}
