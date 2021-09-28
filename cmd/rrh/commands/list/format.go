package list

import (
	"bytes"
	"io"
	"strings"

	"github.com/tamada/rrh/cmd/rrh/commands/common"
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
	return common.ValidateValue(formatter, availables)
}

func newFormatter(formatter string) (formatter, error) {
	if err := validateFormat(formatter); err != nil {
		return nil, err
	}
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
		panic("never reach this line!")
	}
}

type csvFormat struct {
}
type tableFormat struct {
}

func (cf *csvFormat) Format(w io.Writer, r []*Result, li Entries, noAbbrevFlag bool) error {
	return nil
}
func (tf *tableFormat) Format(w io.Writer, r []*Result, le Entries, noAbbrevFlag bool) error {
	return nil
}

func summarize(r []*Result) (groupCount, repositoryCount, actualRepositoryCount int) {
	repos := map[string]string{}
	repoCount := 0
	for _, result := range r {
		repoCount = repoCount + len(result.Repos)
		for _, repo := range result.Repos {
			repos[repo.Path] = repo.Name
		}
	}
	return len(r), repoCount, len(repos)
}
