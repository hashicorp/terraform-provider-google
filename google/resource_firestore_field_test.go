package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFirestoreField_firestoreFieldUpdateAddIndexExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_id":    GetTestFirestoreProjectFromEnv(t),
		"random_suffix": RandString(t, 10),
		"resource_name": "add_index",
	}
	testAccFirestoreField_runUpdateTest(testAccFirestoreField_firestoreFieldUpdateAddIndexExample(context), t, context)
}

func TestAccFirestoreField_firestoreFieldUpdateAddTTLExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_id":    GetTestFirestoreProjectFromEnv(t),
		"random_suffix": RandString(t, 10),
		"resource_name": "add_ttl",
	}
	testAccFirestoreField_runUpdateTest(testAccFirestoreField_firestoreFieldUpdateAddTTLExample(context), t, context)
}

func testAccFirestoreField_runUpdateTest(updateConfig string, t *testing.T, context map[string]interface{}) {
	resourceName := context["resource_name"].(string)

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckFirestoreFieldDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFirestoreField_firestoreFieldUpdateInitialExample(context),
			},
			{
				ResourceName:            fmt.Sprintf("google_firestore_field.%s", resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"database", "collection", "field"},
			},
			{
				Config: updateConfig,
			},
			{
				ResourceName:            fmt.Sprintf("google_firestore_field.%s", resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"database", "collection", "field"},
			},
			{
				Config: testAccFirestoreField_firestoreFieldUpdateInitialExample(context),
			},
			{
				ResourceName:            fmt.Sprintf("google_firestore_field.%s", resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"database", "collection", "field"},
			},
		},
	})
}

func testAccFirestoreField_firestoreFieldUpdateInitialExample(context map[string]interface{}) string {
	return Nprintf(`
resource "google_firestore_field" "%{resource_name}" {
	project = "%{project_id}"
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
	return Nprintf(`
resource "google_firestore_field" "%{resource_name}" {
	project = "%{project_id}"
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

	ttl_config {}
}
`, context)
}

func testAccFirestoreField_firestoreFieldUpdateAddIndexExample(context map[string]interface{}) string {
	return Nprintf(`
resource "google_firestore_field" "%{resource_name}" {
	project = "%{project_id}"
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
