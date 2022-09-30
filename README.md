# Liferay Client Extension Control CLI

Tool for performing Liferay Client Extension related operations from the command line.

## Install

_TODO_

## Commands

See `lcectl`

## Building

* install `make`
  * Linux
    * (Debian flavours) `sudo apt install build-essential`
  * Windows ([Reference](https://www.technewstoday.com/install-and-use-make-in-windows/))

    * Open the command line and run `winget install --accept-package-agreements --accept-source-agreements gnuwin32.make`

    * Add `C:\Program Files (x86)\GnuWin32\bin` to the system path.
* Run `make all`
* to build and run do `go run main.go [command]`

## Running locally

To directly run the project:
* in linux/mac run `./gow run main.go [command]`
* in windows run `.\gow.cmd run main.go [command]`

## Adding additional commands

`lcectl` uses [`cobra-cli`](https://github.com/spf13/cobra-cli) for generating commands. Before being able to create new commands install `cobra-cli` (it is not required for building):
* on linux/mac run `./gow install github.com/spf13/cobra-cli@latest`
* on windows run `.\gow.cmd install github.com/spf13/cobra-cli@latest`

### Root command
To add a root command run:
* linux/mac:
  ```bash
  ./cobra-cliw add <command>
  ```
* windows:
  ```bash
  .\cobra-cliw.cmd add <command>
  ```

_e.g._ to add the command
  ```bash
  lcectl init
  ```
  run:
  * linux/mac:
    ```bash
    ./cobra-cliw add init
    ```
  * windows:
    ```bash
    .\cobra-cliw.cmd add init
    ```


### Sub-command
To add a sub-command run:
* linux/mac:
  ```bash
  ./cobra-cliw add <subcommand> -p <parent>Cmd
  ```
* windows:
  ```bash
  .\cobra-cliw.cmd add <subcommand> -p <parent>Cmd
  ```

_e.g._ to add the sub-command
  ```bash
  lcectl init extension
  ```
  run:
  * linux/mac:
    ```bash
    ./cobra-cliw add extension -p initCmd
    ```
  * windows:
    ```bash
    .\cobra-cliw.cmd add extension -p initCmd
    ```

### Using cobra

See [the cobra documentation here](https://github.com/spf13/cobra/blob/main/user_guide.md#using-the-cobra-library).

### Using viper

See [the viper documentation here](https://github.com/spf13/viper#readme)

## Planned Commands

#### commands

* `config` - show config help **DONE**
  * `get KEY` - get a config value **DONE**
  * `set KEY VALUE` - set a config value **DONE**
  * `delete KEY` - delete a value **DONE**
  * `list` - show current keys and values **DONE**
* `runtime` - shows runtime help
  * `create` - creates (if not already) and starts (if not already) the runtime **DONE**
    * `-n|--no-start`
    * `-v|--verbose` **DONE**
  * `start` - starts already created runtime
  * `stop` - stops the runtime without deleting it **DONE**
    * `-v|--verbose` **DONE**
  * `delete` - deletes the runtime and all resources **DONE**
    * `-v|--verbose` **DONE**
  * `status` - shows the status of runtime resources
* `ext(ension)` - (context/directory sensitive) shows extension help
  * `create` - create a new extension using a wizard
  * `build [extension]` - build extention(s)
  * `up` - brings up the extension(s) watcher (requires a running runtime) **DONE**
    * `-b|--browser` - (default is false) opens the browser
    * `-s|--stream` - (default is false) stay connected to streamed logs
    * `-v|--verbose` **DONE**
  * `down` - brings down the extension(s) watcher (requires a running runtime) **DONE**
    * `-v|--verbose` **DONE**
  * `refresh` - (temporary) until we have a live refresh **DONE**
    * `-v|--verbose` **DONE**
  * `status` - shows the status of extensions
    * `-w|--watch`
* `lxc/auth` - shows lxc help
  * `login` - login to an lxc account, creates the profile if not already present
  * `logout` - logout of the lxc account
  * `status` - show the current login profile details
  * `deploy` - deploy an extension to a specific env
  * `delete` - delete a login profile

### Notes

Mounting a directory in windows:

* `mkdir client-extensions`
* `lcectl ext start --dir %cd%\client-extensions`