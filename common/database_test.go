package common

import (
	"os"
	"testing"
)

func openDatabase() *Database {
	os.Setenv(RrhDatabasePath, "../testdata/database.json")
	os.Setenv(RrhConfigPath, "../testdata/config.json")
	var config = OpenConfig()
	var db, _ = Open(config)
	return db
}

func TestOpenBrokenJson(t *testing.T) {
	os.Setenv(RrhDatabasePath, "../testdata/broken.json")
	var config = OpenConfig()
	var _, err = Open(config)
	if err == nil {
		t.Error("it can read broken json!?")
	}
}

func TestAutoCreateGroup(t *testing.T) {
	var db = openDatabase()
	var group, err = db.AutoCreateGroup("newgroup", "desc", true)
	if err != nil {
		t.Errorf("auto create group failed: %s", err.Error())
	}
	if !db.HasGroup("newgroup") && group.Description != "desc" && group.OmitList {
		t.Errorf("auto created group did not match, wont: %v, got: %v", Group{"newgroup", "desc", true}, group)
	}
	var group2, err2 = db.AutoCreateGroup("no-group", "description1", false)
	if err2 != nil {
		t.Errorf("auto create group failed: %s", err.Error())
	}
	if !db.HasGroup("no-group") && group.Description != "" && group.OmitList {
		t.Errorf("existing group did not match, wont: %v, got: %v", Group{"group1", "desc1", true}, group2)
	}

	db.Config.Update(RrhAutoCreateGroup, "false")
	var _, err3 = db.AutoCreateGroup("failgroup", "desc", true)
	if err3 == nil {
		t.Errorf("auto create group should fail: %s", err3.Error())
	}
}

func TestOpenNonExistFile(t *testing.T) {
	os.Setenv(RrhDatabasePath, "../testdata/not-exist-file.json")
	var config = OpenConfig()
	var db, _ = Open(config)

	if len(db.Repositories) != 0 {
		t.Error("null db have no repository entries")
	}
	if len(db.Groups) != 0 {
		t.Error("null db have no group entries")
	}
}

func TestOpenNullDatabase(t *testing.T) {
	os.Setenv(RrhDatabasePath, "../testdata/nulldb.json")
	var config = OpenConfig()
	var db, _ = Open(config)

	if len(db.Repositories) != 0 {
		t.Error("null db have no repository entries")
	}
	if len(db.Groups) != 0 {
		t.Error("null db have no group entries")
	}
	if len(db.Relations) != 0 {
		t.Error("null db have no group entries")
	}
}

func TestStore(t *testing.T) {
	Rollback("../testdata/tmp.json", "../testdata/config.json", func() {
		var config = OpenConfig()
		var db, _ = Open(config)

		db.CreateGroup("group1", "desc1", false)
		db.CreateGroup("group2", "desc2", false)
		db.CreateRepository("repo1", "path1", "desc1", []Remote{})
		db.CreateRepository("repo2", "path2", "desc2", []Remote{})
		db.Relate("group1", "repo1")
		db.StoreAndClose()

		var db2, _ = Open(config)
		if !db2.HasGroup("group1") {
			t.Error("group1 not found!")
		}
		if !db2.HasGroup("group2") {
			t.Error("group2 not found!")
		}
		if !db2.HasRepository("repo1") {
			t.Error("repo1 not found!")
		}
		if !db2.HasRepository("repo2") {
			t.Error("repo2 not found!")
		}
		if !db2.HasRelation("group1", "repo1") {
			t.Error("group1 does not relate with repo1")
		}
	})
}

func TestPrune(t *testing.T) {
	var db = openDatabase()
	db.CreateGroup("group1", "desc1", false)
	db.CreateGroup("group2", "desc2", false)
	db.CreateRepository("repo1", "path1", "desc1", []Remote{})
	db.CreateRepository("repo2", "path2", "desc2", []Remote{})
	db.Relate("group1", "repo1")
	db.Prune()

	if db.HasGroup("group2") {
		t.Error("group2 was not pruned")
	}
	if db.HasRepository("repo2") {
		t.Error("repo2 was not pruned")
	}
	if !db.HasRelation("group1", "repo1") {
		t.Error("relation is removed?")
	}
}

func TestDeleteGroup(t *testing.T) {
	var db = openDatabase()
	db.CreateGroup("group1", "desc1", false)
	db.CreateGroup("group2", "desc2", false)
	db.CreateRepository("repo1", "path1", "desc1", []Remote{})
	db.CreateRepository("repo2", "path2", "desc2", []Remote{})
	db.Relate("group1", "repo1")

	if err := db.DeleteGroup("unknown"); err == nil {
		t.Error("uknown: group found!")
	}
	if err := db.DeleteGroup("group2"); err != nil {
		t.Error(err.Error())
	}
	if err := db.DeleteGroup("group1"); err == nil {
		t.Error("group1 has one relation.")
	}
	if err := db.ForceDeleteGroup("group1"); err != nil {
		t.Error(err.Error())
	}
	if err := db.ForceDeleteGroup("unknown"); err == nil {
		t.Error("uknown: group found!")
	}
}

func TestDeleteRepository(t *testing.T) {
	var db = openDatabase()
	db.CreateRepository("repo1", "path1", "desc1", []Remote{})
	db.CreateRepository("repo2", "path2", "desc2", []Remote{})
	if err := db.DeleteRepository("unknown"); err == nil {
		t.Error("unknown: repository found!")
	}

	if err := db.DeleteRepository("rrh"); err != nil {
		t.Error(err.Error())
	}
}

func TestUnrelate(t *testing.T) {
	var db = openDatabase()

	db.CreateRepository("somerepo", "unknown", "desc", []Remote{})
	db.CreateGroup("group2", "desc2", false)
	db.Relate("group2", "somerepo")
	db.Relate("no-group", "somerepo")

	db.Unrelate("group2", "Rrh")
	db.Unrelate("no-group", "rrh")

	if !db.HasRelation("group2", "somerepo") {
		t.Error("group2 does not relate with somerepo")
	}
}

func TestCreateRepository(t *testing.T) {
	var db = openDatabase()
	// rrh is already registered repository, therefore, the CreateRepository will fail.
	var r1, err1 = db.CreateRepository("rrh", "unknown", "desc", []Remote{})
	if r1 != nil && err1 == nil {
		t.Error(err1.Error())
	}

	var r2, err2 = db.CreateRepository("somerepo", "unknown", "desc", []Remote{{"name1", "url1"}, {"name2", "url2"}})
	if r2 == nil && err2 != nil {
		t.Error("somerepo: cannot create repository")
	}
	if len(r2.Remotes) != 2 {
		t.Error("remotes were not match.")
	}
	assert(t, r2.ID, "somerepo")
	assert(t, r2.Path, "unknown")
	assert(t, r2.Remotes[0].Name, "name1")
	assert(t, r2.Remotes[0].URL, "url1")
	assert(t, r2.Remotes[1].Name, "name2")
	assert(t, r2.Remotes[1].URL, "url2")
}

func TestCreateGroupRelateAndUnrelate(t *testing.T) {
	var db = openDatabase()

	var g1, err1 = db.CreateGroup("newGroup1", "desc1", false)
	if err1 != nil {
		t.Error(err1.Error())
	}
	if g1.Name != "newGroup1" {
		t.Error("the name of created group is different")
	}
	if g1.Description != "desc1" {
		t.Error("the description of created group is different")
	}

	var g2, err2 = db.CreateGroup("newGroup1", "desc2", false)
	if err2 == nil || g2 != nil {
		t.Error("cannot create same name group")
	}

	if err := db.Relate("no-group", "rrh"); err != nil {
		t.Error("existing relation was never error")
	}
	if err := db.Relate("newGroup1", "rrh"); err != nil {
		t.Error(err.Error())
	}
	if !db.HasRelation("newGroup1", "rrh") {
		t.Error("created relation was not found!")
	}

	db.Unrelate("no-group", "rrh")
	if db.HasRelation("no-group", "rrh") {
		t.Error("deleted relation was not found!")
	}

	if err := db.Relate("unknown", "rrh"); err != nil {
		t.Error(err.Error())
	}

}

func TestUpdateGroup(t *testing.T) {
	var db = openDatabase()

	db.UpdateGroup("no-group", Group{"updated-group", "description", false})
	var group = db.FindGroup("updated-group")
	if group.Name != "updated-group" {
		t.Error("Update is failed (group name was not updated)")
	}
	if group.Description != "description" {
		t.Error("Update is failed (description was not updated)")
	}

	if db.UpdateGroup("unknown", Group{"never used", "never used2", false}) {
		t.Error("unknown group is successfully updated.")
	}
}

func TestFindFunction(t *testing.T) {
	var db = openDatabase()
	var group1 = db.FindGroup("no-group")
	if group1 == nil {
		t.Error("no-group: not found")
	}
	var group2 = db.FindGroup("unknown")
	if group2 != nil {
		t.Error("unknown: found!")
	}

	var repo1 = db.FindRepository("rrh")
	var repo2 = db.FindRepository("unknown")
	if repo1 == nil {
		t.Error("rrh: not found!")
	}
	if repo2 != nil {
		t.Error("rrh: found!")
	}
}

func TestHasGroup(t *testing.T) {
	var db = openDatabase()

	if !db.HasGroup("no-group") {
		t.Error("no-group: group not found")
	}
	if db.HasGroup("unknown") {
		t.Error("unknown: exist not existing group")
	}

	if !db.HasRepository("rrh") {
		t.Error("rrh: repository not found")
	}
	if db.HasRepository("unknown") {
		t.Error("unknown: found!")
	}

	if !db.HasRelation("no-group", "rrh") {
		t.Error("rrh: no relation with no-group")
	}
	if db.HasRelation("unknown", "rrh") {
		t.Error("rrh: unknown relation found!")
	}
	if db.HasRelation("no-group", "unknown") {
		t.Error("unknown relation in no-group found!")
	}
}

func TestFindRelations(t *testing.T) {
	var db = openDatabase()

	var repos = db.FindRelationsOfGroup("no-group")
	if len(repos) != 1 || repos[0] != "rrh" {
		t.Errorf("find relations: wont: [\"rrh\"], got: %v", repos)
	}
	var groups = db.FindRelationsOfRepository("rrh")
	if len(groups) != 1 || groups[0] != "no-group" {
		t.Errorf("find relations: wont: [\"no-group\"], got: %v", groups)
	}
}

func TestUpdateRepository(t *testing.T) {
	Rollback("../testdata/tmp.json", "../testdata/config.json", func() {
		var config = OpenConfig()
		var db, _ = Open(config)

		if !db.UpdateRepository("repo1", Repository{ID: "newRepo1", Path: "newPath1", Description: "desc1"}) {
			t.Errorf("Update failed")
		}
		if !db.HasRepository("newRepo1") || !db.HasRelation("group1", "newRepo1") {
			t.Errorf("update repository failed.")
		}
		if db.HasRepository("repo1") || db.HasRelation("group1", "repo1") {
			t.Errorf("old information are remained.")
		}

		if db.UpdateRepository("unknownrepo", Repository{}) {
			t.Errorf("missing repository updation succeeded.")
		}
	})
}

func TestCounting(t *testing.T) {
	var db = openDatabase()

	if count := db.BelongingCount("rrh"); count != 1 {
		t.Errorf("belonging count of %s: wont: 1, got: %d", "rrh", count)
	}
	if count := db.ContainsCount("no-group"); count != 1 {
		t.Errorf("%s contains: wont: 1, got: %d", "no-group", count)
	}
}
