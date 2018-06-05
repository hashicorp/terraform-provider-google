package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccRedisInstance_basic(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeAddressDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRedisInstance_basic(name),
			},
			resource.TestStep{
				ResourceName:      "google_redis_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRedisInstance_full(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("tf-test")
	network := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeAddressDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRedisInstance_full(name, network),
			},
			resource.TestStep{
				ResourceName:      "google_redis_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRedisInstance_basic(name string) string {
	return fmt.Sprintf(`
resource "google_redis_instance" "test" {
	name           = "%s"
	memory_size_gb = 1
}`, name)
}

func testAccRedisInstance_full(name, network string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "test" {
	name = "%s"
}

resource "google_redis_instance" "test" {
	name           = "%s"
	tier           = "STANDARD_HA"
	memory_size_gb = 1

	region                  = "us-central1"
	location_id             = "us-central1-a"
	alternative_location_id = "us-central1-f"

	redis_version     = "REDIS_3_2"
	display_name      = "Terraform Test Instance"
	reserved_ip_range = "192.168.0.0/29"

	labels {
		my_key    = "my_val"
		other_key = "other_val"
	}
}`, name, network)
}
