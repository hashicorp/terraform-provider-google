package google

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

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
				Type:       schema.TypeList,
				Optional:   true,
				Computed:   true,
				Deprecated: "This field is in beta and will be removed from this provider. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/provider_versions.html for more details.",
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

			"region": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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
				Default:      "NONE",
				Deprecated:   "This field will have no functionality in 2.0.0, and will be removed. If you're using ROLLING_UPDATE, use the google-beta provider. See https://terraform.io/docs/providers/google/provider_versions.html for more details.",
				ValidateFunc: validation.StringInSlice([]string{"NONE", "ROLLING_UPDATE"}, false),
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

			// If true, the resource will report ready only after no instances are being created.
			// This will not block future reads if instances are being recreated, and it respects
			// the "createNoRetry" parameter that's available for this resource.
			"wait_for_instances": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"auto_healing_policies": &schema.Schema{
				Type:       schema.TypeList,
				Optional:   true,
				MaxItems:   1,
				Deprecated: "This field is in beta and will be removed from this provider. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/provider_versions.html for more details.",
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

			"distribution_policy_zones": &schema.Schema{
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
							Default:       0,
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
							Default:       0,
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
		},
	}
}

func resourceComputeRegionInstanceGroupManagerCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	if _, ok := d.GetOk("rolling_update_policy"); d.Get("update_strategy") == "ROLLING_UPDATE" && !ok {
		return fmt.Errorf("[rolling_update_policy] must be set when 'update_strategy' is set to 'ROLLING_UPDATE'")
	}

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
		DistributionPolicy:  expandDistributionPolicy(d.Get("distribution_policy_zones").(*schema.Set)),
		// Force send TargetSize to allow size of 0.
		ForceSendFields: []string{"TargetSize"},
	}

	op, err := config.clientComputeBeta.RegionInstanceGroupManagers.Insert(project, d.Get("region").(string), manager).Do()

	if err != nil {
		return fmt.Errorf("Error creating RegionInstanceGroupManager: %s", err)
	}

	d.SetId(manager.Name)

	// Wait for the operation to complete
	err = computeSharedOperationWait(config.clientCompute, op, project, "Creating InstanceGroupManager")
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

	manager, err := config.clientComputeBeta.RegionInstanceGroupManagers.Get(project, region, d.Id()).Do()
	if err != nil {
		return nil, handleNotFoundError(err, d, fmt.Sprintf("Region Instance Manager %q", d.Get("name").(string)))
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
		if done := m.CurrentActions.None; done < m.TargetSize {
			return done, "creating", nil
		} else {
			return done, "created", nil
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
	d.Set("instance_template", ConvertSelfLinkToV1(manager.InstanceTemplate))
	if err := d.Set("version", flattenVersions(manager.Versions)); err != nil {
		return err
	}
	d.Set("name", manager.Name)
	d.Set("region", GetResourceNameFromSelfLink(manager.Region))
	d.Set("description", manager.Description)
	d.Set("project", project)
	d.Set("target_size", manager.TargetSize)
	if err := d.Set("target_pools", manager.TargetPools); err != nil {
		return fmt.Errorf("Error setting target_pools in state: %s", err.Error())
	}
	if err := d.Set("named_port", flattenNamedPortsBeta(manager.NamedPorts)); err != nil {
		return fmt.Errorf("Error setting named_port in state: %s", err.Error())
	}
	d.Set("fingerprint", manager.Fingerprint)
	d.Set("instance_group", ConvertSelfLinkToV1(manager.InstanceGroup))
	if err := d.Set("auto_healing_policies", flattenAutoHealingPolicies(manager.AutoHealingPolicies)); err != nil {
		return fmt.Errorf("Error setting auto_healing_policies in state: %s", err.Error())
	}
	if err := d.Set("distribution_policy_zones", flattenDistributionPolicy(manager.DistributionPolicy)); err != nil {
		return err
	}
	d.Set("self_link", ConvertSelfLinkToV1(manager.SelfLink))
	update_strategy, ok := d.GetOk("update_strategy")
	if !ok {
		update_strategy = "NONE"
	}
	d.Set("update_strategy", update_strategy.(string))

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

// Updates an instance group manager by applying the update strategy (REPLACE, RESTART)
// and rolling update policy (PROACTIVE, OPPORTUNISTIC). Updates performed by API
// are OPPORTUNISTIC by default.
func performRegionUpdate(config *Config, id string, updateStrategy string, rollingUpdatePolicy *computeBeta.InstanceGroupManagerUpdatePolicy, versions []*computeBeta.InstanceGroupManagerVersion, project string, region string) error {
	if updateStrategy == "RESTART" {
		managedInstances, err := config.clientComputeBeta.RegionInstanceGroupManagers.ListManagedInstances(project, region, id).Do()
		if err != nil {
			return fmt.Errorf("Error getting region instance group managers instances: %s", err)
		}

		managedInstanceCount := len(managedInstances.ManagedInstances)
		instances := make([]string, managedInstanceCount)
		for i, v := range managedInstances.ManagedInstances {
			instances[i] = v.Instance
		}

		recreateInstances := &computeBeta.RegionInstanceGroupManagersRecreateRequest{
			Instances: instances,
		}

		op, err := config.clientComputeBeta.RegionInstanceGroupManagers.RecreateInstances(project, region, id, recreateInstances).Do()
		if err != nil {
			return fmt.Errorf("Error restarting region instance group managers instances: %s", err)
		}

		// Wait for the operation to complete
		err = computeSharedOperationWaitTime(config.clientCompute, op, project, managedInstanceCount*4, "Restarting RegionInstanceGroupManagers instances")
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

		op, err := config.clientComputeBeta.RegionInstanceGroupManagers.Patch(project, region, id, manager).Do()
		if err != nil {
			return fmt.Errorf("Error updating region managed group instances: %s", err)
		}

		err = computeSharedOperationWait(config.clientCompute, op, project, "Updating region managed group instances")
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

	region := d.Get("region").(string)

	d.Partial(true)

	if _, ok := d.GetOk("rolling_update_policy"); d.Get("update_strategy") == "ROLLING_UPDATE" && !ok {
		return fmt.Errorf("[rolling_update_policy] must be set when 'update_strategy' is set to 'ROLLING_UPDATE'")
	}

	if d.HasChange("target_pools") {
		targetPools := convertStringSet(d.Get("target_pools").(*schema.Set))

		// Build the parameter
		setTargetPools := &computeBeta.RegionInstanceGroupManagersSetTargetPoolsRequest{
			Fingerprint: d.Get("fingerprint").(string),
			TargetPools: targetPools,
		}

		op, err := config.clientComputeBeta.RegionInstanceGroupManagers.SetTargetPools(
			project, region, d.Id(), setTargetPools).Do()

		if err != nil {
			return fmt.Errorf("Error updating RegionInstanceGroupManager: %s", err)
		}

		// Wait for the operation to complete
		err = computeSharedOperationWait(config.clientCompute, op, project, "Updating RegionInstanceGroupManager")
		if err != nil {
			return err
		}

		d.SetPartial("target_pools")
	}

	if d.HasChange("instance_template") {
		// Build the parameter
		setInstanceTemplate := &computeBeta.RegionInstanceGroupManagersSetTemplateRequest{
			InstanceTemplate: d.Get("instance_template").(string),
		}

		op, err := config.clientComputeBeta.RegionInstanceGroupManagers.SetInstanceTemplate(
			project, region, d.Id(), setInstanceTemplate).Do()

		if err != nil {
			return fmt.Errorf("Error updating RegionInstanceGroupManager: %s", err)
		}

		// Wait for the operation to complete
		err = computeSharedOperationWait(config.clientCompute, op, project, "Updating InstanceGroupManager")
		if err != nil {
			return err
		}

		updateStrategy := d.Get("update_strategy").(string)
		rollingUpdatePolicy := expandUpdatePolicy(d.Get("rolling_update_policy").([]interface{}))
		err = performRegionUpdate(config, d.Id(), updateStrategy, rollingUpdatePolicy, nil, project, region)
		d.SetPartial("instance_template")
	}

	// If version changes then update
	if d.HasChange("version") {
		updateStrategy := d.Get("update_strategy").(string)
		rollingUpdatePolicy := expandUpdatePolicy(d.Get("rolling_update_policy").([]interface{}))
		versions := expandVersions(d.Get("version").([]interface{}))
		err = performRegionUpdate(config, d.Id(), updateStrategy, rollingUpdatePolicy, versions, project, region)
		if err != nil {
			return err
		}

		d.SetPartial("version")
	}

	if d.HasChange("named_port") {
		// Build the parameters for a "SetNamedPorts" request:
		namedPorts := getNamedPortsBeta(d.Get("named_port").([]interface{}))
		setNamedPorts := &computeBeta.RegionInstanceGroupsSetNamedPortsRequest{
			NamedPorts: namedPorts,
		}

		// Make the request:
		op, err := config.clientComputeBeta.RegionInstanceGroups.SetNamedPorts(
			project, region, d.Id(), setNamedPorts).Do()

		if err != nil {
			return fmt.Errorf("Error updating RegionInstanceGroupManager: %s", err)
		}

		// Wait for the operation to complete:
		err = computeSharedOperationWait(config.clientCompute, op, project, "Updating RegionInstanceGroupManager")
		if err != nil {
			return err
		}

		d.SetPartial("named_port")
	}

	if d.HasChange("target_size") {
		targetSize := int64(d.Get("target_size").(int))
		op, err := config.clientComputeBeta.RegionInstanceGroupManagers.Resize(
			project, region, d.Id(), targetSize).Do()

		if err != nil {
			return fmt.Errorf("Error resizing RegionInstanceGroupManager: %s", err)
		}

		// Wait for the operation to complete
		err = computeSharedOperationWait(config.clientCompute, op, project, "Resizing RegionInstanceGroupManager")
		if err != nil {
			return err
		}

		d.SetPartial("target_size")
	}

	if d.HasChange("auto_healing_policies") {
		setAutoHealingPoliciesRequest := &computeBeta.RegionInstanceGroupManagersSetAutoHealingRequest{}
		if v, ok := d.GetOk("auto_healing_policies"); ok {
			setAutoHealingPoliciesRequest.AutoHealingPolicies = expandAutoHealingPolicies(v.([]interface{}))
		}

		op, err := config.clientComputeBeta.RegionInstanceGroupManagers.SetAutoHealingPolicies(
			project, region, d.Id(), setAutoHealingPoliciesRequest).Do()

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

	d.Partial(false)

	return resourceComputeRegionInstanceGroupManagerRead(d, meta)
}

func resourceComputeRegionInstanceGroupManagerDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region := d.Get("region").(string)

	op, err := config.clientComputeBeta.RegionInstanceGroupManagers.Delete(project, region, d.Id()).Do()

	if err != nil {
		return fmt.Errorf("Error deleting region instance group manager: %s", err)
	}

	// Wait for the operation to complete
	err = computeSharedOperationWaitTime(config.clientCompute, op, project, int(d.Timeout(schema.TimeoutDelete).Minutes()), "Deleting RegionInstanceGroupManager")

	d.SetId("")
	return nil
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
			zones = append(zones, zone.Zone)
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
	return []*schema.ResourceData{d}, nil
}
