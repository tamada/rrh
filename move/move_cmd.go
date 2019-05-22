package move

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	flag "github.com/spf13/pflag"
	"github.com/tamada/rrh/common"
)

/*
Command represents a command.
*/
type Command struct {
	options *options
}

type options struct {
	inquiry bool
	verbose bool
	from    []string
	to      string
}

/*
CommandFactory returns an instance of the MoveCommand.
*/
func CommandFactory() (cli.Command, error) {
	return &Command{&options{false, false, []string{}, ""}}, nil
}

func (options *options) printIfNeeded(message string) {
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
		return 4
	}
	return 0
}

/*
the target type values
*/
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

type targets struct {
	froms []target
	to    target
}

func parseCompound(db *common.Database, types []string, original string) (target, error) {
	var groupFound = db.HasGroup(types[0])
	var repoFound = db.HasRepository(types[1])
	if !groupFound && !repoFound {
		return target{Unknown, types[0], types[1], original}, fmt.Errorf("%s: group %s and repository %s not found", original, types[0], types[1])
	} else if !groupFound && repoFound {
		return target{GroupAndRepoType, types[0], types[1], original}, fmt.Errorf("%s: group %s not found", original, types[0])
	} else if groupFound && !repoFound {
		return target{GroupAndRepoType, types[0], types[1], original}, fmt.Errorf("%s: repository %s not found", original, types[1])
	} else if !db.HasRelation(types[0], types[1]) {
		return target{GroupAndRepoType, types[0], types[1], original}, fmt.Errorf("%s and %s: no relation", types[0], types[1])
	}
	return target{GroupAndRepoType, types[0], types[1], original}, nil
}

func parseEither(db *common.Database, typeString string) (target, error) {
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

func parseType(db *common.Database, typeString string) (target, error) {
	if strings.Contains(typeString, "/") {
		var types = strings.SplitN(typeString, "/", 2)
		return parseCompound(db, types, typeString)
	}
	return parseEither(db, typeString)
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

/*
Move type values
*/
const (
	GroupToGroup = iota
	GroupsToGroup
	RepositoryToRepository
	RepositoriesToGroup
	Invalid
)

func isGroupToGroup(fromType int, toType int) bool {
	return fromType == GroupType && (toType == GroupType || toType == GroupOrRepoType)
}

func isRepositoriesToGroup(fromType int, toType int) bool {
	var flag = (toType == GroupType || toType == GroupOrRepoType)
	return fromType == GroupAndRepoType && flag ||
		fromType == RepositoryType && flag
}

func isRepositoryToRepository(fromType int, toType int) bool {
	return fromType == RepositoryType || fromType == GroupAndRepoType &&
		toType == GroupAndRepoType
}

// func isNotGroupType(fromType int, toType int) bool {
// 	return toType != GroupType && toType != GroupOrRepoType
// }

func verifyArgumentsOneToOne(db *common.Database, from target, to target) (int, error) {
	if from.targetType == Unknown {
		return Invalid, fmt.Errorf("%s: unknown type not acceptable", from.original)
	}
	if isGroupToGroup(from.targetType, to.targetType) {
		return GroupToGroup, nil
	} else if isRepositoriesToGroup(from.targetType, to.targetType) {
		return RepositoriesToGroup, nil
	} else if isRepositoryToRepository(from.targetType, to.targetType) {
		return RepositoryToRepository, nil
		//	never reach this part?
		//	} else if isNotGroupType(from.targetType, to.targetType) {
		//		return Invalid, fmt.Errorf("%s: not group", to.original)
	}
	return Invalid, fmt.Errorf("Specifying arguments did not accept")
}

func findFromTypes(froms []target) (int, error) {
	var fromTypes = []int{}
	for _, from := range froms {
		fromTypes = append(fromTypes, from.targetType)
	}
	return mergeType(fromTypes)
}

func verifyArgumentsMoreToOne(db *common.Database, froms []target, to target) (int, error) {
	if to.targetType != GroupType && to.targetType != GroupOrRepoType {
		return Invalid, fmt.Errorf("types of froms and to did not match: from: %v, to: %v (%d)", froms, to.original, to.targetType)
	}

	var fromType, err2 = findFromTypes(froms)
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

func (mv *Command) performImpl(db *common.Database, targets targets, executionType int) []error {
	switch executionType {
	case GroupToGroup:
		return mv.moveGroupToGroup(db, targets.froms[0], targets.to)
	case GroupsToGroup:
		return mv.moveGroupsToGroup(db, targets.froms, targets.to)
	case RepositoriesToGroup:
		return mv.moveRepositoriesToGroup(db, targets.froms, targets.to)
	case RepositoryToRepository:
		var err = mv.moveRepositoryToRepository(db, targets.froms[0], targets.to)
		if err != nil {
			return []error{err}
		}
	default:
		return []error{fmt.Errorf("%d: unknown execution type", executionType)}
	}
	return []error{}
}

func (mv *Command) perform(db *common.Database) int {
	var from, to = convertToTarget(db, mv.options.from, mv.options.to)
	var executionType, err = verifyArguments(db, from, to)
	if err != nil {
		return printError(db.Config, []error{err})
	}
	var list = mv.performImpl(db, targets{from, to}, executionType)
	var statusCode = printError(db.Config, list)
	if statusCode == 0 {
		db.StoreAndClose()
	}
	return statusCode
}

/*
Run performs the command.
*/
func (mv *Command) Run(args []string) int {
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

func buildFlagSet(mv *Command) (*flag.FlagSet, *options) {
	var options = options{false, false, []string{}, ""}
	flags := flag.NewFlagSet("mv", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(mv.Help()) }
	flags.BoolVarP(&options.verbose, "verbose", "v", false, "verbose mode")
	return flags, &options
}

func (mv *Command) parse(args []string) (*options, error) {
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
	mv.options = options
	return options, nil
}

/*
Help function shows the help message.
*/
func (mv *Command) Help() string {
	return `rrh mv [OPTIONS] <FROMS...> <TO>
OPTIONS
    -v, --verbose   verbose mode

ARGUMENTS
    FROMS...        specifies move from, formatted in <GROUP_NAME/REPO_ID>, or <GROUP_NAME>
    TO              specifies move to, formatted in <GROUP_NAME>`
}

/*
Synopsis returns the help message of the command.
*/
func (mv *Command) Synopsis() string {
	return "move the repositories from groups to another group."
}
