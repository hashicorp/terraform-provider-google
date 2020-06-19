package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRedisInstanceDatasource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstanceDatasourceConfig(randString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState("data.google_redis_instance.redis", "google_redis_instance.redis"),
				),
			},
		},
	})
}

func testAccRedisInstanceDatasourceConfig(suffix string) string {
	return fmt.Sprintf(`
resource "google_redis_instance" "redis" {
  name               = "redis-test-%s"
  memory_size_gb     = 1
}

data "google_redis_instance" "redis" {
  name = "${google_redis_instance.redis.name}"
}
`, suffix)
}
