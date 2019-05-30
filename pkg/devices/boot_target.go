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
}

// bmcBootTargetRead is the read handler function for bmc-boot-target devices.
func bmcBootTargetRead(device *sdk.Device) ([]*output.Reading, error) {
	target, err := protocol.GetChassisBootTarget(device.Data)
	if err != nil {
		return nil, err
	}

	bootTarget := output.Status.MakeReading(target).WithContext(map[string]string{
		"info": "boot target",
	})

	return []*output.Reading{
		bootTarget,
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

		switch strings.ToLower(cmd) {
		case "none":
			log.Info("Setting Boot Target -> NONE")
			target = ipmi.BootDeviceNone
		case "pxe":
			log.Info("Setting Boot Target -> PXE")
			target = ipmi.BootDevicePxe
		case "disk":
			log.Info("Setting Boot Target -> DISK")
			target = ipmi.BootDeviceDisk
		case "safe":
			log.Info("Setting Boot Target -> SAFE")
			target = ipmi.BootDeviceSafe
		case "diag":
			log.Info("Setting Boot Target -> DIAG")
			target = ipmi.BootDeviceDiag
		case "cdrom":
			log.Info("Setting Boot Target -> CDROM")
			target = ipmi.BootDeviceCdrom
		case "bios":
			log.Info("Setting Boot Target -> BIOS")
			target = ipmi.BootDeviceBios
		case "rfloppy":
			log.Info("Setting Boot Target -> RFLOPPY")
			target = ipmi.BootDeviceRemoteFloppy
		case "rprimary":
			log.Info("Setting Boot Target -> RPRIMARY")
			target = ipmi.BootDeviceRemotePrimary
		case "rcdrom":
			log.Info("Setting Boot Target -> RCDROM")
			target = ipmi.BootDeviceRemoteCdrom
		case "rdisk":
			log.Info("Setting Boot Target -> RDISK")
			target = ipmi.BootDeviceRemoteDisk
		case "floppy":
			log.Info("Setting Boot Target -> FLOPPY")
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
