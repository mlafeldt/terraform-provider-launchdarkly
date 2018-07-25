package launchdarkly

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"sdk_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LAUNCHDARKLY_SDK_KEY", nil),
				Description: "The LaunchDarkly SDK key",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"launchdarkly_feature_flag": resourceFeatureFlag(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		SDKKey: d.Get("sdk_key").(string),
	}
	return config.Client()
}
