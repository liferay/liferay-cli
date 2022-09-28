/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package ext

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/docker/docker/api/types/container"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"liferay.com/lcectl/constants"
	"liferay.com/lcectl/docker"
)

var dir string

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Tilts down all client-extension workloads",
	Long:  `Stops localdev server and DXP after shutting down all client-extension workloads.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		wg.Add(1)

		config := container.Config{
			Image: "localdev-server",
			Cmd:   []string{"tilt", "down", "-f", "/repo/tilt/Tiltfile"},
		}
		host := container.HostConfig{
			Binds: []string{
				fmt.Sprintf("%s:%s", viper.GetString(constants.Const.RepoDir), "/repo"),
				"/var/run/docker.sock:/var/run/docker.sock",
				fmt.Sprintf("%s:/workspace/client-extensions", dir),
				"localdevGradleCache:/root/.gradle",
				"localdevLiferayCache:/root/.liferay",
			},
			NetworkMode: container.NetworkMode(viper.GetString(constants.Const.DockerNetwork)),
		}

		docker.InvokeCommandInLocaldev("localdev-down", config, host, Verbose, &wg)

		wg.Wait()
	},
}

func init() {
	downCmd.Flags().BoolVarP(&Verbose, "verbose", "v", false, "enable verbose output")
	extCmd.AddCommand(downCmd)

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting working dir")
	}
	downCmd.Flags().StringVarP(&dir, "dir", "d", wd, "Set the base dir for down command")
}
