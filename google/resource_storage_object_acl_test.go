package google

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var tfObjectAcl, errObjectAcl = ioutil.TempFile("", "tf-gce-test")

func testAclObjectName() string {
	return fmt.Sprintf("%s-%d", "tf-test-acl-object",
		rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}

func TestAccStorageObjectAcl_basic(t *testing.T) {
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
		CheckDestroy: testAccStorageObjectAclDestroy,
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageObjectsAclBasic1(bucketName, objectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObjectAcl(bucketName,
						objectName, roleEntityBasic1),
					testAccCheckGoogleStorageObjectAcl(bucketName,
						objectName, roleEntityBasic2),
				),
			},
		},
	})
}

func TestAccStorageObjectAcl_upgrade(t *testing.T) {
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
		CheckDestroy: testAccStorageObjectAclDestroy,
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageObjectsAclBasic1(bucketName, objectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObjectAcl(bucketName,
						objectName, roleEntityBasic1),
					testAccCheckGoogleStorageObjectAcl(bucketName,
						objectName, roleEntityBasic2),
				),
			},

			{
				Config: testGoogleStorageObjectsAclBasic2(bucketName, objectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObjectAcl(bucketName,
						objectName, roleEntityBasic2),
					testAccCheckGoogleStorageObjectAcl(bucketName,
						objectName, roleEntityBasic3_owner),
				),
			},

			{
				Config: testGoogleStorageObjectsAclBasicDelete(bucketName, objectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObjectAclDelete(bucketName,
						objectName, roleEntityBasic1),
					testAccCheckGoogleStorageObjectAclDelete(bucketName,
						objectName, roleEntityBasic2),
					testAccCheckGoogleStorageObjectAclDelete(bucketName,
						objectName, roleEntityBasic3_reader),
				),
			},
		},
	})
}

func TestAccStorageObjectAcl_downgrade(t *testing.T) {
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
		CheckDestroy: testAccStorageObjectAclDestroy,
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageObjectsAclBasic2(bucketName, objectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObjectAcl(bucketName,
						objectName, roleEntityBasic2),
					testAccCheckGoogleStorageObjectAcl(bucketName,
						objectName, roleEntityBasic3_owner),
				),
			},

			{
				Config: testGoogleStorageObjectsAclBasic3(bucketName, objectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObjectAcl(bucketName,
						objectName, roleEntityBasic2),
					testAccCheckGoogleStorageObjectAcl(bucketName,
						objectName, roleEntityBasic3_reader),
				),
			},

			{
				Config: testGoogleStorageObjectsAclBasicDelete(bucketName, objectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObjectAclDelete(bucketName,
						objectName, roleEntityBasic1),
					testAccCheckGoogleStorageObjectAclDelete(bucketName,
						objectName, roleEntityBasic2),
					testAccCheckGoogleStorageObjectAclDelete(bucketName,
						objectName, roleEntityBasic3_reader),
				),
			},
		},
	})
}

func TestAccStorageObjectAcl_predefined(t *testing.T) {
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
		CheckDestroy: testAccStorageObjectAclDestroy,
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageObjectsAclPredefined(bucketName, objectName),
			},
		},
	})
}

func TestAccStorageObjectAcl_predefinedToExplicit(t *testing.T) {
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
		CheckDestroy: testAccStorageObjectAclDestroy,
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageObjectsAclPredefined(bucketName, objectName),
			},
			{
				Config: testGoogleStorageObjectsAclBasic1(bucketName, objectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObjectAcl(bucketName,
						objectName, roleEntityBasic1),
					testAccCheckGoogleStorageObjectAcl(bucketName,
						objectName, roleEntityBasic2),
				),
			},
		},
	})
}

func TestAccStorageObjectAcl_explicitToPredefined(t *testing.T) {
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
		CheckDestroy: testAccStorageObjectAclDestroy,
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageObjectsAclBasic1(bucketName, objectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObjectAcl(bucketName,
						objectName, roleEntityBasic1),
					testAccCheckGoogleStorageObjectAcl(bucketName,
						objectName, roleEntityBasic2),
				),
			},
			{
				Config: testGoogleStorageObjectsAclPredefined(bucketName, objectName),
			},
		},
	})
}

// Test that we allow the API to reorder our role entities without perma-diffing.
func TestAccStorageObjectAcl_unordered(t *testing.T) {
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
		CheckDestroy: testAccStorageObjectAclDestroy,
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageObjectAclUnordered(bucketName, objectName),
			},
		},
	})
}

func testAccCheckGoogleStorageObjectAcl(bucket, object, roleEntityS string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		roleEntity, _ := getRoleEntityPair(roleEntityS)
		config := testAccProvider.Meta().(*Config)

		res, err := config.clientStorage.ObjectAccessControls.Get(bucket,
			object, roleEntity.Entity).Do()

		if err != nil {
			return fmt.Errorf("Error retrieving contents of acl for bucket %s: %s", bucket, err)
		}

		if res.Role != roleEntity.Role {
			return fmt.Errorf("Error, Role mismatch %s != %s", res.Role, roleEntity.Role)
		}

		return nil
	}
}

func testAccCheckGoogleStorageObjectAclDelete(bucket, object, roleEntityS string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		roleEntity, _ := getRoleEntityPair(roleEntityS)
		config := testAccProvider.Meta().(*Config)

		_, err := config.clientStorage.ObjectAccessControls.Get(bucket,
			object, roleEntity.Entity).Do()

		if err != nil {
			return nil
		}

		return fmt.Errorf("Error, Entity still exists %s", roleEntity.Entity)
	}
}

func testAccStorageObjectAclDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_storage_bucket_acl" {
			continue
		}

		bucket := rs.Primary.Attributes["bucket"]
		object := rs.Primary.Attributes["object"]

		_, err := config.clientStorage.ObjectAccessControls.List(bucket, object).Do()

		if err == nil {
			return fmt.Errorf("Acl for bucket %s still exists", bucket)
		}
	}

	return nil
}

func testGoogleStorageObjectsAclBasicDelete(bucketName string, objectName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
}

resource "google_storage_bucket_object" "object" {
  name   = "%s"
  bucket = google_storage_bucket.bucket.name
  source = "%s"
}

resource "google_storage_object_acl" "acl" {
  object      = google_storage_bucket_object.object.name
  bucket      = google_storage_bucket.bucket.name
  role_entity = []
}
`, bucketName, objectName, tfObjectAcl.Name())
}

func testGoogleStorageObjectsAclBasic1(bucketName string, objectName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
}

resource "google_storage_bucket_object" "object" {
  name   = "%s"
  bucket = google_storage_bucket.bucket.name
  source = "%s"
}

resource "google_storage_object_acl" "acl" {
  object      = google_storage_bucket_object.object.name
  bucket      = google_storage_bucket.bucket.name
  role_entity = ["%s", "%s"]
}
`, bucketName, objectName, tfObjectAcl.Name(),
		roleEntityBasic1, roleEntityBasic2)
}

func testGoogleStorageObjectsAclBasic2(bucketName string, objectName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
}

resource "google_storage_bucket_object" "object" {
  name   = "%s"
  bucket = google_storage_bucket.bucket.name
  source = "%s"
}

resource "google_storage_object_acl" "acl" {
  object      = google_storage_bucket_object.object.name
  bucket      = google_storage_bucket.bucket.name
  role_entity = ["%s", "%s"]
}
`, bucketName, objectName, tfObjectAcl.Name(),
		roleEntityBasic2, roleEntityBasic3_owner)
}

func testGoogleStorageObjectsAclBasic3(bucketName string, objectName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
}

resource "google_storage_bucket_object" "object" {
  name   = "%s"
  bucket = google_storage_bucket.bucket.name
  source = "%s"
}

resource "google_storage_object_acl" "acl" {
  object      = google_storage_bucket_object.object.name
  bucket      = google_storage_bucket.bucket.name
  role_entity = ["%s", "%s"]
}
`, bucketName, objectName, tfObjectAcl.Name(),
		roleEntityBasic2, roleEntityBasic3_reader)
}

func testGoogleStorageObjectsAclPredefined(bucketName string, objectName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
}

resource "google_storage_bucket_object" "object" {
  name   = "%s"
  bucket = google_storage_bucket.bucket.name
  source = "%s"
}

resource "google_storage_object_acl" "acl" {
  object         = google_storage_bucket_object.object.name
  bucket         = google_storage_bucket.bucket.name
  predefined_acl = "projectPrivate"
}
`, bucketName, objectName, tfObjectAcl.Name())
}

func testGoogleStorageObjectAclUnordered(bucketName, objectName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
}

resource "google_storage_bucket_object" "object" {
  name   = "%s"
  bucket = google_storage_bucket.bucket.name
  source = "%s"
}

resource "google_storage_object_acl" "acl" {
  object      = google_storage_bucket_object.object.name
  bucket      = google_storage_bucket.bucket.name
  role_entity = ["%s", "%s", "%s", "%s", "%s"]
}
`, bucketName, objectName, tfObjectAcl.Name(), roleEntityBasic1, roleEntityViewers, roleEntityOwners, roleEntityBasic2, roleEntityEditors)
}
