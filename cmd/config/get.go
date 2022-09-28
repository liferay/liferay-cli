/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get KEY",
	Short: "Get a config value",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		value := viper.GetString(args[0])

		if value != "" {
			fmt.Println(value)
		}
	},
}

func init() {
	configCmd.AddCommand(getCmd)
}
