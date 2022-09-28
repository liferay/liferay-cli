/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package config

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set KEY VALUE",
	Short: "Set a config value",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			log.Fatal("lcectl config set requires two arguments")
		}
		viper.Set(args[0], args[1])
		viper.WriteConfig()
	},
}

func init() {
	configCmd.AddCommand(setCmd)
}
