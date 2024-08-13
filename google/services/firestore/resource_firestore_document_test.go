// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package firestore_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccFirestoreDocument_update(t *testing.T) {
	t.Parallel()

	orgId := envvar.GetTestOrgFromEnv(t)
	randomSuffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccFirestoreDocument_update(randomSuffix, orgId, "OPTIMISTIC", "val1"),
			},
			{
				ResourceName:      "google_firestore_document.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFirestoreDocument_update(randomSuffix, orgId, "OPTIMISTIC", "val2"),
			},
			{
				ResourceName:      "google_firestore_document.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccFirestoreDocument_update_basicDeps(randomSuffix, orgId string) string {
	return fmt.Sprintf(`
resource "google_project" "project" {
	project_id = "tf-test%s"
	name       = "tf-test%s"
	org_id     = "%s"
	deletion_policy = "DELETE"
}

resource "time_sleep" "wait_60_seconds" {
	depends_on = [google_project.project]

	create_duration = "60s"
}

resource "google_project_service" "firestore" {
	project = google_project.project.project_id
	service = "firestore.googleapis.com"

	# Needed for CI tests for permissions to propagate, should not be needed for actual usage
	depends_on = [time_sleep.wait_60_seconds]
}

resource "google_firestore_database" "database" {
	project     = google_project.project.project_id
	name        = "(default)"
	location_id = "nam5"
	type        = "FIRESTORE_NATIVE"

	depends_on = [google_project_service.firestore]
}
`, randomSuffix, randomSuffix, orgId)
}

func testAccFirestoreDocument_update(randomSuffix, orgId, name, val string) string {
	return testAccFirestoreDocument_update_basicDeps(randomSuffix, orgId) + fmt.Sprintf(`
resource "google_firestore_document" "instance" {
	project     = google_project.project.project_id
	database    = google_firestore_database.database.name
	collection  = "somenewcollection"
	document_id = "%s"
	fields      = "{\"something\":{\"mapValue\":{\"fields\":{\"yo\":{\"stringValue\":\"%s\"}}}}}"
}
`, name, val)
}
