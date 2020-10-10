package main

import (
	"github.com/configcat/terraform-provider-configcat/configcat"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return configcat.Provider()
		},
	})
}
