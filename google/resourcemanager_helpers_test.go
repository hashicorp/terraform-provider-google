package google

import (
	"testing"

	"google.golang.org/api/cloudresourcemanager/v1"
)

func TestGetResourceName(t *testing.T) {
	cases := map[string]struct {
		ResourceId           *cloudresourcemanager.ResourceId
		ExpectedResourceName string
	}{
		"nil resource ID": {
			ResourceId:           nil,
			ExpectedResourceName: "",
		},
		"valid resource ID": {
			ResourceId: &cloudresourcemanager.ResourceId{
				Type: "project",
				Id:   "abcd1234",
			},
			ExpectedResourceName: "project/abcd1234",
		},
	}

	for tn, tc := range cases {
		if rn := getResourceName(tc.ResourceId); rn != tc.ExpectedResourceName {
			t.Fatalf("bad: %s, expected resource name to be '%s' but got '%s'", tn, tc.ExpectedResourceName, rn)
		}
	}
}
