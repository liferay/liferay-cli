/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package runtime

import (
	"os"

	"github.com/spf13/cobra"
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
}

func AddRuntimeCmd(cmd *cobra.Command) {
	cmd.AddCommand(runtimeCmd)
}
