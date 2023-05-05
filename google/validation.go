package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/verify"
)

const (
	// Deprecated: For backward compatibility ProjectRegex is still working,
	// but all new code should use ProjectRegex in the verify package instead.
	// Copied from the official Google Cloud auto-generated client.
	ProjectRegex = verify.ProjectRegex
	// Deprecated: For backward compatibility ProjectRegexWildCard is still working,
	// but all new code should use ProjectRegexWildCard in the verify package instead.
	ProjectRegexWildCard = verify.ProjectRegexWildCard
	// Deprecated: For backward compatibility RegionRegex is still working,
	// but all new code should use RegionRegex in the verify package instead.
	RegionRegex = verify.RegionRegex
	// Deprecated: For backward compatibility SubnetworkRegex is still working,
	// but all new code should use SubnetworkRegex in the verify package instead.
	SubnetworkRegex = verify.SubnetworkRegex

	// Deprecated: For backward compatibility SubnetworkLinkRegex is still working,
	// but all new code should use SubnetworkLinkRegex in the verify package instead.
	SubnetworkLinkRegex = verify.SubnetworkLinkRegex

	// Deprecated: For backward compatibility RFC1035NameTemplate is still working,
	// but all new code should use RFC1035NameTemplate in the verify package instead.
	RFC1035NameTemplate = verify.RFC1035NameTemplate
	// Deprecated: For backward compatibility CloudIoTIdRegex is still working,
	// but all new code should use CloudIoTIdRegex in the verify package instead.
	CloudIoTIdRegex = verify.CloudIoTIdRegex

	// Deprecated: For backward compatibility ComputeServiceAccountNameRegex is still working,
	// but all new code should use ComputeServiceAccountNameRegex in the verify package instead.
	// Format of default Compute service accounts created by Google
	// ${PROJECT_ID}-compute@developer.gserviceaccount.com where PROJECT_ID is an int64 (max 20 digits)
	ComputeServiceAccountNameRegex = verify.ComputeServiceAccountNameRegex

	// Deprecated: For backward compatibility IAMCustomRoleIDRegex is still working,
	// but all new code should use IAMCustomRoleIDRegex in the verify package instead.
	// https://cloud.google.com/iam/docs/understanding-custom-roles#naming_the_role
	IAMCustomRoleIDRegex = verify.IAMCustomRoleIDRegex

	// Deprecated: For backward compatibility ADDomainNameRegex is still working,
	// but all new code should use ADDomainNameRegex in the verify package instead.
	// https://cloud.google.com/managed-microsoft-ad/reference/rest/v1/projects.locations.global.domains/create#query-parameters
	ADDomainNameRegex = verify.ADDomainNameRegex
)

var (
	// Deprecated: For backward compatibility ServiceAccountNameRegex is still working,
	// but all new code should use ServiceAccountNameRegex in the verify package instead.
	// Service account name must have a length between 6 and 30.
	// The first and last characters have different restrictions, than
	// the middle characters. The middle characters length must be between
	// 4 and 28 since the first and last character are excluded.
	ServiceAccountNameRegex = verify.ServiceAccountNameRegex

	// Deprecated: For backward compatibility ServiceAccountLinkRegexPrefix is still working,
	// but all new code should use ServiceAccountLinkRegexPrefix in the verify package instead.
	ServiceAccountLinkRegexPrefix = verify.ServiceAccountLinkRegexPrefix
	// Deprecated: For backward compatibility PossibleServiceAccountNames is still working,
	// but all new code should use PossibleServiceAccountNames in the verify package instead.
	PossibleServiceAccountNames = verify.PossibleServiceAccountNames
	// Deprecated: For backward compatibility ServiceAccountLinkRegex is still working,
	// but all new code should use ServiceAccountLinkRegex in the verify package instead.
	ServiceAccountLinkRegex = verify.ServiceAccountLinkRegex

	// Deprecated: For backward compatibility ServiceAccountKeyNameRegex is still working,
	// but all new code should use ServiceAccountKeyNameRegex in the verify package instead.
	ServiceAccountKeyNameRegex = verify.ServiceAccountKeyNameRegex

	// Deprecated: For backward compatibility CreatedServiceAccountNameRegex is still working,
	// but all new code should use CreatedServiceAccountNameRegex in the verify package instead.
	// Format of service accounts created through the API
	CreatedServiceAccountNameRegex = verify.CreatedServiceAccountNameRegex

	// Deprecated: For backward compatibility ServiceDefaultAccountNameRegex is still working,
	// but all new code should use ServiceDefaultAccountNameRegex in the verify package instead.
	// Format of service-created service account
	// examples are:
	// 		$PROJECTID@cloudbuild.gserviceaccount.com
	// 		$PROJECTID@cloudservices.gserviceaccount.com
	// 		$PROJECTID@appspot.gserviceaccount.com
	ServiceDefaultAccountNameRegex = verify.ServiceDefaultAccountNameRegex

	// Deprecated: For backward compatibility ProjectNameInDNSFormRegex is still working,
	// but all new code should use ProjectNameInDNSFormRegex in the verify package instead.
	ProjectNameInDNSFormRegex = verify.ProjectNameInDNSFormRegex
	// Deprecated: For backward compatibility ProjectNameRegex is still working,
	// but all new code should use ProjectNameRegex in the verify package instead.
	ProjectNameRegex = verify.ProjectNameRegex

	// Valid range for Cloud Router ASN values as per RFC6996
	// https://tools.ietf.org/html/rfc6996
	// Must be explicitly int64 to avoid overflow when building Terraform for 32bit architectures
	// Deprecated: For backward compatibility Rfc6996Asn16BitMin is still working,
	// but all new code should use Rfc6996Asn16BitMin in the verify package instead.
	Rfc6996Asn16BitMin = verify.Rfc6996Asn16BitMin
	// Deprecated: For backward compatibility Rfc6996Asn16BitMax is still working,
	// but all new code should use Rfc6996Asn16BitMax in the verify package instead.
	Rfc6996Asn16BitMax = verify.Rfc6996Asn16BitMax
	// Deprecated: For backward compatibility Rfc6996Asn32BitMin is still working,
	// but all new code should use Rfc6996Asn32BitMin in the verify package instead.
	Rfc6996Asn32BitMin = verify.Rfc6996Asn32BitMin
	// Deprecated: For backward compatibility Rfc6996Asn32BitMax is still working,
	// but all new code should use Rfc6996Asn32BitMax in the verify package instead.
	Rfc6996Asn32BitMax = verify.Rfc6996Asn32BitMax
	// Deprecated: For backward compatibility GcpRouterPartnerAsn is still working,
	// but all new code should use Rfc6996Asn16BitMin in the verify package instead.
	GcpRouterPartnerAsn = verify.GcpRouterPartnerAsn
)

// Deprecated: For backward compatibility rfc1918Networks is still working,
// but all new code should use Rfc1918Networks in the verify package instead.
var rfc1918Networks = verify.Rfc1918Networks

// Deprecated: For backward compatibility validateGCEName is still working,
// but all new code should use ValidateGCEName in the verify package instead.
// validateGCEName ensures that a field matches the requirements for Compute Engine resource names
// https://cloud.google.com/compute/docs/naming-resources#resource-name-format
func validateGCEName(v interface{}, k string) (ws []string, errors []error) {
	return verify.ValidateGCEName(v, k)
}

// Deprecated: For backward compatibility validateRFC6996Asn is still working,
// but all new code should use ValidateRFC6996Asn in the verify package instead.
// Ensure that the BGP ASN value of Cloud Router is a valid value as per RFC6996 or a value of 16550
func validateRFC6996Asn(v interface{}, k string) (ws []string, errors []error) {
	return verify.ValidateRFC6996Asn(v, k)
}

// Deprecated: For backward compatibility validateRegexp is still working,
// but all new code should use ValidateRegexp in the verify package instead.
func validateRegexp(re string) schema.SchemaValidateFunc {
	return verify.ValidateRegexp(re)
}

// Deprecated: For backward compatibility validateEnum is still working,
// but all new code should use ValidateEnum in the verify package instead.
func validateEnum(values []string) schema.SchemaValidateFunc {
	return verify.ValidateEnum(values)
}

// Deprecated: For backward compatibility validateRFC1918Network is still working,
// but all new code should use ValidateRFC1918Network in the verify package instead.
func validateRFC1918Network(min, max int) schema.SchemaValidateFunc {
	return verify.ValidateRFC1918Network(min, max)
}

// Deprecated: For backward compatibility validateRFC3339Time is still working,
// but all new code should use ValidateRFC3339Time in the verify package instead.
func validateRFC3339Time(v interface{}, k string) (warnings []string, errors []error) {
	return verify.ValidateRFC3339Time(v, k)
}

// Deprecated: For backward compatibility validateRFC1035Name is still working,
// but all new code should use ValidateRFC1035Name in the verify package instead.
func validateRFC1035Name(min, max int) schema.SchemaValidateFunc {
	return verify.ValidateRFC1035Name(min, max)
}

// Deprecated: For backward compatibility validateIpCidrRange is still working,
// but all new code should use ValidateIpCidrRange in the verify package instead.
func validateIpCidrRange(v interface{}, k string) (warnings []string, errors []error) {
	return verify.ValidateIpCidrRange(v, k)
}

// Deprecated: For backward compatibility validateIAMCustomRoleID is still working,
// but all new code should use ValidateIAMCustomRoleID in the verify package instead.
func validateIAMCustomRoleID(v interface{}, k string) (warnings []string, errors []error) {
	return verify.ValidateIAMCustomRoleID(v, k)
}

// Deprecated: For backward compatibility orEmpty is still working,
// but all new code should use OrEmpty in the verify package instead.
func orEmpty(f schema.SchemaValidateFunc) schema.SchemaValidateFunc {
	return verify.OrEmpty(f)
}

// Deprecated: For backward compatibility validateProjectID is still working,
// but all new code should use ValidateProjectID in the verify package instead.
func validateProjectID() schema.SchemaValidateFunc {
	return verify.ValidateProjectID()
}

// Deprecated: For backward compatibility validateDSProjectID is still working,
// but all new code should use ValidateDSProjectID in the verify package instead.
func validateDSProjectID() schema.SchemaValidateFunc {
	return verify.ValidateDSProjectID()
}

// Deprecated: For backward compatibility validateProjectName is still working,
// but all new code should use ValidateProjectName in the verify package instead.
func validateProjectName() schema.SchemaValidateFunc {
	return verify.ValidateProjectName()
}

// Deprecated: For backward compatibility validateDuration is still working,
// but all new code should use ValidateDuration in the verify package instead.
func validateDuration() schema.SchemaValidateFunc {
	return verify.ValidateDuration()
}

// Deprecated: For backward compatibility validateNonNegativeDuration is still working,
// but all new code should use ValidateNonNegativeDuration in the verify package instead.
func validateNonNegativeDuration() schema.SchemaValidateFunc {
	return verify.ValidateNonNegativeDuration()
}

// Deprecated: For backward compatibility validateIpAddress is still working,
// but all new code should use ValidateIpAddress in the verify package instead.
func validateIpAddress(i interface{}, val string) ([]string, []error) {
	return verify.ValidateIpAddress(i, val)
}

// Deprecated: For backward compatibility validateBase64String is still working,
// but all new code should use ValidateBase64String in the verify package instead.
func validateBase64String(i interface{}, val string) ([]string, []error) {
	return verify.ValidateBase64String(i, val)
}

// Deprecated: For backward compatibility StringNotInSlice is still working,
// but all new code should use StringNotInSlice in the verify package instead.
// StringNotInSlice returns a SchemaValidateFunc which tests if the provided value
// is of type string and that it matches none of the element in the invalid slice.
// if ignorecase is true, case is ignored.
func StringNotInSlice(invalid []string, ignoreCase bool) schema.SchemaValidateFunc {
	return verify.StringNotInSlice(invalid, ignoreCase)
}

// Deprecated: For backward compatibility validateHourlyOnly is still working,
// but all new code should use ValidateHourlyOnly in the verify package instead.
// Ensure that hourly timestamp strings "HH:MM" have the minutes zeroed out for hourly only inputs
func validateHourlyOnly(val interface{}, key string) (warns []string, errs []error) {
	return verify.ValidateHourlyOnly(val, key)
}

// Deprecated: For backward compatibility validateRFC3339Date is still working,
// but all new code should use ValidateRFC3339Date in the verify package instead.
func validateRFC3339Date(v interface{}, k string) (warnings []string, errors []error) {
	return verify.ValidateRFC3339Date(v, k)
}

// Deprecated: For backward compatibility validateADDomainName is still working,
// but all new code should use ValidateADDomainName in the verify package instead.
func validateADDomainName() schema.SchemaValidateFunc {
	return verify.ValidateADDomainName()
}
