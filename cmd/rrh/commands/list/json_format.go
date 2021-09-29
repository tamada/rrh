package list

import (
	"io"

	"github.com/karlseguin/jsonwriter"
)

type jsonFormat struct {
}

func (jf *jsonFormat) Format(w io.Writer, r []*Result, li Entries, noAbbrevFlag bool) error {
	writer := jsonwriter.New(w)
	writer.RootObject(func() {
		writer.Array("groups", func() {
			formatResults(writer, r, li, noAbbrevFlag)
		})
		if li.IsSummary() {
			writer.Object("summary", func() {
				formatSummary(writer, r)
			})
		}
	})
	return nil
}

func formatResults(writer *jsonwriter.Writer, r []*Result, li Entries, noAbbrevFlag bool) {
	for _, result := range r {
		writer.ArrayObject(func() {
			formatResult(writer, result, li, noAbbrevFlag)
		})
	}
}

func formatResult(writer *jsonwriter.Writer, r *Result, li Entries, noAbbrevFlag bool) {
	if li.IsGroupName() {
		writer.KeyValue("group-name", r.GroupName)
	}
	if li.IsRepositoryCount() {
		writer.KeyValue("repository-count", len(r.Repos))
	}
	if li.IsNote() {
		writer.KeyValue("note", r.Note)
	}
	if (noAbbrevFlag || !r.Abbrev) && len(r.Repos) > 0 {
		writer.Array("repositories", func() {
			for _, repo := range r.Repos {
				writer.ArrayObject(func() {
					formatJsonEachRepo(writer, repo, li, noAbbrevFlag)
				})
			}
		})
	} else {
		writer.KeyValue("repositories", "abbreviated")
	}
}

func formatJsonEachRepo(writer *jsonwriter.Writer, repo *Repo, li Entries, noAbbrevFlag bool) {
	if li.IsRepositoryId() {
		writer.KeyValue("id", repo.Name)
	}
	if li.IsRepositoryDesc() {
		writer.KeyValue("description", repo.Desc)
	}
	if li.IsRepositoryPath() {
		writer.KeyValue("path", repo.Path)
	}
	if li.IsRepositoryRemotes() && len(repo.Remotes) > 0 {
		writer.Array("remotes", func() {
			for _, r := range repo.Remotes {
				writer.ArrayObject(func() {
					writer.KeyValue("name", r.Name)
					writer.KeyValue("url", r.URL)
				})
			}
		})
	}
}

func formatSummary(writer *jsonwriter.Writer, r []*Result) {
	groupCount, repoCount, actualRepoCount := summarize(r)
	writer.KeyValue("group-count", groupCount)
	writer.KeyValue("repository-count", repoCount)
	writer.KeyValue("actural-repository-count", actualRepoCount)
}
