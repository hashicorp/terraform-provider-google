// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/compute/v1"
)

func DataSourceGoogleComputeMachineTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGoogleComputeMachineTypesRead,

		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"machine_types": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of machine types`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the machine type.`,
						},
						"guest_cpus": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of virtual CPUs that are available to the instance.`,
						},
						"memory_mb": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The amount of physical memory available to the instance, defined in MB.`,
						},
						"deprecated": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: `The deprecation status associated with this machine type. Only applicable if the machine type is unavailable.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"replacement": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The URL of the suggested replacement for a deprecated machine type.`,
									},
									"state": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The deprecation state of this resource. This can be ACTIVE, DEPRECATED, OBSOLETE, or DELETED.`,
									},
								},
							},
						},
						"maximum_persistent_disks": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The maximum persistent disks allowed.`,
						},
						"maximum_persistent_disks_size_gb": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The maximum total persistent disks size (GB) allowed.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `A textual description of the machine type.`,
						},
						"is_shared_cpus": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether this machine type has a shared CPU.`,
						},
						"accelerators": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `A list of accelerator configurations assigned to this machine type.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"guest_accelerator_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The accelerator type resource name, not a full URL, e.g. nvidia-tesla-t4.`,
									},
									"guest_accelerator_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `Number of accelerator cards exposed to the guest.`,
									},
								},
							},
						},
						"self_link": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The server-defined URL for the machine type.`,
						},
					},
				},
			},
			"zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the zone for this request.`,
				Optional:    true,
			},

			"project": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Project ID for this request.`,
				Optional:    true,
			},
		},
	}
}

func dataSourceGoogleComputeMachineTypesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return diag.FromErr(err)
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return diag.FromErr(err)
	}

	filter := d.Get("filter").(string)
	zone := d.Get("zone").(string)

	machineTypes := make([]map[string]interface{}, 0)
	token := ""

	for paginate := true; paginate; {
		resp, err := config.NewComputeClient(userAgent).MachineTypes.List(project, zone).Context(ctx).Filter(filter).PageToken(token).Do()
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error retrieving machine types: %w", err))

		}
		pageMachineTypes := flattenDatasourceGoogleComputeMachineTypesList(ctx, resp.Items)
		machineTypes = append(machineTypes, pageMachineTypes...)

		token = resp.NextPageToken
		paginate = token != ""
	}

	if err := d.Set("machine_types", machineTypes); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting machine_types: %w", err))
	}

	if err := d.Set("project", project); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting project: %w", err))
	}
	if err := d.Set("zone", zone); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting zone: %w", err))
	}

	id := fmt.Sprintf("projects/%s/zones/%s/machineTypes/filters/%s", project, zone, filter)
	d.SetId(id)

	return diag.Diagnostics{}
}

func flattenDatasourceGoogleComputeMachineTypesList(ctx context.Context, v []*compute.MachineType) []map[string]interface{} {
	if v == nil {
		return make([]map[string]interface{}, 0)
	}

	machineTypes := make([]map[string]interface{}, 0, len(v))
	for _, mt := range v {
		accelerators := make([]map[string]interface{}, len(mt.Accelerators))
		for i, a := range mt.Accelerators {
			accelerators[i] = map[string]interface{}{
				"guest_accelerator_type":  a.GuestAcceleratorType,
				"guest_accelerator_count": a.GuestAcceleratorCount,
			}
		}
		machineType := map[string]interface{}{
			"name":                             mt.Name,
			"guest_cpus":                       mt.GuestCpus,
			"memory_mb":                        mt.MemoryMb,
			"maximum_persistent_disks":         mt.MaximumPersistentDisks,
			"maximum_persistent_disks_size_gb": mt.MaximumPersistentDisksSizeGb,
			"description":                      mt.Description,
			"is_shared_cpus":                   mt.IsSharedCpu,
			"accelerators":                     accelerators,
			"self_link":                        mt.SelfLink,
		}
		if dep := mt.Deprecated; dep != nil {
			d := map[string]interface{}{
				"replacement": dep.Replacement,
				"state":       dep.State,
			}
			machineType["deprecated"] = []map[string]interface{}{d}
		}
		machineTypes = append(machineTypes, machineType)
	}

	return machineTypes
}
