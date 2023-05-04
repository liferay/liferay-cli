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

var buildCmd = &cobra.Command{
	Use:   "build [override default build cmd]",
	Short: "Executes the client extension workspace build to generate deployable artifacts (zip files)",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmdArgs := append([]string{"/repo/scripts/ext/build.sh"}, args...)

		config := container.Config{
			Image: viper.GetString(constants.Const.DockerLocaldevServerImage),
			Cmd:   cmdArgs,
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
				fmt.Sprintf("%s:/workspace/client-extensions", flags.ClientExtensionDir),
				"localdevGradleCache:/home/localdev/.gradle",
				"localdevLiferayCache:/home/localdev/.liferay",
				"localdevNodeModulesCache:/workspace/node_modules_cache",
			},
			NetworkMode: container.NetworkMode(viper.GetString(constants.Const.DockerNetwork)),
		}
		docker.PerformOSSpecificAdjustments(&config, &host)

		exitCode := spinner.Spin(
			spinner.SpinOptions{
				Doing: "Build", Done: "is running", On: "'localdev' extension environment", Enable: !flags.Verbose,
			},
			func(fior func(io.ReadCloser, bool, string) int) int {
				return docker.InvokeCommandInLocaldev("localdev-build", config, host, true, flags.Verbose, fior, "")
			})

		os.Exit(exitCode)
	},
}

func init() {
	extCmd.AddCommand(buildCmd)
}
