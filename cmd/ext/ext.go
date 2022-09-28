/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package ext

import (
	"os"

	"github.com/spf13/cobra"
)

var Verbose bool

// extCmd represents the ext command
var extCmd = &cobra.Command{
	Use:   "ext COMMAND [OPTIONS] [ARG...]",
	Short: "Operations to control client extension builds",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
}

func AddExtCmd(cmd *cobra.Command) {
	cmd.AddCommand(extCmd)
}
