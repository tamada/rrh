package group

import (
	"testing"

	"github.com/tamada/rrh"
)

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
		err := ValidateEntries(td.args)
		if err != nil && !td.iserror || err == nil && td.iserror {
			t.Errorf("ValidateGroupEntries(%v) wont error: %v, but got %v", td.args, td.iserror, err)
		}
		if err != nil {
			continue
		}
		ge, err := NewEntries(td.args)
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
