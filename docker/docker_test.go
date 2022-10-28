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
	"strings"
	"testing"
)

func TestGetDockerSocket(t *testing.T) {
	expectedPrefix := "unix://"
	expectedSuffix := "docker.sock"
	socketLocation := GetDockerSocket()

	if !strings.HasPrefix(socketLocation, expectedPrefix) {
		t.Errorf("expected %q to have %q prefix", socketLocation, expectedPrefix)
	}

	if !strings.HasSuffix(socketLocation, expectedSuffix) {
		t.Errorf("expected %q to have %q suffix", socketLocation, expectedSuffix)
	}
}
