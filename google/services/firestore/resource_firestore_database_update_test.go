// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package firestore_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFirestoreDatabase_updateConcurrencyMode(t *testing.T) {
	t.Parallel()

	orgId := envvar.GetTestOrgFromEnv(t)
	billingAccount := envvar.GetTestBillingAccountFromEnv(t)
	randomSuffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccFirestoreDatabase_concurrencyMode(orgId, billingAccount, randomSuffix, "OPTIMISTIC"),
			},
			{
				ResourceName:            "google_firestore_database.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "project"},
			},
			{
				Config: testAccFirestoreDatabase_concurrencyMode(orgId, billingAccount, randomSuffix, "PESSIMISTIC"),
			},
			{
				ResourceName:            "google_firestore_database.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "project"},
			},
		},
	})
}

func TestAccFirestoreDatabase_updatePitrEnablement(t *testing.T) {
	t.Parallel()

	orgId := envvar.GetTestOrgFromEnv(t)
	billingAccount := envvar.GetTestBillingAccountFromEnv(t)
	randomSuffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccFirestoreDatabase_pitrEnablement(orgId, billingAccount, randomSuffix, "POINT_IN_TIME_RECOVERY_ENABLED"),
			},
			{
				ResourceName:            "google_firestore_database.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "project"},
			},
			{
				Config: testAccFirestoreDatabase_pitrEnablement(orgId, billingAccount, randomSuffix, "POINT_IN_TIME_RECOVERY_DISABLED"),
			},
			{
				ResourceName:            "google_firestore_database.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "project"},
			},
		},
	})
}

func testAccFirestoreDatabase_basicDependencies(orgId, billingAccount string, randomSuffix string) string {
	return fmt.Sprintf(`
resource "google_project" "default" {
  project_id      = "tf-test%s"
  name            = "tf-test%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "time_sleep" "wait_60_seconds" {
  depends_on = [google_project.default]

  create_duration = "60s"
}

resource "google_project_service" "firestore" {
  project = google_project.default.project_id
  service = "firestore.googleapis.com"

  # Needed for CI tests for permissions to propagate, should not be needed for actual usage
  depends_on = [time_sleep.wait_60_seconds]
}
`, randomSuffix, randomSuffix, orgId, billingAccount)
}

func testAccFirestoreDatabase_concurrencyMode(orgId, billingAccount string, randomSuffix string, concurrencyMode string) string {
	return testAccFirestoreDatabase_basicDependencies(orgId, billingAccount, randomSuffix) + fmt.Sprintf(`

resource "google_firestore_database" "default" {
  name             = "(default)"
  type             = "DATASTORE_MODE"
  location_id      = "nam5"
  concurrency_mode = "%s"

  project = google_project.default.project_id

  depends_on = [google_project_service.firestore]
}
`, concurrencyMode)
}

func testAccFirestoreDatabase_pitrEnablement(orgId, billingAccount string, randomSuffix string, pointInTimeRecoveryEnablement string) string {
	return testAccFirestoreDatabase_basicDependencies(orgId, billingAccount, randomSuffix) + fmt.Sprintf(`

resource "google_firestore_database" "default" {
  name                              = "(default)"
  type                              = "DATASTORE_MODE"
  location_id                       = "nam5"
  point_in_time_recovery_enablement = "%s"

  project = google_project.default.project_id

  depends_on = [google_project_service.firestore]
}
`, pointInTimeRecoveryEnablement)
}
