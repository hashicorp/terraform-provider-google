package google

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/compute/v1"
)

func dataSourceGoogleComputeNodeTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeNodeTypesRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceGoogleComputeNodeTypesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone, err := getZone(d, config)
	if err != nil {
		return fmt.Errorf("Please specify zone to get appropriate node types for zone. Unable to get zone: %s", err)
	}

	resp, err := config.clientCompute.NodeTypes.List(project, zone).Do()
	if err != nil {
		return err
	}
	nodeTypes := flattenComputeNodeTypes(resp.Items)
	log.Printf("[DEBUG] Received Google Compute Regions: %q", nodeTypes)

	d.Set("names", nodeTypes)
	d.Set("project", project)
	d.Set("zone", zone)
	d.SetId(time.Now().UTC().String())

	return nil
}

func flattenComputeNodeTypes(nodeTypes []*compute.NodeType) []string {
	result := make([]string, len(nodeTypes), len(nodeTypes))
	for i, nodeType := range nodeTypes {
		result[i] = nodeType.Name
	}
	sort.Strings(result)
	return result
}
