package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
)

/*
Remote represents remote of the git repository.
*/
type Remote struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

/*
Repository represents a Git repository.
*/
type Repository struct {
	ID          string   `json:"repository_id"`
	Path        string   `json:"repository_path"`
	Description string   `json:"repository_desc"`
	Remotes     []Remote `json:"remotes"`
}

/*
Group represents the groups of the Git repositories.
*/
type Group struct {
	Name        string `json:"group_name"`
	Description string `json:"group_desc"`
	OmitList    bool   `json:"omit_list"`
}

/*
Relation represents the relation between group and the repository.
*/
type Relation struct {
	RepositoryID string `json:"repository_id"`
	GroupName    string `json:"group_name"`
}

/*
Database represents the whole database of RRH.
*/
type Database struct {
	Timestamp    RrhTime      `json:"last_modified"`
	Repositories []Repository `json:"repositories"`
	Groups       []Group      `json:"groups"`
	Relations    []Relation   `json:"relations"`
	Config       *Config      `json:"-"`
}

func groupFrequencies(db *Database) map[string]int {
	var groupMap = map[string]int{}
	for _, group := range db.Groups {
		groupMap[group.Name] = 0
	}
	for _, relation := range db.Relations {
		groupMap[relation.GroupName] = groupMap[relation.GroupName] + 1
	}
	return groupMap
}

func pruneGroups(db *Database) int {
	var groupMap = groupFrequencies(db)
	var newGroups = []Group{}
	var count = 0
	for _, group := range db.Groups {
		if groupMap[group.Name] != 0 {
			newGroups = append(newGroups, group)
		} else {
			count++
		}
	}
	db.Groups = newGroups

	return count
}

func repositoryFrequencies(db *Database) map[string]int {
	var repoFlags = map[string]int{}
	for _, repo := range db.Repositories {
		repoFlags[repo.ID] = 0
	}
	for _, relation := range db.Relations {
		repoFlags[relation.RepositoryID] = repoFlags[relation.RepositoryID] + 1
	}
	return repoFlags
}

func pruneRepositories(db *Database) int {
	var repoFlags = repositoryFrequencies(db)
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
	return prunedReposCount
}

/*
Prune eliminates unnecessary repositories, and groups from db.
*/
func (db *Database) Prune() (int, int) {
	var prunedGroupCount = pruneGroups(db)
	var prunedReposCount = pruneRepositories(db)
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

func sortIfNeeded(db *Database) {
	if !db.Config.IsSet(RrhSortOnUpdating) {
		return
	}
	sort.Slice(db.Repositories, func(i, j int) bool {
		return db.Repositories[i].ID < db.Repositories[j].ID
	})
	sort.Slice(db.Groups, func(i, j int) bool {
		return db.Groups[i].Name < db.Groups[j].Name
	})
	sort.Slice(db.Relations, func(i, j int) bool {
		if db.Relations[i].GroupName == db.Relations[j].GroupName {
			return db.Relations[i].RepositoryID < db.Relations[j].RepositoryID
		}
		return db.Relations[i].GroupName < db.Relations[j].GroupName
	})
}

/*
CreateRepository returns the repository by creating the given parameters and store it to database.
*/
func (db *Database) CreateRepository(repoID string, path string, desc string, remotes []Remote) (*Repository, error) {
	if db.HasRepository(repoID) {
		return nil, fmt.Errorf("%s: already registered repository", repoID)
	}
	var repo = Repository{repoID, path, desc, remotes}
	db.Repositories = append(db.Repositories, repo)
	sortIfNeeded(db)

	return &repo, nil
}

/*
AutoCreateGroup returns the group by creating the given parameters and store it to the database
by satisfying the config value of RrhAutoCreateGroup is true.
*/
func (db *Database) AutoCreateGroup(groupID string, description string, omitList bool) (*Group, error) {
	if db.HasGroup(groupID) {
		return db.FindGroup(groupID), nil
	}
	if db.Config.IsSet(RrhAutoCreateGroup) {
		return db.CreateGroup(groupID, description, omitList)
	}
	return nil, fmt.Errorf("%s: could not create group", groupID)
}

/*
CreateGroup returns the group by creating the given parameters and store it to database.
If db has the group with the given groupID, this method returns error.
*/
func (db *Database) CreateGroup(groupID string, description string, omitList bool) (*Group, error) {
	if db.HasGroup(groupID) {
		return nil, fmt.Errorf("%s: already registered group", groupID)
	}
	var group = Group{groupID, description, omitList}
	db.Groups = append(db.Groups, group)
	sortIfNeeded(db)

	return &group, nil
}

/*
UpdateGroup updates found group with `newGroupID` and `newDescription`.
The return value is that the update is success or not.
*/
func (db *Database) UpdateGroup(groupID string, newGroup Group) bool {
	if !db.HasGroup(groupID) {
		return false
	}
	for i, group := range db.Groups {
		if group.Name == groupID {
			updateGroupImpl(db, i, newGroup)
		}
	}
	sortIfNeeded(db)

	return true
}

func updateGroupImpl(db *Database, index int, newGroup Group) {
	var oldGroup = db.Groups[index].Name
	db.Groups[index].Name = newGroup.Name
	db.Groups[index].Description = newGroup.Description
	db.Groups[index].OmitList = newGroup.OmitList
	for i, relation := range db.Relations {
		if relation.GroupName == oldGroup {
			db.Relations[i].GroupName = newGroup.Name
		}
	}
}

/*
UpdateRepository updates the repository information.
*/
func (db *Database) UpdateRepository(repositoryID string, newRepo Repository) bool {
	if !db.HasRepository(repositoryID) {
		return false
	}
	for i, repo := range db.Repositories {
		if repo.ID == repositoryID {
			updateRepositoryImpl(db, i, newRepo)
		}
	}
	sortIfNeeded(db)
	return true
}

func updateRepositoryImpl(db *Database, index int, newRepo Repository) {
	var oldID = db.Repositories[index].ID
	db.Repositories[index].ID = newRepo.ID
	db.Repositories[index].Description = newRepo.Description
	db.Repositories[index].Path = newRepo.Path
	for i, rel := range db.Relations {
		if rel.RepositoryID == oldID {
			db.Relations[i].RepositoryID = newRepo.ID
		}
	}

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
	db.Relations = append(db.Relations, Relation{repoID, groupID})
	sortIfNeeded(db)

	return nil
}

/*
BelongingCount returns the number of groups belonging given repoID.
*/
func (db *Database) BelongingCount(repoID string) int {
	var repos = repositoryFrequencies(db)
	return repos[repoID]
}

/*
ContainsCount returns the number of repositories in the given groupID.
*/
func (db *Database) ContainsCount(groupID string) int {
	var groups = groupFrequencies(db)
	return groups[groupID]
}

/*
FindRelationsOfGroup returns the repository ids belonging to the specified group.
*/
func (db *Database) FindRelationsOfGroup(groupID string) []string {
	var repositories = []string{}
	for _, relation := range db.Relations {
		if relation.GroupName == groupID {
			repositories = append(repositories, relation.RepositoryID)
		}
	}
	return repositories
}

/*
FindRelationsOfRepository returns the group names of the specified repository.
*/
func (db *Database) FindRelationsOfRepository(repositoryID string) []string {
	var groups = []string{}
	for _, relation := range db.Relations {
		if relation.RepositoryID == repositoryID {
			groups = append(groups, relation.GroupName)
		}
	}
	return groups
}

/*
HasRelation returns true if the group and the repository has relation.
The group and the repository are specified by the given parameters.
*/
func (db *Database) HasRelation(groupID string, repoID string) bool {
	for _, relation := range db.Relations {
		if relation.GroupName == groupID && relation.RepositoryID == repoID {
			return true
		}
	}
	return false
}

/*
Unrelate deletes the relation between the group and the repository.
The group and the repository are specified by the given parameters.
If the group and the repository do not have the relation, this function returns `nil` (successfully delete relation).
*/
func (db *Database) Unrelate(groupID string, repoID string) {
	if !db.HasRelation(groupID, repoID) {
		return
	}
	var newRelations = []Relation{}
	for _, relation := range db.Relations {
		if !(relation.GroupName == groupID && relation.RepositoryID == repoID) {
			newRelations = append(newRelations, relation)
		}
	}
	db.Relations = newRelations
}

/*
UnrelateRepository deletes all relations of the specified repository.
*/
func (db *Database) UnrelateRepository(repoID string) {
	var newRelations = []Relation{}
	for _, relation := range db.Relations {
		if relation.RepositoryID != repoID {
			newRelations = append(newRelations, relation)
		}
	}
	db.Relations = newRelations
}

/*
UnrelateFromGroup deletes all relations about the specified group.
*/
func (db *Database) UnrelateFromGroup(groupID string) {
	var newRelations = []Relation{}
	for _, relation := range db.Relations {
		if relation.GroupName != groupID {
			newRelations = append(newRelations, relation)
		}
	}
	db.Relations = newRelations
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
	db.UnrelateRepository(repoID)
	var newRepositories = []Repository{}
	for _, repo := range db.Repositories {
		if repo.ID != repoID {
			newRepositories = append(newRepositories, repo)
		}
	}
	db.Repositories = newRepositories

	return nil
}

func deleteGroup(db *Database, groupID string) error {
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
	var groups = groupFrequencies(db)
	if groups[groupID] != 0 {
		return fmt.Errorf("%s: group has %d relatins", groupID, groups[groupID])
	}
	return deleteGroup(db, groupID)
}

/*
ForceDeleteGroup removes the group of the given groupID from DB.
Even if the group has some repositories, the function forcely remove the group.
*/
func (db *Database) ForceDeleteGroup(groupID string) error {
	if !db.HasGroup(groupID) {
		return fmt.Errorf("%s: group not found", groupID)
	}
	db.UnrelateFromGroup(groupID)
	return deleteGroup(db, groupID)
}

func databasePath(config *Config) string {
	return config.GetValue(RrhDatabasePath)
}

/*
StoreAndClose stores the database to file and close the database.
The database path is defined in RRH_DATABASE_PATH of config.
*/
func (db *Database) StoreAndClose() error {
	db.Timestamp = Now()
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
	var db = Database{Unix(0, 0), []Repository{}, []Group{}, []Relation{}, config}
	if err != nil {
		return &db, nil
	}

	if err := json.Unmarshal(bytes, &db); err != nil {
		return nil, err
	}
	db.Config = config
	return &db, nil
}
