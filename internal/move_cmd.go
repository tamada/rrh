package internal

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	flag "github.com/spf13/pflag"
	"github.com/tamada/rrh/lib"
)

/*
MoveCommand represents a command.
*/
type MoveCommand struct {
	options *moveOptions
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

func parseCompound(db *lib.Database, types []string, original string) (target, error) {
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

func parseEither(db *lib.Database, typeString string) (target, error) {
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

func parseType(db *lib.Database, typeString string) (target, error) {
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

func verifyArgumentsOneToOne(db *lib.Database, from target, to target) (int, error) {
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

func verifyArgumentsMoreToOne(db *lib.Database, froms []target, to target) (int, error) {
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

func verifyArguments(db *lib.Database, froms []target, to target) (int, error) {
	if len(froms) == 1 {
		return verifyArgumentsOneToOne(db, froms[0], to)
	}
	return verifyArgumentsMoreToOne(db, froms, to)
}

func convertToTarget(db *lib.Database, froms []string, to string) ([]target, target) {
	var targetFrom = []target{}
	for _, from := range froms {
		var f, _ = parseType(db, from)
		targetFrom = append(targetFrom, f)
	}
	var targetTo, _ = parseType(db, to)
	return targetFrom, targetTo
}

func (mv *MoveCommand) performImpl(db *lib.Database, targets targets, executionType int) []error {
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

func (mv *MoveCommand) moveRepositoryToRepository(db *lib.Database, from target, to target) error {
	if from.repositoryName != to.repositoryName {
		return fmt.Errorf("repository name did not match: %s, %s", from.original, to.original)
	}
	if _, err := db.AutoCreateGroup(to.groupName, "", false); err != nil {
		return err
	}
	if from.targetType == GroupAndRepoType {
		db.Unrelate(from.groupName, from.repositoryName)
		mv.options.printIfNeeded(fmt.Sprintf("unrelate group %s and repository %s", from.groupName, from.repositoryName))
	}
	db.Relate(to.groupName, to.repositoryName)
	mv.options.printIfNeeded(fmt.Sprintf("relate group %s and repository %s", to.groupName, to.repositoryName))
	return nil
}

func (mv *MoveCommand) moveRepositoryToGroup(db *lib.Database, from target, to target) error {
	if to.targetType == GroupType || to.targetType == GroupOrRepoType {
		if _, err := db.AutoCreateGroup(to.original, "", false); err != nil {
			return err
		}
	}
	if from.targetType == GroupAndRepoType {
		db.Unrelate(from.groupName, from.repositoryName)
	}
	db.Relate(to.original, from.repositoryName)
	return nil
}
func (mv *MoveCommand) moveRepositoriesToGroup(db *lib.Database, froms []target, to target) []error {
	var list = []error{}
	for _, from := range froms {
		var err = mv.moveRepositoryToGroup(db, from, to)
		if err != nil {
			if isFailImmediately(db.Config) {
				return []error{err}
			}
			list = append(list, err)
		}
	}
	return list
}

func (mv *MoveCommand) moveGroupsToGroup(db *lib.Database, froms []target, to target) []error {
	var list = []error{}
	for _, from := range froms {
		var errs = mv.moveGroupToGroup(db, from, to)
		if len(errs) != 0 {
			if isFailImmediately(db.Config) {
				return errs
			}
			list = append(list, errs...)
		}
	}
	return list
}

func (mv *MoveCommand) moveGroupToGroup(db *lib.Database, from target, to target) []error {
	if _, err := db.AutoCreateGroup(to.groupName, "", false); err != nil {
		return []error{err}
	}
	var repos = db.FindRelationsOfGroup(from.groupName)
	for _, repo := range repos {
		db.Unrelate(from.groupName, repo)
		mv.options.printIfNeeded(fmt.Sprintf("unrelate group %s and repository %s", from.groupName, repo))
		db.Relate(to.groupName, repo)
		mv.options.printIfNeeded(fmt.Sprintf("relate group %s and repository %s", to.groupName, repo))
	}
	return []error{}
}

func (mv *MoveCommand) perform(db *lib.Database) int {
	var from, to = convertToTarget(db, mv.options.from, mv.options.to)
	var executionType, err = verifyArguments(db, from, to)
	if err != nil {
		return printErrors(db.Config, []error{err})
	}
	var list = mv.performImpl(db, targets{from, to}, executionType)
	var statusCode = printErrors(db.Config, list)
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
	var config = lib.OpenConfig()
	var db, err2 = lib.Open(config)
	if err2 != nil {
		fmt.Println(err2.Error())
		return 2
	}

	return mv.perform(db)
}

func (mv *MoveCommand) buildFlagSet() (*flag.FlagSet, *moveOptions) {
	var options = moveOptions{false, false, []string{}, ""}
	flags := flag.NewFlagSet("mv", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(mv.Help()) }
	flags.BoolVarP(&options.verbose, "verbose", "v", false, "verbose mode")
	return flags, &options
}

func (mv *MoveCommand) parse(args []string) (*moveOptions, error) {
	var flagSet, options = mv.buildFlagSet()
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
func (mv *MoveCommand) Help() string {
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
func (mv *MoveCommand) Synopsis() string {
	return "move the repositories from groups to another group."
}