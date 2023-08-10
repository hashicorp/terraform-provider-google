// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package securitycenter_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSecurityCenterSourceIamBinding(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"role":          "roles/securitycenter.sourcesViewer",
		"org_id":        envvar.GetTestOrgFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterSourceIamBinding_basic(context),
			},
			{
				ResourceName: "google_scc_source_iam_binding.foo",
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					// This has to be a function because sources only use numeric IDs
					id := state.RootModule().Resources["google_scc_source.custom_source"].Primary.Attributes["id"]
					return fmt.Sprintf("%s %s",
						id,
						context["role"],
					), nil
				},
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Binding update
				Config: testAccSecurityCenterSourceIamBinding_update(context),
			},
			{
				ResourceName: "google_scc_source_iam_binding.foo",
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					// This has to be a function because sources only use numeric IDs
					id := state.RootModule().Resources["google_scc_source.custom_source"].Primary.Attributes["id"]
					return fmt.Sprintf("%s %s",
						id,
						context["role"],
					), nil
				},
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSecurityCenterSourceIamMember(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"role":          "roles/securitycenter.sourcesViewer",
		"org_id":        envvar.GetTestOrgFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccSecurityCenterSourceIamMember_basic(context),
			},
			{
				ResourceName: "google_scc_source_iam_member.foo",
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					// This has to be a function because sources only use numeric IDs
					id := state.RootModule().Resources["google_scc_source.custom_source"].Primary.Attributes["id"]
					return fmt.Sprintf("%s %s user:admin@hashicorptest.com",
						id,
						context["role"],
					), nil
				},
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSecurityCenterSourceIamPolicy(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"role":          "roles/securitycenter.sourcesViewer",
		"org_id":        envvar.GetTestOrgFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterSourceIamPolicy_basic(context),
			},
			{
				ResourceName:      "google_scc_source_iam_policy.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSecurityCenterSourceIamPolicy_emptyBinding(context),
			},
			{
				ResourceName:      "google_scc_source_iam_policy.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSecurityCenterSourceIamMember_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_scc_source" "custom_source" {
  display_name = "tf-test-source%{random_suffix}"
  organization = "%{org_id}"
  description  = "My custom Cloud Security Command Center Finding Source"
}

resource "google_scc_source_iam_member" "foo" {
  source       = google_scc_source.custom_source.id
  organization = "%{org_id}"
  role         = "%{role}"
  member       = "user:admin@hashicorptest.com"
}
`, context)
}

func testAccSecurityCenterSourceIamPolicy_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_scc_source" "custom_source" {
  display_name = "tf-test-source%{random_suffix}"
  organization = "%{org_id}"
  description  = "My custom Cloud Security Command Center Finding Source"
}

data "google_iam_policy" "foo" {
  binding {
    role    = "%{role}"
    members = ["user:admin@hashicorptest.com"]
  }
}

resource "google_scc_source_iam_policy" "foo" {
  source       = google_scc_source.custom_source.id
  organization = "%{org_id}"
  policy_data  = data.google_iam_policy.foo.policy_data
}
`, context)
}

func testAccSecurityCenterSourceIamPolicy_emptyBinding(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_scc_source" "custom_source" {
  display_name = "tf-test-source%{random_suffix}"
  organization = "%{org_id}"
  description  = "My custom Cloud Security Command Center Finding Source"
}

data "google_iam_policy" "foo" {
}

resource "google_scc_source_iam_policy" "foo" {
  source       = google_scc_source.custom_source.id
  organization = "%{org_id}"
  policy_data  = data.google_iam_policy.foo.policy_data
}
`, context)
}

func testAccSecurityCenterSourceIamBinding_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_scc_source" "custom_source" {
  display_name = "tf-test-source%{random_suffix}"
  organization = "%{org_id}"
  description  = "My custom Cloud Security Command Center Finding Source"
}

resource "google_scc_source_iam_binding" "foo" {
  source       = google_scc_source.custom_source.id
  organization = "%{org_id}"
  role         = "%{role}"
  members      = ["user:admin@hashicorptest.com"]
}
`, context)
}

func testAccSecurityCenterSourceIamBinding_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_scc_source" "custom_source" {
  display_name = "tf-test-source%{random_suffix}"
  organization = "%{org_id}"
  description  = "My custom Cloud Security Command Center Finding Source"
}

resource "google_scc_source_iam_binding" "foo" {
  source       = google_scc_source.custom_source.id
  organization = "%{org_id}"
  role         = "%{role}"
  members      = ["user:admin@hashicorptest.com", "user:gterraformtest1@gmail.com"]
}
`, context)
}
