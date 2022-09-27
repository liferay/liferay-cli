/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates the runtime environment for Liferay Client Extension development",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		dockerClient := InitDocker()

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Color("green")
		s.Suffix = " Synchronizing localdev sources..."
		s.FinalMSG = fmt.Sprintf("\u2705 Synced localdev sources.\n")
		s.Start()

		SyncGit()

		s.Stop()
		s.Suffix = " Building localdev images..."
		s.FinalMSG = fmt.Sprintf("\u2705 Built localdev images.\n")
		s.Restart()

		var wg sync.WaitGroup
		wg.Add(1)
		go buildImage("dxp-server", path.Join(
			viper.GetString(Const.repoDir), "docker", "images", "dxp-server"),
			dockerClient, &wg)

		wg.Add(1)
		go buildImage("localdev-server", path.Join(
			viper.GetString(Const.repoDir), "docker", "images", "localdev-server"),
			dockerClient, &wg)

		wg.Wait()

		s.Stop()
		s.Suffix = " Creating localdev environment..."
		s.FinalMSG = fmt.Sprintf("\u2705 Created localdev environment.\n")
		s.Restart()

		wg.Add(1)
		runLocaldevClusterStart("localdev-server", dockerClient, &wg)

		wg.Wait()
		s.Stop()
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func buildImage(
	imageTag string, dockerFileDir string,
	dockerClient *client.Client, wg *sync.WaitGroup) {

	/*log.Println("Building ", imageTag)*/

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
	ioutil.ReadAll(response.Body)
	wg.Done()
}

func getDockerfileBytes(path string) ([]byte, error) {
	dockerfile, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer dockerfile.Close()

	// get the file size
	stat, err := dockerfile.Stat()
	if err != nil {
		return nil, err
	}

	// read the file
	bs := make([]byte, stat.Size())
	_, err = dockerfile.Read(bs)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func runLocaldevClusterStart(imageTag string, dockerClient *client.Client, wg *sync.WaitGroup) {
	ctx := context.Background()

	// out, err := dockerClient.ImagePull(ctx, imageTag, types.ImagePullOptions{})
	// if err != nil {
	// 	log.Printf("Failed to pull image %s: %s\n", imageTag, err)
	// } else {
	// 	defer out.Close()
	// 	io.Copy(os.Stdout, out)
	// }

	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{},
	}
	networkConfig.EndpointsConfig[viper.GetString(Const.dockerNetwork)] =
		&network.EndpointSettings{}

	resp, err := dockerClient.ContainerCreate(
		ctx,
		&container.Config{
			Image: imageTag,
			Cmd:   []string{"/repo/scripts/cluster-start.sh"},
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
					Source: viper.GetString(Const.repoDir),
					Target: "/repo",
				},
			},
			AutoRemove: true,
		},
		networkConfig,
		nil,
		"localdev-start")

	if err != nil {
		log.Fatalf("Failed to create container %s: %s", imageTag, err)
	}

	/*fmt.Println(resp.Warnings, resp.ID)*/

	err = dockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})

	if err != nil {
		log.Fatalf("Failed to start container %s: %s", imageTag, err)
	}

	statusCh, errCh := dockerClient.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			log.Fatalf("Failed to wait for container %s: %s", imageTag, err)
		}
	case <-statusCh:
	}

	/*out, err := dockerClient.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})*/

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
