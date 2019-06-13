package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

/*
StatusOption represents for getting status of branches in a repository.
*/
type StatusOption struct {
	BranchStatus bool
	RemoteStatus bool
}

/*
Status shows the result of the `rrh status` command.
*/
type Status struct {
	Relation     *Relation
	BranchName   string
	LastModified *time.Time
	Description  string
}

/*
NewStatusOption generates an instance of StatusOption.
*/
func NewStatusOption() *StatusOption {
	return &StatusOption{false, false}
}

func openRepository(db *Database, repoID string) (*git.Repository, error) {
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

func checkUpdateFlag(status git.StatusCode) bool {
	return status != git.Unmodified && status != git.Untracked
}

func findStatus(r *git.Repository) (git.Status, error) {
	var worktree, err = r.Worktree()
	if err != nil {
		return nil, err
	}
	var s, err2 = worktree.Status()
	if err2 != nil {
		return nil, err2
	}
	return s, nil
}

func findTime(db *Database, path string, repoID string) time.Time {
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

func flagChecker(time *time.Time, lastModified *time.Time) *time.Time {
	if lastModified == nil || time.After(*lastModified) {
		return time
	}
	return lastModified
}

func findWorktree(name *Relation, r *git.Repository, db *Database) (*Status, error) {
	var s, err = findStatus(r)
	if err != nil {
		return nil, err
	}
	var lastModified *time.Time
	var staging, changesNotAdded = false, false
	for key, value := range s {
		staging = staging || checkUpdateFlag(value.Staging)
		changesNotAdded = changesNotAdded || checkUpdateFlag(value.Worktree)
		var time = findTime(db, key, name.RepositoryID)
		lastModified = flagChecker(&time, lastModified)
	}
	return &Status{name, "WORKTREE", lastModified, generateMessage(staging, changesNotAdded)}, nil
}

func (status *StatusOption) isRemoteTarget(name plumbing.ReferenceName) bool {
	return status.RemoteStatus && name.IsRemote()
}

func (status *StatusOption) isBranchTarget(name plumbing.ReferenceName) bool {
	return status.BranchStatus && name.IsBranch()
}

func (status *StatusOption) isTarget(ref *plumbing.Reference) bool {
	var refName = ref.Name()
	return refName.String() == "HEAD" || status.isRemoteTarget(refName) || status.isBranchTarget(refName)
}

func lastCommitOnLocalBranch(name *Relation, r *git.Repository, ref *plumbing.Reference) (*Status, error) {
	iter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, err
	}
	commit, err := iter.Next()
	if err != nil {
		return nil, err
	}
	var signature = commit.Author
	return &Status{name, ref.Name().String(), &signature.When, ""}, nil
}

func generateMessage(staging bool, changesNotAdded bool) string {
	if staging && changesNotAdded {
		return "Changes in staging"
	} else if !staging && changesNotAdded {
		return "Changes in workspace"
	}
	return "No changes"
}

func (status *StatusOption) findLocalBranches(name *Relation, r *git.Repository) ([]Status, error) {
	var results = []Status{}
	var iter, err2 = r.References()
	if err2 != nil {
		return results, err2
	}

	iter.ForEach(func(ref *plumbing.Reference) error {
		if status.isTarget(ref) {
			var branchResult, err = lastCommitOnLocalBranch(name, r, ref)
			if err != nil {
				return err
			}
			if branchResult.BranchName == "HEAD" {
				var others = []Status{*branchResult}
				results = append(others, results...)
			} else {
				results = append(results, *branchResult)
			}
		}
		return nil
	})
	return results, nil
}

/*
StatusOfRepository returns statuses of a given repository.
*/
func (status *StatusOption) StatusOfRepository(db *Database, name *Relation) ([]Status, error) {
	var r, err = openRepository(db, name.RepositoryID)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", name.RepositoryID, err.Error())
	}
	var results = []Status{}
	var worktree, err2 = findWorktree(name, r, db)
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

/*
FindRemotes function returns the remote of the given git repository.
*/
func FindRemotes(path string) ([]Remote, error) {
	var repo, err = git.PlainOpen(path)
	if err != nil {
		return nil, err
	}
	remotes, err := repo.Remotes()
	if err != nil {
		return nil, err
	}
	var crs = []Remote{}
	for _, remote := range remotes {
		var config = remote.Config()
		crs = append(crs, Remote{Name: config.Name, URL: config.URLs[0]})
	}
	return crs, nil
}
