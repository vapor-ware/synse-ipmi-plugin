package protocol

import (
	log "github.com/sirupsen/logrus"
	ipmi "github.com/vapor-ware/goipmi"
)

// GetChassisBootTarget gets the current boot device for the chassis.
func GetChassisBootTarget(config map[string]interface{}) (string, error) {
	client, err := newClientFromConfig(config)
	if err != nil {
		return "", err
	}

	request := &ipmi.Request{
		NetworkFunction: ipmi.NetworkFunctionChassis,
		Command:         ipmi.CommandGetSystemBootOptions,
		Data: &ipmi.SystemBootOptionsRequest{
			Param: ipmi.BootParamBootFlags,
		},
	}
	response := &ipmi.SystemBootOptionsResponse{}

	err = client.Send(request, response)
	if err != nil {
		return "", err
	}

	// As per Section 28.13 Get System Boot Options, Table 28-, Boot Option Parameters,
	// index 1 of the configuration data holds the service partition selector, i.e.
	// the boot target.
	target := ipmi.BootDevice(response.Data[1])
	return target.String(), nil
}

// SetChassisBootTarget sets the boot device for the chassis.
func SetChassisBootTarget(config map[string]interface{}, target ipmi.BootDevice) error {
	client, err := newClientFromConfig(config)
	if err != nil {
		return err
	}
	log.Debugf("Setting boot target to: %s", target.String())
	return client.SetBootDevice(target)
}
