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
	Use:   "runtime COMMAND [OPTIONS] [ARG...]",
	Short: "Operations to control the local runtime for Liferay Client Extensions",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
}

func AddRuntimeCmd(cmd *cobra.Command) {
	cmd.AddCommand(runtimeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runtimeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runtimeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
