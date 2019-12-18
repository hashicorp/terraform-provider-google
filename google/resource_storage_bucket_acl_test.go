package google

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var (
	roleEntityBasic1        = "OWNER:user-paddy@hashicorp.com"
	roleEntityBasic2        = "READER:user-paddy@carvers.co"
	roleEntityBasic3_owner  = "OWNER:user-paddy@paddy.io"
	roleEntityBasic3_reader = "READER:user-foran.paddy@gmail.com"

	roleEntityOwners  = "OWNER:project-owners-" + os.Getenv("GOOGLE_PROJECT_NUMBER")
	roleEntityEditors = "OWNER:project-editors-" + os.Getenv("GOOGLE_PROJECT_NUMBER")
	roleEntityViewers = "READER:project-viewers-" + os.Getenv("GOOGLE_PROJECT_NUMBER")
)

func TestAccStorageBucketAcl_basic(t *testing.T) {
	t.Parallel()

	bucketName := testBucketName()
	skipIfEnvNotSet(t, "GOOGLE_PROJECT_NUMBER")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketAclDestroy,
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsAclBasic1(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageBucketAcl(bucketName, roleEntityBasic1),
					testAccCheckGoogleStorageBucketAcl(bucketName, roleEntityBasic2),
				),
			},
		},
	})
}

func TestAccStorageBucketAcl_upgrade(t *testing.T) {
	t.Parallel()

	bucketName := testBucketName()
	skipIfEnvNotSet(t, "GOOGLE_PROJECT_NUMBER")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketAclDestroy,
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsAclBasic1(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageBucketAcl(bucketName, roleEntityBasic1),
					testAccCheckGoogleStorageBucketAcl(bucketName, roleEntityBasic2),
				),
			},

			{
				Config: testGoogleStorageBucketsAclBasic2(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageBucketAcl(bucketName, roleEntityBasic2),
					testAccCheckGoogleStorageBucketAcl(bucketName, roleEntityBasic3_owner),
				),
			},

			{
				Config: testGoogleStorageBucketsAclBasicDelete(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageBucketAclDelete(bucketName, roleEntityBasic1),
					testAccCheckGoogleStorageBucketAclDelete(bucketName, roleEntityBasic2),
					testAccCheckGoogleStorageBucketAclDelete(bucketName, roleEntityBasic3_owner),
				),
			},
		},
	})
}

func TestAccStorageBucketAcl_downgrade(t *testing.T) {
	t.Parallel()

	bucketName := testBucketName()
	skipIfEnvNotSet(t, "GOOGLE_PROJECT_NUMBER")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketAclDestroy,
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsAclBasic2(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageBucketAcl(bucketName, roleEntityBasic2),
					testAccCheckGoogleStorageBucketAcl(bucketName, roleEntityBasic3_owner),
				),
			},

			{
				Config: testGoogleStorageBucketsAclBasic3(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageBucketAcl(bucketName, roleEntityBasic2),
					testAccCheckGoogleStorageBucketAcl(bucketName, roleEntityBasic3_reader),
				),
			},

			{
				Config: testGoogleStorageBucketsAclBasicDelete(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageBucketAclDelete(bucketName, roleEntityBasic1),
					testAccCheckGoogleStorageBucketAclDelete(bucketName, roleEntityBasic2),
					testAccCheckGoogleStorageBucketAclDelete(bucketName, roleEntityBasic3_owner),
				),
			},
		},
	})
}

func TestAccStorageBucketAcl_predefined(t *testing.T) {
	t.Parallel()

	bucketName := testBucketName()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketAclDestroy,
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsAclPredefined(bucketName),
			},
		},
	})
}

// Test that we allow the API to reorder our role entities without perma-diffing.
func TestAccStorageBucketAcl_unordered(t *testing.T) {
	t.Parallel()

	bucketName := testBucketName()
	skipIfEnvNotSet(t, "GOOGLE_PROJECT_NUMBER")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketAclDestroy,
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsAclUnordered(bucketName),
			},
		},
	})
}

func testAccCheckGoogleStorageBucketAclDelete(bucket, roleEntityS string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		roleEntity, _ := getRoleEntityPair(roleEntityS)
		config := testAccProvider.Meta().(*Config)

		_, err := config.clientStorage.BucketAccessControls.Get(bucket, roleEntity.Entity).Do()

		if err != nil {
			return nil
		}

		return fmt.Errorf("Error, entity %s still exists", roleEntity.Entity)
	}
}

func testAccCheckGoogleStorageBucketAcl(bucket, roleEntityS string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		roleEntity, _ := getRoleEntityPair(roleEntityS)
		config := testAccProvider.Meta().(*Config)

		res, err := config.clientStorage.BucketAccessControls.Get(bucket, roleEntity.Entity).Do()

		if err != nil {
			return fmt.Errorf("Error retrieving contents of acl for bucket %s: %s", bucket, err)
		}

		if res.Role != roleEntity.Role {
			return fmt.Errorf("Error, Role mismatch %s != %s", res.Role, roleEntity.Role)
		}

		return nil
	}
}

func testAccStorageBucketAclDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_storage_bucket_acl" {
			continue
		}

		bucket := rs.Primary.Attributes["bucket"]

		_, err := config.clientStorage.BucketAccessControls.List(bucket).Do()

		if err == nil {
			return fmt.Errorf("Acl for bucket %s still exists", bucket)
		}
	}

	return nil
}

func testGoogleStorageBucketsAclBasic1(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
}

resource "google_storage_bucket_acl" "acl" {
  bucket      = google_storage_bucket.bucket.name
  role_entity = ["%s", "%s", "%s", "%s", "%s"]
}
`, bucketName, roleEntityOwners, roleEntityEditors, roleEntityViewers, roleEntityBasic1, roleEntityBasic2)
}

func testGoogleStorageBucketsAclBasic2(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
}

resource "google_storage_bucket_acl" "acl" {
  bucket      = google_storage_bucket.bucket.name
  role_entity = ["%s", "%s", "%s", "%s", "%s"]
}
`, bucketName, roleEntityOwners, roleEntityEditors, roleEntityViewers, roleEntityBasic2, roleEntityBasic3_owner)
}

func testGoogleStorageBucketsAclBasicDelete(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
}

resource "google_storage_bucket_acl" "acl" {
  bucket      = google_storage_bucket.bucket.name
  role_entity = []
}
`, bucketName)
}

func testGoogleStorageBucketsAclBasic3(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
}

resource "google_storage_bucket_acl" "acl" {
  bucket      = google_storage_bucket.bucket.name
  role_entity = ["%s", "%s", "%s", "%s", "%s"]
}
`, bucketName, roleEntityOwners, roleEntityEditors, roleEntityViewers, roleEntityBasic2, roleEntityBasic3_reader)
}

func testGoogleStorageBucketsAclUnordered(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
}

resource "google_storage_bucket_acl" "acl" {
  bucket      = google_storage_bucket.bucket.name
  role_entity = ["%s", "%s", "%s", "%s", "%s"]
}
`, bucketName, roleEntityBasic1, roleEntityViewers, roleEntityOwners, roleEntityBasic2, roleEntityEditors)
}

func testGoogleStorageBucketsAclPredefined(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
}

resource "google_storage_bucket_acl" "acl" {
  bucket         = google_storage_bucket.bucket.name
  predefined_acl = "projectPrivate"
  default_acl    = "projectPrivate"
}
`, bucketName)
}
