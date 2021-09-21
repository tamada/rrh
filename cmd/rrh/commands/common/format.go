package common

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/karlseguin/jsonwriter"
	"github.com/olekukonko/tablewriter"
)

type Printer interface {
	Print(i ...interface{})
	Printf(format string, i ...interface{})
	Println(i ...interface{})
	PrintErr(i ...interface{})
	PrintErrf(format string, i ...interface{})
	PrintErrln(i ...interface{})
}

type Formatter interface {
	Format(writer io.Writer, headers []string, items [][]string) error
	Print(printer Printer, headers []string, items [][]string) error
}

func printByFormatter(printer Printer, f Formatter, headers []string, items [][]string) error {
	buffer := &bytes.Buffer{}
	err := f.Format(buffer, headers, items)
	if err != nil {
		return err
	}
	printer.Print(buffer.String())
	return nil
}

func ValidateFormatter(formatter string) error {
	availables := []string{"table", "csv", "json"}
	newFormat := strings.ToLower(formatter)
	for _, a := range availables {
		if newFormat == a {
			return nil
		}
	}
	return fmt.Errorf("%s: unknown format. available values: %s", formatter, strings.Join(availables, ","))

}

func NewFormatter(formatter string, withHeader bool) (Formatter, error) {
	switch strings.ToLower(formatter) {
	case "json":
		return &jsonFormat{}, nil
	case "csv":
		return &csvFormat{header: withHeader}, nil
	case "table":
		return &tableFormat{header: withHeader}, nil
	default:
		return nil, fmt.Errorf("%s: unknown format. available values: table, csv, and json", formatter)
	}
}

type jsonFormat struct {
}
type csvFormat struct{ header bool }
type tableFormat struct{ header bool }

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
	table.AppendBulk(values)
	table.Render()
	return nil
}

func (jf *jsonFormat) Print(printer Printer, headers []string, items [][]string) error {
	return printByFormatter(printer, jf, headers, items)
}
func (cf *csvFormat) Print(printer Printer, headers []string, items [][]string) error {
	return printByFormatter(printer, cf, headers, items)
}
func (tf *tableFormat) Print(printer Printer, headers []string, items [][]string) error {
	return printByFormatter(printer, tf, headers, items)
}
