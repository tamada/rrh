package common

import (
	"fmt"
	"strings"

	"github.com/tamada/rrh"
)

type RepositoryEntry int

const (
	repositoryId   RepositoryEntry = 1
	repositoryDesc                 = 2
	repositoryPath                 = 4
	remotes                        = 8
	groups                         = 16
	groupCount                     = 32
	repositoryAll                  = repositoryId | repositoryDesc | repositoryPath | remotes | groups
)

func (re RepositoryEntry) StringArray() []string {
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
		results = append(results, "remotes")
	}
	if re.IsGroup() {
		results = append(results, "groups")
	}
	if re.IsCount() {
		results = append(results, "group count")
	}
	return results
}

func (re RepositoryEntry) IsId() bool {
	return re&repositoryId == repositoryId
}
func (re RepositoryEntry) IsDesc() bool {
	return re&repositoryDesc == repositoryDesc
}
func (re RepositoryEntry) IsPath() bool {
	return re&repositoryPath == repositoryPath
}
func (re RepositoryEntry) IsRemotes() bool {
	return re&remotes == remotes
}
func (re RepositoryEntry) IsGroup() bool {
	return re&groups == groups
}
func (re RepositoryEntry) IsCount() bool {
	return re&groupCount == groupCount
}

func NewRepositoryEntries(entries []string) (RepositoryEntry, error) {
	var result RepositoryEntry = 0
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
			return 0, fmt.Errorf("%s: invalid entry, availables are: all, id, desc, path, group, remote and count", entry)
		}
	}
	return result, nil
}

func ValidateRepositoryEntries(entries []string) error {
	availables := []string{"all", "id", "count", "desc", "path", "group", "remote"}
	messages := []string{}
	for _, entry := range entries {
		entry = strings.ToLower(entry)
		if !rrh.FindIn(entry, availables) {
			messages = append(messages, entry)
		}
	}
	if len(messages) == 0 {
		return nil
	} else if len(messages) == 1 {
		return fmt.Errorf("%s: unknown entry. available values: %s", messages[0], strings.Join(availables, ","))
	}
	return fmt.Errorf("%s: unknown entries. available values: %s", strings.Join(messages, ","), strings.Join(availables, ", "))
}

type GroupEntry int

const (
	groupName       GroupEntry = 1
	repositoryCount            = 2
	groupDesc                  = 4
	repositories               = 8
	abbrevFlag                 = 16
	groupAll                   = groupName | groupDesc | repositories | abbrevFlag
)

func NewGroupEntries(entries []string) (GroupEntry, error) {
	var result GroupEntry = 0
	for _, entry := range entries {
		switch strings.ToLower(entry) {
		case "all":
			result |= groupAll
		case "name":
			result |= groupName
		case "repo":
			result |= repositories
		case "desc":
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

func (ge GroupEntry) StringArray() []string {
	headers := []string{}
	if ge.IsName() {
		headers = append(headers, "name")
	}
	if ge.IsDesc() {
		headers = append(headers, "description")
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

func ValidateGroupEntries(entries []string) error {
	availables := []string{"all", "name", "count", "desc", "repo", "abbrev"}
	messages := []string{}
	for _, entry := range entries {
		entry = strings.ToLower(entry)
		if !rrh.FindIn(entry, availables) {
			messages = append(messages, entry)
		}
	}
	if len(messages) == 0 {
		return nil
	} else if len(messages) == 1 {
		return fmt.Errorf("%s: unknown entry. available values: %s", messages[0], strings.Join(availables, ","))
	}
	return fmt.Errorf("%s: unknown entries. available values: %s", strings.Join(messages, ","), strings.Join(availables, ", "))
}

func (ge GroupEntry) IsName() bool {
	return ge&groupName == groupName
}

func (ge GroupEntry) IsCount() bool {
	return ge&repositoryCount == repositoryCount
}

func (ge GroupEntry) IsDesc() bool {
	return ge&groupDesc == groupDesc
}

func (ge GroupEntry) IsRepo() bool {
	return ge&repositories == repositories
}

func (ge GroupEntry) IsAbbrev() bool {
	return ge&abbrevFlag == abbrevFlag
}
