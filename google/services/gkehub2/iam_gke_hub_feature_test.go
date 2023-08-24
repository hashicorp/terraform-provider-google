// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package gkehub2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccGKEHub2FeatureIamBindingGenerated(t *testing.T) {
	// VCR fails to handle batched project services
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":   acctest.RandString(t, 10),
		"role":            "roles/viewer",
		"project_id":      fmt.Sprintf("tf-test-gkehub-%s", acctest.RandString(t, 10)),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGKEHub2FeatureIamBinding_basicGenerated(context),
			},
			{
				ResourceName:      "google_gke_hub_feature_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/features/%s roles/viewer", context["project_id"], "global", fmt.Sprint("multiclusterservicediscovery")),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Binding update
				Config: testAccGKEHub2FeatureIamBinding_updateGenerated(context),
			},
			{
				ResourceName:      "google_gke_hub_feature_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/features/%s roles/viewer", context["project_id"], "global", fmt.Sprint("multiclusterservicediscovery")),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGKEHub2FeatureIamMemberGenerated(t *testing.T) {
	// VCR fails to handle batched project services
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":   acctest.RandString(t, 10),
		"role":            "roles/viewer",
		"project_id":      fmt.Sprintf("tf-test-gkehub-%s", acctest.RandString(t, 10)),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccGKEHub2FeatureIamMember_basicGenerated(context),
			},
			{
				ResourceName:      "google_gke_hub_feature_iam_member.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/features/%s roles/viewer user:admin@hashicorptest.com", context["project_id"], "global", fmt.Sprint("multiclusterservicediscovery")),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGKEHub2FeatureIamPolicyGenerated(t *testing.T) {
	// VCR fails to handle batched project services
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":   acctest.RandString(t, 10),
		"role":            "roles/viewer",
		"project_id":      fmt.Sprintf("tf-test-gkehub-%s", acctest.RandString(t, 10)),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGKEHub2FeatureIamPolicy_basicGenerated(context),
				Check:  resource.TestCheckResourceAttrSet("data.google_gke_hub_feature_iam_policy.foo", "policy_data"),
			},
			{
				ResourceName:      "google_gke_hub_feature_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/features/%s", context["project_id"], "global", fmt.Sprint("multiclusterservicediscovery")),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGKEHub2FeatureIamPolicy_emptyBinding(context),
			},
			{
				ResourceName:      "google_gke_hub_feature_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/features/%s", context["project_id"], "global", fmt.Sprint("multiclusterservicediscovery")),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGKEHub2FeatureIamMember_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  name            = "%{project_id}"
  project_id      = "%{project_id}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}
resource "google_project_service" "mcsd" {
  project = google_project.project.project_id
  service = "multiclusterservicediscovery.googleapis.com"
}
resource "google_project_service" "gkehub" {
  project = google_project.project.project_id
  service = "gkehub.googleapis.com"
  disable_on_destroy = false
}
resource "google_gke_hub_feature" "feature" {
  name = "multiclusterservicediscovery"
  location = "global"
  project = google_project.project.project_id
  labels = {
    foo = "bar"
  }
  depends_on = [google_project_service.mcsd, google_project_service.gkehub]
}
resource "google_gke_hub_feature_iam_member" "foo" {
  project = google_gke_hub_feature.feature.project
  location = google_gke_hub_feature.feature.location
  name = google_gke_hub_feature.feature.name
  role = "%{role}"
  member = "user:admin@hashicorptest.com"
}
`, context)
}

func testAccGKEHub2FeatureIamPolicy_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  name            = "%{project_id}"
  project_id      = "%{project_id}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}
resource "google_project_service" "mcsd" {
  project = google_project.project.project_id
  service = "multiclusterservicediscovery.googleapis.com"
}
resource "google_project_service" "gkehub" {
  project = google_project.project.project_id
  service = "gkehub.googleapis.com"
  disable_on_destroy = false
}
resource "google_gke_hub_feature" "feature" {
  name = "multiclusterservicediscovery"
  location = "global"
  project = google_project.project.project_id
  labels = {
    foo = "bar"
  }
  depends_on = [google_project_service.mcsd, google_project_service.gkehub]
}
data "google_iam_policy" "foo" {
  binding {
    role = "%{role}"
    members = ["user:admin@hashicorptest.com"]
  }
}
resource "google_gke_hub_feature_iam_policy" "foo" {
  project = google_gke_hub_feature.feature.project
  location = google_gke_hub_feature.feature.location
  name = google_gke_hub_feature.feature.name
  policy_data = data.google_iam_policy.foo.policy_data
}
data "google_gke_hub_feature_iam_policy" "foo" {
  project = google_gke_hub_feature.feature.project
  location = google_gke_hub_feature.feature.location
  name = google_gke_hub_feature.feature.name
  depends_on = [
    google_gke_hub_feature_iam_policy.foo
  ]
}
`, context)
}

func testAccGKEHub2FeatureIamPolicy_emptyBinding(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  name            = "%{project_id}"
  project_id      = "%{project_id}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}
resource "google_project_service" "mcsd" {
  project = google_project.project.project_id
  service = "multiclusterservicediscovery.googleapis.com"
}
resource "google_project_service" "gkehub" {
  project = google_project.project.project_id
  service = "gkehub.googleapis.com"
  disable_on_destroy = false
}
resource "google_gke_hub_feature" "feature" {
  name = "multiclusterservicediscovery"
  location = "global"
  project = google_project.project.project_id
  labels = {
    foo = "bar"
  }
  depends_on = [google_project_service.mcsd, google_project_service.gkehub]
}
data "google_iam_policy" "foo" {
}
resource "google_gke_hub_feature_iam_policy" "foo" {
  project = google_gke_hub_feature.feature.project
  location = google_gke_hub_feature.feature.location
  name = google_gke_hub_feature.feature.name
  policy_data = data.google_iam_policy.foo.policy_data
}
`, context)
}

func testAccGKEHub2FeatureIamBinding_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  name            = "%{project_id}"
  project_id      = "%{project_id}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}
resource "google_project_service" "mcsd" {
  project = google_project.project.project_id
  service = "multiclusterservicediscovery.googleapis.com"
}
resource "google_project_service" "gkehub" {
  project = google_project.project.project_id
  service = "gkehub.googleapis.com"
  disable_on_destroy = false
}
resource "google_gke_hub_feature" "feature" {
  name = "multiclusterservicediscovery"
  location = "global"
  project = google_project.project.project_id
  labels = {
    foo = "bar"
  }
  depends_on = [google_project_service.mcsd, google_project_service.gkehub]
}
resource "google_gke_hub_feature_iam_binding" "foo" {
  project = google_gke_hub_feature.feature.project
  location = google_gke_hub_feature.feature.location
  name = google_gke_hub_feature.feature.name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com"]
}
`, context)
}

func testAccGKEHub2FeatureIamBinding_updateGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  name            = "%{project_id}"
  project_id      = "%{project_id}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}
resource "google_project_service" "mcsd" {
  project = google_project.project.project_id
  service = "multiclusterservicediscovery.googleapis.com"
}
resource "google_project_service" "gkehub" {
  project = google_project.project.project_id
  service = "gkehub.googleapis.com"
  disable_on_destroy = false
}
resource "google_gke_hub_feature" "feature" {
  name = "multiclusterservicediscovery"
  location = "global"
  project = google_project.project.project_id
  labels = {
    foo = "bar"
  }
  depends_on = [google_project_service.mcsd, google_project_service.gkehub]
}
resource "google_gke_hub_feature_iam_binding" "foo" {
  project = google_gke_hub_feature.feature.project
  location = google_gke_hub_feature.feature.location
  name = google_gke_hub_feature.feature.name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com", "user:gterraformtest1@gmail.com"]
}
`, context)
}
