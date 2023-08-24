// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package logging_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccLoggingMetric_update(t *testing.T) {
	t.Parallel()

	suffix := acctest.RandString(t, 10)
	filter := "resource.type=gae_app AND severity>=ERROR"
	updatedFilter := "resource.type=gae_app AND severity=ERROR"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingMetricDestroyProducer(t),
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

	suffix := acctest.RandString(t, 10)
	filter := "resource.type=gae_app AND severity>=ERROR"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingMetricDestroyProducer(t),
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

func TestAccLoggingMetric_loggingBucket(t *testing.T) {
	t.Parallel()

	filter := "resource.type=gae_app AND severity>=ERROR"
	project_id := envvar.GetTestProjectFromEnv()
	suffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingMetricDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingMetric_loggingBucketBase(suffix, filter),
			},
			{
				ResourceName:      "google_logging_metric.logging_metric",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLoggingMetric_loggingBucket(suffix, filter, project_id),
			},
			{
				ResourceName:      "google_logging_metric.logging_metric",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLoggingMetric_loggingBucketBase(suffix, filter),
			},
			{
				ResourceName:      "google_logging_metric.logging_metric",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingMetric_descriptionUpdated(t *testing.T) {
	t.Parallel()

	suffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingMetricDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingMetric_descriptionUpdated(suffix, "original"),
			},
			{
				ResourceName:      "google_logging_metric.logging_metric",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLoggingMetric_descriptionUpdated(suffix, "Updated"),
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

func testAccLoggingMetric_loggingBucketBase(suffix string, filter string) string {
	return fmt.Sprintf(`
resource "google_logging_metric" "logging_metric" {
  name        = "my-custom-metric-%s"
  filter      = "%s"
}
`, suffix, filter)
}

func testAccLoggingMetric_loggingBucket(suffix string, filter string, project_id string) string {
	return fmt.Sprintf(`
resource "google_logging_project_bucket_config" "logging_bucket" {
  location  = "global"
  project   = "%s"
  bucket_id = "_Default"
}

resource "google_logging_metric" "logging_metric" {
  name        = "my-custom-metric-%s"
  bucket_name = google_logging_project_bucket_config.logging_bucket.id
  filter      = "%s"
}
`, project_id, suffix, filter)
}

func testAccLoggingMetric_descriptionUpdated(suffix, description string) string {
	return fmt.Sprintf(`
resource "google_logging_metric" "logging_metric" {
	name        = "my-custom-metric-%s"
	description = "Counter for VM instances that have hostError's"
	filter      = "resource.type=gce_instance AND protoPayload.methodName=compute.instances.hostError"
	metric_descriptor {
	  metric_kind = "DELTA"
	  value_type  = "INT64"
	  labels {
		key         = "instance"
		value_type  = "STRING"
		description = "%s"
	  }
	  labels {
		key         = "zone"
		value_type  = "STRING"
		description = "Availability zone of instance"
	  }
	  display_name = "VM Host Errors"
	}
	label_extractors = {
	  "instance" = "REGEXP_EXTRACT(protoPayload.resourceName, \"projects/.+/zones/.+/instances/(.+)\")"
	  "zone"     = "EXTRACT(resource.labels.zone)"
	}
  }
`, suffix, description)
}
