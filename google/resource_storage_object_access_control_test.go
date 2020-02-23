package google

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccStorageObjectAccessControl_update(t *testing.T) {
	t.Parallel()

	bucketName := testBucketName()
	objectName := testAclObjectName()
	objectData := []byte("data data data")
	if err := ioutil.WriteFile(tfObjectAcl.Name(), objectData, 0644); err != nil {
		t.Errorf("error writing file: %v", err)
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if errObjectAcl != nil {
				panic(errObjectAcl)
			}
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStorageObjectAccessControlDestroy,
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageObjectAccessControlBasic(bucketName, objectName, "READER", "allUsers"),
			},
			{
				ResourceName:      "google_storage_object_access_control.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testGoogleStorageObjectAccessControlBasic(bucketName, objectName, "OWNER", "allUsers"),
			},
			{
				ResourceName:      "google_storage_object_access_control.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testGoogleStorageObjectAccessControlBasic(bucketName, objectName, role, entity string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
}

resource "google_storage_bucket_object" "object" {
  name   = "%s"
  bucket = google_storage_bucket.bucket.name
  source = "%s"
}

resource "google_storage_object_access_control" "default" {
  object = google_storage_bucket_object.object.name
  bucket = google_storage_bucket.bucket.name
  role   = "%s"
  entity = "%s"
}
`, bucketName, objectName, tfObjectAcl.Name(), role, entity)
}
