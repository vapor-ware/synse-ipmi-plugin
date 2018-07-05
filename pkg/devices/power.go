package devices

import (
	"fmt"
	"strings"

	"github.com/vapor-ware/goipmi"
	"github.com/vapor-ware/synse-ipmi-plugin/pkg/protocol"
	"github.com/vapor-ware/synse-sdk/sdk"
)

// ChassisPower is the handler for the bmc-power device.
var ChassisPower = sdk.DeviceHandler{
	Name:  "chassis.power",
	Read:  bmcPowerRead,
	Write: bmcPowerWrite,
}

// bmcPowerRead is the read handler function for bmc-power devices.
func bmcPowerRead(device *sdk.Device) ([]*sdk.Reading, error) {
	state, err := protocol.GetChassisPowerState(device.Data)
	if err != nil {
		return nil, err
	}

	powerState, err := device.GetOutput("chassis.power.state").MakeReading(state)
	if err != nil {
		return nil, err
	}

	return []*sdk.Reading{
		powerState,
	}, nil
}

// bmcPowerWrite is the write handler function for bmc-power devices.
func bmcPowerWrite(device *sdk.Device, data *sdk.WriteData) error {
	action := data.Action
	raw := data.Data

	// When writing to a BMC Power device, we always expect there to be
	// raw data specified. If there isn't, we return an error.
	if len(raw) == 0 {
		return fmt.Errorf("no values specified for 'raw', but required")
	}

	if action == "state" {
		cmd := string(raw)

		var state ipmi.ChassisControl
		switch strings.ToLower(cmd) {
		case "on":
			state = ipmi.ControlPowerUp
		case "off":
			state = ipmi.ControlPowerDown
		case "reset":
			state = ipmi.ControlPowerHardReset
		case "cycle":
			state = ipmi.ControlPowerCycle
		default:
			return fmt.Errorf("unsupported command for bmc power 'state' action: %s", cmd)
		}

		err := protocol.SetChassisPowerState(device.Data, state)
		if err != nil {
			return err
		}

	} else {
		// If we reach here, then the specified action is not supported.
		return fmt.Errorf("action '%s' is not supported for bmc power devices", action)
	}

	return nil
}
