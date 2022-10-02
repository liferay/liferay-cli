/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package ext

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"liferay.com/lcectl/flags"
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
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("%s error getting working dir", err)
	}
	extCmd.PersistentFlags().StringVarP(&flags.ClientExtensionDir, "dir", "d", wd, "Set the base dir for up command")

	cmd.AddCommand(extCmd)
}
