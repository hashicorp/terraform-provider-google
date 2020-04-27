package google

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	computeBeta "google.golang.org/api/compute/v0.beta"
)

func resourceComputeRegionInstanceGroupManager() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeRegionInstanceGroupManagerCreate,
		Read:   resourceComputeRegionInstanceGroupManagerRead,
		Update: resourceComputeRegionInstanceGroupManagerUpdate,
		Delete: resourceComputeRegionInstanceGroupManagerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRegionInstanceGroupManagerStateImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"base_instance_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"instance_template": {
				Type:     schema.TypeString,
				Computed: true,
				Removed:  "This field has been replaced by `version.instance_template` in 3.0.0",
			},

			"version": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"instance_template": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: compareSelfLinkRelativePaths,
						},

						"target_size": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"fixed": {
										Type:     schema.TypeInt,
										Optional: true,
									},

									"percent": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(0, 100),
									},
								},
							},
						},
					},
				},
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"instance_group": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"named_port": {
				Type:     schema.TypeSet,
				Optional: true,
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

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"update_strategy": {
				Type:     schema.TypeString,
				Removed:  "This field is removed.",
				Optional: true,
				Computed: true,
			},

			"target_pools": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: selfLinkRelativePathHash,
			},
			"target_size": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},

			// If true, the resource will report ready only after no instances are being created.
			// This will not block future reads if instances are being recreated, and it respects
			// the "createNoRetry" parameter that's available for this resource.
			"wait_for_instances": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"auto_healing_policies": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health_check": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: compareSelfLinkRelativePaths,
						},

						"initial_delay_sec": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 3600),
						},
					},
				},
			},

			"distribution_policy_zones": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Set:      hashZoneFromSelfLinkOrResourceName,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					DiffSuppressFunc: compareSelfLinkOrResourceName,
				},
			},

			"update_policy": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"minimal_action": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"RESTART", "REPLACE"}, false),
						},

						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"OPPORTUNISTIC", "PROACTIVE"}, false),
						},

						"max_surge_fixed": {
							Type:          schema.TypeInt,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"update_policy.0.max_surge_percent"},
						},

						"max_surge_percent": {
							Type:          schema.TypeInt,
							Optional:      true,
							ConflictsWith: []string{"update_policy.0.max_surge_fixed"},
							ValidateFunc:  validation.IntBetween(0, 100),
						},

						"max_unavailable_fixed": {
							Type:          schema.TypeInt,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"update_policy.0.max_unavailable_percent"},
						},

						"max_unavailable_percent": {
							Type:          schema.TypeInt,
							Optional:      true,
							ConflictsWith: []string{"update_policy.0.max_unavailable_fixed"},
							ValidateFunc:  validation.IntBetween(0, 100),
						},

						"min_ready_sec": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 3600),
						},
						"instance_redistribution_type": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateFunc:     validation.StringInSlice([]string{"PROACTIVE", "NONE", ""}, false),
							DiffSuppressFunc: emptyOrDefaultStringSuppress("PROACTIVE"),
						},
					},
				},
			},
		},
	}
}

func resourceComputeRegionInstanceGroupManagerCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	manager := &computeBeta.InstanceGroupManager{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		BaseInstanceName:    d.Get("base_instance_name").(string),
		TargetSize:          int64(d.Get("target_size").(int)),
		NamedPorts:          getNamedPortsBeta(d.Get("named_port").(*schema.Set).List()),
		TargetPools:         convertStringSet(d.Get("target_pools").(*schema.Set)),
		AutoHealingPolicies: expandAutoHealingPolicies(d.Get("auto_healing_policies").([]interface{})),
		Versions:            expandVersions(d.Get("version").([]interface{})),
		UpdatePolicy:        expandRegionUpdatePolicy(d.Get("update_policy").([]interface{})),
		DistributionPolicy:  expandDistributionPolicy(d.Get("distribution_policy_zones").(*schema.Set)),
		// Force send TargetSize to allow size of 0.
		ForceSendFields: []string{"TargetSize"},
	}

	op, err := config.clientComputeBeta.RegionInstanceGroupManagers.Insert(project, region, manager).Do()

	if err != nil {
		return fmt.Errorf("Error creating RegionInstanceGroupManager: %s", err)
	}

	id, err := replaceVars(d, config, "projects/{{project}}/regions/{{region}}/instanceGroupManagers/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	// Wait for the operation to complete
	err = computeOperationWaitTime(config, op, project, "Creating InstanceGroupManager", d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}
	return resourceComputeRegionInstanceGroupManagerRead(d, config)
}

type getInstanceManagerFunc func(*schema.ResourceData, interface{}) (*computeBeta.InstanceGroupManager, error)

func getRegionalManager(d *schema.ResourceData, meta interface{}) (*computeBeta.InstanceGroupManager, error) {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return nil, err
	}

	name := d.Get("name").(string)
	manager, err := config.clientComputeBeta.RegionInstanceGroupManagers.Get(project, region, name).Do()
	if err != nil {
		return nil, handleNotFoundError(err, d, fmt.Sprintf("Region Instance Manager %q", name))
	}

	return manager, nil
}

func waitForInstancesRefreshFunc(f getInstanceManagerFunc, d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		m, err := f(d, meta)
		if err != nil {
			log.Printf("[WARNING] Error in fetching manager while waiting for instances to come up: %s\n", err)
			return nil, "error", err
		}
		if m.Status.IsStable {
			return true, "created", nil
		} else {
			return false, "creating", nil
		}
	}
}

func resourceComputeRegionInstanceGroupManagerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	manager, err := getRegionalManager(d, meta)
	if err != nil {
		return err
	}
	if manager == nil {
		log.Printf("[WARN] Region Instance Group Manager %q not found, removing from state.", d.Id())
		d.SetId("")
		return nil
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	d.Set("base_instance_name", manager.BaseInstanceName)
	d.Set("name", manager.Name)
	d.Set("region", GetResourceNameFromSelfLink(manager.Region))
	d.Set("description", manager.Description)
	d.Set("project", project)
	d.Set("target_size", manager.TargetSize)
	if err := d.Set("target_pools", mapStringArr(manager.TargetPools, ConvertSelfLinkToV1)); err != nil {
		return fmt.Errorf("Error setting target_pools in state: %s", err.Error())
	}
	if err := d.Set("named_port", flattenNamedPortsBeta(manager.NamedPorts)); err != nil {
		return fmt.Errorf("Error setting named_port in state: %s", err.Error())
	}
	d.Set("fingerprint", manager.Fingerprint)
	d.Set("instance_group", ConvertSelfLinkToV1(manager.InstanceGroup))
	if err := d.Set("distribution_policy_zones", flattenDistributionPolicy(manager.DistributionPolicy)); err != nil {
		return err
	}
	d.Set("self_link", ConvertSelfLinkToV1(manager.SelfLink))

	if err := d.Set("auto_healing_policies", flattenAutoHealingPolicies(manager.AutoHealingPolicies)); err != nil {
		return fmt.Errorf("Error setting auto_healing_policies in state: %s", err.Error())
	}
	if err := d.Set("version", flattenVersions(manager.Versions)); err != nil {
		return err
	}
	if err := d.Set("update_policy", flattenRegionUpdatePolicy(manager.UpdatePolicy)); err != nil {
		return fmt.Errorf("Error setting update_policy in state: %s", err.Error())
	}

	if d.Get("wait_for_instances").(bool) {
		conf := resource.StateChangeConf{
			Pending: []string{"creating", "error"},
			Target:  []string{"created"},
			Refresh: waitForInstancesRefreshFunc(getRegionalManager, d, meta),
			Timeout: d.Timeout(schema.TimeoutCreate),
		}
		_, err := conf.WaitForState()
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceComputeRegionInstanceGroupManagerUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	updatedManager := &computeBeta.InstanceGroupManager{
		Fingerprint: d.Get("fingerprint").(string),
	}
	var change bool

	if d.HasChange("target_pools") {
		updatedManager.TargetPools = convertStringSet(d.Get("target_pools").(*schema.Set))
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

	if change {
		op, err := config.clientComputeBeta.RegionInstanceGroupManagers.Patch(project, region, d.Get("name").(string), updatedManager).Do()
		if err != nil {
			return fmt.Errorf("Error updating region managed group instances: %s", err)
		}

		err = computeOperationWaitTime(config, op, project, "Updating region managed group instances", d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	}

	// named ports can't be updated through PATCH
	// so we call the update method on the region instance group, instead of the rigm
	if d.HasChange("named_port") {
		d.Partial(true)
		namedPorts := getNamedPortsBeta(d.Get("named_port").(*schema.Set).List())
		setNamedPorts := &computeBeta.RegionInstanceGroupsSetNamedPortsRequest{
			NamedPorts: namedPorts,
		}

		op, err := config.clientComputeBeta.RegionInstanceGroups.SetNamedPorts(
			project, region, d.Get("name").(string), setNamedPorts).Do()

		if err != nil {
			return fmt.Errorf("Error updating RegionInstanceGroupManager: %s", err)
		}

		err = computeOperationWaitTime(config, op, project, "Updating RegionInstanceGroupManager", d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
		d.SetPartial("named_port")
	}

	// target size should use resize
	if d.HasChange("target_size") {
		d.Partial(true)
		targetSize := int64(d.Get("target_size").(int))
		op, err := config.clientComputeBeta.RegionInstanceGroupManagers.Resize(
			project, region, d.Get("name").(string), targetSize).Do()

		if err != nil {
			return fmt.Errorf("Error resizing RegionInstanceGroupManager: %s", err)
		}

		err = computeOperationWaitTime(config, op, project, "Resizing RegionInstanceGroupManager", d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
		d.SetPartial("target_size")
	}

	d.Partial(false)

	return resourceComputeRegionInstanceGroupManagerRead(d, meta)
}

func resourceComputeRegionInstanceGroupManagerDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	op, err := config.clientComputeBeta.RegionInstanceGroupManagers.Delete(project, region, name).Do()

	if err != nil {
		return fmt.Errorf("Error deleting region instance group manager: %s", err)
	}

	// Wait for the operation to complete
	err = computeOperationWaitTime(config, op, project, "Deleting RegionInstanceGroupManager", d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf("Error waiting for delete to complete: %s", err)
	}

	d.SetId("")
	return nil
}

func expandRegionUpdatePolicy(configured []interface{}) *computeBeta.InstanceGroupManagerUpdatePolicy {
	updatePolicy := &computeBeta.InstanceGroupManagerUpdatePolicy{}

	for _, raw := range configured {
		data := raw.(map[string]interface{})

		updatePolicy.MinimalAction = data["minimal_action"].(string)
		updatePolicy.Type = data["type"].(string)
		updatePolicy.InstanceRedistributionType = data["instance_redistribution_type"].(string)

		// percent and fixed values are conflicting
		// when the percent values are set, the fixed values will be ignored
		if v := data["max_surge_percent"]; v.(int) > 0 {
			updatePolicy.MaxSurge = &computeBeta.FixedOrPercent{
				Percent:    int64(v.(int)),
				NullFields: []string{"Fixed"},
			}
		} else {
			updatePolicy.MaxSurge = &computeBeta.FixedOrPercent{
				Fixed: int64(data["max_surge_fixed"].(int)),
				// allow setting this value to 0
				ForceSendFields: []string{"Fixed"},
				NullFields:      []string{"Percent"},
			}
		}

		if v := data["max_unavailable_percent"]; v.(int) > 0 {
			updatePolicy.MaxUnavailable = &computeBeta.FixedOrPercent{
				Percent:    int64(v.(int)),
				NullFields: []string{"Fixed"},
			}
		} else {
			updatePolicy.MaxUnavailable = &computeBeta.FixedOrPercent{
				Fixed: int64(data["max_unavailable_fixed"].(int)),
				// allow setting this value to 0
				ForceSendFields: []string{"Fixed"},
				NullFields:      []string{"Percent"},
			}
		}

		if v, ok := data["min_ready_sec"]; ok {
			updatePolicy.MinReadySec = int64(v.(int))
		}
	}
	return updatePolicy
}

func flattenRegionUpdatePolicy(updatePolicy *computeBeta.InstanceGroupManagerUpdatePolicy) []map[string]interface{} {
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
		up["min_ready_sec"] = updatePolicy.MinReadySec
		up["minimal_action"] = updatePolicy.MinimalAction
		up["type"] = updatePolicy.Type
		up["instance_redistribution_type"] = updatePolicy.InstanceRedistributionType

		results = append(results, up)
	}
	return results
}

func expandDistributionPolicy(configured *schema.Set) *computeBeta.DistributionPolicy {
	if configured.Len() == 0 {
		return nil
	}

	distributionPolicyZoneConfigs := make([]*computeBeta.DistributionPolicyZoneConfiguration, 0, configured.Len())
	for _, raw := range configured.List() {
		data := raw.(string)
		distributionPolicyZoneConfig := computeBeta.DistributionPolicyZoneConfiguration{
			Zone: "zones/" + data,
		}

		distributionPolicyZoneConfigs = append(distributionPolicyZoneConfigs, &distributionPolicyZoneConfig)
	}
	return &computeBeta.DistributionPolicy{Zones: distributionPolicyZoneConfigs}
}

func flattenDistributionPolicy(distributionPolicy *computeBeta.DistributionPolicy) []string {
	zones := make([]string, 0)

	if distributionPolicy != nil {
		for _, zone := range distributionPolicy.Zones {
			zones = append(zones, GetResourceNameFromSelfLink(zone.Zone))
		}
	}

	return zones
}

func hashZoneFromSelfLinkOrResourceName(value interface{}) int {
	parts := strings.Split(value.(string), "/")
	resource := parts[len(parts)-1]

	return hashcode.String(resource)
}

func resourceRegionInstanceGroupManagerStateImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	d.Set("wait_for_instances", false)
	config := meta.(*Config)
	if err := parseImportId([]string{"projects/(?P<project>[^/]+)/regions/(?P<region>[^/]+)/instanceGroupManagers/(?P<name>[^/]+)", "(?P<project>[^/]+)/(?P<region>[^/]+)/(?P<name>[^/]+)", "(?P<region>[^/]+)/(?P<name>[^/]+)", "(?P<name>[^/]+)"}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "projects/{{project}}/regions/{{region}}/instanceGroupManagers/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
