package google

import (
	"fmt"
	"log"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"google.golang.org/api/compute/v1"
)

func DataSourceGoogleComputeNodeTypes() *schema.Resource {
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
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone, err := getZone(d, config)
	if err != nil {
		return fmt.Errorf("Please specify zone to get appropriate node types for zone. Unable to get zone: %s", err)
	}

	resp, err := config.NewComputeClient(userAgent).NodeTypes.List(project, zone).Do()
	if err != nil {
		return err
	}
	nodeTypes := flattenComputeNodeTypes(resp.Items)
	log.Printf("[DEBUG] Received Google Compute Regions: %q", nodeTypes)

	if err := d.Set("names", nodeTypes); err != nil {
		return fmt.Errorf("Error setting names: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("zone", zone); err != nil {
		return fmt.Errorf("Error setting zone: %s", err)
	}
	d.SetId(fmt.Sprintf("projects/%s/zones/%s", project, zone))

	return nil
}

func flattenComputeNodeTypes(nodeTypes []*compute.NodeType) []string {
	result := make([]string, len(nodeTypes))
	for i, nodeType := range nodeTypes {
		result[i] = nodeType.Name
	}
	sort.Strings(result)
	return result
}
