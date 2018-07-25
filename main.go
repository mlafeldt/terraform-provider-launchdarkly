package main

import (
	"github.com/hashicorp/terraform/plugin"

	"github.com/mlafeldt/terraform-provider-launchdarkly/launchdarkly"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: launchdarkly.Provider})
}
