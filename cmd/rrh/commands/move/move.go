package move

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
	"github.com/tamada/rrh/cmd/rrh/commands/common"
)

type moveOptions struct {
	dryRunFlag bool
}

var moveOpts = &moveOptions{}

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mv <FROM...> <TO>",
		Short: "move the repositories from groups to another group",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(c *cobra.Command, args []string) error {
			return common.PerformRrhCommand(c, args, performMove)
		},
	}
	cmd.Flags().BoolVarP(&moveOpts.dryRunFlag, "dry-run", "D", false, "dry-run mode")
	return cmd
}

type targetKind int

const (
	GroupType targetKind = iota + 1
	RepositoryType
	GroupAndRepoType
	GroupOrRepoType
	Unknown
)

type target struct {
	kind           targetKind
	groupName      string
	repositoryName string
	original       string
}

type targets struct {
	froms []target
	to    target
}

func parseCompound(db *rrh.Database, types []string, original string) (target, error) {
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

func parseEither(db *rrh.Database, typeString string) (target, error) {
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

func parseType(db *rrh.Database, typeString string) (target, error) {
	if strings.Contains(typeString, "/") {
		var types = strings.SplitN(typeString, "/", 2)
		return parseCompound(db, types, typeString)
	}
	return parseEither(db, typeString)
}

func mergeType(types []targetKind) (targetKind, error) {
	var t = types[0]
	for _, target := range types {
		if t != target {
			return Unknown, fmt.Errorf("types of froms contain the different types: %v", types)
		}
	}
	return t, nil
}

type executionType int

/*
Move type values
*/
const (
	GroupToGroup executionType = iota + 1
	GroupsToGroup
	RepositoryToRepository
	RepositoriesToGroup
	Invalid
)

func isGroupToGroup(fromType targetKind, toType targetKind) bool {
	return fromType == GroupType && (toType == GroupType || toType == GroupOrRepoType)
}

func isRepositoriesToGroup(fromType targetKind, toType targetKind) bool {
	var flag = (toType == GroupType || toType == GroupOrRepoType)
	return fromType == GroupAndRepoType && flag ||
		fromType == RepositoryType && flag
}

func isRepositoryToRepository(fromType targetKind, toType targetKind) bool {
	return fromType == RepositoryType || fromType == GroupAndRepoType &&
		toType == GroupAndRepoType
}

func verifyArgumentsOneToOne(db *rrh.Database, from target, to target) (executionType, error) {
	if from.kind == Unknown {
		return Invalid, fmt.Errorf("%s: unknown type not acceptable", from.original)
	}
	if isGroupToGroup(from.kind, to.kind) {
		return GroupToGroup, nil
	} else if isRepositoriesToGroup(from.kind, to.kind) {
		return RepositoriesToGroup, nil
	} else if isRepositoryToRepository(from.kind, to.kind) {
		return RepositoryToRepository, nil
		//	never reach this part?
		//	} else if isNotGroupType(from.targetType, to.targetType) {
		//		return Invalid, fmt.Errorf("%s: not group", to.original)
	}
	return Invalid, fmt.Errorf("invalid arguments")
}

func findFromTypes(froms []target) (targetKind, error) {
	var fromTypes = []targetKind{}
	for _, from := range froms {
		fromTypes = append(fromTypes, from.kind)
	}
	return mergeType(fromTypes)
}

func isNotGroupAndGroupOrRepoType(kind targetKind) bool {
	return kind != GroupType && kind != GroupOrRepoType
}

func isGroupAndRepoOrRepoType(kind targetKind) bool {
	return kind == GroupAndRepoType || kind == RepositoryType
}

func verifyArgumentsMoreToOne(db *rrh.Database, froms []target, to target) (executionType, error) {
	if isNotGroupAndGroupOrRepoType(to.kind) {
		return Invalid, fmt.Errorf("types of froms and to did not match: from: %v, to: %v (%d)", froms, to.original, to.kind)
	}

	var fromType, err2 = findFromTypes(froms)
	if err2 != nil {
		return Invalid, err2
	}
	if isGroupAndRepoOrRepoType(fromType) {
		return RepositoriesToGroup, nil
	}
	return GroupsToGroup, nil
}

func verifyArguments(db *rrh.Database, froms []target, to target) (executionType, error) {
	if len(froms) == 1 {
		return verifyArgumentsOneToOne(db, froms[0], to)
	}
	return verifyArgumentsMoreToOne(db, froms, to)
}

func convertToTarget(db *rrh.Database, froms []string, to string) ([]target, target) {
	targetFrom := []target{}
	for _, from := range froms {
		f, _ := parseType(db, from)
		targetFrom = append(targetFrom, f)
	}
	targetTo, _ := parseType(db, to)
	return targetFrom, targetTo
}

func moveRepositoryToRepository(db *rrh.Database, from target, to target) error {
	if from.repositoryName != to.repositoryName {
		return fmt.Errorf("repository name did not match: %s, %s", from.original, to.original)
	}
	if _, err := db.AutoCreateGroup(to.groupName, "", false); err != nil {
		return err
	}
	if from.kind == GroupAndRepoType {
		db.Unrelate(from.groupName, from.repositoryName)
	}
	db.Relate(to.groupName, to.repositoryName)
	return nil
}

func moveRepositoryToGroup(db *rrh.Database, from target, to target) error {
	if to.kind == GroupType || to.kind == GroupOrRepoType {
		if _, err := db.AutoCreateGroup(to.original, "", false); err != nil {
			return err
		}
	}
	if from.kind == GroupAndRepoType {
		db.Unrelate(from.groupName, from.repositoryName)
	}
	db.Relate(to.original, from.repositoryName)
	return nil
}

func moveRepositoriesToGroup(db *rrh.Database, froms []target, to target) error {
	el := common.NewErrorList()
	for _, from := range froms {
		err := moveRepositoryToGroup(db, from, to)
		el = el.Append(err)
	}
	return el.NilOrThis()
}

func moveGroupsToGroup(db *rrh.Database, froms []target, to target) error {
	el := common.NewErrorList()
	for _, from := range froms {
		err := moveGroupToGroup(db, from, to)
		el = el.Append(err)
	}
	return el.NilOrThis()
}

func moveGroupToGroup(db *rrh.Database, from target, to target) error {
	el := common.NewErrorList()
	if _, err := db.AutoCreateGroup(to.groupName, "", false); err != nil {
		return err
	}
	var repos = db.FindRelationsOfGroup(from.groupName)
	for _, repo := range repos {
		db.Unrelate(from.groupName, repo)
		db.Relate(to.groupName, repo)
	}
	return el.NilOrThis()
}

func performImpl(db *rrh.Database, targets targets, execType executionType) error {
	switch execType {
	case GroupToGroup:
		return moveGroupToGroup(db, targets.froms[0], targets.to)
	case GroupsToGroup:
		return moveGroupsToGroup(db, targets.froms, targets.to)
	case RepositoriesToGroup:
		return moveRepositoriesToGroup(db, targets.froms, targets.to)
	case RepositoryToRepository:
		return moveRepositoryToRepository(db, targets.froms[0], targets.to)
	default:
		return fmt.Errorf("%d: unknown execution type", execType)
	}
}

func performAndStore(db *rrh.Database, targets targets, execType executionType) error {
	err := performImpl(db, targets, execType)
	if err != nil {
		return err
	}
	if !moveOpts.dryRunFlag {
		db.StoreAndClose()
	}
	return nil
}

func performMove(c *cobra.Command, args []string, db *rrh.Database) error {
	len := len(args) - 1
	froms, to := convertToTarget(db, args[:len], args[len])
	execType, err := verifyArguments(db, froms, to)
	if err != nil {
		return err
	}
	return performAndStore(db, targets{froms, to}, execType)
}
