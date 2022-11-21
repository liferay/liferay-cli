/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package runtime

import (
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"liferay.com/liferay/cli/constants"
	"liferay.com/liferay/cli/docker"
	"liferay.com/liferay/cli/ext"
	"liferay.com/liferay/cli/flags"
	"liferay.com/liferay/cli/spinner"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts a stopped runtime environment for Liferay Client Extension development",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		config := container.Config{
			Image: viper.GetString(constants.Const.DockerLocaldevServerImage),
			Cmd:   []string{"/repo/scripts/runtime/start.sh"},
			Env: []string{
				"CLIENT_EXTENSION_DIR_KEY=" + ext.GetExtensionDirKey(),
				"LOCALDEV_REPO=/repo",
				"LFRDEV_DOMAIN=" + viper.GetString(constants.Const.TlsLfrdevDomain),
			},
		}
		host := container.HostConfig{
			Binds: []string{
				fmt.Sprintf("%s:%s", viper.GetString(constants.Const.RepoDir), "/repo"),
				docker.GetDockerSocket() + ":/var/run/docker.sock",
			},
			NetworkMode: container.NetworkMode(viper.GetString(constants.Const.DockerNetwork)),
		}
		docker.PerformOSSpecificAdjustments(&config, &host)

		exitCode := spinner.Spin(
			spinner.SpinOptions{
				Doing: "Starting", Done: "started", On: "'localdev' runtime environment", Enable: !flags.Verbose,
			},
			func(fior func(io.ReadCloser, bool, string) int) int {
				return docker.InvokeCommandInLocaldev("localdev-start", config, host, true, flags.Verbose, fior, "")
			})

		os.Exit(exitCode)
	},
}

func init() {
	runtimeCmd.AddCommand(startCmd)
}
