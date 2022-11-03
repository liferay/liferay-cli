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
	"runtime"

	"github.com/docker/docker/api/types/container"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"liferay.com/liferay/cli/constants"
	"liferay.com/liferay/cli/docker"
	"liferay.com/liferay/cli/flags"
	"liferay.com/liferay/cli/spinner"
	"liferay.com/liferay/cli/user"
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
			Image: "localdev-server",
			Cmd:   []string{"/repo/scripts/runtime/kubeconfig.sh"},
			Env: []string{
				"LOCALDEV_REPO=/repo",
				"KUBECONFIG=/var/run/.kube/config",
			},
		}
		if runtime.GOOS == "linux" {
			config.User = user.UserUidAndGuidString()
		}
		host := container.HostConfig{
			Binds: []string{
				fmt.Sprintf("%s:%s", viper.GetString(constants.Const.RepoDir), "/repo"),
				docker.GetDockerSocket() + ":/var/run/docker.sock",
				filepath.Join(userHomeDir, ".kube") + ":/var/run/.kube",
			},
			NetworkMode: container.NetworkMode(viper.GetString(constants.Const.DockerNetwork)),
		}
		if runtime.GOOS == "linux" {
			host.GroupAdd = []string{"docker"}
		}

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
