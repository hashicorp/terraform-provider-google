package google

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	computeBeta "google.golang.org/api/compute/v0.beta"
)

var (
	regionInstanceGroupManagerIdRegex     = regexp.MustCompile("^" + ProjectRegex + "/[a-z0-9-]+/[a-z0-9-]+$")
	regionInstanceGroupManagerIdNameRegex = regexp.MustCompile("^[a-z0-9-]+$")
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
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: compareSelfLinkRelativePaths,
			},

			"version": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Removed:  "This field is in beta. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/provider_versions.html for more details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							Removed:  "This field is in beta. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/provider_versions.html for more details.",
						},

						"instance_template": {
							Type:             schema.TypeString,
							Required:         true,
							Removed:          "This field is in beta. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/provider_versions.html for more details.",
							DiffSuppressFunc: compareSelfLinkRelativePaths,
						},

						"target_size": {
							Type:     schema.TypeList,
							Optional: true,
							Removed:  "This field is in beta. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/provider_versions.html for more details.",
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"fixed": {
										Type:     schema.TypeInt,
										Optional: true,
										Removed:  "This field is in beta. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/provider_versions.html for more details.",
									},

									"percent": {
										Type:         schema.TypeInt,
										Optional:     true,
										Removed:      "This field is in beta. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/provider_versions.html for more details.",
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
				Type:       schema.TypeString,
				Deprecated: "This field is removed.",
				Optional:   true,
				Computed:   true,
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
				Removed:  "This field is in beta. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/provider_versions.html for more details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health_check": {
							Type:             schema.TypeString,
							Required:         true,
							Removed:          "This field is in beta. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/provider_versions.html for more details.",
							DiffSuppressFunc: compareSelfLinkRelativePaths,
						},

						"initial_delay_sec": {
							Type:         schema.TypeInt,
							Required:     true,
							Removed:      "This field is in beta. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/provider_versions.html for more details.",
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

			"rolling_update_policy": {
				Removed:  "This field is in beta. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/provider_versions.html for more details.",
				Computed: true,
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"minimal_action": {
							Type:         schema.TypeString,
							Required:     true,
							Removed:      "This field is in beta. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/provider_versions.html for more details.",
							ValidateFunc: validation.StringInSlice([]string{"RESTART", "REPLACE"}, false),
						},

						"type": {
							Type:         schema.TypeString,
							Required:     true,
							Removed:      "This field is in beta. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/provider_versions.html for more details.",
							ValidateFunc: validation.StringInSlice([]string{"OPPORTUNISTIC", "PROACTIVE"}, false),
						},

						"max_surge_fixed": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							Removed:  "This field is in beta. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/provider_versions.html for more details.",
						},

						"max_surge_percent": {
							Type:         schema.TypeInt,
							Optional:     true,
							Removed:      "This field is in beta. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/provider_versions.html for more details.",
							ValidateFunc: validation.IntBetween(0, 100),
						},

						"max_unavailable_fixed": {
							Type:     schema.TypeInt,
							Optional: true,
							Removed:  "This field is in beta. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/provider_versions.html for more details.",
						},

						"max_unavailable_percent": {
							Type:         schema.TypeInt,
							Optional:     true,
							Removed:      "This field is in beta. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/provider_versions.html for more details.",
							ValidateFunc: validation.IntBetween(0, 100),
						},

						"min_ready_sec": {
							Type:         schema.TypeInt,
							Removed:      "This field is in beta. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/provider_versions.html for more details.",
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

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	manager := &computeBeta.InstanceGroupManager{
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		BaseInstanceName:   d.Get("base_instance_name").(string),
		InstanceTemplate:   d.Get("instance_template").(string),
		TargetSize:         int64(d.Get("target_size").(int)),
		NamedPorts:         getNamedPortsBeta(d.Get("named_port").(*schema.Set).List()),
		TargetPools:        convertStringSet(d.Get("target_pools").(*schema.Set)),
		DistributionPolicy: expandDistributionPolicy(d.Get("distribution_policy_zones").(*schema.Set)),
		// Force send TargetSize to allow size of 0.
		ForceSendFields: []string{"TargetSize"},
	}

	op, err := config.clientComputeBeta.RegionInstanceGroupManagers.Insert(project, region, manager).Do()

	if err != nil {
		return fmt.Errorf("Error creating RegionInstanceGroupManager: %s", err)
	}

	d.SetId(regionInstanceGroupManagerId{Project: project, Region: region, Name: manager.Name}.terraformId())

	// Wait for the operation to complete
	timeoutInMinutes := int(d.Timeout(schema.TimeoutCreate).Minutes())
	err = computeSharedOperationWaitTime(config.clientCompute, op, project, timeoutInMinutes, "Creating InstanceGroupManager")
	if err != nil {
		return err
	}
	return resourceComputeRegionInstanceGroupManagerRead(d, config)
}

type getInstanceManagerFunc func(*schema.ResourceData, interface{}) (*computeBeta.InstanceGroupManager, error)

func getRegionalManager(d *schema.ResourceData, meta interface{}) (*computeBeta.InstanceGroupManager, error) {
	config := meta.(*Config)

	regionalID, err := parseRegionInstanceGroupManagerId(d.Id())
	if err != nil {
		return nil, err
	}

	if regionalID.Project == "" {
		regionalID.Project, err = getProject(d, config)
		if err != nil {
			return nil, err
		}
	}

	if regionalID.Region == "" {
		regionalID.Region, err = getRegion(d, config)
		if err != nil {
			return nil, err
		}
	}

	manager, err := config.clientComputeBeta.RegionInstanceGroupManagers.Get(regionalID.Project, regionalID.Region, regionalID.Name).Do()
	if err != nil {
		return nil, handleNotFoundError(err, d, fmt.Sprintf("Region Instance Manager %q", regionalID.Name))
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

	regionalID, err := parseRegionInstanceGroupManagerId(d.Id())
	if err != nil {
		return err
	}
	if regionalID.Project == "" {
		regionalID.Project, err = getProject(d, config)
		if err != nil {
			return err
		}
	}

	d.Set("base_instance_name", manager.BaseInstanceName)
	d.Set("instance_template", ConvertSelfLinkToV1(manager.InstanceTemplate))

	d.Set("name", manager.Name)
	d.Set("region", GetResourceNameFromSelfLink(manager.Region))
	d.Set("description", manager.Description)
	d.Set("project", regionalID.Project)
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
	// When we make a list Removed, we see a permadiff from `field_name.#: "" => "<computed>"`. Set to nil in Read so we see no diff.
	d.Set("version", nil)
	d.Set("rolling_update_policy", nil)

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

	d.Partial(true)

	if d.HasChange("target_pools") {
		targetPools := convertStringSet(d.Get("target_pools").(*schema.Set))

		// Build the parameter
		setTargetPools := &computeBeta.RegionInstanceGroupManagersSetTargetPoolsRequest{
			Fingerprint: d.Get("fingerprint").(string),
			TargetPools: targetPools,
		}

		op, err := config.clientComputeBeta.RegionInstanceGroupManagers.SetTargetPools(
			project, region, d.Get("name").(string), setTargetPools).Do()

		if err != nil {
			return fmt.Errorf("Error updating RegionInstanceGroupManager: %s", err)
		}

		// Wait for the operation to complete
		timeoutInMinutes := int(d.Timeout(schema.TimeoutUpdate).Minutes())
		err = computeSharedOperationWaitTime(config.clientCompute, op, project, timeoutInMinutes, "Updating RegionInstanceGroupManager")
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
			project, region, d.Get("name").(string), setInstanceTemplate).Do()

		if err != nil {
			return fmt.Errorf("Error updating RegionInstanceGroupManager: %s", err)
		}

		// Wait for the operation to complete
		timeoutInMinutes := int(d.Timeout(schema.TimeoutUpdate).Minutes())
		err = computeSharedOperationWaitTime(config.clientCompute, op, project, timeoutInMinutes, "Updating InstanceGroupManager")
		if err != nil {
			return err
		}

		d.SetPartial("instance_template")
	}

	if d.HasChange("named_port") {
		// Build the parameters for a "SetNamedPorts" request:
		namedPorts := getNamedPortsBeta(d.Get("named_port").(*schema.Set).List())
		setNamedPorts := &computeBeta.RegionInstanceGroupsSetNamedPortsRequest{
			NamedPorts: namedPorts,
		}

		// Make the request:
		op, err := config.clientComputeBeta.RegionInstanceGroups.SetNamedPorts(
			project, region, d.Get("name").(string), setNamedPorts).Do()

		if err != nil {
			return fmt.Errorf("Error updating RegionInstanceGroupManager: %s", err)
		}

		// Wait for the operation to complete:
		timeoutInMinutes := int(d.Timeout(schema.TimeoutUpdate).Minutes())
		err = computeSharedOperationWaitTime(config.clientCompute, op, project, timeoutInMinutes, "Updating RegionInstanceGroupManager")
		if err != nil {
			return err
		}

		d.SetPartial("named_port")
	}

	if d.HasChange("target_size") {
		targetSize := int64(d.Get("target_size").(int))
		op, err := config.clientComputeBeta.RegionInstanceGroupManagers.Resize(
			project, region, d.Get("name").(string), targetSize).Do()

		if err != nil {
			return fmt.Errorf("Error resizing RegionInstanceGroupManager: %s", err)
		}

		// Wait for the operation to complete
		timeoutInMinutes := int(d.Timeout(schema.TimeoutUpdate).Minutes())
		err = computeSharedOperationWaitTime(config.clientCompute, op, project, timeoutInMinutes, "Resizing RegionInstanceGroupManager")
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

	regionalID, err := parseRegionInstanceGroupManagerId(d.Id())
	if err != nil {
		return err
	}

	if regionalID.Project == "" {
		regionalID.Project, err = getProject(d, config)
		if err != nil {
			return err
		}
	}

	if regionalID.Region == "" {
		regionalID.Region, err = getRegion(d, config)
		if err != nil {
			return err
		}
	}

	op, err := config.clientComputeBeta.RegionInstanceGroupManagers.Delete(regionalID.Project, regionalID.Region, regionalID.Name).Do()

	if err != nil {
		return fmt.Errorf("Error deleting region instance group manager: %s", err)
	}

	// Wait for the operation to complete
	timeoutInMinutes := int(d.Timeout(schema.TimeoutDelete).Minutes())
	err = computeSharedOperationWaitTime(config.clientCompute, op, regionalID.Project, timeoutInMinutes, "Deleting RegionInstanceGroupManager")
	if err != nil {
		return fmt.Errorf("Error waiting for delete to complete: %s", err)
	}

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
	regionalID, err := parseRegionInstanceGroupManagerId(d.Id())
	if err != nil {
		return nil, err
	}
	d.Set("project", regionalID.Project)
	d.Set("region", regionalID.Region)
	d.Set("name", regionalID.Name)
	return []*schema.ResourceData{d}, nil
}

type regionInstanceGroupManagerId struct {
	Project string
	Region  string
	Name    string
}

func (r regionInstanceGroupManagerId) terraformId() string {
	return fmt.Sprintf("%s/%s/%s", r.Project, r.Region, r.Name)
}

func parseRegionInstanceGroupManagerId(id string) (*regionInstanceGroupManagerId, error) {
	switch {
	case regionInstanceGroupManagerIdRegex.MatchString(id):
		parts := strings.Split(id, "/")
		return &regionInstanceGroupManagerId{
			Project: parts[0],
			Region:  parts[1],
			Name:    parts[2],
		}, nil
	case regionInstanceGroupManagerIdNameRegex.MatchString(id):
		return &regionInstanceGroupManagerId{
			Name: id,
		}, nil
	default:
		return nil, fmt.Errorf("Invalid region instance group manager specifier. Expecting either {projectId}/{region}/{name} or {name}, where {projectId} and {region} will be derived from the provider.")
	}
}
