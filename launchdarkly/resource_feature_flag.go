package launchdarkly

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mlafeldt/go-launchdarkly/client/feature_flags"
	"github.com/mlafeldt/go-launchdarkly/models"
)

const defaultProject = "default"

func resourceFeatureFlag() *schema.Resource {
	return &schema.Resource{
		Create: resourceFeatureFlagCreate,
		Read:   resourceFeatureFlagRead,
		Update: resourceFeatureFlagUpdate,
		Delete: resourceFeatureFlagDelete,
		Exists: resourceFeatureFlagExists,

		Importer: &schema.ResourceImporter{
			State: resourceFeatureFlagImport,
		},

		Schema: map[string]*schema.Schema{
			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  defaultProject,
			},
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"include_in_snippet": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"temporary": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceFeatureFlagCreate(d *schema.ResourceData, metaRaw interface{}) error {
	meta := metaRaw.(*Meta)
	project := d.Get("project").(string)
	key := d.Get("key").(string)
	name := d.Get("name").(string)
	tags := stringList(d.Get("tags").([]interface{}))

	params := feature_flags.NewPostFeatureFlagParams().
		WithProjectKey(project).
		WithFeatureFlagBody(feature_flags.PostFeatureFlagBody{
			Key:              &key,
			Name:             &name,
			Description:      d.Get("description").(string),
			Tags:             tags,
			IncludeInSnippet: d.Get("include_in_snippet").(bool),
			Temporary:        d.Get("temporary").(bool),
		})

	_, err := meta.LaunchDarkly.FeatureFlags.PostFeatureFlag(params, meta.AuthInfo)
	if err != nil {
		return fmt.Errorf("Failed to create flag %q in project %q: %s", key, project, err)
	}

	d.SetId(key)
	return resourceFeatureFlagRead(d, metaRaw)
}

func resourceFeatureFlagRead(d *schema.ResourceData, metaRaw interface{}) error {
	meta := metaRaw.(*Meta)
	project := d.Get("project").(string)
	key := d.Get("key").(string)

	params := feature_flags.NewGetFeatureFlagParams().
		WithProjectKey(project).
		WithFeatureFlagKey(key)

	flag, err := meta.LaunchDarkly.FeatureFlags.GetFeatureFlag(params, meta.AuthInfo)
	if err != nil {
		return fmt.Errorf("Failed to get flag %q of project %q: %s", key, project, err)
	}

	d.Set("key", flag.Payload.Key)
	d.Set("name", flag.Payload.Name)
	d.Set("description", flag.Payload.Description)
	d.Set("tags", flag.Payload.Tags)
	d.Set("include_in_snippet", flag.Payload.IncludeInSnippet)
	d.Set("temporary", flag.Payload.Temporary)
	return nil
}

func resourceFeatureFlagUpdate(d *schema.ResourceData, metaRaw interface{}) error {
	meta := metaRaw.(*Meta)
	project := d.Get("project").(string)
	key := d.Get("key").(string)
	name := d.Get("name").(string)
	tags := stringList(d.Get("tags").([]interface{}))

	patch := feature_flags.PatchFeatureFlagBody{
		Comment: "Terraform",
		Patch: []*models.PatchOperation{
			{
				Op:    stringPtr("replace"),
				Path:  stringPtr("/name"),
				Value: name,
			},
			{
				Op:    stringPtr("replace"),
				Path:  stringPtr("/description"),
				Value: d.Get("description").(string),
			},
			{
				Op:    stringPtr("replace"),
				Path:  stringPtr("/tags"),
				Value: tags,
			},
			{
				Op:    stringPtr("replace"),
				Path:  stringPtr("/includeInSnippet"),
				Value: d.Get("include_in_snippet").(bool),
			},
			{
				Op:    stringPtr("replace"),
				Path:  stringPtr("/temporary"),
				Value: d.Get("temporary").(bool),
			},
		},
	}

	params := feature_flags.NewPatchFeatureFlagParams().
		WithProjectKey(project).
		WithFeatureFlagKey(key).
		WithPatchComment(patch)

	_, err := meta.LaunchDarkly.FeatureFlags.PatchFeatureFlag(params, meta.AuthInfo)
	if err != nil {
		return fmt.Errorf("Failed to update flag %q in project %q: %s", key, project, err)
	}

	return resourceFeatureFlagRead(d, metaRaw)
}

func resourceFeatureFlagDelete(d *schema.ResourceData, metaRaw interface{}) error {
	meta := metaRaw.(*Meta)
	project := d.Get("project").(string)
	key := d.Get("key").(string)

	params := feature_flags.NewDeleteFeatureFlagParams().
		WithProjectKey(project).
		WithFeatureFlagKey(key)

	_, err := meta.LaunchDarkly.FeatureFlags.DeleteFeatureFlag(params, meta.AuthInfo)
	if err != nil {
		return fmt.Errorf("Failed to delete flag %q from project %q: %s", key, project, err)
	}

	return nil
}

func resourceFeatureFlagExists(d *schema.ResourceData, metaRaw interface{}) (bool, error) {
	meta := metaRaw.(*Meta)
	project := d.Get("project").(string)
	key := d.Get("key").(string)

	params := feature_flags.NewGetFeatureFlagParams().
		WithProjectKey(project).
		WithFeatureFlagKey(key)

	_, err := meta.LaunchDarkly.FeatureFlags.GetFeatureFlag(params, meta.AuthInfo)
	if err == nil {
		return true, nil
	}
	if _, notFound := err.(*feature_flags.GetFeatureFlagNotFound); notFound {
		return false, nil
	}
	return false, fmt.Errorf("Failed to check flag %q of project %q: %s", key, project, err)
}

func resourceFeatureFlagImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	project := defaultProject
	key := d.Id()

	if strings.Contains(d.Id(), "/") {
		parts := strings.SplitN(d.Id(), "/", 2)

		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("ID must have format <key> or <project>/<key>")
		}

		project, key = parts[0], parts[1]
	}

	d.Set("project", project)
	d.Set("key", key)
	d.SetId(key)

	if err := resourceFeatureFlagRead(d, meta); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
