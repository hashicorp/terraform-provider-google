// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package datacatalog_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataCatalogEntryGroupIamBindingGenerated(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"role":          "roles/viewer",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataCatalogEntryGroupIamBinding_basicGenerated(context),
			},
			{
				ResourceName:      "google_data_catalog_entry_group_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/entryGroups/%s roles/viewer", envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv(), fmt.Sprintf("tf_test_my_group%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Binding update
				Config: testAccDataCatalogEntryGroupIamBinding_updateGenerated(context),
			},
			{
				ResourceName:      "google_data_catalog_entry_group_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/entryGroups/%s roles/viewer", envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv(), fmt.Sprintf("tf_test_my_group%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDataCatalogEntryGroupIamMemberGenerated(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"role":          "roles/viewer",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccDataCatalogEntryGroupIamMember_basicGenerated(context),
			},
			{
				ResourceName:      "google_data_catalog_entry_group_iam_member.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/entryGroups/%s roles/viewer user:admin@hashicorptest.com", envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv(), fmt.Sprintf("tf_test_my_group%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDataCatalogEntryGroupIamPolicyGenerated(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"role":          "roles/viewer",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataCatalogEntryGroupIamPolicy_basicGenerated(context),
				Check:  resource.TestCheckResourceAttrSet("data.google_data_catalog_entry_group_iam_policy.foo", "policy_data"),
			},
			{
				ResourceName:      "google_data_catalog_entry_group_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/entryGroups/%s", envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv(), fmt.Sprintf("tf_test_my_group%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataCatalogEntryGroupIamPolicy_emptyBinding(context),
			},
			{
				ResourceName:      "google_data_catalog_entry_group_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/entryGroups/%s", envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv(), fmt.Sprintf("tf_test_my_group%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataCatalogEntryGroupIamMember_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_catalog_entry_group" "basic_entry_group" {
  entry_group_id = "tf_test_my_group%{random_suffix}"
}

resource "google_data_catalog_entry_group_iam_member" "foo" {
  entry_group = google_data_catalog_entry_group.basic_entry_group.name
  role = "%{role}"
  member = "user:admin@hashicorptest.com"
}
`, context)
}

func testAccDataCatalogEntryGroupIamPolicy_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_catalog_entry_group" "basic_entry_group" {
  entry_group_id = "tf_test_my_group%{random_suffix}"
}

data "google_iam_policy" "foo" {
  binding {
    role = "%{role}"
    members = ["user:admin@hashicorptest.com"]
  }
}

resource "google_data_catalog_entry_group_iam_policy" "foo" {
  entry_group = google_data_catalog_entry_group.basic_entry_group.name
  policy_data = data.google_iam_policy.foo.policy_data
}

data "google_data_catalog_entry_group_iam_policy" "foo" {
  entry_group = google_data_catalog_entry_group.basic_entry_group.name
  depends_on = [
    google_data_catalog_entry_group_iam_policy.foo
  ]
}
`, context)
}

func testAccDataCatalogEntryGroupIamPolicy_emptyBinding(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_catalog_entry_group" "basic_entry_group" {
  entry_group_id = "tf_test_my_group%{random_suffix}"
}

data "google_iam_policy" "foo" {
}

resource "google_data_catalog_entry_group_iam_policy" "foo" {
  entry_group = google_data_catalog_entry_group.basic_entry_group.name
  policy_data = data.google_iam_policy.foo.policy_data
}
`, context)
}

func testAccDataCatalogEntryGroupIamBinding_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_catalog_entry_group" "basic_entry_group" {
  entry_group_id = "tf_test_my_group%{random_suffix}"
}

resource "google_data_catalog_entry_group_iam_binding" "foo" {
  entry_group = google_data_catalog_entry_group.basic_entry_group.name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com"]
}
`, context)
}

func testAccDataCatalogEntryGroupIamBinding_updateGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_catalog_entry_group" "basic_entry_group" {
  entry_group_id = "tf_test_my_group%{random_suffix}"
}

resource "google_data_catalog_entry_group_iam_binding" "foo" {
  entry_group = google_data_catalog_entry_group.basic_entry_group.name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com", "user:gterraformtest1@gmail.com"]
}
`, context)
}
