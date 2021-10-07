package list

import (
	"io"
	"strings"

	"github.com/tamada/rrh/cmd/rrh/commands/common"
)

type formatter interface {
	Format(w io.Writer, r []*Result, li Entries, noAbbrevFlag bool) error
}

func validateFormat(formatter string) error {
	availables := []string{"default", "json", "csv", "table"}
	return common.ValidateValue(formatter, availables)
}

func newFormatter(formatter string, headerFlag bool) (formatter, error) {
	if err := validateFormat(formatter); err != nil {
		return nil, err
	}
	switch strings.ToLower(formatter) {
	case "default":
		return &defaultFormat{}, nil
	case "json":
		return &jsonFormat{}, nil
	case "csv":
		return &csvFormat{headerFlag: headerFlag}, nil
	case "table":
		return &tableFormat{headerFlag: headerFlag}, nil
	default:
		panic("never reach this line!")
	}
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