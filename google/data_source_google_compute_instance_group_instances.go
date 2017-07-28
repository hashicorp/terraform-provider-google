package google

import (
	"log"
	"sort"
	"strings"
	"time"

	compute "google.golang.org/api/compute/v1"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGoogleComputeInstanceGroupInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeInstanceGroupInstancesRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceGoogleComputeInstanceGroupInstancesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	name := d.Get("name").(string)
	zone := d.Get("zone").(string)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	resp, err := config.clientCompute.InstanceGroups.ListInstances(project, zone, name, nil).Do()
	if err != nil {
		return err
	}

	instances := flattenInstanceGroupInstances(resp.Items, false)
	instanceNames := flattenInstanceGroupInstances(resp.Items, true)
	log.Printf("[DEBUG] Received Google Compute Instance Group List Instances: %q", instances)

	d.Set("instances", instances)
	d.Set("names", instanceNames)
	d.SetId(time.Now().UTC().String())

	return nil
}

func flattenInstanceGroupInstances(instances []*compute.InstanceWithNamedPorts, extractName bool) []string {
	result := make([]string, len(instances), len(instances))
	for i, instance := range instances {
		if extractName {
			result[i] = extractInstanceNameFromURL(instance.Instance)
		} else {
			result[i] = instance.Instance
		}
	}
	sort.Strings(result)
	return result
}

func extractInstanceNameFromURL(instanceURL string) string {
	paths := strings.Split(instanceURL, "/")
	return paths[len(paths)-1]
}
