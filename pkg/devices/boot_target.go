package devices

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	ipmi "github.com/vapor-ware/goipmi"
	"github.com/vapor-ware/synse-ipmi-plugin/pkg/protocol"
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/output"
)

// ChassisBootTarget is the handler for the bmc-boot-target device.
var ChassisBootTarget = sdk.DeviceHandler{
	Name:  "boot_target",
	Read:  bmcBootTargetRead,
	Write: bmcBootTargetWrite,

	Actions: []string{
		"target",
	},
}

// bmcBootTargetRead is the read handler function for bmc-boot-target devices.
func bmcBootTargetRead(device *sdk.Device) ([]*output.Reading, error) {
	target, err := protocol.GetChassisBootTarget(device.Data)
	if err != nil {
		return nil, err
	}

	return []*output.Reading{
		output.Status.MakeReading(target).WithContext(map[string]string{
			"info": "boot target",
		}),
	}, nil
}

// bmcBootTargetWrite is the write handler function for bmc-boot-target devices.
func bmcBootTargetWrite(device *sdk.Device, data *sdk.WriteData) error {
	action := data.Action
	raw := data.Data

	// When writing to a BMC boot target device, we always expect there to be
	// raw data specified. If there isn't, we return an error.
	if len(raw) == 0 {
		return fmt.Errorf("no values specified for 'raw', but required")
	}

	if action == "target" {
		cmd := string(raw)

		var target ipmi.BootDevice

		log.WithFields(log.Fields{
			"target": strings.ToUpper(cmd),
			"device": device.GetID(),
		}).Info("[ipmi] setting boot target")

		switch strings.ToLower(cmd) {
		case "none":
			target = ipmi.BootDeviceNone
		case "pxe":
			target = ipmi.BootDevicePxe
		case "disk":
			target = ipmi.BootDeviceDisk
		case "safe":
			target = ipmi.BootDeviceSafe
		case "diag":
			target = ipmi.BootDeviceDiag
		case "cdrom":
			target = ipmi.BootDeviceCdrom
		case "bios":
			target = ipmi.BootDeviceBios
		case "rfloppy":
			target = ipmi.BootDeviceRemoteFloppy
		case "rprimary":
			target = ipmi.BootDeviceRemotePrimary
		case "rcdrom":
			target = ipmi.BootDeviceRemoteCdrom
		case "rdisk":
			target = ipmi.BootDeviceRemoteDisk
		case "floppy":
			target = ipmi.BootDeviceFloppy
		default:
			return fmt.Errorf("unsupported command for bmc boot target 'target' action: %s", cmd)
		}

		err := protocol.SetChassisBootTarget(device.Data, target)
		if err != nil {
			return err
		}

	} else {
		// If we reach here, then the specified action is not supported.
		return fmt.Errorf("action '%s' is not supported for bmc boot target devices", action)
	}

	return nil
}
