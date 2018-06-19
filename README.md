# Synse IPMI Plugin
A general-purpose IPMI plugin for [Synse Server][synse-server].

> **NOTE**: This plugin is still in development and is considered to be
> in an alpha state. While the plugin should be functional, it requires
> updates in the underlying SDK to function better. See the [Caveats](#caveats)
> section, below.

This plugin provides read support for:
- Current power status (on/off)
- Current boot device
- Current identify state (IPMI 2.0, if supported by BMC)

And write support for:
- BMC Power control *(on/off/reset/cycle)*
- BMC Boot target *(none/pxe/disk/safe/diag/cdrom/bios/floppy/rfloppy/rprimary/rcdrom/rdisk)*
- BMC Identify *(on - 15s/off)*

## Getting Started

### Getting the Plugin
You can get the Synse IPMI plugin either by cloning this repo and running one of:
```console
# Build the IPMI plugin binary locally
$ make setup build

# Build the IPMI plugin Docker image locally
$ make docker
```

You can also use a pre-built docker image from [DockerHub][plugin-dockerhub]
```console
$ docker pull vaporio/ipmi-plugin
```

Or a pre-built binary from the latest [release][plugin-release].

### Running the Plugin
If you are using the plugin binary:
```console
# The name of the plugin binary may differ depending on whether it is built
# locally or a pre-built binary is used.
$ ./plugin
```

If you are using the Docker image:
```console
$ docker run vaporio/ipmi-plugin
```

In either case, the plugin should run, but you should not see any devices configured,
and you should see errors in the logs saying that various configurations were not found.
See the next section for how to configure your plugin. The [Example Deployment](#example-deployment)
section describes how to run a functional end-to-end example, right out of the box.

## Configuring the Plugin for your deployment
Plugins have three different types of configurations - these are all described in detail
in the [SDK Configuration Documentation][sdk-config-docs].

For your own deployment, you will need to provide your own plugin config, `config.yml`.
The IPMI plugin dynamically generates the device configuration records from BMC data, so
there is no need to specify device instance configuration files here.

### plugin config
After reading through the docs linked above for the plugin config, take a look at the [example
plugin config](example/config.yml). This can be used as a reference. To specify your own BMCs, you
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
- Dynamic `DeviceConfig` generation is a relatively new feature. As it is used, it will
  become more developed and easier to use.
- `Location` information needs to be manually specified for each dynamically generated
  device. It may not always be clear what values to use for the `rack` and `board` ids.
  This plugin uses `ipmi` as the rack id, and the BMC host as the `board`.
- Each of the devices (`bmc-boot-target`, `bmc-power`, `bmc-chassis-led`) are defined fairly
  statically - they require their current type/model to stay the same. If you wish to supply
  your own configs (e.g. to correctly identify the manufacturer and other meta-info), you
  will need to keep the same type/model.

## Feedback
Feedback for this plugin, or any component of the Synse ecosystem, is greatly appreciated!
If you experience any issues, find the documentation unclear, have requests for features,
or just have questions about it, we'd love to know. Feel free to open an issue for any
feedback you may have.

## Contributing
We welcome contributions to the project. The project maintainers actively manage the issues
and pull requests. If you choose to contribute, we ask that you either comment on an existing

The Synse IPMI Plugin, and all other components of the Synse ecosystem, is released under the
[GPL-3.0](LICENSE) license.


[synse-server-api]: https://vapor-ware.github.io/synse-server
[synse-server]: https://github.com/vapor-ware/synse-server
[synse-sdk]: https://github.com/vapor-ware/synse-sdk
[synse-cli]: https://github.com/vapor-ware/synse-cli
[plugin-dockerhub]: https://hub.docker.com/r/vaporio/ipmi-plugin
[plugin-release]: https://github.com/vapor-ware/synse-ipmi-plugin/releases
[sdk-config-docs]: http://synse-sdk.readthedocs.io/en/latest/user/configuration.html