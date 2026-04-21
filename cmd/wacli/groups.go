package main

import "github.com/spf13/cobra"

func newGroupsCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "groups",
		Short: "Group management (read-only)",
	}
	cmd.AddCommand(newGroupsListCmd(flags))
	cmd.AddCommand(newGroupsRefreshCmd(flags))
	cmd.AddCommand(newGroupsInfoCmd(flags))
	return cmd
}
