/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package config

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// removeCmd represents the remove command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all config keys and values",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		keys := viper.AllKeys()
		sort.Sort(sort.StringSlice(keys))
		for _, key := range keys {
			fmt.Printf("%s=%s\n", key, viper.GetString(key))
		}
	},
}

func init() {
	configCmd.AddCommand(listCmd)
}
