package protocol

import (
	"github.com/vapor-ware/goipmi"
	"github.com/vapor-ware/synse-sdk/sdk/logger"
)

const (
	powerOn = 0x01
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

	// Check the power state. According to the IPMI spec, Section 28.2, table 28:
	// Get Chassis Status Command, power on/off state is held in bit 0 of the
	// Current Power State byte, where 1b = system power is on, 0b = system power
	// is off
	if response.PowerState&1 == powerOn {
		state = "on"
	} else {
		state = "off"
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
