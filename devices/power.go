package devices

import (
	"github.com/mitchellh/mapstructure"
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/logger"
	"github.com/vmware/goipmi"
	"strconv"
	"fmt"
	"strings"
)

// BmcPower is the handler for the bmc-power device.
var BmcPower = sdk.DeviceHandler{
	Type:  "power",
	Model: "bmc-power",

	Read: bmcPowerRead,
	Write: bmcPowerWrite,
}


func makeConnection(data map[string]string) (*ipmi.Connection, error) {
	// FIXME (etd): need to do some type casting because the device
	// data is a map[string]string, but we need some values as ints
	var tmp = make(map[string]interface{})
	for k, v := range data {
		tmp[k] = v
	}

	port, err := strconv.Atoi(data["port"])
	if err != nil {
		return nil, err
	}
	tmp["port"] = port

	conn := &ipmi.Connection{}
	err = mapstructure.Decode(tmp, conn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}


func bmcPowerRead(device *sdk.Device) ([]*sdk.Reading, error) {

	conn, err := makeConnection(device.Data)
	if err != nil {
		return nil, err
	}

	logger.Infof("bmcPowerRead Conn: %+v", conn)

	client, err := ipmi.NewClient(conn)
	if err != nil {
		return nil, err
	}

	ipmiReq := &ipmi.Request{
		NetworkFunction: ipmi.NetworkFunctionChassis,
		Command:         ipmi.CommandChassisStatus,
		Data:            &ipmi.ChassisStatusRequest{},
	}
	ipmiRes := &ipmi.ChassisStatusResponse{}

	err = client.Send(ipmiReq, ipmiRes)
	if err != nil {
		return nil, err
	}

	logger.Infof("IPMI power response state: %+v", ipmiRes.State)
	logger.Infof("IPMI power response powerState: %+v", ipmiRes.PowerState)

	var state string
	switch ipmiRes.PowerState {
	case 0:
		state = "off"
	case 1:
		state = "on"
	default:
		return nil, fmt.Errorf("unknown power state response: %v", ipmiRes.PowerState)
	}

	ret := []*sdk.Reading{
		sdk.NewReading("state", state),
	}

	return ret, nil
}


func bmcPowerWrite(device *sdk.Device, data *sdk.WriteData) error {

	action := data.Action
	raw := data.Raw

	conn, err := makeConnection(device.Data)
	if err != nil {
		return err
	}

	client, err := ipmi.NewClient(conn)
	if err != nil {
		return err
	}

	// When writing to a BMC Power device, we always expect there to be
	// raw data specified. If there isn't, we return an error.
	if len(raw) == 0 {
		return fmt.Errorf("no values specified for 'raw', but required")
	}

	if action == "state" {
		cmd := string(raw[0])

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

		err = client.Control(state)
		if err != nil {
			return err
		}

	} else {
		// If we reach here, then the specified action is not supported.
		return fmt.Errorf("action '%s' is not supported for bmc power devices", action)
	}

	return nil
}
