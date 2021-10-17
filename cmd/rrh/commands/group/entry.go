package group

import (
	"fmt"
	"strings"

	"github.com/tamada/rrh/cmd/rrh/commands/utils"
)

type Entries int

const (
	groupName       Entries = 1
	repositoryCount         = 2
	groupDesc               = 4
	repositories            = 8
	abbrevFlag              = 16
	groupAll                = groupName | groupDesc | repositories | abbrevFlag
)

func NewEntries(entries []string) (Entries, error) {
	var result Entries = 0
	for _, entry := range entries {
		switch strings.ToLower(entry) {
		case "all":
			result |= groupAll
		case "name":
			result |= groupName
		case "repo":
			result |= repositories
		case "note":
			result |= groupDesc
		case "abbrev":
			result |= abbrevFlag
		case "count":
			result |= repositoryCount
		default:
			return 0, fmt.Errorf("%s: invalid entry, availables are: all, name, desc, repo, abbrev and count", entry)
		}
	}
	return result, nil
}

func (ge Entries) StringArray() []string {
	headers := []string{}
	if ge.IsName() {
		headers = append(headers, "name")
	}
	if ge.IsDesc() {
		headers = append(headers, "note")
	}
	if ge.IsAbbrev() {
		headers = append(headers, "abbrev")
	}
	if ge.IsRepo() {
		headers = append(headers, "repositories")
	}
	if ge.IsCount() {
		headers = append(headers, "repository count")
	}
	return headers
}

func ValidateEntries(entries []string) error {
	availables := []string{"all", "name", "count", "note", "repo", "abbrev"}
	return utils.ValidateValues(entries, availables)
}

func (ge Entries) IsName() bool {
	return ge&groupName == groupName
}

func (ge Entries) IsCount() bool {
	return ge&repositoryCount == repositoryCount
}

func (ge Entries) IsDesc() bool {
	return ge&groupDesc == groupDesc
}

func (ge Entries) IsRepo() bool {
	return ge&repositories == repositories
}

func (ge Entries) IsAbbrev() bool {
	return ge&abbrevFlag == abbrevFlag
}
