//
package google

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
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
							DiffSuppressFunc: compareSelfLinkRelativePaths,
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

			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The zone that instances in this group should be created in.`,
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
				Set:         selfLinkRelativePathHash,
				Description: `The full URL of all target pools to which new instances in the group are added. Updating the target pools attribute does not affect existing instances.`,
			},

			"target_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: `The target number of running instances for this managed instance group. This value should always be explicitly set unless this resource is attached to an autoscaler, in which case it should never be set. Defaults to 0.`,
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
							DiffSuppressFunc: compareSelfLinkRelativePaths,
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

			"update_policy": {
				Computed:    true,
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: `The update policy for this managed instance group.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"minimal_action": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"RESTART", "REPLACE"}, false),
							Description:  `Minimal action to be taken on an instance. You can specify either RESTART to restart existing instances or REPLACE to delete and create new instances from the target template. If you specify a RESTART, the Updater will attempt to perform that action only. However, if the Updater determines that the minimal action you specify is not enough to perform the update, it might perform a more disruptive action.`,
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
							Description:   `The maximum number of instances that can be created above the specified targetSize during the update process. Conflicts with max_surge_percent. If neither is set, defaults to 1`,
						},

						"max_surge_percent": {
							Type:          schema.TypeInt,
							Optional:      true,
							ConflictsWith: []string{"update_policy.0.max_surge_fixed"},
							ValidateFunc:  validation.IntBetween(0, 100),
							Description:   `The maximum number of instances(calculated as percentage) that can be created above the specified targetSize during the update process. Conflicts with max_surge_fixed.`,
						},

						"max_unavailable_fixed": {
							Type:          schema.TypeInt,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"update_policy.0.max_unavailable_percent"},
							Description:   `The maximum number of instances that can be unavailable during the update process. Conflicts with max_unavailable_percent. If neither is set, defaults to 1.`,
						},

						"max_unavailable_percent": {
							Type:          schema.TypeInt,
							Optional:      true,
							ConflictsWith: []string{"update_policy.0.max_unavailable_fixed"},
							ValidateFunc:  validation.IntBetween(0, 100),
							Description:   `The maximum number of instances(calculated as percentage) that can be unavailable during the update process. Conflicts with max_unavailable_fixed.`,
						},

						"min_ready_sec": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 3600),
							Description:  `Minimum number of seconds to wait for after a newly created instance becomes available. This value must be from range [0, 3600].`,
						},
					},
				},
			},

			"wait_for_instances": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: `Whether to wait for all instances to be created/updated before returning. Note that if this is set to true and the operation does not succeed, Terraform will continue trying until it times out.`,
			},
			"operation": {
				Type:     schema.TypeString,
				Computed: true,
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

	// Build the parameter
	manager := &computeBeta.InstanceGroupManager{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		BaseInstanceName:    d.Get("base_instance_name").(string),
		TargetSize:          int64(d.Get("target_size").(int)),
		NamedPorts:          getNamedPortsBeta(d.Get("named_port").(*schema.Set).List()),
		TargetPools:         convertStringSet(d.Get("target_pools").(*schema.Set)),
		AutoHealingPolicies: expandAutoHealingPolicies(d.Get("auto_healing_policies").([]interface{})),
		Versions:            expandVersions(d.Get("version").([]interface{})),
		UpdatePolicy:        expandUpdatePolicy(d.Get("update_policy").([]interface{})),
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
	id, err := replaceVars(d, config, "projects/{{project}}/zones/{{zone}}/instanceGroupManagers/{{name}}")
	if err != nil {
		return err
	}
	d.SetId(id)

	// Wait for the operation to complete
	err = computeOperationWaitTime(config, op, project, "Creating InstanceGroupManager", d.Timeout(schema.TimeoutCreate))
	if err != nil {
		// Check if the create operation failed because Terraform was prematurely terminated. If it was we can persist the
		// operation id to state so that a subsequent refresh of this resource will wait until the operation has terminated
		// before attempting to Read the state of the manager. This allows a graceful resumption of a Create that was killed
		// by the upstream Terraform process exiting early such as a sigterm.
		select {
		case <-config.context.Done():
			log.Printf("[DEBUG] Persisting %s so this operation can be resumed \n", op.Name)
			d.Set("operation", op.Name)
			return nil
		default:
			// leaving default case to ensure this is non blocking
		}
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

	zone, _ := getZone(d, config)
	name := d.Get("name").(string)

	manager, err := config.clientComputeBeta.InstanceGroupManagers.Get(project, zone, name).Do()
	if err != nil {
		return nil, handleNotFoundError(err, d, fmt.Sprintf("Instance Group Manager %q", name))
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

	operation := d.Get("operation").(string)
	if operation != "" {
		log.Printf("[DEBUG] in progress operation detected at %v, attempting to resume", operation)
		zone, _ := getZone(d, config)
		op := &computeBeta.Operation{
			Name: operation,
			Zone: zone,
		}
		d.Set("operation", "")
		err = computeOperationWaitTime(config, op, project, "Creating InstanceGroupManager", d.Timeout(schema.TimeoutCreate))
		if err != nil {
			// remove from state to allow refresh to finish
			log.Printf("[DEBUG] Resumed operation returned an error, removing from state: %s", err)
			d.SetId("")
			return nil
		}
	}

	manager, err := getManager(d, meta)
	if err != nil {
		return err
	}
	if manager == nil {
		log.Printf("[WARN] Instance Group Manager %q not found, removing from state.", d.Id())
		d.SetId("")
		return nil
	}

	if err := d.Set("base_instance_name", manager.BaseInstanceName); err != nil {
		return fmt.Errorf("Error reading base_instance_name: %s", err)
	}
	if err := d.Set("name", manager.Name); err != nil {
		return fmt.Errorf("Error reading name: %s", err)
	}
	if err := d.Set("zone", GetResourceNameFromSelfLink(manager.Zone)); err != nil {
		return fmt.Errorf("Error reading zone: %s", err)
	}
	if err := d.Set("description", manager.Description); err != nil {
		return fmt.Errorf("Error reading description: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading project: %s", err)
	}
	if err := d.Set("target_size", manager.TargetSize); err != nil {
		return fmt.Errorf("Error reading target_size: %s", err)
	}
	if err = d.Set("target_pools", mapStringArr(manager.TargetPools, ConvertSelfLinkToV1)); err != nil {
		return fmt.Errorf("Error setting target_pools in state: %s", err.Error())
	}
	if err = d.Set("named_port", flattenNamedPortsBeta(manager.NamedPorts)); err != nil {
		return fmt.Errorf("Error setting named_port in state: %s", err.Error())
	}
	if err := d.Set("fingerprint", manager.Fingerprint); err != nil {
		return fmt.Errorf("Error reading fingerprint: %s", err)
	}
	if err := d.Set("instance_group", ConvertSelfLinkToV1(manager.InstanceGroup)); err != nil {
		return fmt.Errorf("Error reading instance_group: %s", err)
	}
	if err := d.Set("self_link", ConvertSelfLinkToV1(manager.SelfLink)); err != nil {
		return fmt.Errorf("Error reading self_link: %s", err)
	}

	if err = d.Set("auto_healing_policies", flattenAutoHealingPolicies(manager.AutoHealingPolicies)); err != nil {
		return fmt.Errorf("Error setting auto_healing_policies in state: %s", err.Error())
	}
	if err := d.Set("version", flattenVersions(manager.Versions)); err != nil {
		return err
	}
	if err = d.Set("update_policy", flattenUpdatePolicy(manager.UpdatePolicy)); err != nil {
		return fmt.Errorf("Error setting update_policy in state: %s", err.Error())
	}

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

	updatedManager := &computeBeta.InstanceGroupManager{
		Fingerprint: d.Get("fingerprint").(string),
	}
	var change bool

	if d.HasChange("target_pools") {
		updatedManager.TargetPools = convertStringSet(d.Get("target_pools").(*schema.Set))
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
		updatedManager.UpdatePolicy = expandUpdatePolicy(d.Get("update_policy").([]interface{}))
		change = true
	}

	if change {
		op, err := config.clientComputeBeta.InstanceGroupManagers.Patch(project, zone, d.Get("name").(string), updatedManager).Do()
		if err != nil {
			return fmt.Errorf("Error updating managed group instances: %s", err)
		}

		err = computeOperationWaitTime(config, op, project, "Updating managed group instances", d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	}

	// named ports can't be updated through PATCH
	// so we call the update method on the instance group, instead of the igm
	if d.HasChange("named_port") {
		d.Partial(true)

		// Build the parameters for a "SetNamedPorts" request:
		namedPorts := getNamedPortsBeta(d.Get("named_port").(*schema.Set).List())
		setNamedPorts := &computeBeta.InstanceGroupsSetNamedPortsRequest{
			NamedPorts: namedPorts,
		}

		// Make the request:
		op, err := config.clientComputeBeta.InstanceGroups.SetNamedPorts(
			project, zone, d.Get("name").(string), setNamedPorts).Do()

		if err != nil {
			return fmt.Errorf("Error updating InstanceGroupManager: %s", err)
		}

		// Wait for the operation to complete:
		err = computeOperationWaitTime(config, op, project, "Updating InstanceGroupManager", d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	}

	// target_size should be updated through resize
	if d.HasChange("target_size") {
		d.Partial(true)

		targetSize := int64(d.Get("target_size").(int))
		op, err := config.clientComputeBeta.InstanceGroupManagers.Resize(
			project, zone, d.Get("name").(string), targetSize).Do()

		if err != nil {
			return fmt.Errorf("Error updating InstanceGroupManager: %s", err)
		}

		// Wait for the operation to complete
		err = computeOperationWaitTime(config, op, project, "Updating InstanceGroupManager", d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
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

	zone, _ := getZone(d, config)
	name := d.Get("name").(string)

	op, err := config.clientComputeBeta.InstanceGroupManagers.Delete(project, zone, name).Do()
	attempt := 0
	for err != nil && attempt < 20 {
		attempt++
		time.Sleep(2000 * time.Millisecond)
		op, err = config.clientComputeBeta.InstanceGroupManagers.Delete(project, zone, name).Do()
	}

	if err != nil {
		return fmt.Errorf("Error deleting instance group manager: %s", err)
	}

	currentSize := int64(d.Get("target_size").(int))

	// Wait for the operation to complete
	err = computeOperationWaitTime(config, op, project, "Deleting InstanceGroupManager", d.Timeout(schema.TimeoutDelete))

	for err != nil && currentSize > 0 {
		if !strings.Contains(err.Error(), "timeout") {
			return err
		}

		instanceGroup, igErr := config.clientComputeBeta.InstanceGroups.Get(
			project, zone, name).Do()
		if igErr != nil {
			return fmt.Errorf("Error getting instance group size: %s", err)
		}

		instanceGroupSize := instanceGroup.Size

		if instanceGroupSize >= currentSize {
			return fmt.Errorf("Error, instance group isn't shrinking during delete")
		}

		log.Printf("[INFO] timeout occurred, but instance group is shrinking (%d < %d)", instanceGroupSize, currentSize)
		currentSize = instanceGroupSize
		err = computeOperationWaitTime(config, op, project, "Deleting InstanceGroupManager", d.Timeout(schema.TimeoutDelete))
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

func flattenUpdatePolicy(updatePolicy *computeBeta.InstanceGroupManagerUpdatePolicy) []map[string]interface{} {
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
		results = append(results, up)
	}
	return results
}

func resourceInstanceGroupManagerStateImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := d.Set("wait_for_instances", false); err != nil {
		return nil, fmt.Errorf("Error reading wait_for_instances: %s", err)
	}
	config := meta.(*Config)
	if err := parseImportId([]string{"projects/(?P<project>[^/]+)/zones/(?P<zone>[^/]+)/instanceGroupManagers/(?P<name>[^/]+)", "(?P<project>[^/]+)/(?P<zone>[^/]+)/(?P<name>[^/]+)", "(?P<project>[^/]+)/(?P<name>[^/]+)", "(?P<name>[^/]+)"}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "projects/{{project}}/zones/{{zone}}/instanceGroupManagers/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
