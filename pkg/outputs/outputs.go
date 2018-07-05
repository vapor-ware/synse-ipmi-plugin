package outputs

import "github.com/vapor-ware/synse-sdk/sdk"

var (
	// ChassisPowerState is the output type for chassis power (on/off).
	ChassisPowerState = sdk.OutputType{
		Name: "chassis.power.state",
	}

	// ChassisLedState is the output type for chassis identify LED state (on/off).
	ChassisLedState = sdk.OutputType{
		Name: "chassis.led.state",
	}

	// ChassisBootTarget is the output type for chassis boot target settings.
	ChassisBootTarget = sdk.OutputType{
		Name: "chassis.boot.target",
	}
)
