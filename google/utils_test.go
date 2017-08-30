package google

import "testing"

func TestIpCidrRangeDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New          string
		ExpectDiffSupress bool
	}{
		"single ip address": {
			Old:               "10.2.3.4",
			New:               "10.2.3.5",
			ExpectDiffSupress: false,
		},
		"cidr format string": {
			Old:               "10.1.2.0/24",
			New:               "10.1.3.0/24",
			ExpectDiffSupress: false,
		},
		"netmask same mask": {
			Old:               "10.1.2.0/24",
			New:               "/24",
			ExpectDiffSupress: true,
		},
		"netmask different mask": {
			Old:               "10.1.2.0/24",
			New:               "/32",
			ExpectDiffSupress: false,
		},
		"add netmask": {
			Old:               "",
			New:               "/24",
			ExpectDiffSupress: false,
		},
		"remove netmask": {
			Old:               "/24",
			New:               "",
			ExpectDiffSupress: false,
		},
	}

	for tn, tc := range cases {
		if ipCidrRangeDiffSuppress("ip_cidr_range", tc.Old, tc.New, nil) != tc.ExpectDiffSupress {
			t.Fatalf("bad: %s, '%s' => '%s' expect %t", tn, tc.Old, tc.New, tc.ExpectDiffSupress)
		}
	}
}
