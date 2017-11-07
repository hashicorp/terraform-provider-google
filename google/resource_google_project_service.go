package google

import (
	"fmt"
	"strings"

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

	d.SetId(projectServiceId{project, srv}.terraformId())
	return resourceGoogleProjectServiceRead(d, meta)
}

func resourceGoogleProjectServiceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	id, err := parseProjectServiceId(d.Id())
	if err != nil {
		return err
	}

	services, err := getApiServices(project, config, map[string]struct{}{})
	if err != nil {
		return err
	}

	for _, s := range services {
		if s == id.service {
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

	id, err := parseProjectServiceId(d.Id())
	if err != nil {
		return err
	}

	if err = disableService(id.service, project, config); err != nil {
		return fmt.Errorf("Error disabling service: %s", err)
	}

	d.SetId("")
	return nil
}

// Parts that make up the id of a `google_project_service` resource.
// Project is included in order to allow multiple projects to enable the same service within the same Terraform state
type projectServiceId struct {
	project string
	service string
}

func (id projectServiceId) terraformId() string {
	return fmt.Sprintf("%s/%s", id.project, id.service)
}

func parseProjectServiceId(id string) (*projectServiceId, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid google_project_service id format, expecting `{project}/{service}`, found %s", id)
	}

	return &projectServiceId{parts[0], parts[1]}, nil
}
