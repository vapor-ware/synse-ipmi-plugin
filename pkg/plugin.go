package pkg

import (
	log "github.com/Sirupsen/logrus"
	"github.com/vapor-ware/synse-ipmi-plugin/pkg/devices"
	"github.com/vapor-ware/synse-ipmi-plugin/pkg/outputs"
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/policies"
)

// MakePlugin creates a new instance of the IPMI plugin.
func MakePlugin() *sdk.Plugin {
	plugin := sdk.NewPlugin(
		sdk.CustomDeviceIdentifier(deviceIdentifier),
		sdk.CustomDynamicDeviceConfigRegistration(dynamicDeviceConfig),
	)

	policies.Add(policies.DeviceConfigDynamicRequired)
	policies.Add(policies.DeviceConfigFileOptional)

	err := plugin.RegisterOutputTypes(
		&outputs.ChassisBootTarget,
		&outputs.ChassisLedState,
		&outputs.ChassisPowerState,
	)
	if err != nil {
		log.Fatal(err)
	}

	plugin.RegisterDeviceHandlers(
		&devices.ChassisBootTarget,
		&devices.ChassisLed,
		&devices.ChassisPower,
	)

	return plugin
}
