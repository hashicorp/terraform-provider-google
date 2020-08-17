// Contains common diff suppress functions.

package google

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func optionalPrefixSuppress(prefix string) schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		return prefix+old == new || prefix+new == old
	}
}

func optionalSurroundingSpacesSuppress(k, old, new string, d *schema.ResourceData) bool {
	return strings.TrimSpace(old) == strings.TrimSpace(new)
}

func emptyOrDefaultStringSuppress(defaultVal string) schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		return (old == "" && new == defaultVal) || (new == "" && old == defaultVal)
	}
}

func ipCidrRangeDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	// The range may be a:
	// A) single IP address (e.g. 10.2.3.4)
	// B) CIDR format string (e.g. 10.1.2.0/24)
	// C) netmask (e.g. /24)
	//
	// For A) and B), no diff to suppress, they have to match completely.
	// For C), The API picks a network IP address and this creates a diff of the form:
	// network_interface.0.alias_ip_range.0.ip_cidr_range: "10.128.1.0/24" => "/24"
	// We should only compare the mask portion for this case.
	if len(new) > 0 && new[0] == '/' {
		oldNetmaskStartPos := strings.LastIndex(old, "/")

		if oldNetmaskStartPos != -1 {
			oldNetmask := old[strings.LastIndex(old, "/"):]
			if oldNetmask == new {
				return true
			}
		}
	}

	return false
}

// sha256DiffSuppress
// if old is the hex-encoded sha256 sum of new, treat them as equal
func sha256DiffSuppress(_, old, new string, _ *schema.ResourceData) bool {
	return hex.EncodeToString(sha256.New().Sum([]byte(old))) == new
}

func caseDiffSuppress(_, old, new string, _ *schema.ResourceData) bool {
	return strings.ToUpper(old) == strings.ToUpper(new)
}

// Port range '80' and '80-80' is equivalent.
// `old` is read from the server and always has the full range format (e.g. '80-80', '1024-2048').
// `new` can be either a single port or a port range.
func portRangeDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	return old == new+"-"+new
}

// Single-digit hour is equivalent to hour with leading zero e.g. suppress diff 1:00 => 01:00.
// Assume either value could be in either format.
func rfc3339TimeDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	if (len(old) == 4 && "0"+old == new) || (len(new) == 4 && "0"+new == old) {
		return true
	}
	return false
}
