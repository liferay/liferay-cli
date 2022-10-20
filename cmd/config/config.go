/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package config

import (
	"os"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:     "config COMMAND [OPTIONS] [ARG...]",
	Short:   "Operations for configuration of liferay cli",
	Aliases: []string{"cfg"},
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
}

func AddConfigCmd(cmd *cobra.Command) {
	cmd.AddCommand(configCmd)
}
