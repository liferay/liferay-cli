/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package ext

import (
	"os"

	"github.com/spf13/cobra"
)

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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// extCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// extCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
