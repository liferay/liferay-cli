/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package ext

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"liferay.com/lcectl/cetypes"
	"liferay.com/lcectl/prereq"
)

// refreshCmd represents the refresh command
var createCmd = &cobra.Command{
	Use:   "create [OPTIONS] [FLAGS]",
	Short: "Creates new Client Extensions using a wizard-like interface",
	Run: func(cmd *cobra.Command, args []string) {
		prereq.Prereq(Verbose)

		dat, err := cetypes.ClientExtensionTypeKeys(Verbose)

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		prompt := promptui.Select{
			Label: "Select Type",
			Items: dat,
		}

		_, result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		fmt.Printf("You choose %q\n", result)
	},
}

func init() {
	createCmd.Flags().BoolVarP(&Verbose, "verbose", "v", false, "enable verbose output")
	extCmd.AddCommand(createCmd)
}
