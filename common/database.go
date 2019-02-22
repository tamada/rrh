package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type Remote struct {
	Name string
	URL  string
}

/*
Repository represents a Git repository.
*/
type Repository struct {
	ID      string   `json:"repository_id"`
	Path    string   `json:"repository_path"`
	Remotes []Remote `json:"remotes"`
}

/*
Group represents the groups of the Git repositories.
*/
type Group struct {
	Name        string   `json:"group_name"`
	Description string   `json:"group_desc"`
	Items       []string `json:"group_items"`
}

/*
Database represents the whole database of RRH.
*/
type Database struct {
	Timestamp    time.Time    `json:"last_modified"`
	Repositories []Repository `json:"repositories"`
	Groups       []Group      `json:"groups"`
	Config       *Config      `json:"-"`
}

func (db *Database) countReposInGroups() map[string]int {
	var repoFlags = map[string]int{}
	for _, repo := range db.Repositories {
		repoFlags[repo.ID] = 0
	}
	for _, group := range db.Groups {
		for _, item := range group.Items {
			repoFlags[item] = repoFlags[item] + 1
		}
	}
	return repoFlags
}

func (db *Database) pruneGroup() int {
	var newGroups = []Group{}
	var count = 0
	for _, group := range db.Groups {
		if len(group.Items) != 0 {
			newGroups = append(newGroups, group)
		} else {
			count++
		}
	}
	db.Groups = newGroups

	return count
}

/*
Prune eliminates unnecessary repositories, and groups from db.
*/
func (db *Database) Prune() (int, int) {
	var prunedGroupCount = db.pruneGroup()
	var repoFlags = db.countReposInGroups()
	var repos = []Repository{}
	var prunedReposCount = 0
	for _, repo := range db.Repositories {
		if repoFlags[repo.ID] != 0 {
			repos = append(repos, repo)
		} else {
			prunedReposCount++
		}
	}
	db.Repositories = repos
	return prunedGroupCount, prunedReposCount
}

/*
FindRepository returns the repository which ID is given `repoID.`
*/
func (db *Database) FindRepository(repoID string) *Repository {
	for _, repo := range db.Repositories {
		if repo.ID == repoID {
			return &repo
		}
	}
	return nil
}

/*
FindGroup returns the group which name is given `groupID.`
*/
func (db *Database) FindGroup(groupID string) *Group {
	for _, group := range db.Groups {
		if group.Name == groupID {
			return &group
		}
	}
	return nil
}

/*
CreateRepository returns the repository by creating the given parameters and store it to database.
*/
func (db *Database) CreateRepository(repoID string, path string, remotes []Remote) (*Repository, error) {
	if db.HasRepository(repoID) {
		return nil, fmt.Errorf("%s: already registered repository", repoID)
	}
	var repo = Repository{repoID, path, remotes}
	db.Repositories = append(db.Repositories, repo)

	return &repo, nil
}

/*
CreateGroup returns the group by creating the given parameters and store it to database.
*/
func (db *Database) CreateGroup(groupID string, description string) (*Group, error) {
	if db.HasGroup(groupID) {
		return nil, fmt.Errorf("%s: already registered group", groupID)
	}
	var group = Group{groupID, description, []string{}}
	db.Groups = append(db.Groups, group)

	return &group, nil
}

/*
UpdateGroup updates found group with `newGroupID` and `newDescription`.
The return value is that the update is success or not.
*/
func (db *Database) UpdateGroup(groupID string, newGroupID string, newDescription string) bool {
	if !db.HasGroup(groupID) {
		return false
	}
	for i, group := range db.Groups {
		if group.Name == groupID {
			db.Groups[i].Name = newGroupID
			db.Groups[i].Description = newDescription
		}
	}

	return true
}

/*
Relate create the relation between the group and the repository.
The group and the repository are specified by the given parameters.
If the group and the repository have the relation, this function returns `nil` (successfully create relation).
*/
func (db *Database) Relate(groupID string, repoID string) error {
	if db.HasRelation(groupID, repoID) {
		return nil
	}
	for i, group := range db.Groups {
		if group.Name == groupID {
			db.Groups[i].Items = append(group.Items, repoID)
		}
	}

	return nil
}

/*
HasRelation returns true if the group and the repository has relation.
The group and the repository are specified by the given parameters.
*/
func (db *Database) HasRelation(groupID string, repoID string) bool {
	var group = db.FindGroup(groupID)
	if group == nil {
		return false
	}
	for _, item := range group.Items {
		if item == repoID {
			return true
		}
	}
	return false
}

/*
Unrelate delete the relation between the group and the repository.
The group and the repository are specified by the given parameters.
If the group and the repository do not have the relation, this function returns `nil` (successfully delete relation).
*/
func (db *Database) Unrelate(groupID string, repoID string) error {
	if !db.HasRelation(groupID, repoID) {
		return nil
	}
	for i, group := range db.Groups {
		if group.Name == groupID {
			var newItems = []string{}
			for _, item := range group.Items {
				if item != repoID {
					newItems = append(newItems, item)
				}
			}
			db.Groups[i].Items = newItems
		}
	}
	return nil
}

/*
HasRepository returns true if the db has the repository of repoID.
*/
func (db *Database) HasRepository(repoID string) bool {
	for _, repo := range db.Repositories {
		if repo.ID == repoID {
			return true
		}
	}
	return false
}

/*
HasGroup returns true if the db has the group of groupID.
*/
func (db *Database) HasGroup(groupID string) bool {
	for _, group := range db.Groups {
		if group.Name == groupID {
			return true
		}
	}
	return false
}

/*
DeleteRepository function removes the repository of given repoID from DB.
Also, the relation between the repository and groups are removed.
*/
func (db *Database) DeleteRepository(repoID string) error {
	if !db.HasRepository(repoID) {
		return fmt.Errorf("%s: repository not found", repoID)
	}
	for _, group := range db.Groups {
		db.Unrelate(group.Name, repoID)
	}
	var newRepositories = []Repository{}
	for _, repo := range db.Repositories {
		if repo.ID != repoID {
			newRepositories = append(newRepositories, repo)
		}
	}
	db.Repositories = newRepositories

	return nil
}

func (db *Database) deleteGroup(groupID string) error {
	var groups = []Group{}
	for _, group := range db.Groups {
		if group.Name != groupID {
			groups = append(groups, group)
		}
	}
	db.Groups = groups

	return nil
}

/*
DeleteGroup removes the group of the given groupId from DB.
If the group has some repositories, the function fails to remove.
*/
func (db *Database) DeleteGroup(groupID string) error {
	if !db.HasGroup(groupID) {
		return fmt.Errorf("%s: group not found", groupID)
	}
	var group = db.FindGroup(groupID)
	if len(group.Items) != 0 {
		return fmt.Errorf("%s: group has %d relatins", groupID, len(group.Items))
	}
	var err = db.deleteGroup(groupID)
	return err
}

/*
ForceDeleteGroup removes the group of the given groupID from DB.
Even if the group has some repositories, the function forcely remove the group.
*/
func (db *Database) ForceDeleteGroup(groupID string) error {
	if !db.HasGroup(groupID) {
		return fmt.Errorf("%s: group not found", groupID)
	}
	return db.deleteGroup(groupID)
}

func databasePath(config *Config) string {
	return config.GetValue(RrhDatabasePath)
}

/*
StoreAndClose stores the database to file and close the database.
The database path is defined in RRH_DATABASE_PATH of config.
*/
func (db *Database) StoreAndClose() error {
	db.Timestamp = time.Now()
	var bytes, err = json.Marshal(db)
	if err != nil {
		return err
	}
	var databasePath = databasePath(db.Config)
	var err1 = CreateParentDir(databasePath)
	if err1 != nil {
		return err1
	}
	return ioutil.WriteFile(databasePath, bytes, 0644)
}

/*
Open function is to read rrh database from a certain path.

How to call this function

    var db *Database
	db = common.Open()
*/
func Open(config *Config) (*Database, error) {
	bytes, err := ioutil.ReadFile(databasePath(config))
	if err != nil {
		var db = Database{time.Now(), []Repository{}, []Group{}, config}
		return &db, nil
	}

	var db Database
	if err := json.Unmarshal(bytes, &db); err != nil {
		return nil, err
	}
	if db.Repositories == nil {
		db.Repositories = []Repository{}
	}
	if db.Groups == nil {
		db.Groups = []Group{}
	}
	db.Config = config
	return &db, nil
}
