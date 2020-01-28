package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/xanderflood/fruit-pi-server/provider"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
