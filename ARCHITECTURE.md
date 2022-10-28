## Building

* install `make`
  * Linux
    * (Debian flavours) `sudo apt install build-essential`
  * Windows ([Reference](https://www.technewstoday.com/install-and-use-make-in-windows/))
    * Open the command line and run `winget install --accept-package-agreements --accept-source-agreements gnuwin32.make`
    * Add `C:\Program Files (x86)\GnuWin32\bin` to the system path.
* run `make`
* to build and run do
  * in linux/mac: `./gow run main.go [command]`
  * in windows: `.\gow.cmd run main.go [command]`
* to install into a temporary location run `make install`

### Reference Documentation

Cobra (command framework) [the cobra documentation here](https://github.com/spf13/cobra/blob/main/user_guide.md#using-the-cobra-library).

Viper (config framework) [the viper documentation here](https://github.com/spf13/viper#readme)

### Notes

Mounting a directory in windows.

A directory in the same current path:
```cmd
mkdir client-extensions
liferay ext start --dir %cd%\client-extensions
```

A directory in the user home:
```cmd
mkdir %homedrive%%homepath%\client-extensions
liferay ext start --dir %homedrive%%homepath%\client-extensions
```

## TODOs
  * Release requirments
    * batch mode (pass direct arguments to create command) (Ray)
    * light user docs (Greg & Ray)
    * set the DXP image to 'u' release (or latest nightly) (Ray)
    * acid test using wizard (Greg)
    * finish template metadata file (Greg)

  * Later
    * tilt args (harness capability of)
    * add support for globs in the runtime: watch: client extension properties
    * move current verbose to second level of verbose (-vv)
    * move git sync and docker image build output to a first level of verbose (-v)
    * allow specialized command outputs
    * add an export command which in association with a CasC client extension project can export existing batch models
    * detect if we need to build images and skip if not

### Commands

* `ext(ension)` - shows extension help
  * `refresh [extension]` - (temporary) until we have a live refresh **DONE**
    * `-v|--verbose` **DONE**
  * `status` - shows the status of extensions
    * `-w|--watch`
* `lxc/auth` - shows lxc help
  * `login` - login to an lxc account, creates the profile if not already present
  * `logout` - logout of the lxc account
  * `status` - show the current login profile details
  * `deploy` - deploy an extension to a specific env
  * `delete` - delete a login profile

### Built in demo

* object def
* object action
* custom element UI
* object data
* ?page definition?