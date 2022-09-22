// Contains common diff suppress functions.

package google

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func optionalPrefixSuppress(prefix string) schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		return prefix+old == new || prefix+new == old
	}
}

func ignoreMissingKeyInMap(key string) schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		log.Printf("[DEBUG] - suppressing diff %q with old %q, new %q", k, old, new)
		if strings.HasSuffix(k, ".%") {
			oldNum, err := strconv.Atoi(old)
			if err != nil {
				log.Printf("[ERROR] could not parse %q as number, no longer attempting diff suppress", old)
				return false
			}
			newNum, err := strconv.Atoi(new)
			if err != nil {
				log.Printf("[ERROR] could not parse %q as number, no longer attempting diff suppress", new)
				return false
			}
			return oldNum+1 == newNum
		} else if strings.HasSuffix(k, "."+key) {
			return old == ""
		}
		return false
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

// Suppress diffs for blocks where one version is completely unset and the other is set
// to an empty block. This might occur in situations where removing a block completely
// is impossible (if it's computed or part of an AtLeastOneOf), so instead the user sets
// its values to empty.
// NOTE: Using Optional + Computed is *strongly* preferred to this DSF, as it's
// more well understood and resilient to API changes.
func emptyOrUnsetBlockDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	o, n := d.GetChange(strings.TrimSuffix(k, ".#"))
	return emptyOrUnsetBlockDiffSuppressLogic(k, old, new, o, n)
}

// The core logic for emptyOrUnsetBlockDiffSuppress, in a format that is more conducive
// to unit testing.
func emptyOrUnsetBlockDiffSuppressLogic(k, old, new string, o, n interface{}) bool {
	if !strings.HasSuffix(k, ".#") {
		return false
	}
	var l []interface{}
	if old == "0" && new == "1" {
		l = n.([]interface{})
	} else if new == "0" && old == "1" {
		l = o.([]interface{})
	} else {
		// we don't have one set and one unset, so don't suppress the diff
		return false
	}

	contents, ok := l[0].(map[string]interface{})
	if !ok {
		return false
	}
	for _, v := range contents {
		if !isEmptyValue(reflect.ValueOf(v)) {
			return false
		}
	}
	return true
}

// Suppress diffs for values that are equivalent except for their use of the words "location"
// compared to "region" or "zone"
func locationDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	return locationDiffSuppressHelper(old, new) || locationDiffSuppressHelper(new, old)
}

func locationDiffSuppressHelper(a, b string) bool {
	return strings.Replace(a, "/locations/", "/regions/", 1) == b ||
		strings.Replace(a, "/locations/", "/zones/", 1) == b
}

// For managed SSL certs, if new is an absolute FQDN (trailing '.') but old isn't, treat them as equals.
func absoluteDomainSuppress(k, old, new string, _ *schema.ResourceData) bool {
	if strings.HasPrefix(k, "managed.0.domains.") {
		return old == strings.TrimRight(new, ".") || new == strings.TrimRight(old, ".")
	}
	return false
}

func timestampDiffSuppress(format string) schema.SchemaDiffSuppressFunc {
	return func(_, old, new string, _ *schema.ResourceData) bool {
		oldT, err := time.Parse(format, old)
		if err != nil {
			return false
		}

		newT, err := time.Parse(format, new)
		if err != nil {
			return false
		}

		return oldT == newT
	}
}

// suppress diff when saved is Ipv4 format while new is required a reference
// this happens for an internal ip for Private Services Connect
func internalIpDiffSuppress(_, old, new string, _ *schema.ResourceData) bool {
	return (net.ParseIP(old) != nil) && (net.ParseIP(new) == nil)
}

// Suppress diffs for duration format. ex "60.0s" and "60s" same
// https://developers.google.com/protocol-buffers/docs/reference/google.protobuf#duration
func durationDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	oDuration, err := time.ParseDuration(old)
	if err != nil {
		return false
	}
	nDuration, err := time.ParseDuration(new)
	if err != nil {
		return false
	}
	return oDuration == nDuration
}

// Use this method when the field accepts either an IP address or a
// self_link referencing a resource (such as google_compute_route's
// next_hop_ilb)
func compareIpAddressOrSelfLinkOrResourceName(_, old, new string, _ *schema.ResourceData) bool {
	// if we can parse `new` as an IP address, then compare as strings
	if net.ParseIP(new) != nil {
		return new == old
	}

	// otherwise compare as self links
	return compareSelfLinkOrResourceName("", old, new, nil)
}

// Use this method when subnet is optioanl and auto_create_subnetworks = true
// API sometimes choose a subnet so the diff needs to be ignored
func compareOptionalSubnet(_, old, new string, _ *schema.ResourceData) bool {
	if isEmptyValue(reflect.ValueOf(new)) {
		return true
	}
	// otherwise compare as self links
	return compareSelfLinkOrResourceName("", old, new, nil)
}

// Suppress diffs in below cases
// "https://hello-rehvs75zla-uc.a.run.app/" -> "https://hello-rehvs75zla-uc.a.run.app"
// "https://hello-rehvs75zla-uc.a.run.app" -> "https://hello-rehvs75zla-uc.a.run.app/"
func lastSlashDiffSuppress(_, old, new string, _ *schema.ResourceData) bool {
	if last := len(new) - 1; last >= 0 && new[last] == '/' {
		new = new[:last]
	}

	if last := len(old) - 1; last >= 0 && old[last] == '/' {
		old = old[:last]
	}
	return new == old
}

// Suppress diffs when the value read from api
// has the project number instead of the project name
func projectNumberDiffSuppress(_, old, new string, _ *schema.ResourceData) bool {
	var a2, b2 string
	reN := regexp.MustCompile("projects/\\d+")
	re := regexp.MustCompile("projects/[^/]+")
	replacement := []byte("projects/equal")
	a2 = string(reN.ReplaceAll([]byte(old), replacement))
	b2 = string(re.ReplaceAll([]byte(new), replacement))
	return a2 == b2
}
