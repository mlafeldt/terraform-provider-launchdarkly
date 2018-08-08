package launchdarkly

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mlafeldt/go-launchdarkly/client/projects"
	"github.com/mlafeldt/go-launchdarkly/models"
)

const defaultProject = "default"

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,
		Exists: resourceProjectExists,

		Importer: &schema.ResourceImporter{
			State: resourceProjectImport,
		},

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
		},
	}
}

func resourceProjectCreate(d *schema.ResourceData, metaRaw interface{}) error {
	meta := metaRaw.(*Meta)
	key := d.Get("key").(string)

	params := projects.NewPostProjectParams().
		WithProjectBody(projects.PostProjectBody{
			Key:  &key,
			Name: stringPtr(d.Get("name").(string)),
		})

	_, err := meta.LaunchDarkly.Projects.PostProject(params, meta.AuthInfo)
	if err != nil {
		return fmt.Errorf("Failed to create project with key %q: %s", key, err)
	}

	d.SetId(key)
	return resourceProjectRead(d, metaRaw)
}

func resourceProjectRead(d *schema.ResourceData, metaRaw interface{}) error {
	meta := metaRaw.(*Meta)
	key := d.Get("key").(string)

	params := projects.NewGetProjectParams().WithProjectKey(key)

	project, err := meta.LaunchDarkly.Projects.GetProject(params, meta.AuthInfo)
	if err != nil {
		return fmt.Errorf("Failed to get project with key %q: %s", key, err)
	}

	d.Set("key", project.Payload.Key)
	d.Set("name", project.Payload.Name)
	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, metaRaw interface{}) error {
	meta := metaRaw.(*Meta)
	key := d.Get("key").(string)

	patch := []*models.PatchOperation{
		{
			Op:    stringPtr("replace"),
			Path:  stringPtr("/name"),
			Value: d.Get("name").(string),
		},
	}

	params := projects.NewPatchProjectParams().
		WithProjectKey(key).
		WithPatchDelta(patch)

	_, err := meta.LaunchDarkly.Projects.PatchProject(params, meta.AuthInfo)
	if err != nil {
		return fmt.Errorf("Failed to update project with key %q: %s", key, err)
	}

	return resourceProjectRead(d, metaRaw)
}

func resourceProjectDelete(d *schema.ResourceData, metaRaw interface{}) error {
	meta := metaRaw.(*Meta)
	key := d.Get("key").(string)

	params := projects.NewDeleteProjectParams().WithProjectKey(key)

	_, err := meta.LaunchDarkly.Projects.DeleteProject(params, meta.AuthInfo)
	if err != nil {
		return fmt.Errorf("Failed to delete project with key %q: %s", key, err)
	}

	return nil
}

func resourceProjectExists(d *schema.ResourceData, metaRaw interface{}) (bool, error) {
	meta := metaRaw.(*Meta)
	key := d.Get("key").(string)

	params := projects.NewGetProjectParams().WithProjectKey(key)

	_, err := meta.LaunchDarkly.Projects.GetProject(params, meta.AuthInfo)
	if err == nil {
		return true, nil
	}
	if _, notFound := err.(*projects.GetProjectNotFound); notFound {
		return false, nil
	}
	return false, fmt.Errorf("Failed to check if project with key %q exists: %s", key, err)
}

func resourceProjectImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	d.SetId(d.Id())
	d.Set("key", d.Id())

	if err := resourceProjectRead(d, meta); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
