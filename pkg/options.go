package pkg

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/mapstructure"
	"github.com/vapor-ware/goipmi"
	"github.com/vapor-ware/synse-sdk/sdk"
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
func dynamicDeviceConfig(data map[string]interface{}) ([]*sdk.DeviceConfig, error) {

	// FIXME (etd): creating the connection and the client here are not totally
	// necessary right now other than to validate that the Connection info is
	// generally correct. We will need this later, e.g. when scanning the SDR
	// for devices.
	conn := &ipmi.Connection{}
	err := mapstructure.Decode(data, conn)
	if err != nil {
		return nil, err
	}

	log.Debugf("Connection: %+v", conn)

	// FIXME (etd): see the FIXME above - we do not currently use this.
	_, err = ipmi.NewClient(conn)
	if err != nil {
		return nil, err
	}

	// Make new power device for the BMC. This device would be akin to
	// `ipmitool [options] chassis power ...` commands.
	cfg := sdk.DeviceConfig{
		SchemeVersion: sdk.SchemeVersion{Version: "1.0"},
		Locations: []*sdk.LocationConfig{
			{
				Name:  "ipmi",
				Rack:  &sdk.LocationData{Name: "ipmi"},
				Board: &sdk.LocationData{Name: conn.Hostname},
			},
		},
		Devices: []*sdk.DeviceKind{
			{
				Name: "chassis.power",
				Outputs: []*sdk.DeviceOutput{
					{Type: "chassis.power.state"},
				},
				Instances: []*sdk.DeviceInstance{
					{
						Info:     "BMC chassis power",
						Location: "ipmi",
						Data: map[string]interface{}{
							// FIXME (etd): is there anything unique that we can use instead of hardcoding?
							// if not, find a better way than manually specifying ids...
							"path":      conn.Path,
							"hostname":  conn.Hostname,
							"port":      conn.Port,
							"username":  conn.Username,
							"password":  conn.Password,
							"interface": conn.Interface,
						},
					},
				},
			},
			{
				Name: "boot_target",
				Outputs: []*sdk.DeviceOutput{
					{Type: "chassis.boot.target"},
				},
				Instances: []*sdk.DeviceInstance{
					{
						Info:     "BMC chassis boot target",
						Location: "ipmi",
						Data: map[string]interface{}{
							"id":        "2", // FIXME (etd): see above
							"path":      conn.Path,
							"hostname":  conn.Hostname,
							"port":      conn.Port,
							"username":  conn.Username,
							"password":  conn.Password,
							"interface": conn.Interface,
						},
					},
				},
			},
			{
				Name: "chassis.led",
				Outputs: []*sdk.DeviceOutput{
					{Type: "chassis.led.state"},
				},
				Instances: []*sdk.DeviceInstance{
					{
						Info:     "BMC chassis identify LED",
						Location: "ipmi",
						Data: map[string]interface{}{
							"id":        "3", // FIXME (etd): see above
							"path":      conn.Path,
							"hostname":  conn.Hostname,
							"port":      conn.Port,
							"username":  conn.Username,
							"password":  conn.Password,
							"interface": conn.Interface,
						},
					},
				},
			},
		},
	}

	return []*sdk.DeviceConfig{&cfg}, nil
}
