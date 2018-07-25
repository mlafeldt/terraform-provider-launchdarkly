package launchdarkly

import "github.com/hashicorp/terraform/helper/schema"

func resourceFeatureFlag() *schema.Resource {
	return &schema.Resource{
		Create: resourceFeatureFlagCreate,
		Read:   resourceFeatureFlagRead,
		Update: resourceFeatureFlagUpdate,
		Delete: resourceFeatureFlagDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     schema.TypeString,
				Optional: true,
			},

			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceFeatureFlagCreate(d *schema.ResourceData, metaRaw interface{}) error {
	return nil
}

func resourceFeatureFlagRead(d *schema.ResourceData, metaRaw interface{}) error {
	return nil
}

func resourceFeatureFlagUpdate(d *schema.ResourceData, metaRaw interface{}) error {
	return nil
}

func resourceFeatureFlagDelete(d *schema.ResourceData, metaRaw interface{}) error {
	return nil
}
