// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccProjectIamMemberRemove_basic(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	randomSuffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleProjectIamCustomRoleDestroyProducer(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config:             testAccCheckGoogleProjectIamMemberRemove_basic(randomSuffix, org),
				ExpectNonEmptyPlan: true, // Due to adding in binding, then removing in remove resource
			},
			{
				Config:   testAccCheckGoogleProjectIamMemberRemove_basic2(randomSuffix, org),
				PlanOnly: true, // binding expects the membership to be removed. Any diff will fail the test due to PlanOnly.
			},
		},
	})
}

func TestAccProjectIamMemberRemove_multipleMembersInBinding(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	randomSuffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleProjectIamCustomRoleDestroyProducer(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config:             testAccCheckGoogleProjectIamMemberRemove_multipleMemberBinding(randomSuffix, org),
				ExpectNonEmptyPlan: true, // Due to adding in binding, then removing in remove resource
			},
			{
				Config:   testAccCheckGoogleProjectIamMemberRemove_multipleMemberBinding2(randomSuffix, org),
				PlanOnly: true, // binding expects the membership to be removed. Any diff will fail the test due to PlanOnly.
			},
		},
	})
}

func TestAccProjectIamMemberRemove_memberInMultipleBindings(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	randomSuffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleProjectIamCustomRoleDestroyProducer(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config:             testAccCheckGoogleProjectIamMemberRemove_multipleMemberBinding(randomSuffix, org),
				ExpectNonEmptyPlan: true, // Due to adding in binding, then removing in remove resource
			},
			{
				Config:   testAccCheckGoogleProjectIamMemberRemove_multipleMemberBinding2(randomSuffix, org),
				PlanOnly: true, // binding expects the membership to be removed. Any diff will fail the test due to PlanOnly.
			},
		},
	})
}

func testAccCheckGoogleProjectIamMemberRemove_basic(randomSuffix, org string) string {
	return fmt.Sprintf(`
resource "google_project" "project" {
  project_id = "tf-test-%s"
  name       = "tf-test-%s"
  org_id     = "%s"
}

resource "google_project_iam_binding" "bar" {
  project = google_project.project.project_id
  members = ["user:gterraformtest1@gmail.com"]
  role    = "roles/editor"
}

resource "time_sleep" "wait_20s" {
  depends_on = [google_project_iam_binding.bar]
  create_duration = "20s"
}

resource "google_project_iam_member_remove" "foo" {
  role     = "roles/editor"
  project  = google_project.project.project_id
  member  = "user:gterraformtest1@gmail.com"
  depends_on = [time_sleep.wait_20s]
}
`, randomSuffix, randomSuffix, org)
}

func testAccCheckGoogleProjectIamMemberRemove_basic2(randomSuffix, org string) string {
	return fmt.Sprintf(`
resource "google_project" "project" {
  project_id = "tf-test-%s"
  name       = "tf-test-%s"
  org_id     = "%s"
}

resource "google_project_iam_binding" "bar" {
  project = google_project.project.project_id
  members = []
  role    = "roles/editor"
}

resource "time_sleep" "wait_20s" {
  depends_on = [google_project_iam_binding.bar]
  create_duration = "20s"
}

resource "google_project_iam_member_remove" "foo" {
  role     = "roles/editor"
  project  = google_project.project.project_id
  member  = "user:gterraformtest1@gmail.com"
  depends_on = [time_sleep.wait_20s]
}
`, randomSuffix, randomSuffix, org)
}

func testAccCheckGoogleProjectIamMemberRemove_multipleMemberBinding(random_suffix, org string) string {
	return fmt.Sprintf(`
resource "google_project" "project" {
  project_id = "tf-test-%s"
  name       = "tf-test-%s"
  org_id     = "%s"
}

resource "google_project_iam_binding" "bar" {
  project = google_project.project.project_id
  members = ["user:gterraformtest1@gmail.com", "user:gterraformtest2@gmail.com"]
  role    = "roles/editor"
}

resource "time_sleep" "wait_20s" {
  depends_on = [google_project_iam_binding.bar]
  create_duration = "20s"
}

resource "google_project_iam_member_remove" "foo" {
  role     = "roles/editor"
  project  = google_project.project.project_id
  member  = "user:gterraformtest1@gmail.com"
  depends_on = [time_sleep.wait_20s]
}
`, random_suffix, random_suffix, org)
}

func testAccCheckGoogleProjectIamMemberRemove_multipleMemberBinding2(random_suffix, org string) string {
	return fmt.Sprintf(`
resource "google_project" "project" {
  project_id = "tf-test-%s"
  name       = "tf-test-%s"
  org_id     = "%s"
}

resource "google_project_iam_binding" "bar" {
  project = google_project.project.project_id
  members = ["user:gterraformtest2@gmail.com"]
  role    = "roles/editor"
}

resource "time_sleep" "wait_20s" {
  depends_on = [google_project_iam_binding.bar]
  create_duration = "20s"
}

resource "google_project_iam_member_remove" "foo" {
  role     = "roles/editor"
  project  = google_project.project.project_id
  member  = "user:gterraformtest1@gmail.com"
  depends_on = [time_sleep.wait_20s]
}
`, random_suffix, random_suffix, org)
}

func testAccProjectIamMemberRemove_memberInMultipleBindings(random_suffix, org string) string {
	return fmt.Sprintf(`
resource "google_project" "project" {
  project_id = "tf-test-%s"
  name       = "tf-test-%s"
  org_id     = "%s"
}

resource "google_project_iam_binding" "bar" {
  project = google_project.project.project_id
  members = ["user:gterraformtest1@gmail.com"]
  role    = "roles/editor"
}

resource "google_project_iam_binding" "baz" {
  project = google_project.project.project_id
  members = ["user:gterraformtest1@gmail.com"]
  role    = "roles/viewer"
}

resource "time_sleep" "wait_20s" {
  depends_on = [google_project_iam_binding.bar, google_project_iam_binding.baz]
  create_duration = "20s"
}

resource "google_project_iam_member_remove" "foo" {
  role     = "roles/editor"
  project  = google_project.project.project_id
  member  = "user:gterraformtest1@gmail.com"
  depends_on = [time_sleep.wait_20s]
}
`, random_suffix, random_suffix, org)
}

func testAccProjectIamMemberRemove_memberInMultipleBindings2(random_suffix, org string) string {
	return fmt.Sprintf(`
resource "google_project" "project" {
  project_id = "tf-test-%s"
  name       = "tf-test-%s"
  org_id     = "%s"
}

resource "google_project_iam_binding" "bar" {
  project = google_project.project.project_id
  members = []
  role    = "roles/editor"
}

resource "google_project_iam_binding" "baz" {
  project = google_project.project.project_id
  members = ["user:gterraformtest1@gmail.com"]
  role    = "roles/viewer"
}

resource "time_sleep" "wait_20s" {
  depends_on = [google_project_iam_binding.bar, google_project_iam_binding.baz]
  create_duration = "20s"
}

resource "google_project_iam_member_remove" "foo" {
  role     = "roles/editor"
  project  = google_project.project.project_id
  member  = "user:gterraformtest1@gmail.com"
  depends_on = [time_sleep.wait_20s]
}
`, random_suffix, random_suffix, org)
}
