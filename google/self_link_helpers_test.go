package google

import "testing"

func TestCompareSelfLinkOrResourceName(t *testing.T) {
	cases := map[string]struct {
		Old, New string
		Expect   bool
	}{
		"name only, same": {
			Old:    "https://www.googleapis.com/compute/v1/projects/your-project/global/networks/a-network",
			New:    "a-network",
			Expect: true,
		},
		"name only, different": {
			Old:    "https://www.googleapis.com/compute/v1/projects/your-project/global/networks/a-network",
			New:    "another-network",
			Expect: false,
		},
		"partial path, same": {
			Old:    "https://www.googleapis.com/compute/v1/projects/your-project/global/networks/a-network",
			New:    "projects/your-project/global/networks/a-network",
			Expect: true,
		},
		"partial path, different name": {
			Old:    "https://www.googleapis.com/compute/v1/projects/your-project/global/networks/a-network",
			New:    "projects/your-project/global/networks/another-network",
			Expect: false,
		},
		"partial path, different project": {
			Old:    "https://www.googleapis.com/compute/v1/projects/your-project/global/networks/a-network",
			New:    "projects/another-project/global/networks/a-network",
			Expect: false,
		},
		"full path, different name": {
			Old:    "https://www.googleapis.com/compute/v1/projects/your-project/global/networks/a-network",
			New:    "https://www.googleapis.com/compute/v1/projects/your-project/global/networks/another-network",
			Expect: false,
		},
		"full path, different project": {
			Old:    "https://www.googleapis.com/compute/v1/projects/your-project/global/networks/a-network",
			New:    "https://www.googleapis.com/compute/v1/projects/another-project/global/networks/a-network",
			Expect: false,
		},
		"beta full path, same": {
			Old:    "https://www.googleapis.com/compute/v1/projects/your-project/global/networks/a-network",
			New:    "https://www.googleapis.com/compute/beta/projects/your-project/global/networks/a-network",
			Expect: true,
		},
		"beta full path, different name": {
			Old:    "https://www.googleapis.com/compute/v1/projects/your-project/global/networks/a-network",
			New:    "https://www.googleapis.com/compute/beta/projects/your-project/global/networks/another-network",
			Expect: false,
		},
		"beta full path, different project": {
			Old:    "https://www.googleapis.com/compute/v1/projects/your-project/global/networks/a-network",
			New:    "https://www.googleapis.com/compute/beta/projects/another-project/global/networks/a-network",
			Expect: false,
		},
	}

	for tn, tc := range cases {
		if compareSelfLinkOrResourceName("", tc.Old, tc.New, nil) != tc.Expect {
			t.Errorf("bad: %s, expected %t for old = %q and new = %q", tn, tc.Expect, tc.Old, tc.New)
		}
	}
}

func TestGetResourceNameFromSelfLink(t *testing.T) {
	cases := map[string]struct {
		SelfLink, ExpectedName string
	}{
		"name is extracted from self_link": {
			SelfLink:     "http://something.com/one/two/three",
			ExpectedName: "three",
		},
		"name is returned if the self_link only contains the name": {
			SelfLink:     "resource_name",
			ExpectedName: "resource_name",
		},
	}

	for tn, tc := range cases {
		if n := GetResourceNameFromSelfLink(tc.SelfLink); n != tc.ExpectedName {
			t.Errorf("%s: expected resource name %q; got %q", tn, tc.ExpectedName, n)
		}
	}
}

func TestSelfLinkNameHash(t *testing.T) {
	cases := map[string]struct {
		SelfLink, Name string
		Expect         bool
	}{
		"same": {
			SelfLink: "https://www.googleapis.com/compute/v1/projects/your-project/global/networks/a-network",
			Name:     "a-network",
			Expect:   true,
		},
		"different": {
			SelfLink: "https://www.googleapis.com/compute/v1/projects/your-project/global/networks/a-network",
			Name:     "another-network",
			Expect:   false,
		},
	}

	for tn, tc := range cases {
		if (selfLinkNameHash(tc.SelfLink) == selfLinkNameHash(tc.Name)) != tc.Expect {
			t.Errorf("%s: expected %t for whether hashes matched for self link = %q, name = %q", tn, tc.Expect, tc.SelfLink, tc.Name)
		}
	}
}

func TestGetRegionFromRegionSelfLink(t *testing.T) {
	cases := map[string]string{
		"https://www.googleapis.com/compute/v1/projects/test/regions/europe-west3": "europe-west3",
		"europe-west3": "europe-west3",
	}
	for input, expected := range cases {
		if result := GetRegionFromRegionSelfLink(input); result != expected {
			t.Errorf("expected to get %q from %q, got %q", expected, input, result)
		}
	}
}
