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

package firestore_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"google.golang.org/api/googleapi"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccFirestoreField_firestoreFieldBasicExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_id":              envvar.GetTestProjectFromEnv(),
		"delete_protection_state": "DELETE_PROTECTION_DISABLED",
		"random_suffix":           acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckFirestoreFieldDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFirestoreField_firestoreFieldBasicExample(context),
			},
			{
				ResourceName:            "google_firestore_field.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"collection", "database", "field"},
			},
		},
	})
}

func testAccFirestoreField_firestoreFieldBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_firestore_database" "database" {
  project     = "%{project_id}"
  name        = "tf-test-database-id%{random_suffix}"
  location_id = "nam5"
  type        = "FIRESTORE_NATIVE"

  delete_protection_state = "%{delete_protection_state}"
  deletion_policy         = "DELETE"
}

resource "google_firestore_field" "basic" {
  project    = "%{project_id}"
  database   = google_firestore_database.database.name
  collection = "chatrooms_%{random_suffix}"
  field      = "basic"

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

func TestAccFirestoreField_firestoreFieldTimestampExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_id":              envvar.GetTestProjectFromEnv(),
		"delete_protection_state": "DELETE_PROTECTION_DISABLED",
		"random_suffix":           acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckFirestoreFieldDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFirestoreField_firestoreFieldTimestampExample(context),
			},
			{
				ResourceName:            "google_firestore_field.timestamp",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"collection", "database", "field"},
			},
		},
	})
}

func testAccFirestoreField_firestoreFieldTimestampExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_firestore_database" "database" {
  project     = "%{project_id}"
  name        = "tf-test-database-id%{random_suffix}"
  location_id = "nam5"
  type        = "FIRESTORE_NATIVE"

  delete_protection_state = "%{delete_protection_state}"
  deletion_policy         = "DELETE"
}

resource "google_firestore_field" "timestamp" {
  project    = "%{project_id}"
  database   = google_firestore_database.database.name
  collection = "chatrooms"
  field      = "timestamp"

  # enables a TTL policy for the document based on the value of entries with this field
  ttl_config {}

  // Disable all single field indexes for the timestamp property.
  index_config {}
}
`, context)
}

func TestAccFirestoreField_firestoreFieldMatchOverrideExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_id":              envvar.GetTestProjectFromEnv(),
		"delete_protection_state": "DELETE_PROTECTION_DISABLED",
		"random_suffix":           acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckFirestoreFieldDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFirestoreField_firestoreFieldMatchOverrideExample(context),
			},
			{
				ResourceName:            "google_firestore_field.match_override",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"collection", "database", "field"},
			},
		},
	})
}

func testAccFirestoreField_firestoreFieldMatchOverrideExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_firestore_database" "database" {
  project     = "%{project_id}"
  name        = "tf-test-database-id%{random_suffix}"
  location_id = "nam5"
  type        = "FIRESTORE_NATIVE"

  delete_protection_state = "%{delete_protection_state}"
  deletion_policy         = "DELETE"
}

resource "google_firestore_field" "match_override" {
  project    = "%{project_id}"
  database   = google_firestore_database.database.name
  collection = "chatrooms_%{random_suffix}"
  field      = "field_with_same_configuration_as_ancestor"

  index_config {
    indexes {
        order = "ASCENDING"
    }
    indexes {
        order = "DESCENDING"
    }
    indexes {
        array_config = "CONTAINS"
    }
  }
}
`, context)
}

func TestAccFirestoreField_firestoreFieldWildcardExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_id":              envvar.GetTestProjectFromEnv(),
		"delete_protection_state": "DELETE_PROTECTION_DISABLED",
		"random_suffix":           acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckFirestoreFieldDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFirestoreField_firestoreFieldWildcardExample(context),
			},
			{
				ResourceName:            "google_firestore_field.wildcard",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"collection", "database", "field"},
			},
		},
	})
}

func testAccFirestoreField_firestoreFieldWildcardExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_firestore_database" "database" {
	project     = "%{project_id}"
	name        = "tf-test-database-id%{random_suffix}"
	location_id = "nam5"
	type        = "FIRESTORE_NATIVE"

	delete_protection_state = "%{delete_protection_state}"
	deletion_policy         = "DELETE"
  }

  resource "google_firestore_field" "wildcard" {
	project    = "%{project_id}"
	database   = google_firestore_database.database.name
	collection = "chatrooms_%{random_suffix}"
	field      = "*"

	index_config {
	  indexes {
		  order       = "ASCENDING"
		  query_scope = "COLLECTION_GROUP"
	  }
	  indexes {
		  array_config = "CONTAINS"
	  }
	}
  }
`, context)
}

func testAccCheckFirestoreFieldDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_firestore_field" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			// Firestore fields are not deletable. We consider the field deleted if:
			// 1) the index configuration has no overrides and matches the ancestor configuration.
			// 2) the ttl configuration is unset.

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{FirestoreBasePath}}projects/{{project}}/databases/{{database}}/collectionGroups/{{collection}}/fields/{{field}}")
			if err != nil {
				return err
			}

			res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err != nil {
				e := err.(*googleapi.Error)
				if e.Code == 403 && strings.Contains(e.Message, "Cloud Firestore API has not been used in project") {
					// The acceptance test has provisioned the resources under test in a new project, and the destroy check is seeing the
					// effects of the project not existing. This means the service isn't enabled, and that the resource is definitely destroyed.
					// We do not return the error in this case - destroy was successful
					return nil
				}

				// Return err in all other cases
				return err
			}

			if v := res["indexConfig"]; v != nil {
				indexConfig := v.(map[string]interface{})

				usesAncestorConfig, ok := indexConfig["usesAncestorConfig"].(bool)

				if !ok || !usesAncestorConfig {
					return fmt.Errorf("Index configuration is not using the ancestor config %s.", url)
				}
			}

			if res["ttlConfig"] != nil {
				return fmt.Errorf("TTL configuration was not deleted at %s.", url)
			}
		}

		return nil
	}
}
