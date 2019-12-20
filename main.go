package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/rancher/terraform-provider-rke/rke"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: rke.Provider,
	})
}
