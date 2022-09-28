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
	"github.com/docker/go-connections/nat"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"liferay.com/lcectl/constants"
	"liferay.com/lcectl/docker"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Tilts up all client-extension workloads",
	Long:  "Starts up localdev server including DXP server and monitors client-extension workspace to build and deploy workloads",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		wg.Add(1)

		tiltPort, err := nat.NewPort("tcp", "10350")

		if err != nil {
			fmt.Println("Unable to create tilt port")
			return
		}

		exposedPorts := map[nat.Port]struct{}{
			tiltPort: {},
		}

		config := container.Config{
			Image:        "localdev-server",
			Cmd:          []string{"tilt", "up", "-f", "/repo/tilt/Tiltfile", "--stream"},
			Env:          []string{"DO_NOT_TRACK=1"},
			ExposedPorts: exposedPorts,
			AttachStdout: true,
			AttachStderr: true,
			Tty:          true,
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
			PortBindings: nat.PortMap{
				tiltPort: []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: "10350",
					},
				},
			},
		}

		docker.InvokeCommandInLocaldev("localdev-up", config, host, Verbose, &wg, nil)

		wg.Wait()

		browser.OpenURL("http://localhost:10350/r/(all)/overview")
	},
}

func init() {
	upCmd.Flags().BoolVarP(&Verbose, "verbose", "v", false, "enable verbose output")
	extCmd.AddCommand(upCmd)

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("%s error getting working dir", err)
	}
	upCmd.Flags().StringVarP(&dir, "dir", "d", wd, "Set the base dir for up command")
}
