/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package runtime

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"liferay.com/liferay/cli/constants"
	"liferay.com/liferay/cli/docker"
	"liferay.com/liferay/cli/ext"
	"liferay.com/liferay/cli/flags"
	"liferay.com/liferay/cli/spinner"
)

// kubeconfigCmd represents the kubeconfig command
var kubeconfigCmd = &cobra.Command{
	Use:   "kubeconfig",
	Short: "Generates the kubeconfig necessary to talk to k8s cluster context",
	Run: func(cmd *cobra.Command, args []string) {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		config := container.Config{
			Image: viper.GetString(constants.Const.DockerLocaldevServerImage),
			Cmd:   []string{"/repo/scripts/runtime/kubeconfig.sh"},
			Env: []string{
				"CLIENT_EXTENSION_DIR_KEY=" + ext.GetExtensionDirKey(),
				"LOCALDEV_REPO=/repo",
				"KUBECONFIG=/var/run/.kube/config",
			},
		}
		host := container.HostConfig{
			Binds: []string{
				fmt.Sprintf("%s:%s", viper.GetString(constants.Const.RepoDir), "/repo"),
				docker.GetDockerSocket() + ":/var/run/docker.sock",
				filepath.Join(userHomeDir, ".kube") + ":/var/run/.kube",
			},
			NetworkMode: container.NetworkMode(viper.GetString(constants.Const.DockerNetwork)),
		}
		docker.PerformOSSpecificAdjustments(&config, &host)

		exitCode := spinner.Spin(
			spinner.SpinOptions{
				Doing: "Writing", Done: "was written", On: "'kubeconfig' config file", Enable: !flags.Verbose,
			},
			func(fior func(io.ReadCloser, bool, string) int) int {
				return docker.InvokeCommandInLocaldev("localdev-kubeconfig", config, host, true, flags.Verbose, fior, "")
			})

		os.Exit(exitCode)
	},
}

func init() {
	runtimeCmd.AddCommand(kubeconfigCmd)
}
