# `liferay` - Liferay Client Extension Control CLI

Tool for performing Liferay Client Extension related operations from the command line.

## Manuall Installation

### Manuall Installation On MacOS using `curl`

1. Download the binary using curl
    1. Apple Silicon
        ```bash
        curl -fsSL https://github.com/liferay/liferay-cli/releases/latest/download/liferay-darwin-arm64 -O
        ```
    1. Intel
        ```bash
        curl -fsSL https://github.com/liferay/liferay-cli/releases/latest/download/liferay-darwin-amd64 -O
        ```
1. Validate the binary (optional)
    Download the checksum file
    ```bash
    curl -fsSL https://github.com/liferay/liferay-cli/releases/latest/download/checksums.txt -o checksum.txt
    ```
    Validate the binary against the checksum file
    1. Apple Silicon
        ```bash
        shasum -c <(grep liferay-darwin-arm64 checksum.txt)
        ```
    1. Intel
        ```bash
        shasum -c <(grep liferay-darwin-amd64 checksum.txt)
        ```
    If valid, the output is:
    ```bash
    <binary>: OK
    ```
    If the check fails, `shasum` exits with nonzero status and prints output similar to:
    ```bash
    <binary>: FAILED
    shasum: WARNING: 1 computed checksum did NOT match
    ```
1. Make the binary executable.
    ```bash
    chmod +x ./liferay-*
    ```
1. Move the binary to a file location on your system `PATH` and rename it to `liferay` for convenience.
    ```bash
    sudo mv ./liferay-* /usr/local/bin/liferay
    sudo chown root: /usr/local/bin/liferay
    ```
    _Make sure `/usr/local/bin` is in your `PATH` environment variable._
1. Test to ensure the version you installed is up-to-date:
    ```bash
    liferay --version
    ```

### Manuall Installation On Linux using `curl`

1. Download the binary using curl
    ```bash
    curl -fsSL https://github.com/liferay/liferay-cli/releases/latest/download/liferay-linux-amd64 -O
    ```
1. Validate the binary (optional)
    Download the checksum file
    ```bash
    curl -fsSL https://github.com/liferay/liferay-cli/releases/latest/download/checksums.txt -o checksum.txt
    ```
    Validate the binary against the checksum file
    ```bash
    shasum -c <(grep liferay-linux-amd64 checksum.txt)
    ```
    If valid, the output is:
    ```bash
    <binary>: OK
    ```
    If the check fails, `shasum` exits with nonzero status and prints output similar to:
    ```bash
    <binary>: FAILED
    shasum: WARNING: 1 computed checksum did NOT match
    ```
1. Make the binary executable.
    ```bash
    chmod +x ./liferay-linux-amd64
    ```
1. Move the binary to a file location on your system `PATH` and rename it to `liferay` for convenience.
    ```bash
    sudo mv ./liferay-linux-amd64 /usr/local/bin/liferay
    sudo chown root: /usr/local/bin/liferay
    ```
    _Make sure `/usr/local/bin` is in your `PATH` environment variable._
1. Test to ensure the version you installed is up-to-date:
    ```bash
    liferay --version
    ```

### Manuall Installation On Windows using `curl`

1. Download the binary using curl
    1. ARM
        ```bash
        curl.exe -fsSL "https://github.com/liferay/liferay-cli/releases/latest/download/liferay-windows-arm64.exe" -O
        ```
    1. Intel
        ```bash
        curl.exe -fsSL "https://github.com/liferay/liferay-cli/releases/latest/download/liferay-windows-amd64.exe" -O
        ```
1. Validate the binary (optional)
    Download the checksum file
    ```bash
    curl -fsSL https://github.com/liferay/liferay-cli/releases/latest/download/checksums.txt -o checksum.txt
    ```
    Validate the binary against the checksum file
    1. Using Command Prompt to manually compare `CertUtil`'s output to the checksum file downloaded:
        1. ARM
            ```cmd
            CertUtil -hashfile liferay-windows-arm64.exe SHA256
            findstr liferay-windows-arm64.exe checksum.txt
            ```
        1. Intel
            ```cmd
            CertUtil -hashfile liferay-windows-amd64.exe SHA256
            findstr liferay-windows-amd64.exe checksum.txt
            ```
1. Rename it to `liferay` for convenience and move the binary to a location which can be added to your system `PATH`.
    1. ARM
        ```cmd
        ren "liferay-windows-arm64.exe" "liferay.exe" & move /Y "liferay.exe" "%USERPROFILE%\AppData\Local\Programs\Common"
        ```
    1. Intel
        ```cmd
        ren "liferay-windows-amd64.exe" "liferay.exe" & move /Y "liferay.exe" "%USERPROFILE%\AppData\Local\Programs\Common"
        ```
1. Add `%USERPROFILE%\AppData\Local\Programs\Common` to your `PATH` system variable.
    1. Press the Windows **⊞** key and type `env`.
    1. In the result pane select **Edit the system environment variables** to open the **System Properites** widget.
    1. Click **|Environment Variables...|** button.
    1. Under **User variables for `%user%`** click the `Path` entry and select **|Edit|**.
    1. Click **|New|** and paste `%USERPROFILE%\AppData\Local\Programs\Common`
    1. Click **|OK|** and close all the windows.
    1. Logout and back in.
1. Test to ensure the version you installed is up-to-date open a terminal and execute:
    ```bash
    liferay --version
    ```

## Automated Installation

### Installation on Mac or Linux using the `install.sh` script

1. Execute:
    ```bash
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/liferay/liferay-cli/HEAD/install.sh)"
    ```
1. Test to ensure the version you installed is up-to-date open a terminal and execute:
    ```bash
    liferay --version
    ```

## Run Prerequisits

* Docker (Desktop)
* the `liferay` platform specific binary

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
  del /q /s %USERPROFILE%\.liferay/cli.yaml
  rd /q /s %USERPROFILE%\.liferay/cli
  ```