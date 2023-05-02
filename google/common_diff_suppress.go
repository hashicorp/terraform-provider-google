// Contains common diff suppress functions.

package google

import (
	"net"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

// Deprecated: For backward compatibility OptionalPrefixSuppress is still working,
// but all new code should use OptionalPrefixSuppress in the tpgresource package instead.
func OptionalPrefixSuppress(prefix string) schema.SchemaDiffSuppressFunc {
	return tpgresource.OptionalPrefixSuppress(prefix)
}

// Deprecated: For backward compatibility IgnoreMissingKeyInMap is still working,
// but all new code should use IgnoreMissingKeyInMap in the tpgresource package instead.
func IgnoreMissingKeyInMap(key string) schema.SchemaDiffSuppressFunc {
	return tpgresource.IgnoreMissingKeyInMap(key)
}

// Deprecated: For backward compatibility OptionalSurroundingSpacesSuppress is still working,
// but all new code should use OptionalSurroundingSpacesSuppress in the tpgresource package instead.
func OptionalSurroundingSpacesSuppress(k, old, new string, d *schema.ResourceData) bool {
	return tpgresource.OptionalSurroundingSpacesSuppress(k, old, new, d)
}

// Deprecated: For backward compatibility EmptyOrDefaultStringSuppress is still working,
// but all new code should use EmptyOrDefaultStringSuppress in the tpgresource package instead.
func EmptyOrDefaultStringSuppress(defaultVal string) schema.SchemaDiffSuppressFunc {
	return tpgresource.EmptyOrDefaultStringSuppress(defaultVal)
}

// Deprecated: For backward compatibility IpCidrRangeDiffSuppress is still working,
// but all new code should use IpCidrRangeDiffSuppress in the tpgresource package instead.
func IpCidrRangeDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	return tpgresource.IpCidrRangeDiffSuppress(k, old, new, d)
}

// Deprecated: For backward compatibility Sha256DiffSuppress is still working,
// but all new code should use Sha256DiffSuppress in the tpgresource package instead.
// Sha256DiffSuppress
// if old is the hex-encoded sha256 sum of new, treat them as equal
func Sha256DiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	return tpgresource.Sha256DiffSuppress(k, old, new, d)
}

// Deprecated: For backward compatibility CaseDiffSuppress is still working,
// but all new code should use CaseDiffSuppress in the tpgresource package instead.
func CaseDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	return tpgresource.CaseDiffSuppress(k, old, new, d)
}

// Deprecated: For backward compatibility PortRangeDiffSuppress is still working,
// but all new code should use PortRangeDiffSuppress in the tpgresource package instead.
// Port range '80' and '80-80' is equivalent.
// `old` is read from the server and always has the full range format (e.g. '80-80', '1024-2048').
// `new` can be either a single port or a port range.
func PortRangeDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	return tpgresource.PortRangeDiffSuppress(k, old, new, d)
}

// Deprecated: For backward compatibility Rfc3339TimeDiffSuppress is still working,
// but all new code should use Rfc3339TimeDiffSuppress in the tpgresource package instead.
// Single-digit hour is equivalent to hour with leading zero e.g. suppress diff 1:00 => 01:00.
// Assume either value could be in either format.
func Rfc3339TimeDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	return tpgresource.Rfc3339TimeDiffSuppress(k, old, new, d)
}

// Deprecated: For backward compatibility EmptyOrUnsetBlockDiffSuppress is still working,
// but all new code should use EmptyOrUnsetBlockDiffSuppress in the tpgresource package instead.
// Suppress diffs for blocks where one version is completely unset and the other is set
// to an empty block. This might occur in situations where removing a block completely
// is impossible (if it's computed or part of an AtLeastOneOf), so instead the user sets
// its values to empty.
// NOTE: Using Optional + Computed is *strongly* preferred to this DSF, as it's
// more well understood and resilient to API changes.
func EmptyOrUnsetBlockDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	return tpgresource.EmptyOrUnsetBlockDiffSuppress(k, old, new, d)
}

// Deprecated: For backward compatibility EmptyOrUnsetBlockDiffSuppressLogic is still working,
// but all new code should use EmptyOrUnsetBlockDiffSuppressLogic in the tpgresource package instead.
// The core logic for EmptyOrUnsetBlockDiffSuppress, in a format that is more conducive
// to unit testing.
func EmptyOrUnsetBlockDiffSuppressLogic(k, old, new string, o, n interface{}) bool {
	return tpgresource.EmptyOrUnsetBlockDiffSuppressLogic(k, old, new, o, n)
}

// Deprecated: For backward compatibility LocationDiffSuppress is still working,
// but all new code should use LocationDiffSuppress in the tpgresource package instead.
// Suppress diffs for values that are equivalent except for their use of the words "location"
// compared to "region" or "zone"
func LocationDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	return tpgresource.LocationDiffSuppress(k, old, new, d)
}

// Deprecated: For backward compatibility locationDiffSuppressHelper is still working,
// but all new code should use LocationDiffSuppressHelper in the tpgresource package instead.
func locationDiffSuppressHelper(a, b string) bool {
	return tpgresource.LocationDiffSuppressHelper(a, b)
}

// Deprecated: For backward compatibility AbsoluteDomainSuppress is still working,
// but all new code should use AbsoluteDomainSuppress in the tpgresource package instead.
// For managed SSL certs, if new is an absolute FQDN (trailing '.') but old isn't, treat them as equals.
func AbsoluteDomainSuppress(k, old, new string, d *schema.ResourceData) bool {
	return tpgresource.AbsoluteDomainSuppress(k, old, new, d)
}

// Deprecated: For backward compatibility TimestampDiffSuppress is still working,
// but all new code should use TimestampDiffSuppress in the tpgresource package instead.
func TimestampDiffSuppress(format string) schema.SchemaDiffSuppressFunc {
	return tpgresource.TimestampDiffSuppress(format)
}

// Deprecated: For backward compatibility InternalIpDiffSuppress is still working,
// but all new code should use InternalIpDiffSuppress in the tpgresource package instead.
// suppress diff when saved is Ipv4 format while new is required a reference
// this happens for an internal ip for Private Services Connect
func InternalIpDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	return tpgresource.InternalIpDiffSuppress(k, old, new, d)
}

// Deprecated: For backward compatibility DurationDiffSuppress is still working,
// but all new code should use DurationDiffSuppress in the tpgresource package instead.
// Suppress diffs for duration format. ex "60.0s" and "60s" same
// https://developers.google.com/protocol-buffers/docs/reference/google.protobuf#duration
func DurationDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	return tpgresource.DurationDiffSuppress(k, old, new, d)
}

// Use this method when the field accepts either an IP address or a
// self_link referencing a resource (such as google_compute_route's
// next_hop_ilb)
func CompareIpAddressOrSelfLinkOrResourceName(_, old, new string, _ *schema.ResourceData) bool {
	// if we can parse `new` as an IP address, then compare as strings
	if net.ParseIP(new) != nil {
		return new == old
	}

	// otherwise compare as self links
	return compareSelfLinkOrResourceName("", old, new, nil)
}

// Use this method when subnet is optioanl and auto_create_subnetworks = true
// API sometimes choose a subnet so the diff needs to be ignored
func CompareOptionalSubnet(_, old, new string, _ *schema.ResourceData) bool {
	if tpgresource.IsEmptyValue(reflect.ValueOf(new)) {
		return true
	}
	// otherwise compare as self links
	return compareSelfLinkOrResourceName("", old, new, nil)
}

// Deprecated: For backward compatibility LastSlashDiffSuppress is still working,
// but all new code should use LastSlashDiffSuppress in the tpgresource package instead.
// Suppress diffs in below cases
// "https://hello-rehvs75zla-uc.a.run.app/" -> "https://hello-rehvs75zla-uc.a.run.app"
// "https://hello-rehvs75zla-uc.a.run.app" -> "https://hello-rehvs75zla-uc.a.run.app/"
func LastSlashDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	return tpgresource.LastSlashDiffSuppress(k, old, new, d)
}

// Deprecated: For backward compatibility ProjectNumberDiffSuppress is still working,
// but all new code should use ProjectNumberDiffSuppress in the tpgresource package instead.
// Suppress diffs when the value read from api
// has the project number instead of the project name
func ProjectNumberDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	return tpgresource.ProjectNumberDiffSuppress(k, old, new, d)
}
