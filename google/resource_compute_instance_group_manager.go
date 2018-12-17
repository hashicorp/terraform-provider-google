package google

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
)

var (
	instanceGroupManagerIdRegex     = regexp.MustCompile("^" + ProjectRegex + "/[a-z0-9-]+/[a-z0-9-]+$")
	instanceGroupManagerIdNameRegex = regexp.MustCompile("^[a-z0-9-]+$")
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
			"base_instance_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"instance_template": {
				Type:             schema.TypeString,
				Optional:         true,
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

			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

			"wait_for_instances": {
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

	// Build the parameter
	manager := &computeBeta.InstanceGroupManager{
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		BaseInstanceName: d.Get("base_instance_name").(string),
		InstanceTemplate: d.Get("instance_template").(string),
		TargetSize:       int64(d.Get("target_size").(int)),
		NamedPorts:       getNamedPortsBeta(d.Get("named_port").(*schema.Set).List()),
		TargetPools:      convertStringSet(d.Get("target_pools").(*schema.Set)),
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
	d.SetId(instanceGroupManagerId{Project: project, Zone: zone, Name: manager.Name}.terraformId())

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

func getManager(d *schema.ResourceData, meta interface{}) (*computeBeta.InstanceGroupManager, error) {
	config := meta.(*Config)
	zonalID, err := parseInstanceGroupManagerId(d.Id())
	if err != nil {
		return nil, err
	}

	if zonalID.Project == "" {
		project, err := getProject(d, config)
		if err != nil {
			return nil, err
		}
		zonalID.Project = project
	}

	if zonalID.Zone == "" {
		zonalID.Zone, _ = getZone(d, config)
	}

	manager, err := config.clientComputeBeta.InstanceGroupManagers.Get(zonalID.Project, zonalID.Zone, zonalID.Name).Do()
	if err != nil {
		return nil, handleNotFoundError(err, d, fmt.Sprintf("Instance Group Manager %q", zonalID.Name))
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
	if err != nil {
		return err
	}
	if manager == nil {
		log.Printf("[WARN] Instance Group Manager %q not found, removing from state.", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("base_instance_name", manager.BaseInstanceName)
	d.Set("instance_template", ConvertSelfLinkToV1(manager.InstanceTemplate))
	d.Set("name", manager.Name)
	d.Set("zone", GetResourceNameFromSelfLink(manager.Zone))
	d.Set("description", manager.Description)
	d.Set("project", project)
	d.Set("target_size", manager.TargetSize)
	if err = d.Set("target_pools", manager.TargetPools); err != nil {
		return fmt.Errorf("Error setting target_pools in state: %s", err.Error())
	}
	if err = d.Set("named_port", flattenNamedPortsBeta(manager.NamedPorts)); err != nil {
		return fmt.Errorf("Error setting named_port in state: %s", err.Error())
	}
	d.Set("fingerprint", manager.Fingerprint)
	d.Set("instance_group", ConvertSelfLinkToV1(manager.InstanceGroup))
	d.Set("self_link", ConvertSelfLinkToV1(manager.SelfLink))

	update_strategy, ok := d.GetOk("update_strategy")
	if !ok {
		update_strategy = "REPLACE"
	}
	d.Set("update_strategy", update_strategy.(string))

	// When we make a list Removed, we see a permadiff from `field_name.#: "" => "<computed>"`. Set to nil in Read so we see no diff.
	d.Set("version", nil)
	d.Set("rolling_update_policy", nil)

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
func performZoneUpdate(config *Config, id string, updateStrategy string, project string, zone string) error {
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

	return nil
}

func resourceComputeInstanceGroupManagerUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	zonalID, err := parseInstanceGroupManagerId(d.Id())
	if err != nil {
		return err
	}
	if zonalID.Project == "" {
		zonalID.Project, err = getProject(d, config)
		if err != nil {
			return err
		}
	}
	if zonalID.Zone == "" {
		zonalID.Zone, err = getZone(d, config)
		if err != nil {
			return err
		}
	}

	d.Partial(true)

	// If target_pools changes then update
	if d.HasChange("target_pools") {
		targetPools := convertStringSet(d.Get("target_pools").(*schema.Set))

		// Build the parameter
		setTargetPools := &computeBeta.InstanceGroupManagersSetTargetPoolsRequest{
			Fingerprint: d.Get("fingerprint").(string),
			TargetPools: targetPools,
		}

		op, err := config.clientComputeBeta.InstanceGroupManagers.SetTargetPools(
			zonalID.Project, zonalID.Zone, zonalID.Name, setTargetPools).Do()

		if err != nil {
			return fmt.Errorf("Error updating InstanceGroupManager: %s", err)
		}

		// Wait for the operation to complete
		err = computeSharedOperationWait(config.clientCompute, op, zonalID.Project, "Updating InstanceGroupManager")
		if err != nil {
			return err
		}

		d.SetPartial("target_pools")
	}

	// If named_port changes then update:
	if d.HasChange("named_port") {

		// Build the parameters for a "SetNamedPorts" request:
		namedPorts := getNamedPortsBeta(d.Get("named_port").(*schema.Set).List())
		setNamedPorts := &computeBeta.InstanceGroupsSetNamedPortsRequest{
			NamedPorts: namedPorts,
		}

		// Make the request:
		op, err := config.clientComputeBeta.InstanceGroups.SetNamedPorts(
			zonalID.Project, zonalID.Zone, zonalID.Name, setNamedPorts).Do()

		if err != nil {
			return fmt.Errorf("Error updating InstanceGroupManager: %s", err)
		}

		// Wait for the operation to complete:
		err = computeSharedOperationWait(config.clientCompute, op, zonalID.Project, "Updating InstanceGroupManager")
		if err != nil {
			return err
		}

		d.SetPartial("named_port")
	}

	if d.HasChange("target_size") {
		targetSize := int64(d.Get("target_size").(int))
		op, err := config.clientComputeBeta.InstanceGroupManagers.Resize(
			zonalID.Project, zonalID.Zone, zonalID.Name, targetSize).Do()

		if err != nil {
			return fmt.Errorf("Error updating InstanceGroupManager: %s", err)
		}

		// Wait for the operation to complete
		err = computeSharedOperationWait(config.clientCompute, op, zonalID.Project, "Updating InstanceGroupManager")
		if err != nil {
			return err
		}

		d.SetPartial("target_size")
	}

	// If instance_template changes then update
	if d.HasChange("instance_template") {
		// Build the parameter
		setInstanceTemplate := &computeBeta.InstanceGroupManagersSetInstanceTemplateRequest{
			InstanceTemplate: d.Get("instance_template").(string),
		}

		op, err := config.clientComputeBeta.InstanceGroupManagers.SetInstanceTemplate(zonalID.Project, zonalID.Zone, zonalID.Name, setInstanceTemplate).Do()

		if err != nil {
			return fmt.Errorf("Error updating InstanceGroupManager: %s", err)
		}

		// Wait for the operation to complete
		err = computeSharedOperationWait(config.clientCompute, op, zonalID.Project, "Updating InstanceGroupManager")
		if err != nil {
			return err
		}

		updateStrategy := d.Get("update_strategy").(string)
		err = performZoneUpdate(config, zonalID.Name, updateStrategy, zonalID.Project, zonalID.Zone)
		d.SetPartial("instance_template")
	}

	d.Partial(false)

	return resourceComputeInstanceGroupManagerRead(d, meta)
}

func resourceComputeInstanceGroupManagerDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	zonalID, err := parseInstanceGroupManagerId(d.Id())
	if err != nil {
		return err
	}

	if zonalID.Project == "" {
		zonalID.Project, err = getProject(d, config)
		if err != nil {
			return err
		}
	}

	if zonalID.Zone == "" {
		zonalID.Zone, err = getZone(d, config)
		if err != nil {
			return err
		}
	}

	op, err := config.clientComputeBeta.InstanceGroupManagers.Delete(zonalID.Project, zonalID.Zone, zonalID.Name).Do()
	attempt := 0
	for err != nil && attempt < 20 {
		attempt++
		time.Sleep(2000 * time.Millisecond)
		op, err = config.clientComputeBeta.InstanceGroupManagers.Delete(zonalID.Project, zonalID.Zone, zonalID.Name).Do()
	}

	if err != nil {
		return fmt.Errorf("Error deleting instance group manager: %s", err)
	}

	currentSize := int64(d.Get("target_size").(int))

	// Wait for the operation to complete
	err = computeSharedOperationWait(config.clientCompute, op, zonalID.Project, "Deleting InstanceGroupManager")

	for err != nil && currentSize > 0 {
		if !strings.Contains(err.Error(), "timeout") {
			return err
		}

		instanceGroup, err := config.clientComputeBeta.InstanceGroups.Get(
			zonalID.Project, zonalID.Zone, zonalID.Name).Do()
		if err != nil {
			return fmt.Errorf("Error getting instance group size: %s", err)
		}

		instanceGroupSize := instanceGroup.Size

		if instanceGroupSize >= currentSize {
			return fmt.Errorf("Error, instance group isn't shrinking during delete")
		}

		log.Printf("[INFO] timeout occurred, but instance group is shrinking (%d < %d)", instanceGroupSize, currentSize)
		currentSize = instanceGroupSize
		err = computeSharedOperationWait(config.clientCompute, op, zonalID.Project, "Deleting InstanceGroupManager")
	}

	d.SetId("")
	return nil
}

func resourceInstanceGroupManagerStateImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	d.Set("wait_for_instances", false)
	zonalID, err := parseInstanceGroupManagerId(d.Id())
	if err != nil {
		return nil, err
	}
	if zonalID.Zone == "" || zonalID.Project == "" {
		return nil, fmt.Errorf("Invalid instance group manager import ID. Expecting {projectId}/{zone}/{name}.")
	}
	d.Set("project", zonalID.Project)
	d.Set("zone", zonalID.Zone)
	d.Set("name", zonalID.Name)
	return []*schema.ResourceData{d}, nil
}

type instanceGroupManagerId struct {
	Project string
	Zone    string
	Name    string
}

func (i instanceGroupManagerId) terraformId() string {
	return fmt.Sprintf("%s/%s/%s", i.Project, i.Zone, i.Name)
}

func parseInstanceGroupManagerId(id string) (*instanceGroupManagerId, error) {
	switch {
	case instanceGroupManagerIdRegex.MatchString(id):
		parts := strings.Split(id, "/")
		return &instanceGroupManagerId{
			Project: parts[0],
			Zone:    parts[1],
			Name:    parts[2],
		}, nil
	case instanceGroupManagerIdNameRegex.MatchString(id):
		return &instanceGroupManagerId{
			Name: id,
		}, nil
	default:
		return nil, fmt.Errorf("Invalid instance group manager specifier. Expecting either {projectId}/{zone}/{name} or {name}, where {projectId} and {zone} will be derived from the provider.")
	}
}
