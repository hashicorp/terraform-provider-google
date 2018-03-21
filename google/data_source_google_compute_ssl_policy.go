package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGoogleComputeSslPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeSslPolicyRead,

		Schema: map[string]*schema.Schema{

			"custom_features": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"fingerprint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"min_tls_version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"profile": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleComputeSslPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	sslPolicy, err := config.clientComputeBeta.SslPolicies.Get(project, name).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("SSL Policy %q", name))
	}

	d.Set("name", sslPolicy.Name)
	d.Set("description", sslPolicy.Description)
	d.Set("min_tls_version", sslPolicy.MinTlsVersion)
	d.Set("profile", sslPolicy.Profile)
	d.Set("fingerprint", sslPolicy.Fingerprint)
	d.Set("project", project)
	d.Set("self_link", ConvertSelfLinkToV1(sslPolicy.SelfLink))

	if sslPolicy.CustomFeatures != nil {
		d.Set("custom_features", convertStringArrToInterface(sslPolicy.CustomFeatures))
	}
	d.SetId(sslPolicy.Name)
	return nil
}
