// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package firestore_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFirestoreDatabase_updateConcurrencyMode(t *testing.T) {
	t.Parallel()

	projectId := envvar.GetTestProjectFromEnv()
	randomSuffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccFirestoreDatabase_concurrencyMode(projectId, randomSuffix, "OPTIMISTIC"),
			},
			{
				ResourceName:            "google_firestore_database.database",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "project"},
			},
			{
				Config: testAccFirestoreDatabase_concurrencyMode(projectId, randomSuffix, "PESSIMISTIC"),
			},
			{
				ResourceName:            "google_firestore_database.database",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "project"},
			},
		},
	})
}

func TestAccFirestoreDatabase_updatePitrEnablement(t *testing.T) {
	t.Parallel()

	projectId := envvar.GetTestProjectFromEnv()
	randomSuffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccFirestoreDatabase_pitrEnablement(projectId, randomSuffix, "POINT_IN_TIME_RECOVERY_ENABLED"),
			},
			{
				ResourceName:            "google_firestore_database.database",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "project"},
			},
			{
				Config: testAccFirestoreDatabase_pitrEnablement(projectId, randomSuffix, "POINT_IN_TIME_RECOVERY_DISABLED"),
			},
			{
				ResourceName:            "google_firestore_database.database",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "project"},
			},
		},
	})
}

func TestAccFirestoreDatabase_updateDeleteProtectionState(t *testing.T) {
	t.Parallel()

	projectId := envvar.GetTestProjectFromEnv()
	randomSuffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccFirestoreDatabase_deleteProtectionState(projectId, randomSuffix, "DELETE_PROTECTION_ENABLED"),
			},
			{
				ResourceName:            "google_firestore_database.database",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "project"},
			},
			{
				Config: testAccFirestoreDatabase_deleteProtectionState(projectId, randomSuffix, "DELETE_PROTECTION_DISABLED"),
			},
			{
				ResourceName:            "google_firestore_database.database",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "project"},
			},
		},
	})
}

func testAccFirestoreDatabase_concurrencyMode(projectId string, randomSuffix string, concurrencyMode string) string {
	return fmt.Sprintf(`
resource "google_firestore_database" "database" {
  project          = "%s"
  name             = "tf-test-%s"
  type             = "DATASTORE_MODE"
  location_id      = "nam5"
  concurrency_mode = "%s"
}
`, projectId, randomSuffix, concurrencyMode)
}

func testAccFirestoreDatabase_pitrEnablement(projectId string, randomSuffix string, pointInTimeRecoveryEnablement string) string {
	return fmt.Sprintf(`
resource "google_firestore_database" "database" {
  project                           = "%s"
  name                              = "tf-test-%s"
  type                              = "DATASTORE_MODE"
  location_id                       = "nam5"
  point_in_time_recovery_enablement = "%s"
}
`, projectId, randomSuffix, pointInTimeRecoveryEnablement)
}

func testAccFirestoreDatabase_deleteProtectionState(projectId string, randomSuffix string, deleteProtectionState string) string {
	return fmt.Sprintf(`
resource "google_firestore_database" "database" {
  project                 = "%s"
  name                    = "tf-test-%s"
  type                    = "DATASTORE_MODE"
  location_id             = "nam5"
  delete_protection_state = "%s"
}
`, projectId, randomSuffix, deleteProtectionState)
}
