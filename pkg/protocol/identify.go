package protocol

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	ipmi "github.com/vapor-ware/goipmi"
)

// FIXME - now that we have a fork of goipmi, some of this stuff should
// probably live there

// IdentifyState represents the state of the identify device
// on the BMC managed chassis.
type IdentifyState uint8

// IPMI command for chassis identify.
const (
	CommandChassisIdentify = ipmi.Command(0x04)
)

// Chassis identify IPMI byte constants.
const (
	ChassisIdentifySupported    = 0x40
	ChassisIdentifyMask         = 0x30
	ChassisIdentifyOff          = 0x00
	ChassisIdentifyOn           = 0x10
	ChassisIdentifyOnIndefinite = 0x20
	ChassisIdentifyReserved     = 0x30
)

// States for the chassis identify device.
const (
	IdentifyOff         = IdentifyState(0x0)
	IdentifyOn          = IdentifyState(0x1)
	IdentifyIndefinite  = IdentifyState(0x2)
	IdentifyReserved    = IdentifyState(0x3)
	IdentifyUnsupported = IdentifyState(0x4)
)

// ChassisIdentifyRequest per Section 28.5
type ChassisIdentifyRequest struct {
	Interval uint8
	Force    uint8
}

// ChassisIdentifyResponse per Section 28.5
type ChassisIdentifyResponse struct {
	ipmi.CompletionCode
}

// GetIdentifyState gets the current state of the chassis identify.
func GetIdentifyState(status *ipmi.ChassisStatusResponse) IdentifyState {
	if (status.State & ChassisIdentifySupported) == ChassisIdentifySupported {
		switch status.State & ChassisIdentifyMask {
		case ChassisIdentifyOff:
			return IdentifyOff
		case ChassisIdentifyOn:
			return IdentifyOn
		case ChassisIdentifyOnIndefinite:
			return IdentifyIndefinite
		case ChassisIdentifyReserved:
			return IdentifyReserved
		}
	}
	return IdentifyUnsupported
}

// Identify issues a request to enable/disable the chassis identify device.
func Identify(c *ipmi.Client, time int) error {
	if time > 255 || time < 0 {
		return fmt.Errorf("invalid time value: %d, must be between 0 and 255", time)
	}

	request := &ipmi.Request{
		NetworkFunction: ipmi.NetworkFunctionChassis,
		Command:         CommandChassisIdentify,
		Data: &ChassisIdentifyRequest{
			Interval: uint8(time),
		},
	}
	return c.Send(request, &ChassisIdentifyResponse{})
}

// GetChassisIdentify gets the current identify state from the chassis.
func GetChassisIdentify(client *ipmi.Client) (string, error) {
	request := &ipmi.Request{
		NetworkFunction: ipmi.NetworkFunctionChassis,
		Command:         ipmi.CommandChassisStatus,
		Data:            &ipmi.ChassisStatusRequest{},
	}
	response := &ipmi.ChassisStatusResponse{}

	if err := client.Send(request, response); err != nil {
		return "", err
	}

	identifyState := GetIdentifyState(response)
	var state string

	switch identifyState {
	case IdentifyOn:
	case IdentifyIndefinite:
		state = "on"
	case IdentifyOff:
		state = "off"
	case IdentifyUnsupported:
		// Unclear that this should be an error, since this is optionally
		// specified. This setting is optional for IPMI 2.0 according to the
		// spec and it just means that the "chassis identify command support
		// [is] unspecified via this command", not that identify support is
		// unsupported by the BMC.
		return "", fmt.Errorf("chassis identify not specified as supported (identify state: %v)", identifyState)
	default:
		return "", fmt.Errorf("unsupported identify state: %v", identifyState)
	}

	return state, nil
}

// SetChassisIdentify sets the identify state of the chassis
func SetChassisIdentify(client *ipmi.Client, state IdentifyState) error {
	var time int
	switch state {
	case IdentifyOn:
		time = 15 // 15s is the default interval
	case IdentifyOff:
		time = 0
	default:
		return fmt.Errorf("identify state unsupported for setting: %v", state)
	}

	log.WithFields(log.Fields{
		"duration": time,
	}).Info("[ipmi] setting chassis identify")

	return Identify(
		client,
		time,
	)
}
