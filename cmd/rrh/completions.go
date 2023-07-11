package main

import "github.com/spf13/cobra"

func generateCompletionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "completion <bash|zsh|fish|powershell>",
		Short:                 "generate completion script",
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		RunE:                  generateCompletion,
	}
	cmd.Hidden = true
	return cmd
}

func generateCompletion(c *cobra.Command, args []string) error {
	switch args[0] {
	case "bash":
		c.Root().GenBashCompletion(c.OutOrStdout())
	case "zsh":
		c.Root().GenZshCompletion(c.OutOrStdout())
	case "fish":
		c.Root().GenFishCompletion(c.OutOrStdout(), true)
	case "powershell":
		c.Root().GenPowerShellCompletion(c.OutOrStdout())
	}
	return nil
}
