package common

import (
	"testing"

	"github.com/tamada/rrh"
)

func TestRepositoryEntriesTest(t *testing.T) {
	testdata := []struct {
		args    []string
		iserror bool
		array   []string
	}{
		{[]string{"all"}, false, []string{"id", "description", "path", "remotes", "groups"}},
		{[]string{"id", "desc"}, false, []string{"id", "description"}},
		{[]string{"hoge"}, true, []string{}},
		{[]string{"hoge", "fuga"}, true, []string{}},
		{[]string{"id", "desc", "path", "remote", "group", "count"}, false, []string{"id", "description", "path", "remotes", "groups", "group count"}},
		{[]string{"count"}, false, []string{"group count"}},
	}
	for _, td := range testdata {
		err := ValidateRepositoryEntries(td.args)
		if err != nil && !td.iserror || err == nil && td.iserror {
			t.Errorf("ValidateRepositoryEntries(%v) wont error: %v, but got %v", td.args, td.iserror, err)
		}
		if err != nil {
			continue
		}
		re, err := NewRepositoryEntries(td.args)
		if err == nil && td.iserror || err != nil && !td.iserror {
			t.Errorf("NewRepositoryEntries(%v): %d wont error %v, but got error %v", td.args, re, td.iserror, err == nil)
		}
		array := re.StringArray()
		for _, item := range array {
			if !rrh.FindIn(item, td.array) {
				t.Errorf("NewRepositoryEntries(%v): %d do not find %s", td.args, re, item)
			}
		}
	}
}

func TestGroupEntriesTest(t *testing.T) {
	testdata := []struct {
		args    []string
		iserror bool
		array   []string
	}{
		{[]string{"all"}, false, []string{"name", "description", "abbrev", "repositories"}},
		{[]string{"name", "desc"}, false, []string{"name", "description"}},
		{[]string{"hoge"}, true, []string{}},
		{[]string{"hoge", "fuga"}, true, []string{}},
		{[]string{"name", "desc", "abbrev", "repo", "count"}, false, []string{"name", "description", "abbrev", "repositories", "repository count"}},
	}
	for _, td := range testdata {
		err := ValidateGroupEntries(td.args)
		if err != nil && !td.iserror || err == nil && td.iserror {
			t.Errorf("ValidateGroupEntries(%v) wont error: %v, but got %v", td.args, td.iserror, err)
		}
		if err != nil {
			continue
		}
		ge, err := NewGroupEntries(td.args)
		if err == nil && td.iserror || err != nil && !td.iserror {
			t.Errorf("NewGroupEntries(%v) wont error %v, but got error %v", td.args, td.iserror, err == nil)
		}
		array := ge.StringArray()
		for _, item := range array {
			if !rrh.FindIn(item, td.array) {
				t.Errorf("NewGroupEntries(%v) do not find %s", td.args, item)
			}
		}
	}
}
