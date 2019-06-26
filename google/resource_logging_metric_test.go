package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccLoggingMetric_update(t *testing.T) {
	t.Parallel()

	suffix := acctest.RandString(10)
	filter := "resource.type=gae_app AND severity>=ERROR"
	updatedFilter := "resource.type=gae_app AND severity=ERROR"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingMetricDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingMetric_update(suffix, filter),
			},
			{
				ResourceName:      "google_logging_metric.logging_metric",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLoggingMetric_update(suffix, updatedFilter),
			},
			{
				ResourceName:      "google_logging_metric.logging_metric",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccLoggingMetric_update(suffix string, filter string) string {
	return fmt.Sprintf(`
resource "google_logging_metric" "logging_metric" {
	name = "my-custom-metric-%s"
	filter = "%s"
	metric_descriptor {
		metric_kind = "DELTA"
		value_type = "INT64"
	}
}`, suffix, filter)
}
