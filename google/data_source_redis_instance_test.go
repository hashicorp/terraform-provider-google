package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccRedisInstanceDatasource_basic(t *testing.T) {
	t.Parallel()

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstanceDatasourceConfig(RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_redis_instance.redis", "google_redis_instance.redis"),
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
  name = google_redis_instance.redis.name
}
`, suffix)
}
