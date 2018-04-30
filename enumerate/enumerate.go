package enumerate

import (
	"strconv"

	"github.com/mitchellh/mapstructure"
	"github.com/vapor-ware/synse-sdk/sdk/config"
	"github.com/vapor-ware/synse-sdk/sdk/logger"
	"github.com/vmware/goipmi"
)

// FIXME (etd): the SDK needs some improvements to how dynamic device
// creation happens.. this is all a somewhat roundabout workaround for the
// time being..

func DeviceEnumerator(data map[string]interface{}) ([]*config.DeviceConfig, error) {
	var devices []*config.DeviceConfig

	conn := &ipmi.Connection{}
	err := mapstructure.Decode(data, conn)
	if err != nil {
		return nil, err
	}

	logger.Infof("Connection: %+v", conn)

	client, err := ipmi.NewClient(conn)
	if err != nil {
		return nil, err
	}

	// need to do something with the client here?
	logger.Infof("Client: %v", client)

	// Make new power device for the BMC
	power := config.DeviceConfig{
		Version: "1",
		Type:    "power",
		Model:   "bmc-power",
		// FIXME (etd): how do we determine the location for dynamic devices?
		// These values are placeholders until we figure this out.
		Location: config.Location{
			Rack:  "ipmi",
			Board: conn.Hostname,
		},
		// We can put the connection info here.. sorta (needs to be string->interface{},
		// or interface{}->interface{}
		Data: map[string]string{
			"path":      conn.Path,
			"hostname":  conn.Hostname,
			"port":      strconv.Itoa(conn.Port),
			"username":  conn.Username,
			"password":  conn.Password,
			"interface": conn.Interface,
		},
	}
	devices = append(devices, &power)

	// Make new boot target device for the BMC (todo)
	bootTarget := config.DeviceConfig{
		Version: "1",
		Type:    "boot_target",
		Model:   "bmc-boot-target",
		Location: config.Location{
			Rack:  "ipmi",
			Board: conn.Hostname,
		},
		Data: map[string]string{
			"path":      conn.Path,
			"hostname":  conn.Hostname,
			"port":      strconv.Itoa(conn.Port),
			"username":  conn.Username,
			"password":  conn.Password,
			"interface": conn.Interface,
		},
	}
	devices = append(devices, &bootTarget)

	// Make new identify device for the BMC (todo)

	return devices, nil
}
