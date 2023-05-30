// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccStorageDefaultObjectAccessControl_update(t *testing.T) {
	t.Parallel()

	bucketName := testBucketName(t)
	VcrTest(t, resource.TestCase{
		PreCheck: func() {
			if errObjectAcl != nil {
				panic(errObjectAcl)
			}
			acctest.AccTestPreCheck(t)
		},
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckStorageDefaultObjectAccessControlDestroyProducer(t),
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
  name     = "%s"
  location = "US"
}

resource "google_storage_default_object_access_control" "default" {
  bucket = google_storage_bucket.bucket.name
  role   = "%s"
  entity = "%s"
}
`, bucketName, role, entity)
}
