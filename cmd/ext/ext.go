/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package ext

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"liferay.com/lcectl/ansicolor"
	"liferay.com/lcectl/constants"
	"liferay.com/lcectl/flags"
	lio "liferay.com/lcectl/io"
)

// extCmd represents the ext command
var extCmd = &cobra.Command{
	Use:     "extension COMMAND [OPTIONS] [ARG...]",
	Short:   "Operations to control extension environment",
	Aliases: []string{"ext"},
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var specifiedDir string

		dirFlag := cmd.Flags().Lookup("dir")

		if dirFlag != nil {
			specifiedDir = dirFlag.Value.String()
		} else if !viper.GetBool(constants.Const.ExtClientExtensionDirSpecified) {
			specifiedDir = confirmUseOfDefaultDir()
		}

		if specifiedDir != "" {
			setClientExtensionDir(specifiedDir)
		}
	},
}

func init() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("%s error getting working dir", err)
	}
	viper.SetDefault(constants.Const.ExtClientExtensionDir, wd)
	viper.SetDefault(constants.Const.ExtClientExtensionDirSpecified, false)
}

func AddExtCmd(cmd *cobra.Command) {
	extCmd.PersistentFlags().StringVarP(&flags.ClientExtensionDir, "dir", "d", viper.GetString(constants.Const.ExtClientExtensionDir), "Set the base dir for up command")
	viper.BindPFlag("dir", extCmd.Flags().Lookup("dir"))

	cmd.AddCommand(extCmd)
}

func confirmUseOfDefaultDir() string {
	fmt.Println(ansicolor.Bold("It looks like the default Client Extension directory was never specified. The current default is"), flags.ClientExtensionDir)

	if !lio.IsDirEmpty(flags.ClientExtensionDir) {
		fmt.Println(ansicolor.Bold("However, this directory is not empty. It would be preferrable to start with an empty directory."))
	}

	fmt.Println(ansicolor.Bold("Please confirm if this directory should be used."))

	validate := func(input string) error {
		if len(input) > 0 && input != "N" && input != "n" && input != "Y" && input != "y" {
			return errors.New("Please specify Yes or No")
		}
		return nil
	}

	useDefaultDirPrompt := promptui.Prompt{
		Label:    "Use default? [Y]es | [N]o | (Enter=No)",
		Validate: validate,
	}

	result, err := useDefaultDirPrompt.Run()

	if err != nil {
		log.Fatal(err)
	}

	if result == "" || result == "N" || result == "n" {
		fmt.Println(ansicolor.Bold("Please specify a directory (created if missing.) A relative path will be relative to your user home directory."))

		validate = func(input string) error {
			if len(input) <= 0 {
				return errors.New("Directory name must not be empty")
			}
			return nil
		}

		clientExtenionDirPrompt := promptui.Prompt{
			Label:    "Directory",
			Validate: validate,
		}

		result, err = clientExtenionDirPrompt.Run()

		if err != nil {
			log.Fatal(err)
		}

		return result
	}

	return ""
}

func setClientExtensionDir(dir string) {
	if !filepath.IsAbs(dir) {
		dirname, err := os.UserHomeDir()

		if err != nil {
			log.Fatal(err)
		}

		dir = filepath.Join(dirname, dir)
	}

	if !lio.Exists(dir) {
		err := os.MkdirAll(dir, 0644)

		if err != nil {
			log.Fatal(err)
		}
	}

	currentDir := viper.GetString(constants.Const.ExtClientExtensionDir)

	if dir != currentDir {
		fmt.Println("Specified Client Extension directory is ", dir)

		viper.Set(constants.Const.ExtClientExtensionDir, dir)

		flags.ClientExtensionDir = dir

		viper.Set(constants.Const.ExtClientExtensionDirSpecified, true)
		viper.WriteConfig()
	}
}
