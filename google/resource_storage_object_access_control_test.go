package google

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"io/ioutil"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccStorageObjectAccessControl_update(t *testing.T) {
	t.Parallel()

	bucketName := testBucketName(t)
	objectName := fmt.Sprintf("%s-%d", "tf-test-acl-object", RandInt(t))
	objectData := []byte("data data data")
	if err := ioutil.WriteFile(tfObjectAcl.Name(), objectData, 0644); err != nil {
		t.Errorf("error writing file: %v", err)
	}
	VcrTest(t, resource.TestCase{
		PreCheck: func() {
			if errObjectAcl != nil {
				panic(errObjectAcl)
			}
			acctest.AccTestPreCheck(t)
		},
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckStorageObjectAccessControlDestroyProducer(t),
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

func TestAccStorageObjectAccessControl_updateWithSlashes(t *testing.T) {
	t.Parallel()

	bucketName := testBucketName(t)
	objectName := fmt.Sprintf("%s-%d", "tf-test/acl/object", RandInt(t))
	objectData := []byte("data data data")
	if err := ioutil.WriteFile(tfObjectAcl.Name(), objectData, 0644); err != nil {
		t.Errorf("error writing file: %v", err)
	}
	VcrTest(t, resource.TestCase{
		PreCheck: func() {
			if errObjectAcl != nil {
				panic(errObjectAcl)
			}
			acctest.AccTestPreCheck(t)
		},
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckStorageObjectAccessControlDestroyProducer(t),
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
  name     = "%s"
  location = "US"
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
