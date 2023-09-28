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

package certificatemanager

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"
)

func ResourceCertificateManagerCertificateIssuanceConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceCertificateManagerCertificateIssuanceConfigCreate,
		Read:   resourceCertificateManagerCertificateIssuanceConfigRead,
		Delete: resourceCertificateManagerCertificateIssuanceConfigDelete,

		Importer: &schema.ResourceImporter{
			State: resourceCertificateManagerCertificateIssuanceConfigImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			tpgresource.SetLabelsDiff,
			tpgresource.DefaultProviderProject,
		),

		Schema: map[string]*schema.Schema{
			"certificate_authority_config": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: `The CA that issues the workload certificate. It includes the CA address, type, authentication to CA service, etc.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"certificate_authority_service_config": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: `Defines a CertificateAuthorityServiceConfig.`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ca_pool": {
										Type:             schema.TypeString,
										Required:         true,
										ForceNew:         true,
										DiffSuppressFunc: tpgresource.CompareResourceNames,
										Description: `A CA pool resource used to issue a certificate.
The CA pool string has a relative resource path following the form
"projects/{project}/locations/{location}/caPools/{caPool}".`,
									},
								},
							},
						},
					},
				},
			},
			"key_algorithm": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidateEnum([]string{"RSA_2048", "ECDSA_P256"}),
				Description:  `Key algorithm to use when generating the private key. Possible values: ["RSA_2048", "ECDSA_P256"]`,
			},
			"lifetime": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `Lifetime of issued certificates. A duration in seconds with up to nine fractional digits, ending with 's'.
Example: "1814400s". Valid values are from 21 days (1814400s) to 30 days (2592000s)`,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `A user-defined name of the certificate issuance config.
CertificateIssuanceConfig names must be unique globally.`,
			},
			"rotation_window_percentage": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
				Description: `It specifies the percentage of elapsed time of the certificate lifetime to wait before renewing the certificate.
Must be a number between 1-99, inclusive.
You must set the rotation window percentage in relation to the certificate lifetime so that certificate renewal occurs at least 7 days after
the certificate has been issued and at least 7 days before it expires.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `One or more paragraphs of text description of a CertificateIssuanceConfig.`,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Description: `'Set of label tags associated with the CertificateIssuanceConfig resource.
 An object containing a list of "key": value pairs. Example: { "name": "wrench", "count": "3" }.


**Note**: This field is non-authoritative, and will only manage the labels present in your configuration.
Please refer to the field 'effective_labels' for all of the labels present on the resource.`,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The Certificate Manager location. If not specified, "global" is used.`,
				Default:     "global",
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The creation timestamp of a CertificateIssuanceConfig. Timestamp is in RFC3339 UTC "Zulu" format,
accurate to nanoseconds with up to nine fractional digits.
Examples: "2014-10-02T15:01:23Z" and "2014-10-02T15:01:23.045123456Z".`,
			},
			"effective_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				ForceNew:    true,
				Description: `All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"terraform_labels": {
				Type:     schema.TypeMap,
				Computed: true,
				Description: `The combination of labels configured directly on the resource
 and default labels configured on the provider.`,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The last update timestamp of a CertificateIssuanceConfig. Timestamp is in RFC3339 UTC "Zulu" format,
accurate to nanoseconds with up to nine fractional digits.
Examples: "2014-10-02T15:01:23Z" and "2014-10-02T15:01:23.045123456Z".`,
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

func resourceCertificateManagerCertificateIssuanceConfigCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	descriptionProp, err := expandCertificateManagerCertificateIssuanceConfigDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	rotationWindowPercentageProp, err := expandCertificateManagerCertificateIssuanceConfigRotationWindowPercentage(d.Get("rotation_window_percentage"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("rotation_window_percentage"); !tpgresource.IsEmptyValue(reflect.ValueOf(rotationWindowPercentageProp)) && (ok || !reflect.DeepEqual(v, rotationWindowPercentageProp)) {
		obj["rotationWindowPercentage"] = rotationWindowPercentageProp
	}
	keyAlgorithmProp, err := expandCertificateManagerCertificateIssuanceConfigKeyAlgorithm(d.Get("key_algorithm"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("key_algorithm"); !tpgresource.IsEmptyValue(reflect.ValueOf(keyAlgorithmProp)) && (ok || !reflect.DeepEqual(v, keyAlgorithmProp)) {
		obj["keyAlgorithm"] = keyAlgorithmProp
	}
	lifetimeProp, err := expandCertificateManagerCertificateIssuanceConfigLifetime(d.Get("lifetime"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("lifetime"); !tpgresource.IsEmptyValue(reflect.ValueOf(lifetimeProp)) && (ok || !reflect.DeepEqual(v, lifetimeProp)) {
		obj["lifetime"] = lifetimeProp
	}
	certificateAuthorityConfigProp, err := expandCertificateManagerCertificateIssuanceConfigCertificateAuthorityConfig(d.Get("certificate_authority_config"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("certificate_authority_config"); !tpgresource.IsEmptyValue(reflect.ValueOf(certificateAuthorityConfigProp)) && (ok || !reflect.DeepEqual(v, certificateAuthorityConfigProp)) {
		obj["certificateAuthorityConfig"] = certificateAuthorityConfigProp
	}
	labelsProp, err := expandCertificateManagerCertificateIssuanceConfigEffectiveLabels(d.Get("effective_labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("effective_labels"); !tpgresource.IsEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{CertificateManagerBasePath}}projects/{{project}}/locations/{{location}}/certificateIssuanceConfigs?certificateIssuanceConfigId={{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new CertificateIssuanceConfig: %#v", obj)
	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for CertificateIssuanceConfig: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutCreate),
	})
	if err != nil {
		return fmt.Errorf("Error creating CertificateIssuanceConfig: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/certificateIssuanceConfigs/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = CertificateManagerOperationWaitTime(
		config, res, project, "Creating CertificateIssuanceConfig", userAgent,
		d.Timeout(schema.TimeoutCreate))

	if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create CertificateIssuanceConfig: %s", err)
	}

	log.Printf("[DEBUG] Finished creating CertificateIssuanceConfig %q: %#v", d.Id(), res)

	return resourceCertificateManagerCertificateIssuanceConfigRead(d, meta)
}

func resourceCertificateManagerCertificateIssuanceConfigRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{CertificateManagerBasePath}}projects/{{project}}/locations/{{location}}/certificateIssuanceConfigs/{{name}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for CertificateIssuanceConfig: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("CertificateManagerCertificateIssuanceConfig %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading CertificateIssuanceConfig: %s", err)
	}

	if err := d.Set("description", flattenCertificateManagerCertificateIssuanceConfigDescription(res["description"], d, config)); err != nil {
		return fmt.Errorf("Error reading CertificateIssuanceConfig: %s", err)
	}
	if err := d.Set("rotation_window_percentage", flattenCertificateManagerCertificateIssuanceConfigRotationWindowPercentage(res["rotationWindowPercentage"], d, config)); err != nil {
		return fmt.Errorf("Error reading CertificateIssuanceConfig: %s", err)
	}
	if err := d.Set("key_algorithm", flattenCertificateManagerCertificateIssuanceConfigKeyAlgorithm(res["keyAlgorithm"], d, config)); err != nil {
		return fmt.Errorf("Error reading CertificateIssuanceConfig: %s", err)
	}
	if err := d.Set("lifetime", flattenCertificateManagerCertificateIssuanceConfigLifetime(res["lifetime"], d, config)); err != nil {
		return fmt.Errorf("Error reading CertificateIssuanceConfig: %s", err)
	}
	if err := d.Set("create_time", flattenCertificateManagerCertificateIssuanceConfigCreateTime(res["createTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading CertificateIssuanceConfig: %s", err)
	}
	if err := d.Set("update_time", flattenCertificateManagerCertificateIssuanceConfigUpdateTime(res["updateTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading CertificateIssuanceConfig: %s", err)
	}
	if err := d.Set("labels", flattenCertificateManagerCertificateIssuanceConfigLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading CertificateIssuanceConfig: %s", err)
	}
	if err := d.Set("certificate_authority_config", flattenCertificateManagerCertificateIssuanceConfigCertificateAuthorityConfig(res["certificateAuthorityConfig"], d, config)); err != nil {
		return fmt.Errorf("Error reading CertificateIssuanceConfig: %s", err)
	}
	if err := d.Set("terraform_labels", flattenCertificateManagerCertificateIssuanceConfigTerraformLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading CertificateIssuanceConfig: %s", err)
	}
	if err := d.Set("effective_labels", flattenCertificateManagerCertificateIssuanceConfigEffectiveLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading CertificateIssuanceConfig: %s", err)
	}

	return nil
}

func resourceCertificateManagerCertificateIssuanceConfigDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for CertificateIssuanceConfig: %s", err)
	}
	billingProject = project

	url, err := tpgresource.ReplaceVars(d, config, "{{CertificateManagerBasePath}}projects/{{project}}/locations/{{location}}/certificateIssuanceConfigs/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting CertificateIssuanceConfig %q", d.Id())

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "DELETE",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutDelete),
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, "CertificateIssuanceConfig")
	}

	err = CertificateManagerOperationWaitTime(
		config, res, project, "Deleting CertificateIssuanceConfig", userAgent,
		d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting CertificateIssuanceConfig %q: %#v", d.Id(), res)
	return nil
}

func resourceCertificateManagerCertificateIssuanceConfigImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/certificateIssuanceConfigs/(?P<name>[^/]+)$",
		"^(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<name>[^/]+)$",
		"^(?P<location>[^/]+)/(?P<name>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/certificateIssuanceConfigs/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenCertificateManagerCertificateIssuanceConfigDescription(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenCertificateManagerCertificateIssuanceConfigRotationWindowPercentage(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenCertificateManagerCertificateIssuanceConfigKeyAlgorithm(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenCertificateManagerCertificateIssuanceConfigLifetime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenCertificateManagerCertificateIssuanceConfigCreateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenCertificateManagerCertificateIssuanceConfigUpdateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenCertificateManagerCertificateIssuanceConfigLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenCertificateManagerCertificateIssuanceConfigCertificateAuthorityConfig(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["certificate_authority_service_config"] =
		flattenCertificateManagerCertificateIssuanceConfigCertificateAuthorityConfigCertificateAuthorityServiceConfig(original["certificateAuthorityServiceConfig"], d, config)
	return []interface{}{transformed}
}
func flattenCertificateManagerCertificateIssuanceConfigCertificateAuthorityConfigCertificateAuthorityServiceConfig(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["ca_pool"] =
		flattenCertificateManagerCertificateIssuanceConfigCertificateAuthorityConfigCertificateAuthorityServiceConfigCaPool(original["caPool"], d, config)
	return []interface{}{transformed}
}
func flattenCertificateManagerCertificateIssuanceConfigCertificateAuthorityConfigCertificateAuthorityServiceConfigCaPool(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenCertificateManagerCertificateIssuanceConfigTerraformLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenCertificateManagerCertificateIssuanceConfigEffectiveLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandCertificateManagerCertificateIssuanceConfigDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandCertificateManagerCertificateIssuanceConfigRotationWindowPercentage(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandCertificateManagerCertificateIssuanceConfigKeyAlgorithm(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandCertificateManagerCertificateIssuanceConfigLifetime(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandCertificateManagerCertificateIssuanceConfigCertificateAuthorityConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedCertificateAuthorityServiceConfig, err := expandCertificateManagerCertificateIssuanceConfigCertificateAuthorityConfigCertificateAuthorityServiceConfig(original["certificate_authority_service_config"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedCertificateAuthorityServiceConfig); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["certificateAuthorityServiceConfig"] = transformedCertificateAuthorityServiceConfig
	}

	return transformed, nil
}

func expandCertificateManagerCertificateIssuanceConfigCertificateAuthorityConfigCertificateAuthorityServiceConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedCaPool, err := expandCertificateManagerCertificateIssuanceConfigCertificateAuthorityConfigCertificateAuthorityServiceConfigCaPool(original["ca_pool"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedCaPool); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["caPool"] = transformedCaPool
	}

	return transformed, nil
}

func expandCertificateManagerCertificateIssuanceConfigCertificateAuthorityConfigCertificateAuthorityServiceConfigCaPool(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandCertificateManagerCertificateIssuanceConfigEffectiveLabels(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}
