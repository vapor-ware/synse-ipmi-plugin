package devices

import (
	"fmt"
	"strings"

	ipmi "github.com/vapor-ware/goipmi"
	"github.com/vapor-ware/synse-ipmi-plugin/pkg/bmcs"
	"github.com/vapor-ware/synse-ipmi-plugin/pkg/protocol"
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/output"
)

// ChassisPower is the handler for the bmc-power device.
var ChassisPower = sdk.DeviceHandler{
	Name:  "chassis.power",
	Read:  bmcPowerRead,
	Write: bmcPowerWrite,

	Actions: []string{
		"state",
	},
}

// bmcPowerRead is the read handler function for bmc-power devices.
func bmcPowerRead(device *sdk.Device) ([]*output.Reading, error) {
	bmcID, err := bmcs.GetIDFromConfig(device.Data)
	if err != nil {
		return nil, err
	}

	client, err := bmcs.Get(bmcID)
	if err != nil {
		return nil, err
	}

	state, err := protocol.GetChassisPowerState(client)
	if err != nil {
		return nil, err
	}

	return []*output.Reading{
		output.State.MakeReading(state).WithContext(map[string]string{
			"info": "chassis power state",
		}),
	}, nil
}

// bmcPowerWrite is the write handler function for bmc-power devices.
func bmcPowerWrite(device *sdk.Device, data *sdk.WriteData) error {
	bmcID, err := bmcs.GetIDFromConfig(device.Data)
	if err != nil {
		return err
	}

	client, err := bmcs.Get(bmcID)
	if err != nil {
		return err
	}

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

		err := protocol.SetChassisPowerState(client, state)
		if err != nil {
			return err
		}

	} else {
		// If we reach here, then the specified action is not supported.
		return fmt.Errorf("action '%s' is not supported for bmc power devices", action)
	}

	return nil
}
