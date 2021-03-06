package internal

import (
	"os"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/tamada/rrh"
)

func Example_repositoryUpdateRemotesCommand_Run() {
	var dbFile = rrh.Rollback("../testdata/remotes.json", "../testdata/config.json", func(config *rrh.Config, oldDB *rrh.Database) {
		var command, _ = repositoryUpdateRemotesCommandFactory()
		command.Run([]string{"--verbose", "--dry-run"})
	})
	defer os.Remove(dbFile)
	// Output:
	// { origin:git@github.com:tamada/fibonacci.git } -> { origin:https://htamada@bitbucket.org/htamada/fibonacci.git }
	// {  } -> { origin:https://htamada@bitbucket.org/htamada/helloworld.git }
}

func TestRepository(t *testing.T) {
	var testcases = []struct {
		args         []string
		status       int
		output       string
		ignoreOutput bool
	}{
		{[]string{}, 0, "rrh repository <SUBCOMMAND>+SUBCOMMAND+    info [OPTIONS] <REPO...>     shows repository information.+    update [OPTIONS] <REPO...>   updates repository information.+    update-remotes [OPTIONS]     updates remotes of all repositories.", false},
		{[]string{"unknown-command"}, 127, "", true},
		{[]string{"list"}, 0, "", false},
		{[]string{"list", "--id"}, 0, "repo1+repo2", false},
		{[]string{"list", "--path", "repo2"}, 0, "path2", false},
		{[]string{"list", "--with-group", "repo1"}, 0, "group1/repo1", false},
	}
	for _, tc := range testcases {
		var dbFile = rrh.Rollback("../testdata/test_db.json", "../testdata/config.json", func(config *rrh.Config, oldDB *rrh.Database) {
			var output = rrh.CaptureStdout(func() {
				var command, _ = RepositoryCommandFactory()
				var status = command.Run(tc.args)
				if status != tc.status {
					t.Errorf("%v: status code did not match, wont: %d, got: %d", tc.args, tc.status, status)
				}
			})
			if !tc.ignoreOutput {
				output = strings.TrimSpace(output)
				output = rrh.ReplaceNewline(output, "+")
				if output != tc.output {
					t.Errorf("%v: output did not match, wont: %s, got: %s", tc.args, tc.output, output)
				}
			}
		})
		defer os.Remove(dbFile)
	}
}

func TestListRepository(t *testing.T) {
	var testcases = []struct {
		args         []string
		status       int
		output       string
		ignoreOutput bool
	}{
		{[]string{"--id"}, 0, "repo1+repo2", false},
		{[]string{"--path"}, 0, "path1+path2", false},
		{[]string{"--with-group"}, 0, "group1/repo1+group3/repo2", false},
		{[]string{"--id", "repo2"}, 0, "repo2", false},
		{[]string{"--path", "repo1"}, 0, "path1", false},
		{[]string{"--with-group", "repo2"}, 0, "group3/repo2", false},
		{[]string{}, 0, "", false},
		{[]string{"--invalid-option"}, 1, "", true},
	}
	for _, tc := range testcases {
		var dbFile = rrh.Rollback("../testdata/test_db.json", "../testdata/config.json", func(config *rrh.Config, oldDB *rrh.Database) {
			var output = rrh.CaptureStdout(func() {
				var listCommand, _ = repositoryListCommandFactory()
				var status = listCommand.Run(tc.args)
				if status != tc.status {
					t.Errorf("%v: status code did not match, wont: %d, got: %d", tc.args, tc.status, status)
				}
			})
			if !tc.ignoreOutput {
				output = strings.TrimSpace(output)
				output = rrh.ReplaceNewline(output, "+")
				if output != tc.output {
					t.Errorf("%v: output did not match, wont: %s, got: %s", tc.args, tc.output, output)
				}
			}
		})
		defer os.Remove(dbFile)
	}
}

func TestInfoRepository(t *testing.T) {
	var testcases = []struct {
		args         []string
		status       int
		output       string
		ignoreOutput bool
	}{
		{[]string{"--csv", "repo1"}, 0, "repo1,,path1", false},
		{[]string{}, 1, "missing arguments", false},
		{[]string{"repo1"}, 0, `ID:          repo1+Groups:      group1+Description: +Path:        path1`, false},
		{[]string{"--color", "repo2"}, 0, `ID:          repo2+Groups:      group3+Description: +Path:        path2+Remote:     +    origin: git@github.com:example/repo2.git`, false},
		{[]string{"--invalid-option"}, 1, "", true},
		{[]string{"repo4"}, 0, "repo4: repository not found", false},
	}

	for _, tc := range testcases {
		var dbFile = rrh.Rollback("../testdata/test_db.json", "../testdata/config.json", func(config *rrh.Config, oldDB *rrh.Database) {
			var output = rrh.CaptureStdout(func() {
				var infoCommand, _ = repositoryInfoCommandFactory()
				var status = infoCommand.Run(tc.args)
				if status != tc.status {
					t.Errorf("%v: status code did not match, wont: %d, got: %d", tc.args, tc.status, status)
				}
			})
			if !tc.ignoreOutput {
				output = strings.TrimSpace(output)
				output = rrh.ReplaceNewline(output, "+")
				if output != tc.output {
					t.Errorf("%v: result did not match, wont: \"%s\", got: \"%s\"", tc.args, tc.output, output)
				}
			}
		})
		defer os.Remove(dbFile)
	}
}

func TestUpdateRepository(t *testing.T) {
	var testcases = []struct {
		args       []string
		statusCode int
		newRepoID  string
		wontRepo   *rrh.Repository
	}{
		{[]string{"--id", "newRepo1", "--path", "newPath1", "--desc", "desc1", "repo1"}, 0, "newRepo1", &rrh.Repository{ID: "newRepo1", Description: "desc1", Path: "newPath1"}},
		{[]string{"-d", "desc2", "repo2"}, 0, "repo2", &rrh.Repository{ID: "repo2", Description: "desc2", Path: "path2"}},
		{[]string{"repo4"}, 3, "repo4", nil},                             // unknown repository
		{[]string{"--invalid-option"}, 1, "never used", nil},             // invalid option
		{[]string{}, 1, "never used", nil},                               // missing arguments.
		{[]string{"-d", "desc", "repo1", "repo3"}, 1, "never used", nil}, // too many arguments.
	}

	for _, tc := range testcases {
		var dbFile = rrh.Rollback("../testdata/test_db.json", "../testdata/config.json", func(config *rrh.Config, oldDB *rrh.Database) {
			var updateCommand, _ = repositoryUpdateCommandFactory()
			var status = updateCommand.Run(tc.args)
			if status != tc.statusCode {
				t.Errorf("%v: status code did not match, wont: %d, got: %d", tc.args, tc.statusCode, status)
			}
			if status != 0 {
				return
			}
			var db, _ = rrh.Open(config)
			var repo = db.FindRepository(tc.newRepoID)
			if repo == nil {
				t.Errorf("%s: new repository do not found", tc.newRepoID)
				return
			}
			if repo.ID != tc.wontRepo.ID {
				t.Errorf("%v: id did not match: wont: %s, got: %s", tc.args, tc.wontRepo.ID, repo.ID)
			}
			if repo.Path != tc.wontRepo.Path {
				t.Errorf("%v: path did not match: wont: %s, got: %s", tc.args, tc.wontRepo.Path, repo.Path)
			}
			if repo.Description != tc.wontRepo.Description {
				t.Errorf("%v: description did not match: wont: %s, got: %s", tc.args, tc.wontRepo.Description, repo.Description)
			}
		})
		defer os.Remove(dbFile)
	}
}

func TestHelpOfRepository(t *testing.T) {
	var commandHelp = `rrh repository <SUBCOMMAND>
SUBCOMMAND
    info [OPTIONS] <REPO...>     shows repository information.
    update [OPTIONS] <REPO...>   updates repository information.
    update-remotes [OPTIONS]     updates remotes of all repositories.`

	var infoCommandHelp = `rrh repository info [OPTIONS] [REPOSITORIES...]
    -G, --color     prints the results with color.
    -c, --csv       prints the results in the csv format.
ARGUMENTS
    REPOSITORIES    target repositories.  If no repositories are specified,
                    this sub command failed.`

	var listCommandHelp = `rrh repository list [OPTIONS] [ARGUMENTS...]
OPTIONS
    --id            prints ids in the results.
    --path          prints paths in the results.
    --with-group    prints the results in "GROUP/REPOSITORY" format.
Note:
    This sub command is used for a completion target generation.`

	var updateCommandHelp = `rrh repository update [OPTIONS] <REPOSITORY>
OPTIONS
    -i, --id <NEWID>     specifies new repository id.
    -d, --desc <DESC>    specifies new description.
    -p, --path <PATH>    specifies new path.
ARGUMENTS
    REPOSITORY           specifies the repository id.`
	var updateRemotesCommandHelp = `rrh repository update-remotes [OPTIONS]
OPTIONS
    -d, --dry-run    dry-run mode.
    -v, --verbose    verbose mode.`

	var infoCommand, _ = repositoryInfoCommandFactory()
	var listCommand, _ = repositoryListCommandFactory()
	var updateCommand, _ = repositoryUpdateCommandFactory()
	var updateRemotesCommand, _ = repositoryUpdateRemotesCommandFactory()
	var command, _ = RepositoryCommandFactory()

	if infoCommand.Help() != infoCommandHelp {
		t.Errorf("infoCommand help did not match")
	}
	if listCommand.Help() != listCommandHelp {
		t.Errorf("listCommand help did not match")
	}
	if updateCommand.Help() != updateCommandHelp {
		t.Errorf("updateCommand help did not match")
	}
	if updateRemotesCommand.Help() != updateRemotesCommandHelp {
		t.Errorf("updateCommand help did not match")
	}
	if command.Help() != commandHelp {
		t.Errorf("command help did not match")
	}
}

func TestSynopsisOfRepository(t *testing.T) {
	var infoCommand, _ = repositoryInfoCommandFactory()
	var listCommand, _ = repositoryListCommandFactory()
	var updateCommand, _ = repositoryUpdateCommandFactory()
	var updateRemotesCommand, _ = repositoryUpdateRemotesCommandFactory()
	var command, _ = RepositoryCommandFactory()

	if infoCommand.Synopsis() != "prints information of the specified repositories." {
		t.Errorf("infoCommand synopsis did not match")
	}
	if listCommand.Synopsis() != "lists repositories." {
		t.Errorf("listCommand synopsis did not match")
	}
	if updateCommand.Synopsis() != "update information of the specified repository." {
		t.Errorf("updateCommand synopsis did not match")
	}
	if updateRemotesCommand.Synopsis() != "update remote entries of all repositories." {
		t.Errorf("updateCommand synopsis did not match")
	}
	if command.Synopsis() != "manages repositories." {
		t.Errorf("command synopsis did not match")
	}
}

func TestRepositoryCommandRunFailedByBrokenDBFile(t *testing.T) {
	os.Setenv(rrh.DatabasePath, "../testdata/broken.json")

	var testcases = []struct {
		comGenerator func() (cli.Command, error)
		args         []string
		statusCode   int
	}{
		{repositoryInfoCommandFactory, []string{"group1"}, 2},
		{repositoryListCommandFactory, []string{""}, 2},
		{repositoryUpdateCommandFactory, []string{""}, 2},
	}
	for _, tc := range testcases {
		var com, _ = tc.comGenerator()
		var status = com.Run(tc.args)
		if status != tc.statusCode {
			t.Errorf("%v status code did not match, wont: %d, got: %d", tc.args, tc.statusCode, status)
		}
	}
}
