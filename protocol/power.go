package protocol

import (
	"github.com/vmware/goipmi"
	"fmt"
)


func GetChassisPowerState(config map[string]string) (string, error) {
	client, err := newClientFromConfig(config)
	if err != nil {
		return "", err
	}

	request := &ipmi.Request{
		NetworkFunction: ipmi.NetworkFunctionChassis,
		Command: ipmi.CommandChassisStatus,
		Data: &ipmi.ChassisStatusRequest{},
	}
	response := &ipmi.ChassisStatusResponse{}

	err = client.Send(request, response)
	if err != nil {
		return "", err
	}

	var state string
	switch response.PowerState {
	case 0:
		state = "off"
	case 1:
		state = "on"
	default:
		return "", fmt.Errorf("unknown power state response: %v", response.PowerState)
	}

	return state, nil
}


func SetChassisPowerState(config map[string]string, control ipmi.ChassisControl) error {
	client, err := newClientFromConfig(config)
	if err != nil {
		return err
	}

	return client.Control(control)
}