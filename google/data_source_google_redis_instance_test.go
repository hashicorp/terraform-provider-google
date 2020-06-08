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
				Config: testAccRedisInstanceDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleRedisInstanceCheck("data.google_redis_instance.redis", "google_redis_instance.redis"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleRedisInstanceCheck(dataSourceName string, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", dataSourceName)
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", resourceName)
		}

		dsAttr := ds.Primary.Attributes
		rsAttr := rs.Primary.Attributes

		clusterAttrToCheck := []string{
			"name",
			"region",
			"host",
			"port",
		}

		for _, attr := range clusterAttrToCheck {
			if dsAttr[attr] != rsAttr[attr] {
				return fmt.Errorf(
					"%s is %s; want %s",
					attr,
					dsAttr[attr],
					rsAttr[attr],
				)
			}
		}

		return nil
	}
}

var testAccRedisInstanceDatasourceConfig = fmt.Sprintf(`
	resource "google_redis_instance" "redis" {
		name               = "redis-test-%s"
    memory_size_gb     = 1
    region             = "europe-west1"
	}

	data "google_redis_instance" "redis" {
		name = "${google_redis_instance.redis.name}"
	}
`, acctest.RandString(10))
