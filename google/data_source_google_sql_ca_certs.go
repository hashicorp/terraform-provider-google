package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleSQLCaCerts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleSQLCaCertsRead,

		Schema: map[string]*schema.Schema{
			"instance": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},
			"project": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"active_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certs": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cert": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"common_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expiration_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sha1_fingerprint": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleSQLCaCertsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	fv, err := parseProjectFieldValue("instances", d.Get("instance").(string), "project", d, config, false)
	if err != nil {
		return err
	}
	project := fv.Project
	instance := fv.Name

	log.Printf("[DEBUG] Fetching CA certs from instance %s", instance)

	response, err := config.NewSqlAdminClient(userAgent).Instances.ListServerCas(project, instance).Do()
	if err != nil {
		return fmt.Errorf("error retrieving CA certs: %s", err)
	}

	log.Printf("[DEBUG] Fetched CA certs from instance %s", instance)

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("certs", flattenServerCaCerts(response.Certs)); err != nil {
		return fmt.Errorf("Error setting certs: %s", err)
	}
	if err := d.Set("active_version", response.ActiveVersion); err != nil {
		return fmt.Errorf("Error setting active_version: %s", err)
	}
	d.SetId(fmt.Sprintf("projects/%s/instance/%s", project, instance))

	return nil
}
