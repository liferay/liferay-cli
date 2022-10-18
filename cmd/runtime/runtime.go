/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package runtime

import (
	"os"

	"github.com/spf13/cobra"
	"liferay.com/lcectl/flags"
	"liferay.com/lcectl/mkcert"
	"liferay.com/lcectl/prereq"
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
		prereq.Prereq(flags.Verbose)

		if cmd.Name() != "mkcert" {
			mkcert.CopyCerts()
		}
	},
}

func AddRuntimeCmd(cmd *cobra.Command) {
	cmd.AddCommand(runtimeCmd)
}
