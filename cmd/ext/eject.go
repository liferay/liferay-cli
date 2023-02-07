/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package ext

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

var ejectCmd = &cobra.Command{
	Use:   "eject [extract workspace resources for standalone build]",
	Short: "Extracts workspace resources and client-extensions projects for standalone build",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmdArgs := append([]string{"/repo/scripts/ext/eject.py"}, args...)

		config := container.Config{
			Image: viper.GetString(constants.Const.DockerLocaldevServerImage),
			Cmd:   cmdArgs,
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
				fmt.Sprintf("%s:/workspace/client-extensions", flags.ClientExtensionDir),
			},
			NetworkMode: container.NetworkMode(viper.GetString(constants.Const.DockerNetwork)),
		}
		docker.PerformOSSpecificAdjustments(&config, &host)

		exitCode := spinner.Spin(
			spinner.SpinOptions{
				Doing: "Eject", Done: "is running", On: "'localdev' extension environment", Enable: !flags.Verbose,
			},
			func(fior func(io.ReadCloser, bool, string) int) int {
				return docker.InvokeCommandInLocaldev("localdev-build", config, host, true, flags.Verbose, fior, "")
			})

		os.Exit(exitCode)
	},
}

func init() {
	extCmd.AddCommand(ejectCmd)
}
