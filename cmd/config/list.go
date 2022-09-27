/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// removeCmd represents the remove command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all config keys and values",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		for _, key := range viper.AllKeys() {
			fmt.Printf("%s=%s\n", key, viper.GetString(key))
		}
	},
}

func init() {
	configCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
