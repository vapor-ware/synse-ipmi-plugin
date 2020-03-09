package pkg

import (
	"fmt"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-ipmi-plugin/pkg/bmcs"
	"github.com/vapor-ware/synse-ipmi-plugin/pkg/devices"
	"github.com/vapor-ware/synse-ipmi-plugin/pkg/protocol"
	"github.com/vapor-ware/synse-sdk/sdk/config"
	"github.com/vapor-ware/synse-sdk/sdk/utils"
)

// deviceIdentifier defines the IPMI-specific way of uniquely identifying a device
// through its device configuration.
//
// Since the currently supported devices do not have a unique identifier beyond the
// BMC (they are really just interfaces for the BMC chassis), we supply our own unique
// ID in the "id" field. This is liable to change in the future.
func deviceIdentifier(data map[string]interface{}) string {
	return fmt.Sprint(data["id"])
}

// dynamicDeviceConfig is the custom override option to enable the IPMI plugin to
// dynamically register device configs at runtime for the configured BMCs.
//
// Currently, this will not scan the SDR or otherwise search for devices exposed
// by the BMC. Instead, it will just create a few higher-level devices for chassis
// control for each configured BMC.
func dynamicDeviceConfig(data map[string]interface{}) ([]*config.DeviceProto, error) {

	// Create a new client for the configured BMC.
	client, err := protocol.NewClientFromConfig(data)
	if err != nil {
		log.WithFields(log.Fields{
			"data": utils.RedactPasswords(data),
		}).Error("[ipmi] failed to create a new client from dynamic config")
		return nil, err
	}

	// Generate an internal ID used to reference the BMC client.
	bmcID := uuid.New().String()
	bmcs.Add(bmcID, client)

	log.WithFields(log.Fields{
		"interface": client.Interface,
		"port":      client.Port,
		"path":      client.Path,
		"host":      client.Hostname,
		"user":      client.Username,
	}).Debug("[ipmi] created client from dynamic config")

	// FIXME (etd): The device IDs are hardcoded here incrementally. While not incorrect, it doesn't
	// 	 feel like the best solution. Investigate alternatives to this manual definition of IDs.

	cfg := []*config.DeviceProto{

		// Chassis Power Device
		{
			Type:    "power",
			Handler: devices.ChassisPower.Name,
			Context: map[string]string{
				"location": "chassis",
			},
			Instances: []*config.DeviceInstance{
				{
					Info: "BMC chassis power",
					Data: map[string]interface{}{
						"id":  "1",
						"bmc": bmcID,
					},
				},
			},
		},

		// Chassis Boot Target Device
		{
			Type:    "boot_target",
			Handler: devices.ChassisBootTarget.Name,
			Instances: []*config.DeviceInstance{
				{
					Info: "BMC boot target",
					Data: map[string]interface{}{
						"id":  "2",
						"bmc": bmcID,
					},
				},
			},
		},

		// Chassis LED (Identify) Device
		{
			Type:    "led",
			Handler: devices.ChassisLed.Name,
			Context: map[string]string{
				"location": "chassis",
			},
			Instances: []*config.DeviceInstance{
				{
					Info: "BMC chassis identify LED",
					Data: map[string]interface{}{
						"id":  "3",
						"bmc": bmcID,
					},
				},
			},
		},
	}

	return cfg, nil
}
