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
package cmd

import (
	"fmt"
	"log"
	"os/exec"
)

var dockerPath string

// Check to see if the docker command is on the executable PATH
func InitDocker() {
	path, err := exec.LookPath("docker")

	if err != nil {
		log.Fatal(`Could not find docker on the system PATH!

Please install docker and make sure it is added to the system PATH.`)
	}

	dockerPath = path
	fmt.Println("docker found at", dockerPath)
}
