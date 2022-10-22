/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"liferay.com/liferay/cli/cmd/config"
	"liferay.com/liferay/cli/cmd/ext"
	"liferay.com/liferay/cli/cmd/runtime"
	"liferay.com/liferay/cli/cmd/sync"
	"liferay.com/liferay/cli/docker"
	"liferay.com/liferay/cli/flags"
)

var Version = "development"

// liferayCmd represents the base command when called without any subcommands
var liferayCmd = &cobra.Command{
	Use:              "liferay [OPTIONS] COMMAND [ARG...]",
	Short:            "Tool for performing Liferay Client Extension related operations",
	SilenceErrors:    true,
	TraverseChildren: true,
	Version:          Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := liferayCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	initConfig()

	_, err := docker.GetDockerClient()

	if err != nil {
		log.Fatalf("%s getting dockerclient", err)
	}

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	liferayCmd.PersistentFlags().StringVar(&flags.ConfigFile, "config", "", "config file (default is $HOME/.liferay/cli.yaml)")
	liferayCmd.PersistentFlags().BoolVarP(&flags.Verbose, "verbose", "v", false, "enable verbose output")

	// add sub-commands
	config.AddConfigCmd(liferayCmd)
	ext.AddExtCmd(liferayCmd)
	runtime.AddRuntimeCmd(liferayCmd)
	sync.AddSyncCmd(liferayCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := os.UserHomeDir()
	liferayPath := filepath.Join(home, ".liferay")

	if flags.ConfigFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(flags.ConfigFile)
	} else {
		cobra.CheckErr(err)

		// Search config in home directory with name ".liferay/cli.yaml".
		viper.AddConfigPath(liferayPath)
		viper.SetConfigType("yaml")
		viper.SetConfigName("cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err = viper.ReadInConfig()

	if err != nil {
		// ensure .liferay directory exists
		os.MkdirAll(liferayPath, os.ModePerm)
		err = viper.SafeWriteConfig()

		if err != nil {
			log.Fatal("Could not write config", err)
		}
	}
}
