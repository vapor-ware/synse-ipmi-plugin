package protocol

import (
	"fmt"

	"github.com/vapor-ware/goipmi"
	"github.com/vapor-ware/synse-sdk/sdk/logger"
)

// GetChassisPowerState gets the current state (on/off) of the chassis.
func GetChassisPowerState(config map[string]string) (string, error) {
	client, err := newClientFromConfig(config)
	if err != nil {
		return "", err
	}

	request := &ipmi.Request{
		NetworkFunction: ipmi.NetworkFunctionChassis,
		Command:         ipmi.CommandChassisStatus,
		Data:            &ipmi.ChassisStatusRequest{},
	}
	response := &ipmi.ChassisStatusResponse{}

	err = client.Send(request, response)
	if err != nil {
		return "", err
	}

	var state string
	switch uint8(response.PowerState) & 1 {
	case 0:
		state = "off"
	case 1:
		state = "on"
	default:
		return "", fmt.Errorf("unknown power state response: %v", response.PowerState)
	}

	return state, nil
}

// SetChassisPowerState sets the state of the chassis.
func SetChassisPowerState(config map[string]string, control ipmi.ChassisControl) error {
	client, err := newClientFromConfig(config)
	if err != nil {
		return err
	}

	logger.Debugf("Setting power state to: %s", control.String())
	return client.Control(control)
}
