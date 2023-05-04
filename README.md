# `cx` - CLI for Liferay Client Extension Control

Command line tool for running Liferay Client Extension in high-fidelity, locally.

## Install from Script

### Installation on Mac or Linux using the `install.sh` script

Execute:
```bash
curl -fsSL https://raw.githubusercontent.com/liferay/liferay-cli/HEAD/install.sh | bash
```

### Manuall Installation

1. Download the appropriate binary for your OS from https://github.com/liferay/liferay-cli/releases
1. Copy the binary to a directory
1. Rename it `cx`
1. Make it executable
1. Add the directory to your `PATH` environment variable

### Validating the binary

1. The checksum for validating the downloaded binary can be located [here](https://github.com/liferay/liferay-cli/releases/latest/download/checksums.txt)

## Run Prerequisites

* a Docker engine (e.g. Docker Desktop, Rancher Desktop, etc.)

## Getting Started / Onboarding steps

We have a new `Getting Started` guide [available here.](https://github.com/liferay/liferay-cli/blob/main/docs/GETTING_STARTED.markdown)

## Advanced: How to customize the DXP Image used in localdev

* Run `LOCALDEV_RESOURCES_DIR=$(liferay config get localdev.resources.dir)` to obtain the path where localdev resources are synced
* Edit `${LOCALDEV_RESOURCES_DIR}/docker/images/localdev-server/workspace/gradle.properties` file to set the the docker image or product key.
* If localdev runtime is already started
  * Run `liferay ext refresh`
* If localdev runtime is not already started
  * Run `liferay ext start`

## Getting productive with Tilt

* show logs
* refreshing resources
* disabling resources
* status bars
* ...

## Cleanup liferay/cli

* linux/mac:
  ```
  rm -rf ~/.liferay/cli*
  ```
* windows:
  ```
  del /q /s %USERPROFILE%\.liferay/cli.yaml
  rd /q /s %USERPROFILE%\.liferay/cli
  ```