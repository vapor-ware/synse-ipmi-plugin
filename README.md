[![Build Status](https://build.vio.sh/buildStatus/icon?job=vapor-ware/synse-ipmi-plugin/master)](https://build.vio.sh/blue/organizations/jenkins/vapor-ware%2Fsynse-ipmi-plugin/activity)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fvapor-ware%2Fsynse-ipmi-plugin.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fvapor-ware%2Fsynse-ipmi-plugin?ref=badge_shield)

# Synse IPMI Plugin

A general-purpose IPMI plugin for [Synse Server][synse-server].

## Getting Started

### Getting

You can get the IPMI plugin as a pre-built binary from a [release][plugin-release], or
as a Docker image.

```bash
docker pull vaporio/ipmi-plugin
```

If you wish to use a development build, fork/clone the repo and build the plugin
from source.

### Running

The IPMI plugin requires the IPMI-enabled servers it will communicate with to be configured.
As such, running the plugin without additional configuration will cause it to fail. As an
example of how to configure and get started with running the IPMI plugin, a simple example
deployment exists within the [example](example) directory. It runs Synse Server, the IPMI plugin,
and a basic IPMI simulator.

To run it,

```bash
cd example
docker-compose up -d
```

You can then use Synse's HTTP API or the [Synse CLI][synse-cli] to query Synse for plugin data.
Additionally, if you have `ipmitool`, you can use that to interface with the IPMI Simulator
used in the deployment.

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

## Plugin Configuration

Plugin and device configuration are described in detail in the [SDK Documentation][sdk-docs].

When deploying, you will need to provide your own plugin configuration (`config.yaml`)
with dynamic configuration defined. This is how the IPMI plugin knows about which BMCs
to communicate with. It will query the configured BMC(s) at runtime to determine their
capabilities and any devices they may have.

As an example:

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

### Dynamic Registration Options

Below are the fields that are expected in each of the dynamic registration items.
If no default is specified (`-`), the field is required.

| Field     | Description | Default |
| --------- | ----------- | --- |
| path      | The path to the `ipmitool` binary. Since this is run as a container, this generally never needs to be set in configuration, as `ipmitool` is on the container PATH. | `ipmitool` |
| hostname  | The address/hostname of the BMC to connect to.  | `-` |
| port      | The port of the BMC to connect to. This is usually 623. | `-` |
| username  | The user on the BMC to run commands as. | `-` |
| password  | The password for user authentication on the BMC, if set. | `-` |
| interface | The transport interface. Must be one of: `lan`, `lanplus` | `-` |

### Reading Outputs

Outputs are referenced by name. A single device may have more than one instance
of an output type. A value of `-` in the table below indicates that there is no value
set for that field. The *built-in* section describes outputs this plugin uses which are [built-in to
the SDK](https://synse.readthedocs.io/en/latest/sdk/concepts/reading_outputs/#built-ins).

**Built-in**

| Name          | Description                          | Unit  | Type     | Precision |
| ------------- | ------------------------------------ | :---: | -------- | :-------: |
| state         | Power or LED state. (on/off)         | `-`   | `state`  | `-`       |
| status        | Status of the boot target selection. | `-`   | `status` | `-`       |

### Device Handlers

Device Handlers are referenced by name.

| Name          | Description                                      | Outputs  | Read  | Write | Bulk Read | Listen |
| ------------- | ------------------------------------------------ | -------- | :---: | :---: | :-------: | :----: |
| boot_target   | A handler for chassis boot target.               | `status` | ✓     | ✓     | ✗         | ✗      |
| chassis.led   | A handler for chassis identify, commonly an LED. | `state`  | ✓     | ✓     | ✗         | ✗      |
| chassis.power | A handler for chassis power.                     | `state`  | ✓     | ✓     | ✗         | ✗      |

### Write Values

This plugin supports the following values when writing to a device via a handler.

| Handler       | Write Action  | Write Data | Description |
| ------------- | :-----------: | :--------: | ----------- |
| boot_target   | `target`      | `none`, `pxe`, `disk`, `safe`, `diag`, `cdrom`, `bios`, `rfloppy`, `rprimary`, `rcdrom`, `rdisk`, `floppy` | The boot target selection for the chassis. |
| chassis.led   | `state`       | `on`, `off` | The power state to put the identify LED into. |
| chassis.power | `state`       | `on`, `off`, `reset`, `cycle` | The power state to put the chassis into. |

## Tested BMCs

This plugin has been tested against the `vaporio/ipmi-simulator` image. It has also been
tested against the following hardware:

- HPE Cloudline CL2200 G3 Server

If you have tested this on other hardware and found it to work, let us know! Open a PR
and add to the list.

## Compatibility

Below is a table describing the compatibility of plugin versions with Synse platform versions.

|             | Synse v2 | Synse v3 |
| ----------- | -------- | -------- |
| plugin v1.x | ✓        | ✗        |
| plugin v2.x | ✗        | ✓        |

## Troubleshooting

### Debugging

The plugin can be run in debug mode for additional logging. This is done by:

- Setting the `debug` option  to `true` in the plugin configuration YAML

  ```yaml
  debug: true
  ```

- Passing the `--debug` flag when running the binary/image

  ```
  docker run vaporio/ipmi-plugin --debug
  ```

- Running the image with the `PLUGIN_DEBUG` environment variable set to `true`

  ```
  docker run -e PLUGIN_DEBUG=true vaporio/ipmi-plugin
  ```

### Developing

A [development/debug Dockerfile](Dockerfile.dev) is provided in the project repository to enable
building image which may be useful when developing or debugging a plugin. The development image
uses an ubuntu base, bringing with it all the standard command line  tools one would expect.
To build a development image:

```
make docker-dev
```

The built image will be tagged using the format `dev-{COMMIT}`, where `COMMIT` is the short commit for
the repository at the time. This image is not published as part of the CI pipeline, but those with access
to the Docker Hub repo may publish manually.

## Contributing / Reporting

If you experience a bug, would like to ask a question, or request a feature, open a
[new issue](https://github.com/vapor-ware/synse-ipmi-plugin/issues) and provide as much
context as possible. All contributions, questions, and feedback are welcomed and appreciated.

## License

The Synse IPMI Plugin is licensed under GPLv3. See [LICENSE](LICENSE) for more info.

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fvapor-ware%2Fsynse-ipmi-plugin.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fvapor-ware%2Fsynse-ipmi-plugin?ref=badge_large)

[synse-server]: https://github.com/vapor-ware/synse-server
[synse-sdk]: https://github.com/vapor-ware/synse-sdk
[synse-cli]: https://github.com/vapor-ware/synse-cli
[plugin-release]: https://github.com/vapor-ware/synse-ipmi-plugin/releases
[sdk-docs]: https://synse.readthedocs.io/en/latest/sdk/intro/
