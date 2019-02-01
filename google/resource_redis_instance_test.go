package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
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
}`, name)
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
}`, name)
}
