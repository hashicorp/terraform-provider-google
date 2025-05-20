// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storagecontrol_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccStorageControlOrganizationIntelligenceConfig_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgTargetFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageControlOrganizationIntelligenceConfig_basic(context),
			},
			{
				ResourceName:            "google_storage_control_organization_intelligence_config.organization_intelligence_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
			{
				Config: testAccStorageControlOrganizationIntelligenceConfig_update_with_filter(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_control_organization_intelligence_config.organization_intelligence_config", "filter.0.excluded_cloud_storage_buckets.0.bucket_id_regexes.0", "random-test-*"),
					resource.TestCheckResourceAttr(
						"google_storage_control_organization_intelligence_config.organization_intelligence_config", "filter.0.excluded_cloud_storage_buckets.0.bucket_id_regexes.1", "random-test2-*"),
					resource.TestCheckResourceAttr(
						"google_storage_control_organization_intelligence_config.organization_intelligence_config", "filter.0.included_cloud_storage_locations.0.locations.0", "us-east-1"),
					resource.TestCheckResourceAttr(
						"google_storage_control_organization_intelligence_config.organization_intelligence_config", "filter.0.included_cloud_storage_locations.0.locations.1", "us-east-2"),
				),
			},
			{
				ResourceName:            "google_storage_control_organization_intelligence_config.organization_intelligence_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
			{
				Config: testAccStorageControlOrganizationIntelligenceConfig_update_with_empty_filter_fields(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_control_organization_intelligence_config.organization_intelligence_config", "filter.0.excluded_cloud_storage_buckets.#", "0"),
					resource.TestCheckResourceAttr(
						"google_storage_control_organization_intelligence_config.organization_intelligence_config", "filter.0.included_cloud_storage_locations.#", "0"),
				),
			},
			{
				ResourceName:            "google_storage_control_organization_intelligence_config.organization_intelligence_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
			{
				Config: testAccStorageControlOrganizationIntelligenceConfig_update_with_filter2(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_control_organization_intelligence_config.organization_intelligence_config", "filter.0.included_cloud_storage_buckets.0.bucket_id_regexes.0", "random-test-*"),
					resource.TestCheckResourceAttr(
						"google_storage_control_organization_intelligence_config.organization_intelligence_config", "filter.0.included_cloud_storage_buckets.0.bucket_id_regexes.1", "random-test2-*"),
					resource.TestCheckResourceAttr(
						"google_storage_control_organization_intelligence_config.organization_intelligence_config", "filter.0.excluded_cloud_storage_locations.0.locations.0", "us-east-1"),
					resource.TestCheckResourceAttr(
						"google_storage_control_organization_intelligence_config.organization_intelligence_config", "filter.0.excluded_cloud_storage_locations.0.locations.1", "us-east-2"),
				),
			},
			{
				ResourceName:            "google_storage_control_organization_intelligence_config.organization_intelligence_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
			{
				Config: testAccStorageControlOrganizationIntelligenceConfig_update_with_empty_filter_fields2(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_control_organization_intelligence_config.organization_intelligence_config", "filter.0.excluded_cloud_storage_buckets.#", "0"),
					resource.TestCheckResourceAttr(
						"google_storage_control_organization_intelligence_config.organization_intelligence_config", "filter.0.included_cloud_storage_locations.#", "0"),
				),
			},
			{
				ResourceName:            "google_storage_control_organization_intelligence_config.organization_intelligence_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
			{
				Config: testAccStorageControlOrganizationIntelligenceConfig_update_mode_disable(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_control_organization_intelligence_config.organization_intelligence_config", "edition_config", "DISABLED"),
				),
			},
			{
				ResourceName:            "google_storage_control_organization_intelligence_config.organization_intelligence_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
			{
				Config: testAccStorageControlOrganizationIntelligenceConfig_update_mode_inherit(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_control_organization_intelligence_config.organization_intelligence_config", "edition_config", "INHERIT"),
				),
			},
			{
				ResourceName:            "google_storage_control_organization_intelligence_config.organization_intelligence_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
		},
	})
}

func testAccStorageControlOrganizationIntelligenceConfig_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_control_organization_intelligence_config" "organization_intelligence_config" {
  name = "%{org_id}"
  edition_config = "STANDARD"
}
`, context)
}

func testAccStorageControlOrganizationIntelligenceConfig_update_with_filter(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_control_organization_intelligence_config" "organization_intelligence_config" {
  name = "%{org_id}"
  edition_config = "STANDARD"
  filter {
    excluded_cloud_storage_buckets{
      bucket_id_regexes = ["random-test-*", "random-test2-*"]
    }
    included_cloud_storage_locations{
      locations = ["us-east-1", "us-east-2"]
    }
  }
}
`, context)
}

func testAccStorageControlOrganizationIntelligenceConfig_update_with_empty_filter_fields(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_control_organization_intelligence_config" "organization_intelligence_config" {
  name = "%{org_id}"
  edition_config = "STANDARD"
  filter {
    excluded_cloud_storage_buckets{
      bucket_id_regexes = []
    }
    included_cloud_storage_locations{
      locations = []
    }
  }
}
`, context)
}

func testAccStorageControlOrganizationIntelligenceConfig_update_with_filter2(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_control_organization_intelligence_config" "organization_intelligence_config" {
  name = "%{org_id}"
  edition_config = "STANDARD"
  filter {
    included_cloud_storage_buckets{
      bucket_id_regexes = ["random-test-*", "random-test2-*"]
    }
    excluded_cloud_storage_locations{
      locations = ["us-east-1", "us-east-2"]
    }
  }
}
`, context)
}

func testAccStorageControlOrganizationIntelligenceConfig_update_with_empty_filter_fields2(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_control_organization_intelligence_config" "organization_intelligence_config" {
  name = "%{org_id}"
  edition_config = "STANDARD"
  filter {
    included_cloud_storage_buckets{
      bucket_id_regexes = []
    }
    excluded_cloud_storage_locations{
      locations = []
    }
  }
}
`, context)
}

func testAccStorageControlOrganizationIntelligenceConfig_update_mode_disable(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_control_organization_intelligence_config" "organization_intelligence_config" {
  name = "%{org_id}"
  edition_config = "DISABLED"
}
`, context)
}

func testAccStorageControlOrganizationIntelligenceConfig_update_mode_inherit(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_control_organization_intelligence_config" "organization_intelligence_config" {
  name = "%{org_id}"
  edition_config = "INHERIT"
}
`, context)
}
