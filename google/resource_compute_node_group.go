// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package google

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceComputeNodeGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeNodeGroupCreate,
		Read:   resourceComputeNodeGroupRead,
		Update: resourceComputeNodeGroupUpdate,
		Delete: resourceComputeNodeGroupDelete,

		Importer: &schema.ResourceImporter{
			State: resourceComputeNodeGroupImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"node_template": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      `The URL of the node template to which this node group belongs.`,
			},
			"autoscaling_policy": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				ForceNew: true,
				Description: `If you use sole-tenant nodes for your workloads, you can use the node
group autoscaler to automatically manage the sizes of your node groups.`,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_nodes": {
							Type:     schema.TypeInt,
							Computed: true,
							Optional: true,
							ForceNew: true,
							Description: `Maximum size of the node group. Set to a value less than or equal
to 100 and greater than or equal to min-nodes.`,
						},
						"mode": {
							Type:         schema.TypeString,
							Computed:     true,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validateEnum([]string{"OFF", "ON", "ONLY_SCALE_OUT"}),
							Description: `The autoscaling mode. Set to one of the following:
  - OFF: Disables the autoscaler.
  - ON: Enables scaling in and scaling out.
  - ONLY_SCALE_OUT: Enables only scaling out.
  You must use this mode if your node groups are configured to
  restart their hosted VMs on minimal servers. Possible values: ["OFF", "ON", "ONLY_SCALE_OUT"]`,
						},
						"min_nodes": {
							Type:     schema.TypeInt,
							Computed: true,
							Optional: true,
							ForceNew: true,
							Description: `Minimum size of the node group. Must be less
than or equal to max-nodes. The default value is 0.`,
						},
					},
				},
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `An optional textual description of the resource.`,
			},
			"initial_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Description:  `The initial number of nodes in the node group. One of 'initial_size' or 'size' must be specified.`,
				ExactlyOneOf: []string{"size", "initial_size"},
			},
			"maintenance_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies how to handle instances when a node in the group undergoes maintenance. Set to one of: DEFAULT, RESTART_IN_PLACE, or MIGRATE_WITHIN_NODE_GROUP. The default value is DEFAULT.`,
				Default:     "DEFAULT",
			},
			"maintenance_window": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: `contains properties for the timeframe of maintenance`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: `instances.start time of the window. This must be in UTC format that resolves to one of 00:00, 04:00, 08:00, 12:00, 16:00, or 20:00. For example, both 13:00-5 and 08:00 are valid.`,
						},
					},
				},
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Name of the resource.`,
			},
			"size": {
				Type:         schema.TypeInt,
				Computed:     true,
				Optional:     true,
				ForceNew:     true,
				Description:  `The total number of nodes in the node group. One of 'initial_size' or 'size' must be specified.`,
				ExactlyOneOf: []string{"size", "initial_size"},
			},
			"zone": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      `Zone where this node group is located`,
			},
			"creation_timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Creation timestamp in RFC3339 text format.`,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceComputeNodeGroupCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	descriptionProp, err := expandComputeNodeGroupDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	nameProp, err := expandComputeNodeGroupName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !isEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	nodeTemplateProp, err := expandComputeNodeGroupNodeTemplate(d.Get("node_template"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("node_template"); !isEmptyValue(reflect.ValueOf(nodeTemplateProp)) && (ok || !reflect.DeepEqual(v, nodeTemplateProp)) {
		obj["nodeTemplate"] = nodeTemplateProp
	}
	sizeProp, err := expandComputeNodeGroupSize(d.Get("size"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("size"); ok || !reflect.DeepEqual(v, sizeProp) {
		obj["size"] = sizeProp
	}
	maintenancePolicyProp, err := expandComputeNodeGroupMaintenancePolicy(d.Get("maintenance_policy"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("maintenance_policy"); !isEmptyValue(reflect.ValueOf(maintenancePolicyProp)) && (ok || !reflect.DeepEqual(v, maintenancePolicyProp)) {
		obj["maintenancePolicy"] = maintenancePolicyProp
	}
	maintenanceWindowProp, err := expandComputeNodeGroupMaintenanceWindow(d.Get("maintenance_window"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("maintenance_window"); !isEmptyValue(reflect.ValueOf(maintenanceWindowProp)) && (ok || !reflect.DeepEqual(v, maintenanceWindowProp)) {
		obj["maintenanceWindow"] = maintenanceWindowProp
	}
	autoscalingPolicyProp, err := expandComputeNodeGroupAutoscalingPolicy(d.Get("autoscaling_policy"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("autoscaling_policy"); !isEmptyValue(reflect.ValueOf(autoscalingPolicyProp)) && (ok || !reflect.DeepEqual(v, autoscalingPolicyProp)) {
		obj["autoscalingPolicy"] = autoscalingPolicyProp
	}
	zoneProp, err := expandComputeNodeGroupZone(d.Get("zone"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("zone"); !isEmptyValue(reflect.ValueOf(zoneProp)) && (ok || !reflect.DeepEqual(v, zoneProp)) {
		obj["zone"] = zoneProp
	}

	url, err := replaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/zones/{{zone}}/nodeGroups?initialNodeCount=PRE_CREATE_REPLACE_ME")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new NodeGroup: %#v", obj)
	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for NodeGroup: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	var sizeParam string
	if v, ok := d.GetOkExists("size"); ok {
		sizeParam = fmt.Sprintf("%v", v)
	} else if v, ok := d.GetOkExists("initial_size"); ok {
		sizeParam = fmt.Sprintf("%v", v)
	}

	url = regexp.MustCompile("PRE_CREATE_REPLACE_ME").ReplaceAllLiteralString(url, sizeParam)
	res, err := sendRequestWithTimeout(config, "POST", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating NodeGroup: %s", err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "projects/{{project}}/zones/{{zone}}/nodeGroups/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = computeOperationWaitTime(
		config, res, project, "Creating NodeGroup", userAgent,
		d.Timeout(schema.TimeoutCreate))

	if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create NodeGroup: %s", err)
	}

	log.Printf("[DEBUG] Finished creating NodeGroup %q: %#v", d.Id(), res)

	return resourceComputeNodeGroupRead(d, meta)
}

func resourceComputeNodeGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/zones/{{zone}}/nodeGroups/{{name}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for NodeGroup: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequest(config, "GET", billingProject, url, userAgent, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("ComputeNodeGroup %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading NodeGroup: %s", err)
	}

	if err := d.Set("creation_timestamp", flattenComputeNodeGroupCreationTimestamp(res["creationTimestamp"], d, config)); err != nil {
		return fmt.Errorf("Error reading NodeGroup: %s", err)
	}
	if err := d.Set("description", flattenComputeNodeGroupDescription(res["description"], d, config)); err != nil {
		return fmt.Errorf("Error reading NodeGroup: %s", err)
	}
	if err := d.Set("name", flattenComputeNodeGroupName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading NodeGroup: %s", err)
	}
	if err := d.Set("node_template", flattenComputeNodeGroupNodeTemplate(res["nodeTemplate"], d, config)); err != nil {
		return fmt.Errorf("Error reading NodeGroup: %s", err)
	}
	if err := d.Set("size", flattenComputeNodeGroupSize(res["size"], d, config)); err != nil {
		return fmt.Errorf("Error reading NodeGroup: %s", err)
	}
	if err := d.Set("maintenance_policy", flattenComputeNodeGroupMaintenancePolicy(res["maintenancePolicy"], d, config)); err != nil {
		return fmt.Errorf("Error reading NodeGroup: %s", err)
	}
	if err := d.Set("maintenance_window", flattenComputeNodeGroupMaintenanceWindow(res["maintenanceWindow"], d, config)); err != nil {
		return fmt.Errorf("Error reading NodeGroup: %s", err)
	}
	if err := d.Set("autoscaling_policy", flattenComputeNodeGroupAutoscalingPolicy(res["autoscalingPolicy"], d, config)); err != nil {
		return fmt.Errorf("Error reading NodeGroup: %s", err)
	}
	if err := d.Set("zone", flattenComputeNodeGroupZone(res["zone"], d, config)); err != nil {
		return fmt.Errorf("Error reading NodeGroup: %s", err)
	}
	if err := d.Set("self_link", ConvertSelfLinkToV1(res["selfLink"].(string))); err != nil {
		return fmt.Errorf("Error reading NodeGroup: %s", err)
	}

	return nil
}

func resourceComputeNodeGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for NodeGroup: %s", err)
	}
	billingProject = project

	d.Partial(true)

	if d.HasChange("node_template") {
		obj := make(map[string]interface{})

		nodeTemplateProp, err := expandComputeNodeGroupNodeTemplate(d.Get("node_template"), d, config)
		if err != nil {
			return err
		} else if v, ok := d.GetOkExists("node_template"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, nodeTemplateProp)) {
			obj["nodeTemplate"] = nodeTemplateProp
		}

		url, err := replaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/zones/{{zone}}/nodeGroups/{{name}}/setNodeTemplate")
		if err != nil {
			return err
		}

		// err == nil indicates that the billing_project value was found
		if bp, err := getBillingProject(d, config); err == nil {
			billingProject = bp
		}

		res, err := sendRequestWithTimeout(config, "POST", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("Error updating NodeGroup %q: %s", d.Id(), err)
		} else {
			log.Printf("[DEBUG] Finished updating NodeGroup %q: %#v", d.Id(), res)
		}

		err = computeOperationWaitTime(
			config, res, project, "Updating NodeGroup", userAgent,
			d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceComputeNodeGroupRead(d, meta)
}

func resourceComputeNodeGroupDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for NodeGroup: %s", err)
	}
	billingProject = project

	url, err := replaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/zones/{{zone}}/nodeGroups/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting NodeGroup %q", d.Id())

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "DELETE", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return handleNotFoundError(err, d, "NodeGroup")
	}

	err = computeOperationWaitTime(
		config, res, project, "Deleting NodeGroup", userAgent,
		d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting NodeGroup %q: %#v", d.Id(), res)
	return nil
}

func resourceComputeNodeGroupImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/zones/(?P<zone>[^/]+)/nodeGroups/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<zone>[^/]+)/(?P<name>[^/]+)",
		"(?P<zone>[^/]+)/(?P<name>[^/]+)",
		"(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "projects/{{project}}/zones/{{zone}}/nodeGroups/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenComputeNodeGroupCreationTimestamp(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenComputeNodeGroupDescription(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenComputeNodeGroupName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenComputeNodeGroupNodeTemplate(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return v
	}
	return ConvertSelfLinkToV1(v.(string))
}

func flattenComputeNodeGroupSize(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := stringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeNodeGroupMaintenancePolicy(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenComputeNodeGroupMaintenanceWindow(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["start_time"] =
		flattenComputeNodeGroupMaintenanceWindowStartTime(original["startTime"], d, config)
	return []interface{}{transformed}
}
func flattenComputeNodeGroupMaintenanceWindowStartTime(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenComputeNodeGroupAutoscalingPolicy(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["mode"] =
		flattenComputeNodeGroupAutoscalingPolicyMode(original["mode"], d, config)
	transformed["min_nodes"] =
		flattenComputeNodeGroupAutoscalingPolicyMinNodes(original["minNodes"], d, config)
	transformed["max_nodes"] =
		flattenComputeNodeGroupAutoscalingPolicyMaxNodes(original["maxNodes"], d, config)
	return []interface{}{transformed}
}
func flattenComputeNodeGroupAutoscalingPolicyMode(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenComputeNodeGroupAutoscalingPolicyMinNodes(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := stringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeNodeGroupAutoscalingPolicyMaxNodes(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := stringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeNodeGroupZone(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return v
	}
	return NameFromSelfLinkStateFunc(v)
}

func expandComputeNodeGroupDescription(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeNodeGroupName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeNodeGroupNodeTemplate(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	f, err := parseRegionalFieldValue("nodeTemplates", v.(string), "project", "region", "zone", d, config, true)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for node_template: %s", err)
	}
	return f.RelativeLink(), nil
}

func expandComputeNodeGroupSize(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeNodeGroupMaintenancePolicy(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeNodeGroupMaintenanceWindow(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedStartTime, err := expandComputeNodeGroupMaintenanceWindowStartTime(original["start_time"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedStartTime); val.IsValid() && !isEmptyValue(val) {
		transformed["startTime"] = transformedStartTime
	}

	return transformed, nil
}

func expandComputeNodeGroupMaintenanceWindowStartTime(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeNodeGroupAutoscalingPolicy(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedMode, err := expandComputeNodeGroupAutoscalingPolicyMode(original["mode"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMode); val.IsValid() && !isEmptyValue(val) {
		transformed["mode"] = transformedMode
	}

	transformedMinNodes, err := expandComputeNodeGroupAutoscalingPolicyMinNodes(original["min_nodes"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMinNodes); val.IsValid() && !isEmptyValue(val) {
		transformed["minNodes"] = transformedMinNodes
	}

	transformedMaxNodes, err := expandComputeNodeGroupAutoscalingPolicyMaxNodes(original["max_nodes"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMaxNodes); val.IsValid() && !isEmptyValue(val) {
		transformed["maxNodes"] = transformedMaxNodes
	}

	return transformed, nil
}

func expandComputeNodeGroupAutoscalingPolicyMode(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeNodeGroupAutoscalingPolicyMinNodes(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeNodeGroupAutoscalingPolicyMaxNodes(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeNodeGroupZone(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	f, err := parseGlobalFieldValue("zones", v.(string), "project", d, config, true)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for zone: %s", err)
	}
	return f.RelativeLink(), nil
}
