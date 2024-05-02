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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"
)

func ResourceNetworkConnectivityRegionalEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkConnectivityRegionalEndpointCreate,
		Read:   resourceNetworkConnectivityRegionalEndpointRead,
		Update: resourceNetworkConnectivityRegionalEndpointUpdate,
		Delete: resourceNetworkConnectivityRegionalEndpointDelete,

		Importer: &schema.ResourceImporter{
			State: resourceNetworkConnectivityRegionalEndpointImport,
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
			"access_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidateEnum([]string{"GLOBAL", "REGIONAL"}),
				Description:  `The access type of this regional endpoint. This field is reflected in the PSC Forwarding Rule configuration to enable global access. Possible values: ["GLOBAL", "REGIONAL"]`,
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The location of the RegionalEndpoint.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the RegionalEndpoint.`,
			},
			"target_google_api": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The service endpoint this private regional endpoint connects to. Format: '{apiname}.{region}.p.rep.googleapis.com' Example: \"cloudkms.us-central1.p.rep.googleapis.com\".`,
			},
			"address": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
				Description: `The IP Address of the Regional Endpoint. When no address is provided, an IP from the subnetwork is allocated. Use one of the following formats: * IPv4 address as in '10.0.0.1' * Address resource URI as in 'projects/{project}/regions/{region}/addresses/{address_name}'

~> **Note:** This field accepts both a reference to a Compute Address resource, which is the resource name of which format is given in the description, and IP literal value. If the user chooses to input a reserved address value; they need to make sure that the reserved address is in IPv4 version, its purpose is GCE_ENDPOINT, its type is INTERNAL and its status is RESERVED. If the user chooses to input an IP literal, they need to make sure that it's a valid IPv4 address (x.x.x.x) within the subnetwork.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `A description of this resource.`,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Description: `User-defined labels.


**Note**: This field is non-authoritative, and will only manage the labels present in your configuration.
Please refer to the field 'effective_labels' for all of the labels present on the resource.`,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"network": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: `The name of the VPC network for this private regional endpoint. Format: 'projects/{project}/global/networks/{network}'`,
			},
			"subnetwork": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: `The name of the subnetwork from which the IP address will be allocated. Format: 'projects/{project}/regions/{region}/subnetworks/{subnetwork}'`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Time when the RegionalEndpoint was created.`,
			},
			"effective_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				ForceNew:    true,
				Description: `All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"psc_forwarding_rule": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource reference of the PSC Forwarding Rule created on behalf of the customer. Format: '//compute.googleapis.com/projects/{project}/regions/{region}/forwardingRules/{forwarding_rule_name}'`,
			},
			"terraform_labels": {
				Type:     schema.TypeMap,
				Computed: true,
				Description: `The combination of labels configured directly on the resource
 and default labels configured on the provider.`,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Time when the RegionalEndpoint was updated.`,
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

func resourceNetworkConnectivityRegionalEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	descriptionProp, err := expandNetworkConnectivityRegionalEndpointDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	targetGoogleApiProp, err := expandNetworkConnectivityRegionalEndpointTargetGoogleApi(d.Get("target_google_api"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("target_google_api"); !tpgresource.IsEmptyValue(reflect.ValueOf(targetGoogleApiProp)) && (ok || !reflect.DeepEqual(v, targetGoogleApiProp)) {
		obj["targetGoogleApi"] = targetGoogleApiProp
	}
	networkProp, err := expandNetworkConnectivityRegionalEndpointNetwork(d.Get("network"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("network"); !tpgresource.IsEmptyValue(reflect.ValueOf(networkProp)) && (ok || !reflect.DeepEqual(v, networkProp)) {
		obj["network"] = networkProp
	}
	subnetworkProp, err := expandNetworkConnectivityRegionalEndpointSubnetwork(d.Get("subnetwork"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("subnetwork"); !tpgresource.IsEmptyValue(reflect.ValueOf(subnetworkProp)) && (ok || !reflect.DeepEqual(v, subnetworkProp)) {
		obj["subnetwork"] = subnetworkProp
	}
	accessTypeProp, err := expandNetworkConnectivityRegionalEndpointAccessType(d.Get("access_type"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("access_type"); !tpgresource.IsEmptyValue(reflect.ValueOf(accessTypeProp)) && (ok || !reflect.DeepEqual(v, accessTypeProp)) {
		obj["accessType"] = accessTypeProp
	}
	addressProp, err := expandNetworkConnectivityRegionalEndpointAddress(d.Get("address"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("address"); !tpgresource.IsEmptyValue(reflect.ValueOf(addressProp)) && (ok || !reflect.DeepEqual(v, addressProp)) {
		obj["address"] = addressProp
	}
	labelsProp, err := expandNetworkConnectivityRegionalEndpointEffectiveLabels(d.Get("effective_labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("effective_labels"); !tpgresource.IsEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{NetworkConnectivityBasePath}}projects/{{project}}/locations/{{location}}/regionalEndpoints?regional_endpoint_id={{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new RegionalEndpoint: %#v", obj)
	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for RegionalEndpoint: %s", err)
	}
	billingProject = project

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
		return fmt.Errorf("Error creating RegionalEndpoint: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/regionalEndpoints/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = NetworkConnectivityOperationWaitTime(
		config, res, project, "Creating RegionalEndpoint", userAgent,
		d.Timeout(schema.TimeoutCreate))

	if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create RegionalEndpoint: %s", err)
	}

	log.Printf("[DEBUG] Finished creating RegionalEndpoint %q: %#v", d.Id(), res)

	return resourceNetworkConnectivityRegionalEndpointRead(d, meta)
}

func resourceNetworkConnectivityRegionalEndpointRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{NetworkConnectivityBasePath}}projects/{{project}}/locations/{{location}}/regionalEndpoints/{{name}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for RegionalEndpoint: %s", err)
	}
	billingProject = project

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
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("NetworkConnectivityRegionalEndpoint %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading RegionalEndpoint: %s", err)
	}

	if err := d.Set("create_time", flattenNetworkConnectivityRegionalEndpointCreateTime(res["createTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading RegionalEndpoint: %s", err)
	}
	if err := d.Set("update_time", flattenNetworkConnectivityRegionalEndpointUpdateTime(res["updateTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading RegionalEndpoint: %s", err)
	}
	if err := d.Set("labels", flattenNetworkConnectivityRegionalEndpointLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading RegionalEndpoint: %s", err)
	}
	if err := d.Set("description", flattenNetworkConnectivityRegionalEndpointDescription(res["description"], d, config)); err != nil {
		return fmt.Errorf("Error reading RegionalEndpoint: %s", err)
	}
	if err := d.Set("target_google_api", flattenNetworkConnectivityRegionalEndpointTargetGoogleApi(res["targetGoogleApi"], d, config)); err != nil {
		return fmt.Errorf("Error reading RegionalEndpoint: %s", err)
	}
	if err := d.Set("network", flattenNetworkConnectivityRegionalEndpointNetwork(res["network"], d, config)); err != nil {
		return fmt.Errorf("Error reading RegionalEndpoint: %s", err)
	}
	if err := d.Set("subnetwork", flattenNetworkConnectivityRegionalEndpointSubnetwork(res["subnetwork"], d, config)); err != nil {
		return fmt.Errorf("Error reading RegionalEndpoint: %s", err)
	}
	if err := d.Set("access_type", flattenNetworkConnectivityRegionalEndpointAccessType(res["accessType"], d, config)); err != nil {
		return fmt.Errorf("Error reading RegionalEndpoint: %s", err)
	}
	if err := d.Set("psc_forwarding_rule", flattenNetworkConnectivityRegionalEndpointPscForwardingRule(res["pscForwardingRule"], d, config)); err != nil {
		return fmt.Errorf("Error reading RegionalEndpoint: %s", err)
	}
	if err := d.Set("address", flattenNetworkConnectivityRegionalEndpointAddress(res["address"], d, config)); err != nil {
		return fmt.Errorf("Error reading RegionalEndpoint: %s", err)
	}
	if err := d.Set("terraform_labels", flattenNetworkConnectivityRegionalEndpointTerraformLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading RegionalEndpoint: %s", err)
	}
	if err := d.Set("effective_labels", flattenNetworkConnectivityRegionalEndpointEffectiveLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading RegionalEndpoint: %s", err)
	}

	return nil
}

func resourceNetworkConnectivityRegionalEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	// Only the root field "labels" and "terraform_labels" are mutable
	return resourceNetworkConnectivityRegionalEndpointRead(d, meta)
}

func resourceNetworkConnectivityRegionalEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for RegionalEndpoint: %s", err)
	}
	billingProject = project

	url, err := tpgresource.ReplaceVars(d, config, "{{NetworkConnectivityBasePath}}projects/{{project}}/locations/{{location}}/regionalEndpoints/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)

	log.Printf("[DEBUG] Deleting RegionalEndpoint %q", d.Id())
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
		return transport_tpg.HandleNotFoundError(err, d, "RegionalEndpoint")
	}

	err = NetworkConnectivityOperationWaitTime(
		config, res, project, "Deleting RegionalEndpoint", userAgent,
		d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting RegionalEndpoint %q: %#v", d.Id(), res)
	return nil
}

func resourceNetworkConnectivityRegionalEndpointImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/regionalEndpoints/(?P<name>[^/]+)$",
		"^(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<name>[^/]+)$",
		"^(?P<location>[^/]+)/(?P<name>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/regionalEndpoints/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenNetworkConnectivityRegionalEndpointCreateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetworkConnectivityRegionalEndpointUpdateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetworkConnectivityRegionalEndpointLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenNetworkConnectivityRegionalEndpointDescription(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetworkConnectivityRegionalEndpointTargetGoogleApi(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetworkConnectivityRegionalEndpointNetwork(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetworkConnectivityRegionalEndpointSubnetwork(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetworkConnectivityRegionalEndpointAccessType(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetworkConnectivityRegionalEndpointPscForwardingRule(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetworkConnectivityRegionalEndpointAddress(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetworkConnectivityRegionalEndpointTerraformLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenNetworkConnectivityRegionalEndpointEffectiveLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandNetworkConnectivityRegionalEndpointDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNetworkConnectivityRegionalEndpointTargetGoogleApi(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNetworkConnectivityRegionalEndpointNetwork(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNetworkConnectivityRegionalEndpointSubnetwork(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNetworkConnectivityRegionalEndpointAccessType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNetworkConnectivityRegionalEndpointAddress(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNetworkConnectivityRegionalEndpointEffectiveLabels(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}
