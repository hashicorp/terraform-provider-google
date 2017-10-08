package google

import "testing"

func TestParseNetworkFieldValue(t *testing.T) {
	cases := map[string]struct {
		Network              string
		ExpectedRelativeLink string
		Config               *Config
	}{
		"network is a full self link": {
			Network:              "https://www.googleapis.com/compute/v1/projects/myproject/global/networks/my-network",
			ExpectedRelativeLink: "projects/myproject/global/networks/my-network",
		},
		"network is a relative self link": {
			Network:              "projects/myproject/global/networks/my-network",
			ExpectedRelativeLink: "projects/myproject/global/networks/my-network",
		},
		"network is a partial relative self link": {
			Network:              "global/networks/my-network",
			ExpectedRelativeLink: "projects/default-project/global/networks/my-network",
			Config:               &Config{Project: "default-project"},
		},
		"network is the name only": {
			Network:              "my-network",
			ExpectedRelativeLink: "projects/default-project/global/networks/my-network",
			Config:               &Config{Project: "default-project"},
		},
	}

	for tn, tc := range cases {
		if fieldValue := ParseNetworkFieldValue(tc.Network, tc.Config); fieldValue.RelativeLink() != tc.ExpectedRelativeLink {
			t.Fatalf("bad: %s, expected relative link to be '%s' but got '%s'", tn, tc.ExpectedRelativeLink, fieldValue.RelativeLink())
		}
	}
}
