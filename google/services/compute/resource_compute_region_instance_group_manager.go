// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"google.golang.org/api/compute/v1"
)

func ResourceComputeRegionInstanceGroupManager() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeRegionInstanceGroupManagerCreate,
		Read:   resourceComputeRegionInstanceGroupManagerRead,
		Update: resourceComputeRegionInstanceGroupManagerUpdate,
		Delete: resourceComputeRegionInstanceGroupManagerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRegionInstanceGroupManagerStateImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"base_instance_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The base instance name to use for instances in this group. The value must be a valid RFC1035 name. Supported characters are lowercase letters, numbers, and hyphens (-). Instances are named by appending a hyphen and a random four-character string to the base instance name.`,
			},

			"version": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `Application versions managed by this instance group. Each version deals with a specific instance template, allowing canary release scenarios.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Version name.`,
						},

						"instance_template": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: compareSelfLinkRelativePathsIgnoreParams,
							Description:      `The full URL to an instance template from which all new instances of this version will be created.`,
						},

						"target_size": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: `The number of instances calculated as a fixed number or a percentage depending on the settings.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"fixed": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: `The number of instances which are managed for this version. Conflicts with percent.`,
									},

									"percent": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(0, 100),
										Description:  `The number of instances (calculated as percentage) which are managed for this version. Conflicts with fixed. Note that when using percent, rounding will be in favor of explicitly set target_size values; a managed instance group with 2 instances and 2 versions, one of which has a target_size.percent of 60 will create 2 instances of that version.`,
									},
								},
							},
						},
					},
				},
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the instance group manager. Must be 1-63 characters long and comply with RFC1035. Supported characters include lowercase letters, numbers, and hyphens.`,
			},

			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: `The region where the managed instance group resides.`,
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `An optional textual description of the instance group manager.`,
			},

			"fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The fingerprint of the instance group manager.`,
			},

			"instance_group": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The full URL of the instance group created by the manager.`,
			},

			"named_port": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: `The named port configuration.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of the port.`,
						},

						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `The port number.`,
						},
					},
				},
			},

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The ID of the project in which the resource belongs. If it is not provided, the provider project is used.`,
			},

			"self_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The URL of the created resource.`,
			},

			"target_pools": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         tpgresource.SelfLinkRelativePathHash,
				Description: `The full URL of all target pools to which new instances in the group are added. Updating the target pools attribute does not affect existing instances.`,
			},
			"target_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: `The target number of running instances for this managed instance group. This value should always be explicitly set unless this resource is attached to an autoscaler, in which case it should never be set. Defaults to 0.`,
			},

			"list_managed_instances_results": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "PAGELESS",
				ValidateFunc: validation.StringInSlice([]string{"PAGELESS", "PAGINATED"}, false),
				Description:  `Pagination behavior of the listManagedInstances API method for this managed instance group. Valid values are: "PAGELESS", "PAGINATED". If PAGELESS (default), Pagination is disabled for the group's listManagedInstances API method. maxResults and pageToken query parameters are ignored and all instances are returned in a single response. If PAGINATED, pagination is enabled, maxResults and pageToken query parameters are respected.`,
			},

			// If true, the resource will report ready only after no instances are being created.
			// This will not block future reads if instances are being recreated, and it respects
			// the "createNoRetry" parameter that's available for this resource.
			"wait_for_instances": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: `Whether to wait for all instances to be created/updated before returning. Note that if this is set to true and the operation does not succeed, Terraform will continue trying until it times out.`,
			},
			"wait_for_instances_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "STABLE",
				ValidateFunc: validation.StringInSlice([]string{"STABLE", "UPDATED"}, false),

				Description: `When used with wait_for_instances specifies the status to wait for. When STABLE is specified this resource will wait until the instances are stable before returning. When UPDATED is set, it will wait for the version target to be reached and any per instance configs to be effective as well as all instances to be stable before returning.`,
			},

			"auto_healing_policies": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: `The autohealing policies for this managed instance group. You can specify only one value.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health_check": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: tpgresource.CompareSelfLinkRelativePaths,
							Description:      `The health check resource that signals autohealing.`,
						},

						"initial_delay_sec": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 3600),
							Description:  `The number of seconds that the managed instance group waits before it applies autohealing policies to new instances or recently recreated instances. Between 0 and 3600.`,
						},
					},
				},
			},

			"distribution_policy_zones": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The distribution policy for this managed instance group. You can specify one or more values.`,
				Set:         hashZoneFromSelfLinkOrResourceName,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				},
			},

			"distribution_policy_target_shape": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The shape to which the group converges either proactively or on resize events (depending on the value set in updatePolicy.instanceRedistributionType).`,
			},

			"instance_lifecycle_policy": {
				Computed:    true,
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: `The instance lifecycle policy for this managed instance group.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"force_update_on_repair": {
							Type:         schema.TypeString,
							Default:      "NO",
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"YES", "NO"}, false),
							Description:  `Specifies whether to apply the group's latest configuration when repairing a VM. Valid options are: YES, NO. If YES and you updated the group's instance template or per-instance configurations after the VM was created, then these changes are applied when VM is repaired. If NO (default), then updates are applied in accordance with the group's update policy type.`,
						},
					},
				},
			},

			"update_policy": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				MaxItems:    1,
				Description: `The update policy for this managed instance group.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"minimal_action": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"REFRESH", "RESTART", "REPLACE"}, false),
							Description:  `Minimal action to be taken on an instance. You can specify either REFRESH to update without stopping instances, RESTART to restart existing instances or REPLACE to delete and create new instances from the target template. If you specify a REFRESH, the Updater will attempt to perform that action only. However, if the Updater determines that the minimal action you specify is not enough to perform the update, it might perform a more disruptive action.`,
						},

						"most_disruptive_allowed_action": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"NONE", "REFRESH", "RESTART", "REPLACE"}, false),
							Description:  `Most disruptive action that is allowed to be taken on an instance. You can specify either NONE to forbid any actions, REFRESH to allow actions that do not need instance restart, RESTART to allow actions that can be applied without instance replacing or REPLACE to allow all possible actions. If the Updater determines that the minimal update action needed is more disruptive than most disruptive allowed action you specify it will not perform the update at all.`,
						},

						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"OPPORTUNISTIC", "PROACTIVE"}, false),
							Description:  `The type of update process. You can specify either PROACTIVE so that the instance group manager proactively executes actions in order to bring instances to their target versions or OPPORTUNISTIC so that no action is proactively executed but the update will be performed as part of other actions (for example, resizes or recreateInstances calls).`,
						},

						"max_surge_fixed": {
							Type:          schema.TypeInt,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"update_policy.0.max_surge_percent"},
							Description:   `The maximum number of instances that can be created above the specified targetSize during the update process. Conflicts with max_surge_percent. It has to be either 0 or at least equal to the number of zones. If fixed values are used, at least one of max_unavailable_fixed or max_surge_fixed must be greater than 0.`,
						},

						"max_surge_percent": {
							Type:          schema.TypeInt,
							Optional:      true,
							ConflictsWith: []string{"update_policy.0.max_surge_fixed"},
							Description:   `The maximum number of instances(calculated as percentage) that can be created above the specified targetSize during the update process. Conflicts with max_surge_fixed. Percent value is only allowed for regional managed instance groups with size at least 10.`,
							ValidateFunc:  validation.IntBetween(0, 100),
						},

						"max_unavailable_fixed": {
							Type:          schema.TypeInt,
							Optional:      true,
							Computed:      true,
							Description:   `The maximum number of instances that can be unavailable during the update process. Conflicts with max_unavailable_percent. It has to be either 0 or at least equal to the number of zones. If fixed values are used, at least one of max_unavailable_fixed or max_surge_fixed must be greater than 0.`,
							ConflictsWith: []string{"update_policy.0.max_unavailable_percent"},
						},

						"max_unavailable_percent": {
							Type:          schema.TypeInt,
							Optional:      true,
							ConflictsWith: []string{"update_policy.0.max_unavailable_fixed"},
							ValidateFunc:  validation.IntBetween(0, 100),
							Description:   `The maximum number of instances(calculated as percentage) that can be unavailable during the update process. Conflicts with max_unavailable_fixed. Percent value is only allowed for regional managed instance groups with size at least 10.`,
						},

						"instance_redistribution_type": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateFunc:     validation.StringInSlice([]string{"PROACTIVE", "NONE", ""}, false),
							DiffSuppressFunc: tpgresource.EmptyOrDefaultStringSuppress("PROACTIVE"),
							Description:      `The instance redistribution policy for regional managed instance groups. Valid values are: "PROACTIVE", "NONE". If PROACTIVE (default), the group attempts to maintain an even distribution of VM instances across zones in the region. If NONE, proactive redistribution is disabled.`,
						},
						"replacement_method": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateFunc:     validation.StringInSlice([]string{"RECREATE", "SUBSTITUTE", ""}, false),
							DiffSuppressFunc: tpgresource.EmptyOrDefaultStringSuppress("SUBSTITUTE"),
							Description:      `The instance replacement method for regional managed instance groups. Valid values are: "RECREATE", "SUBSTITUTE". If SUBSTITUTE (default), the group replaces VM instances with new instances that have randomly generated names. If RECREATE, instance names are preserved.  You must also set max_unavailable_fixed or max_unavailable_percent to be greater than 0.`,
						},
					},
				},
			},
			"stateful_disk": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: `Disks created on the instances that will be preserved on instance delete, update, etc. Structure is documented below. For more information see the official documentation. Proactive cross zone instance redistribution must be disabled before you can update stateful disks on existing instance group managers. This can be controlled via the update_policy.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"device_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The device name of the disk to be attached.`,
						},

						"delete_rule": {
							Type:         schema.TypeString,
							Default:      "NEVER",
							Optional:     true,
							Description:  `A value that prescribes what should happen to the stateful disk when the VM instance is deleted. The available options are NEVER and ON_PERMANENT_INSTANCE_DELETION. NEVER - detach the disk when the VM is deleted, but do not delete the disk. ON_PERMANENT_INSTANCE_DELETION will delete the stateful disk when the VM is permanently deleted from the instance group. The default is NEVER.`,
							ValidateFunc: validation.StringInSlice([]string{"NEVER", "ON_PERMANENT_INSTANCE_DELETION"}, true),
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The status of this managed instance group.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_stable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `A bit indicating whether the managed instance group is in a stable state. A stable state means that: none of the instances in the managed instance group is currently undergoing any type of change (for example, creation, restart, or deletion); no future changes are scheduled for instances in the managed instance group; and the managed instance group itself is not being modified.`,
						},

						"version_target": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `A status of consistency of Instances' versions with their target version specified by version field on Instance Group Manager.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_reached": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `A bit indicating whether version target has been reached in this managed instance group, i.e. all instances are in their target version. Instances' target version are specified by version field on Instance Group Manager.`,
									},
								},
							},
						},
						"stateful": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Stateful status of the given Instance Group Manager.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"has_stateful_config": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `A bit indicating whether the managed instance group has stateful configuration, that is, if you have configured any items in a stateful policy or in per-instance configs. The group might report that it has no stateful config even when there is still some preserved state on a managed instance, for example, if you have deleted all PICs but not yet applied those deletions.`,
									},
									"per_instance_configs": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `Status of per-instance configs on the instance.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"all_effective": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: `A bit indicating if all of the group's per-instance configs (listed in the output of a listPerInstanceConfigs API call) have status EFFECTIVE or there are no per-instance-configs.`,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		UseJSONNumber: true,
	}
}

func resourceComputeRegionInstanceGroupManagerCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return err
	}

	manager := &compute.InstanceGroupManager{
		Name:                        d.Get("name").(string),
		Description:                 d.Get("description").(string),
		BaseInstanceName:            d.Get("base_instance_name").(string),
		TargetSize:                  int64(d.Get("target_size").(int)),
		ListManagedInstancesResults: d.Get("list_managed_instances_results").(string),
		NamedPorts:                  getNamedPortsBeta(d.Get("named_port").(*schema.Set).List()),
		TargetPools:                 tpgresource.ConvertStringSet(d.Get("target_pools").(*schema.Set)),
		AutoHealingPolicies:         expandAutoHealingPolicies(d.Get("auto_healing_policies").([]interface{})),
		Versions:                    expandVersions(d.Get("version").([]interface{})),
		UpdatePolicy:                expandRegionUpdatePolicy(d.Get("update_policy").([]interface{})),
		InstanceLifecyclePolicy:     expandInstanceLifecyclePolicy(d.Get("instance_lifecycle_policy").([]interface{})),
		DistributionPolicy:          expandDistributionPolicy(d),
		StatefulPolicy:              expandStatefulPolicy(d),
		// Force send TargetSize to allow size of 0.
		ForceSendFields: []string{"TargetSize"},
	}

	op, err := config.NewComputeClient(userAgent).RegionInstanceGroupManagers.Insert(project, region, manager).Do()

	if err != nil {
		return fmt.Errorf("Error creating RegionInstanceGroupManager: %s", err)
	}

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/regions/{{region}}/instanceGroupManagers/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	// Wait for the operation to complete
	err = ComputeOperationWaitTime(config, op, project, "Creating InstanceGroupManager", userAgent, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	if d.Get("wait_for_instances").(bool) {
		err := computeRIGMWaitForInstanceStatus(d, meta)
		if err != nil {
			return err
		}
	}

	return resourceComputeRegionInstanceGroupManagerRead(d, config)
}

func computeRIGMWaitForInstanceStatus(d *schema.ResourceData, meta interface{}) error {
	waitForUpdates := d.Get("wait_for_instances_status").(string) == "UPDATED"
	conf := resource.StateChangeConf{
		Pending: []string{"creating", "error", "updating per instance configs", "reaching version target"},
		Target:  []string{"created"},
		Refresh: waitForInstancesRefreshFunc(getRegionalManager, waitForUpdates, d, meta),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}
	_, err := conf.WaitForState()
	if err != nil {
		return err
	}
	return nil
}

type getInstanceManagerFunc func(*schema.ResourceData, interface{}) (*compute.InstanceGroupManager, error)

func getRegionalManager(d *schema.ResourceData, meta interface{}) (*compute.InstanceGroupManager, error) {
	config := meta.(*transport_tpg.Config)

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return nil, err
	}

	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return nil, err
	}

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return nil, err
	}

	name := d.Get("name").(string)
	manager, err := config.NewComputeClient(userAgent).RegionInstanceGroupManagers.Get(project, region, name).Do()
	if err != nil {
		return nil, transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Region Instance Manager %q", name))
	}

	return manager, nil
}

func waitForInstancesRefreshFunc(f getInstanceManagerFunc, waitForUpdates bool, d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		m, err := f(d, meta)
		if err != nil {
			log.Printf("[WARNING] Error in fetching manager while waiting for instances to come up: %s\n", err)
			return nil, "error", err
		}
		if m == nil {
			// getManager/getRegional manager call handleNotFoundError, which will return a nil error and nil object in the case
			// that the original error was a 404. if m == nil here, we will assume that it was not found return an "instance manager not found"
			// error so that we can parse it later on and handle it there
			return nil, "error", fmt.Errorf("instance manager not found")
		}
		if m.Status.IsStable {
			if waitForUpdates {
				// waitForUpdates waits for versions to be reached and per instance configs to be updated (if present)
				if m.Status.Stateful.HasStatefulConfig {
					if !m.Status.Stateful.PerInstanceConfigs.AllEffective {
						return false, "updating per instance configs", nil
					}
				}
				if !m.Status.VersionTarget.IsReached {
					return false, "reaching version target", nil
				}
				if !m.Status.VersionTarget.IsReached {
					return false, "reaching version target", nil
				}
			}
			return true, "created", nil
		} else {
			return false, "creating", nil
		}
	}
}

func resourceComputeRegionInstanceGroupManagerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	manager, err := getRegionalManager(d, meta)
	if err != nil {
		return err
	}
	if manager == nil {
		log.Printf("[WARN] Region Instance Group Manager %q not found, removing from state.", d.Id())
		d.SetId("")
		return nil
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	if err := d.Set("base_instance_name", manager.BaseInstanceName); err != nil {
		return fmt.Errorf("Error setting base_instance_name: %s", err)
	}
	if err := d.Set("name", manager.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("region", tpgresource.GetResourceNameFromSelfLink(manager.Region)); err != nil {
		return fmt.Errorf("Error setting region: %s", err)
	}
	if err := d.Set("description", manager.Description); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("target_size", manager.TargetSize); err != nil {
		return fmt.Errorf("Error setting target_size: %s", err)
	}
	if err := d.Set("list_managed_instances_results", manager.ListManagedInstancesResults); err != nil {
		return fmt.Errorf("Error setting list_managed_instances_results: %s", err)
	}
	if err := d.Set("target_pools", tpgresource.MapStringArr(manager.TargetPools, tpgresource.ConvertSelfLinkToV1)); err != nil {
		return fmt.Errorf("Error setting target_pools in state: %s", err.Error())
	}
	if err := d.Set("named_port", flattenNamedPortsBeta(manager.NamedPorts)); err != nil {
		return fmt.Errorf("Error setting named_port in state: %s", err.Error())
	}
	if err := d.Set("fingerprint", manager.Fingerprint); err != nil {
		return fmt.Errorf("Error setting fingerprint: %s", err)
	}
	if err := d.Set("instance_group", tpgresource.ConvertSelfLinkToV1(manager.InstanceGroup)); err != nil {
		return fmt.Errorf("Error setting instance_group: %s", err)
	}
	if err := d.Set("distribution_policy_zones", flattenDistributionPolicy(manager.DistributionPolicy)); err != nil {
		return err
	}
	if err := d.Set("distribution_policy_target_shape", manager.DistributionPolicy.TargetShape); err != nil {
		return err
	}
	if err := d.Set("self_link", tpgresource.ConvertSelfLinkToV1(manager.SelfLink)); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}

	if err := d.Set("auto_healing_policies", flattenAutoHealingPolicies(manager.AutoHealingPolicies)); err != nil {
		return fmt.Errorf("Error setting auto_healing_policies in state: %s", err.Error())
	}
	if err := d.Set("version", flattenVersions(manager.Versions)); err != nil {
		return err
	}
	if err := d.Set("update_policy", flattenRegionUpdatePolicy(manager.UpdatePolicy)); err != nil {
		return fmt.Errorf("Error setting update_policy in state: %s", err.Error())
	}
	if err = d.Set("instance_lifecycle_policy", flattenInstanceLifecyclePolicy(manager.InstanceLifecyclePolicy)); err != nil {
		return fmt.Errorf("Error setting instance lifecycle policy in state: %s", err.Error())
	}
	if err = d.Set("stateful_disk", flattenStatefulPolicy(manager.StatefulPolicy)); err != nil {
		return fmt.Errorf("Error setting stateful_disk in state: %s", err.Error())
	}
	if err = d.Set("status", flattenStatus(manager.Status)); err != nil {
		return fmt.Errorf("Error setting status in state: %s", err.Error())
	}
	// If unset in state set to default value
	if d.Get("wait_for_instances_status").(string) == "" {
		if err = d.Set("wait_for_instances_status", "STABLE"); err != nil {
			return fmt.Errorf("Error setting wait_for_instances_status in state: %s", err.Error())
		}
	}

	return nil
}

func resourceComputeRegionInstanceGroupManagerUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return err
	}

	updatedManager := &compute.InstanceGroupManager{
		Fingerprint: d.Get("fingerprint").(string),
	}
	var change bool

	if d.HasChange("target_pools") {
		updatedManager.TargetPools = tpgresource.ConvertStringSet(d.Get("target_pools").(*schema.Set))
		updatedManager.ForceSendFields = append(updatedManager.ForceSendFields, "TargetPools")
		change = true
	}

	if d.HasChange("auto_healing_policies") {
		updatedManager.AutoHealingPolicies = expandAutoHealingPolicies(d.Get("auto_healing_policies").([]interface{}))
		updatedManager.ForceSendFields = append(updatedManager.ForceSendFields, "AutoHealingPolicies")
		change = true
	}

	if d.HasChange("version") {
		updatedManager.Versions = expandVersions(d.Get("version").([]interface{}))
		change = true
	}

	if d.HasChange("update_policy") {
		updatedManager.UpdatePolicy = expandRegionUpdatePolicy(d.Get("update_policy").([]interface{}))
		change = true
	}

	if d.HasChange("instance_lifecycle_policy") {
		updatedManager.InstanceLifecyclePolicy = expandInstanceLifecyclePolicy(d.Get("instance_lifecycle_policy").([]interface{}))
		change = true
	}

	if d.HasChange("stateful_disk") {
		updatedManager.StatefulPolicy = expandStatefulPolicy(d)
		change = true
	}

	if d.HasChange("list_managed_instances_results") {
		updatedManager.ListManagedInstancesResults = d.Get("list_managed_instances_results").(string)
		change = true
	}

	if change {
		op, err := config.NewComputeClient(userAgent).RegionInstanceGroupManagers.Patch(project, region, d.Get("name").(string), updatedManager).Do()
		if err != nil {
			return fmt.Errorf("Error updating region managed group instances: %s", err)
		}

		err = ComputeOperationWaitTime(config, op, project, "Updating region managed group instances", userAgent, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	}

	// named ports can't be updated through PATCH
	// so we call the update method on the region instance group, instead of the rigm
	if d.HasChange("named_port") {
		d.Partial(true)
		namedPorts := getNamedPortsBeta(d.Get("named_port").(*schema.Set).List())
		setNamedPorts := &compute.RegionInstanceGroupsSetNamedPortsRequest{
			NamedPorts: namedPorts,
		}

		op, err := config.NewComputeClient(userAgent).RegionInstanceGroups.SetNamedPorts(
			project, region, d.Get("name").(string), setNamedPorts).Do()

		if err != nil {
			return fmt.Errorf("Error updating RegionInstanceGroupManager: %s", err)
		}

		err = ComputeOperationWaitTime(config, op, project, "Updating RegionInstanceGroupManager", userAgent, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	}

	// target size should use resize
	if d.HasChange("target_size") {
		d.Partial(true)
		targetSize := int64(d.Get("target_size").(int))
		op, err := config.NewComputeClient(userAgent).RegionInstanceGroupManagers.Resize(
			project, region, d.Get("name").(string), targetSize).Do()

		if err != nil {
			return fmt.Errorf("Error resizing RegionInstanceGroupManager: %s", err)
		}

		err = ComputeOperationWaitTime(config, op, project, "Resizing RegionInstanceGroupManager", userAgent, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	}

	d.Partial(false)

	if d.Get("wait_for_instances").(bool) {
		err := computeRIGMWaitForInstanceStatus(d, meta)
		if err != nil {
			return err
		}
	}

	return resourceComputeRegionInstanceGroupManagerRead(d, meta)
}

func resourceComputeRegionInstanceGroupManagerDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	if d.Get("wait_for_instances").(bool) {
		err := computeRIGMWaitForInstanceStatus(d, meta)
		if err != nil {
			notFound, reErr := regexp.MatchString(`not found`, err.Error())
			if reErr != nil {
				return reErr
			}
			if notFound {
				// manager was not found, we can exit gracefully
				return nil
			}
			return err
		}
	}

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	op, err := config.NewComputeClient(userAgent).RegionInstanceGroupManagers.Delete(project, region, name).Do()

	if err != nil {
		return fmt.Errorf("Error deleting region instance group manager: %s", err)
	}

	// Wait for the operation to complete
	err = ComputeOperationWaitTime(config, op, project, "Deleting RegionInstanceGroupManager", userAgent, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf("Error waiting for delete to complete: %s", err)
	}

	d.SetId("")
	return nil
}

func expandRegionUpdatePolicy(configured []interface{}) *compute.InstanceGroupManagerUpdatePolicy {
	updatePolicy := &compute.InstanceGroupManagerUpdatePolicy{}

	for _, raw := range configured {
		data := raw.(map[string]interface{})

		updatePolicy.MinimalAction = data["minimal_action"].(string)
		mostDisruptiveAllowedAction := data["most_disruptive_allowed_action"].(string)
		if mostDisruptiveAllowedAction != "" {
			updatePolicy.MostDisruptiveAllowedAction = mostDisruptiveAllowedAction
		} else {
			updatePolicy.NullFields = append(updatePolicy.NullFields, "MostDisruptiveAllowedAction")
		}
		updatePolicy.Type = data["type"].(string)
		updatePolicy.InstanceRedistributionType = data["instance_redistribution_type"].(string)
		updatePolicy.ReplacementMethod = data["replacement_method"].(string)

		// percent and fixed values are conflicting
		// when the percent values are set, the fixed values will be ignored
		if v := data["max_surge_percent"]; v.(int) > 0 {
			updatePolicy.MaxSurge = &compute.FixedOrPercent{
				Percent:    int64(v.(int)),
				NullFields: []string{"Fixed"},
			}
		} else {
			updatePolicy.MaxSurge = &compute.FixedOrPercent{
				Fixed: int64(data["max_surge_fixed"].(int)),
				// allow setting this value to 0
				ForceSendFields: []string{"Fixed"},
				NullFields:      []string{"Percent"},
			}
		}

		if v := data["max_unavailable_percent"]; v.(int) > 0 {
			updatePolicy.MaxUnavailable = &compute.FixedOrPercent{
				Percent:    int64(v.(int)),
				NullFields: []string{"Fixed"},
			}
		} else {
			updatePolicy.MaxUnavailable = &compute.FixedOrPercent{
				Fixed: int64(data["max_unavailable_fixed"].(int)),
				// allow setting this value to 0
				ForceSendFields: []string{"Fixed"},
				NullFields:      []string{"Percent"},
			}
		}
	}
	return updatePolicy
}

func flattenRegionUpdatePolicy(updatePolicy *compute.InstanceGroupManagerUpdatePolicy) []map[string]interface{} {
	results := []map[string]interface{}{}
	if updatePolicy != nil {
		up := map[string]interface{}{}
		if updatePolicy.MaxSurge != nil {
			up["max_surge_fixed"] = updatePolicy.MaxSurge.Fixed
			up["max_surge_percent"] = updatePolicy.MaxSurge.Percent
		} else {
			up["max_surge_fixed"] = 0
			up["max_surge_percent"] = 0
		}
		if updatePolicy.MaxUnavailable != nil {
			up["max_unavailable_fixed"] = updatePolicy.MaxUnavailable.Fixed
			up["max_unavailable_percent"] = updatePolicy.MaxUnavailable.Percent
		} else {
			up["max_unavailable_fixed"] = 0
			up["max_unavailable_percent"] = 0
		}
		up["minimal_action"] = updatePolicy.MinimalAction
		up["most_disruptive_allowed_action"] = updatePolicy.MostDisruptiveAllowedAction
		up["type"] = updatePolicy.Type
		up["instance_redistribution_type"] = updatePolicy.InstanceRedistributionType
		up["replacement_method"] = updatePolicy.ReplacementMethod

		results = append(results, up)
	}
	return results
}

func expandDistributionPolicy(d *schema.ResourceData) *compute.DistributionPolicy {
	dpz := d.Get("distribution_policy_zones").(*schema.Set)
	dpts := d.Get("distribution_policy_target_shape").(string)
	if dpz.Len() == 0 && dpts == "" {
		return nil
	}

	distributionPolicyZoneConfigs := make([]*compute.DistributionPolicyZoneConfiguration, 0, dpz.Len())
	for _, raw := range dpz.List() {
		data := raw.(string)
		distributionPolicyZoneConfig := compute.DistributionPolicyZoneConfiguration{
			Zone: "zones/" + data,
		}

		distributionPolicyZoneConfigs = append(distributionPolicyZoneConfigs, &distributionPolicyZoneConfig)
	}

	return &compute.DistributionPolicy{Zones: distributionPolicyZoneConfigs, TargetShape: dpts}
}

func flattenDistributionPolicy(distributionPolicy *compute.DistributionPolicy) []string {
	zones := make([]string, 0)

	if distributionPolicy != nil {
		for _, zone := range distributionPolicy.Zones {
			zones = append(zones, tpgresource.GetResourceNameFromSelfLink(zone.Zone))
		}
	}

	return zones
}

func hashZoneFromSelfLinkOrResourceName(value interface{}) int {
	parts := strings.Split(value.(string), "/")
	resource := parts[len(parts)-1]

	return tpgresource.Hashcode(resource)
}

func resourceRegionInstanceGroupManagerStateImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := d.Set("wait_for_instances", false); err != nil {
		return nil, fmt.Errorf("Error setting wait_for_instances: %s", err)
	}
	if err := d.Set("wait_for_instances_status", "STABLE"); err != nil {
		return nil, fmt.Errorf("Error setting wait_for_instances_status: %s", err)
	}
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{"projects/(?P<project>[^/]+)/regions/(?P<region>[^/]+)/instanceGroupManagers/(?P<name>[^/]+)", "(?P<project>[^/]+)/(?P<region>[^/]+)/(?P<name>[^/]+)", "(?P<region>[^/]+)/(?P<name>[^/]+)", "(?P<name>[^/]+)"}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/regions/{{region}}/instanceGroupManagers/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
