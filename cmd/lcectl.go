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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"liferay.com/lcectl/cmd/config"
	"liferay.com/lcectl/cmd/ext"
	"liferay.com/lcectl/cmd/runtime"
	"liferay.com/lcectl/docker"
)

var cfgFile string

// lcectlCmd represents the base command when called without any subcommands
var lcectlCmd = &cobra.Command{
	Use:              "lcectl [OPTIONS] COMMAND [ARG...]",
	Short:            "Tool for performing Liferay Client Extension related operations",
	SilenceErrors:    true,
	TraverseChildren: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := lcectlCmd.Execute()
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

	lcectlCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.lcectl.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	lcectlCmd.Flags().BoolP("verbose", "v", false, "enable verbose output")

	// add sub-commands
	config.AddConfigCmd(lcectlCmd)
	runtime.AddRuntimeCmd(lcectlCmd)
	ext.AddExtCmd(lcectlCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".lcectl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".lcectl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()

	if err != nil {
		err = viper.SafeWriteConfig()

		if err != nil {
			log.Fatal("Could not write config", err)
		}
	}
}
