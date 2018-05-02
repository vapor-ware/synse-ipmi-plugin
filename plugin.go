package main

import (
	"log"

	"github.com/vapor-ware/synse-ipmi-plugin/devices"
	"github.com/vapor-ware/synse-ipmi-plugin/enumerate"
	"github.com/vapor-ware/synse-sdk/sdk"
)

// Build time variables for setting the version info of a Plugin.
var (
	BuildDate     string
	GitCommit     string
	GitTag        string
	GoVersion     string
	VersionString string
)

// DeviceIdentifier defines the IPMI-specific way of uniquely identifying a device
// through its device configuration.
//
// Since the currently supported devices do not have a unique identifier beyond the
// BMC (they are really just interfaces for the BMC chassis), we supply our own unique
// ID in the "id" field. This is liable to change in the future.
func DeviceIdentifier(data map[string]string) string {
	return data["id"]
}

func main() {

	// Create the handlers for the IPMI plugin. The DeviceEnumerator will
	// create the device instance configurations at runtime.
	handlers, err := sdk.NewHandlers(DeviceIdentifier, enumerate.DeviceEnumerator)
	if err != nil {
		log.Fatal(err)
	}

	plugin, err := sdk.NewPlugin(handlers, nil)
	if err != nil {
		log.Fatal(err)
	}

	plugin.RegisterDeviceHandlers(
		// BMC Chassis Power, e.g. ipmitool [options] chassis power ...
		&devices.BmcPower,

		// BMC Chassis Boot Target, e.g. ipmitool [options] chassis bootdev ...
		&devices.BmcBootTarget,

		// BMC Chassis LED (identify), e.g. ipmitool [options] chassis identify ...
		&devices.BmcChassisLed,
	)

	// Set build-time version info.
	plugin.SetVersion(sdk.VersionInfo{
		BuildDate:     BuildDate,
		GitCommit:     GitCommit,
		GitTag:        GitTag,
		GoVersion:     GoVersion,
		VersionString: VersionString,
	})

	// Run the plugin.
	err = plugin.Run()
	if err != nil {
		log.Fatal(err)
	}
}
