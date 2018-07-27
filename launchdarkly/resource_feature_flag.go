package launchdarkly

import (
	"fmt"

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
	key := d.Get("key").(string)
	name := d.Get("name").(string)
	project := d.Get("project").(string)

	params := feature_flags.NewPostFeatureFlagParams().
		WithProjectKey(project).
		WithFeatureFlagBody(feature_flags.PostFeatureFlagBody{
			Key:       &key,
			Name:      &name,
			Temporary: !d.Get("permanent").(bool),
		})

	_, err := meta.LaunchDarkly.FeatureFlags.PostFeatureFlag(params, meta.AuthInfo)
	if err != nil {
		return fmt.Errorf("Failed to create flag %q in project %q: %s", key, project, err)
	}

	d.SetId(key)
	return nil
}

func resourceFeatureFlagRead(d *schema.ResourceData, metaRaw interface{}) error {
	meta := metaRaw.(*Meta)
	key := d.Get("key").(string)
	project := d.Get("project").(string)

	params := feature_flags.NewGetFeatureFlagParams().
		WithFeatureFlagKey(key).
		WithProjectKey(project)

	flag, err := meta.LaunchDarkly.FeatureFlags.GetFeatureFlag(params, meta.AuthInfo)
	if err != nil {
		return fmt.Errorf("Failed to get flag %q in project %q: %s", key, project, err)
	}

	d.SetId(flag.Payload.Key)
	return nil
}

func resourceFeatureFlagUpdate(d *schema.ResourceData, metaRaw interface{}) error {
	// TODO
	return nil
}

func resourceFeatureFlagDelete(d *schema.ResourceData, metaRaw interface{}) error {
	meta := metaRaw.(*Meta)
	key := d.Get("key").(string)
	project := d.Get("project").(string)

	params := feature_flags.NewDeleteFeatureFlagParams().
		WithFeatureFlagKey(key).
		WithProjectKey(project)

	_, err := meta.LaunchDarkly.FeatureFlags.DeleteFeatureFlag(params, meta.AuthInfo)
	if err != nil {
		return fmt.Errorf("Failed to delete flag %q from project %q: %s", key, project, err)
	}

	return nil
}
