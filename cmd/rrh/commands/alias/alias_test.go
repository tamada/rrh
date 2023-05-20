package alias

import (
	"os"
	"testing"

	"github.com/tamada/rrh"
)

func TestValidateArguments(t *testing.T) {
	testdata := []struct {
		giveArgs  []string
		wontError bool
	}{
		{[]string{"ls"}, false},
		{[]string{"ls", "|", "wc"}, false},
		{[]string{"ls", ">", "list.txt"}, true},
	}
	for _, td := range testdata {
		err := validateWithoutRedirectSymbol(td.giveArgs)
		if err == nil && td.wontError || err != nil && !td.wontError {
			t.Errorf("validateWithoutRedirectSymbol(%v) wont error %v, but got %v", td.giveArgs, td.wontError, err)
		}
	}
}

func Example_ListAlias() {
	os.Setenv(rrh.AliasPath, "../../../../testdata/alias.json")
	cmd := New()
	cmd.SetOut(os.Stdout)
	cmd.SetArgs([]string{})
	cmd.Execute()
	// Output:
	// grlist=repository list --entry group,id
}

func TestRemoveAlias(t *testing.T) {
	dbFile := rrh.RollbackAlias("../../../../testdata/database.json", "../../../../testdata/config.json", "../../../../testdata/alias.json", func(config *rrh.Config, oldDB *rrh.Database) {
		cmd := New()
		cmd.SetArgs([]string{"--remove", "grlist"})
		err := cmd.Execute()
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		commands, err := LoadAliases(config)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
		if len(commands) > 0 {
			t.Errorf("alias list wont 0, but got %d (%v)", len(commands), commands[0])
		}
	})
	defer os.Remove(dbFile)
}

func TestRegisterNewAlias(t *testing.T) {
	dbFile := rrh.RollbackAlias("../../../../testdata/database.json", "../../../../testdata/config.json", "../../../../testdata/alias.json", func(config *rrh.Config, oldDB *rrh.Database) {
		cmd := New()
		cmd.SetArgs([]string{"gl", "group", "list"})
		err := cmd.Execute()
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
		commands, err := LoadAliases(config)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
		if len(commands) != 2 {
			t.Errorf("alias list wont 2, but got %d (%v)", len(commands), commands[0])
		}
	})
	defer os.Remove(dbFile)
}
