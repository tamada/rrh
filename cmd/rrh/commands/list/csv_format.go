package list

import (
	"encoding/csv"
	"io"

	"github.com/tamada/rrh"
)

type csvFormat struct {
	headerFlag bool
}

func (cf *csvFormat) Format(w io.Writer, r []*Result, li Entries, noAbbrevFlag bool) error {
	writer := csv.NewWriter(w)
	if cf.headerFlag {
		header := li.StringArray()
		writer.Write(header)
	}
	for _, result := range r {
		cf.formatEach(writer, result, li)
	}
	writer.Flush()
	return nil
}

func (cf *csvFormat) formatEach(writer *csv.Writer, r *Result, li Entries) {
	for _, repo := range r.Repos {
		for _, remote := range repo.Remotes {
			array := cf.formatCsvEachRepo(r, repo, remote, li)
			writer.Write(array)
		}
		if len(repo.Remotes) == 0 {
			array := cf.formatCsvEachRepo(r, repo, nil, li)
			writer.Write(array)
		}
	}
}

func (cf *csvFormat) formatCsvEachRepo(r *Result, repo *Repo, remote *rrh.Remote, li Entries) []string {
	result := []string{}
	if li.IsGroupName() {
		result = append(result, r.GroupName)
	}
	if li.IsNote() {
		result = append(result, r.Note)
	}
	if li.IsRepositoryId() {
		result = append(result, repo.Name)
	}
	if li.IsRepositoryDesc() {
		result = append(result, repo.Desc)
	}
	if li.IsRepositoryPath() {
		result = append(result, repo.Path)
	}
	if li.IsRepositoryRemotes() {
		if remote != nil {
			result = append(result, remote.Name)
			result = append(result, remote.URL)
		} else {
			result = append(result, "")
			result = append(result, "")
		}
	}
	return result
}
