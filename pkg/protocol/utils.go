package protocol

import (
	"github.com/mitchellh/mapstructure"
	ipmi "github.com/vapor-ware/goipmi"
)

// NewClientFromConfig is a utility function to create a new IPMI client
// using the configuration specified in a Device's Data field.
func NewClientFromConfig(config map[string]interface{}) (*ipmi.Client, error) {
	conn := &ipmi.Connection{}

	if err := mapstructure.Decode(config, conn); err != nil {
		return nil, err
	}
	return ipmi.NewClient(conn)
}
