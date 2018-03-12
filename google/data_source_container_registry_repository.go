package google

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGoogleContainerRepo() *schema.Resource {
	return &schema.Resource{
		Read: containerRegistryRepoRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"repository_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func containerRegistryRepoRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	d.Set("project", project)
	region, ok := d.GetOk("region")
	if ok && region != nil && region != "" {
		d.Set("repository_url", fmt.Sprintf("%s.gcr.io/%s", region, project))
	} else {
		d.Set("repository_url", fmt.Sprintf("gcr.io/%s", project))
	}
	d.SetId(d.Get("repository_url").(string))
	return nil
}
