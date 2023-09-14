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

package apigee

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"
)

func ResourceApigeeTargetServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceApigeeTargetServerCreate,
		Read:   resourceApigeeTargetServerRead,
		Update: resourceApigeeTargetServerUpdate,
		Delete: resourceApigeeTargetServerDelete,

		Importer: &schema.ResourceImporter{
			State: resourceApigeeTargetServerImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"env_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `The Apigee environment group associated with the Apigee environment,
in the format 'organizations/{{org_name}}/environments/{{env_name}}'.`,
			},
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The host name this target connects to. Value must be a valid hostname as described by RFC-1123.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The resource id of this reference. Values must match the regular expression [\w\s-.]+.`,
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The port number this target connects to on the given host. Value must be between 1 and 65535, inclusive.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `A human-readable description of this TargetServer.`,
			},
			"is_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Enabling/disabling a TargetServer is useful when TargetServers are used in load balancing configurations, and one or more TargetServers need to taken out of rotation periodically. Defaults to true.`,
				Default:     true,
			},
			"protocol": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidateEnum([]string{"HTTP", "HTTP2", "GRPC_TARGET", "GRPC", "EXTERNAL_CALLOUT", ""}),
				Description:  `Immutable. The protocol used by this TargetServer. Possible values: ["HTTP", "HTTP2", "GRPC_TARGET", "GRPC", "EXTERNAL_CALLOUT"]`,
			},
			"s_sl_info": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies TLS configuration info for this TargetServer. The JSON name is sSLInfo for legacy/backwards compatibility reasons -- Edge originally supported SSL, and the name is still used for TLS configuration.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `Enables TLS. If false, neither one-way nor two-way TLS will be enabled.`,
						},
						"ciphers": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `The SSL/TLS cipher suites to be used. For programmable proxies, it must be one of the cipher suite names listed in: http://docs.oracle.com/javase/8/docs/technotes/guides/security/StandardNames.html#ciphersuites. For configurable proxies, it must follow the configuration specified in: https://commondatastorage.googleapis.com/chromium-boringssl-docs/ssl.h.html#Cipher-suite-configuration. This setting has no effect for configurable proxies when negotiating TLS 1.3.`,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"client_auth_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Enables two-way TLS.`,
						},
						"common_name": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `The TLS Common Name of the certificate.`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The TLS Common Name string of the certificate.`,
									},
									"wildcard_match": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: `Indicates whether the cert should be matched against as a wildcard cert.`,
									},
								},
							},
						},
						"ignore_validation_errors": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `If true, Edge ignores TLS certificate errors. Valid when configuring TLS for target servers and target endpoints, and when configuring virtual hosts that use 2-way TLS. When used with a target endpoint/target server, if the backend system uses SNI and returns a cert with a subject Distinguished Name (DN) that does not match the hostname, there is no way to ignore the error and the connection fails.`,
						},
						"key_alias": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Required if clientAuthEnabled is true. The resource ID for the alias containing the private key and cert.`,
						},
						"key_store": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Required if clientAuthEnabled is true. The resource ID of the keystore.`,
						},
						"protocols": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `The TLS versioins to be used.`,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"trust_store": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The resource ID of the truststore.`,
						},
					},
				},
			},
		},
		UseJSONNumber: true,
	}
}

func resourceApigeeTargetServerCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	nameProp, err := expandApigeeTargetServerName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	descriptionProp, err := expandApigeeTargetServerDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	hostProp, err := expandApigeeTargetServerHost(d.Get("host"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("host"); !tpgresource.IsEmptyValue(reflect.ValueOf(hostProp)) && (ok || !reflect.DeepEqual(v, hostProp)) {
		obj["host"] = hostProp
	}
	portProp, err := expandApigeeTargetServerPort(d.Get("port"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("port"); !tpgresource.IsEmptyValue(reflect.ValueOf(portProp)) && (ok || !reflect.DeepEqual(v, portProp)) {
		obj["port"] = portProp
	}
	isEnabledProp, err := expandApigeeTargetServerIsEnabled(d.Get("is_enabled"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("is_enabled"); !tpgresource.IsEmptyValue(reflect.ValueOf(isEnabledProp)) && (ok || !reflect.DeepEqual(v, isEnabledProp)) {
		obj["isEnabled"] = isEnabledProp
	}
	sSLInfoProp, err := expandApigeeTargetServerSSLInfo(d.Get("s_sl_info"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("s_sl_info"); !tpgresource.IsEmptyValue(reflect.ValueOf(sSLInfoProp)) && (ok || !reflect.DeepEqual(v, sSLInfoProp)) {
		obj["sSLInfo"] = sSLInfoProp
	}
	protocolProp, err := expandApigeeTargetServerProtocol(d.Get("protocol"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("protocol"); !tpgresource.IsEmptyValue(reflect.ValueOf(protocolProp)) && (ok || !reflect.DeepEqual(v, protocolProp)) {
		obj["protocol"] = protocolProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ApigeeBasePath}}{{env_id}}/targetservers")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new TargetServer: %#v", obj)
	billingProject := ""

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
		return fmt.Errorf("Error creating TargetServer: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "{{env_id}}/targetservers/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating TargetServer %q: %#v", d.Id(), res)

	return resourceApigeeTargetServerRead(d, meta)
}

func resourceApigeeTargetServerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ApigeeBasePath}}{{env_id}}/targetservers/{{name}}")
	if err != nil {
		return err
	}

	billingProject := ""

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
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("ApigeeTargetServer %q", d.Id()))
	}

	if err := d.Set("name", flattenApigeeTargetServerName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading TargetServer: %s", err)
	}
	if err := d.Set("description", flattenApigeeTargetServerDescription(res["description"], d, config)); err != nil {
		return fmt.Errorf("Error reading TargetServer: %s", err)
	}
	if err := d.Set("host", flattenApigeeTargetServerHost(res["host"], d, config)); err != nil {
		return fmt.Errorf("Error reading TargetServer: %s", err)
	}
	if err := d.Set("port", flattenApigeeTargetServerPort(res["port"], d, config)); err != nil {
		return fmt.Errorf("Error reading TargetServer: %s", err)
	}
	if err := d.Set("is_enabled", flattenApigeeTargetServerIsEnabled(res["isEnabled"], d, config)); err != nil {
		return fmt.Errorf("Error reading TargetServer: %s", err)
	}
	if err := d.Set("s_sl_info", flattenApigeeTargetServerSSLInfo(res["sSLInfo"], d, config)); err != nil {
		return fmt.Errorf("Error reading TargetServer: %s", err)
	}
	if err := d.Set("protocol", flattenApigeeTargetServerProtocol(res["protocol"], d, config)); err != nil {
		return fmt.Errorf("Error reading TargetServer: %s", err)
	}

	return nil
}

func resourceApigeeTargetServerUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	obj := make(map[string]interface{})
	nameProp, err := expandApigeeTargetServerName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	descriptionProp, err := expandApigeeTargetServerDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	hostProp, err := expandApigeeTargetServerHost(d.Get("host"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("host"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, hostProp)) {
		obj["host"] = hostProp
	}
	portProp, err := expandApigeeTargetServerPort(d.Get("port"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("port"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, portProp)) {
		obj["port"] = portProp
	}
	isEnabledProp, err := expandApigeeTargetServerIsEnabled(d.Get("is_enabled"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("is_enabled"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, isEnabledProp)) {
		obj["isEnabled"] = isEnabledProp
	}
	sSLInfoProp, err := expandApigeeTargetServerSSLInfo(d.Get("s_sl_info"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("s_sl_info"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, sSLInfoProp)) {
		obj["sSLInfo"] = sSLInfoProp
	}
	protocolProp, err := expandApigeeTargetServerProtocol(d.Get("protocol"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("protocol"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, protocolProp)) {
		obj["protocol"] = protocolProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ApigeeBasePath}}{{env_id}}/targetservers/{{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating TargetServer %q: %#v", d.Id(), obj)

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "PUT",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutUpdate),
	})

	if err != nil {
		return fmt.Errorf("Error updating TargetServer %q: %s", d.Id(), err)
	} else {
		log.Printf("[DEBUG] Finished updating TargetServer %q: %#v", d.Id(), res)
	}

	return resourceApigeeTargetServerRead(d, meta)
}

func resourceApigeeTargetServerDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	url, err := tpgresource.ReplaceVars(d, config, "{{ApigeeBasePath}}{{env_id}}/targetservers/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting TargetServer %q", d.Id())

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
		return transport_tpg.HandleNotFoundError(err, d, "TargetServer")
	}

	log.Printf("[DEBUG] Finished deleting TargetServer %q: %#v", d.Id(), res)
	return nil
}

func resourceApigeeTargetServerImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	// current import_formats cannot import fields with forward slashes in their value
	if err := tpgresource.ParseImportId([]string{
		"(?P<env_id>.+)/targetservers/(?P<name>.+)",
		"(?P<env_id>.+)/(?P<name>.+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "{{env_id}}/targetservers/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenApigeeTargetServerName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeTargetServerDescription(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeTargetServerHost(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeTargetServerPort(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenApigeeTargetServerIsEnabled(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeTargetServerSSLInfo(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["enabled"] =
		flattenApigeeTargetServerSSLInfoEnabled(original["enabled"], d, config)
	transformed["client_auth_enabled"] =
		flattenApigeeTargetServerSSLInfoClientAuthEnabled(original["clientAuthEnabled"], d, config)
	transformed["key_store"] =
		flattenApigeeTargetServerSSLInfoKeyStore(original["keyStore"], d, config)
	transformed["key_alias"] =
		flattenApigeeTargetServerSSLInfoKeyAlias(original["keyAlias"], d, config)
	transformed["trust_store"] =
		flattenApigeeTargetServerSSLInfoTrustStore(original["trustStore"], d, config)
	transformed["ignore_validation_errors"] =
		flattenApigeeTargetServerSSLInfoIgnoreValidationErrors(original["ignoreValidationErrors"], d, config)
	transformed["protocols"] =
		flattenApigeeTargetServerSSLInfoProtocols(original["protocols"], d, config)
	transformed["ciphers"] =
		flattenApigeeTargetServerSSLInfoCiphers(original["ciphers"], d, config)
	transformed["common_name"] =
		flattenApigeeTargetServerSSLInfoCommonName(original["commonName"], d, config)
	return []interface{}{transformed}
}
func flattenApigeeTargetServerSSLInfoEnabled(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeTargetServerSSLInfoClientAuthEnabled(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeTargetServerSSLInfoKeyStore(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeTargetServerSSLInfoKeyAlias(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeTargetServerSSLInfoTrustStore(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeTargetServerSSLInfoIgnoreValidationErrors(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeTargetServerSSLInfoProtocols(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeTargetServerSSLInfoCiphers(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeTargetServerSSLInfoCommonName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["value"] =
		flattenApigeeTargetServerSSLInfoCommonNameValue(original["value"], d, config)
	transformed["wildcard_match"] =
		flattenApigeeTargetServerSSLInfoCommonNameWildcardMatch(original["wildcardMatch"], d, config)
	return []interface{}{transformed}
}
func flattenApigeeTargetServerSSLInfoCommonNameValue(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeTargetServerSSLInfoCommonNameWildcardMatch(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeTargetServerProtocol(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandApigeeTargetServerName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeTargetServerDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeTargetServerHost(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeTargetServerPort(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeTargetServerIsEnabled(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeTargetServerSSLInfo(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedEnabled, err := expandApigeeTargetServerSSLInfoEnabled(original["enabled"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedEnabled); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["enabled"] = transformedEnabled
	}

	transformedClientAuthEnabled, err := expandApigeeTargetServerSSLInfoClientAuthEnabled(original["client_auth_enabled"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedClientAuthEnabled); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["clientAuthEnabled"] = transformedClientAuthEnabled
	}

	transformedKeyStore, err := expandApigeeTargetServerSSLInfoKeyStore(original["key_store"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedKeyStore); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["keyStore"] = transformedKeyStore
	}

	transformedKeyAlias, err := expandApigeeTargetServerSSLInfoKeyAlias(original["key_alias"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedKeyAlias); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["keyAlias"] = transformedKeyAlias
	}

	transformedTrustStore, err := expandApigeeTargetServerSSLInfoTrustStore(original["trust_store"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedTrustStore); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["trustStore"] = transformedTrustStore
	}

	transformedIgnoreValidationErrors, err := expandApigeeTargetServerSSLInfoIgnoreValidationErrors(original["ignore_validation_errors"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedIgnoreValidationErrors); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["ignoreValidationErrors"] = transformedIgnoreValidationErrors
	}

	transformedProtocols, err := expandApigeeTargetServerSSLInfoProtocols(original["protocols"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedProtocols); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["protocols"] = transformedProtocols
	}

	transformedCiphers, err := expandApigeeTargetServerSSLInfoCiphers(original["ciphers"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedCiphers); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["ciphers"] = transformedCiphers
	}

	transformedCommonName, err := expandApigeeTargetServerSSLInfoCommonName(original["common_name"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedCommonName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["commonName"] = transformedCommonName
	}

	return transformed, nil
}

func expandApigeeTargetServerSSLInfoEnabled(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeTargetServerSSLInfoClientAuthEnabled(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeTargetServerSSLInfoKeyStore(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeTargetServerSSLInfoKeyAlias(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeTargetServerSSLInfoTrustStore(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeTargetServerSSLInfoIgnoreValidationErrors(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeTargetServerSSLInfoProtocols(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeTargetServerSSLInfoCiphers(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeTargetServerSSLInfoCommonName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedValue, err := expandApigeeTargetServerSSLInfoCommonNameValue(original["value"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedValue); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["value"] = transformedValue
	}

	transformedWildcardMatch, err := expandApigeeTargetServerSSLInfoCommonNameWildcardMatch(original["wildcard_match"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedWildcardMatch); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["wildcardMatch"] = transformedWildcardMatch
	}

	return transformed, nil
}

func expandApigeeTargetServerSSLInfoCommonNameValue(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeTargetServerSSLInfoCommonNameWildcardMatch(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeTargetServerProtocol(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
