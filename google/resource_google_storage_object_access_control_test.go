package google

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccStorageObjectAccessControl_basic(t *testing.T) {
	t.Parallel()

	bucketName := testBucketName()
	objectName := testAclObjectName()
	objectData := []byte("data data data")
	ioutil.WriteFile(tfObjectAcl.Name(), objectData, 0644)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if errObjectAcl != nil {
				panic(errObjectAcl)
			}
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageObjectAccessControlDestroy,
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageObjectAccessControlBasic(bucketName, objectName, "READER", "allUsers"),
			},
			{
				ResourceName:      "google_storage_object_access_control.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageObjectAccessControl_update(t *testing.T) {
	t.Parallel()

	bucketName := testBucketName()
	objectName := testAclObjectName()
	objectData := []byte("data data data")
	ioutil.WriteFile(tfObjectAcl.Name(), objectData, 0644)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if errObjectAcl != nil {
				panic(errObjectAcl)
			}
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageObjectAccessControlDestroy,
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

func testAccStorageObjectAccessControlDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_storage_bucket_acl" {
			continue
		}

		bucket := rs.Primary.Attributes["bucket"]
		object := rs.Primary.Attributes["object"]
		entity := rs.Primary.Attributes["entity"]

		rePairs, err := config.clientStorage.ObjectAccessControls.List(bucket, object).Do()
		if err != nil {
			return fmt.Errorf("Can't list role entity acl for object %s in bucket %s", object, bucket)
		}

		for _, v := range rePairs.Items {
			if v.Entity == entity {
				return fmt.Errorf("found entity %s as role entity acl entry for object %s in bucket %s", entity, object, bucket)
			}
		}

	}

	return nil
}

func testGoogleStorageObjectAccessControlBasic(bucketName, objectName, role, entity string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
	name = "%s"
}

resource "google_storage_bucket_object" "object" {
	name = "%s"
	bucket = "${google_storage_bucket.bucket.name}"
	source = "%s"
}

resource "google_storage_object_access_control" "default" {
	object = "${google_storage_bucket_object.object.name}"
	bucket = "${google_storage_bucket.bucket.name}"
	role   = "%s"
	entity = "%s"
}
`, bucketName, objectName, tfObjectAcl.Name(), role, entity)
}
