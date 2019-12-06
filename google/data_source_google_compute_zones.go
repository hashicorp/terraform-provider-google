package google

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"google.golang.org/api/compute/v1"
)

func dataSourceGoogleComputeZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeZonesRead,
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
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"UP", "DOWN"}, false),
			},
		},
	}
}

func dataSourceGoogleComputeZonesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	region := config.Region
	if r, ok := d.GetOk("region"); ok {
		region = r.(string)
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// we want to share exactly the same base path as the compute client or the
	// region string may mismatch, giving us no results
	// note that the client's BasePath includes a `projects/` suffix, so that'll
	// need to be added to the URL below if the source changes
	computeClientBasePath := config.clientCompute.BasePath

	regionUrl, err := replaceVars(d, config, fmt.Sprintf("%s%s/regions/%s", computeClientBasePath, project, region))
	if err != nil {
		return err
	}
	filter := fmt.Sprintf("(region eq %s)", regionUrl)

	if s, ok := d.GetOk("status"); ok {
		filter += fmt.Sprintf(" (status eq %s)", s)
	}

	call := config.clientCompute.Zones.List(project).Filter(filter)

	resp, err := call.Do()
	if err != nil {
		return err
	}

	zones := flattenZones(resp.Items)
	log.Printf("[DEBUG] Received Google Compute Zones: %q", zones)

	d.Set("names", zones)
	d.Set("region", region)
	d.Set("project", project)
	d.SetId(time.Now().UTC().String())

	return nil
}

func flattenZones(zones []*compute.Zone) []string {
	result := make([]string, len(zones))
	for i, zone := range zones {
		result[i] = zone.Name
	}
	sort.Strings(result)
	return result
}
