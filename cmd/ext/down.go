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

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Tilts down all client-extension workloads",
	Long:  `Stops localdev server and DXP after shutting down all client-extension workloads.`,
	Run: func(cmd *cobra.Command, args []string) {
		dockerClient, err := docker.GetDockerClient()
		if err != nil {
			log.Fatalf("%s error getting dockerclient", err)
		}

		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			log.Fatalf("%s error getting dir", err)
		}
		runLocaldevDown("localdev-server", dockerClient, dir)
		stopRemoveLocaldevServer("localdev-server", dockerClient)
	},
}

func init() {
	extCmd.AddCommand(downCmd)

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting working dir")
	}
	downCmd.Flags().String("dir", wd, "Set the base dir for down command")
}

func runLocaldevDown(imageTag string, dockerClient *client.Client, wd string) error {
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
			Cmd:          []string{"tilt", "down", "-f", "/repo/tilt/Tiltfile"},
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
		"localdev-server-down")

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

func stopRemoveLocaldevServer(containerName string, dockerClient *client.Client) error {
	ctx := context.Background()

	if err := dockerClient.ContainerStop(ctx, containerName, nil); err != nil {
		log.Printf("Unable to stop container %s: %s", containerName, err)
		return err
	}

	return nil
}
