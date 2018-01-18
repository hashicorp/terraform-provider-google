package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccGoogleStorageDefaultObjectAcl_basic(t *testing.T) {
	t.Parallel()

	bucketName := testBucketName()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleStorageDefaultObjectAclDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testGoogleStorageDefaultObjectsAclBasic1(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageDefaultObjectAcl(bucketName, roleEntityBasic1),
					testAccCheckGoogleStorageDefaultObjectAcl(bucketName, roleEntityBasic2),
				),
			},
		},
	})
}

func TestAccGoogleStorageDefaultObjectAcl_upgrade(t *testing.T) {
	t.Parallel()

	bucketName := testBucketName()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleStorageDefaultObjectAclDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testGoogleStorageDefaultObjectsAclBasic1(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageDefaultObjectAcl(bucketName, roleEntityBasic1),
					testAccCheckGoogleStorageDefaultObjectAcl(bucketName, roleEntityBasic2),
				),
			},

			resource.TestStep{
				Config: testGoogleStorageDefaultObjectsAclBasic2(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageDefaultObjectAcl(bucketName, roleEntityBasic2),
					testAccCheckGoogleStorageDefaultObjectAcl(bucketName, roleEntityBasic3_owner),
				),
			},

			resource.TestStep{
				Config: testGoogleStorageDefaultObjectsAclBasicDelete(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageDefaultObjectAclDelete(bucketName, roleEntityBasic1),
					testAccCheckGoogleStorageDefaultObjectAclDelete(bucketName, roleEntityBasic2),
					testAccCheckGoogleStorageDefaultObjectAclDelete(bucketName, roleEntityBasic3_reader),
				),
			},
		},
	})
}

func TestAccGoogleStorageDefaultObjectAcl_downgrade(t *testing.T) {
	t.Parallel()

	bucketName := testBucketName()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleStorageDefaultObjectAclDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testGoogleStorageDefaultObjectsAclBasic2(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageDefaultObjectAcl(bucketName, roleEntityBasic2),
					testAccCheckGoogleStorageDefaultObjectAcl(bucketName, roleEntityBasic3_owner),
				),
			},

			resource.TestStep{
				Config: testGoogleStorageDefaultObjectsAclBasic3(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageDefaultObjectAcl(bucketName, roleEntityBasic2),
					testAccCheckGoogleStorageDefaultObjectAcl(bucketName, roleEntityBasic3_reader),
				),
			},

			resource.TestStep{
				Config: testGoogleStorageDefaultObjectsAclBasicDelete(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageDefaultObjectAclDelete(bucketName, roleEntityBasic1),
					testAccCheckGoogleStorageDefaultObjectAclDelete(bucketName, roleEntityBasic2),
					testAccCheckGoogleStorageDefaultObjectAclDelete(bucketName, roleEntityBasic3_reader),
				),
			},
		},
	})
}

func testAccCheckGoogleStorageDefaultObjectAcl(bucket, roleEntityS string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		roleEntity, _ := getRoleEntityPair(roleEntityS)
		config := testAccProvider.Meta().(*Config)

		res, err := config.clientStorage.DefaultObjectAccessControls.Get(bucket,
			roleEntity.Entity).Do()

		if err != nil {
			return fmt.Errorf("Error retrieving contents of storage default Acl for bucket %s: %s", bucket, err)
		}

		if res.Role != roleEntity.Role {
			return fmt.Errorf("Error, Role mismatch %s != %s", res.Role, roleEntity.Role)
		}

		return nil
	}
}

func testAccGoogleStorageDefaultObjectAclDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "google_storage_default_object_acl" {
			continue
		}

		bucket := rs.Primary.Attributes["bucket"]

		_, err := config.clientStorage.DefaultObjectAccessControls.List(bucket).Do()
		if err == nil {
			return fmt.Errorf("Default Storage Object Acl for bucket %s still exists", bucket)
		}
	}
	return nil
}

func testAccCheckGoogleStorageDefaultObjectAclDelete(bucket, roleEntityS string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		roleEntity, _ := getRoleEntityPair(roleEntityS)
		config := testAccProvider.Meta().(*Config)

		_, err := config.clientStorage.DefaultObjectAccessControls.Get(bucket, roleEntity.Entity).Do()

		if err != nil {
			return nil
		}

		return fmt.Errorf("Error, Object Default Acl Entity still exists %s for bucket %s",
			roleEntity.Entity, bucket)
	}
}

func testGoogleStorageDefaultObjectsAclBasicDelete(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
	name = "%s"
}

resource "google_storage_default_object_acl" "acl" {
	bucket = "${google_storage_bucket.bucket.name}"
	role_entity = []
}
`, bucketName)
}

func testGoogleStorageDefaultObjectsAclBasic1(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
	name = "%s"
}

resource "google_storage_default_object_acl" "acl" {
	bucket = "${google_storage_bucket.bucket.name}"
	role_entity = ["%s", "%s"]
}
`, bucketName, roleEntityBasic1, roleEntityBasic2)
}

func testGoogleStorageDefaultObjectsAclBasic2(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
	name = "%s"
}

resource "google_storage_default_object_acl" "acl" {
	bucket = "${google_storage_bucket.bucket.name}"
	role_entity = ["%s", "%s"]
}
`, bucketName, roleEntityBasic2, roleEntityBasic3_owner)
}

func testGoogleStorageDefaultObjectsAclBasic3(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
	name = "%s"
}

resource "google_storage_default_object_acl" "acl" {
	bucket = "${google_storage_bucket.bucket.name}"
	role_entity = ["%s", "%s"]
}
`, bucketName, roleEntityBasic2, roleEntityBasic3_reader)
}
