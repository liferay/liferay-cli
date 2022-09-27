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
	"runtime"

	"github.com/docker/docker/client"
	"github.com/spf13/viper"
	"liferay.com/lcectl/constants"
)

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

var dockerClient *client.Client

// Check to see if the docker command is on the executable PATH
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
