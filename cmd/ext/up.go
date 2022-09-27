/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package ext

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
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
		dockerClient, err := docker.GetDockerClient()

		if err != nil {
			log.Fatalf("%s error dockerclient", err)
		}

		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			log.Fatalf("%s error getting dir", err)
		}
		runLocaldevUp("localdev-server", dockerClient, dir)
	},
}

func init() {
	extCmd.AddCommand(upCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("%s error getting working dir", err)
	}
	upCmd.Flags().String("dir", wd, "Set the base dir for up command")
}

func runLocaldevUp(imageTag string, dockerClient *client.Client, wd string) {
	ctx := context.Background()

	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{},
	}
	networkConfig.EndpointsConfig[viper.GetString(constants.Const.DockerNetwork)] =
		&network.EndpointSettings{}

	tiltPort, err := nat.NewPort("tcp", "10350")

	if err != nil {
		fmt.Println("Unable to create tilt port")
		return
	}

	exposedPorts := map[nat.Port]struct{}{
		tiltPort: {},
	}

	resp, err := dockerClient.ContainerCreate(
		ctx,
		&container.Config{
			Image:        imageTag,
			Cmd:          []string{"tilt", "up", "-f", "/repo/tilt/Tiltfile", "--stream"},
			Env:          []string{"DO_NOT_TRACK=1"},
			ExposedPorts: exposedPorts,
			AttachStdout: true,
			AttachStderr: true,
			Tty:          true,
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				tiltPort: []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: "10350",
					},
				},
			},
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: "/var/run/docker.sock",
					Target: "/var/run/docker.sock",
				},
				{
					Type:   mount.TypeBind,
					Source: viper.GetString(constants.Const.RepoDir),
					Target: "/repo",
				},
				{
					Type:   mount.TypeBind,
					Source: wd,
					Target: "/workspace/client-extensions",
				},
				{
					Type:   mount.TypeVolume,
					Source: "localdevGradleCache",
					Target: "/root/.gradle",
				},
				{
					Type:   mount.TypeVolume,
					Source: "localdevLiferayCache",
					Target: "/root/.liferay",
				},
			},
			AutoRemove: true,
		},
		networkConfig,
		nil,
		"localdev-server")

	if err != nil {
		log.Fatalf("Failed to create container %s: %s", imageTag, err)
	}

	err = dockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})

	if err != nil {
		log.Fatalf("Failed to start container %s: %s", imageTag, err)
	}

	hijacked, err := dockerClient.ContainerAttach(ctx, resp.ID, types.ContainerAttachOptions{
		Stderr: true,
		Stdout: true,
		Stream: true,
	})

	if err != nil {
		log.Fatalf("Failed to attach to container %s", resp.ID)
	}

	go io.Copy(os.Stdout, hijacked.Reader)
	go io.Copy(os.Stderr, hijacked.Reader)

	time.Sleep(4 * time.Second)
	browser.OpenURL("http://localhost:10350/r/(all)/overview")

	statusCh, errCh := dockerClient.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}
}
