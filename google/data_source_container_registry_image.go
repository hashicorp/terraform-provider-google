package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceGoogleContainerImage() *schema.Resource {
	return &schema.Resource{
		Read: containerRegistryImageRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"digest": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"image_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func containerRegistryImageRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	d.Set("project", project)
	region, ok := d.GetOk("region")
	var url_base string
	escapedProject := strings.Replace(project, ":", "/", -1)
	if ok && region != nil && region != "" {
		url_base = fmt.Sprintf("%s.gcr.io/%s", region, escapedProject)
	} else {
		url_base = fmt.Sprintf("gcr.io/%s", escapedProject)
	}
	tag, t_ok := d.GetOk("tag")
	digest, d_ok := d.GetOk("digest")
	if t_ok && tag != nil && tag != "" {
		d.Set("image_url", fmt.Sprintf("%s/%s:%s", url_base, d.Get("name").(string), tag))
	} else if d_ok && digest != nil && digest != "" {
		d.Set("image_url", fmt.Sprintf("%s/%s@%s", url_base, d.Get("name").(string), digest))
	} else {
		d.Set("image_url", fmt.Sprintf("%s/%s", url_base, d.Get("name").(string)))
	}
	d.SetId(d.Get("image_url").(string))
	return nil
}
