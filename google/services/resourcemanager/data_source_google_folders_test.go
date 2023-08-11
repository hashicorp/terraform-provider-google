// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceGoogleFolders_basic(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	parent := fmt.Sprintf("organizations/%s", org)
	displayName := "tf-test-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleFoldersConfig(parent, displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_folders.root-test", "folders.0.name"),
					resource.TestCheckResourceAttrSet("data.google_folders.root-test", "folders.0.display_name"),
					resource.TestCheckResourceAttrSet("data.google_folders.root-test", "folders.0.state"),
					resource.TestCheckResourceAttrSet("data.google_folders.root-test", "folders.0.create_time"),
					resource.TestCheckResourceAttrSet("data.google_folders.root-test", "folders.0.update_time"),
					resource.TestCheckResourceAttrSet("data.google_folders.root-test", "folders.0.etag"),
				),
			},
		},
	})
}

func testAccCheckGoogleFoldersConfig(parent string, displayName string) string {
	return fmt.Sprintf(`
resource "google_folder" "foobar" {
		parent       = "%s"
		display_name = "%s"
}

data "google_folders" "root-test" {
  parent_id = "%s"
}
`, parent, displayName, parent)
}
