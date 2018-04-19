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
	"google.golang.org/api/compute/v1"
)

var RegionInstanceGroupManagerBaseApiVersion = v1
var RegionInstanceGroupManagerVersionedFeatures = []Feature{
	Feature{Version: v0beta, Item: "auto_healing_policies"},
	Feature{Version: v0beta, Item: "distribution_policy_zones"},
	Feature{Version: v0beta, Item: "rolling_update_policy"},
}

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
				Required:         true,
				DiffSuppressFunc: compareSelfLinkRelativePaths,
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
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
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
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
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
	computeApiVersion := getComputeApiVersion(d, RegionInstanceGroupManagerBaseApiVersion, RegionInstanceGroupManagerVersionedFeatures)
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
		DistributionPolicy:  expandDistributionPolicy(d.Get("distribution_policy_zones").(*schema.Set)),
		// Force send TargetSize to allow size of 0.
		ForceSendFields: []string{"TargetSize"},
	}

	var op interface{}
	switch computeApiVersion {
	case v1:
		managerV1 := &compute.InstanceGroupManager{}
		err = Convert(manager, managerV1)
		if err != nil {
			return err
		}
		managerV1.ForceSendFields = manager.ForceSendFields
		op, err = config.clientCompute.RegionInstanceGroupManagers.Insert(project, d.Get("region").(string), managerV1).Do()
	case v0beta:
		op, err = config.clientComputeBeta.RegionInstanceGroupManagers.Insert(project, d.Get("region").(string), manager).Do()
	}

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
	computeApiVersion := getComputeApiVersion(d, RegionInstanceGroupManagerBaseApiVersion, RegionInstanceGroupManagerVersionedFeatures)
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	region := d.Get("region").(string)
	manager := &computeBeta.InstanceGroupManager{}

	switch computeApiVersion {
	case v1:
		v1Manager := &compute.InstanceGroupManager{}
		v1Manager, err = config.clientCompute.RegionInstanceGroupManagers.Get(project, region, d.Id()).Do()

		if v1Manager == nil {
			log.Printf("[WARN] Removing Region Instance Group Manager %q because it's gone", d.Get("name").(string))

			// The resource doesn't exist anymore
			d.SetId("")
			return nil, nil
		}

		err = Convert(v1Manager, manager)
		if err != nil {
			return nil, err
		}
	case v0beta:
		manager, err = config.clientComputeBeta.RegionInstanceGroupManagers.Get(project, region, d.Id()).Do()
	}

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
	if err != nil || manager == nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	d.Set("base_instance_name", manager.BaseInstanceName)
	d.Set("instance_template", manager.InstanceTemplate)
	d.Set("name", manager.Name)
	d.Set("region", GetResourceNameFromSelfLink(manager.Region))
	d.Set("description", manager.Description)
	d.Set("project", project)
	d.Set("target_size", manager.TargetSize)
	d.Set("target_pools", manager.TargetPools)
	d.Set("named_port", flattenNamedPortsBeta(manager.NamedPorts))
	d.Set("fingerprint", manager.Fingerprint)
	d.Set("instance_group", manager.InstanceGroup)
	d.Set("auto_healing_policies", flattenAutoHealingPolicies(manager.AutoHealingPolicies))
	if err := d.Set("distribution_policy_zones", flattenDistributionPolicy(manager.DistributionPolicy)); err != nil {
		return err
	}
	d.Set("self_link", ConvertSelfLinkToV1(manager.SelfLink))

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
	computeApiVersion := getComputeApiVersionUpdate(d, RegionInstanceGroupManagerBaseApiVersion, RegionInstanceGroupManagerVersionedFeatures, []Feature{})
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

		var op interface{}
		switch computeApiVersion {
		case v1:
			setTargetPoolsV1 := &compute.RegionInstanceGroupManagersSetTargetPoolsRequest{}
			err = Convert(setTargetPools, setTargetPoolsV1)
			if err != nil {
				return err
			}

			op, err = config.clientCompute.RegionInstanceGroupManagers.SetTargetPools(
				project, region, d.Id(), setTargetPoolsV1).Do()
		case v0beta:
			setTargetPoolsV0beta := &computeBeta.RegionInstanceGroupManagersSetTargetPoolsRequest{}
			err = Convert(setTargetPools, setTargetPoolsV0beta)
			if err != nil {
				return err
			}

			op, err = config.clientComputeBeta.RegionInstanceGroupManagers.SetTargetPools(
				project, region, d.Id(), setTargetPoolsV0beta).Do()
		}

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

		var op interface{}
		switch computeApiVersion {
		case v1:
			setInstanceTemplateV1 := &compute.RegionInstanceGroupManagersSetTemplateRequest{}
			err = Convert(setInstanceTemplate, setInstanceTemplateV1)
			if err != nil {
				return err
			}

			op, err = config.clientCompute.RegionInstanceGroupManagers.SetInstanceTemplate(
				project, region, d.Id(), setInstanceTemplateV1).Do()
		case v0beta:
			setInstanceTemplateV0beta := &computeBeta.RegionInstanceGroupManagersSetTemplateRequest{}
			err = Convert(setInstanceTemplate, setInstanceTemplateV0beta)
			if err != nil {
				return err
			}

			op, err = config.clientComputeBeta.RegionInstanceGroupManagers.SetInstanceTemplate(
				project, region, d.Id(), setInstanceTemplateV0beta).Do()
		}

		if err != nil {
			return fmt.Errorf("Error updating RegionInstanceGroupManager: %s", err)
		}

		// Wait for the operation to complete
		err = computeSharedOperationWait(config.clientCompute, op, project, "Updating InstanceGroupManager")
		if err != nil {
			return err
		}

		if d.Get("update_strategy").(string) == "ROLLING_UPDATE" {
			// UpdatePolicy is set for InstanceGroupManager on update only, because it is only relevant for `Patch` calls.
			// Other tools(gcloud and UI) capable of executing the same `ROLLING UPDATE` call
			// expect those values to be provided by user as part of the call
			// or provide their own defaults without respecting what was previously set on UpdateManager.
			// To follow the same logic, we provide policy values on relevant update change only.
			manager := &computeBeta.InstanceGroupManager{
				UpdatePolicy: expandUpdatePolicy(d.Get("rolling_update_policy").([]interface{})),
			}

			op, err = config.clientComputeBeta.RegionInstanceGroupManagers.Patch(
				project, region, d.Id(), manager).Do()
			if err != nil {
				return fmt.Errorf("Error updating managed group instances: %s", err)
			}

			err = computeSharedOperationWait(config.clientCompute, op, project, "Updating managed group instances")
			if err != nil {
				return err
			}
		}

		d.SetPartial("instance_template")
	}

	if d.HasChange("named_port") {
		// Build the parameters for a "SetNamedPorts" request:
		namedPorts := getNamedPortsBeta(d.Get("named_port").([]interface{}))
		setNamedPorts := &computeBeta.RegionInstanceGroupsSetNamedPortsRequest{
			NamedPorts: namedPorts,
		}

		// Make the request:
		var op interface{}
		switch computeApiVersion {
		case v1:
			setNamedPortsV1 := &compute.RegionInstanceGroupsSetNamedPortsRequest{}
			err = Convert(setNamedPorts, setNamedPortsV1)
			if err != nil {
				return err
			}

			op, err = config.clientCompute.RegionInstanceGroups.SetNamedPorts(
				project, region, d.Id(), setNamedPortsV1).Do()
		case v0beta:
			setNamedPortsV0beta := &computeBeta.RegionInstanceGroupsSetNamedPortsRequest{}
			err = Convert(setNamedPorts, setNamedPortsV0beta)
			if err != nil {
				return err
			}

			op, err = config.clientComputeBeta.RegionInstanceGroups.SetNamedPorts(
				project, region, d.Id(), setNamedPortsV0beta).Do()
		}

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
		var op interface{}
		switch computeApiVersion {
		case v1:
			op, err = config.clientCompute.RegionInstanceGroupManagers.Resize(
				project, region, d.Id(), targetSize).Do()
		case v0beta:
			op, err = config.clientComputeBeta.RegionInstanceGroupManagers.Resize(
				project, region, d.Id(), targetSize).Do()
		}

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
	computeApiVersion := getComputeApiVersion(d, RegionInstanceGroupManagerBaseApiVersion, RegionInstanceGroupManagerVersionedFeatures)
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region := d.Get("region").(string)

	var op interface{}
	switch computeApiVersion {
	case v1:
		op, err = config.clientCompute.RegionInstanceGroupManagers.Delete(project, region, d.Id()).Do()
	case v0beta:
		op, err = config.clientComputeBeta.RegionInstanceGroupManagers.Delete(project, region, d.Id()).Do()
	}

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
