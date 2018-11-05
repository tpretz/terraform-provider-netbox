package main

import (
	"github.com/Preskton/terraform-provider-netbox/plugin/providers/netbox"
	"github.com/hashicorp/terraform/plugin"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	log.Info("Loading terraform-provider-netbox plugin")

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: netbox.Provider,
	})
}
