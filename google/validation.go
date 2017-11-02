package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"net"
	"regexp"
)

const (
	// Copied from the official Google Cloud auto-generated client.
	ProjectRegex    = "(?:(?:[-a-z0-9]{1,63}\\.)*(?:[a-z](?:[-a-z0-9]{0,61}[a-z0-9])?):)?(?:[0-9]{1,19}|(?:[a-z0-9](?:[-a-z0-9]{0,61}[a-z0-9])?))"
	RegionRegex     = "[a-z](?:[-a-z0-9]{0,61}[a-z0-9])?"
	SubnetworkRegex = "[a-z](?:[-a-z0-9]{0,61}[a-z0-9])?"

	SubnetworkLinkRegex = "projects/(" + ProjectRegex + ")/regions/(" + RegionRegex + ")/subnetworks/(" + SubnetworkRegex + ")$"
)

var rfc1918Networks = []string{
	"10.0.0.0/8",
	"172.16.0.0/12",
	"192.168.0.0/16",
}

func validateGCPName(v interface{}, k string) (ws []string, errors []error) {
	re := `^(?:[a-z](?:[-a-z0-9]{0,61}[a-z0-9])?)$`
	return validateRegexp(re)(v, k)
}

func validateRegexp(re string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(string)
		if !regexp.MustCompile(re).MatchString(value) {
			errors = append(errors, fmt.Errorf(
				"%q (%q) doesn't match regexp %q", k, value, re))
		}

		return
	}
}

func validateRFC1918Network(min, max int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {

		s, es = validation.CIDRNetwork(min, max)(i, k)
		if len(es) > 0 {
			return
		}

		v, _ := i.(string)
		ip, _, _ := net.ParseCIDR(v)
		for _, c := range rfc1918Networks {
			if _, ipnet, _ := net.ParseCIDR(c); ipnet.Contains(ip) {
				return
			}
		}

		es = append(es, fmt.Errorf("expected %q to be an RFC1918-compliant CIDR, got: %s", k, v))

		return
	}
}
