package google

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	compute "google.golang.org/api/compute/v1"
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
	var project, region, name string
	if self_link, ok := d.GetOk("self_link"); ok {
		parsed, err := url.Parse(self_link.(string))
		if err != nil {
			return err
		}
		s := strings.Split(parsed.Path, "/")
		project, region, name = s[4], s[6], s[8]
		// e.g. https://www.googleapis.com/compute/beta/projects/project_name/regions/region_name/instanceGroups/foobarbaz

	} else {
		var err error
		project, err = getProject(d, config)
		if err != nil {
			return err
		}

		region, err = getRegion(d, config)
		if err != nil {
			return err
		}
		n, ok := d.GetOk("name")
		name = n.(string)
		if !ok {
			return errors.New("Must provide either `self_link` or `name`.")
		}
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
	d.SetId(strconv.FormatUint(instanceGroup.Id, 16))
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
