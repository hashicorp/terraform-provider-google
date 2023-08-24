// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package monitoring_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccMonitoringMetricDescriptor_update(t *testing.T) {
	// TODO: Fix requires a breaking change https://github.com/hashicorp/terraform-provider-google/issues/12139
	t.Skip()

	t.Parallel()
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringMetricDescriptorDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringMetricDescriptor_update("key1", "STRING",
					"description1", "30s", "30s"),
			},
			{
				ResourceName:            "google_monitoring_metric_descriptor.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata", "launch_stage"},
			},
			{
				Config: testAccMonitoringMetricDescriptor_update("key2", "INT64",
					"description2", "60s", "60s"),
			},
			{
				ResourceName:            "google_monitoring_metric_descriptor.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata", "launch_stage"},
			},
		},
	})
}

func testAccMonitoringMetricDescriptor_update(key, valueType, description,
	samplePeriod, ingestDelay string) string {
	return fmt.Sprintf(`
resource "google_monitoring_metric_descriptor" "basic" {
	description = "Daily sales records from all branch stores."
	display_name = "daily sales"
	type = "custom.googleapis.com/stores/daily_sales"
	metric_kind = "GAUGE"
	value_type = "DOUBLE"
	unit = "{USD}"
	labels {
		key = "%s"
		value_type = "%s"
		description = "%s"
	}
	launch_stage = "BETA"
	metadata {
		sample_period = "%s"
		ingest_delay = "%s"
	}
}
`, key, valueType, description, samplePeriod, ingestDelay,
	)
}
