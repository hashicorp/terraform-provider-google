// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Contains common diff suppress functions.

package tpgresource

import "testing"

func TestOptionalPrefixSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		Prefix             string
		ExpectDiffSuppress bool
	}{
		"with same prefix": {
			Old:                "my-folder",
			New:                "folders/my-folder",
			Prefix:             "folders/",
			ExpectDiffSuppress: true,
		},
		"with different prefix": {
			Old:                "folders/my-folder",
			New:                "organizations/my-folder",
			Prefix:             "folders/",
			ExpectDiffSuppress: false,
		},
		"same without prefix": {
			Old:                "my-folder",
			New:                "my-folder",
			Prefix:             "folders/",
			ExpectDiffSuppress: false,
		},
		"different without prefix": {
			Old:                "my-folder",
			New:                "my-new-folder",
			Prefix:             "folders/",
			ExpectDiffSuppress: false,
		},
	}

	for tn, tc := range cases {
		if OptionalPrefixSuppress(tc.Prefix)("folder", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Fatalf("bad: %s, '%s' => '%s' expect %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}

func TestIgnoreMissingKeyInMap(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		Key                string
		ExpectDiffSuppress bool
	}{
		"missing key in map": {
			Old:                "",
			New:                "v1",
			Key:                "x-goog-version",
			ExpectDiffSuppress: true,
		},
		"different values": {
			Old:                "v1",
			New:                "v2",
			Key:                "x-goog-version",
			ExpectDiffSuppress: false,
		},
		"same values": {
			Old:                "v1",
			New:                "v1",
			Key:                "x-goog-version",
			ExpectDiffSuppress: false,
		},
	}

	for tn, tc := range cases {
		if IgnoreMissingKeyInMap(tc.Key)("push_config.0.attributes."+tc.Key, tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Fatalf("bad: %s, '%s' => '%s' expect %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}

func TestOptionalSurroundingSpacesSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"surrounding spaces": {
			Old:                "value",
			New:                " value ",
			ExpectDiffSuppress: true,
		},
		"no surrounding spaces": {
			Old:                "value",
			New:                "value",
			ExpectDiffSuppress: true,
		},
		"one space each": {
			Old:                " value",
			New:                "value ",
			ExpectDiffSuppress: true,
		},
		"different values": {
			Old:                " different",
			New:                "values ",
			ExpectDiffSuppress: false,
		},
	}

	for tn, tc := range cases {
		if OptionalSurroundingSpacesSuppress("filter", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Fatalf("bad: %s, '%s' => '%s' expect %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}

func TestCaseDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"differents cases": {
			Old:                "Value",
			New:                "value",
			ExpectDiffSuppress: true,
		},
		"different values": {
			Old:                "value",
			New:                "NewValue",
			ExpectDiffSuppress: false,
		},
		"same cases": {
			Old:                "value",
			New:                "value",
			ExpectDiffSuppress: true,
		},
	}

	for tn, tc := range cases {
		if CaseDiffSuppress("key", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Fatalf("bad: %s, '%s' => '%s' expect %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}

func TestPortRangeDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"different single values": {
			Old:                "80-80",
			New:                "443",
			ExpectDiffSuppress: false,
		},
		"different ranges": {
			Old:                "80-80",
			New:                "443-444",
			ExpectDiffSuppress: false,
		},
		"same single values": {
			Old:                "80-80",
			New:                "80",
			ExpectDiffSuppress: true,
		},
		"same ranges": {
			Old:                "80-80",
			New:                "80-80",
			ExpectDiffSuppress: false,
		},
	}

	for tn, tc := range cases {
		if PortRangeDiffSuppress("ports", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Fatalf("bad: %s, '%s' => '%s' expect %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}

func TestLocationDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"locations to zones": {
			Old:                "projects/x/locations/y/resource/z",
			New:                "projects/x/zones/y/resource/z",
			ExpectDiffSuppress: true,
		},
		"regions to locations": {
			Old:                "projects/x/regions/y/resource/z",
			New:                "projects/x/locations/y/resource/z",
			ExpectDiffSuppress: true,
		},
		"locations to locations": {
			Old:                "projects/x/locations/y/resource/z",
			New:                "projects/x/locations/y/resource/z",
			ExpectDiffSuppress: false,
		},
		"zones to regions": {
			Old:                "projects/x/zones/y/resource/z",
			New:                "projects/x/regions/y/resource/z",
			ExpectDiffSuppress: false,
		},
		"different locations": {
			Old:                "projects/x/locations/a/resource/z",
			New:                "projects/x/locations/b/resource/z",
			ExpectDiffSuppress: false,
		},
	}

	for tn, tc := range cases {
		if LocationDiffSuppress("policy_uri", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Fatalf("bad: %s, '%s' => '%s' expect %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}

func TestAbsoluteDomainSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"new trailing dot": {
			Old:                "sslcert.tf-test.club",
			New:                "sslcert.tf-test.club.",
			ExpectDiffSuppress: true,
		},
		"old trailing dot": {
			Old:                "sslcert.tf-test.club.",
			New:                "sslcert.tf-test.club",
			ExpectDiffSuppress: true,
		},
		"same trailing dot": {
			Old:                "sslcert.tf-test.club.",
			New:                "sslcert.tf-test.club.",
			ExpectDiffSuppress: false,
		},
		"different trailing dot": {
			Old:                "sslcert.tf-test.club.",
			New:                "sslcert.tf-test.clubs.",
			ExpectDiffSuppress: false,
		},
		"different no trailing dot": {
			Old:                "sslcert.tf-test.club",
			New:                "sslcert.tf-test.clubs",
			ExpectDiffSuppress: false,
		},
	}

	for tn, tc := range cases {
		if AbsoluteDomainSuppress("managed.0.domains.", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Fatalf("bad: %s, '%s' => '%s' expect %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}

func TestDurationDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"different values": {
			Old:                "60s",
			New:                "65s",
			ExpectDiffSuppress: false,
		},
		"same values": {
			Old:                "60s",
			New:                "60s",
			ExpectDiffSuppress: true,
		},
		"different values, different formats": {
			Old:                "65s",
			New:                "60.0s",
			ExpectDiffSuppress: false,
		},
		"same values, different formats": {
			Old:                "60.0s",
			New:                "60s",
			ExpectDiffSuppress: true,
		},
	}

	for tn, tc := range cases {
		if DurationDiffSuppress("duration", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Fatalf("bad: %s, '%s' => '%s' expect %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}

func TestInternalIpDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"suppress - same long and short ipv6 IPs without netmask": {
			Old:                "2600:1900:4020:31cd:8000:0:0:0",
			New:                "2600:1900:4020:31cd:8000::",
			ExpectDiffSuppress: true,
		},
		"suppress - long and short ipv6 IPs with netmask": {
			Old:                "2600:1900:4020:31cd:8000:0:0:0/96",
			New:                "2600:1900:4020:31cd:8000::/96",
			ExpectDiffSuppress: true,
		},
		"suppress - long ipv6 IP with netmask and short ipv6 IP without netmask": {
			Old:                "2600:1900:4020:31cd:8000:0:0:0/96",
			New:                "2600:1900:4020:31cd:8000::",
			ExpectDiffSuppress: true,
		},
		"suppress - long ipv6 IP without netmask and short ipv6 IP with netmask": {
			Old:                "2600:1900:4020:31cd:8000:0:0:0",
			New:                "2600:1900:4020:31cd:8000::/96",
			ExpectDiffSuppress: true,
		},
		"suppress - long ipv6 IP with netmask and reference": {
			Old:                "2600:1900:4020:31cd:8000:0:0:0/96",
			New:                "projects/project_id/regions/region/addresses/address-name",
			ExpectDiffSuppress: true,
		},
		"suppress - long ipv6 IP without netmask and reference": {
			Old:                "2600:1900:4020:31cd:8000:0:0:0",
			New:                "projects/project_id/regions/region/addresses/address-name",
			ExpectDiffSuppress: true,
		},
		"do not suppress - ipv6 IPs different netmask": {
			Old:                "2600:1900:4020:31cd:8000:0:0:0/96",
			New:                "2600:1900:4020:31cd:8000:0:0:0/95",
			ExpectDiffSuppress: false,
		},
		"do not suppress - reference and ipv6 IP with netmask": {
			Old:                "projects/project_id/regions/region/addresses/address-name",
			New:                "2600:1900:4020:31cd:8000:0:0:0/96",
			ExpectDiffSuppress: false,
		},
		"do not suppress - ipv6 IPs - 1": {
			Old:                "2600:1900:4020:31cd:8000:0:0:0",
			New:                "2600:1900:4020:31cd:8001::",
			ExpectDiffSuppress: false,
		},
		"do not suppress - ipv6 IPs - 2": {
			Old:                "2600:1900:4020:31cd:8000:0:0:0",
			New:                "2600:1900:4020:31cd:8000:0:0:8000",
			ExpectDiffSuppress: false,
		},
		"suppress - ipv4 IPs": {
			Old:                "1.2.3.4",
			New:                "1.2.3.4",
			ExpectDiffSuppress: true,
		},
		"suppress - ipv4 IP without netmask and ipv4 IP with netmask": {
			Old:                "1.2.3.4",
			New:                "1.2.3.4/24",
			ExpectDiffSuppress: true,
		},
		"suppress - ipv4 IP without netmask and reference": {
			Old:                "1.2.3.4",
			New:                "projects/project_id/regions/region/addresses/address-name",
			ExpectDiffSuppress: true,
		},
		"do not suppress - reference and ipv4 IP without netmask": {
			Old:                "projects/project_id/regions/region/addresses/address-name",
			New:                "1.2.3.4",
			ExpectDiffSuppress: false,
		},
		"do not suppress - different ipv4 IPs": {
			Old:                "1.2.3.4",
			New:                "1.2.3.5",
			ExpectDiffSuppress: false,
		},
		"do not suppress - ipv4 IPs different netmask": {
			Old:                "1.2.3.4/24",
			New:                "1.2.3.5/25",
			ExpectDiffSuppress: false,
		},
		"do not suppress - different references": {
			Old:                "projects/project_id/regions/region/addresses/address-name",
			New:                "projects/project_id/regions/region/addresses/address-name-1",
			ExpectDiffSuppress: false,
		},
		"do not suppress - same references": {
			Old:                "projects/project_id/regions/region/addresses/address-name",
			New:                "projects/project_id/regions/region/addresses/address-name",
			ExpectDiffSuppress: false,
		},
	}

	for tn, tc := range cases {
		if InternalIpDiffSuppress("ipv4/v6_compare", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Fatalf("bad: %s, '%s' => '%s' expect %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}

func TestLastSlashDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"slash to no slash": {
			Old:                "https://hello-rehvs75zla-uc.a.run.app/",
			New:                "https://hello-rehvs75zla-uc.a.run.app",
			ExpectDiffSuppress: true,
		},
		"no slash to slash": {
			Old:                "https://hello-rehvs75zla-uc.a.run.app",
			New:                "https://hello-rehvs75zla-uc.a.run.app/",
			ExpectDiffSuppress: true,
		},
		"slash to slash": {
			Old:                "https://hello-rehvs75zla-uc.a.run.app/",
			New:                "https://hello-rehvs75zla-uc.a.run.app/",
			ExpectDiffSuppress: true,
		},
		"no slash to no slash": {
			Old:                "https://hello-rehvs75zla-uc.a.run.app",
			New:                "https://hello-rehvs75zla-uc.a.run.app",
			ExpectDiffSuppress: true,
		},
		"different domains": {
			Old:                "https://x.a.run.app/",
			New:                "https://y.a.run.app",
			ExpectDiffSuppress: false,
		},
	}

	for tn, tc := range cases {
		if LastSlashDiffSuppress("uri", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Fatalf("bad: %s, '%s' => '%s' expect %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}

func TestEmptyOrUnsetBlockDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Key, Old, New      string
		OldVal, NewVal     interface{}
		ExpectDiffSuppress bool
	}{
		"empty block vs. block containing empty string": {
			Key:                "example_block.#",
			Old:                "0",
			New:                "1",
			OldVal:             []interface{}{},
			NewVal:             []interface{}{map[string]interface{}{"empty_string": ""}},
			ExpectDiffSuppress: true,
		},
		"empty block vs. block containing false bool": {
			Key:                "example_block.#",
			Old:                "0",
			New:                "1",
			OldVal:             []interface{}{},
			NewVal:             []interface{}{map[string]interface{}{"false_bool": false}},
			ExpectDiffSuppress: true,
		},
		"empty block vs. block containing empty list": {
			Key:                "example_block.#",
			Old:                "0",
			New:                "1",
			OldVal:             []interface{}{},
			NewVal:             []interface{}{map[string]interface{}{"example_list": []interface{}{}}},
			ExpectDiffSuppress: true,
		},
		// If a parent block returns an empty sub-block in lieu of nil or an empty map, the values of the undefined
		// parent block and an empty, but defined block will be identical while the array count will have changed
		"nested block, defined empty vs. undefined": {
			Key:                "example_block.#",
			Old:                "1",
			New:                "0",
			OldVal:             []interface{}{map[string]interface{}{"nested_block": []interface{}{}}},
			NewVal:             []interface{}{map[string]interface{}{"nested_block": []interface{}{}}},
			ExpectDiffSuppress: true,
		},
		"nested block, defined empty vs. nil": {
			Key:                "node_pool_auto_config.#",
			Old:                "1",
			New:                "0",
			OldVal:             []interface{}{map[string]interface{}{"network_tags": []interface{}{}}},
			NewVal:             nil,
			ExpectDiffSuppress: true,
		},
		"nested block, empty vs. non-empty list": {
			Key:                "node_pool_auto_config.#",
			Old:                "0",
			New:                "1",
			OldVal:             []interface{}{},
			NewVal:             []interface{}{map[string]interface{}{"network_tags": []interface{}{map[string]interface{}{"tags": []interface{}{"test-network-tag"}}}}},
			ExpectDiffSuppress: false,
		},
		"nested block with nil list": {
			Key:                "node_pool_auto_config.#",
			Old:                "0",
			New:                "1",
			OldVal:             nil,
			NewVal:             []interface{}{map[string]interface{}{"network_tags": []interface{}{map[string]interface{}{"tags": nil}}}},
			ExpectDiffSuppress: false,
		},
		"nested block with empty list": {
			Key:                "node_pool_auto_config.#",
			Old:                "0",
			New:                "1",
			OldVal:             nil,
			NewVal:             []interface{}{map[string]interface{}{"network_tags": []interface{}{map[string]interface{}{"tags": []interface{}{}}}}},
			ExpectDiffSuppress: false,
		},
		"list inside nested optional block": {
			Key:                "node_pool_auto_config.0.network_tags.0.tags.#",
			Old:                "0",
			New:                "1",
			OldVal:             []interface{}{},
			NewVal:             []interface{}{"test-network-tag"},
			ExpectDiffSuppress: false,
		},
		"list item inside optional block": {
			Key:                "node_pool_auto_config.0.network_tags.0.tags.0",
			Old:                "",
			New:                "test-network-tag",
			OldVal:             "",
			NewVal:             "test-network-tag",
			ExpectDiffSuppress: false,
		},
	}

	for tn, tc := range cases {
		if EmptyOrUnsetBlockDiffSuppressLogic(tc.Key, tc.Old, tc.New, tc.OldVal, tc.NewVal) != tc.ExpectDiffSuppress {
			t.Fatalf("bad: %s, '%s' => '%s' expect %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
		if EmptyOrUnsetBlockDiffSuppressLogic(tc.Key, tc.New, tc.Old, tc.NewVal, tc.OldVal) != tc.ExpectDiffSuppress {
			t.Fatalf("bad: %s (reverse check), '%s' => '%s' expect %t", tn, tc.New, tc.Old, tc.ExpectDiffSuppress)
		}
	}
}
