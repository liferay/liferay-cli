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

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates the runtime environment for Liferay Client Extension development",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		config := container.Config{
			Image: viper.GetString(constants.Const.DockerLocaldevServerImage),
			Cmd:   []string{"/repo/scripts/runtime/create.sh"},
			Env: []string{
				"CLIENT_EXTENSION_DIR_KEY=" + ext.GetExtensionDirKey(),
				"LOCALDEV_REPO=/repo",
				"LFRDEV_DOMAIN=" + viper.GetString(constants.Const.TlsLfrdevDomain),
				"DOCKER_HOST=unix:///home/localdev/.local/cx/docker.sock",
			},
		}
		host := container.HostConfig{
			Binds: []string{
				fmt.Sprintf("%s:%s", viper.GetString(constants.Const.RepoDir), "/repo"),
				docker.GetDockerSocketPath() + ":/home/localdev/.local/cx/docker.sock",
			},
			NetworkMode: container.NetworkMode(viper.GetString(constants.Const.DockerNetwork)),
		}
		docker.PerformOSSpecificAdjustments(&config, &host)

		exitCode := spinner.Spin(
			spinner.SpinOptions{
				Doing: "Creating", Done: "created", On: "'localdev' runtime environment", Enable: !flags.Verbose,
			},
			func(fior func(io.ReadCloser, bool, string) int) int {
				return docker.InvokeCommandInLocaldev("localdev-runtime-create", config, host, true, flags.Verbose, fior, "")
			})

		os.Exit(exitCode)
	},
}

func init() {
	runtimeCmd.AddCommand(createCmd)
}
