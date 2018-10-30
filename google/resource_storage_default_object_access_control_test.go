package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccStorageDefaultObjectAccessControl_basic(t *testing.T) {
	t.Parallel()

	bucketName := testBucketName()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if errObjectAcl != nil {
				panic(errObjectAcl)
			}
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageDefaultObjectAccessControlDestroy,
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageDefaultObjectAccessControlBasic(bucketName, "READER", "allUsers"),
			},
			{
				ResourceName:      "google_storage_default_object_access_control.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageDefaultObjectAccessControl_update(t *testing.T) {
	t.Parallel()

	bucketName := testBucketName()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if errObjectAcl != nil {
				panic(errObjectAcl)
			}
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageDefaultObjectAccessControlDestroy,
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

func testAccStorageDefaultObjectAccessControlDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_storage_bucket_acl" {
			continue
		}

		bucket := rs.Primary.Attributes["bucket"]
		entity := rs.Primary.Attributes["entity"]

		rePairs, err := config.clientStorage.DefaultObjectAccessControls.List(bucket).Do()
		if err != nil {
			return fmt.Errorf("Can't list role entity acl for bucket %s", bucket)
		}

		for _, v := range rePairs.Items {
			if v.Entity == entity {
				return fmt.Errorf("found entity %s as role entity acl entry in bucket %s", entity, bucket)
			}
		}

	}

	return nil
}

func testGoogleStorageDefaultObjectAccessControlBasic(bucketName, role, entity string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
	name = "%s"
}

resource "google_storage_default_object_access_control" "default" {
	bucket = "${google_storage_bucket.bucket.name}"
	role   = "%s"
	entity = "%s"
}
`, bucketName, role, entity)
}
