package devices

import (
	"fmt"
	"strings"

	"github.com/vapor-ware/synse-ipmi-plugin/protocol"
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/logger"
	"github.com/vmware/goipmi"
)

// BmcBootTarget is the handler for the bmc-boot-target device.
var BmcBootTarget = sdk.DeviceHandler{
	Type:  "boot_target",
	Model: "bmc-boot-target",

	Read:  bmcBootTargetRead,
	Write: bmcBootTargetWrite,
}

func bmcBootTargetRead(device *sdk.Device) ([]*sdk.Reading, error) {

	target, err := protocol.GetChassisBootTarget(device.Data)
	if err != nil {
		return nil, err
	}

	ret := []*sdk.Reading{
		sdk.NewReading("target", target),
	}

	return ret, nil
}

func bmcBootTargetWrite(device *sdk.Device, data *sdk.WriteData) error {

	action := data.Action
	raw := data.Raw

	// When writing to a BMC Power device, we always expect there to be
	// raw data specified. If there isn't, we return an error.
	if len(raw) == 0 {
		return fmt.Errorf("no values specified for 'raw', but required")
	}

	if action == "target" {
		cmd := string(raw[0])

		var target ipmi.BootDevice

		switch strings.ToLower(cmd) {
		case "none":
			logger.Info("Setting Boot Target -> NONE")
			target = ipmi.BootDeviceNone
		case "pxe":
			logger.Info("Setting Boot Target -> PXE")
			target = ipmi.BootDevicePxe
		case "disk":
			logger.Info("Setting Boot Target -> DISK")
			target = ipmi.BootDeviceDisk
		case "safe":
			logger.Info("Setting Boot Target -> SAFE")
			target = ipmi.BootDeviceSafe
		case "diag":
			logger.Info("Setting Boot Target -> DIAG")
			target = ipmi.BootDeviceDiag
		case "cdrom":
			logger.Info("Setting Boot Target -> CDROM")
			target = ipmi.BootDeviceCdrom
		case "bios":
			logger.Info("Setting Boot Target -> BIOS")
			target = ipmi.BootDeviceBios
		case "rfloppy":
			logger.Info("Setting Boot Target -> RFLOPPY")
			target = ipmi.BootDeviceRemoteFloppy
		case "rprimary":
			logger.Info("Setting Boot Target -> RPRIMARY")
			target = ipmi.BootDeviceRemotePrimary
		case "rcdrom":
			logger.Info("Setting Boot Target -> RCDROM")
			target = ipmi.BootDeviceRemoteCdrom
		case "rdisk":
			logger.Info("Setting Boot Target -> RDISK")
			target = ipmi.BootDeviceRemoteDisk
		case "floppy":
			logger.Info("Setting Boot Target -> FLOPPY")
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


// curl -H "Content-Type: application/json" -X POST -d '{"action": "target", "raw": "pxe"}' "http://localhost:5000/synse/2.0/write/ipmi/ipmisim/dbe630b6dd49d7b9d4f435349bc2cccc"