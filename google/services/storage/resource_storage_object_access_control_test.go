// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccStorageObjectAccessControl_update(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	objectName := fmt.Sprintf("%s-%d", "tf-test-acl-object", acctest.RandInt(t))
	objectData := []byte("data data data")
	if err := ioutil.WriteFile(tfObjectAcl.Name(), objectData, 0644); err != nil {
		t.Errorf("error writing file: %v", err)
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck: func() {
			if errObjectAcl != nil {
				panic(errObjectAcl)
			}
			acctest.AccTestPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
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

	bucketName := acctest.TestBucketName(t)
	objectName := fmt.Sprintf("%s-%d", "tf-test/acl/object", acctest.RandInt(t))
	objectData := []byte("data data data")
	if err := ioutil.WriteFile(tfObjectAcl.Name(), objectData, 0644); err != nil {
		t.Errorf("error writing file: %v", err)
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck: func() {
			if errObjectAcl != nil {
				panic(errObjectAcl)
			}
			acctest.AccTestPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
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
