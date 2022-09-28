/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package docker

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/spf13/viper"
	"liferay.com/lcectl/constants"
)

var dockerClient *client.Client

func init() {
	var defaultNetwork string
	switch runtime.GOOS {
	case "linux":
		defaultNetwork = "host"
	default:
		defaultNetwork = "bridge"
	}
	viper.SetDefault(constants.Const.DockerNetwork, defaultNetwork)
	viper.SetDefault(constants.Const.DockerLocaldevServerImage, "localdev-server")
}

func GetDockerClient() (*client.Client, error) {
	if dockerClient != nil {
		return dockerClient, nil
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	return dockerClient, nil
}

func BuildImage(
	imageTag string, dockerFileDir string, verbose bool, wg *sync.WaitGroup) {

	dockerClient, err := GetDockerClient()

	if err != nil {
		log.Fatalf("%s getting dockerclient", err)
	}

	ctx := context.Background()
	buff := bytes.NewBuffer(nil)
	Tar(dockerFileDir, buff)

	response, err := dockerClient.ImageBuild(
		ctx, buff, types.ImageBuildOptions{
			Tags: []string{imageTag},
		})

	if err != nil {
		log.Fatal("Error during docker build: ", err)
	}

	defer response.Body.Close()

	if verbose {
		_, err = io.Copy(os.Stdout, response.Body)
	} else {
		io.ReadAll(response.Body)
	}
	wg.Done()
}

func InvokeCommandInLocaldev(
	containerName string, command []string, verbose bool, wg *sync.WaitGroup) {

	dockerClient, err := GetDockerClient()

	if err != nil {
		log.Fatalf("%s getting dockerclient", err)
	}

	ctx := context.Background()

	// out, err := dockerClient.ImagePull(ctx, imageTag, types.ImagePullOptions{})
	// if err != nil {
	// 	log.Printf("Failed to pull image %s: %s\n", imageTag, err)
	// } else {
	// 	defer out.Close()
	// 	io.Copy(os.Stdout, out)
	// }

	resp, err := dockerClient.ContainerCreate(
		ctx,
		&container.Config{
			Image: "localdev-server",
			Cmd:   command,
		},
		&container.HostConfig{
			Binds: []string{
				fmt.Sprintf("%s:%s", viper.GetString(constants.Const.RepoDir), "/repo"),
				"/var/run/docker.sock:/var/run/docker.sock",
			},
			NetworkMode: container.NetworkMode(viper.GetString(constants.Const.DockerNetwork)),
		},
		nil,
		nil,
		containerName)

	if err != nil {
		log.Fatalf("Failed to create container %s: %s", containerName, err)
	}

	if verbose {
		fmt.Printf("Built container with id: %s\n", resp.ID)
	}

	statusChan := waitExitOrRemoved(ctx, dockerClient, resp.ID, true)
	defer func() { <-statusChan }()

	err = dockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})

	if err != nil {
		log.Fatalf("Failed to start container %s: %s", containerName, err)
	}

	if verbose {
		out, err := dockerClient.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true})
		if err != nil {
			log.Fatalf("%s getting container logs", err)
		}

		stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	}

	wg.Done()
}

func Tar(src string, writers ...io.Writer) error {
	// ensure the src actually exists before trying to tar it
	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("Unable to tar files - %v", err.Error())
	}

	mw := io.MultiWriter(writers...)

	gzw := gzip.NewWriter(mw)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	// walk path
	return filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !fi.Mode().IsRegular() {
			return nil
		}

		// create a new dir/file header
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(strings.Replace(file, src, "", -1), string(filepath.Separator))

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}

		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		f.Close()

		return nil
	})
}

func waitExitOrRemoved(ctx context.Context, dockerClient *client.Client, containerID string, waitRemove bool) <-chan int {
	if len(containerID) == 0 {
		// containerID can never be empty
		panic("Internal Error: waitExitOrRemoved needs a containerID as parameter")
	}

	condition := container.WaitConditionNextExit
	if waitRemove {
		condition = container.WaitConditionRemoved
	}

	resultC, errC := dockerClient.ContainerWait(ctx, containerID, condition)

	statusC := make(chan int)
	go func() {
		select {
		case result := <-resultC:
			if result.Error != nil {
				log.Printf("Error waiting for container: %v\n", result.Error.Message)
				statusC <- 125
			} else {
				statusC <- int(result.StatusCode)
			}
		case err := <-errC:
			log.Printf("error waiting for container: %v\n", err)
			statusC <- 125
		}
	}()

	return statusC
}
