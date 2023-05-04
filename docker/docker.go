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
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/briandowns/spinner"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/fileutils"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/progress"
	"github.com/docker/docker/pkg/streamformatter"
	"github.com/moby/buildkit/frontend/dockerfile/dockerignore"
	"github.com/spf13/viper"
	"liferay.com/liferay/cli/ansicolor"
	"liferay.com/liferay/cli/constants"
	lstrings "liferay.com/liferay/cli/strings"
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

type spinnerWriter struct {
	callback func(p []byte) (int, error)
}

func (e spinnerWriter) Write(p []byte) (int, error) {
	return e.callback(p)
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

func GetDockerSocketPath() string {
	dockerClient := GetDockerClient()

	daemonHost := dockerClient.DaemonHost()

	fmt.Printf("daemonHost: %v\n", daemonHost)

	protoAddrParts := strings.SplitN(daemonHost, "://", 2)

	if len(protoAddrParts) == 1 {
		log.Fatalf("unable to parse docker host `%s`", daemonHost)
	}

	return protoAddrParts[1]
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
	imageTag string, dockerFileDir string, verbose bool, s *spinner.Spinner) error {

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

	var out io.Writer = os.Stdout
	var fd = os.Stdout.Fd()
	isTerminal := true

	if !verbose {
		out = spinnerWriter{
			callback: func(bytes []byte) (int, error) {
				msg := ansicolor.StripCodes(
					strings.TrimSpace(
						string(TrimLogHeader(bytes))))

				if msg != "" {
					s.Suffix = fmt.Sprintf(
						" Building 'localdev' %s",
						lstrings.StripNewlines(lstrings.TruncateText(msg, 80)))
				}

				return 0, nil
			},
		}
		fd = 0
		isTerminal = true
	}

	err = jsonmessage.DisplayJSONMessagesStream(response.Body, out, fd, isTerminal, nil)
	if err != nil {
		return err
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

func PerformOSSpecificAdjustments(config *container.Config, host *container.HostConfig) {
	if runtime.GOOS == "linux" {
		config.User = user.UserUidAndGuidString()
	}
}

type DockerEndpoints struct {
	Host string
}

type DockerContext struct {
	Name      string
	Endpoints map[string]DockerEndpoints
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
