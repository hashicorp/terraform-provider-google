package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccStorageBucketAccessControl_update(t *testing.T) {
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
		CheckDestroy: testAccCheckStorageObjectAccessControlDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketAccessControlBasic(bucketName, "READER", "allUsers"),
			},
			{
				ResourceName:      "google_storage_bucket_access_control.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testGoogleStorageBucketAccessControlBasic(bucketName, "OWNER", "allUsers"),
			},
			{
				ResourceName:      "google_storage_bucket_access_control.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testGoogleStorageBucketAccessControlBasic(bucketName, role, entity string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket_access_control" "default" {
  bucket = google_storage_bucket.bucket.name
  role   = "%s"
  entity = "%s"
}

resource "google_storage_bucket" "bucket" {
	name = "%s"
}
`, role, entity, bucketName)
}
