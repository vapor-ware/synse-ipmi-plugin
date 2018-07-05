package main

import (
	log "github.com/Sirupsen/logrus"

	"github.com/vapor-ware/synse-ipmi-plugin/pkg"
	"github.com/vapor-ware/synse-sdk/sdk"
)

const (
	pluginName       = "ipmi"
	pluginMaintainer = "vaporio"
	pluginDesc       = "A general-purpose IPMI plugin"
	pluginVcs        = "https://github.com/vapor-ware/synse-ipmi-plugin"
)

func main() {
	// Set the plugin metadata
	sdk.SetPluginMeta(
		pluginName,
		pluginMaintainer,
		pluginDesc,
		pluginVcs,
	)

	plugin := pkg.MakePlugin()

	// Run the plugin
	if err := plugin.Run(); err != nil {
		log.Fatal(err)
	}
}
