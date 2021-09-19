package group

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
	"github.com/tamada/rrh/cmd/rrh/commands/common"
)

func New() *cobra.Command {
	groupCommand := &cobra.Command{
		Use:   "group [subcommand]",
		Short: "manage the groups for the rrh database",
		RunE: func(c *cobra.Command, args []string) error {
			return common.PerformRrhCommand(c, args, listGroups)
		},
	}
	registerGroupCommands(groupCommand)
	return groupCommand
}

func registerGroupCommands(c *cobra.Command) {
	c.AddCommand(createGroupAddCommand())
	c.AddCommand(createGroupInfoCommand())
	c.AddCommand(createGroupOfCommand())
	c.AddCommand(createGroupListCommand())
	c.AddCommand(createGroupUpdateCommand())
	c.AddCommand(createGroupRemoveCommand())
}

func createGroupAddCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "add <GROUP_NAME>",
		Short: "add groups to the rrh database",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return common.PerformRrhCommand(c, args, func(c *cobra.Command, args []string, db *rrh.Database) error {
				if err := createAndAddGroup(c, args[0], db); err != nil {
					return err
				}
				if !createOpts.dryRunFlag {
					return db.StoreAndClose()
				}
				c.Printf("create group Name: %s, Description: %s, Abbrev: %v\n", args[0], createOpts.desc, createOpts.abbrevFlag)
				return nil
			})
		},
	}
	flags := command.Flags()
	flags.StringVarP(&createOpts.desc, "desc", "d", "", "specifies the description")
	flags.StringVarP(&createOpts.abbrevFlag, "abbrev", "a", "false", "set the group as the abbrev")
	flags.BoolVarP(&createOpts.dryRunFlag, "dry-run", "D", false, "dry-run mode")
	return command
}

func createAndAddGroup(c *cobra.Command, groupName string, db *rrh.Database) error {
	abbrevFlag, err := strconv.ParseBool(createOpts.abbrevFlag)
	if err != nil {
		return err
	}
	_, err = db.CreateGroup(groupName, createOpts.desc, abbrevFlag)
	return err
}

type createOptions struct {
	desc       string
	abbrevFlag string
	dryRunFlag bool
}

var createOpts = &createOptions{}

func createGroupInfoCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "info <GROUP>",
		Short: "print the information of the specified group",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return common.PerformRrhCommand(c, args, printGroupInfos)
		},
	}
	return command
}

func printGroupInfos(c *cobra.Command, args []string, db *rrh.Database) error {
	for _, groupName := range args {
		err := printGroupInfo(c, groupName, db)
		if err != nil {
			return err
		}
	}
	return nil
}

func printGroupInfo(c *cobra.Command, groupName string, db *rrh.Database) error {
	group := db.FindGroup(groupName)
	if group == nil {
		return fmt.Errorf("%s: group not found", groupName)
	}
	count := db.ContainsCount(group.Name)
	unit := "repositories"
	if count == 1 {
		unit = "repository"
	}
	fmt.Printf("%s: %s (%d %s, abbrev: %v)\n", group.Name, group.Description, count, unit, group.OmitList)
	return nil
}

func createGroupOfCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "of <REPOSITORY_ID>",
		Short: "print the group name of the specified repository",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return common.PerformRrhCommand(c, args, func(c *cobra.Command, args []string, db *rrh.Database) error {
				if !db.HasRepository(args[0]) {
					return fmt.Errorf("%s: repository not found", args[0])
				}
				var groups = db.FindRelationsOfRepository(args[0])
				fmt.Printf("%s, %v\n", args[0], groups)
				return nil
			})
		},
	}
	return command
}
func createGroupListCommand() *cobra.Command {
	command := &cobra.Command{
		Use:  "list",
		Args: cobra.NoArgs,
		PreRunE: func(c *cobra.Command, args []string) error {
			return validateListOpts(listOpts)
		},
		Short: "print the group list",
		RunE: func(c *cobra.Command, args []string) error {
			return common.PerformRrhCommand(c, args, listGroups)
		},
	}
	flags := command.Flags()
	flags.StringSliceVarP(&listOpts.entries, "entry", "e", []string{"name", "count"}, "specifies the printing entries separated with comma. Available vaues: all, name, desc, repo, abbrev, count")

	return command
}

func validateListOpts(opts *listOptions) error {
	availables := []string{"all", "name", "count", "desc", "repo", "abbrev"}
	messages := []string{}
	for _, entry := range opts.entries {
		entry = strings.ToLower(entry)
		if !rrh.FindIn(entry, availables) {
			messages = append(messages, entry)
		}
	}
	if len(messages) == 0 {
		return nil
	} else if len(messages) == 1 {
		return fmt.Errorf("%s: unknown entry. available values: %s", messages[0], strings.Join(availables, ", "))
	}
	return fmt.Errorf("%s: unknown entries. available values: %s", strings.Join(messages, ","), strings.Join(availables, ", "))
}

type listOptions struct {
	entries []string
}

var listOpts = &listOptions{}

func convertToFlags(entries []string) (groupEntry, error) {
	var result groupEntry = 0
	for _, entry := range entries {
		switch strings.ToLower(entry) {
		case "all":
			result |= all
		case "name":
			result |= name
		case "repo":
			result |= repo
		case "desc":
			result |= desc
		case "abbrev":
			result |= abbrev
		case "count":
			result |= count
		default:
			return 0, fmt.Errorf("%s: invalid entry, availables are: all, name, desc, repo, abbrev, count", entry)
		}
	}
	return result, nil
}

type groupEntry int

const (
	name   groupEntry = 1
	count             = 2
	desc              = 4
	repo              = 8
	abbrev            = 16
	all               = name | desc | repo | abbrev
)

func listGroups(c *cobra.Command, args []string, db *rrh.Database) error {
	if len(listOpts.entries) == 0 {
		listOpts.entries = []string{"name", "count"}
	}
	entry, err := convertToFlags(listOpts.entries)
	if err != nil {
		return err
	}
	return printGroupList(c, entry, db)
}

func printGroupList(c *cobra.Command, printTarget groupEntry, db *rrh.Database) error {
	for _, group := range db.Groups {
		resultItems := []string{}
		if printTarget&name == name {
			resultItems = append(resultItems, group.Name)
		}
		if printTarget&count == count {
			resultItems = append(resultItems, repositoryCount(db, group.Name))
		}
		if printTarget&desc == desc {
			resultItems = append(resultItems, group.Description)
		}
		if printTarget&abbrev == abbrev {
			resultItems = append(resultItems, strconv.FormatBool(group.OmitList))
		}
		if printTarget&repo == repo {
			resultItems = append(resultItems, strings.Join(db.FindRelationsOfGroup(group.Name), ","))
		}
		c.Println(strings.Join(resultItems, ","))
	}
	return nil
}

type updateOptions struct {
	name       string
	desc       string
	abbrev     string
	dryRunFlag bool
}

var updateOpts = &updateOptions{}

func createGroupUpdateCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "update <GROUP>",
		Short: "update the information of the specified group",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return common.PerformRrhCommand(c, args, updateGroup)
		},
	}
	flags := command.Flags()
	flags.StringVarP(&updateOpts.name, "name", "n", "", "change the group name to the given string")
	flags.StringVarP(&updateOpts.desc, "desc", "d", "", "change the description to the given string")
	flags.StringVarP(&updateOpts.abbrev, "abbrev", "a", "false", "change abbrev flag of the group, the given string must be true or false")
	flags.BoolVarP(&updateOpts.dryRunFlag, "dry-run", "D", false, "dry-run mode")
	return command
}

func updateGroup(c *cobra.Command, args []string, db *rrh.Database) error {
	abbrevFlag, err := strconv.ParseBool(updateOpts.abbrev)
	if err != nil {
		return err
	}
	group := &rrh.Group{Name: updateOpts.name, Description: updateOpts.desc, OmitList: abbrevFlag}
	if !db.UpdateGroup(args[0], group) {
		return fmt.Errorf("%s: update failed", args[0])
	}
	c.Printf("update %s -> %v", args[0], group)
	if updateOpts.dryRunFlag {
		c.Println("(dry run mode)")
	} else {
		c.Println()
	}
	return nil
}

func repositoryCount(db *rrh.Database, groupName string) string {
	count := db.ContainsCount(groupName)
	if count == 1 {
		return ",1 repository"
	}
	return fmt.Sprintf(",%d repositories", count)
}

type removeOptions struct {
	inquiry    bool
	force      bool
	dryRunFlag bool
}

var removeOpts = &removeOptions{}

func createGroupRemoveCommand() *cobra.Command {
	c := &cobra.Command{
		Use:  "rm <GROUP NAMEs...>",
		Args: cobra.MinimumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return common.PerformRrhCommand(c, args, executeGroupRemove)
		},
	}
	flags := c.Flags()
	flags.BoolVarP(&removeOpts.inquiry, "inquiry", "i", false, "inquiry mode")
	flags.BoolVarP(&removeOpts.force, "force", "f", false, "force remove")
	flags.BoolVarP(&removeOpts.dryRunFlag, "dry-run", "D", false, "dry-run mode")

	return c
}

func executeGroupRemove(c *cobra.Command, args []string, db *rrh.Database) error {
	// verbose, _ := c.PersistentFlags().GetBool("verbose")
	for _, groupName := range args {
		group := db.FindGroup(groupName)
		if group == nil && !removeOpts.force {
			return fmt.Errorf("%s: group not found", groupName)
		}
		if group != nil || !inquiryRemovingGroup(removeOpts.inquiry, groupName) {
			return nil
		}
		if err := removeGroupsImpl(c, db, groupName); err != nil {
			return err
		}
	}
	return nil
}

func removeGroupsImpl(c *cobra.Command, db *rrh.Database, groupName string) error {
	if removeOpts.force {
		db.ForceDeleteGroup(groupName)
		c.Printf("%s: group removed", groupName)
	} else if db.ContainsCount(groupName) == 0 {
		db.DeleteGroup(groupName)
		c.Printf("%s: group removed", groupName)
	} else {
		return fmt.Errorf("%s: cannot remove group. the group has relations", groupName)
	}
	return nil
}

func inquiryRemovingGroup(inquiryFlag bool, groupName string) bool {
	if !inquiryFlag {
		return true
	}
	return rrh.IsInputYes(fmt.Sprintf("%s: remove group? [yN]", groupName))
}
