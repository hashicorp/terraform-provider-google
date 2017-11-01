package google

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGoogleProjectService() *schema.Resource {
	return &schema.Resource{
		Create: resourceGoogleProjectServiceCreate,
		Read:   resourceGoogleProjectServiceRead,
		Delete: resourceGoogleProjectServiceDelete,

		Schema: map[string]*schema.Schema{
			"service": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceGoogleProjectServiceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	srv := d.Get("service").(string)

	if err = enableService(srv, project, config); err != nil {
		return fmt.Errorf("Error enabling service: %s", err)
	}

	d.SetId(srv)
	return resourceGoogleProjectServiceRead(d, meta)
}

func resourceGoogleProjectServiceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	services, err := getApiServices(project, config, map[string]struct{}{})
	if err != nil {
		return err
	}

	for _, s := range services {
		if s == d.Id() {
			d.Set("service", s)
			return nil
		}
	}

	// The service is not enabled server-side, so remove it from state
	d.SetId("")
	return nil
}

func resourceGoogleProjectServiceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	if err = disableService(d.Id(), project, config); err != nil {
		return fmt.Errorf("Error disabling service: %s", err)
	}

	d.SetId("")
	return nil
}
