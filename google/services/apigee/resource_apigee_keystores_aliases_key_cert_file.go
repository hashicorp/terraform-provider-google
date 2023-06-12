// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package apigee

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceApigeeKeystoresAliasesKeyCertFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceApigeeKeystoresAliasesKeyCertFileCreate,
		Read:   resourceApigeeKeystoresAliasesKeyCertFileRead,
		Update: resourceApigeeKeystoresAliasesKeyCertFileUpdate,
		Delete: resourceApigeeKeystoresAliasesKeyCertFileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceApigeeKeystoresAliasesKeyCertFileImport,
		},

		CustomizeDiff: customdiff.All(
			/*
				If cert is changed then an update is expected, so we tell Terraform core to expect update on certs_info
			*/

			customdiff.ComputedIf("certs_info", func(_ context.Context, diff *schema.ResourceDiff, v interface{}) bool {
				return diff.HasChange("cert")
			}),
		),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Read:   schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"alias": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: `Alias Name`,
			},
			"cert": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Cert content`,
			},
			"environment": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: `Environment associated with the alias`,
			},
			"keystore": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: `Keystore Name`,
			},
			"org_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: `Organization ID associated with the alias`,
			},
			"certs_info": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: `Chain of certificates under this alias.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cert_info": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: `List of all properties in the object.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"basic_constraints": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: `X.509 basic constraints extension.`,
									},
									"expiry_date": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: `X.509 notAfter validity period in milliseconds since epoch.`,
									},
									"is_valid": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										Description: `Flag that specifies whether the certificate is valid. 
Flag is set to Yes if the certificate is valid, No if expired, or Not yet if not yet valid.`,
									},
									"issuer": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: `X.509 issuer.`,
									},
									"public_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: `Public key component of the X.509 subject public key info.`,
									},
									"serial_number": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: `X.509 serial number.`,
									},
									"sig_alg_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: `X.509 signatureAlgorithm.`,
									},
									"subject": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: `X.509 subject.`,
									},
									"subject_alternative_names": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: `X.509 subject alternative names (SANs) extension.`,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"valid_from": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: `X.509 notBefore validity period in milliseconds since epoch.`,
									},
									"version": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: `X.509 version.`,
									},
								},
							},
						},
					},
				},
			},
			"key": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Sensitive:   true,
				Description: `Private Key content, omit if uploading to truststore`,
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: `Password for the Private Key if it's encrypted`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Optional.Type of Alias`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceApigeeKeystoresAliasesKeyCertFileCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	bw := multipart.NewWriter(buf)
	if key, ok := d.GetOkExists("key"); ok {
		keyFilePartWriter, _ := bw.CreateFormField("keyFile")
		keyFilePartWriter.Write([]byte(key.(string)))
	}
	if password, ok := d.GetOkExists("password"); ok {
		keyFilePartWriter, _ := bw.CreateFormField("password")
		keyFilePartWriter.Write([]byte(password.(string)))
	}
	certFilePartWriter, _ := bw.CreateFormField("certFile")
	certFilePartWriter.Write([]byte(d.Get("cert").(string)))
	bw.Close()

	url, err := tpgresource.ReplaceVars(d, config, "{{ApigeeBasePath}}organizations/{{org_id}}/environments/{{environment}}/keystores/{{keystore}}/aliases?format=keycertfile&alias={{alias}}&ignoreExpiryValidation=true")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new KeystoresAliasesKeyCertFile")
	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestRawBodyWithTimeout(config, "POST", billingProject, url, userAgent, buf, "multipart/form-data; boundary="+bw.Boundary(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating KeystoresAliasesKeyCertFile: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "organizations/{{org_id}}/environments/{{environment}}/keystores/{{keystore}}/aliases/{{alias}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating KeystoresAliasesKeyCertFile %q: %#v", d.Id(), res)

	return resourceApigeeKeystoresAliasesKeyCertFileRead(d, meta)
}

func resourceApigeeKeystoresAliasesKeyCertFileRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ApigeeBasePath}}organizations/{{org_id}}/environments/{{environment}}/keystores/{{keystore}}/aliases/{{alias}}")
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
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("ApigeeKeystoresAliasesKeyCertFile %q", d.Id()))
	}

	if err := d.Set("alias", flattenApigeeKeystoresAliasesKeyCertFileAlias(res["alias"], d, config)); err != nil {
		return fmt.Errorf("Error reading KeystoresAliasesKeyCertFile: %s", err)
	}

	if err := d.Set("certs_info", flattenApigeeKeystoresAliasesKeyCertFileCertsInfo(res["certsInfo"], d, config)); err != nil {
		return fmt.Errorf("Error reading KeystoresAliasesKeyCertFile: %s", err)
	}
	if err := d.Set("type", flattenApigeeKeystoresAliasesKeyCertFileType(res["type"], d, config)); err != nil {
		return fmt.Errorf("Error reading KeystoresAliasesKeyCertFile: %s", err)
	}

	return nil
}

func resourceApigeeKeystoresAliasesKeyCertFileUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	url, err := tpgresource.ReplaceVars(d, config, "{{ApigeeBasePath}}organizations/{{org_id}}/environments/{{environment}}/keystores/{{keystore}}/aliases/{{alias}}?ignoreExpiryValidation=true")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating KeystoresAliasesKeyCertFile %q", d.Id())

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	buf := new(bytes.Buffer)
	bw := multipart.NewWriter(buf)
	certFilePartWriter, _ := bw.CreateFormField("certFile")
	certFilePartWriter.Write([]byte(d.Get("cert").(string)))
	bw.Close()

	res, err := sendRequestRawBodyWithTimeout(config, "PUT", billingProject, url, userAgent, buf, "multipart/form-data; boundary="+bw.Boundary(), d.Timeout(schema.TimeoutCreate))

	if err != nil {
		return fmt.Errorf("Error updating KeystoresAliasesKeyCertFile %q: %s", d.Id(), err)
	} else {
		log.Printf("[DEBUG] Finished updating KeystoresAliasesKeyCertFile %q: %#v", d.Id(), res)
	}

	return resourceApigeeKeystoresAliasesKeyCertFileRead(d, meta)
}

func resourceApigeeKeystoresAliasesKeyCertFileDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	url, err := tpgresource.ReplaceVars(d, config, "{{ApigeeBasePath}}organizations/{{org_id}}/environments/{{environment}}/keystores/{{keystore}}/aliases/{{alias}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting KeystoresAliasesKeyCertFile %q", d.Id())

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
		return transport_tpg.HandleNotFoundError(err, d, "KeystoresAliasesKeyCertFile")
	}

	log.Printf("[DEBUG] Finished deleting KeystoresAliasesKeyCertFile %q: %#v", d.Id(), res)
	return nil
}

func resourceApigeeKeystoresAliasesKeyCertFileImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"organizations/(?P<org_id>[^/]+)/environments/(?P<environment>[^/]+)/keystores/(?P<keystore>[^/]+)/aliases/(?P<alias>[^/]+)",
		"(?P<org_id>[^/]+)/(?P<environment>[^/]+)/(?P<keystore>[^/]+)/(?P<alias>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "organizations/{{org_id}}/environments/{{environment}}/keystores/{{keystore}}/aliases/{{alias}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenApigeeKeystoresAliasesKeyCertFileOrgId(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoresAliasesKeyCertFileEnvironment(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoresAliasesKeyCertFileKeystore(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoresAliasesKeyCertFileAlias(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoresAliasesKeyCertFileKey(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoresAliasesKeyCertFilePassword(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoresAliasesKeyCertFileCert(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoresAliasesKeyCertFileCertsInfo(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["cert_info"] =
		flattenApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfo(original["certInfo"], d, config)
	return []interface{}{transformed}
}
func flattenApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfo(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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
			"version":                   flattenApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoVersion(original["version"], d, config),
			"subject":                   flattenApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoSubject(original["subject"], d, config),
			"issuer":                    flattenApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoIssuer(original["issuer"], d, config),
			"expiry_date":               flattenApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoExpiryDate(original["expiryDate"], d, config),
			"valid_from":                flattenApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoValidFrom(original["validFrom"], d, config),
			"is_valid":                  flattenApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoIsValid(original["isValid"], d, config),
			"subject_alternative_names": flattenApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoSubjectAlternativeNames(original["subjectAlternativeNames"], d, config),
			"sig_alg_name":              flattenApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoSigAlgName(original["sigAlgName"], d, config),
			"public_key":                flattenApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoPublicKey(original["publicKey"], d, config),
			"basic_constraints":         flattenApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoBasicConstraints(original["basicConstraints"], d, config),
			"serial_number":             flattenApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoSerialNumber(original["serialNumber"], d, config),
		})
	}
	return transformed
}
func flattenApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoVersion(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoSubject(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoIssuer(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoExpiryDate(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoValidFrom(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoIsValid(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoSubjectAlternativeNames(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoSigAlgName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoPublicKey(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoBasicConstraints(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoSerialNumber(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoresAliasesKeyCertFileType(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandApigeeKeystoresAliasesKeyCertFileOrgId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoresAliasesKeyCertFileEnvironment(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoresAliasesKeyCertFileKeystore(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoresAliasesKeyCertFileAlias(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoresAliasesKeyCertFileKey(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoresAliasesKeyCertFilePassword(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoresAliasesKeyCertFileCert(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoresAliasesKeyCertFileCertsInfo(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedCertInfo, err := expandApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfo(original["cert_info"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedCertInfo); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["certInfo"] = transformedCertInfo
	}

	return transformed, nil
}

func expandApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfo(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedVersion, err := expandApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoVersion(original["version"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedVersion); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["version"] = transformedVersion
		}

		transformedSubject, err := expandApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoSubject(original["subject"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedSubject); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["subject"] = transformedSubject
		}

		transformedIssuer, err := expandApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoIssuer(original["issuer"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedIssuer); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["issuer"] = transformedIssuer
		}

		transformedExpiryDate, err := expandApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoExpiryDate(original["expiry_date"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedExpiryDate); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["expiryDate"] = transformedExpiryDate
		}

		transformedValidFrom, err := expandApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoValidFrom(original["valid_from"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedValidFrom); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["validFrom"] = transformedValidFrom
		}

		transformedIsValid, err := expandApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoIsValid(original["is_valid"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedIsValid); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["isValid"] = transformedIsValid
		}

		transformedSubjectAlternativeNames, err := expandApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoSubjectAlternativeNames(original["subject_alternative_names"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedSubjectAlternativeNames); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["subjectAlternativeNames"] = transformedSubjectAlternativeNames
		}

		transformedSigAlgName, err := expandApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoSigAlgName(original["sig_alg_name"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedSigAlgName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["sigAlgName"] = transformedSigAlgName
		}

		transformedPublicKey, err := expandApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoPublicKey(original["public_key"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedPublicKey); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["publicKey"] = transformedPublicKey
		}

		transformedBasicConstraints, err := expandApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoBasicConstraints(original["basic_constraints"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedBasicConstraints); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["basicConstraints"] = transformedBasicConstraints
		}

		transformedSerialNumber, err := expandApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoSerialNumber(original["serial_number"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedSerialNumber); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["serialNumber"] = transformedSerialNumber
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoVersion(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoSubject(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoIssuer(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoExpiryDate(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoValidFrom(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoIsValid(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoSubjectAlternativeNames(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoSigAlgName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoPublicKey(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoBasicConstraints(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoresAliasesKeyCertFileCertsInfoCertInfoSerialNumber(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
