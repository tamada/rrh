package move

import (
	"flag"
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/tamada/rrh/common"
)

/*
PruneCommand represents a command.
*/
type MoveCommand struct {
	Options *moveOptions
}

type moveOptions struct {
	inquiry bool
	verbose bool
	from    []string
	to      string
}

/*
MoveCommandFactory returns an instance of the MoveCommand.
*/
func MoveCommandFactory() (cli.Command, error) {
	return &MoveCommand{&moveOptions{false, false, []string{}, ""}}, nil
}

func (options *moveOptions) printIfNeeded(message string) {
	if options.verbose {
		fmt.Println(message)
	}
}

func printError(config *common.Config, errs []error) int {
	var onError = config.GetValue(common.RrhOnError)
	if onError != common.Ignore {
		for _, err := range errs {
			fmt.Println(err.Error())
		}
	}
	if len(errs) > 0 && (onError == common.Fail || onError == common.FailImmediately) {
		return 1
	}
	return 0
}

const (
	GroupType = iota
	RepositoryType
	GroupAndRepoType
	GroupOrRepoType
	Unknown
)

type target struct {
	targetType     int
	groupName      string
	repositoryName string
	original       string
}

func parseType(db *common.Database, typeString string) (target, error) {
	if strings.Contains(typeString, "/") {
		var types = strings.SplitN(typeString, "/", 2)
		var groupFound = db.HasGroup(types[0])
		var repoFound = db.HasRepository(types[1])
		if !groupFound && !repoFound {
			return target{Unknown, types[0], types[1], typeString}, fmt.Errorf("%s: group %s and repository %s not found", typeString, types[0], types[1])
		} else if !groupFound && repoFound {
			return target{GroupAndRepoType, types[0], types[1], typeString}, fmt.Errorf("%s: group %s not found", typeString, types[0])
		} else if groupFound && !repoFound {
			return target{GroupAndRepoType, types[0], types[1], typeString}, fmt.Errorf("%s: repository %s not found", typeString, types[1])
		} else if !db.HasRelation(types[0], types[1]) {
			return target{GroupAndRepoType, types[0], types[1], typeString}, fmt.Errorf("%s and %s: no relation", types[0], types[1])
		}
		return target{GroupAndRepoType, types[0], types[1], typeString}, nil
	}
	var groupFound = db.HasGroup(typeString)
	var repositoryFound = db.HasRepository(typeString)

	if groupFound && repositoryFound {
		return target{Unknown, "", "", typeString}, fmt.Errorf("%s: group and repository both exist", typeString)
	} else if groupFound && !repositoryFound {
		return target{GroupType, typeString, "", typeString}, nil
	} else if !groupFound && repositoryFound {
		return target{RepositoryType, "", typeString, typeString}, nil
	}
	return target{GroupOrRepoType, typeString, "", typeString}, nil
}

func mergeType(types []int) (int, error) {
	var t = types[0]
	for _, target := range types {
		if t != target {
			return Unknown, fmt.Errorf("types of froms contain the different types: %v", types)
		}
	}
	return t, nil
}

const (
	GroupToGroup = iota
	GroupsToGroup
	RepositoryToRepository
	RepositoriesToGroup
	Invalid
)

func verifyArgumentsOneToOne(db *common.Database, from target, to target) (int, error) {
	if from.targetType == Unknown {
		return Invalid, fmt.Errorf("%s: unknown type not acceptable", from.original)
	}
	if from.targetType == GroupType && (to.targetType == GroupType || to.targetType == GroupOrRepoType) {
		return GroupToGroup, nil
	} else if from.targetType == GroupAndRepoType && (to.targetType == GroupType || to.targetType == GroupOrRepoType) {
		return RepositoriesToGroup, nil
	} else if from.targetType == RepositoryType && (to.targetType == GroupType || to.targetType == GroupOrRepoType) {
		return RepositoriesToGroup, nil
	} else if (from.targetType == RepositoryType || from.targetType == GroupAndRepoType) && to.targetType == GroupAndRepoType {
		return RepositoryToRepository, nil
	} else if to.targetType != GroupType && to.targetType != GroupOrRepoType {
		return Invalid, fmt.Errorf("%s: not group", to.original)
	}
	return Invalid, fmt.Errorf("Specifying arguments did not accept")
}

func verifyArgumentsMoreToOne(db *common.Database, froms []target, to target) (int, error) {
	if to.targetType != GroupType && to.targetType != GroupOrRepoType {
		return Invalid, fmt.Errorf("types of froms and to did not match: from: %v, to: %v (%d)", froms, to.original, to.targetType)
	}

	var fromTypes = []int{}
	for _, from := range froms {
		fromTypes = append(fromTypes, from.targetType)
	}
	var fromType, err2 = mergeType(fromTypes)
	if err2 != nil {
		return Invalid, err2
	}
	if fromType == GroupAndRepoType || fromType == RepositoryType {
		return RepositoriesToGroup, nil
	}
	return GroupsToGroup, nil
}

func verifyArguments(db *common.Database, froms []target, to target) (int, error) {
	if len(froms) == 1 {
		return verifyArgumentsOneToOne(db, froms[0], to)
	}
	return verifyArgumentsMoreToOne(db, froms, to)
}

func convertToTarget(db *common.Database, froms []string, to string) ([]target, target) {
	var targetFrom = []target{}
	for _, from := range froms {
		var f, _ = parseType(db, from)
		targetFrom = append(targetFrom, f)
	}
	var targetTo, _ = parseType(db, to)
	return targetFrom, targetTo
}

func (mv *MoveCommand) perform(db *common.Database) int {
	var list = []error{}
	var from, to = convertToTarget(db, mv.Options.from, mv.Options.to)
	var executionType, err = verifyArguments(db, from, to)
	if err != nil {
		return printError(db.Config, []error{err})
	}
	switch executionType {
	case GroupToGroup:
		var errs = mv.moveGroupToGroup(db, from[0], to)
		if len(errs) > 0 {
			list = append(list, errs...)
		}
	case GroupsToGroup:
		var errs = mv.moveGroupsToGroup(db, from, to)
		if len(errs) > 0 {
			list = append(list, errs...)
		}
	case RepositoriesToGroup:
		var errs = mv.moveRepositoriesToGroup(db, from, to)
		if len(errs) > 0 {
			list = append(list, errs...)
		}
	case RepositoryToRepository:
		var err = mv.moveRepositoryToRepository(db, from[0], to)
		if err != nil {
			list = append(list, err)
		}
	default:
		list = append(list, fmt.Errorf("%d: unknown execution type", executionType))
	}

	var statusCode = printError(db.Config, list)
	if statusCode == 0 {
		db.StoreAndClose()
	}
	return statusCode
}

/*
Run performs the command.
*/
func (mv *MoveCommand) Run(args []string) int {
	var _, err1 = mv.parse(args)
	if err1 != nil {
		fmt.Println(err1.Error())
		return 1
	}
	var config = common.OpenConfig()
	var db, err2 = common.Open(config)
	if err2 != nil {
		fmt.Println(err2.Error())
		return 2
	}

	return mv.perform(db)
}

func buildFlagSet(mv *MoveCommand) (*flag.FlagSet, *moveOptions) {
	var options = moveOptions{false, false, []string{}, ""}
	flags := flag.NewFlagSet("mv", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(mv.Help()) }
	flags.BoolVar(&options.inquiry, "i", false, "inquiry mode")
	flags.BoolVar(&options.inquiry, "inquiry", false, "inquiry mode")
	flags.BoolVar(&options.verbose, "v", false, "verbose mode")
	flags.BoolVar(&options.verbose, "verbose", false, "verbose mode")
	return flags, &options
}

func (mv *MoveCommand) parse(args []string) (*moveOptions, error) {
	var flagSet, options = buildFlagSet(mv)
	if err := flagSet.Parse(args); err != nil {
		return nil, err
	}
	var newArgs = flagSet.Args()
	if len(newArgs) < 2 {
		return nil, fmt.Errorf("too few arguments: %v", newArgs)
	}
	var len = len(newArgs) - 1
	options.from = newArgs[:len]
	options.to = newArgs[len]
	mv.Options = options
	return options, nil
}

/*
Help function shows the help message.
*/
func (mv *MoveCommand) Help() string {
	return `rrh mv [OPTIONS] <FROMS...> <TO>
OPTIONS
    -v, --verbose   verbose mode
    -i, --inquiry   inquiry mode

ARGUMENTS
    FROMS...        specifies move from, formatted in <GROUP_NAME/REPO_ID>, or <GROUP_NAME>
    TO              specifies move to, formatted in <GROUP_NAME>`
}

/*
Synopsis returns the help message of the command.
*/
func (mv *MoveCommand) Synopsis() string {
	return "move the repositories from groups to another group"
}
