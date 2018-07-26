package launchdarkly

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mlafeldt/go-launchdarkly/client/feature_flags"
)

func resourceFeatureFlag() *schema.Resource {
	return &schema.Resource{
		Create: resourceFeatureFlagCreate,
		Read:   resourceFeatureFlagRead,
		Update: resourceFeatureFlagUpdate,
		Delete: resourceFeatureFlagDelete,

		Schema: map[string]*schema.Schema{
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"permanent": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},
		},
	}
}

func resourceFeatureFlagCreate(d *schema.ResourceData, metaRaw interface{}) error {
	meta := metaRaw.(*Meta)

	params := feature_flags.NewPostFeatureFlagParams()
	params.SetProjectKey(d.Get("project").(string))

	key := d.Get("key").(string)
	name := d.Get("name").(string)
	params.SetFeatureFlagBody(feature_flags.PostFeatureFlagBody{
		Key:       &key,
		Name:      &name,
		Temporary: !d.Get("permanent").(bool),
	})

	_, err := meta.LaunchDarkly.FeatureFlags.PostFeatureFlag(params, meta.AuthInfo)
	return err
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
