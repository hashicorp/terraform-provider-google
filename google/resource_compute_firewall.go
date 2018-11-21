// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
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
	"bytes"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	compute "google.golang.org/api/compute/v1"
)

func resourceComputeFirewallRuleHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["protocol"].(string)))

	// We need to make sure to sort the strings below so that we always
	// generate the same hash code no matter what is in the set.
	if v, ok := m["ports"]; ok {
		s := convertStringArr(v.([]interface{}))
		sort.Strings(s)

		for _, v := range s {
			buf.WriteString(fmt.Sprintf("%s-", v))
		}
	}

	return hashcode.String(buf.String())
}

func resourceComputeFirewall() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeFirewallCreate,
		Read:   resourceComputeFirewallRead,
		Update: resourceComputeFirewallUpdate,
		Delete: resourceComputeFirewallDelete,

		Importer: &schema.ResourceImporter{
			State: resourceComputeFirewallImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(240 * time.Second),
			Update: schema.DefaultTimeout(240 * time.Second),
			Delete: schema.DefaultTimeout(240 * time.Second),
		},
		SchemaVersion: 1,
		MigrateState:  resourceComputeFirewallMigrateState,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateGCPName,
			},
			"network": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},
			"allow": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ports": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
				Set:           resourceComputeFirewallRuleHash,
				ConflictsWith: []string{"deny"},
			},
			"deny": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ports": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
				Set:           resourceComputeFirewallRuleHash,
				ConflictsWith: []string{"allow"},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_ranges": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:           schema.HashString,
				ConflictsWith: []string{"source_ranges", "source_tags"},
			},
			"direction": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"INGRESS", "EGRESS", ""}, false),
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"priority": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 65535),
				Default:      1000,
			},
			"source_ranges": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"source_service_accounts": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:           schema.HashString,
				ConflictsWith: []string{"source_tags", "target_tags"},
			},
			"source_tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"target_service_accounts": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:           schema.HashString,
				ConflictsWith: []string{"source_tags", "target_tags"},
			},
			"target_tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"creation_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
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
	}
}

func resourceComputeFirewallCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := make(map[string]interface{})
	allowedProp, err := expandComputeFirewallAllow(d.Get("allow"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("allow"); !isEmptyValue(reflect.ValueOf(allowedProp)) && (ok || !reflect.DeepEqual(v, allowedProp)) {
		obj["allowed"] = allowedProp
	}
	deniedProp, err := expandComputeFirewallDeny(d.Get("deny"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("deny"); !isEmptyValue(reflect.ValueOf(deniedProp)) && (ok || !reflect.DeepEqual(v, deniedProp)) {
		obj["denied"] = deniedProp
	}
	descriptionProp, err := expandComputeFirewallDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	destinationRangesProp, err := expandComputeFirewallDestinationRanges(d.Get("destination_ranges"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("destination_ranges"); !isEmptyValue(reflect.ValueOf(destinationRangesProp)) && (ok || !reflect.DeepEqual(v, destinationRangesProp)) {
		obj["destinationRanges"] = destinationRangesProp
	}
	directionProp, err := expandComputeFirewallDirection(d.Get("direction"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("direction"); !isEmptyValue(reflect.ValueOf(directionProp)) && (ok || !reflect.DeepEqual(v, directionProp)) {
		obj["direction"] = directionProp
	}
	disabledProp, err := expandComputeFirewallDisabled(d.Get("disabled"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("disabled"); ok || !reflect.DeepEqual(v, disabledProp) {
		obj["disabled"] = disabledProp
	}
	nameProp, err := expandComputeFirewallName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !isEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	networkProp, err := expandComputeFirewallNetwork(d.Get("network"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("network"); !isEmptyValue(reflect.ValueOf(networkProp)) && (ok || !reflect.DeepEqual(v, networkProp)) {
		obj["network"] = networkProp
	}
	priorityProp, err := expandComputeFirewallPriority(d.Get("priority"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("priority"); !isEmptyValue(reflect.ValueOf(priorityProp)) && (ok || !reflect.DeepEqual(v, priorityProp)) {
		obj["priority"] = priorityProp
	}
	sourceRangesProp, err := expandComputeFirewallSourceRanges(d.Get("source_ranges"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("source_ranges"); !isEmptyValue(reflect.ValueOf(sourceRangesProp)) && (ok || !reflect.DeepEqual(v, sourceRangesProp)) {
		obj["sourceRanges"] = sourceRangesProp
	}
	sourceServiceAccountsProp, err := expandComputeFirewallSourceServiceAccounts(d.Get("source_service_accounts"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("source_service_accounts"); !isEmptyValue(reflect.ValueOf(sourceServiceAccountsProp)) && (ok || !reflect.DeepEqual(v, sourceServiceAccountsProp)) {
		obj["sourceServiceAccounts"] = sourceServiceAccountsProp
	}
	sourceTagsProp, err := expandComputeFirewallSourceTags(d.Get("source_tags"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("source_tags"); !isEmptyValue(reflect.ValueOf(sourceTagsProp)) && (ok || !reflect.DeepEqual(v, sourceTagsProp)) {
		obj["sourceTags"] = sourceTagsProp
	}
	targetServiceAccountsProp, err := expandComputeFirewallTargetServiceAccounts(d.Get("target_service_accounts"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("target_service_accounts"); !isEmptyValue(reflect.ValueOf(targetServiceAccountsProp)) && (ok || !reflect.DeepEqual(v, targetServiceAccountsProp)) {
		obj["targetServiceAccounts"] = targetServiceAccountsProp
	}
	targetTagsProp, err := expandComputeFirewallTargetTags(d.Get("target_tags"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("target_tags"); !isEmptyValue(reflect.ValueOf(targetTagsProp)) && (ok || !reflect.DeepEqual(v, targetTagsProp)) {
		obj["targetTags"] = targetTagsProp
	}

	url, err := replaceVars(d, config, "https://www.googleapis.com/compute/v1/projects/{{project}}/global/firewalls")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Firewall: %#v", obj)
	res, err := sendRequest(config, "POST", url, obj)
	if err != nil {
		return fmt.Errorf("Error creating Firewall: %s", err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	op := &compute.Operation{}
	err = Convert(res, op)
	if err != nil {
		return err
	}

	waitErr := computeOperationWaitTime(
		config.clientCompute, op, project, "Creating Firewall",
		int(d.Timeout(schema.TimeoutCreate).Minutes()))

	if waitErr != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create Firewall: %s", waitErr)
	}

	log.Printf("[DEBUG] Finished creating Firewall %q: %#v", d.Id(), res)

	return resourceComputeFirewallRead(d, meta)
}

func resourceComputeFirewallRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	url, err := replaceVars(d, config, "https://www.googleapis.com/compute/v1/projects/{{project}}/global/firewalls/{{name}}")
	if err != nil {
		return err
	}

	res, err := sendRequest(config, "GET", url, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("ComputeFirewall %q", d.Id()))
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading Firewall: %s", err)
	}

	if err := d.Set("allow", flattenComputeFirewallAllow(res["allowed"], d)); err != nil {
		return fmt.Errorf("Error reading Firewall: %s", err)
	}
	if err := d.Set("creation_timestamp", flattenComputeFirewallCreationTimestamp(res["creationTimestamp"], d)); err != nil {
		return fmt.Errorf("Error reading Firewall: %s", err)
	}
	if err := d.Set("deny", flattenComputeFirewallDeny(res["denied"], d)); err != nil {
		return fmt.Errorf("Error reading Firewall: %s", err)
	}
	if err := d.Set("description", flattenComputeFirewallDescription(res["description"], d)); err != nil {
		return fmt.Errorf("Error reading Firewall: %s", err)
	}
	if err := d.Set("destination_ranges", flattenComputeFirewallDestinationRanges(res["destinationRanges"], d)); err != nil {
		return fmt.Errorf("Error reading Firewall: %s", err)
	}
	if err := d.Set("direction", flattenComputeFirewallDirection(res["direction"], d)); err != nil {
		return fmt.Errorf("Error reading Firewall: %s", err)
	}
	if err := d.Set("disabled", flattenComputeFirewallDisabled(res["disabled"], d)); err != nil {
		return fmt.Errorf("Error reading Firewall: %s", err)
	}
	if err := d.Set("name", flattenComputeFirewallName(res["name"], d)); err != nil {
		return fmt.Errorf("Error reading Firewall: %s", err)
	}
	if err := d.Set("network", flattenComputeFirewallNetwork(res["network"], d)); err != nil {
		return fmt.Errorf("Error reading Firewall: %s", err)
	}
	if err := d.Set("priority", flattenComputeFirewallPriority(res["priority"], d)); err != nil {
		return fmt.Errorf("Error reading Firewall: %s", err)
	}
	if err := d.Set("source_ranges", flattenComputeFirewallSourceRanges(res["sourceRanges"], d)); err != nil {
		return fmt.Errorf("Error reading Firewall: %s", err)
	}
	if err := d.Set("source_service_accounts", flattenComputeFirewallSourceServiceAccounts(res["sourceServiceAccounts"], d)); err != nil {
		return fmt.Errorf("Error reading Firewall: %s", err)
	}
	if err := d.Set("source_tags", flattenComputeFirewallSourceTags(res["sourceTags"], d)); err != nil {
		return fmt.Errorf("Error reading Firewall: %s", err)
	}
	if err := d.Set("target_service_accounts", flattenComputeFirewallTargetServiceAccounts(res["targetServiceAccounts"], d)); err != nil {
		return fmt.Errorf("Error reading Firewall: %s", err)
	}
	if err := d.Set("target_tags", flattenComputeFirewallTargetTags(res["targetTags"], d)); err != nil {
		return fmt.Errorf("Error reading Firewall: %s", err)
	}
	if err := d.Set("self_link", ConvertSelfLinkToV1(res["selfLink"].(string))); err != nil {
		return fmt.Errorf("Error reading Firewall: %s", err)
	}

	return nil
}

func resourceComputeFirewallUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := make(map[string]interface{})
	allowedProp, err := expandComputeFirewallAllow(d.Get("allow"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("allow"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, allowedProp)) {
		obj["allowed"] = allowedProp
	}
	deniedProp, err := expandComputeFirewallDeny(d.Get("deny"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("deny"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, deniedProp)) {
		obj["denied"] = deniedProp
	}
	descriptionProp, err := expandComputeFirewallDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	destinationRangesProp, err := expandComputeFirewallDestinationRanges(d.Get("destination_ranges"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("destination_ranges"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, destinationRangesProp)) {
		obj["destinationRanges"] = destinationRangesProp
	}
	directionProp, err := expandComputeFirewallDirection(d.Get("direction"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("direction"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, directionProp)) {
		obj["direction"] = directionProp
	}
	disabledProp, err := expandComputeFirewallDisabled(d.Get("disabled"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("disabled"); ok || !reflect.DeepEqual(v, disabledProp) {
		obj["disabled"] = disabledProp
	}
	networkProp, err := expandComputeFirewallNetwork(d.Get("network"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("network"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, networkProp)) {
		obj["network"] = networkProp
	}
	priorityProp, err := expandComputeFirewallPriority(d.Get("priority"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("priority"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, priorityProp)) {
		obj["priority"] = priorityProp
	}
	sourceRangesProp, err := expandComputeFirewallSourceRanges(d.Get("source_ranges"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("source_ranges"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, sourceRangesProp)) {
		obj["sourceRanges"] = sourceRangesProp
	}
	sourceServiceAccountsProp, err := expandComputeFirewallSourceServiceAccounts(d.Get("source_service_accounts"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("source_service_accounts"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, sourceServiceAccountsProp)) {
		obj["sourceServiceAccounts"] = sourceServiceAccountsProp
	}
	sourceTagsProp, err := expandComputeFirewallSourceTags(d.Get("source_tags"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("source_tags"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, sourceTagsProp)) {
		obj["sourceTags"] = sourceTagsProp
	}
	targetServiceAccountsProp, err := expandComputeFirewallTargetServiceAccounts(d.Get("target_service_accounts"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("target_service_accounts"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, targetServiceAccountsProp)) {
		obj["targetServiceAccounts"] = targetServiceAccountsProp
	}
	targetTagsProp, err := expandComputeFirewallTargetTags(d.Get("target_tags"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("target_tags"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, targetTagsProp)) {
		obj["targetTags"] = targetTagsProp
	}

	url, err := replaceVars(d, config, "https://www.googleapis.com/compute/v1/projects/{{project}}/global/firewalls/{{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating Firewall %q: %#v", d.Id(), obj)
	res, err := sendRequest(config, "PATCH", url, obj)

	if err != nil {
		return fmt.Errorf("Error updating Firewall %q: %s", d.Id(), err)
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	op := &compute.Operation{}
	err = Convert(res, op)
	if err != nil {
		return err
	}

	err = computeOperationWaitTime(
		config.clientCompute, op, project, "Updating Firewall",
		int(d.Timeout(schema.TimeoutUpdate).Minutes()))

	if err != nil {
		return err
	}

	return resourceComputeFirewallRead(d, meta)
}

func resourceComputeFirewallDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	url, err := replaceVars(d, config, "https://www.googleapis.com/compute/v1/projects/{{project}}/global/firewalls/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting Firewall %q", d.Id())
	res, err := sendRequest(config, "DELETE", url, obj)
	if err != nil {
		return handleNotFoundError(err, d, "Firewall")
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	op := &compute.Operation{}
	err = Convert(res, op)
	if err != nil {
		return err
	}

	err = computeOperationWaitTime(
		config.clientCompute, op, project, "Deleting Firewall",
		int(d.Timeout(schema.TimeoutDelete).Minutes()))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting Firewall %q: %#v", d.Id(), res)
	return nil
}

func resourceComputeFirewallImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	parseImportId([]string{"projects/(?P<project>[^/]+)/global/firewalls/(?P<name>[^/]+)", "(?P<project>[^/]+)/(?P<name>[^/]+)", "(?P<name>[^/]+)"}, d, config)

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenComputeFirewallAllow(v interface{}, d *schema.ResourceData) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"protocol": flattenComputeFirewallAllowProtocol(original["IPProtocol"], d),
			"ports":    flattenComputeFirewallAllowPorts(original["ports"], d),
		})
	}
	return transformed
}
func flattenComputeFirewallAllowProtocol(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenComputeFirewallAllowPorts(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenComputeFirewallCreationTimestamp(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenComputeFirewallDeny(v interface{}, d *schema.ResourceData) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"protocol": flattenComputeFirewallDenyProtocol(original["IPProtocol"], d),
			"ports":    flattenComputeFirewallDenyPorts(original["ports"], d),
		})
	}
	return transformed
}
func flattenComputeFirewallDenyProtocol(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenComputeFirewallDenyPorts(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenComputeFirewallDescription(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenComputeFirewallDestinationRanges(v interface{}, d *schema.ResourceData) interface{} {
	if v == nil {
		return v
	}
	return schema.NewSet(schema.HashString, v.([]interface{}))
}

func flattenComputeFirewallDirection(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenComputeFirewallDisabled(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenComputeFirewallName(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenComputeFirewallNetwork(v interface{}, d *schema.ResourceData) interface{} {
	if v == nil {
		return v
	}
	return ConvertSelfLinkToV1(v.(string))
}

func flattenComputeFirewallPriority(v interface{}, d *schema.ResourceData) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := strconv.ParseInt(strVal, 10, 64); err == nil {
			return intVal
		} // let terraform core handle it if we can't convert the string to an int.
	}
	return v
}

func flattenComputeFirewallSourceRanges(v interface{}, d *schema.ResourceData) interface{} {
	if v == nil {
		return v
	}
	return schema.NewSet(schema.HashString, v.([]interface{}))
}

func flattenComputeFirewallSourceServiceAccounts(v interface{}, d *schema.ResourceData) interface{} {
	if v == nil {
		return v
	}
	return schema.NewSet(schema.HashString, v.([]interface{}))
}

func flattenComputeFirewallSourceTags(v interface{}, d *schema.ResourceData) interface{} {
	if v == nil {
		return v
	}
	return schema.NewSet(schema.HashString, v.([]interface{}))
}

func flattenComputeFirewallTargetServiceAccounts(v interface{}, d *schema.ResourceData) interface{} {
	if v == nil {
		return v
	}
	return schema.NewSet(schema.HashString, v.([]interface{}))
}

func flattenComputeFirewallTargetTags(v interface{}, d *schema.ResourceData) interface{} {
	if v == nil {
		return v
	}
	return schema.NewSet(schema.HashString, v.([]interface{}))
}

func expandComputeFirewallAllow(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedProtocol, err := expandComputeFirewallAllowProtocol(original["protocol"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedProtocol); val.IsValid() && !isEmptyValue(val) {
			transformed["IPProtocol"] = transformedProtocol
		}

		transformedPorts, err := expandComputeFirewallAllowPorts(original["ports"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedPorts); val.IsValid() && !isEmptyValue(val) {
			transformed["ports"] = transformedPorts
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandComputeFirewallAllowProtocol(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeFirewallAllowPorts(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeFirewallDeny(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedProtocol, err := expandComputeFirewallDenyProtocol(original["protocol"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedProtocol); val.IsValid() && !isEmptyValue(val) {
			transformed["IPProtocol"] = transformedProtocol
		}

		transformedPorts, err := expandComputeFirewallDenyPorts(original["ports"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedPorts); val.IsValid() && !isEmptyValue(val) {
			transformed["ports"] = transformedPorts
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandComputeFirewallDenyProtocol(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeFirewallDenyPorts(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeFirewallDescription(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeFirewallDestinationRanges(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	return v, nil
}

func expandComputeFirewallDirection(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeFirewallDisabled(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeFirewallName(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeFirewallNetwork(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	f, err := parseGlobalFieldValue("networks", v.(string), "project", d, config, true)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for network: %s", err)
	}
	return f.RelativeLink(), nil
}

func expandComputeFirewallPriority(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeFirewallSourceRanges(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	return v, nil
}

func expandComputeFirewallSourceServiceAccounts(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	return v, nil
}

func expandComputeFirewallSourceTags(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	return v, nil
}

func expandComputeFirewallTargetServiceAccounts(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	return v, nil
}

func expandComputeFirewallTargetTags(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	return v, nil
}
