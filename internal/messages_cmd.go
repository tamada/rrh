package internal

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/tamada/rrh/lib"
)

/*
BuildCommandFactoryMap builds a map of CommandFactories of rrh commands.
*/
func BuildCommandFactoryMap() map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"add":        AddCommandFactory,
		"clone":      CloneCommandFactory,
		"config":     ConfigCommandFactory,
		"export":     ExportCommandFactory,
		"fetch":      FetchCommandFactory,
		"fetch-all":  FetchAllCommandFactory,
		"group":      GroupCommandFactory,
		"help":       HelpCommandFactory,
		"import":     ImportCommandFactory,
		"list":       ListCommandFactory,
		"mv":         MoveCommandFactory,
		"open":       OpenCommandFactory,
		"prune":      PruneCommandFactory,
		"repository": RepositoryCommandFactory,
		"rm":         RemoveCommandFactory,
		"version":    VersionCommandFactory,
		"status":     StatusCommandFactory,
	}
}

/*
HelpCommand shows the struct for help command.
*/
type HelpCommand struct {
}

/*
VersionCommand shows the struct for version command.
*/
type VersionCommand struct {
}

/*
GenerateDefaultHelp generates the help message string.
*/
func GenerateDefaultHelp() string {
	var commands = BuildCommandFactoryMap()
	var maxLength = findMaxLength(commands)
	var messages = convertToHelpMessage(commands, maxLength)
	var preface = `rrh [GLOBAL OPTIONS] <SUB COMMANDS> [ARGUMENTS]
GLOBAL OPTIONS
    -h, --help                        print this message.
    -v, --version                     print version.
    -c, --config-file <CONFIG_FILE>   specifies the config file path.
AVAILABLE SUB COMMANDS:`
	// insert preface into the first element of messages
	messages = append([]string{preface}, messages...)
	return strings.Join(messages, "\n")
}

const formatterString = "    %%-%ds   %%s"

func convertToHelpMessage(commands map[string]cli.CommandFactory, max int) []string {
	var results = []string{}
	var formatter = fmt.Sprintf(formatterString, max)
	for key, value := range commands {
		var factory, _ = value()
		var synopsis = factory.Synopsis()
		results = append(results, fmt.Sprintf(formatter, key, synopsis))
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i] < results[j]
	})
	return results
}

func findMaxLength(commands map[string]cli.CommandFactory) int {
	var maxLength = 0
	for key := range commands {
		if len(key) > maxLength {
			maxLength = len(key)
		}
	}
	return maxLength
}

func printHelpOfGivenCommands(args []string) {
	var commands = BuildCommandFactoryMap()
	for _, arg := range args {
		var value = commands[arg]
		if value == nil {
			fmt.Printf("%s: subcommand not found\n", arg)
		} else {
			var com, _ = value()
			fmt.Println(com.Help())
		}
	}
}

/*
Run performs the command.
*/
func (help *HelpCommand) Run(args []string) int {
	if len(args) == 0 {
		fmt.Println(GenerateDefaultHelp())
	} else {
		printHelpOfGivenCommands(args)
	}
	return 0
}

/*
Run performs the command.
*/
func (version *VersionCommand) Run(args []string) int {
	fmt.Printf("rrh version %s\n", lib.VERSION)
	return 0
}

/*
HelpCommandFactory returns an instance of the HelpCommand.
*/
func HelpCommandFactory() (cli.Command, error) {
	return &HelpCommand{}, nil
}

/*
VersionCommandFactory returns an instance of the VersionCommand.
*/
func VersionCommandFactory() (cli.Command, error) {
	return &VersionCommand{}, nil
}

/*
Synopsis returns the help message of the command.
*/
func (help *HelpCommand) Synopsis() string {
	return "print this message."
}

/*
Help returns the help message.
*/
func (help *HelpCommand) Help() string {
	return `rrh help [ARGUMENTS...]
ARGUMENTS
    print help message of target command.`
}

/*
Synopsis returns the help message of the command.
*/
func (version *VersionCommand) Synopsis() string {
	return "show version."
}

/*
Help returns the help message.
*/
func (version *VersionCommand) Help() string {
	return `rrh version`
}
