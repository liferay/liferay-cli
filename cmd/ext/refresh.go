/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package ext

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"liferay.com/lcectl/constants"
	"liferay.com/lcectl/docker"
)

// refreshCmd represents the refresh command
var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refreshes client-extension workload resources in localdev server",
	Run: func(cmd *cobra.Command, args []string) {
		dockerClient, err := docker.GetDockerClient()

		if err != nil {
			log.Fatalf("%s getting docker client", err)
		}

		runLocaldevRefresh("localdev-server", dockerClient)
	},
}

func init() {
	extCmd.AddCommand(refreshCmd)
}

func runLocaldevRefresh(imageTag string, dockerClient *client.Client) error {
	ctx := context.Background()

	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{},
	}
	networkConfig.EndpointsConfig[viper.GetString(constants.Const.DockerNetwork)] =
		&network.EndpointSettings{}

	resp, err := dockerClient.ContainerCreate(
		ctx,
		&container.Config{
			Image:        imageTag,
			Cmd:          []string{"tilt", "trigger", "(Tiltfile)", "--host", "host.docker.internal"},
			AttachStdout: true,
			AttachStderr: true,
			Tty:          true,
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: "/var/run/docker.sock",
					Target: "/var/run/docker.sock",
				},
			},
			AutoRemove: true,
		},
		networkConfig,
		nil,
		"localdev-server-refresh")

	if err != nil {
		log.Fatalf("Failed to create container %s: %s", imageTag, err)
		return err
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

	statusCh, errCh := dockerClient.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}
	return nil
}
