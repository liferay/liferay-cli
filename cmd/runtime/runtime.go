/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package runtime

import (
	"os"

	"github.com/spf13/cobra"
	"liferay.com/lcectl/docker"
	"liferay.com/lcectl/flags"
	"liferay.com/lcectl/git"
	"liferay.com/lcectl/mkcert"
)

// runtimeCmd represents the runtime command
var runtimeCmd = &cobra.Command{
	Use:     "runtime COMMAND [OPTIONS] [ARG...]",
	Short:   "Operations to control the runtime environment",
	Aliases: []string{"rt"},
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		git.SyncGit(flags.Verbose)
		if cmd.Name() != "mkcert" {
			mkcert.CopyCerts(flags.Verbose)
			docker.BuildImages(flags.Verbose)
		}
	},
}

func AddRuntimeCmd(cmd *cobra.Command) {
	cmd.AddCommand(runtimeCmd)
}
