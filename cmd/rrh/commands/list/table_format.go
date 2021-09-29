package list

import (
	"io"

	"github.com/olekukonko/tablewriter"
	"github.com/tamada/rrh"
)

type tableFormat struct {
}

func (tf *tableFormat) Format(w io.Writer, r []*Result, le Entries, noAbbrevFlag bool) error {
	writer := tablewriter.NewWriter(w)
	writer.SetHeader(tf.header(le))
	for _, result := range r {
		tf.formatEach(writer, result, le)
	}
	writer.Render()
	return nil
}

func (tf *tableFormat) formatEach(w *tablewriter.Table, r *Result, le Entries) {
	for _, repo := range r.Repos {
		for _, remote := range repo.Remotes {
			array := tf.formatEachRepo(r, repo, remote, le)
			w.Append(array)
		}
		if len(repo.Remotes) == 0 {
			array := tf.formatEachRepo(r, repo, nil, le)
			w.Append(array)
		}
	}
}

func (tf *tableFormat) header(le Entries) []string {
	labels := []string{}
	if le.IsGroupName() {
		labels = append(labels, "group")
	}
	if le.IsNote() {
		labels = append(labels, "note")
	}
	if le.IsRepositoryId() {
		labels = append(labels, "repository")
	}
	if le.IsRepositoryDesc() {
		labels = append(labels, "description")
	}
	if le.IsRepositoryPath() {
		labels = append(labels, "path")
	}
	if le.IsRepositoryRemotes() {
		labels = append(labels, "remote name")
		labels = append(labels, "remote url")
	}
	return labels
}

func (tf *tableFormat) formatEachRepo(r *Result, repo *Repo, remote *rrh.Remote, le Entries) []string {
	result := []string{}
	if le.IsGroupName() {
		result = append(result, r.GroupName)
	}
	if le.IsNote() {
		result = append(result, r.Note)
	}
	if le.IsRepositoryId() {
		result = append(result, repo.Name)
	}
	if le.IsRepositoryDesc() {
		result = append(result, repo.Desc)
	}
	if le.IsRepositoryPath() {
		result = append(result, repo.Path)
	}
	if le.IsRepositoryRemotes() {
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
