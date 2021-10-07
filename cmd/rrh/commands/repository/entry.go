package repository

import (
	"fmt"
	"strings"

	"github.com/tamada/rrh/cmd/rrh/commands/common"
)

type Entries int

const (
	repositoryId   Entries = 1
	repositoryDesc         = 2
	repositoryPath         = 4
	remotes                = 8
	groups                 = 16
	groupCount             = 32
	repositoryAll          = repositoryId | repositoryDesc | repositoryPath | remotes | groups
)

func (re Entries) StringArray() []string {
	results := []string{}
	if re.IsId() {
		results = append(results, "id")
	}
	if re.IsDesc() {
		results = append(results, "description")
	}
	if re.IsPath() {
		results = append(results, "path")
	}
	if re.IsRemotes() {
		results = append(results, "remote name")
		results = append(results, "remote url")
	}
	if re.IsGroup() {
		results = append(results, "groups")
	}
	if re.IsCount() {
		results = append(results, "group count")
	}
	return results
}

func (re Entries) IsId() bool {
	return re&repositoryId == repositoryId
}
func (re Entries) IsDesc() bool {
	return re&repositoryDesc == repositoryDesc
}
func (re Entries) IsPath() bool {
	return re&repositoryPath == repositoryPath
}
func (re Entries) IsRemotes() bool {
	return re&remotes == remotes
}
func (re Entries) IsGroup() bool {
	return re&groups == groups
}
func (re Entries) IsCount() bool {
	return re&groupCount == groupCount
}

func NewEntries(entries []string) (Entries, error) {
	var result Entries = 0
	for _, entry := range entries {
		switch strings.ToLower(entry) {
		case "all":
			result = result | repositoryAll
		case "id":
			result = result | repositoryId
		case "count":
			result = result | groupCount
		case "desc":
			result = result | repositoryDesc
		case "path":
			result = result | repositoryPath
		case "group":
			result = result | groups
		case "remote":
			result = result | remotes
		default:
			return 0, fmt.Errorf("%s: invalid entry, availables are: all, id, desc, path, remote, group, and count", entry)
		}
	}
	return result, nil
}

func ValidateEntries(entries []string) error {
	availables := []string{"all", "id", "desc", "path", "remote", "group", "count"}
	return common.ValidateValues(entries, availables)
}
