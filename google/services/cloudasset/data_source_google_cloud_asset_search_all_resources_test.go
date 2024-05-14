// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudasset_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceGoogleCloudAssetSearchAllResources_basic(t *testing.T) {
	t.Parallel()

	project := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleCloudAssetProjectResourcesList(project),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.google_cloud_asset_search_all_resources.resources",
						"results.0.asset_type", regexp.MustCompile("cloudresourcemanager.googleapis.com/Project")),
					resource.TestMatchResourceAttr("data.google_cloud_asset_search_all_resources.resources",
						"results.0.display_name", regexp.MustCompile(project)),
					resource.TestMatchResourceAttr("data.google_cloud_asset_search_all_resources.resources",
						"results.0.name", regexp.MustCompile(fmt.Sprintf("//cloudresourcemanager.googleapis.com/projects/%s", project))),
					resource.TestCheckResourceAttrSet("data.google_cloud_asset_search_all_resources.resources", "results.0.location"),
					resource.TestCheckResourceAttrSet("data.google_cloud_asset_search_all_resources.resources", "results.0.project"),
				),
			},
		},
	})
}

func testAccCheckGoogleCloudAssetProjectResourcesList(project string) string {
	return fmt.Sprintf(`
data google_cloud_asset_search_all_resources resources {
	scope = "projects/%s"
	asset_types = [
		"cloudresourcemanager.googleapis.com/Project"
	]
}
`, project)
}
