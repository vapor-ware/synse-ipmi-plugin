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

		// Chassis LED (identify) is disabled. The underlying IPMI library
		// wraps ipmitool for the lanplus interface. The ipmitool (version
		// 1.8.16) prints out the response bytes on new lines to keep 16-byte
		// width rows, e.g.
		//
		//	$ ipmitool -H 127.0.0.1 -U ADMIN -P ADMIN -I lanplus raw 0x00 0x04 0x0f 0x00
		//  7f 00 00 90 cb 68 d5 ff 7f 00 00 b0 ca 68 d5 ff
		//	7f 00 00 90 cb 68 d5 ff 7f 00 00 20 91 88 01 ec
		//  55 00 00
		//
		// The IPMI library does not properly handle the newline in the case
		// of more than one line of raw response, causing a panic.
		//&devices.BmcChassisLed,
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
