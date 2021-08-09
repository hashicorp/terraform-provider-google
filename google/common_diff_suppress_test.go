// Contains common diff suppress functions.

package google

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
		if optionalPrefixSuppress(tc.Prefix)("folder", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
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
		if ignoreMissingKeyInMap(tc.Key)("push_config.0.attributes."+tc.Key, tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
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
		if optionalSurroundingSpacesSuppress("filter", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
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
		if caseDiffSuppress("key", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
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
		if portRangeDiffSuppress("ports", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
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
		if locationDiffSuppress("policy_uri", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
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
		if absoluteDomainSuppress("managed.0.domains.", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
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
		if durationDiffSuppress("duration", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Fatalf("bad: %s, '%s' => '%s' expect %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}
