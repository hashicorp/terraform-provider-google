// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFolderServiceIdentity_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 5),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleFolderServiceIdentity_basic(context),
				Check: resource.ComposeTestCheckFunc(
					// Email field for osconfig service account should be non-empty and contain at least an "@".
					resource.TestCheckResourceAttrWith("google_folder_service_identity.osconfig_sa", "email", func(value string) error {
						if strings.Contains(value, "@") {
							return nil
						}
						return fmt.Errorf("osconfig_sa service identity email value was %s, expected a valid email", value)
					}),
					// Member field for osconfig service account should start with "serviceAccount:service-folder-" and end with "@gcp-sa-osconfig.iam.gserviceaccount.com"
					resource.TestCheckResourceAttrWith("google_folder_service_identity.osconfig_sa", "member", func(value string) error {
						if !strings.HasPrefix(value, "serviceAccount:service-folder-") {
							return fmt.Errorf("osconfig_sa folder service identity member value %q does not start with 'serviceAccount:service-folder-'", value)
						}
						if !strings.HasSuffix(value, "@gcp-sa-osconfig.iam.gserviceaccount.com") {
							return fmt.Errorf("osconfig_sa folder service identity member value %q does not end with '@gcp-sa-osconfig.iam.gserviceaccount.com'", value)
						}
						return nil
					}),
				),
			},
		},
	})
}

func testGoogleFolderServiceIdentity_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_folder" "folder" {
  parent = "organizations/%{org_id}"
  display_name = "test-folder-%{random_suffix}"
  deletion_protection = false
}

resource "google_folder_service_identity" "osconfig_sa" {
  folder = google_folder.folder.folder_id
  service = "osconfig.googleapis.com"
}
`, context)
}
