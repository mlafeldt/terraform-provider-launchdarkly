package launchdarkly

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mlafeldt/go-launchdarkly/client/environments"
	"github.com/mlafeldt/go-launchdarkly/models"
)

func resourceEnvironment() *schema.Resource {
	return &schema.Resource{
		Create: resourceEnvironmentCreate,
		Read:   resourceEnvironmentRead,
		Update: resourceEnvironmentUpdate,
		Delete: resourceEnvironmentDelete,
		Exists: resourceEnvironmentExists,

		Importer: &schema.ResourceImporter{
			State: resourceEnvironmentImport,
		},

		Schema: map[string]*schema.Schema{
			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  defaultProject,
				ForceNew: true,
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
			"color": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"default_ttl": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"api_key": &schema.Schema{
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"mobile_key": &schema.Schema{
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceEnvironmentCreate(d *schema.ResourceData, metaRaw interface{}) error {
	meta := metaRaw.(*Meta)
	project := d.Get("project").(string)
	key := d.Get("key").(string)

	params := environments.NewPostEnvironmentParams().
		WithProjectKey(project).
		WithEnvironmentBody(&models.EnvironmentPost{
			Key:        &key,
			Name:       stringPtr(d.Get("name").(string)),
			Color:      stringPtr(d.Get("color").(string)),
			DefaultTTL: float64(d.Get("default_ttl").(int)),
		})

	_, err := meta.LaunchDarkly.Environments.PostEnvironment(params, meta.AuthInfo)
	if err != nil {
		return fmt.Errorf("Failed to create env %q in project %q: %s", key, project, err)
	}

	d.SetId(key)
	return resourceEnvironmentRead(d, metaRaw)
}

func resourceEnvironmentRead(d *schema.ResourceData, metaRaw interface{}) error {
	meta := metaRaw.(*Meta)
	project := d.Get("project").(string)
	key := d.Get("key").(string)

	params := environments.NewGetEnvironmentParams().
		WithProjectKey(project).
		WithEnvironmentKey(key)

	env, err := meta.LaunchDarkly.Environments.GetEnvironment(params, meta.AuthInfo)
	if err != nil {
		return fmt.Errorf("Failed to get env %q of project %q: %s", key, project, err)
	}

	d.Set("key", env.Payload.Key)
	d.Set("name", env.Payload.Name)
	d.Set("color", env.Payload.Color)
	d.Set("default_ttl", env.Payload.DefaultTTL)
	d.Set("api_key", env.Payload.APIKey)
	d.Set("mobile_key", env.Payload.MobileKey)
	return nil
}

func resourceEnvironmentUpdate(d *schema.ResourceData, metaRaw interface{}) error {
	meta := metaRaw.(*Meta)
	project := d.Get("project").(string)
	key := d.Get("key").(string)

	patch := []*models.PatchOperation{
		{
			Op:    stringPtr("replace"),
			Path:  stringPtr("/name"),
			Value: d.Get("name").(string),
		},
		{
			Op:    stringPtr("replace"),
			Path:  stringPtr("/color"),
			Value: d.Get("color").(string),
		},
		{
			Op:    stringPtr("replace"),
			Path:  stringPtr("/defaultTtl"),
			Value: d.Get("default_ttl").(int),
		},
	}

	params := environments.NewPatchEnvironmentParams().
		WithProjectKey(project).
		WithEnvironmentKey(key).
		WithPatchDelta(patch)

	_, err := meta.LaunchDarkly.Environments.PatchEnvironment(params, meta.AuthInfo)
	if err != nil {
		return fmt.Errorf("Failed to update env %q in project %q: %s", key, project, err)
	}

	return resourceEnvironmentRead(d, metaRaw)
}

func resourceEnvironmentDelete(d *schema.ResourceData, metaRaw interface{}) error {
	meta := metaRaw.(*Meta)
	project := d.Get("project").(string)
	key := d.Get("key").(string)

	params := environments.NewDeleteEnvironmentParams().
		WithProjectKey(project).
		WithEnvironmentKey(key)

	_, err := meta.LaunchDarkly.Environments.DeleteEnvironment(params, meta.AuthInfo)
	if err != nil {
		return fmt.Errorf("Failed to delete env %q from project %q: %s", key, project, err)
	}

	return nil
}

func resourceEnvironmentExists(d *schema.ResourceData, metaRaw interface{}) (bool, error) {
	meta := metaRaw.(*Meta)
	project := d.Get("project").(string)
	key := d.Get("key").(string)

	params := environments.NewGetEnvironmentParams().
		WithProjectKey(project).
		WithEnvironmentKey(key)

	_, err := meta.LaunchDarkly.Environments.GetEnvironment(params, meta.AuthInfo)
	if err == nil {
		return true, nil
	}
	if _, notFound := err.(*environments.GetEnvironmentNotFound); notFound {
		return false, nil
	}
	return false, fmt.Errorf("Failed to check if env %q exists in project %q: %s", key, project, err)
}

func resourceEnvironmentImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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

	if err := resourceEnvironmentRead(d, meta); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
