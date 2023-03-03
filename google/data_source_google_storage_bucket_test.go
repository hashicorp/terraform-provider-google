package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGoogleStorageBucket_basic(t *testing.T) {
	t.Parallel()

	bucket := "tf-bucket-" + RandString(t, 10)

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccStorageBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleStorageBucketConfig(bucket),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceStateWithIgnores("data.google_storage_bucket.bar", "google_storage_bucket.foo", map[string]struct{}{"force_destroy": {}}),
				),
			},
		},
	})
}

func testAccDataSourceGoogleStorageBucketConfig(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "foo" {
  name     = "%s"
  location = "US"
}

data "google_storage_bucket" "bar" {
  name = google_storage_bucket.foo.name
  depends_on = [
    google_storage_bucket.foo,
  ]
}
`, bucketName)
}
