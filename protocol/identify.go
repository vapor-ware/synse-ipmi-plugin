package protocol

import (
	"fmt"

	"github.com/vmware/goipmi"
)

type IdentifyState uint8

const (
	CommandChassisIdentify = ipmi.Command(0x04)
)

const (
	ChassisIdentifySupported    = 0x40
	ChassisIdentifyMask         = 0x30
	ChassisIdentifyOff          = 0x00
	ChassisIdentifyOn           = 0x10
	ChassisIdentifyOnIndefinite = 0x20
	ChassisIdentifyReserved     = 0x30
)

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

func Identify(c *ipmi.Client, time int, indefinately bool) error {
	if time > 255 || time < 0 {
		return fmt.Errorf("invalid time value: %d", time)
	}

	var force uint8
	if indefinately {
		force = 1
	} else {
		force = 0
	}

	request := &ipmi.Request{
		NetworkFunction: ipmi.NetworkFunctionChassis,
		Command:         CommandChassisIdentify,
		Data: &ChassisIdentifyRequest{
			Interval: uint8(time),
			Force:    force,
		},
	}
	return c.Send(request, &ChassisIdentifyResponse{})
}

func GetChassisIdentify(config map[string]string) (string, error) {
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

	identifyState := GetIdentifyState(response)
	var state string

	switch identifyState {
	case IdentifyOn:
	case IdentifyIndefinite:
		state = "on"
	case IdentifyOff:
		state = "off"
	case IdentifyUnsupported:
		return "", fmt.Errorf("chassis identify unsupported (identify state: %v)", identifyState)
	default:
		return "", fmt.Errorf("unsupported identify state: %v", identifyState)
	}

	return state, err
}

func SetChassisIdentify(config map[string]string, state IdentifyState) error {
	client, err := newClientFromConfig(config)
	if err != nil {
		return err
	}

	var time int
	switch state {
	case IdentifyOn:
		time = 15 // 15s is the default interval
	case IdentifyOff:
		time = 0
	default:
		return fmt.Errorf("identify state unsupported for setting: %v", state)
	}

	return Identify(
		client,
		time,
		false, // for now, never turn it on indefinitely
	)
}
