// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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

package networkconnectivity

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"
)

func ResourceNetworkConnectivityHub() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkConnectivityHubCreate,
		Read:   resourceNetworkConnectivityHubRead,
		Update: resourceNetworkConnectivityHubUpdate,
		Delete: resourceNetworkConnectivityHubDelete,

		Importer: &schema.ResourceImporter{
			State: resourceNetworkConnectivityHubImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			tpgresource.SetLabelsDiff,
			tpgresource.DefaultProviderProject,
		),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      `Immutable. The name of the hub. Hub names must be unique. They use the following form: 'projects/{project_number}/locations/global/hubs/{hub_id}'`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `An optional description of the hub.`,
			},
			"export_psc": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: `Whether Private Service Connect transitivity is enabled for the hub. If true, Private Service Connect endpoints in VPC spokes attached to the hub are made accessible to other VPC spokes attached to the hub. The default value is false.`,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Description: `Optional labels in key:value format. For more information about labels, see [Requirements for labels](https://cloud.google.com/resource-manager/docs/creating-managing-labels#requirements).

**Note**: This field is non-authoritative, and will only manage the labels present in your configuration.
Please refer to the field 'effective_labels' for all of the labels present on the resource.`,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"preset_topology": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidateEnum([]string{"MESH", "STAR", ""}),
				Description:  `Optional. The topology implemented in this hub. Currently, this field is only used when policyMode = PRESET. The available preset topologies are MESH and STAR. If presetTopology is unspecified and policyMode = PRESET, the presetTopology defaults to MESH. When policyMode = CUSTOM, the presetTopology is set to PRESET_TOPOLOGY_UNSPECIFIED. Possible values: ["MESH", "STAR"]`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Output only. The time the hub was created.`,
			},
			"effective_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: `All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"routing_vpcs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The VPC network associated with this hub's spokes. All of the VPN tunnels, VLAN attachments, and router appliance instances referenced by this hub's spokes must belong to this VPC network. This field is read-only. Network Connectivity Center automatically populates it based on the set of spokes attached to the hub.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uri": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The URI of the VPC network.`,
						},
					},
				},
			},
			"state": {
				Type:         schema.TypeString,
				Computed:     true,
				Description:  `Output only. The current lifecycle state of this hub.`,
				ExactlyOneOf: []string{},
			},
			"terraform_labels": {
				Type:     schema.TypeMap,
				Computed: true,
				Description: `The combination of labels configured directly on the resource
 and default labels configured on the provider.`,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"unique_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Output only. The Google-generated UUID for the hub. This value is unique across all hub resources. If a hub is deleted and another with the same name is created, the new hub is assigned a different unique_id.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Output only. The time the hub was last updated.`,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceNetworkConnectivityHubCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	nameProp, err := expandNetworkConnectivityHubName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	descriptionProp, err := expandNetworkConnectivityHubDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	presetTopologyProp, err := expandNetworkConnectivityHubPresetTopology(d.Get("preset_topology"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("preset_topology"); !tpgresource.IsEmptyValue(reflect.ValueOf(presetTopologyProp)) && (ok || !reflect.DeepEqual(v, presetTopologyProp)) {
		obj["presetTopology"] = presetTopologyProp
	}
	exportPscProp, err := expandNetworkConnectivityHubExportPsc(d.Get("export_psc"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("export_psc"); !tpgresource.IsEmptyValue(reflect.ValueOf(exportPscProp)) && (ok || !reflect.DeepEqual(v, exportPscProp)) {
		obj["exportPsc"] = exportPscProp
	}
	labelsProp, err := expandNetworkConnectivityHubEffectiveLabels(d.Get("effective_labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("effective_labels"); !tpgresource.IsEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}

	url, err := tpgresource.ReplaceVarsForId(d, config, "{{NetworkConnectivityBasePath}}projects/{{project}}/locations/global/hubs?hubId={{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Hub: %#v", obj)
	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Hub: %s", err)
	}
	billingProject = strings.TrimPrefix(project, "projects/")

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutCreate),
		Headers:   headers,
	})
	if err != nil {
		return fmt.Errorf("Error creating Hub: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVarsForId(d, config, "projects/{{project}}/locations/global/hubs/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = NetworkConnectivityOperationWaitTime(
		config, res, tpgresource.GetResourceNameFromSelfLink(project), "Creating Hub", userAgent,
		d.Timeout(schema.TimeoutCreate))

	if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create Hub: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Hub %q: %#v", d.Id(), res)

	return resourceNetworkConnectivityHubRead(d, meta)
}

func resourceNetworkConnectivityHubRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVarsForId(d, config, "{{NetworkConnectivityBasePath}}projects/{{project}}/locations/global/hubs/{{name}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Hub: %s", err)
	}
	billingProject = strings.TrimPrefix(project, "projects/")

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Headers:   headers,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("NetworkConnectivityHub %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading Hub: %s", err)
	}

	if err := d.Set("name", flattenNetworkConnectivityHubName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading Hub: %s", err)
	}
	if err := d.Set("create_time", flattenNetworkConnectivityHubCreateTime(res["createTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading Hub: %s", err)
	}
	if err := d.Set("update_time", flattenNetworkConnectivityHubUpdateTime(res["updateTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading Hub: %s", err)
	}
	if err := d.Set("labels", flattenNetworkConnectivityHubLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading Hub: %s", err)
	}
	if err := d.Set("description", flattenNetworkConnectivityHubDescription(res["description"], d, config)); err != nil {
		return fmt.Errorf("Error reading Hub: %s", err)
	}
	if err := d.Set("unique_id", flattenNetworkConnectivityHubUniqueId(res["uniqueId"], d, config)); err != nil {
		return fmt.Errorf("Error reading Hub: %s", err)
	}
	if err := d.Set("state", flattenNetworkConnectivityHubState(res["state"], d, config)); err != nil {
		return fmt.Errorf("Error reading Hub: %s", err)
	}
	if err := d.Set("routing_vpcs", flattenNetworkConnectivityHubRoutingVpcs(res["routingVpcs"], d, config)); err != nil {
		return fmt.Errorf("Error reading Hub: %s", err)
	}
	if err := d.Set("preset_topology", flattenNetworkConnectivityHubPresetTopology(res["presetTopology"], d, config)); err != nil {
		return fmt.Errorf("Error reading Hub: %s", err)
	}
	if err := d.Set("export_psc", flattenNetworkConnectivityHubExportPsc(res["exportPsc"], d, config)); err != nil {
		return fmt.Errorf("Error reading Hub: %s", err)
	}
	if err := d.Set("terraform_labels", flattenNetworkConnectivityHubTerraformLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading Hub: %s", err)
	}
	if err := d.Set("effective_labels", flattenNetworkConnectivityHubEffectiveLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading Hub: %s", err)
	}

	return nil
}

func resourceNetworkConnectivityHubUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Hub: %s", err)
	}
	billingProject = strings.TrimPrefix(project, "projects/")

	obj := make(map[string]interface{})
	descriptionProp, err := expandNetworkConnectivityHubDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	exportPscProp, err := expandNetworkConnectivityHubExportPsc(d.Get("export_psc"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("export_psc"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, exportPscProp)) {
		obj["exportPsc"] = exportPscProp
	}
	labelsProp, err := expandNetworkConnectivityHubEffectiveLabels(d.Get("effective_labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("effective_labels"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}

	url, err := tpgresource.ReplaceVarsForId(d, config, "{{NetworkConnectivityBasePath}}projects/{{project}}/locations/global/hubs/{{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating Hub %q: %#v", d.Id(), obj)
	headers := make(http.Header)
	updateMask := []string{}

	if d.HasChange("description") {
		updateMask = append(updateMask, "description")
	}

	if d.HasChange("export_psc") {
		updateMask = append(updateMask, "exportPsc")
	}

	if d.HasChange("effective_labels") {
		updateMask = append(updateMask, "labels")
	}
	// updateMask is a URL parameter but not present in the schema, so ReplaceVars
	// won't set it
	url, err = transport_tpg.AddQueryParams(url, map[string]string{"updateMask": strings.Join(updateMask, ",")})
	if err != nil {
		return err
	}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	// if updateMask is empty we are not updating anything so skip the post
	if len(updateMask) > 0 {
		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "PATCH",
			Project:   billingProject,
			RawURL:    url,
			UserAgent: userAgent,
			Body:      obj,
			Timeout:   d.Timeout(schema.TimeoutUpdate),
			Headers:   headers,
		})

		if err != nil {
			return fmt.Errorf("Error updating Hub %q: %s", d.Id(), err)
		} else {
			log.Printf("[DEBUG] Finished updating Hub %q: %#v", d.Id(), res)
		}

		err = NetworkConnectivityOperationWaitTime(
			config, res, tpgresource.GetResourceNameFromSelfLink(project), "Updating Hub", userAgent,
			d.Timeout(schema.TimeoutUpdate))

		if err != nil {
			return err
		}
	}

	return resourceNetworkConnectivityHubRead(d, meta)
}

func resourceNetworkConnectivityHubDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Hub: %s", err)
	}
	billingProject = strings.TrimPrefix(project, "projects/")

	url, err := tpgresource.ReplaceVarsForId(d, config, "{{NetworkConnectivityBasePath}}projects/{{project}}/locations/global/hubs/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)

	log.Printf("[DEBUG] Deleting Hub %q", d.Id())
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "DELETE",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutDelete),
		Headers:   headers,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, "Hub")
	}

	err = NetworkConnectivityOperationWaitTime(
		config, res, tpgresource.GetResourceNameFromSelfLink(project), "Deleting Hub", userAgent,
		d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting Hub %q: %#v", d.Id(), res)
	return nil
}

func resourceNetworkConnectivityHubImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^projects/(?P<project>[^/]+)/locations/global/hubs/(?P<name>[^/]+)$",
		"^(?P<project>[^/]+)/(?P<name>[^/]+)$",
		"^(?P<name>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVarsForId(d, config, "projects/{{project}}/locations/global/hubs/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenNetworkConnectivityHubName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetworkConnectivityHubCreateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetworkConnectivityHubUpdateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetworkConnectivityHubLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}

	transformed := make(map[string]interface{})
	if l, ok := d.GetOkExists("labels"); ok {
		for k := range l.(map[string]interface{}) {
			transformed[k] = v.(map[string]interface{})[k]
		}
	}

	return transformed
}

func flattenNetworkConnectivityHubDescription(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetworkConnectivityHubUniqueId(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetworkConnectivityHubState(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetworkConnectivityHubRoutingVpcs(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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
			"uri": flattenNetworkConnectivityHubRoutingVpcsUri(original["uri"], d, config),
		})
	}
	return transformed
}
func flattenNetworkConnectivityHubRoutingVpcsUri(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetworkConnectivityHubPresetTopology(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetworkConnectivityHubExportPsc(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetworkConnectivityHubTerraformLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}

	transformed := make(map[string]interface{})
	if l, ok := d.GetOkExists("terraform_labels"); ok {
		for k := range l.(map[string]interface{}) {
			transformed[k] = v.(map[string]interface{})[k]
		}
	}

	return transformed
}

func flattenNetworkConnectivityHubEffectiveLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandNetworkConnectivityHubName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNetworkConnectivityHubDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNetworkConnectivityHubPresetTopology(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNetworkConnectivityHubExportPsc(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNetworkConnectivityHubEffectiveLabels(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}
