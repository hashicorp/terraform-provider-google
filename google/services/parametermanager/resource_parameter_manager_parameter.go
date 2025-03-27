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
//     Configuration: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/parametermanager/Parameter.yaml
//     Template:      https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/resource.go.tmpl
//
//     DO NOT EDIT this file directly. Any changes made to this file will be
//     overwritten during the next generation cycle.
//
// ----------------------------------------------------------------------------

package parametermanager

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

func ResourceParameterManagerParameter() *schema.Resource {
	return &schema.Resource{
		Create: resourceParameterManagerParameterCreate,
		Read:   resourceParameterManagerParameterRead,
		Update: resourceParameterManagerParameterUpdate,
		Delete: resourceParameterManagerParameterDelete,

		Importer: &schema.ResourceImporter{
			State: resourceParameterManagerParameterImport,
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
			"parameter_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `This must be unique within the project.`,
			},
			"format": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidateEnum([]string{"UNFORMATTED", "YAML", "JSON", ""}),
				Description:  `The format type of the parameter resource. Default value: "UNFORMATTED" Possible values: ["UNFORMATTED", "YAML", "JSON"]`,
				Default:      "UNFORMATTED",
			},
			"kms_key": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `The resource name of the Cloud KMS CryptoKey used to encrypt parameter version payload. Format
'projects/{{project}}/locations/global/keyRings/{{key_ring}}/cryptoKeys/{{crypto_key}}'`,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Description: `The labels assigned to this Parameter.

Label keys must be between 1 and 63 characters long, have a UTF-8 encoding of maximum 128 bytes,
and must conform to the following PCRE regular expression: [\p{Ll}\p{Lo}][\p{Ll}\p{Lo}\p{N}_-]{0,62}

Label values must be between 0 and 63 characters long, have a UTF-8 encoding of maximum 128 bytes,
and must conform to the following PCRE regular expression: [\p{Ll}\p{Lo}\p{N}_-]{0,63}

No more than 64 labels can be assigned to a given resource.

An object containing a list of "key": value pairs. Example:
{ "name": "wrench", "mass": "1.3kg", "count": "3" }.


**Note**: This field is non-authoritative, and will only manage the labels present in your configuration.
Please refer to the field 'effective_labels' for all of the labels present on the resource.`,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time at which the Parameter was created.`,
			},
			"effective_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: `All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The resource name of the Parameter. Format:
'projects/{{project}}/locations/global/parameters/{{parameter_id}}'`,
			},
			"policy_member": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Policy member strings of a Google Cloud resource.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"iam_policy_name_principal": {
							Type:     schema.TypeString,
							Computed: true,
							Description: `IAM policy binding member referring to a Google Cloud resource by user-assigned name. If a
resource is deleted and recreated with the same name, the binding will be applicable to the
new resource. Format:
'principal://parametermanager.googleapis.com/projects/{{project}}/name/locations/global/parameters/{{parameter_id}}'`,
						},
						"iam_policy_uid_principal": {
							Type:     schema.TypeString,
							Computed: true,
							Description: `IAM policy binding member referring to a Google Cloud resource by system-assigned unique identifier.
If a resource is deleted and recreated with the same name, the binding will not be applicable to the
new resource. Format:
'principal://parametermanager.googleapis.com/projects/{{project}}/uid/locations/global/parameters/{{uid}}'`,
						},
					},
				},
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
				Description: `The time at which the Parameter was updated.`,
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

func resourceParameterManagerParameterCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	formatProp, err := expandParameterManagerParameterFormat(d.Get("format"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("format"); !tpgresource.IsEmptyValue(reflect.ValueOf(formatProp)) && (ok || !reflect.DeepEqual(v, formatProp)) {
		obj["format"] = formatProp
	}
	kmsKeyProp, err := expandParameterManagerParameterKmsKey(d.Get("kms_key"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("kms_key"); !tpgresource.IsEmptyValue(reflect.ValueOf(kmsKeyProp)) && (ok || !reflect.DeepEqual(v, kmsKeyProp)) {
		obj["kmsKey"] = kmsKeyProp
	}
	labelsProp, err := expandParameterManagerParameterEffectiveLabels(d.Get("effective_labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("effective_labels"); !tpgresource.IsEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ParameterManagerBasePath}}projects/{{project}}/locations/global/parameters?parameter_id={{parameter_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Parameter: %#v", obj)
	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Parameter: %s", err)
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
		return fmt.Errorf("Error creating Parameter: %s", err)
	}
	if err := d.Set("name", flattenParameterManagerParameterName(res["name"], d, config)); err != nil {
		return fmt.Errorf(`Error setting computed identity field "name": %s`, err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/global/parameters/{{parameter_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating Parameter %q: %#v", d.Id(), res)

	return resourceParameterManagerParameterRead(d, meta)
}

func resourceParameterManagerParameterRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ParameterManagerBasePath}}projects/{{project}}/locations/global/parameters/{{parameter_id}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Parameter: %s", err)
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
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("ParameterManagerParameter %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading Parameter: %s", err)
	}

	if err := d.Set("name", flattenParameterManagerParameterName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading Parameter: %s", err)
	}
	if err := d.Set("create_time", flattenParameterManagerParameterCreateTime(res["createTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading Parameter: %s", err)
	}
	if err := d.Set("update_time", flattenParameterManagerParameterUpdateTime(res["updateTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading Parameter: %s", err)
	}
	if err := d.Set("policy_member", flattenParameterManagerParameterPolicyMember(res["policyMember"], d, config)); err != nil {
		return fmt.Errorf("Error reading Parameter: %s", err)
	}
	if err := d.Set("labels", flattenParameterManagerParameterLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading Parameter: %s", err)
	}
	if err := d.Set("format", flattenParameterManagerParameterFormat(res["format"], d, config)); err != nil {
		return fmt.Errorf("Error reading Parameter: %s", err)
	}
	if err := d.Set("kms_key", flattenParameterManagerParameterKmsKey(res["kmsKey"], d, config)); err != nil {
		return fmt.Errorf("Error reading Parameter: %s", err)
	}
	if err := d.Set("terraform_labels", flattenParameterManagerParameterTerraformLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading Parameter: %s", err)
	}
	if err := d.Set("effective_labels", flattenParameterManagerParameterEffectiveLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading Parameter: %s", err)
	}

	return nil
}

func resourceParameterManagerParameterUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Parameter: %s", err)
	}
	billingProject = project

	obj := make(map[string]interface{})
	kmsKeyProp, err := expandParameterManagerParameterKmsKey(d.Get("kms_key"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("kms_key"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, kmsKeyProp)) {
		obj["kmsKey"] = kmsKeyProp
	}
	labelsProp, err := expandParameterManagerParameterEffectiveLabels(d.Get("effective_labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("effective_labels"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ParameterManagerBasePath}}projects/{{project}}/locations/global/parameters/{{parameter_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating Parameter %q: %#v", d.Id(), obj)
	headers := make(http.Header)
	updateMask := []string{}

	if d.HasChange("kms_key") {
		updateMask = append(updateMask, "kmsKey")
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
			return fmt.Errorf("Error updating Parameter %q: %s", d.Id(), err)
		} else {
			log.Printf("[DEBUG] Finished updating Parameter %q: %#v", d.Id(), res)
		}

	}

	return resourceParameterManagerParameterRead(d, meta)
}

func resourceParameterManagerParameterDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Parameter: %s", err)
	}
	billingProject = project

	url, err := tpgresource.ReplaceVars(d, config, "{{ParameterManagerBasePath}}projects/{{project}}/locations/global/parameters/{{parameter_id}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)

	log.Printf("[DEBUG] Deleting Parameter %q", d.Id())
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
		return transport_tpg.HandleNotFoundError(err, d, "Parameter")
	}

	log.Printf("[DEBUG] Finished deleting Parameter %q: %#v", d.Id(), res)
	return nil
}

func resourceParameterManagerParameterImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^projects/(?P<project>[^/]+)/locations/global/parameters/(?P<parameter_id>[^/]+)$",
		"^(?P<project>[^/]+)/(?P<parameter_id>[^/]+)$",
		"^(?P<parameter_id>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/global/parameters/{{parameter_id}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenParameterManagerParameterName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenParameterManagerParameterCreateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenParameterManagerParameterUpdateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenParameterManagerParameterPolicyMember(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["iam_policy_uid_principal"] =
		flattenParameterManagerParameterPolicyMemberIamPolicyUidPrincipal(original["iamPolicyUidPrincipal"], d, config)
	transformed["iam_policy_name_principal"] =
		flattenParameterManagerParameterPolicyMemberIamPolicyNamePrincipal(original["iamPolicyNamePrincipal"], d, config)
	return []interface{}{transformed}
}
func flattenParameterManagerParameterPolicyMemberIamPolicyUidPrincipal(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenParameterManagerParameterPolicyMemberIamPolicyNamePrincipal(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenParameterManagerParameterLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenParameterManagerParameterFormat(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenParameterManagerParameterKmsKey(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenParameterManagerParameterTerraformLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenParameterManagerParameterEffectiveLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandParameterManagerParameterFormat(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandParameterManagerParameterKmsKey(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandParameterManagerParameterEffectiveLabels(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}
