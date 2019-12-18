package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func dataSourceGoogleComputeRegionInstanceGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceComputeRegionInstanceGroupRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance": {
							Type:     schema.TypeString,
							Required: true,
						},

						"status": {
							Type:     schema.TypeString,
							Required: true,
						},

						"named_ports": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
					},
				},
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"self_link": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceComputeRegionInstanceGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, region, name, err := GetRegionalResourcePropertiesFromSelfLinkOrSchema(d, config)
	if err != nil {
		return err
	}

	instanceGroup, err := config.clientCompute.RegionInstanceGroups.Get(
		project, region, name).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Region Instance Group %q", name))
	}

	members, err := config.clientCompute.RegionInstanceGroups.ListInstances(
		project, region, name, &compute.RegionInstanceGroupsListInstancesRequest{
			InstanceState: "ALL",
		}).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			// The resource doesn't have any instances, which is okay.
			d.Set("instances", nil)
		} else {
			return fmt.Errorf("Error reading RegionInstanceGroup Members: %s", err)
		}
	} else {
		d.Set("instances", flattenInstancesWithNamedPorts(members.Items))
	}
	d.SetId(fmt.Sprintf("projects/%s/regions/%s/instanceGroups/%s", project, region, name))
	d.Set("self_link", instanceGroup.SelfLink)
	d.Set("name", name)
	d.Set("project", project)
	d.Set("region", region)
	return nil
}

func flattenInstancesWithNamedPorts(insts []*compute.InstanceWithNamedPorts) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(insts))
	log.Printf("There were %d instances.\n", len(insts))
	for _, inst := range insts {
		instMap := make(map[string]interface{})
		instMap["instance"] = inst.Instance
		instMap["named_ports"] = flattenNamedPorts(inst.NamedPorts)
		instMap["status"] = inst.Status
		result = append(result, instMap)
	}
	return result
}

func flattenNamedPorts(namedPorts []*compute.NamedPort) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(namedPorts))
	for _, namedPort := range namedPorts {
		namedPortMap := make(map[string]interface{})
		namedPortMap["name"] = namedPort.Name
		namedPortMap["port"] = namedPort.Port
		result = append(result, namedPortMap)
	}
	return result
}
