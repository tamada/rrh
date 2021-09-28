package list

import (
	"fmt"
	"strings"

	"github.com/tamada/rrh/cmd/rrh/commands/common"
)

type Entries int

const (
	groupName         Entries = 1
	note                      = 2
	repositoryId              = 4
	repositoryDesc            = 8
	repositoryCount           = 16
	repositoryPath            = 32
	repositoryRemotes         = 64
	summary                   = 128
	all                       = groupName | note | repositoryId | repositoryDesc | repositoryCount | repositoryPath | repositoryRemotes | summary
)

func (le Entries) IsGroupName() bool {
	return le&groupName == groupName
}
func (le Entries) IsNote() bool {
	return le&note == note
}
func (le Entries) IsRepositoryId() bool {
	return le&repositoryId == repositoryId
}
func (le Entries) IsRepositoryDesc() bool {
	return le&repositoryDesc == repositoryDesc
}
func (le Entries) IsRepositoryCount() bool {
	return le&repositoryCount == repositoryCount
}
func (le Entries) IsRepositoryPath() bool {
	return le&repositoryPath == repositoryPath
}
func (le Entries) IsRepositoryRemotes() bool {
	return le&repositoryRemotes == repositoryRemotes
}
func (le Entries) IsSummary() bool {
	return le&summary == summary
}

func newListEntry(entries []string) (Entries, error) {
	if err := ValidateEntries(entries); err != nil {
		return -1, err
	}
	var result Entries = 0
	for _, entry := range entries {
		switch strings.ToLower(entry) {
		case "group":
			result = result | groupName
		case "note":
			result = result | note
		case "count":
			result = result | repositoryCount
		case "id":
			result = result | repositoryId
		case "desc":
			result = result | repositoryDesc
		case "path":
			result = result | repositoryPath
		case "remote":
			result = result | repositoryRemotes
		case "summary":
			result = result | summary
		case "all":
			result = result | all
		default:
			return -1, fmt.Errorf("%s: unknown entry", entry)
		}
	}
	return result, nil
}

func ValidateEntries(entries []string) error {
	availables := []string{"group", "note", "id", "desc", "count", "path", "summary", "remote", "all"}
	return common.ValidateValues(entries, availables)
}
