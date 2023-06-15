// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package apigee

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceApigeeKeystoresAliasesPkcs12() *schema.Resource {
	return &schema.Resource{
		Create: ResourceApigeeKeystoresAliasesPkcs12Create,
		Read:   ResourceApigeeKeystoresAliasesPkcs12Read,
		Delete: ResourceApigeeKeystoresAliasesPkcs12Delete,

		Importer: &schema.ResourceImporter{
			State: ResourceApigeeKeystoresAliasesPkcs12Import,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"alias": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Alias Name`,
			},
			"file": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
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
				Required:    true,
				ForceNew:    true,
				Description: `Keystore Name`,
			},
			"org_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: `Organization ID associated with the alias`,
			},
			"filehash": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Hash of the pkcs file",
			},
			"certs_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Chain of certificates under this alias.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cert_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `List of all properties in the object.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"basic_constraints": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `X.509 basic constraints extension.`,
									},
									"expiry_date": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `X.509 notAfter validity period in milliseconds since epoch.`,
									},
									"is_valid": {
										Type:     schema.TypeString,
										Computed: true,
										Description: `Flag that specifies whether the certificate is valid. 
Flag is set to Yes if the certificate is valid, No if expired, or Not yet if not yet valid.`,
									},
									"issuer": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `X.509 issuer.`,
									},
									"public_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Public key component of the X.509 subject public key info.`,
									},
									"serial_number": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `X.509 serial number.`,
									},
									"sig_alg_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `X.509 signatureAlgorithm.`,
									},
									"subject": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `X.509 subject.`,
									},
									"subject_alternative_names": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `X.509 subject alternative names (SANs) extension.`,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"valid_from": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `X.509 notBefore validity period in milliseconds since epoch.`,
									},
									"version": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `X.509 version.`,
									},
								},
							},
						},
					},
				},
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
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

func ResourceApigeeKeystoresAliasesPkcs12Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	filePath, _ := d.GetOk("file")
	file, err := os.Open(filePath.(string))
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	bw := multipart.NewWriter(buf)
	if password, ok := d.GetOkExists("password"); ok {
		keyFilePartWriter, _ := bw.CreateFormField("password")
		keyFilePartWriter.Write([]byte(password.(string)))
	}
	certFilePartWriter, _ := bw.CreateFormField("file")
	_, err = io.Copy(certFilePartWriter, file)
	bw.Close()
	file.Close()
	url, err := tpgresource.ReplaceVars(d, config, "{{ApigeeBasePath}}organizations/{{org_id}}/environments/{{environment}}/keystores/{{keystore}}/aliases?format=pkcs12&alias={{alias}}&ignoreExpiryValidation=true")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new KeystoresAliasesPkcs12")
	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestRawBodyWithTimeout(config, "POST", billingProject, url, userAgent, buf, "multipart/form-data; boundary="+bw.Boundary(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating KeystoresAliasesPkcs12: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "organizations/{{org_id}}/environments/{{environment}}/keystores/{{keystore}}/aliases/{{alias}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating KeystoreAliasesPkcs %q: %#v", d.Id(), res)

	return ResourceApigeeKeystoresAliasesPkcs12Read(d, meta)
}

func ResourceApigeeKeystoresAliasesPkcs12Read(d *schema.ResourceData, meta interface{}) error {
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
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("ApigeeKeystoreAliasesPkcs %q", d.Id()))
	}

	if err := d.Set("alias", flattenApigeeKeystoreAliasesPkcsAlias(res["alias"], d, config)); err != nil {
		return fmt.Errorf("Error reading KeystoreAliasesPkcs: %s", err)
	}

	if err := d.Set("certs_info", flattenApigeeKeystoreAliasesPkcsCertsInfo(res["certsInfo"], d, config)); err != nil {
		return fmt.Errorf("Error reading KeystoreAliasesPkcs: %s", err)
	}

	if err := d.Set("type", flattenApigeeKeystoreAliasesPkcsType(res["type"], d, config)); err != nil {
		return fmt.Errorf("Error reading KeystoreAliasesPkcs: %s", err)
	}

	return nil
}

func ResourceApigeeKeystoresAliasesPkcs12Delete(d *schema.ResourceData, meta interface{}) error {
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
	log.Printf("[DEBUG] Deleting KeystoreAliasesPkcs %q", d.Id())

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
		return transport_tpg.HandleNotFoundError(err, d, "KeystoreAliasesPkcs")
	}

	log.Printf("[DEBUG] Finished deleting KeystoreAliasesPkcs %q: %#v", d.Id(), res)
	return nil
}

func ResourceApigeeKeystoresAliasesPkcs12Import(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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

func flattenApigeeKeystoreAliasesPkcsOrgId(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoreAliasesPkcsEnvironment(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoreAliasesPkcsKeystore(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoreAliasesPkcsAlias(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoreAliasesPkcsKey(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoreAliasesPkcsPassword(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoreAliasesPkcsCert(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoreAliasesPkcsCertsInfo(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["cert_info"] =
		flattenApigeeKeystoreAliasesPkcsCertsInfoCertInfo(original["certInfo"], d, config)
	return []interface{}{transformed}
}
func flattenApigeeKeystoreAliasesPkcsCertsInfoCertInfo(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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
			"version":                   flattenApigeeKeystoreAliasesPkcsCertsInfoCertInfoVersion(original["version"], d, config),
			"subject":                   flattenApigeeKeystoreAliasesPkcsCertsInfoCertInfoSubject(original["subject"], d, config),
			"issuer":                    flattenApigeeKeystoreAliasesPkcsCertsInfoCertInfoIssuer(original["issuer"], d, config),
			"expiry_date":               flattenApigeeKeystoreAliasesPkcsCertsInfoCertInfoExpiryDate(original["expiryDate"], d, config),
			"valid_from":                flattenApigeeKeystoreAliasesPkcsCertsInfoCertInfoValidFrom(original["validFrom"], d, config),
			"is_valid":                  flattenApigeeKeystoreAliasesPkcsCertsInfoCertInfoIsValid(original["isValid"], d, config),
			"subject_alternative_names": flattenApigeeKeystoreAliasesPkcsCertsInfoCertInfoSubjectAlternativeNames(original["subjectAlternativeNames"], d, config),
			"sig_alg_name":              flattenApigeeKeystoreAliasesPkcsCertsInfoCertInfoSigAlgName(original["sigAlgName"], d, config),
			"public_key":                flattenApigeeKeystoreAliasesPkcsCertsInfoCertInfoPublicKey(original["publicKey"], d, config),
			"basic_constraints":         flattenApigeeKeystoreAliasesPkcsCertsInfoCertInfoBasicConstraints(original["basicConstraints"], d, config),
			"serial_number":             flattenApigeeKeystoreAliasesPkcsCertsInfoCertInfoSerialNumber(original["serialNumber"], d, config),
		})
	}
	return transformed
}
func flattenApigeeKeystoreAliasesPkcsCertsInfoCertInfoVersion(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenApigeeKeystoreAliasesPkcsCertsInfoCertInfoSubject(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoreAliasesPkcsCertsInfoCertInfoIssuer(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoreAliasesPkcsCertsInfoCertInfoExpiryDate(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoreAliasesPkcsCertsInfoCertInfoValidFrom(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoreAliasesPkcsCertsInfoCertInfoIsValid(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoreAliasesPkcsCertsInfoCertInfoSubjectAlternativeNames(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoreAliasesPkcsCertsInfoCertInfoSigAlgName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoreAliasesPkcsCertsInfoCertInfoPublicKey(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoreAliasesPkcsCertsInfoCertInfoBasicConstraints(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoreAliasesPkcsCertsInfoCertInfoSerialNumber(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeKeystoreAliasesPkcsType(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandApigeeKeystoreAliasesPkcsOrgId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoreAliasesPkcsEnvironment(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoreAliasesPkcsKeystore(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoreAliasesPkcsAlias(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoreAliasesPkcsKey(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoreAliasesPkcsPassword(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoreAliasesPkcsCert(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoreAliasesPkcsCertsInfo(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedCertInfo, err := expandApigeeKeystoreAliasesPkcsCertsInfoCertInfo(original["cert_info"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedCertInfo); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["certInfo"] = transformedCertInfo
	}

	return transformed, nil
}

func expandApigeeKeystoreAliasesPkcsCertsInfoCertInfo(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedVersion, err := expandApigeeKeystoreAliasesPkcsCertsInfoCertInfoVersion(original["version"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedVersion); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["version"] = transformedVersion
		}

		transformedSubject, err := expandApigeeKeystoreAliasesPkcsCertsInfoCertInfoSubject(original["subject"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedSubject); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["subject"] = transformedSubject
		}

		transformedIssuer, err := expandApigeeKeystoreAliasesPkcsCertsInfoCertInfoIssuer(original["issuer"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedIssuer); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["issuer"] = transformedIssuer
		}

		transformedExpiryDate, err := expandApigeeKeystoreAliasesPkcsCertsInfoCertInfoExpiryDate(original["expiry_date"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedExpiryDate); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["expiryDate"] = transformedExpiryDate
		}

		transformedValidFrom, err := expandApigeeKeystoreAliasesPkcsCertsInfoCertInfoValidFrom(original["valid_from"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedValidFrom); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["validFrom"] = transformedValidFrom
		}

		transformedIsValid, err := expandApigeeKeystoreAliasesPkcsCertsInfoCertInfoIsValid(original["is_valid"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedIsValid); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["isValid"] = transformedIsValid
		}

		transformedSubjectAlternativeNames, err := expandApigeeKeystoreAliasesPkcsCertsInfoCertInfoSubjectAlternativeNames(original["subject_alternative_names"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedSubjectAlternativeNames); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["subjectAlternativeNames"] = transformedSubjectAlternativeNames
		}

		transformedSigAlgName, err := expandApigeeKeystoreAliasesPkcsCertsInfoCertInfoSigAlgName(original["sig_alg_name"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedSigAlgName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["sigAlgName"] = transformedSigAlgName
		}

		transformedPublicKey, err := expandApigeeKeystoreAliasesPkcsCertsInfoCertInfoPublicKey(original["public_key"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedPublicKey); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["publicKey"] = transformedPublicKey
		}

		transformedBasicConstraints, err := expandApigeeKeystoreAliasesPkcsCertsInfoCertInfoBasicConstraints(original["basic_constraints"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedBasicConstraints); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["basicConstraints"] = transformedBasicConstraints
		}

		transformedSerialNumber, err := expandApigeeKeystoreAliasesPkcsCertsInfoCertInfoSerialNumber(original["serial_number"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedSerialNumber); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["serialNumber"] = transformedSerialNumber
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandApigeeKeystoreAliasesPkcsCertsInfoCertInfoVersion(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoreAliasesPkcsCertsInfoCertInfoSubject(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoreAliasesPkcsCertsInfoCertInfoIssuer(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoreAliasesPkcsCertsInfoCertInfoExpiryDate(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoreAliasesPkcsCertsInfoCertInfoValidFrom(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoreAliasesPkcsCertsInfoCertInfoIsValid(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoreAliasesPkcsCertsInfoCertInfoSubjectAlternativeNames(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoreAliasesPkcsCertsInfoCertInfoSigAlgName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoreAliasesPkcsCertsInfoCertInfoPublicKey(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoreAliasesPkcsCertsInfoCertInfoBasicConstraints(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeKeystoreAliasesPkcsCertsInfoCertInfoSerialNumber(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
