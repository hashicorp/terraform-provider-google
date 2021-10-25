package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceStorageBucketObjectContent_Basic(t *testing.T) {

	bucket := "tf-bucket-object-content-" + randString(t, 10)
	content := "qwertyuioasdfghjk1234567!!@#$*"

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceStorageBucketObjectContent_Basic(content, bucket),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_storage_bucket_object_content.default", "content"),
					resource.TestCheckResourceAttr("data.google_storage_bucket_object_content.default", "content", content),
				),
			},
		},
	})
}

func testAccDataSourceStorageBucketObjectContent_Basic(content, bucket string) string {
	return fmt.Sprintf(`
data "google_storage_bucket_object_content" "default" {
	bucket = google_storage_bucket.contenttest.name
	name   = google_storage_bucket_object.object.name      
}

resource "google_storage_bucket_object" "object" {
	name    = "butterfly01"
	content = "%s"
	bucket  = google_storage_bucket.contenttest.name
}

resource "google_storage_bucket" "contenttest" {
	name          = "%s"
	location      = "US"
	force_destroy = true
}`, content, bucket)
}
