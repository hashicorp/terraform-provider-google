package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccStorageDefaultObjectAccessControl_update(t *testing.T) {
	t.Parallel()

	bucketName := testBucketName(t)
	vcrTest(t, resource.TestCase{
		PreCheck: func() {
			if errObjectAcl != nil {
				panic(errObjectAcl)
			}
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStorageDefaultObjectAccessControlDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageDefaultObjectAccessControlBasic(bucketName, "READER", "allUsers"),
			},
			{
				ResourceName:      "google_storage_default_object_access_control.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testGoogleStorageDefaultObjectAccessControlBasic(bucketName, "OWNER", "allUsers"),
			},
			{
				ResourceName:      "google_storage_default_object_access_control.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testGoogleStorageDefaultObjectAccessControlBasic(bucketName, role, entity string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
}

resource "google_storage_default_object_access_control" "default" {
  bucket = google_storage_bucket.bucket.name
  role   = "%s"
  entity = "%s"
}
`, bucketName, role, entity)
}
