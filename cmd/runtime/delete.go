/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package runtime

import (
	"fmt"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"liferay.com/lcectl/constants"
	"liferay.com/lcectl/docker"
	lcectldocker "liferay.com/lcectl/docker"
	"liferay.com/lcectl/flags"
	"liferay.com/lcectl/prereq"
	"liferay.com/lcectl/spinner"
)

// createCmd represents the create command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete the runtime environment for Liferay Client Extension development",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		prereq.Prereq(flags.Verbose)

		config := container.Config{
			Image: "localdev-server",
			Cmd:   []string{"/repo/scripts/runtime/delete.sh"},
		}
		host := container.HostConfig{
			Binds: []string{
				fmt.Sprintf("%s:%s", viper.GetString(constants.Const.RepoDir), "/repo"),
				docker.GetDockerSocket() + ":/var/run/docker.sock",
			},
			NetworkMode: container.NetworkMode(viper.GetString(constants.Const.DockerNetwork)),
		}

		spinner.Spin(
			spinner.SpinOptions{
				Doing: "Deleting", Done: "deleted", On: "'localdev' runtime environment", Enable: !flags.Verbose,
			},
			func(fior func(io.ReadCloser, bool)) int {
				return lcectldocker.InvokeCommandInLocaldev("localdev-delete", config, host, true, flags.Verbose, fior)
			})
	},
}

func init() {
	runtimeCmd.AddCommand(deleteCmd)
}
