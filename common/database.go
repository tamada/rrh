package common

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"
)

/*
Repository represents a Git repository.
*/
type Repository struct {
	ID   string `json:"repository_id"`
	Path string `json:"repository_path"`
	URL  string `json:"repository_url"`
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
Database represents the whole database of GRIM.
*/
type Database struct {
	Timestamp    time.Time    `json:"last_modified"`
	Repositories []Repository `json:"repositories"`
	Groups       []Group      `json:"gruops"`
}

/*
FindGroups method finds groups from database filtering by given names.
*/
func (db *Database) FindGroups(args []string) []Group {
	if len(args) == 0 {
		return db.Groups
	}
	var groups = []Group{}
	for _, arg := range args {
		var group = db.FindGroup(arg)
		if group != nil {
			groups = append(groups, *group)
		}
	}
	return groups
}

/*
AddGroup method adds group to the database by given name and description.
If given groupName is already registered, this method fails to add group and returns `false`.
*/
func (db *Database) AddGroup(groupName string, description string) *Group {
	var group = db.FindGroup(groupName)
	if group != nil {
		return nil
	}
	var g = Group{groupName, description, []string{}}
	db.Groups = append(db.Groups, g)
	db.Timestamp = time.Now()
	return &g
}

func (db *Database) FindGroup(groupName string) *Group {
	for _, g := range db.Groups {
		if g.Name == groupName {
			return &g
		}
	}
	return nil
}

func (db *Database) removeRepo(repos []Repository, repo Repository) []Repository {
	var results = []Repository{}
	for _, item := range repos {
		if item.ID != repo.ID {
			results = append(results, item)
		}
	}
	return results
}

func (db *Database) removeRepoFromGroup(group *Group, name string) {
	var results = []string{}
	for _, item := range group.Items {
		if item != name {
			results = append(results, item)
		}
	}
	group.Items = results
}

func (db *Database) RemoveRepository(item string) bool {
	var repo = db.FindRepository(item)
	if repo == nil {
		return false
	}
	db.Repositories = db.removeRepo(db.Repositories, *repo)
	for _, group := range db.Groups {
		db.removeRepoFromGroup(&group, repo.ID)
	}
	db.Timestamp = time.Now()

	return true
}

func (db *Database) removeGroup(group Group) []Group {
	var results = []Group{}
	for _, g := range db.Groups {
		if g.Name != group.Name {
			results = append(results, g)
		}
	}
	return results
}

func (db *Database) RemoveGroup(groupName string) bool {
	var group = db.FindGroup(groupName)
	if group != nil {
		return false
	}
	if len(group.Items) != 0 {
		return false
	}
	db.Groups = db.removeGroup(*group)
	db.Timestamp = time.Now()

	return true
}

func (db *Database) AddRepository(repo *Repository, group *Group) bool {
	var found = db.FindRepository(repo.ID)
	if found != nil {
		return false
	}
	db.Repositories = append(db.Repositories, *repo)
	for i, g := range db.Groups {
		if g.Name == group.Name {
			db.Groups[i].Items = append(db.Groups[i].Items, repo.ID)
		}
	}
	return true
}

func (db *Database) FindRepository(item string) *Repository {
	for _, r := range db.Repositories {
		if r.ID == item {
			return &r
		}
	}
	return nil
}

func databasePath(config *Config) string {
	return config.GetValue(GrimDatabasePath)
}

func (db *Database) StoreAndClose(config *Config) error {
	var bytes, err = json.Marshal(db)
	if err == nil {
		return ioutil.WriteFile(databasePath(config), bytes, 0644)
	}
	return err
}

/*
Open function is to read grim database from a certain path.

How to call this function

    var db *Database
	db = common.Open()
*/
func Open(config *Config) *Database {
	bytes, err := ioutil.ReadFile(databasePath(config))
	if err != nil {
		var db = Database{time.Now(), []Repository{}, []Group{}}
		return &db
	}

	var db Database
	if err := json.Unmarshal(bytes, &db); err != nil {
		log.Fatal(err)
	}
	if db.Repositories == nil {
		db.Repositories = []Repository{}
	}
	if db.Groups == nil {
		db.Groups = []Group{}
	}
	return &db
}
