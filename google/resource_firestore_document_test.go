package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccFirestoreDocument_update(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", RandInt(t))
	project := acctest.GetTestFirestoreProjectFromEnv(t)

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFirestoreDocument_update(project, name),
			},
			{
				ResourceName:      "google_firestore_document.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFirestoreDocument_update2(project, name),
			},
			{
				ResourceName:      "google_firestore_document.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccFirestoreDocument_update(project, name string) string {
	return fmt.Sprintf(`
resource "google_firestore_document" "instance" {
	project     = "%s"
	database    = "(default)"
	collection  = "somenewcollection"
	document_id = "%s"
	fields      = "{\"something\":{\"mapValue\":{\"fields\":{\"yo\":{\"stringValue\":\"val1\"}}}}}"
}
`, project, name)
}

func testAccFirestoreDocument_update2(project, name string) string {
	return fmt.Sprintf(`
resource "google_firestore_document" "instance" {
	project     = "%s"
	database    = "(default)"
	collection  = "somenewcollection"
	document_id = "%s"
	fields      = "{\"something\":{\"mapValue\":{\"fields\":{\"yo\":{\"stringValue\":\"val2\"}}}}}"
}
`, project, name)
}
