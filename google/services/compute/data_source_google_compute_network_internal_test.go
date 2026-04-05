package compute

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func TestDataSourceGoogleComputeNetwork_selfLinkParsing(t *testing.T) {
	tests := []struct {
		name        string
		selfLink    string
		wantProject string
		wantName    string
	}{
		{
			name:        "full self_link",
			selfLink:    "https://www.googleapis.com/compute/v1/projects/my-project/global/networks/my-network",
			wantProject: "my-project",
			wantName:    "my-network",
		},
		{
			name:        "partial path",
			selfLink:    "projects/my-project/global/networks/my-network",
			wantProject: "my-project",
			wantName:    "my-network",
		},
		{
			name:        "beta API version",
			selfLink:    "https://www.googleapis.com/compute/beta/projects/test-proj/global/networks/test-net",
			wantProject: "test-proj",
			wantName:    "test-net",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotName := tpgresource.GetResourceNameFromSelfLink(tc.selfLink)
			if gotName != tc.wantName {
				t.Errorf("GetResourceNameFromSelfLink(%q) = %q, want %q", tc.selfLink, gotName, tc.wantName)
			}

			// Verify project extraction logic matches what's in the data source
			parts := strings.Split(tc.selfLink, "/")
			var gotProject string
			for i, part := range parts {
				if part == "projects" && i+1 < len(parts) {
					gotProject = parts[i+1]
					break
				}
			}
			if gotProject != tc.wantProject {
				t.Errorf("project from %q = %q, want %q", tc.selfLink, gotProject, tc.wantProject)
			}
		})
	}
}
