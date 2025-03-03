// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This code is generated by Magic Modules using the following:
//
//     Configuration: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/networkservices/Mesh.yaml
//     Template:      https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/resource.go.tmpl
//
//     DO NOT EDIT this file directly. Any changes made to this file will be
//     overwritten during the next generation cycle.
//
// ----------------------------------------------------------------------------

package networkservices

import (
	"context"
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

func ResourceNetworkServicesMesh() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkServicesMeshCreate,
		Read:   resourceNetworkServicesMeshRead,
		Update: resourceNetworkServicesMeshUpdate,
		Delete: resourceNetworkServicesMeshDelete,

		Importer: &schema.ResourceImporter{
			State: resourceNetworkServicesMeshImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 1,

		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceNetworkServicesMeshResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: ResourceNetworkServicesMeshUpgradeV0,
				Version: 0,
			},
		},
		CustomizeDiff: customdiff.All(
			tpgresource.SetLabelsDiff,
			tpgresource.DefaultProviderProject,
		),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Short name of the Mesh resource to be created.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `A free-text description of the resource. Max length 1024 characters.`,
			},
			"interception_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: `Optional. If set to a valid TCP port (1-65535), instructs the SIDECAR proxy to listen on the
specified port of localhost (127.0.0.1) address. The SIDECAR proxy will expect all traffic to
be redirected to this port regardless of its actual ip:port destination. If unset, a port
'15001' is used as the interception port. This will is applicable only for sidecar proxy
deployments.`,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Description: `Set of label tags associated with the Mesh resource.

**Note**: This field is non-authoritative, and will only manage the labels present in your configuration.
Please refer to the field 'effective_labels' for all of the labels present on the resource.`,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"location": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidateRegexp(`^global$`),
				Description:  `Location (region) of the Mesh resource to be created. Only the value 'global' is currently allowed; defaults to 'global' if omitted.`,
				Default:      "global",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Time the Mesh was created in UTC.`,
			},
			"effective_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: `All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"self_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Server-defined URL of this resource.`,
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
				Description: `Time the Mesh was updated in UTC.`,
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

func resourceNetworkServicesMeshCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	descriptionProp, err := expandNetworkServicesMeshDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	interceptionPortProp, err := expandNetworkServicesMeshInterceptionPort(d.Get("interception_port"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("interception_port"); !tpgresource.IsEmptyValue(reflect.ValueOf(interceptionPortProp)) && (ok || !reflect.DeepEqual(v, interceptionPortProp)) {
		obj["interceptionPort"] = interceptionPortProp
	}
	labelsProp, err := expandNetworkServicesMeshEffectiveLabels(d.Get("effective_labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("effective_labels"); !tpgresource.IsEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{NetworkServicesBasePath}}projects/{{project}}/locations/{{location}}/meshes?meshId={{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Mesh: %#v", obj)
	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Mesh: %s", err)
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
		return fmt.Errorf("Error creating Mesh: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/meshes/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = NetworkServicesOperationWaitTime(
		config, res, project, "Creating Mesh", userAgent,
		d.Timeout(schema.TimeoutCreate))

	if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create Mesh: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Mesh %q: %#v", d.Id(), res)

	return resourceNetworkServicesMeshRead(d, meta)
}

func resourceNetworkServicesMeshRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{NetworkServicesBasePath}}projects/{{project}}/locations/{{location}}/meshes/{{name}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Mesh: %s", err)
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
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("NetworkServicesMesh %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading Mesh: %s", err)
	}

	if err := d.Set("self_link", flattenNetworkServicesMeshSelfLink(res["selfLink"], d, config)); err != nil {
		return fmt.Errorf("Error reading Mesh: %s", err)
	}
	if err := d.Set("create_time", flattenNetworkServicesMeshCreateTime(res["createTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading Mesh: %s", err)
	}
	if err := d.Set("update_time", flattenNetworkServicesMeshUpdateTime(res["updateTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading Mesh: %s", err)
	}
	if err := d.Set("labels", flattenNetworkServicesMeshLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading Mesh: %s", err)
	}
	if err := d.Set("description", flattenNetworkServicesMeshDescription(res["description"], d, config)); err != nil {
		return fmt.Errorf("Error reading Mesh: %s", err)
	}
	if err := d.Set("interception_port", flattenNetworkServicesMeshInterceptionPort(res["interceptionPort"], d, config)); err != nil {
		return fmt.Errorf("Error reading Mesh: %s", err)
	}
	if err := d.Set("terraform_labels", flattenNetworkServicesMeshTerraformLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading Mesh: %s", err)
	}
	if err := d.Set("effective_labels", flattenNetworkServicesMeshEffectiveLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading Mesh: %s", err)
	}

	return nil
}

func resourceNetworkServicesMeshUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Mesh: %s", err)
	}
	billingProject = project

	obj := make(map[string]interface{})
	descriptionProp, err := expandNetworkServicesMeshDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	interceptionPortProp, err := expandNetworkServicesMeshInterceptionPort(d.Get("interception_port"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("interception_port"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, interceptionPortProp)) {
		obj["interceptionPort"] = interceptionPortProp
	}
	labelsProp, err := expandNetworkServicesMeshEffectiveLabels(d.Get("effective_labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("effective_labels"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{NetworkServicesBasePath}}projects/{{project}}/locations/{{location}}/meshes/{{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating Mesh %q: %#v", d.Id(), obj)
	headers := make(http.Header)
	updateMask := []string{}

	if d.HasChange("description") {
		updateMask = append(updateMask, "description")
	}

	if d.HasChange("interception_port") {
		updateMask = append(updateMask, "interceptionPort")
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
			return fmt.Errorf("Error updating Mesh %q: %s", d.Id(), err)
		} else {
			log.Printf("[DEBUG] Finished updating Mesh %q: %#v", d.Id(), res)
		}

		err = NetworkServicesOperationWaitTime(
			config, res, project, "Updating Mesh", userAgent,
			d.Timeout(schema.TimeoutUpdate))

		if err != nil {
			return err
		}
	}

	return resourceNetworkServicesMeshRead(d, meta)
}

func resourceNetworkServicesMeshDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Mesh: %s", err)
	}
	billingProject = project

	url, err := tpgresource.ReplaceVars(d, config, "{{NetworkServicesBasePath}}projects/{{project}}/locations/{{location}}/meshes/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)

	log.Printf("[DEBUG] Deleting Mesh %q", d.Id())
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
		return transport_tpg.HandleNotFoundError(err, d, "Mesh")
	}

	err = NetworkServicesOperationWaitTime(
		config, res, project, "Deleting Mesh", userAgent,
		d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting Mesh %q: %#v", d.Id(), res)
	return nil
}

func resourceNetworkServicesMeshImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/meshes/(?P<name>[^/]+)$",
		"^(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<name>[^/]+)$",
		"^(?P<location>[^/]+)/(?P<name>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/meshes/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenNetworkServicesMeshSelfLink(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetworkServicesMeshCreateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetworkServicesMeshUpdateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetworkServicesMeshLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenNetworkServicesMeshDescription(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetworkServicesMeshInterceptionPort(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
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

func flattenNetworkServicesMeshTerraformLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenNetworkServicesMeshEffectiveLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandNetworkServicesMeshDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNetworkServicesMeshInterceptionPort(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNetworkServicesMeshEffectiveLabels(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func resourceNetworkServicesMeshResourceV0() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkServicesMeshCreate,
		Read:   resourceNetworkServicesMeshRead,
		Update: resourceNetworkServicesMeshUpdate,
		Delete: resourceNetworkServicesMeshDelete,

		Importer: &schema.ResourceImporter{
			State: resourceNetworkServicesMeshImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			tpgresource.SetLabelsDiff,
			tpgresource.DefaultProviderProject,
		),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Short name of the Mesh resource to be created.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `A free-text description of the resource. Max length 1024 characters.`,
			},
			"interception_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: `Optional. If set to a valid TCP port (1-65535), instructs the SIDECAR proxy to listen on the
specified port of localhost (127.0.0.1) address. The SIDECAR proxy will expect all traffic to
be redirected to this port regardless of its actual ip:port destination. If unset, a port
'15001' is used as the interception port. This will is applicable only for sidecar proxy
deployments.`,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Description: `Set of label tags associated with the Mesh resource.

**Note**: This field is non-authoritative, and will only manage the labels present in your configuration.
Please refer to the field 'effective_labels' for all of the labels present on the resource.`,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Time the Mesh was created in UTC.`,
			},
			"effective_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: `All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"self_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Server-defined URL of this resource.`,
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
				Description: `Time the Mesh was updated in UTC.`,
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

func ResourceNetworkServicesMeshUpgradeV0(_ context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	log.Printf("[DEBUG] Attributes before migration: %#v", rawState)
	if _, ok := rawState["location"]; !ok {
		rawState["location"] = "global"
	}
	log.Printf("[DEBUG] Attributes after migration: %#v", rawState)
	return rawState, nil
}
