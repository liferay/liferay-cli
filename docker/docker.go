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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/builder/dockerignore"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/fileutils"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/progress"
	"github.com/docker/docker/pkg/streamformatter"
	"github.com/spf13/viper"
	"liferay.com/liferay/cli/constants"
	"liferay.com/liferay/cli/user"
)

var (
	STDIN  = [4]byte{0, 0, 0, 0}
	STDOUT = [4]byte{1, 0, 0, 0}
	STDERR = [4]byte{2, 0, 0, 0}
)

func TrimLogHeader(bytes []byte) []byte {
	if len(bytes) <= 8 {
		return bytes
	}

	header := *(*[4]byte)(bytes[0:4])

	if header == STDIN ||
		header == STDOUT ||
		header == STDERR {

		bytes = bytes[8:]
	}

	return bytes
}

// lastProgressOutput is the same as progress.Output except
// that it only output with the last update. It is used in
// non terminal scenarios to suppress verbose messages
type lastProgressOutput struct {
	output progress.Output
}

// WriteProgress formats progress information from a ProgressReader.
func (out *lastProgressOutput) WriteProgress(prog progress.Progress) error {
	if !prog.LastUpdate {
		return nil
	}

	return out.output.WriteProgress(prog)
}

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

func GetDockerSocket() string {
	/*
		// TODO reenable this once bugs can be fixed
		socketLocation, err := lookupSocketLocationFromContext()

		if err == nil {
			return socketLocation
		}
	*/

	if runtime.GOOS == "windows" {
		return "//var/run/docker.sock"
	}

	return "/var/run/docker.sock"
}

func GetDockerClient() *client.Client {
	if dockerClient != nil {
		return dockerClient
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal("Could not create docker client", err)
	}

	return dockerClient
}

func BuildImage(
	imageTag string, dockerFileDir string, verbose bool) error {

	dockerClient := GetDockerClient()

	ctx := context.Background()

	excludes, err := readDockerignore(dockerFileDir)
	if err != nil {
		return err
	}

	excludes = trimBuildFilesFromExcludes(excludes, "Dockerfile", false)
	buildCtx, err := archive.TarWithOptions(dockerFileDir, &archive.TarOptions{
		ExcludePatterns: excludes,
		ChownOpts:       &idtools.Identity{UID: 0, GID: 0},
	})
	if err != nil {
		return err
	}

	if verbose {
		progressOutput := streamformatter.NewProgressOutput(os.Stdout)
		if !verbose {
			progressOutput = &lastProgressOutput{output: progressOutput}
		}

		buildCtx = progress.NewProgressReader(buildCtx, progressOutput, 0, "", "Sending build context to Docker daemon")
	}

	buildArgs := make(map[string]*string)

	if runtime.GOOS != "windows" {
		currentUser := user.CurrentUser()

		buildArgs["UID"] = &currentUser.Uid
		buildArgs["GID"] = &currentUser.Gid

		host := GetDockerClient().DaemonHost()

		if url, _ := client.ParseHostURL(host); url != nil {
			info, _ := os.Stat(url.Host)

			if stat, ok := info.Sys().(*syscall.Stat_t); ok {
				gid := strconv.FormatUint(uint64(stat.Gid), 10)
				buildArgs["DOCKER_GID"] = &gid
			}
		}
	}

	response, err := dockerClient.ImageBuild(
		ctx, buildCtx, types.ImageBuildOptions{
			Tags:        []string{imageTag},
			PullParent:  true,
			NetworkMode: viper.GetString(constants.Const.DockerNetwork),
			BuildArgs:   buildArgs,
		})

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if verbose {
		err = jsonmessage.DisplayJSONMessagesStream(response.Body, os.Stdout, os.Stdout.Fd(), true, nil)
		if err != nil {
			_, err = io.Copy(os.Stdout, response.Body)
		}
	} else {
		io.ReadAll(response.Body)
	}

	return nil
}

func InvokeCommandInLocaldev(
	containerName string, config container.Config, host container.HostConfig, autoremove bool, verbose bool, logPipe func(io.ReadCloser, bool, string) int, exitPattern string) int {

	dockerClient := GetDockerClient()

	ctx := context.Background()

	containers, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		log.Printf("%s error listing containers\n", err)
	}

	// delete any left over container
	match := "/" + containerName
	for _, container := range containers {
		for _, name := range container.Names {
			if name == match {
				if verbose {
					fmt.Printf("deleting lingering container %s (%s)\n", container.Names[0], container.ID)
				}

				dockerClient.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{Force: true})
			}
		}
	}

	resp, err := dockerClient.ContainerCreate(ctx, &config, &host, nil, nil, containerName)

	if err != nil {
		log.Fatalf("Failed to create container %s: %s", containerName, err)
	}

	statusChan := waitExitOrRemoved(ctx, dockerClient, resp.ID, false)

	err = dockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})

	if err != nil {
		log.Fatalf("Failed to start container %s: %s", containerName, err)
	}

	out, err := dockerClient.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true})
	if err != nil {
		log.Fatalf("%s getting container logs", err)
	}

	if autoremove {
		defer dockerClient.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{})
	}

	if logPipe != nil {
		pipeChan := make(chan (int))
		go func() {
			code := logPipe(out, verbose, exitPattern)

			pipeChan <- code
		}()

		pipeCode := <-pipeChan

		if pipeCode < 0 {
			return 0
		}

		return <-statusChan
	}

	return 0
}

type DockerEndpoints struct {
	Host string
}

type DockerContext struct {
	Name      string
	Endpoints map[string]DockerEndpoints
}

func lookupSocketLocationFromContext() (string, error) {
	// call docker context inspect to get info on sock
	cmd := exec.Command("docker", "context", "inspect")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}

	var contexts []DockerContext
	jsonerr := json.Unmarshal([]byte(output), &contexts)

	if jsonerr != nil {
		return "", jsonerr
	}

	if len(contexts) == 0 {
		return "", errors.New("unable to find docker contexts")
	}
	hostString := contexts[0].Endpoints["docker"].Host

	if strings.HasPrefix(hostString, "unix://") {
		return strings.TrimPrefix(hostString, "unix://"), nil
	}

	return "", errors.New("hostString did not have unix:// prefix")
}

// ReadDockerignore reads the .dockerignore file in the context directory and
// returns the list of paths to exclude
func readDockerignore(contextDir string) ([]string, error) {
	var excludes []string

	f, err := os.Open(filepath.Join(contextDir, ".dockerignore"))
	switch {
	case os.IsNotExist(err):
		return excludes, nil
	case err != nil:
		return nil, err
	}
	defer f.Close()

	return dockerignore.ReadAll(f)
}

// TrimBuildFilesFromExcludes removes the named Dockerfile and .dockerignore from
// the list of excluded files. The daemon will remove them from the final context
// but they must be in available in the context when passed to the API.
func trimBuildFilesFromExcludes(excludes []string, dockerfile string, dockerfileFromStdin bool) []string {
	if keep, _ := fileutils.Matches(".dockerignore", excludes); keep {
		excludes = append(excludes, "!.dockerignore")
	}

	// canonicalize dockerfile name to be platform-independent.
	dockerfile = filepath.ToSlash(dockerfile)
	if keep, _ := fileutils.Matches(dockerfile, excludes); keep && !dockerfileFromStdin {
		excludes = append(excludes, "!"+dockerfile)
	}
	return excludes
}

func waitExitOrRemoved(ctx context.Context, dockerClient *client.Client, containerID string, waitRemove bool) chan int {
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
			log.Printf("Error waiting for container: %v\n", err)
			statusC <- 125
		}
	}()

	return statusC
}
