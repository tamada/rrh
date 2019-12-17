package internal

import (
	"strings"
	"testing"

	"github.com/tamada/rrh/lib"
)

const defaultHelpMessage = `rrh [GLOBAL OPTIONS] <SUB COMMANDS> [ARGUMENTS]
GLOBAL OPTIONS
    -h, --help                        print this message.
    -v, --version                     print version.
    -c, --config-file <CONFIG_FILE>   specifies the config file path.
AVAILABLE SUB COMMANDS:
    add          add repositories on the local path to rrh.
    clone        run "git clone" and register it to a group.
    config       set/unset and list configuration of RRH.
    export       export rrh database to stdout.
    fetch        run "git fetch" on the given groups.
    fetch-all    run "git fetch" in the all repositories.
    group        add/list/update/remove groups and show groups of the repository.
    help         print this message.
    import       import the given database.
    list         print managed repositories and their groups.
    mv           move the repositories from groups to another group.
    prune        prune unnecessary repositories and groups.
    repository   manages repositories.
    rm           remove given repository from database.
    status       show git status of repositories.
    version      show version.`

func TestGenerateHelpMessage(t *testing.T) {
	if defaultHelpMessage != GenerateDefaultHelp() {
		t.Errorf("generated help message did not match.")
	}
}

func TestHelpCommand(t *testing.T) {
	var testcases = []struct {
		args        []string
		wontMessage string
	}{
		{[]string{}, defaultHelpMessage},
		{[]string{"prune"}, "rrh prune"},
		{[]string{"unknown_subcommand"}, "unknown_subcommand: subcommand not found"},
	}
	for _, tc := range testcases {
		var command, _ = HelpCommandFactory()
		var message = lib.CaptureStdout(func() {
			command.Run(tc.args)
		})
		message = strings.TrimSpace(message)
		if message != tc.wontMessage {
			t.Errorf("%v: result did not match, wont: %s, got: %s", tc.args, tc.wontMessage, message)
		}
	}
}

func TestHelpOfHelpAndVersionCommand(t *testing.T) {
	var helpCommand, _ = HelpCommandFactory()
	var helpMessage = `rrh help [ARGUMENTS...]
ARGUMENTS
    print help message of target command.`
	if helpCommand.Help() != helpMessage {
		t.Errorf("help message of help command did not match.")
	}
	var versionCommand, _ = VersionCommandFactory()
	if versionCommand.Help() != "rrh version" {
		t.Errorf("help message of version command did not match.")
	}
}

func ExampleVersionCommand_Run() {
	var command, _ = VersionCommandFactory()
	command.Run([]string{})
	// Output:
	// rrh version 1.1.0
}
