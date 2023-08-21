/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package ext

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"liferay.com/liferay/cli/ansicolor"
	"liferay.com/liferay/cli/constants"
	"liferay.com/liferay/cli/docker"
	"liferay.com/liferay/cli/ext"
	"liferay.com/liferay/cli/flags"
	"liferay.com/liferay/cli/spinner"
)

var openBrowser bool
var browserUrl = "http://localhost:10350/r/(all)/overview"

// upCmd represents the up command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts all client-extension workloads",
	Long:  "Starts up localdev server including DXP server and monitors client-extension workspace to build and deploy workloads",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		dockerClient := docker.GetDockerClient()

		ctx := context.Background()

		containers, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{All: true})
		if err != nil {
			log.Printf("%s error listing containers\n", err)
		}

		containerName := "localdev-extension-runtime"

		// if we're already running, short circuit startup
		for _, container := range containers {
			for _, name := range container.Names {
				if name == "/"+containerName && container.State == "running" {
					fmt.Println(ansicolor.Good + " 'localdev' extension environment is running.")
					doBrowser()
					return
				}
			}
		}

		tiltPort, err := nat.NewPort("tcp", "10350")

		if err != nil {
			fmt.Println("Unable to create tilt port")
			return
		}

		exposedPorts := map[nat.Port]struct{}{
			tiltPort: {},
		}

		cmdArgs := append([]string{"/repo/scripts/ext/start.sh"}, args...)

		config := container.Config{
			Image: viper.GetString(constants.Const.DockerLocaldevServerImage),
			Cmd:   cmdArgs,
			Env: []string{
				"CLIENT_EXTENSION_DIR_KEY=" + ext.GetExtensionDirKey(),
				"LOCALDEV_REPO=/repo",
				"LFRDEV_DOMAIN=" + viper.GetString(constants.Const.TlsLfrdevDomain),
			},
			ExposedPorts: exposedPorts,
		}
		host := container.HostConfig{
			Binds: []string{
				fmt.Sprintf("%s:%s", viper.GetString(constants.Const.RepoDir), "/repo"),
				docker.GetDockerSocket() + ":/var/run/docker.sock",
				fmt.Sprintf("%s:/workspace/client-extensions", flags.ClientExtensionDir),
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
		docker.PerformOSSpecificAdjustments(&config, &host)

		exitCode := spinner.Spin(
			spinner.SpinOptions{
				Doing: "Starting", Done: "started", On: "'localdev' extension environment", Enable: !flags.Verbose,
			},
			func(fior func(io.ReadCloser, bool, string) int) int {
				return docker.InvokeCommandInLocaldev(containerName, config, host, false, flags.Verbose, fior, "^Tilt started .*")
			})

		doBrowser()

		os.Exit(exitCode)
	},
}

func init() {
	extCmd.AddCommand(startCmd)
	startCmd.Flags().BoolVarP(&openBrowser, "browser", "b", false, "Open the browser to the management UI")
}

func doBrowser() {
	if openBrowser {
		var d net.Dialer
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		conn, err := d.DialContext(ctx, "tcp", "localhost:10350")
		if err != nil {
			log.Fatalf("%s trying to dial localdev server", err)
		}
		defer conn.Close()

		browser.OpenURL(browserUrl)
	} else {
		fmt.Printf("The management console can be opened at\n\t\n\t\"%s\"\n\n", browserUrl)
	}
}
