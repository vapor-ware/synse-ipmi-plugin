package enumerate

import (
	"strconv"

	"github.com/mitchellh/mapstructure"
	"github.com/vapor-ware/synse-sdk/sdk/config"
	"github.com/vapor-ware/synse-sdk/sdk/logger"
	"github.com/vmware/goipmi"
)

// FIXME (etd): the SDK needs some improvements to how dynamic device
// creation happens. While this implementation works, it is liable to change
// moving forward.

// DeviceEnumerator is the enumeration handler for the IPMI Plugin. It generates
// device instance configurations at runtime for the configured BMCs.
//
// Currently, this will not scan the SDR or otherwise search for devices exposed
// by the BMC. Instead, it will just create a few higher-level devices for chassis
// control for each configured BMC.
func DeviceEnumerator(data map[string]interface{}) ([]*config.DeviceConfig, error) {
	var devices []*config.DeviceConfig

	// FIXME (etd): creating the connection and the client here are not totally
	// necessary right now other than to validate that the Connection info is
	// generally correct. We will need this later, e.g. when scanning the SDR
	// for devices.
	conn := &ipmi.Connection{}
	err := mapstructure.Decode(data, conn)
	if err != nil {
		return nil, err
	}

	logger.Debugf("Connection: %+v", conn)

	// FIXME (etd): see the FIXME above - we do not currently use this.
	_, err = ipmi.NewClient(conn)
	if err != nil {
		return nil, err
	}

	// Make new power device for the BMC. This device would be akin to
	// `ipmitool [options] chassis power ...` commands.
	power := config.DeviceConfig{
		Version: "1",
		Type:    "power",
		Model:   "bmc-power",
		// FIXME (etd): how do we determine the location for dynamic devices?
		// These values are semi-sane placeholders until we figure this out.
		Location: config.Location{
			Rack:  "ipmi",
			Board: conn.Hostname,
		},
		// We can put the connection info here.. sorta (needs to be string->interface{},
		// or interface{}->interface{}
		Data: map[string]string{
			"id":        "1", // FIXME (etd): is there anything unique that we can use instead of hardcoding? if not, find a better way than manually specifying ids...
			"path":      conn.Path,
			"hostname":  conn.Hostname,
			"port":      strconv.Itoa(conn.Port),
			"username":  conn.Username,
			"password":  conn.Password,
			"interface": conn.Interface,
		},
	}
	devices = append(devices, &power)

	// Make new boot target device for the BMC. This device would be akin to
	// `ipmitool [options] chassis bootdev ...` commands.
	bootTarget := config.DeviceConfig{
		Version: "1",
		Type:    "boot_target",
		Model:   "bmc-boot-target",
		Location: config.Location{
			Rack:  "ipmi",
			Board: conn.Hostname,
		},
		Data: map[string]string{
			"id":        "2", // FIXME (etd): see above
			"path":      conn.Path,
			"hostname":  conn.Hostname,
			"port":      strconv.Itoa(conn.Port),
			"username":  conn.Username,
			"password":  conn.Password,
			"interface": conn.Interface,
		},
	}
	devices = append(devices, &bootTarget)

	// Make new identify device for the BMC. This device would be akin to
	// `ipmitool [options] chassis identify ...` commands.
	//
	// NOTE: The prototype for this device is NOT registered with the plugin
	// (in plugin.go). As such, it should never show up as a device for the plugin.
	// See the comment there for more info.
	/*
	identifyLed := config.DeviceConfig{
		Version: "1",
		Type:    "led",
		Model:   "bmc-chassis-led",
		Location: config.Location{
			Rack:  "ipmi",
			Board: conn.Hostname,
		},
		Data: map[string]string{
			"id":        "3", // FIXME (etd): see above
			"path":      conn.Path,
			"hostname":  conn.Hostname,
			"port":      strconv.Itoa(conn.Port),
			"username":  conn.Username,
			"password":  conn.Password,
			"interface": conn.Interface,
		},
	}
	devices = append(devices, &identifyLed)
	*/

	return devices, nil
}
