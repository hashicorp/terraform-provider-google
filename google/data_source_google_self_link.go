package google

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGoogleSelfLink() *schema.Resource {
	return &schema.Resource{
		Read: datasourceGoogleSelfLink,
		Schema: map[string]*schema.Schema{
			"self_link": {
				Type:     schema.TypeString,
				Required: true,
			},
			"relative_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceGoogleSelfLink(d *schema.ResourceData, meta interface{}) error {
	selfLink := d.Get("self_link").(string)

	relativeUri, err := getRelativePath(selfLink)
	if err != nil {
		return err
	}

	d.SetId(selfLink)
	d.Set("self_link", selfLink)
	d.Set("relative_uri", relativeUri)
	d.Set("name", GetResourceNameFromSelfLink(selfLink))

	return nil
}
