package repository

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
		{[]string{"all"}, false, []string{"id", "description", "path", "remote name", "remote url", "groups"}},
		{[]string{"id", "desc"}, false, []string{"id", "description"}},
		{[]string{"hoge"}, true, []string{}},
		{[]string{"hoge", "fuga"}, true, []string{}},
		{[]string{"id", "desc", "path", "remote", "group", "count"}, false, []string{"id", "description", "path", "remote name", "remote url", "groups", "group count"}},
		{[]string{"count"}, false, []string{"group count"}},
	}
	for _, td := range testdata {
		err := ValidateEntries(td.args)
		if err != nil && !td.iserror || err == nil && td.iserror {
			t.Errorf("ValidateRepositoryEntries(%v) wont error: %v, but got %v", td.args, td.iserror, err)
		}
		if err != nil {
			continue
		}
		re, err := NewEntries(td.args)
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
