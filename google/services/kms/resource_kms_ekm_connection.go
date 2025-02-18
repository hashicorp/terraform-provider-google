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
//     Configuration: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/kms/EkmConnection.yaml
//     Template:      https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/resource.go.tmpl
//
//     DO NOT EDIT this file directly. Any changes made to this file will be
//     overwritten during the next generation cycle.
//
// ----------------------------------------------------------------------------

package kms

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

func ResourceKMSEkmConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceKMSEkmConnectionCreate,
		Read:   resourceKMSEkmConnectionRead,
		Update: resourceKMSEkmConnectionUpdate,
		Delete: resourceKMSEkmConnectionDelete,

		Importer: &schema.ResourceImporter{
			State: resourceKMSEkmConnectionImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			tpgresource.DefaultProviderProject,
		),

		Schema: map[string]*schema.Schema{
			"location": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `The location for the EkmConnection.
A full list of valid locations can be found by running 'gcloud kms locations list'.`,
			},
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareResourceNames,
				Description:      `The resource name for the EkmConnection.`,
			},
			"service_resolvers": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `A list of ServiceResolvers where the EKM can be reached. There should be one ServiceResolver per EKM replica. Currently, only a single ServiceResolver is supported`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostname": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Required. The hostname of the EKM replica used at TLS and HTTP layers.`,
						},
						"server_certificates": {
							Type:        schema.TypeList,
							Required:    true,
							Description: `Required. A list of leaf server certificates used to authenticate HTTPS connections to the EKM replica. Currently, a maximum of 10 Certificate is supported.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"raw_der": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Required. The raw certificate bytes in DER format. A base64-encoded string.`,
									},
									"issuer": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Output only. The issuer distinguished name in RFC 2253 format. Only present if parsed is true.`,
									},
									"not_after_time": {
										Type:     schema.TypeString,
										Computed: true,
										Description: `Output only. The certificate is not valid after this time. Only present if parsed is true.
A timestamp in RFC3339 UTC "Zulu" format, with nanosecond resolution and up to nine fractional digits. Examples: "2014-10-02T15:01:23Z" and "2014-10-02T15:01:23.045123456Z".`,
									},
									"not_before_time": {
										Type:     schema.TypeString,
										Computed: true,
										Description: `Output only. The certificate is not valid before this time. Only present if parsed is true.
A timestamp in RFC3339 UTC "Zulu" format, with nanosecond resolution and up to nine fractional digits. Examples: "2014-10-02T15:01:23Z" and "2014-10-02T15:01:23.045123456Z".`,
									},
									"parsed": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Output only. True if the certificate was parsed successfully.`,
									},
									"serial_number": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Output only. The certificate serial number as a hex string. Only present if parsed is true.`,
									},
									"sha256_fingerprint": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Output only. The SHA-256 certificate fingerprint as a hex string. Only present if parsed is true.`,
									},
									"subject": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Output only. The subject distinguished name in RFC 2253 format. Only present if parsed is true.`,
									},
									"subject_alternative_dns_names": {
										Type:        schema.TypeList,
										Computed:    true,
										Optional:    true,
										Description: `Output only. The subject Alternative DNS names. Only present if parsed is true.`,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"service_directory_service": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Required. The resource name of the Service Directory service pointing to an EKM replica, in the format projects/*/locations/*/namespaces/*/services/*`,
						},
						"endpoint_filter": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: `Optional. The filter applied to the endpoints of the resolved service. If no filter is specified, all endpoints will be considered. An endpoint will be chosen arbitrarily from the filtered list for each request. For endpoint filter syntax and examples, see https://cloud.google.com/service-directory/docs/reference/rpc/google.cloud.servicedirectory.v1#resolveservicerequest.`,
						},
					},
				},
			},
			"crypto_space_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: `Optional. Identifies the EKM Crypto Space that this EkmConnection maps to. Note: This field is required if KeyManagementMode is CLOUD_KMS.`,
			},
			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: `Optional. Etag of the currently stored EkmConnection.`,
			},
			"key_management_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: verify.ValidateEnum([]string{"MANUAL", "CLOUD_KMS", ""}),
				Description:  `Optional. Describes who can perform control plane operations on the EKM. If unset, this defaults to MANUAL Default value: "MANUAL" Possible values: ["MANUAL", "CLOUD_KMS"]`,
				Default:      "MANUAL",
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Output only. The time at which the EkmConnection was created.
A timestamp in RFC3339 UTC "Zulu" format, with nanosecond resolution and up to nine fractional digits. Examples: "2014-10-02T15:01:23Z" and "2014-10-02T15:01:23.045123456Z".`,
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

func resourceKMSEkmConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	nameProp, err := expandKMSEkmConnectionName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	serviceResolversProp, err := expandKMSEkmConnectionServiceResolvers(d.Get("service_resolvers"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("service_resolvers"); !tpgresource.IsEmptyValue(reflect.ValueOf(serviceResolversProp)) && (ok || !reflect.DeepEqual(v, serviceResolversProp)) {
		obj["serviceResolvers"] = serviceResolversProp
	}
	keyManagementModeProp, err := expandKMSEkmConnectionKeyManagementMode(d.Get("key_management_mode"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("key_management_mode"); !tpgresource.IsEmptyValue(reflect.ValueOf(keyManagementModeProp)) && (ok || !reflect.DeepEqual(v, keyManagementModeProp)) {
		obj["keyManagementMode"] = keyManagementModeProp
	}
	etagProp, err := expandKMSEkmConnectionEtag(d.Get("etag"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("etag"); !tpgresource.IsEmptyValue(reflect.ValueOf(etagProp)) && (ok || !reflect.DeepEqual(v, etagProp)) {
		obj["etag"] = etagProp
	}
	cryptoSpacePathProp, err := expandKMSEkmConnectionCryptoSpacePath(d.Get("crypto_space_path"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("crypto_space_path"); !tpgresource.IsEmptyValue(reflect.ValueOf(cryptoSpacePathProp)) && (ok || !reflect.DeepEqual(v, cryptoSpacePathProp)) {
		obj["cryptoSpacePath"] = cryptoSpacePathProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{KMSBasePath}}projects/{{project}}/locations/{{location}}/ekmConnections?ekmConnectionId={{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new EkmConnection: %#v", obj)
	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for EkmConnection: %s", err)
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
		return fmt.Errorf("Error creating EkmConnection: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/ekmConnections/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating EkmConnection %q: %#v", d.Id(), res)

	return resourceKMSEkmConnectionRead(d, meta)
}

func resourceKMSEkmConnectionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{KMSBasePath}}projects/{{project}}/locations/{{location}}/ekmConnections/{{name}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for EkmConnection: %s", err)
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
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("KMSEkmConnection %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading EkmConnection: %s", err)
	}

	if err := d.Set("name", flattenKMSEkmConnectionName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading EkmConnection: %s", err)
	}
	if err := d.Set("service_resolvers", flattenKMSEkmConnectionServiceResolvers(res["serviceResolvers"], d, config)); err != nil {
		return fmt.Errorf("Error reading EkmConnection: %s", err)
	}
	if err := d.Set("key_management_mode", flattenKMSEkmConnectionKeyManagementMode(res["keyManagementMode"], d, config)); err != nil {
		return fmt.Errorf("Error reading EkmConnection: %s", err)
	}
	if err := d.Set("etag", flattenKMSEkmConnectionEtag(res["etag"], d, config)); err != nil {
		return fmt.Errorf("Error reading EkmConnection: %s", err)
	}
	if err := d.Set("crypto_space_path", flattenKMSEkmConnectionCryptoSpacePath(res["cryptoSpacePath"], d, config)); err != nil {
		return fmt.Errorf("Error reading EkmConnection: %s", err)
	}
	if err := d.Set("create_time", flattenKMSEkmConnectionCreateTime(res["createTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading EkmConnection: %s", err)
	}

	return nil
}

func resourceKMSEkmConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for EkmConnection: %s", err)
	}
	billingProject = project

	obj := make(map[string]interface{})
	serviceResolversProp, err := expandKMSEkmConnectionServiceResolvers(d.Get("service_resolvers"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("service_resolvers"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, serviceResolversProp)) {
		obj["serviceResolvers"] = serviceResolversProp
	}
	keyManagementModeProp, err := expandKMSEkmConnectionKeyManagementMode(d.Get("key_management_mode"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("key_management_mode"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, keyManagementModeProp)) {
		obj["keyManagementMode"] = keyManagementModeProp
	}
	etagProp, err := expandKMSEkmConnectionEtag(d.Get("etag"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("etag"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, etagProp)) {
		obj["etag"] = etagProp
	}
	cryptoSpacePathProp, err := expandKMSEkmConnectionCryptoSpacePath(d.Get("crypto_space_path"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("crypto_space_path"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, cryptoSpacePathProp)) {
		obj["cryptoSpacePath"] = cryptoSpacePathProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{KMSBasePath}}projects/{{project}}/locations/{{location}}/ekmConnections/{{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating EkmConnection %q: %#v", d.Id(), obj)
	headers := make(http.Header)
	updateMask := []string{}

	if d.HasChange("service_resolvers") {
		updateMask = append(updateMask, "serviceResolvers")
	}

	if d.HasChange("key_management_mode") {
		updateMask = append(updateMask, "keyManagementMode")
	}

	if d.HasChange("etag") {
		updateMask = append(updateMask, "etag")
	}

	if d.HasChange("crypto_space_path") {
		updateMask = append(updateMask, "cryptoSpacePath")
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
			return fmt.Errorf("Error updating EkmConnection %q: %s", d.Id(), err)
		} else {
			log.Printf("[DEBUG] Finished updating EkmConnection %q: %#v", d.Id(), res)
		}

	}

	return resourceKMSEkmConnectionRead(d, meta)
}

func resourceKMSEkmConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARNING] KMS EkmConnection resources"+
		" cannot be deleted from Google Cloud. The resource %s will be removed from Terraform"+
		" state, but will still be present on Google Cloud.", d.Id())
	d.SetId("")

	return nil
}

func resourceKMSEkmConnectionImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/ekmConnections/(?P<name>[^/]+)$",
		"^(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<name>[^/]+)$",
		"^(?P<location>[^/]+)/(?P<name>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/ekmConnections/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenKMSEkmConnectionName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	return tpgresource.NameFromSelfLinkStateFunc(v)
}

func flattenKMSEkmConnectionServiceResolvers(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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
			"service_directory_service": flattenKMSEkmConnectionServiceResolversServiceDirectoryService(original["serviceDirectoryService"], d, config),
			"hostname":                  flattenKMSEkmConnectionServiceResolversHostname(original["hostname"], d, config),
			"server_certificates":       flattenKMSEkmConnectionServiceResolversServerCertificates(original["serverCertificates"], d, config),
			"endpoint_filter":           flattenKMSEkmConnectionServiceResolversEndpointFilter(original["endpointFilter"], d, config),
		})
	}
	return transformed
}
func flattenKMSEkmConnectionServiceResolversServiceDirectoryService(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenKMSEkmConnectionServiceResolversHostname(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenKMSEkmConnectionServiceResolversServerCertificates(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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
			"raw_der":                       flattenKMSEkmConnectionServiceResolversServerCertificatesRawDer(original["rawDer"], d, config),
			"parsed":                        flattenKMSEkmConnectionServiceResolversServerCertificatesParsed(original["parsed"], d, config),
			"issuer":                        flattenKMSEkmConnectionServiceResolversServerCertificatesIssuer(original["issuer"], d, config),
			"subject":                       flattenKMSEkmConnectionServiceResolversServerCertificatesSubject(original["subject"], d, config),
			"not_before_time":               flattenKMSEkmConnectionServiceResolversServerCertificatesNotBeforeTime(original["notBeforeTime"], d, config),
			"not_after_time":                flattenKMSEkmConnectionServiceResolversServerCertificatesNotAfterTime(original["notAfterTime"], d, config),
			"sha256_fingerprint":            flattenKMSEkmConnectionServiceResolversServerCertificatesSha256Fingerprint(original["sha256Fingerprint"], d, config),
			"serial_number":                 flattenKMSEkmConnectionServiceResolversServerCertificatesSerialNumber(original["serialNumber"], d, config),
			"subject_alternative_dns_names": flattenKMSEkmConnectionServiceResolversServerCertificatesSubjectAlternativeDnsNames(original["subjectAlternativeDnsNames"], d, config),
		})
	}
	return transformed
}
func flattenKMSEkmConnectionServiceResolversServerCertificatesRawDer(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenKMSEkmConnectionServiceResolversServerCertificatesParsed(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenKMSEkmConnectionServiceResolversServerCertificatesIssuer(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenKMSEkmConnectionServiceResolversServerCertificatesSubject(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenKMSEkmConnectionServiceResolversServerCertificatesNotBeforeTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenKMSEkmConnectionServiceResolversServerCertificatesNotAfterTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenKMSEkmConnectionServiceResolversServerCertificatesSha256Fingerprint(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenKMSEkmConnectionServiceResolversServerCertificatesSerialNumber(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenKMSEkmConnectionServiceResolversServerCertificatesSubjectAlternativeDnsNames(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenKMSEkmConnectionServiceResolversEndpointFilter(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenKMSEkmConnectionKeyManagementMode(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenKMSEkmConnectionEtag(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenKMSEkmConnectionCryptoSpacePath(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenKMSEkmConnectionCreateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandKMSEkmConnectionName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandKMSEkmConnectionServiceResolvers(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedServiceDirectoryService, err := expandKMSEkmConnectionServiceResolversServiceDirectoryService(original["service_directory_service"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedServiceDirectoryService); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["serviceDirectoryService"] = transformedServiceDirectoryService
		}

		transformedHostname, err := expandKMSEkmConnectionServiceResolversHostname(original["hostname"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedHostname); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["hostname"] = transformedHostname
		}

		transformedServerCertificates, err := expandKMSEkmConnectionServiceResolversServerCertificates(original["server_certificates"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedServerCertificates); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["serverCertificates"] = transformedServerCertificates
		}

		transformedEndpointFilter, err := expandKMSEkmConnectionServiceResolversEndpointFilter(original["endpoint_filter"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedEndpointFilter); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["endpointFilter"] = transformedEndpointFilter
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandKMSEkmConnectionServiceResolversServiceDirectoryService(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandKMSEkmConnectionServiceResolversHostname(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandKMSEkmConnectionServiceResolversServerCertificates(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedRawDer, err := expandKMSEkmConnectionServiceResolversServerCertificatesRawDer(original["raw_der"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedRawDer); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["rawDer"] = transformedRawDer
		}

		transformedParsed, err := expandKMSEkmConnectionServiceResolversServerCertificatesParsed(original["parsed"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedParsed); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["parsed"] = transformedParsed
		}

		transformedIssuer, err := expandKMSEkmConnectionServiceResolversServerCertificatesIssuer(original["issuer"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedIssuer); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["issuer"] = transformedIssuer
		}

		transformedSubject, err := expandKMSEkmConnectionServiceResolversServerCertificatesSubject(original["subject"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedSubject); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["subject"] = transformedSubject
		}

		transformedNotBeforeTime, err := expandKMSEkmConnectionServiceResolversServerCertificatesNotBeforeTime(original["not_before_time"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedNotBeforeTime); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["notBeforeTime"] = transformedNotBeforeTime
		}

		transformedNotAfterTime, err := expandKMSEkmConnectionServiceResolversServerCertificatesNotAfterTime(original["not_after_time"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedNotAfterTime); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["notAfterTime"] = transformedNotAfterTime
		}

		transformedSha256Fingerprint, err := expandKMSEkmConnectionServiceResolversServerCertificatesSha256Fingerprint(original["sha256_fingerprint"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedSha256Fingerprint); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["sha256Fingerprint"] = transformedSha256Fingerprint
		}

		transformedSerialNumber, err := expandKMSEkmConnectionServiceResolversServerCertificatesSerialNumber(original["serial_number"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedSerialNumber); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["serialNumber"] = transformedSerialNumber
		}

		transformedSubjectAlternativeDnsNames, err := expandKMSEkmConnectionServiceResolversServerCertificatesSubjectAlternativeDnsNames(original["subject_alternative_dns_names"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedSubjectAlternativeDnsNames); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["subjectAlternativeDnsNames"] = transformedSubjectAlternativeDnsNames
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandKMSEkmConnectionServiceResolversServerCertificatesRawDer(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandKMSEkmConnectionServiceResolversServerCertificatesParsed(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandKMSEkmConnectionServiceResolversServerCertificatesIssuer(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandKMSEkmConnectionServiceResolversServerCertificatesSubject(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandKMSEkmConnectionServiceResolversServerCertificatesNotBeforeTime(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandKMSEkmConnectionServiceResolversServerCertificatesNotAfterTime(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandKMSEkmConnectionServiceResolversServerCertificatesSha256Fingerprint(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandKMSEkmConnectionServiceResolversServerCertificatesSerialNumber(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandKMSEkmConnectionServiceResolversServerCertificatesSubjectAlternativeDnsNames(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandKMSEkmConnectionServiceResolversEndpointFilter(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandKMSEkmConnectionKeyManagementMode(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandKMSEkmConnectionEtag(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandKMSEkmConnectionCryptoSpacePath(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
