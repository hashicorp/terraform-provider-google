package google

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
)

func resourceComputeInstanceGroupManager() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeInstanceGroupManagerCreate,
		Read:   resourceComputeInstanceGroupManagerRead,
		Update: resourceComputeInstanceGroupManagerUpdate,
		Delete: resourceComputeInstanceGroupManagerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceInstanceGroupManagerStateImporter,
		},

		Schema: map[string]*schema.Schema{
			"base_instance_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"instance_template": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: compareSelfLinkRelativePaths,
			},

			"version": &schema.Schema{
				Deprecated: "This field is in beta and will be removed from this provider. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/provider_versions.html for more details.",
				Type:       schema.TypeList,
				Optional:   true,
				Computed:   true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"instance_template": &schema.Schema{
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: compareSelfLinkRelativePaths,
						},

						"target_size": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"fixed": &schema.Schema{
										Type:     schema.TypeInt,
										Optional: true,
									},

									"percent": &schema.Schema{
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

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"zone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"fingerprint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"instance_group": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"named_port": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"port": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"update_strategy": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "REPLACE",
				ValidateFunc: validation.StringInSlice([]string{"RESTART", "NONE", "ROLLING_UPDATE", "REPLACE"}, false),
				DiffSuppressFunc: func(key, old, new string, d *schema.ResourceData) bool {
					if old == "REPLACE" && new == "RESTART" {
						return true
					}
					if old == "RESTART" && new == "REPLACE" {
						return true
					}
					return false
				},
			},

			"target_pools": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: selfLinkRelativePathHash,
			},

			"target_size": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},

			"auto_healing_policies": &schema.Schema{
				Deprecated: "This field is in beta and will be removed from this provider. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/provider_versions.html for more details.",
				Type:       schema.TypeList,
				Optional:   true,
				MaxItems:   1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health_check": &schema.Schema{
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: compareSelfLinkRelativePaths,
						},

						"initial_delay_sec": &schema.Schema{
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 3600),
						},
					},
				},
			},

			"rolling_update_policy": &schema.Schema{
				Deprecated: "This field is in beta and will be removed from this provider. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/provider_versions.html for more details.",
				Type:       schema.TypeList,
				Optional:   true,
				MaxItems:   1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"minimal_action": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"RESTART", "REPLACE"}, false),
						},

						"type": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"OPPORTUNISTIC", "PROACTIVE"}, false),
						},

						"max_surge_fixed": &schema.Schema{
							Type:          schema.TypeInt,
							Optional:      true,
							Default:       1,
							ConflictsWith: []string{"rolling_update_policy.0.max_surge_percent"},
						},

						"max_surge_percent": &schema.Schema{
							Type:          schema.TypeInt,
							Optional:      true,
							ConflictsWith: []string{"rolling_update_policy.0.max_surge_fixed"},
							ValidateFunc:  validation.IntBetween(0, 100),
						},

						"max_unavailable_fixed": &schema.Schema{
							Type:          schema.TypeInt,
							Optional:      true,
							Default:       1,
							ConflictsWith: []string{"rolling_update_policy.0.max_unavailable_percent"},
						},

						"max_unavailable_percent": &schema.Schema{
							Type:          schema.TypeInt,
							Optional:      true,
							ConflictsWith: []string{"rolling_update_policy.0.max_unavailable_fixed"},
							ValidateFunc:  validation.IntBetween(0, 100),
						},

						"min_ready_sec": &schema.Schema{
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 3600),
						},
					},
				},
			},

			"wait_for_instances": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func getNamedPorts(nps []interface{}) []*compute.NamedPort {
	namedPorts := make([]*compute.NamedPort, 0, len(nps))
	for _, v := range nps {
		np := v.(map[string]interface{})
		namedPorts = append(namedPorts, &compute.NamedPort{
			Name: np["name"].(string),
			Port: int64(np["port"].(int)),
		})
	}

	return namedPorts
}

func getNamedPortsBeta(nps []interface{}) []*computeBeta.NamedPort {
	namedPorts := make([]*computeBeta.NamedPort, 0, len(nps))
	for _, v := range nps {
		np := v.(map[string]interface{})
		namedPorts = append(namedPorts, &computeBeta.NamedPort{
			Name: np["name"].(string),
			Port: int64(np["port"].(int)),
		})
	}

	return namedPorts
}

func resourceComputeInstanceGroupManagerCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone, err := getZone(d, config)
	if err != nil {
		return err
	}

	if _, ok := d.GetOk("rolling_update_policy"); d.Get("update_strategy") == "ROLLING_UPDATE" && !ok {
		return fmt.Errorf("[rolling_update_policy] must be set when 'update_strategy' is set to 'ROLLING_UPDATE'")
	}

	// Build the parameter
	manager := &computeBeta.InstanceGroupManager{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		BaseInstanceName:    d.Get("base_instance_name").(string),
		InstanceTemplate:    d.Get("instance_template").(string),
		TargetSize:          int64(d.Get("target_size").(int)),
		NamedPorts:          getNamedPortsBeta(d.Get("named_port").([]interface{})),
		TargetPools:         convertStringSet(d.Get("target_pools").(*schema.Set)),
		AutoHealingPolicies: expandAutoHealingPolicies(d.Get("auto_healing_policies").([]interface{})),
		Versions:            expandVersions(d.Get("version").([]interface{})),
		// Force send TargetSize to allow a value of 0.
		ForceSendFields: []string{"TargetSize"},
	}

	log.Printf("[DEBUG] InstanceGroupManager insert request: %#v", manager)
	op, err := config.clientComputeBeta.InstanceGroupManagers.Insert(
		project, zone, manager).Do()

	if err != nil {
		return fmt.Errorf("Error creating InstanceGroupManager: %s", err)
	}

	// It probably maybe worked, so store the ID now
	d.SetId(manager.Name)

	// Wait for the operation to complete
	err = computeSharedOperationWait(config.clientCompute, op, project, "Creating InstanceGroupManager")
	if err != nil {
		return err
	}

	return resourceComputeInstanceGroupManagerRead(d, meta)
}

func flattenNamedPortsBeta(namedPorts []*computeBeta.NamedPort) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(namedPorts))
	for _, namedPort := range namedPorts {
		namedPortMap := make(map[string]interface{})
		namedPortMap["name"] = namedPort.Name
		namedPortMap["port"] = namedPort.Port
		result = append(result, namedPortMap)
	}
	return result

}

func flattenVersions(versions []*computeBeta.InstanceGroupManagerVersion) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(versions))
	for _, version := range versions {
		versionMap := make(map[string]interface{})
		versionMap["name"] = version.Name
		versionMap["instance_template"] = ConvertSelfLinkToV1(version.InstanceTemplate)
		versionMap["target_size"] = flattenFixedOrPercent(version.TargetSize)
		result = append(result, versionMap)
	}

	return result
}

func flattenFixedOrPercent(fixedOrPercent *computeBeta.FixedOrPercent) []map[string]interface{} {
	result := make(map[string]interface{})
	if value := fixedOrPercent.Percent; value > 0 {
		result["percent"] = value
	} else if value := fixedOrPercent.Fixed; value > 0 {
		result["fixed"] = fixedOrPercent.Fixed
	} else {
		return []map[string]interface{}{}
	}
	return []map[string]interface{}{result}
}

func getManager(d *schema.ResourceData, meta interface{}) (*computeBeta.InstanceGroupManager, error) {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return nil, err
	}

	getInstanceGroupManager := func(zone string) (interface{}, error) {
		return config.clientComputeBeta.InstanceGroupManagers.Get(project, zone, d.Id()).Do()
	}

	var manager *computeBeta.InstanceGroupManager
	var e error
	if zone, _ := getZone(d, config); zone != "" {
		manager, e = config.clientComputeBeta.InstanceGroupManagers.Get(project, zone, d.Id()).Do()

		if e != nil {
			return nil, handleNotFoundError(e, d, fmt.Sprintf("Instance Group Manager %q", d.Get("name").(string)))
		}
	} else {
		// If the resource was imported, the only info we have is the ID. Try to find the resource
		// by searching in the region of the project.
		var resource interface{}
		resource, e = getZonalBetaResourceFromRegion(getInstanceGroupManager, region, config.clientComputeBeta, project)
		if e != nil {
			return nil, e
		}

		manager = resource.(*computeBeta.InstanceGroupManager)
	}

	if manager == nil {
		log.Printf("[WARN] Removing Instance Group Manager %q because it's gone", d.Get("name").(string))

		// The resource doesn't exist anymore
		d.SetId("")
		return nil, nil
	}

	return manager, nil
}

func resourceComputeInstanceGroupManagerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	manager, err := getManager(d, meta)
	if err != nil || manager == nil {
		return err
	}

	d.Set("base_instance_name", manager.BaseInstanceName)
	d.Set("instance_template", ConvertSelfLinkToV1(manager.InstanceTemplate))
	if err := d.Set("version", flattenVersions(manager.Versions)); err != nil {
		return err
	}
	d.Set("name", manager.Name)
	d.Set("zone", GetResourceNameFromSelfLink(manager.Zone))
	d.Set("description", manager.Description)
	d.Set("project", project)
	d.Set("target_size", manager.TargetSize)
	d.Set("target_pools", manager.TargetPools)
	d.Set("named_port", flattenNamedPortsBeta(manager.NamedPorts))
	d.Set("fingerprint", manager.Fingerprint)
	d.Set("instance_group", ConvertSelfLinkToV1(manager.InstanceGroup))
	d.Set("self_link", ConvertSelfLinkToV1(manager.SelfLink))
	update_strategy, ok := d.GetOk("update_strategy")
	if !ok {
		update_strategy = "REPLACE"
	}
	d.Set("update_strategy", update_strategy.(string))
	d.Set("auto_healing_policies", flattenAutoHealingPolicies(manager.AutoHealingPolicies))

	if d.Get("wait_for_instances").(bool) {
		conf := resource.StateChangeConf{
			Pending: []string{"creating", "error"},
			Target:  []string{"created"},
			Refresh: waitForInstancesRefreshFunc(getManager, d, meta),
			Timeout: d.Timeout(schema.TimeoutCreate),
		}
		_, err := conf.WaitForState()
		if err != nil {
			return err
		}
	}

	return nil
}

// Updates an instance group manager by applying the update strategy (REPLACE, RESTART)
// and rolling update policy (PROACTIVE, OPPORTUNISTIC). Updates performed by API
// are OPPORTUNISTIC by default.
func performZoneUpdate(config *Config, id string, updateStrategy string, rollingUpdatePolicy *computeBeta.InstanceGroupManagerUpdatePolicy, versions []*computeBeta.InstanceGroupManagerVersion, project string, zone string) error {
	if updateStrategy == "RESTART" || updateStrategy == "REPLACE" {
		managedInstances, err := config.clientComputeBeta.InstanceGroupManagers.ListManagedInstances(project, zone, id).Do()
		if err != nil {
			return fmt.Errorf("Error getting instance group managers instances: %s", err)
		}

		managedInstanceCount := len(managedInstances.ManagedInstances)
		instances := make([]string, managedInstanceCount)
		for i, v := range managedInstances.ManagedInstances {
			instances[i] = v.Instance
		}

		recreateInstances := &computeBeta.InstanceGroupManagersRecreateInstancesRequest{
			Instances: instances,
		}

		op, err := config.clientComputeBeta.InstanceGroupManagers.RecreateInstances(project, zone, id, recreateInstances).Do()
		if err != nil {
			return fmt.Errorf("Error restarting instance group managers instances: %s", err)
		}

		// Wait for the operation to complete
		err = computeSharedOperationWaitTime(config.clientCompute, op, project, managedInstanceCount*4, "Restarting InstanceGroupManagers instances")
		if err != nil {
			return err
		}
	}

	if updateStrategy == "ROLLING_UPDATE" {
		// UpdatePolicy is set for InstanceGroupManager on update only, because it is only relevant for `Patch` calls.
		// Other tools(gcloud and UI) capable of executing the same `ROLLING UPDATE` call
		// expect those values to be provided by user as part of the call
		// or provide their own defaults without respecting what was previously set on UpdateManager.
		// To follow the same logic, we provide policy values on relevant update change only.
		manager := &computeBeta.InstanceGroupManager{
			UpdatePolicy: rollingUpdatePolicy,
			Versions:     versions,
		}

		op, err := config.clientComputeBeta.InstanceGroupManagers.Patch(project, zone, id, manager).Do()
		if err != nil {
			return fmt.Errorf("Error updating managed group instances: %s", err)
		}

		err = computeSharedOperationWait(config.clientCompute, op, project, "Updating managed group instances")
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceComputeInstanceGroupManagerUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone, err := getZone(d, config)
	if err != nil {
		return err
	}

	d.Partial(true)

	if _, ok := d.GetOk("rolling_update_policy"); d.Get("update_strategy") == "ROLLING_UPDATE" && !ok {
		return fmt.Errorf("[rolling_update_policy] must be set when 'update_strategy' is set to 'ROLLING_UPDATE'")
	}

	// If target_pools changes then update
	if d.HasChange("target_pools") {
		targetPools := convertStringSet(d.Get("target_pools").(*schema.Set))

		// Build the parameter
		setTargetPools := &computeBeta.InstanceGroupManagersSetTargetPoolsRequest{
			Fingerprint: d.Get("fingerprint").(string),
			TargetPools: targetPools,
		}

		op, err := config.clientComputeBeta.InstanceGroupManagers.SetTargetPools(
			project, zone, d.Id(), setTargetPools).Do()

		if err != nil {
			return fmt.Errorf("Error updating InstanceGroupManager: %s", err)
		}

		// Wait for the operation to complete
		err = computeSharedOperationWait(config.clientCompute, op, project, "Updating InstanceGroupManager")
		if err != nil {
			return err
		}

		d.SetPartial("target_pools")
	}

	// If named_port changes then update:
	if d.HasChange("named_port") {

		// Build the parameters for a "SetNamedPorts" request:
		namedPorts := getNamedPortsBeta(d.Get("named_port").([]interface{}))
		setNamedPorts := &computeBeta.InstanceGroupsSetNamedPortsRequest{
			NamedPorts: namedPorts,
		}

		// Make the request:
		op, err := config.clientComputeBeta.InstanceGroups.SetNamedPorts(
			project, zone, d.Id(), setNamedPorts).Do()

		if err != nil {
			return fmt.Errorf("Error updating InstanceGroupManager: %s", err)
		}

		// Wait for the operation to complete:
		err = computeSharedOperationWait(config.clientCompute, op, project, "Updating InstanceGroupManager")
		if err != nil {
			return err
		}

		d.SetPartial("named_port")
	}

	if d.HasChange("target_size") {
		targetSize := int64(d.Get("target_size").(int))
		op, err := config.clientComputeBeta.InstanceGroupManagers.Resize(
			project, zone, d.Id(), targetSize).Do()

		if err != nil {
			return fmt.Errorf("Error updating InstanceGroupManager: %s", err)
		}

		// Wait for the operation to complete
		err = computeSharedOperationWait(config.clientCompute, op, project, "Updating InstanceGroupManager")
		if err != nil {
			return err
		}

		d.SetPartial("target_size")
	}

	// We will always be in v0beta inside this conditional
	if d.HasChange("auto_healing_policies") {
		setAutoHealingPoliciesRequest := &computeBeta.InstanceGroupManagersSetAutoHealingRequest{}
		if v, ok := d.GetOk("auto_healing_policies"); ok {
			setAutoHealingPoliciesRequest.AutoHealingPolicies = expandAutoHealingPolicies(v.([]interface{}))
		}

		op, err := config.clientComputeBeta.InstanceGroupManagers.SetAutoHealingPolicies(
			project, zone, d.Id(), setAutoHealingPoliciesRequest).Do()

		if err != nil {
			return fmt.Errorf("Error updating AutoHealingPolicies: %s", err)
		}

		// Wait for the operation to complete
		err = computeSharedOperationWait(config.clientCompute, op, project, "Updating AutoHealingPolicies")
		if err != nil {
			return err
		}

		d.SetPartial("auto_healing_policies")
	}

	// If instance_template changes then update
	if d.HasChange("instance_template") {
		// Build the parameter
		setInstanceTemplate := &computeBeta.InstanceGroupManagersSetInstanceTemplateRequest{
			InstanceTemplate: d.Get("instance_template").(string),
		}

		op, err := config.clientComputeBeta.InstanceGroupManagers.SetInstanceTemplate(project, zone, d.Id(), setInstanceTemplate).Do()

		if err != nil {
			return fmt.Errorf("Error updating InstanceGroupManager: %s", err)
		}

		// Wait for the operation to complete
		err = computeSharedOperationWait(config.clientCompute, op, project, "Updating InstanceGroupManager")
		if err != nil {
			return err
		}

		updateStrategy := d.Get("update_strategy").(string)
		rollingUpdatePolicy := expandUpdatePolicy(d.Get("rolling_update_policy").([]interface{}))
		err = performZoneUpdate(config, d.Id(), updateStrategy, rollingUpdatePolicy, nil, project, zone)
		d.SetPartial("instance_template")
	}

	// If version changes then update
	if d.HasChange("version") {
		updateStrategy := d.Get("update_strategy").(string)
		rollingUpdatePolicy := expandUpdatePolicy(d.Get("rolling_update_policy").([]interface{}))
		versions := expandVersions(d.Get("version").([]interface{}))
		err = performZoneUpdate(config, d.Id(), updateStrategy, rollingUpdatePolicy, versions, project, zone)
		if err != nil {
			return err
		}

		d.SetPartial("version")
	}

	d.Partial(false)

	return resourceComputeInstanceGroupManagerRead(d, meta)
}

func resourceComputeInstanceGroupManagerDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone, err := getZone(d, config)
	if err != nil {
		return err
	}

	op, err := config.clientComputeBeta.InstanceGroupManagers.Delete(project, zone, d.Id()).Do()
	attempt := 0
	for err != nil && attempt < 20 {
		attempt++
		time.Sleep(2000 * time.Millisecond)
		op, err = config.clientComputeBeta.InstanceGroupManagers.Delete(project, zone, d.Id()).Do()
	}

	if err != nil {
		return fmt.Errorf("Error deleting instance group manager: %s", err)
	}

	currentSize := int64(d.Get("target_size").(int))

	// Wait for the operation to complete
	err = computeSharedOperationWait(config.clientCompute, op, project, "Deleting InstanceGroupManager")

	for err != nil && currentSize > 0 {
		if !strings.Contains(err.Error(), "timeout") {
			return err
		}

		instanceGroup, err := config.clientComputeBeta.InstanceGroups.Get(
			project, zone, d.Id()).Do()
		if err != nil {
			return fmt.Errorf("Error getting instance group size: %s", err)
		}

		instanceGroupSize := instanceGroup.Size

		if instanceGroupSize >= currentSize {
			return fmt.Errorf("Error, instance group isn't shrinking during delete")
		}

		log.Printf("[INFO] timeout occured, but instance group is shrinking (%d < %d)", instanceGroupSize, currentSize)
		currentSize = instanceGroupSize
		err = computeSharedOperationWait(config.clientCompute, op, project, "Deleting InstanceGroupManager")
	}

	d.SetId("")
	return nil
}

func expandAutoHealingPolicies(configured []interface{}) []*computeBeta.InstanceGroupManagerAutoHealingPolicy {
	autoHealingPolicies := make([]*computeBeta.InstanceGroupManagerAutoHealingPolicy, 0, len(configured))
	for _, raw := range configured {
		data := raw.(map[string]interface{})
		autoHealingPolicy := computeBeta.InstanceGroupManagerAutoHealingPolicy{
			HealthCheck:     data["health_check"].(string),
			InitialDelaySec: int64(data["initial_delay_sec"].(int)),
		}

		autoHealingPolicies = append(autoHealingPolicies, &autoHealingPolicy)
	}
	return autoHealingPolicies
}

func expandVersions(configured []interface{}) []*computeBeta.InstanceGroupManagerVersion {
	versions := make([]*computeBeta.InstanceGroupManagerVersion, 0, len(configured))
	for _, raw := range configured {
		data := raw.(map[string]interface{})

		version := computeBeta.InstanceGroupManagerVersion{
			Name:             data["name"].(string),
			InstanceTemplate: data["instance_template"].(string),
			TargetSize:       expandFixedOrPercent(data["target_size"].([]interface{})),
		}

		versions = append(versions, &version)
	}
	return versions
}

func expandFixedOrPercent(configured []interface{}) *computeBeta.FixedOrPercent {
	fixedOrPercent := &computeBeta.FixedOrPercent{}

	for _, raw := range configured {
		data := raw.(map[string]interface{})
		if percent := data["percent"]; percent.(int) > 0 {
			fixedOrPercent.Percent = int64(percent.(int))
		} else {
			fixedOrPercent.Fixed = int64(data["fixed"].(int))
			fixedOrPercent.ForceSendFields = []string{"Fixed"}
		}
	}
	return fixedOrPercent
}

func expandUpdatePolicy(configured []interface{}) *computeBeta.InstanceGroupManagerUpdatePolicy {
	updatePolicy := &computeBeta.InstanceGroupManagerUpdatePolicy{}

	for _, raw := range configured {
		data := raw.(map[string]interface{})

		updatePolicy.MinimalAction = data["minimal_action"].(string)
		updatePolicy.Type = data["type"].(string)

		// percent and fixed values are conflicting
		// when the percent values are set, the fixed values will be ignored
		if v := data["max_surge_percent"]; v.(int) > 0 {
			updatePolicy.MaxSurge = &computeBeta.FixedOrPercent{
				Percent: int64(v.(int)),
			}
		} else {
			updatePolicy.MaxSurge = &computeBeta.FixedOrPercent{
				Fixed: int64(data["max_surge_fixed"].(int)),
				// allow setting this value to 0
				ForceSendFields: []string{"Fixed"},
			}
		}

		if v := data["max_unavailable_percent"]; v.(int) > 0 {
			updatePolicy.MaxUnavailable = &computeBeta.FixedOrPercent{
				Percent: int64(v.(int)),
			}
		} else {
			updatePolicy.MaxUnavailable = &computeBeta.FixedOrPercent{
				Fixed: int64(data["max_unavailable_fixed"].(int)),
				// allow setting this value to 0
				ForceSendFields: []string{"Fixed"},
			}
		}

		if v, ok := data["min_ready_sec"]; ok {
			updatePolicy.MinReadySec = int64(v.(int))
		}
	}
	return updatePolicy
}

func flattenAutoHealingPolicies(autoHealingPolicies []*computeBeta.InstanceGroupManagerAutoHealingPolicy) []map[string]interface{} {
	autoHealingPoliciesSchema := make([]map[string]interface{}, 0, len(autoHealingPolicies))
	for _, autoHealingPolicy := range autoHealingPolicies {
		data := map[string]interface{}{
			"health_check":      autoHealingPolicy.HealthCheck,
			"initial_delay_sec": autoHealingPolicy.InitialDelaySec,
		}

		autoHealingPoliciesSchema = append(autoHealingPoliciesSchema, data)
	}
	return autoHealingPoliciesSchema
}

func resourceInstanceGroupManagerStateImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	d.Set("wait_for_instances", false)
	return []*schema.ResourceData{d}, nil
}
