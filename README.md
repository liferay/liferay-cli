# Liferay Client Extension Control CLI

Tool for performing Liferay Client Extension related operations from the command line.

## Prerequisits

* Docker (Desktop)
* `make` (GNU Make 3.8+, `xcode-select --install`)
* install a git client
* `wget` or `curl`

## Install

* Clone the CLI repo: `git clone https://github.com/rotty3000/lcectl $CLI_SOURCES`
* `cd $CLI_SOURCES`
* `make install`

### CLI updates
If a CLI update is required
* return to `$CLI_SOURCES` dir
* `git pull`
* `make install`

## Onboarding steps

* Install the developer ROOT CA into your browser
  * We add a ROOT cert in order to support a self signed certificate for *.localdev.me so that as we add new client extensions we can give them valid "local" domains.
  * **Chrome:** Settings  > Privacy and security > Security > Manage Certificates > Authorities > Import > `$HOME/.lcectl/sources/localdev/k8s/tls/ca.crt`
  * **Firefox:** Settings > Privacy & Security > Security > View Certificates... > Authorities > Import > `$HOME/.lcectl/sources/localdev/k8s/tls/ca.crt`
* linux/mac: `lcectl ext start -d ${demodir} -b --demo`
* windows: ``
* LIVE CODING IS NOW ACTIVE! --> sitting [Tilt UI](http://localhost:10350/r/(all)/overview)

### Reproducing what just happened with the `--demo` flag
* Go to [DXP Resource](http://localhost:10350/r/dxp.localdev.me/overview)
* Click the `dxp.localdev.me` link
* Login (`test@dxp.localdev.me`/`test`)
* Define the Object Definition
* Define an Action on the Object definition (use groovy as a placeholder)
* Create an Object defintion client extension project: `lcectl ext create --name=? --type=?`
* Export the Object definition JSON file using the link in the Tilt UI `dxp.localdev.me` resource (e.g. `https://dxp.localdev.me/...`) and save the file into the Object defintion client extension project `src` directory
* Create the Object action client extension project: `lcectl ext create --name=? --type=?`
* Update the Object definition JSON in the Object defintion client extension project with the object action ID (e.g. `"objectActionExecutorKey": "function#<object-action-id>"`)

## Onboarding steps verification

---

*locadev DEMO starts here*

* [Tilt UI](http://localhost:10350/r/(all)/overview) should show all resources are green
  * DXP Instance
  * object definition
  * object action
* Navigate to `https://dxp.localdev.me`
* Show the Object UI (Firebase-like storage engine)
* Triggering the action by performing an operation on an object entry
* Edit logic of action (edit java file) (image will be rebuilt and re-deployed)
* Re-trigger action in DXP by performing an operation on an object entry and show update result of the action

*localdev DEMO ends here*

---

## How to customize the DXP Image used in localdev

* Run `LOCALDEV_RESOURCES_DIR=$(lcectl config get localdev.resources.dir)` to obtain the path where localdev resources are synced
* Edit `${LOCALDEV_RESOURCES_DIR}/docker/images/localdev-server/workspace/gradle.properties` file to set the the docker image or product key.
* If localdev runtime is already started
  * Run `lcectl ext refresh`
* If localdev runtime is not already started
  * Run `lcectl ext start`

## Getting productive with Tilt

* show logs
* refreshing resources
* disabling resources
* status bars
* ...


## Cleanup lcectl

* linux/mac:
  ```
  rm -rf ~/.lcect*
  ```
* windows:
  ```
  del /q /s %homedrive%%homepath%\.lcectl.yaml
  rd /q /s %homedrive%%homepath%\.lcectl
  ```