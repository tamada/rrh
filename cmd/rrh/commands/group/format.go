package group

import (
	"encoding/csv"
	"io"
	"strings"

	"github.com/karlseguin/jsonwriter"
	"github.com/olekukonko/tablewriter"
	"github.com/tamada/rrh/cmd/rrh/commands/utils"
)

type Formatter interface {
	Format(writer io.Writer, headers []string, items [][]string) error
}

func ValidateFormatter(formatter string) error {
	availables := []string{"table", "csv", "json", "default"}
	return utils.ValidateValue(formatter, availables)

}

func NewFormatter(formatter string, withHeader bool) (Formatter, error) {
	if err := ValidateFormatter(formatter); err != nil {
		return nil, err
	}
	switch strings.ToLower(formatter) {
	case "json":
		return &jsonFormat{}, nil
	case "csv":
		return &csvFormat{header: withHeader}, nil
	case "table":
		return &tableFormat{header: withHeader, border: true}, nil
	case "default":
		return &tableFormat{header: false, border: false}, nil
	default:
		panic("never reach this line!")
	}
}

type jsonFormat struct {
}
type csvFormat struct{ header bool }
type tableFormat struct {
	header bool
	border bool
}

func (jf *jsonFormat) Format(w io.Writer, headers []string, values [][]string) error {
	writer := jsonwriter.New(w)
	writer.RootArray(func() {
		for _, line := range values {
			writer.ArrayObject(func() {
				for index, item := range headers {
					if item == "repositories" {
						writer.Array(item, func() { writeRepositories(writer, line[index]) })
					} else {
						writer.KeyString(item, line[index])
					}
				}
			})
		}
	})
	return nil
}

func writeRepositories(writer *jsonwriter.Writer, repos string) {
	repoList := strings.Trim(repos, "[]")
	for _, repo := range strings.Split(repoList, " ") {
		writer.Value(repo)
	}
}

func (cf *csvFormat) Format(w io.Writer, headers []string, values [][]string) error {
	writer := csv.NewWriter(w)
	if cf.header {
		writer.Write(headers)
	}
	for _, line := range values {
		writer.Write(line)
	}
	writer.Flush()
	return writer.Error()
}

func (tf *tableFormat) Format(w io.Writer, headers []string, values [][]string) error {
	table := tablewriter.NewWriter(w)
	if tf.header {
		table.SetHeader(headers)
	}
	table.SetBorder(tf.border)
	if !tf.border {
		table.SetNoWhiteSpace(true)
		table.SetTablePadding("  ")
	}
	table.AppendBulk(values)
	table.Render()
	return nil
}
