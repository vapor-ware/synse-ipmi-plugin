[![Build Status](https://build.vio.sh/buildStatus/icon?job=vapor-ware/synse-ipmi-plugin/master)](https://build.vio.sh/blue/organizations/jenkins/vapor-ware%2Fsynse-ipmi-plugin/activity)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fvapor-ware%2Fsynse-ipmi-plugin.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fvapor-ware%2Fsynse-ipmi-plugin?ref=badge_shield)

# Synse IPMI Plugin
A general-purpose IPMI plugin for [Synse Server][synse-server].

## Plugin Support
### Outputs
Outputs should be referenced by name. A single device can have more than one instance
of an output type. A value of `-` in the table below indicates that there is no value
set for that field.

| Name | Description | Unit | Precision | Scaling Factor |
| ---- | ----------- | ---- | --------- | -------------- |
| chassis.power.state | The power state (on/off) of the chassis. | - | - | - |
| chassis.led.state | The chassis' LED identify state. | - | - | - |
| chassis.boot.target | The boot target set for the chassis. | - | - | - |


### Device Handlers
Device Handlers should be referenced by name.

| Name | Description | Read | Write | Bulk Read |
| ---- | ----------- | ---- | ----- | --------- |
| boot_target | A handler for chassis boot target. | ✓ | ✓ | ✗ |
| chassis.led | A handler for chassis identify, commonly an LED. | ✓ | ✓ | ✗ |
| chassis.power | A handler for chassis power. | ✓ | ✓ | ✗ |


### Write Values
This plugin supports the following values when writing to a device via a handler.

| Handler | Write Action | Write Data |
| ------- | ------------ | ---------- |
| boot_target | `target` | `none`, `pxe`, `disk`, `safe`, `diag`, `cdrom`, `bios`, `rfloppy`, `rprimary`, `rcdrom`, `rdisk`, `floppy` |
| chassis.led | `state` | `on`, `off` |
| chassis.power | `state` | `on`, `off`, `reset`, `cycle` |

## Getting Started

### Getting the Plugin
You can get the Synse IPMI plugin either by cloning this repo, setting up the project dependencies,
and building the binary or docker image:
```bash
# Setup the project
$ make setup

# Build the binary
$ make build

# Build the docker image
$ make docker
```

You can also use a pre-built docker image from [DockerHub][plugin-dockerhub]
```bash
$ docker pull vaporio/ipmi-plugin
```

Or a pre-built binary from the latest [release][plugin-release].

### Running the Plugin
If you are using the plugin binary:
```bash
# The name of the plugin binary may differ depending on whether it is built
# locally or a pre-built binary is used.
$ ./plugin
```

If you are using the Docker image:
```bash
$ docker run vaporio/ipmi-plugin
```

In either case, the plugin should run, but you should not see any devices configured,
and you should see errors in the logs saying that various configurations were not found.
See the next section for how to configure your plugin. The [Example Deployment](#example-deployment)
section describes how to run a functional end-to-end example, right out of the box.

## Configuring the Plugin for your deployment
Plugin and device configuration are described in detail in the [SDK Configuration Documentation][sdk-config-docs].

For your own deployment, you will need to provide your own plugin config, `config.yml`.
The IPMI plugin dynamically generates the device configuration records from BMC data, so
there is no need to specify device instance configuration files here.

### plugin config
After reading through the docs linked above for the plugin config, take a look at the [example
plugin config](example/plugin.yml). This can be used as a reference. To specify your own BMCs, you
will need to list them under the `dynamicRegistration` field, e.g.

```yaml
dynamicRegistration:
  config:
    - hostname: 10.1.2.3
      port: 623
      username: ADMIN
      password: ADMIN
      interface: lanplus
    - hostname: 10.1.2.4
      port: 623
      username: ADMIN
      password: ADMIN
      interface: lanplus
```

Once you have your own plugin config, you can either mount it into the container at `/plugin/config.yml`, 
or mount it anywhere in the container, e.g. `/tmp/cfg/config.yml`, and specify that path in
the plugin config override environment variable, `PLUGIN_CONFIG=/tmp/cfg`.


## Example Deployment
The `deploy` directory contains configuration(s) for simple deployments of the Plugin,
and emulator backing, and Synse Server. These can both serve as examples for how to configure
the plugin with Synse Server, and can also be used to try out the plugin right out of the
box with no hardware dependency.

A Makefile target is provided to run the example deployment via `docker-compose`:
```console
$ make deploy
```

From there, you can either hit Synse Server's [HTTP API][synse-server-api] directly, or
use the [Synse CLI][synse-cli]. Additionally, if you have `ipmitool`, you can use that to
interface with the IPMI Simulator used in the deployment.

```console
$ ipmitool -H 127.0.0.1 -p 623 -U ADMIN -P ADMIN -I lanplus chassis status
  System Power         : on
  Power Overload       : false
  Power Interlock      : inactive
  Main Power Fault     : false
  Power Control Fault  : false
  Power Restore Policy : always-off
  Last Power Event     : 
  Chassis Intrusion    : inactive
  Front-Panel Lockout  : inactive
  Drive Fault          : false
  Cooling/Fan Fault    : false
```

## Tested BMCs
This plugin has been tested against the `vaporio/ipmi-simulator` image. It has also been
tested against the following hardware:

- HPE Cloudline CL2200 G3 Server

If you have tested this on other hardware and found it to work, let us know! Open a PR
and add to the list.

## Caveats
As the Plugin SDK is still being actively developed and improved, there are some under-developed
areas that affect this plugin. 
- `Location` information needs to be manually specified for each dynamically generated
  device. It may not always be clear what values to use for the `rack` and `board` ids.
  This plugin uses `ipmi` as the rack id, and the BMC host as the `board`.

## Feedback
Feedback for this plugin, or any component of the Synse ecosystem, is greatly appreciated!
If you experience any issues, find the documentation unclear, have requests for features,
or just have questions about it, we'd love to know. Feel free to open an issue for any
feedback you may have.

## Contributing
We welcome contributions to the project. The project maintainers actively manage the issues
and pull requests. If you choose to contribute, we ask that you either comment on an existing
issue or open a new one.

The Synse IPMI Plugin, and all other components of the Synse ecosystem, is released under the
[GPL-3.0](LICENSE) license.


[synse-server-api]: https://vapor-ware.github.io/synse-server
[synse-server]: https://github.com/vapor-ware/synse-server
[synse-sdk]: https://github.com/vapor-ware/synse-sdk
[synse-cli]: https://github.com/vapor-ware/synse-cli
[plugin-dockerhub]: https://hub.docker.com/r/vaporio/ipmi-plugin
[plugin-release]: https://github.com/vapor-ware/synse-ipmi-plugin/releases
[sdk-config-docs]: http://synse-sdk.readthedocs.io/en/latest/user/configuration.html

## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fvapor-ware%2Fsynse-ipmi-plugin.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fvapor-ware%2Fsynse-ipmi-plugin?ref=badge_large)