package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRedisInstance_update(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRedisInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstance_update(name),
			},
			{
				ResourceName:      "google_redis_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRedisInstance_update2(name),
			},
			{
				ResourceName:      "google_redis_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRedisInstance_regionFromLocation(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("tf-test")

	// Pick a zone that isn't in the provider-specified region so we know we
	// didn't fall back to that one.
	region := "us-west1"
	zone := "us-west1-a"
	if getTestRegionFromEnv() == "us-west1" {
		region = "us-central1"
		zone = "us-central1-a"
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRedisInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstance_regionFromLocation(name, zone),
				Check:  resource.TestCheckResourceAttr("google_redis_instance.test", "region", region),
			},
			{
				ResourceName:      "google_redis_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRedisInstance_update(name string) string {
	return fmt.Sprintf(`
resource "google_redis_instance" "test" {
  name           = "%s"
  display_name   = "pre-update"
  memory_size_gb = 1
  region         = "us-central1"

  labels = {
    my_key    = "my_val"
    other_key = "other_val"
  }

  redis_configs = {
    maxmemory-policy       = "allkeys-lru"
    notify-keyspace-events = "KEA"
  }
}
`, name)
}

func testAccRedisInstance_update2(name string) string {
	return fmt.Sprintf(`
resource "google_redis_instance" "test" {
  name           = "%s"
  display_name   = "post-update"
  memory_size_gb = 1

  labels = {
    my_key    = "my_val"
    other_key = "new_val"
  }

  redis_configs = {
    maxmemory-policy       = "noeviction"
    notify-keyspace-events = ""
  }
}
`, name)
}

func testAccRedisInstance_regionFromLocation(name, zone string) string {
	return fmt.Sprintf(`
resource "google_redis_instance" "test" {
  name           = "%s"
  memory_size_gb = 1
  location_id    = "%s"
}
`, name, zone)
}
