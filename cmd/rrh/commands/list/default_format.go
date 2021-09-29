package list

import (
	"bufio"
	"fmt"
	"io"

	"github.com/dustin/go-humanize/english"
	"github.com/tamada/rrh"
)

type defaultFormat struct {
}

func (df *defaultFormat) formatEach(writer *bufio.Writer, r *Result, li Entries, noAbbrevFlag bool) error {
	if li.IsGroupName() {
		writer.WriteString(r.GroupName)
	}
	if li.IsRepositoryCount() {
		writer.WriteString(" (" + english.Plural(len(r.Repos), "repository", "repositories") + ")")
	}
	if !noAbbrevFlag && r.Abbrev {
		writer.WriteString(" (abbreviate repositories)")
	} else {
		if li.IsNote() {
			writer.WriteString(fmt.Sprintf("\n    Note: %s", r.Note))
		}
		printRepositoryInfo(writer, r, li)
	}
	writer.WriteString("\n")
	return nil
}

func printRepositoryInfo(writer *bufio.Writer, r *Result, li Entries) {
	repositoryIdFormatter := repositoryIdFormatter(r)
	for _, repo := range r.Repos {
		writer.WriteString("\n")
		if li.IsRepositoryId() {
			writer.WriteString(fmt.Sprintf(repositoryIdFormatter, repo.Name))
		}
		if li.IsRepositoryPath() {
			writer.WriteString(fmt.Sprintf("    %s", repo.Path))
		}
		if li.IsRepositoryDesc() {
			writer.WriteString(fmt.Sprintf("\n        Desc: %s", repo.Desc))
		}
		if li.IsRepositoryRemotes() {
			for _, remote := range repo.Remotes {
				writer.WriteString(fmt.Sprintf("\n        %s\t%s", remote.Name, remote.URL))
			}
		}
	}
}

func repositoryIdFormatter(r *Result) string {
	return fmt.Sprintf("    %%-%ds", computeWidth(r))
}

func computeWidth(r *Result) int {
	max := 0
	for _, repo := range r.Repos {
		max = rrh.MaxInt(len(repo.Name), max)
	}
	return max
}

func (df *defaultFormat) formatSummary(writer *bufio.Writer, r []*Result) {
	groupCount, repoCount, actualRepo := summarize(r)
	writer.WriteString(fmt.Sprintf("%s, and %s", english.Plural(groupCount, "group", ""), english.Plural(repoCount, "repository", "repositories")))
	if actualRepo != repoCount {
		writer.WriteString(fmt.Sprintf(" (actually %s)", english.Plural(actualRepo, "repository", "repositories")))
	}
	writer.WriteString("\n")
}

func (df *defaultFormat) Format(writer io.Writer, r []*Result, li Entries, noAbbrevFlag bool) error {
	w := bufio.NewWriter(writer)
	for _, result := range r {
		df.formatEach(w, result, li, noAbbrevFlag)
	}
	df.formatSummary(w, r)
	w.Flush()
	return nil
}
