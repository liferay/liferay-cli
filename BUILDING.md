### Build Prerequisits
* `git` client
* `go` >= 1.18 (you don't absolutely need this but it simplifies debugging and installing the binary you build)
* `make` (GNU Make 3.8+, `xcode-select --install`)

## Build

1. Clone the CLI repo: `git clone https://github.com/liferay/liferay-cli $CLI_SOURCES`
1. `cd $CLI_SOURCES`
1. Execute `make`
1. To install from source:
    * if you installed `go` outside of the build the following should work:
        *  `make install`
    * if you don't have `go` installed outside the build follow the [Manuall Installation](README.md#manuall-installation) proceedure.

