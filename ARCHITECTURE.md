## Stuff for Ray and Greg to do
  * DONE make dir flag creation relative to home if relative path
  * DONE runtime status
  * DONE add persist the dir flag value if passed on command
  * DONE ext status
  * DONE `liferay ext start` should create/start runtime
  * DONE make Tilt rebuild DXP on refresh
  * DONE rename (`repo`) the config properties to `localdev.resources.*`
  * DONE take care of checking dnsmasq to see if it's working
  * DONE handle the case where a registry (kind) is already running on port 5000
  * DONE handle the case where docker-desktop managed kubernetes is already running
  * DONE switch default registry port to higher port so won't conflict
  * parameterise docker build
  * parameterise tilt build
  * add support for globs in the runtime: watch: client extension properties

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

```bash
mkdir client-extensions
liferay ext start --dir %cd%\client-extensions
```

### Built in demo

* object def
* object action
* custom element UI
* object data
* ?page definition?