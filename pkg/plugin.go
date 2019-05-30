package pkg

import (
	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-ipmi-plugin/pkg/devices"
	"github.com/vapor-ware/synse-sdk/sdk"
)

// MakePlugin creates a new instance of the IPMI plugin.
func MakePlugin() *sdk.Plugin {
	plugin, err := sdk.NewPlugin(
		sdk.CustomDeviceIdentifier(deviceIdentifier),
		sdk.CustomDynamicDeviceConfigRegistration(dynamicDeviceConfig),
		sdk.DeviceConfigOptional(),
		sdk.DynamicConfigRequired(),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Register device handlers.
	err = plugin.RegisterDeviceHandlers(
		&devices.ChassisBootTarget,
		&devices.ChassisLed,
		&devices.ChassisPower,
	)
	if err != nil {
		log.Fatal(err)
	}

	return plugin
}
