package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strings"
)

func resourceGoogleProjectService() *schema.Resource {
	return &schema.Resource{
		Create: resourceGoogleProjectServiceCreate,
		Read:   resourceGoogleProjectServiceRead,
		Delete: resourceGoogleProjectServiceDelete,
		Update: resourceGoogleProjectServiceUpdate,

		Importer: &schema.ResourceImporter{
			State: resourceGoogleProjectServiceImport,
		},

		Schema: map[string]*schema.Schema{
			"service": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"disable_dependent_services": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"disable_on_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceGoogleProjectServiceImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid google_project_service id format for import, expecting `{project}/{service}`, found %s", d.Id())
	}
	d.Set("project", parts[0])
	d.Set("service", parts[1])
	return []*schema.ResourceData{d}, nil
}

func resourceGoogleProjectServiceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	srv := d.Get("service").(string)
	err = enableServiceUsageProjectServices([]string{srv}, project, config)
	if err != nil {
		return err
	}

	id, err := replaceVars(d, config, "{{project}}/{{service}}")
	if err != nil {
		return fmt.Errorf("unable to construct ID: %s", err)
	}
	d.SetId(id)
	return resourceGoogleProjectServiceRead(d, meta)
}

func resourceGoogleProjectServiceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	enabledServices, err := readEnabledServiceUsageProjectServices(project, config)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Project Service %s", d.Id()))
	}

	srv := d.Get("service")
	for _, s := range enabledServices {
		if s == srv {
			d.Set("project", project)
			d.Set("service", s)
			return nil
		}
	}

	// The service is was not found in enabled services - remove it from state
	log.Printf("[DEBUG] service %s not in enabled services for project %s, removing from state", srv, project)
	d.SetId("")
	return nil
}

func resourceGoogleProjectServiceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if disable := d.Get("disable_on_destroy"); !(disable.(bool)) {
		log.Printf("[WARN] Project service %q disable_on_destroy is false, skip disabling service", d.Id())
		d.SetId("")
		return nil
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	service := d.Get("service").(string)
	disableDependencies := d.Get("disable_dependent_services").(bool)
	if err = disableServiceUsageProjectService(service, project, config, disableDependencies); err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Project Service %s", d.Id()))
	}

	d.SetId("")
	return nil
}

func resourceGoogleProjectServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	// This update method is no-op because the only updatable fields
	// are state/config-only, i.e. they aren't sent in requests to the API.
	return nil
}
