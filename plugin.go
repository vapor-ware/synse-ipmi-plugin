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
// FIXME - this is just a stub for framing up the plugin
func DeviceIdentifier(data map[string]string) string {
	return data["id"]
}

func main() {

	handlers, err := sdk.NewHandlers(DeviceIdentifier, enumerate.DeviceEnumerator)
	if err != nil {
		log.Fatal(err)
	}

	plugin, err := sdk.NewPlugin(handlers, nil)
	if err != nil {
		log.Fatal(err)
	}

	plugin.RegisterDeviceHandlers(
		&devices.BmcPower,
		&devices.BmcBootTarget,
		//&devices.BmcChassisLed, // FIXME (etd) - temporarily disabled as this is unhappy with the emulator...
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
