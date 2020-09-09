package google

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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

	filter := ""
	if s, ok := d.GetOk("status"); ok {
		filter += fmt.Sprintf(" (status eq %s)", s)
	}

	zones := []string{}
	err = config.clientCompute.Zones.List(project).Filter(filter).Pages(config.context, func(zl *compute.ZoneList) error {
		for _, zone := range zl.Items {
			// We have no way to guarantee a specific base path for the region, but the built-in API-level filtering
			// only lets us query on exact matches, so we do our own filtering here.
			if strings.HasSuffix(zone.Region, "/"+region) {
				zones = append(zones, zone.Name)
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	sort.Strings(zones)
	log.Printf("[DEBUG] Received Google Compute Zones: %q", zones)

	d.Set("names", zones)
	d.Set("region", region)
	d.Set("project", project)
	d.SetId(time.Now().UTC().String())

	return nil
}
