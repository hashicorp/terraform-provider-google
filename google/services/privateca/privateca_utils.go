// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package privateca

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// This file contains shared flatteners between PrivateCA Certificate, CaPool and CertificateAuthority.
// These resources share the x509Config (Certificate, CertificateAuthorty)/baselineValues (CaPool) object.
// The API does not return this object if it only contains booleans with the default (false) value. This
// causes problems if a user specifies only default values, as Terraform detects that the object has been
// deleted on the API-side. This flattener creates default objects for sub-objects that match this pattern
// to fix perma-diffs on default-only objects. For this to work all objects that are flattened from nil to
// their default object *MUST* be set in the user's config, so they are all marked as Required in the schema.
//
// This file also contains shared expanders between the above resources. The expanders are required in order
// to handle the optional primitive field in CaOptions. By adding a virtual field, the expander can distinguish
// between an unset primitive field and a set primitive field with a default value.

// Expander utilities

func expandPrivatecaCertificateConfigX509ConfigCaOptions(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	// Fields non_ca, zero_max_issuer_path_length are used to distinguish between
	// unset booleans and booleans set with a default value.
	// Unset is_ca or unset max_issuer_path_length either allow any values for these fields when
	// used in an issuance policy, or allow the API to use default values when used in a
	// certificate config. A default value of is_ca=false means that issued certificates cannot
	// be CA certificates. A default value of max_issuer_path_length=0 means that the CA cannot
	// issue CA certificates.
	if v == nil {
		return nil, nil
	}
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})

	nonCa := original["non_ca"].(bool)
	isCa := original["is_ca"].(bool)

	zeroPathLength := original["zero_max_issuer_path_length"].(bool)
	maxIssuerPathLength := original["max_issuer_path_length"].(int)

	transformed := make(map[string]interface{})

	if nonCa && isCa {
		return nil, fmt.Errorf("non_ca, is_ca can not be set to true at the same time.")
	}
	if zeroPathLength && maxIssuerPathLength > 0 {
		return nil, fmt.Errorf("zero_max_issuer_path_length can not be set to true while max_issuer_path_length being set to a positive integer.")
	}

	if isCa || nonCa {
		transformed["isCa"] = original["is_ca"]
	}
	if maxIssuerPathLength > 0 || zeroPathLength {
		transformed["maxIssuerPathLength"] = original["max_issuer_path_length"]
	}
	return transformed, nil
}

func expandPrivatecaCertificateConfigX509ConfigKeyUsage(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	if v == nil {
		return v, nil
	}
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}

	raw := l[0]
	original := raw.(map[string]interface{})
	if len(original) == 0 {
		// Ignore empty KeyUsage
		return nil, nil
	}
	transformed := make(map[string]interface{})
	transformed["baseKeyUsage"] =
		expandPrivatecaCertificateConfigX509ConfigKeyUsageBaseKeyUsage(original["base_key_usage"], d, config)
	transformed["extendedKeyUsage"] =
		expandPrivatecaCertificateConfigX509ConfigKeyUsageExtendedKeyUsage(original["extended_key_usage"], d, config)
	transformed["unknownExtendedKeyUsages"] =
		expandPrivatecaCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages(original["unknown_extended_key_usages"], d, config)

	return transformed, nil
}

func expandPrivatecaCertificateConfigX509ConfigKeyUsageBaseKeyUsage(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	raw := l[0]
	original := raw.(map[string]interface{})
	if len(original) == 0 {
		// Ignore empty BaseKeyUsage
		return nil
	}

	transformed := make(map[string]interface{})
	transformed["digitalSignature"] = original["digital_signature"]
	transformed["contentCommitment"] = original["content_commitment"]
	transformed["keyEncipherment"] = original["key_encipherment"]
	transformed["dataEncipherment"] = original["data_encipherment"]
	transformed["keyAgreement"] = original["key_agreement"]
	transformed["certSign"] = original["cert_sign"]
	transformed["crlSign"] = original["crl_sign"]
	transformed["encipherOnly"] = original["encipher_only"]
	transformed["decipherOnly"] = original["decipher_only"]
	return transformed
}

func expandPrivatecaCertificateConfigX509ConfigKeyUsageExtendedKeyUsage(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	raw := l[0]
	original := raw.(map[string]interface{})
	if len(original) == 0 {
		// Ignore empty ExtendedKeyUsage
		return nil
	}

	transformed := make(map[string]interface{})
	transformed["serverAuth"] = original["server_auth"]
	transformed["clientAuth"] = original["client_auth"]
	transformed["codeSigning"] = original["code_signing"]
	transformed["emailProtection"] = original["email_protection"]
	transformed["timeStamping"] = original["time_stamping"]
	transformed["ocspSigning"] = original["ocsp_signing"]
	return transformed
}

func expandPrivatecaCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	// Parses the list of object IDs
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Ignore empty UnknownExtendedKeyUsages
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"objectIdPath": original["object_id_path"],
		})
	}
	return transformed
}

func expandPrivatecaCertificateConfigX509ConfigPolicyIds(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	if v == nil {
		return v, nil
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	// Parses the list of object IDs
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Ignore empty ObjectId
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"objectIdPath": original["object_id_path"],
		})
	}
	return transformed, nil
}

func expandPrivatecaCertificateConfigX509ConfigAdditionalExtensions(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	if v == nil {
		return v, nil
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Ignore empty AdditionalExtensions
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"critical": original["critical"],
			"value":    original["value"],
			"objectId": expandPrivatecaCertificateConfigX509ConfigAdditionalExtensionsObjectId(original["object_id"], d, config),
		})
	}
	return transformed, nil
}

func expandPrivatecaCertificateConfigX509ConfigAdditionalExtensionsObjectId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	// Expects a single object ID.
	raw := l[0]
	original := raw.(map[string]interface{})
	if len(original) == 0 {
		// Ignore empty ObjectId
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["objectIdPath"] = original["object_id_path"]
	return transformed
}

func expandPrivatecaCertificateConfigX509ConfigAiaOcspServers(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	// List of strings, no processing necessary.
	return v, nil
}

func expandPrivatecaCertificateConfigX509ConfigNameConstraints(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	if v == nil {
		return nil, nil
	}

	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}

	raw := l[0]
	original := raw.(map[string]interface{})
	if len(original) == 0 {
		// Ignore empty name constraints
		return nil, nil
	}

	transformed := make(map[string]interface{})
	transformed["critical"] = original["critical"]
	transformed["permittedDnsNames"] = original["permitted_dns_names"]
	transformed["excludedDnsNames"] = original["excluded_dns_names"]
	transformed["permittedIpRanges"] = original["permitted_ip_ranges"]
	transformed["excludedIpRanges"] = original["excluded_ip_ranges"]
	transformed["permittedEmailAddresses"] = original["permitted_email_addresses"]
	transformed["excludedEmailAddresses"] = original["excluded_email_addresses"]
	transformed["permittedUris"] = original["permitted_uris"]
	transformed["excludedUris"] = original["excluded_uris"]

	return transformed, nil
}

// Flattener utilities

func flattenPrivatecaCertificateConfigX509ConfigAdditionalExtensions(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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
func flattenPrivatecaCertificateConfigX509ConfigAdditionalExtensionsCritical(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigAdditionalExtensionsValue(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigAdditionalExtensionsObjectId(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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
func flattenPrivatecaCertificateConfigX509ConfigAdditionalExtensionsObjectIdObjectIdPath(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigPolicyIds(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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
func flattenPrivatecaCertificateConfigX509ConfigPolicyIdsObjectIdPath(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigAiaOcspServers(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigCaOptions(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Special case here as the CaPool API returns an empty object rather than nil unlike the Certificate
	// and CertificateAuthority APIs.
	if v == nil || len(v.(map[string]interface{})) == 0 {
		v = make(map[string]interface{})
	}
	original := v.(map[string]interface{})
	transformed := make(map[string]interface{})

	val, exists := original["isCa"]
	transformed["is_ca"] =
		flattenPrivatecaCertificateConfigX509ConfigCaOptionsIsCa(val, d, config)
	if exists && !val.(bool) {
		transformed["non_ca"] = true
	}

	val, exists = original["maxIssuerPathLength"]
	transformed["max_issuer_path_length"] =
		flattenPrivatecaCertificateConfigX509ConfigCaOptionsMaxIssuerPathLength(val, d, config)
	if exists && int(val.(float64)) == 0 {
		transformed["zero_max_issuer_path_length"] = true
	}

	return []interface{}{transformed}
}
func flattenPrivatecaCertificateConfigX509ConfigCaOptionsIsCa(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigCaOptionsMaxIssuerPathLength(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenPrivatecaCertificateConfigX509ConfigKeyUsage(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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
func flattenPrivatecaCertificateConfigX509ConfigKeyUsageBaseKeyUsage(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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
func flattenPrivatecaCertificateConfigX509ConfigKeyUsageBaseKeyUsageDigitalSignature(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}
func flattenPrivatecaCertificateConfigX509ConfigKeyUsageBaseKeyUsageContentCommitment(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}
func flattenPrivatecaCertificateConfigX509ConfigKeyUsageBaseKeyUsageKeyEncipherment(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigKeyUsageBaseKeyUsageDataEncipherment(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigKeyUsageBaseKeyUsageKeyAgreement(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigKeyUsageBaseKeyUsageCertSign(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigKeyUsageBaseKeyUsageCrlSign(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigKeyUsageBaseKeyUsageEncipherOnly(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigKeyUsageBaseKeyUsageDecipherOnly(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigKeyUsageExtendedKeyUsage(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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
func flattenPrivatecaCertificateConfigX509ConfigKeyUsageExtendedKeyUsageServerAuth(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigKeyUsageExtendedKeyUsageClientAuth(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigKeyUsageExtendedKeyUsageCodeSigning(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigKeyUsageExtendedKeyUsageEmailProtection(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigKeyUsageExtendedKeyUsageTimeStamping(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigKeyUsageExtendedKeyUsageOcspSigning(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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
func flattenPrivatecaCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesObjectIdPath(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenPrivatecaCertificateConfigX509ConfigNameConstraints(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformed["critical"] = original["critical"]
	transformed["permitted_dns_names"] = original["permittedDnsNames"]
	transformed["excluded_dns_names"] = original["excludedDnsNames"]
	transformed["permitted_ip_ranges"] = original["permittedIpRanges"]
	transformed["excluded_ip_ranges"] = original["excludedIpRanges"]
	transformed["permitted_email_addresses"] = original["permittedEmailAddresses"]
	transformed["excluded_email_addresses"] = original["excludedEmailAddresses"]
	transformed["permitted_uris"] = original["permittedUris"]
	transformed["excluded_uris"] = original["excludedUris"]

	return []interface{}{transformed}
}
