package google

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// This file contains shared flatteners between PrivateCA Certificate, CaPool and CertificateAuthority.
// These resources share the x509Config (Certificate, CertificateAuthorty)/baselineValues (CaPool) object.
// The API does not return this object if it only contains booleans with the default (false) value. This
// causes problems if a user specifies only default values, as Terraform detects that the object has been
// deleted on the API-side. This flattener creates default objects for sub-objects that match this pattern
// to fix perma-diffs on default-only objects. For this to work all objects that are flattened from nil to
// their default object *MUST* be set in the user's config, so they are all marked as Required in the schema.

func flattenPrivatecaCertificateConfigX509ConfigAdditionalExtensions(v interface{}, d *schema.ResourceData, config *Config) interface{} {
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
			"critical":  flattenPrivatecaCertificateConfigX509ConfigAdditionalExtensionsCritical(original["critical"], d, config),
			"value":     flattenPrivatecaCertificateConfigX509ConfigAdditionalExtensionsValue(original["value"], d, config),
			"object_id": flattenPrivatecaCertificateConfigX509ConfigAdditionalExtensionsObjectId(original["objectId"], d, config),
		})
	}
	return transformed
}
func flattenPrivatecaCertificateConfigX509ConfigAdditionalExtensionsCritical(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigAdditionalExtensionsValue(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigAdditionalExtensionsObjectId(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["object_id_path"] =
		flattenPrivatecaCertificateConfigX509ConfigAdditionalExtensionsObjectIdObjectIdPath(original["objectIdPath"], d, config)
	return []interface{}{transformed}
}
func flattenPrivatecaCertificateConfigX509ConfigAdditionalExtensionsObjectIdObjectIdPath(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigPolicyIds(v interface{}, d *schema.ResourceData, config *Config) interface{} {
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
			"object_id_path": flattenPrivatecaCertificateConfigX509ConfigPolicyIdsObjectIdPath(original["objectIdPath"], d, config),
		})
	}
	return transformed
}
func flattenPrivatecaCertificateConfigX509ConfigPolicyIdsObjectIdPath(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigAiaOcspServers(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigCaOptions(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	// Special case here as the CaPool API returns an empty object rather than nil unlike the Certificate
	// and CertificateAuthority APIs.
	if v == nil || len(v.(map[string]interface{})) == 0 {
		v = make(map[string]interface{})
		original := v.(map[string]interface{})
		transformed := make(map[string]interface{})
		transformed["is_ca"] = flattenPrivatecaCertificateConfigX509ConfigCaOptionsIsCa(original["isCa"], d, config)
		return []interface{}{transformed}
	}
	original := v.(map[string]interface{})
	transformed := make(map[string]interface{})
	transformed["is_ca"] =
		flattenPrivatecaCertificateConfigX509ConfigCaOptionsIsCa(original["isCa"], d, config)
	transformed["max_issuer_path_length"] =
		flattenPrivatecaCertificateConfigX509ConfigCaOptionsMaxIssuerPathLength(original["maxIssuerPathLength"], d, config)
	return []interface{}{transformed}
}
func flattenPrivatecaCertificateConfigX509ConfigCaOptionsIsCa(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigCaOptionsMaxIssuerPathLength(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := strconv.ParseInt(strVal, 10, 64); err == nil {
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

func flattenPrivatecaCertificateConfigX509ConfigKeyUsage(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		v = make(map[string]interface{})
	}
	original := v.(map[string]interface{})
	transformed := make(map[string]interface{})
	transformed["base_key_usage"] =
		flattenPrivatecaCertificateConfigX509ConfigKeyUsageBaseKeyUsage(original["baseKeyUsage"], d, config)
	transformed["extended_key_usage"] =
		flattenPrivatecaCertificateConfigX509ConfigKeyUsageExtendedKeyUsage(original["extendedKeyUsage"], d, config)
	transformed["unknown_extended_key_usages"] =
		flattenPrivatecaCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages(original["unknownExtendedKeyUsages"], d, config)
	return []interface{}{transformed}
}
func flattenPrivatecaCertificateConfigX509ConfigKeyUsageBaseKeyUsage(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		v = make(map[string]interface{})
	}
	original := v.(map[string]interface{})
	transformed := make(map[string]interface{})
	transformed["digital_signature"] =
		flattenPrivatecaCertificateConfigX509ConfigKeyUsageBaseKeyUsageDigitalSignature(original["digitalSignature"], d, config)
	transformed["content_commitment"] =
		flattenPrivatecaCertificateConfigX509ConfigKeyUsageBaseKeyUsageContentCommitment(original["contentCommitment"], d, config)
	transformed["key_encipherment"] =
		flattenPrivatecaCertificateConfigX509ConfigKeyUsageBaseKeyUsageKeyEncipherment(original["keyEncipherment"], d, config)
	transformed["data_encipherment"] =
		flattenPrivatecaCertificateConfigX509ConfigKeyUsageBaseKeyUsageDataEncipherment(original["dataEncipherment"], d, config)
	transformed["key_agreement"] =
		flattenPrivatecaCertificateConfigX509ConfigKeyUsageBaseKeyUsageKeyAgreement(original["keyAgreement"], d, config)
	transformed["cert_sign"] =
		flattenPrivatecaCertificateConfigX509ConfigKeyUsageBaseKeyUsageCertSign(original["certSign"], d, config)
	transformed["crl_sign"] =
		flattenPrivatecaCertificateConfigX509ConfigKeyUsageBaseKeyUsageCrlSign(original["crlSign"], d, config)
	transformed["encipher_only"] =
		flattenPrivatecaCertificateConfigX509ConfigKeyUsageBaseKeyUsageEncipherOnly(original["encipherOnly"], d, config)
	transformed["decipher_only"] =
		flattenPrivatecaCertificateConfigX509ConfigKeyUsageBaseKeyUsageDecipherOnly(original["decipherOnly"], d, config)
	return []interface{}{transformed}
}
func flattenPrivatecaCertificateConfigX509ConfigKeyUsageBaseKeyUsageDigitalSignature(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigKeyUsageBaseKeyUsageContentCommitment(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigKeyUsageBaseKeyUsageKeyEncipherment(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigKeyUsageBaseKeyUsageDataEncipherment(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigKeyUsageBaseKeyUsageKeyAgreement(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigKeyUsageBaseKeyUsageCertSign(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigKeyUsageBaseKeyUsageCrlSign(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigKeyUsageBaseKeyUsageEncipherOnly(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigKeyUsageBaseKeyUsageDecipherOnly(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigKeyUsageExtendedKeyUsage(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		v = make(map[string]interface{})
	}
	original := v.(map[string]interface{})
	transformed := make(map[string]interface{})
	transformed["server_auth"] =
		flattenPrivatecaCertificateConfigX509ConfigKeyUsageExtendedKeyUsageServerAuth(original["serverAuth"], d, config)
	transformed["client_auth"] =
		flattenPrivatecaCertificateConfigX509ConfigKeyUsageExtendedKeyUsageClientAuth(original["clientAuth"], d, config)
	transformed["code_signing"] =
		flattenPrivatecaCertificateConfigX509ConfigKeyUsageExtendedKeyUsageCodeSigning(original["codeSigning"], d, config)
	transformed["email_protection"] =
		flattenPrivatecaCertificateConfigX509ConfigKeyUsageExtendedKeyUsageEmailProtection(original["emailProtection"], d, config)
	transformed["time_stamping"] =
		flattenPrivatecaCertificateConfigX509ConfigKeyUsageExtendedKeyUsageTimeStamping(original["timeStamping"], d, config)
	transformed["ocsp_signing"] =
		flattenPrivatecaCertificateConfigX509ConfigKeyUsageExtendedKeyUsageOcspSigning(original["ocspSigning"], d, config)
	return []interface{}{transformed}
}
func flattenPrivatecaCertificateConfigX509ConfigKeyUsageExtendedKeyUsageServerAuth(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigKeyUsageExtendedKeyUsageClientAuth(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigKeyUsageExtendedKeyUsageCodeSigning(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigKeyUsageExtendedKeyUsageEmailProtection(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigKeyUsageExtendedKeyUsageTimeStamping(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigKeyUsageExtendedKeyUsageOcspSigning(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages(v interface{}, d *schema.ResourceData, config *Config) interface{} {
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
			"object_id_path": flattenPrivatecaCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesObjectIdPath(original["objectIdPath"], d, config),
		})
	}
	return transformed
}
func flattenPrivatecaCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesObjectIdPath(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}
