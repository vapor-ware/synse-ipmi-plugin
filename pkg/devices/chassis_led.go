package devices

import (
	"fmt"
	"strings"

	"github.com/vapor-ware/synse-ipmi-plugin/pkg/protocol"
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/output"
)

// ChassisLed is the handler for the bmc-boot-target device.
//
// This is really chassis identify, which according to the IPMI spec:
//
//   "This command causes the chassis to physically identify itself by a mechanism
//   chosen by the system implementation; such as turning on blinking user-visible
//   lights or emitting beeps via a speaker, LCD panel, etc" -- 28.5 Chassis Identify Command
//
// This was considered LED in Synse 1.4 so we will continue to consider it
// an LED device, even though it may not be.
var ChassisLed = sdk.DeviceHandler{
	Name:  "chassis.led",
	Read:  bmcChassisLedRead,
	Write: bmcChassisLedWrite,
}

// bmcChassisLedRead is the read handler function for bmc-chassis-led devices.
func bmcChassisLedRead(device *sdk.Device) ([]*output.Reading, error) {
	state, err := protocol.GetChassisIdentify(device.Data)
	if err != nil {
		return nil, err
	}

	chassisIdentify := output.State.MakeReading(state).WithContext(map[string]string{
		"info": "chassis identify led",
	})

	return []*output.Reading{
		chassisIdentify,
	}, nil
}

// bmcChassisLedWrite is the write handler function for bmc-chassis-led devices.
func bmcChassisLedWrite(device *sdk.Device, data *sdk.WriteData) error {
	action := data.Action
	raw := data.Data

	// When writing to a BMC LED (identify) device, we always expect there to be
	// raw data specified. If there isn't, we return an error.
	if len(raw) == 0 {
		return fmt.Errorf("no values specified for 'raw', but required")
	}

	if action == "state" {
		cmd := string(raw)

		var state protocol.IdentifyState
		// TODO (etd): figure out if we want to support intervals. if so, how? could be
		// its own action (LED interval).. could be a second value in the raw list (["on", "20"]),
		// could be a regular raw value for state here ({"state": "20"})
		switch strings.ToLower(cmd) {
		case "on":
			state = protocol.IdentifyOn
		case "off":
			state = protocol.IdentifyOff
		default:
			return fmt.Errorf("unsupported command for bmc chassis led (identify) 'state' action: %s", cmd)
		}

		err := protocol.SetChassisIdentify(device.Data, state)
		if err != nil {
			return err
		}
	} else {
		// If we reach here, then the specified action is not supported.
		return fmt.Errorf("action '%s' is not supported for bmc chassis led (identify) devices", action)
	}

	return nil
}
