package protocol

import (
	log "github.com/sirupsen/logrus"
	ipmi "github.com/vapor-ware/goipmi"
)

// GetChassisBootTarget gets the current boot device for the chassis.
func GetChassisBootTarget(client *ipmi.Client) (string, error) {
	request := &ipmi.Request{
		NetworkFunction: ipmi.NetworkFunctionChassis,
		Command:         ipmi.CommandGetSystemBootOptions,
		Data: &ipmi.SystemBootOptionsRequest{
			Param: ipmi.BootParamBootFlags,
		},
	}
	response := &ipmi.SystemBootOptionsResponse{}

	if err := client.Send(request, response); err != nil {
		return "", err
	}

	// As per Section 28.13 Get System Boot Options, Table 28-, Boot Option Parameters,
	// index 1 of the configuration data holds the service partition selector, i.e.
	// the boot target.
	target := ipmi.BootDevice(response.Data[1])
	return target.String(), nil
}

// SetChassisBootTarget sets the boot device for the chassis.
func SetChassisBootTarget(client *ipmi.Client, target ipmi.BootDevice) error {
	log.WithFields(log.Fields{
		"target": target.String(),
	}).Info("[ipmi] setting chassis boot target")
	return client.SetBootDevice(target)
}
