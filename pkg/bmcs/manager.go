package bmcs

import (
	"fmt"

	ipmi "github.com/vapor-ware/goipmi"
)

var registeredBMCs = map[string]*ipmi.Client{}

// Get the IPMI client for the given BMC ID.
func Get(bmcID string) (*ipmi.Client, error) {
	client, found := registeredBMCs[bmcID]
	if !found {
		return nil, fmt.Errorf("no BMC with id %s found", bmcID)
	}
	return client, nil
}

// Add a new IPMI client for the given BMC ID.
func Add(bmcID string, client *ipmi.Client) {
	registeredBMCs[bmcID] = client
}
