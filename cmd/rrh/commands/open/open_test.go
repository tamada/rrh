package open

import (
	"os"
	"testing"

	"github.com/tamada/rrh"
)

func TestExecute(t *testing.T) {
	var dbFile = rrh.Rollback("../../../../testdata/test_db.json", "../../../../testdata/config.json", func(config *rrh.Config, oldDB *rrh.Database) {
		cmd := New()
		cmd.SetOut(os.Stdout)
		cmd.SetArgs([]string{"not_exist_repo"})
		err := cmd.Execute()
		if err == nil || err.Error() != "not_exist_repo: repository not found" {
			t.Errorf("found not_exist_repo: %v", err)
		}
	})
	defer os.Remove(dbFile)
}

func TestConvertGitURL(t *testing.T) {
	testdata := []struct {
		giveString string
		errorFlag  bool
		wontString string
	}{
		{"git@github.com:tamada/rrh.git", false, `https://github.com/tamada/rrh`},
		{"git@github.com:tamada/rrh", false, `https://github.com/tamada/rrh`},
	}
	for _, td := range testdata {
		url, err := convertToRepositoryURL(td.giveString)
		if (err == nil) == td.errorFlag {
			t.Errorf("convertToRepositoryURL(%s) should be %v, but %v", td.giveString, td.errorFlag, err)
		}
		if url != td.wontString {
			t.Errorf("convertToRepositoryURL(%s) wont %s, but got %s", td.giveString, td.wontString, url)
		}
	}
}
