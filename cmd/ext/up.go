/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package ext

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"liferay.com/lcectl/constants"
	"liferay.com/lcectl/docker"
	"liferay.com/lcectl/prereq"
	"liferay.com/lcectl/spinner"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Tilts up all client-extension workloads",
	Long:  "Starts up localdev server including DXP server and monitors client-extension workspace to build and deploy workloads",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		prereq.Prereq(Verbose)

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
			Cmd:          []string{"/repo/scripts/ext/up.sh"},
			ExposedPorts: exposedPorts,
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

		spinner.Spin(
			spinner.SpinOptions{
				Doing: "Starting", Done: "Started", On: "'localdev' extension environment", Enable: Verbose,
			},
			func(fior func(io.ReadCloser, bool)) int {
				return docker.InvokeCommandInLocaldev("localdev-up", config, host, false, Verbose, nil)
			})

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
