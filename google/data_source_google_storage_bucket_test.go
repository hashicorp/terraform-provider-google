package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGoogleStorageBucket_basic(t *testing.T) {
	t.Parallel()

	bucket := "tf-bucket-" + randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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
