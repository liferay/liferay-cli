/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package sync

import (
	"github.com/spf13/cobra"
	"liferay.com/liferay/cli/flags"
	"liferay.com/liferay/cli/git"
)

// syncCmd represents the config command
var syncCmd = &cobra.Command{
	Use:   "sync [OPTIONS] [ARG...]",
	Short: "Sync the backing git repo of liferay cli",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		git.SyncGit(flags.Verbose)
	},
}

func AddSyncCmd(cmd *cobra.Command) {
	cmd.AddCommand(syncCmd)
}
