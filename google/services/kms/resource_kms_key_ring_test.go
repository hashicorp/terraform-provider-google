// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package kms_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccKmsKeyRing_basic(t *testing.T) {
	projectId := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	projectOrg := envvar.GetTestOrgFromEnv(t)
	projectBillingAccount := envvar.GetTestBillingAccountFromEnv(t)
	keyRingName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleKmsKeyRingWasRemovedFromState("google_kms_key_ring.key_ring"),
		Steps: []resource.TestStep{
			{
				Config: testGoogleKmsKeyRing_basic(projectId, projectOrg, projectBillingAccount, keyRingName),
			},
			{
				ResourceName:      "google_kms_key_ring.key_ring",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testGoogleKmsKeyRing_removed(projectId, projectOrg, projectBillingAccount),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleKmsKeyRingWasRemovedFromState("google_kms_key_ring.key_ring"),
				),
			},
		},
	})
}

// KMS KeyRings cannot be deleted. This ensures that the KeyRing resource was removed from state,
// even though the server-side resource was not removed.
func testAccCheckGoogleKmsKeyRingWasRemovedFromState(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]

		if ok {
			return fmt.Errorf("Resource was not removed from state: %s", resourceName)
		}

		return nil
	}
}

// This test runs in its own project, otherwise the test project would start to get filled
// with undeletable resources
func testGoogleKmsKeyRing_basic(projectId, projectOrg, projectBillingAccount, keyRingName string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  name            = "%s"
  project_id      = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "acceptance" {
  project = google_project.acceptance.project_id
  service = "cloudkms.googleapis.com"
}

resource "google_kms_key_ring" "key_ring" {
  project  = google_project_service.acceptance.project
  name     = "%s"
  location = "us-central1"
}
`, projectId, projectId, projectOrg, projectBillingAccount, keyRingName)
}

func testGoogleKmsKeyRing_removed(projectId, projectOrg, projectBillingAccount string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  name            = "%s"
  project_id      = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "acceptance" {
  project = google_project.acceptance.project_id
  service = "cloudkms.googleapis.com"
}
`, projectId, projectId, projectOrg, projectBillingAccount)
}
