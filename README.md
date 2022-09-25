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

## Adding additional commands

### Root command
To add a root command run
```bash
cobra-cli add <command>
```

_e.g._ to add the command
  ```bash
  lcectl init
  ```
  run
  ```bash
  cobra-cli add init
  ```

### Sub-command
To add a sub-command run
```bash
cobra-cli add <subcommand> -p <parent>Cmd
```

_e.g._ to add the sub-command
  ```bash
  lcectl init extension
  ```
  run
  ```bash
  cobra-cli add extension -p initCmd
  ```

### Using cobra

See [the cobra documentation here](https://github.com/spf13/cobra/blob/main/user_guide.md#using-the-cobra-library).

### Using viper

See [the viper documentation here](https://github.com/spf13/viper#readme)