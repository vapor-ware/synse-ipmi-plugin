package protocol

import (
	log "github.com/sirupsen/logrus"
	ipmi "github.com/vapor-ware/goipmi"
)

const (
	powerOn = 0x01
)

// GetChassisPowerState gets the current state (on/off) of the chassis.
func GetChassisPowerState(config map[string]interface{}) (string, error) {
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

	log.WithFields(log.Fields{
		"state":  state,
	}).Debug("[ipmi] got chassis power state")

	return state, nil
}

// SetChassisPowerState sets the state of the chassis.
func SetChassisPowerState(config map[string]interface{}, control ipmi.ChassisControl) error {
	client, err := newClientFromConfig(config)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"state":  control.String(),
	}).Info("[ipmi] setting chassis power state")
	return client.Control(control)
}
