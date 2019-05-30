package protocol

import (
	"github.com/mitchellh/mapstructure"
	ipmi "github.com/vapor-ware/goipmi"
)

// newClientFromConfig is a utility function to create a new IPMI client
// using the configuration specified in a Device's Data field.
func newClientFromConfig(config map[string]interface{}) (*ipmi.Client, error) {
	conn, err := makeConnection(config)
	if err != nil {
		return nil, err
	}
	return ipmi.NewClient(conn)
}

// makeConnection is a utility function to initialize the ipmi.Connection
// used for an ipmi.Client by parsing the Device's Data map.
func makeConnection(data map[string]interface{}) (*ipmi.Connection, error) {
	conn := &ipmi.Connection{}
	err := mapstructure.Decode(data, conn)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
