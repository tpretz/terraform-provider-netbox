package main

import (
	"github.com/Preskton/terraform-provider-netbox/plugin/providers/netbox"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: netbox.Provider,
	})
}
