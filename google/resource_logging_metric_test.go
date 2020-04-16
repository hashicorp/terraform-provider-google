package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccLoggingMetric_update(t *testing.T) {
	t.Parallel()

	suffix := randString(t, 10)
	filter := "resource.type=gae_app AND severity>=ERROR"
	updatedFilter := "resource.type=gae_app AND severity=ERROR"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingMetricDestroyProducer(t),
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

func TestAccLoggingMetric_explicitBucket(t *testing.T) {
	t.Parallel()

	suffix := randString(t, 10)
	filter := "resource.type=gae_app AND severity>=ERROR"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingMetricDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingMetric_explicitBucket(suffix, filter),
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
  name   = "my-custom-metric-%s"
  filter = "%s"
  metric_descriptor {
    metric_kind  = "DELTA"
    value_type   = "INT64"
    display_name = "My metric"
  }
}
`, suffix, filter)
}

func testAccLoggingMetric_explicitBucket(suffix string, filter string) string {
	return fmt.Sprintf(`
resource "google_logging_metric" "logging_metric" {
  name   = "my-custom-metric-%s"
  filter = "%s"

  metric_descriptor {
    metric_kind = "DELTA"
    value_type  = "DISTRIBUTION"
  }

  value_extractor = "EXTRACT(jsonPayload.metrics.running_jobs)"

  bucket_options {
    explicit_buckets {
      bounds = [0, 1, 2, 3, 4.2]
    }
  }
}
`, suffix, filter)
}
