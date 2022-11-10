# Getting Started with new Client Extensions Dev Experience aka Localdev

This guide will help you get started with Client Extension development using new Liferay cli and Localdev environment.

## Install and setup
- First [install `liferay` CLI](https://github.com/liferay/liferay-cli/blob/main/README.md#automated-installation)
- If automated install doesn't work, install it using [manual installation method.](https://github.com/liferay/liferay-cli/blob/main/README.md#manuall-installation)
- Install docker as per your OS (Docker desktop 4.13.1 on Windows/MacOS, Linux as per your distro)
- Verify correctly installed (and into your PATH) with `liferay --version`
- Be sure to remove any existing DXP instances running in docker.
- Run `docker system prune` to make plenty of space available for localdev
- If you run into errors with this process, [open an issue.](https://github.com/liferay/liferay-cli/issues)

## Generate localdev certificates

Client Extension development that will support cloud workloads need to use real domain names and TLS 1.3.  Therefore we need to one-time generate a 'wildcard certificate' and a 'rootCA' and install them into your system, so that your browser will trust localdev domains, and DXP can trust your workfloads running in their own domains.

- Execute `liferay runtime mkcert` (this will generate a wildcard certificate for `*.lfr.dev`)
- Execute `liferay runtime mkcert --install` (This will install the rootCA that you just generated into your OS keystores, or Chrome/Firefox/Edge keystores)
- If you have problems with these commands, please [open an issue.](https://github.com/liferay/liferay-cli/issues)

## Bring up Localdev environment

Now that you have `liferay` installed, docker is ready, and your certificates are generated, and docker ready, lets bring up an "empty" localdev environment.

- `mkdir my-client-extensions-workspace`
- `liferay extension start -b -d /full/path/to/my-client-extensions-workspace`
- The first time you run this command, it will take a long time as we have to download and prepare the client extension dev environment.  Subsequent times will be very fast.
- Once the environment is ready, your web browser should be opened to a localdev service called Tilt available at [http://localhost:10350/r/(all)/overview](http://localhost:10350/r/(all)/overview)
- On this screen you should see a "resource" on the left named `dxp.lfr.dev` This is an instance of DXP server running in a docker container on your machine.  YOu can see the logs displayed.  You can visit this DXP instance at `https://dxp.lfr.dev` . Login is `test@dxp.lfr.dev` pw: `test`
- If you have any problems with this, please open an issue on this repo and we will help you out.
- While DXP is starting, go to the next section

## Create New Client Extensions

As soon as the Tilt UI is availalbe, you are ready to start creating your first client extensions.  We have a new CLI wizard experience for this.  It is just a prototype right now, and will be rapidly improving in coming CLI releases.

- Exec `liferay extension create`
- Follow the wizard flow and create a new project from a `sample` called `coupon-with-object-actions`
- When you are fininshed you should have a new projects available in your workspace folder
- Switch to the Tilt UI and you will see the new "client extension" projects are available as "resources" on the left hand side of the screen.  If DXP is still starting, these will be "waiting" for DXP to finish starting.  Once DXP is fully up, the client extension projects will be "built" and "deployed" into localdev environment, and automatically provisioned into DXP.
- One the Client Extension resources are "all green" in the Tilt UI, you can begin using them in DXP (if you have a CSS client extension add it to the UI, if you have a object defintion with object actions you can use them.)

Lets assume that you followed the above steps and create a sample project "coupon-with-object-actions".  If so, you can now go do the following to verify things are working.

- Go to DXP UI select `Control panel > applications > coupons`
- Create a new Coupon, and hit Save
- In this sample our two "object actions" are registered to only fire when a coupon is updated.
- So first go to the Tilt UI and select the `coupon-service-springboot` and clear the logs (link on the top right)
- Then clear the logs on the `coupon-service-nodejs`
- Now go to the DXP UI and edit an existing Coupon.
- Switch back to Tilt and click both the `coupon` functions, and see their log outputs.

Lets explain a bit what just happened...(TBD...)

- Open either the `coupon-service-nodejs/src/app.js` or the `coupon-service-springboot/src/main/java/com/company/service/CouponObjectAction.java` and make a change to one of both of those files and save
- Switch to Tilt UI and see that the resources are being rebuilt and redeploy.  Once they are ready you can update the coupon again and see the changes in the log.

To work more with Object Definitions and Actions view these guides:

* Create an Object Definition (see [Creating and Managing Objects](https://learn.liferay.com/dxp/latest/en/building-applications/objects/creating-and-managing-objects.html))
* Add an Action on the Object definition (see [Defining Object Actions](https://learn.liferay.com/dxp/latest/en/building-applications/objects/creating-and-managing-objects/defining-object-actions.html))

Let remember that you can create other types of Client Extensions in this same workspace

- `liferay extension create`
- Select `from template` by `Name` and select `global-css`
- Once your global CSS extension is deployed, you can go to the DXP site and configure that client extension global css extension in the look and feel prefs for your site.
- Make any change to your global css file in your extension and those changes will be redeployed to DXP

## Stopping client extension localdev environment

- If you want to just stop your current workspace exec `liferay extension stop`
- Once you are ready to "devleop again" you can exec `liferay extension start`
- The localdev environment is not fully down exec `docker ps` and you will see a few containers left
- This is so that it is "fast" the next time you want to "start devleping client extensions"
- If you want to "fully cleanup" exec `liferay runtime delete` and then do a `docker system prune` as well to reclaim all of your disk space

## Feedback

We need your feedback, please connect to the `#liferay-cli-alpha-testing` channel on slack and start a thread and we will try to help you out!