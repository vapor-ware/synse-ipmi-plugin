package protocol

import (
	"strconv"

	"github.com/mitchellh/mapstructure"
	"github.com/vmware/goipmi"
)

// newClientFromConfig is a utility function to create a new IPMI client
// using the configuration specified in a Device's Data field.
func newClientFromConfig(config map[string]string) (*ipmi.Client, error) {
	conn, err := makeConnection(config)
	if err != nil {
		return nil, err
	}
	return ipmi.NewClient(conn)
}

// makeConnection is a utility function to initialize the ipmi.Connection
// used for an ipmi.Client by parsing the Device's Data map.
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
