package list

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/dustin/go-humanize/english"
	"github.com/tamada/rrh"
)

type formatter interface {
	Format(w io.Writer, r []*Result, li Entries, noAbbrevFlag bool) error
}

func toBytes(f formatter, r []*Result, li Entries, noAbbrevFlag bool) (string, error) {
	buffer := &bytes.Buffer{}
	err := f.Format(buffer, r, li, noAbbrevFlag)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func validateFormat(formatter string) error {
	availables := []string{"default", "json", "csv", "table"}
	if !rrh.FindIn(formatter, availables) {
		return fmt.Errorf("%s: unknown format. available values: %s", formatter, strings.Join(availables, ", "))
	}
	return nil
}

func newFormatter(formatter string) (formatter, error) {
	switch strings.ToLower(formatter) {
	case "default":
		return &defaultFormat{}, nil
	case "json":
		return &jsonFormat{}, nil
	case "csv":
		return &csvFormat{}, nil
	case "table":
		return &tableFormat{}, nil
	default:
		return nil, fmt.Errorf("%s: unknown formatter", formatter)
	}
}

type defaultFormat struct {
}
type jsonFormat struct {
}
type csvFormat struct {
}
type tableFormat struct {
}

func (df *defaultFormat) formatEach(writer *bufio.Writer, r *Result, li Entries, noAbbrevFlag bool) error {
	if li.IsGroupName() {
		writer.WriteString(r.GroupName)
	}
	if li.IsRepositoryCount() {
		writer.WriteString(" (" + english.Plural(len(r.Repos), "repository", "repositories") + ")")
	}
	if !noAbbrevFlag && r.Abbrev {
		writer.WriteString(" (abbreviate repositiries)")
	}
	if li.IsNote() {
		writer.WriteString(fmt.Sprintf("\n    Note: %s", r.Note))
	}
	if noAbbrevFlag || !r.Abbrev {
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
	groupCount := len(r)
	repos := map[string]string{}
	repoCount := 0
	for _, result := range r {
		repoCount = repoCount + len(result.Repos)
		for _, repo := range result.Repos {
			repos[repo.Path] = repo.Name
		}
	}
	writer.WriteString(fmt.Sprintf("%s, %s", english.Plural(groupCount, "group", ""), english.Plural(repoCount, "repository", "repositories")))
	if len(repos) != repoCount {
		writer.WriteString(fmt.Sprintf(" (actually %s)", english.Plural(len(repos), "repository", "repositories")))
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
func (jf *jsonFormat) Format(w io.Writer, r []*Result, li Entries, noAbbrevFlag bool) error {
	return nil
}
func (cf *csvFormat) Format(w io.Writer, r []*Result, li Entries, noAbbrevFlag bool) error {
	return nil
}
func (tf *tableFormat) Format(w io.Writer, r []*Result, le Entries, noAbbrevFlag bool) error {
	return nil
}
