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

func TestAccFirestoreField_firestoreFieldUpdateAddIndexExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
		"resource_name": "add_index",
	}
	testAccFirestoreField_runUpdateTest(testAccFirestoreField_firestoreFieldUpdateAddIndexExample(context), true, t, context)
}

func TestAccFirestoreField_firestoreFieldUpdateAddTTLExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_id":    envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
		"resource_name": "add_ttl",
	}
	testAccFirestoreField_runUpdateTest(testAccFirestoreField_firestoreFieldUpdateAddTTLExample(context), false, t, context)
}

func testAccFirestoreField_runUpdateTest(updateConfig string, useOwnProject bool, t *testing.T, context map[string]interface{}) {
	resourceName := context["resource_name"].(string)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckFirestoreFieldDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFirestoreField_firestoreFieldUpdateInitialExample(context, useOwnProject),
			},
			{
				ResourceName:      fmt.Sprintf("google_firestore_field.%s", resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: updateConfig,
			},
			{
				ResourceName:      fmt.Sprintf("google_firestore_field.%s", resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFirestoreField_firestoreFieldUpdateInitialExample(context, useOwnProject),
			},
			{
				ResourceName:      fmt.Sprintf("google_firestore_field.%s", resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccFirestoreField_update_basicDeps(context map[string]interface{}, useOwnProject bool) string {
	// TTls require billing, so don't use their own project
	if useOwnProject {
		return acctest.Nprintf(`
resource "google_project" "project" {
	project_id = "tf-test%{random_suffix}"
	name       = "tf-test%{random_suffix}"
	org_id     = "%{org_id}"
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

	# used to control delete order
	depends_on = [
		google_project_service.firestore,
		google_project.project
	]
}
`, context)
	} else {
		return acctest.Nprintf(`
resource "google_firestore_database" "database" {
	project     = "%{project_id}"
	name        = "tf-test%{random_suffix}"
	location_id = "nam5"
	type        = "FIRESTORE_NATIVE"

	delete_protection_state = "DELETE_PROTECTION_DISABLED"
	deletion_policy         = "DELETE"
}
`, context)
	}
}

func testAccFirestoreField_firestoreFieldUpdateInitialExample(context map[string]interface{}, useOwnProject bool) string {
	return testAccFirestoreField_update_basicDeps(context, useOwnProject) + acctest.Nprintf(`
resource "google_firestore_field" "%{resource_name}" {
	project = google_firestore_database.database.project
	database = google_firestore_database.database.name
	collection = "chatrooms_%{random_suffix}"
	field = "%{resource_name}"

	index_config {
		indexes {
			order = "ASCENDING"
			query_scope = "COLLECTION_GROUP"
		}
		indexes {
			array_config = "CONTAINS"
		}
	}
}
`, context)
}

func testAccFirestoreField_firestoreFieldUpdateAddTTLExample(context map[string]interface{}) string {
	// TTLs need billing, so do not use isolated project
	return testAccFirestoreField_update_basicDeps(context, false) + acctest.Nprintf(`
resource "google_firestore_field" "%{resource_name}" {
	project    = google_firestore_database.database.project
	database   = google_firestore_database.database.name
	collection = "chatrooms_%{random_suffix}"
	field      = "%{resource_name}"

	index_config {
		indexes {
			order = "ASCENDING"
			query_scope = "COLLECTION_GROUP"
		}
		indexes {
			array_config = "CONTAINS"
		}
	}

	ttl_config {}
}
`, context)
}

func testAccFirestoreField_firestoreFieldUpdateAddIndexExample(context map[string]interface{}) string {
	return testAccFirestoreField_update_basicDeps(context, true) + acctest.Nprintf(`
resource "google_firestore_field" "%{resource_name}" {
	project = google_firestore_database.database.project
	database = google_firestore_database.database.name
	collection = "chatrooms_%{random_suffix}"
	field = "%{resource_name}"

	index_config {
		indexes {
			order = "ASCENDING"
			query_scope = "COLLECTION_GROUP"
		}
		indexes {
			array_config = "CONTAINS"
		}
		indexes {
			order = "DESCENDING"
			query_scope = "COLLECTION_GROUP"
		}
	}
}
`, context)
}
