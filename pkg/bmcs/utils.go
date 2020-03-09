package bmcs

import "fmt"

// GetIDFromConfig gets the BMC ID which is embedded in the device config on
// dynamic configuration at startup.
func GetIDFromConfig(config map[string]interface{}) (string, error) {
	bmcID, found := config["bmc"]
	if !found {
		return "", fmt.Errorf("BMC ID not found in device data")
	}
	conv, ok := bmcID.(string)
	if !ok {
		return "", fmt.Errorf("BMC ID stored in device data in invalid format: must be string")
	}
	return conv, nil
}
