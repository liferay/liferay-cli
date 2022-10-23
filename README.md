# `liferay` - Liferay Client Extension Control CLI

Tool for performing Liferay Client Extension related operations from the command line.

## Run Prerequisits

* Docker (Desktop)
* the `liferay` platform specific binary

### Build Prerequisits
* `git` client
* `go` >= 1.18 (you don't absolutely need this but it simplifies debugging and installing the binary you build)
* `make` (GNU Make 3.8+, `xcode-select --install`)

## Build

* Clone the CLI repo: `git clone https://github.com/liferay/liferay-cli $CLI_SOURCES`
* `cd $CLI_SOURCES`
* install
  * if you install `go` outside of the build the following should work:
    *  `make install`
  * if you don't have `go` installed outside the build add the binary built for your platform (under `bin`) into your system `PATH` variable

## Onboarding steps

* create a new directory (say the path of that directory is stored in `${client_extension_dir}`)
* execute `liferay ext start -d ${client_extension_dir} -b`
* LIVE CODING IS NOW ACTIVE!
* Open the [Tilt UI](http://localhost:10350/r/(all)/overview) (http://localhost:10350/r/(all)/overview)

### A basic workflow after `ext start`
* From the [Tilt UI](http://localhost:10350/r/dxp.lfr.dev/overview) click the `dxp.lfr.dev` resource in the left menu
* Once DXP is started click the `dxp.lfr.dev` link found near the top of the page
* Login (`test@dxp.lfr.dev`/`test`)
* Create an Object Definition (see [Creating and Managing Objects](https://learn.liferay.com/dxp/latest/en/building-applications/objects/creating-and-managing-objects.html))
* Add an Action on the Object definition (see [Defining Object Actions](https://learn.liferay.com/dxp/latest/en/building-applications/objects/creating-and-managing-objects/defining-object-actions.html))
  * use groovy as a placeholder
* In your `${client_extension_dir}` create an Object **defintion** client extension project:
  * `liferay ext create --name=? --type=?`
* Export the Object definition JSON file from DXP into the Object defintion client extension project `src` directory
* In your `${client_extension_dir}` create an Object **action** client extension project:
  * `liferay ext create --name=? --type=?`
* Update the Object definition JSON in the Object **defintion** client extension project with the object action ID (e.g. `"objectActionExecutorKey": "function#<object-action-id>"`)

## How to customize the DXP Image used in localdev

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
  rm -rf ~/.lcect*
  ```
* windows:
  ```
  del /q /s %homedrive%%homepath%\.liferay/cli.yaml
  rd /q /s %homedrive%%homepath%\.liferay/cli
  ```