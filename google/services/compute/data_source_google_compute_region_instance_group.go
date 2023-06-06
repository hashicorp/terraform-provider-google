// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func DataSourceGoogleComputeRegionInstanceGroup() *schema.Resource {
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
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, region, name, err := tpgresource.GetRegionalResourcePropertiesFromSelfLinkOrSchema(d, config)
	if err != nil {
		return err
	}

	instanceGroup, err := config.NewComputeClient(userAgent).RegionInstanceGroups.Get(
		project, region, name).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Region Instance Group %q", name))
	}

	members, err := config.NewComputeClient(userAgent).RegionInstanceGroups.ListInstances(
		project, region, name, &compute.RegionInstanceGroupsListInstancesRequest{
			InstanceState: "ALL",
		}).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			// The resource doesn't have any instances, which is okay.
			if err := d.Set("instances", nil); err != nil {
				return fmt.Errorf("Error setting instances: %s", err)
			}
		} else {
			return fmt.Errorf("Error reading RegionInstanceGroup Members: %s", err)
		}
	} else {
		if err := d.Set("instances", flattenInstancesWithNamedPorts(members.Items)); err != nil {
			return fmt.Errorf("Error setting instances: %s", err)
		}
	}
	d.SetId(fmt.Sprintf("projects/%s/regions/%s/instanceGroups/%s", project, region, name))
	if err := d.Set("self_link", instanceGroup.SelfLink); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}
	if err := d.Set("name", name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("region", region); err != nil {
		return fmt.Errorf("Error setting region: %s", err)
	}
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
