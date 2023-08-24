// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package monitoring_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Stackdriver tests cannot be run in parallel otherwise they will error out with:
// Error 503: Too many concurrent edits to the project configuration. Please try again.

func TestAccMonitoringAlertPolicy(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"basic":    testAccMonitoringAlertPolicy_basic,
		"full":     testAccMonitoringAlertPolicy_full,
		"update":   testAccMonitoringAlertPolicy_update,
		"mql":      testAccMonitoringAlertPolicy_mql,
		"log":      testAccMonitoringAlertPolicy_log,
		"forecast": testAccMonitoringAlertPolicy_forecast,
		"promql":   testAccMonitoringAlertPolicy_promql,
	}

	for name, tc := range testCases {
		// shadow the tc variable into scope so that when
		// the loop continues, if t.Run hasn't executed tc(t)
		// yet, we don't have a race condition
		// see https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc(t)
		})
	}
}

func testAccMonitoringAlertPolicy_basic(t *testing.T) {

	alertName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	conditionName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	filter := `metric.type=\"compute.googleapis.com/instance/disk/write_bytes_count\" AND resource.type=\"gce_instance\"`

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlertPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringAlertPolicy_basicCfg(alertName, conditionName, "ALIGN_RATE", filter),
			},
			{
				ResourceName:      "google_monitoring_alert_policy.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMonitoringAlertPolicy_update(t *testing.T) {

	alertName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	conditionName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	filter1 := `metric.type=\"compute.googleapis.com/instance/disk/write_bytes_count\" AND resource.type=\"gce_instance\"`
	aligner1 := "ALIGN_RATE"
	filter2 := `metric.type=\"compute.googleapis.com/instance/cpu/utilization\" AND resource.type=\"gce_instance\"`
	aligner2 := "ALIGN_MAX"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlertPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringAlertPolicy_basicCfg(alertName, conditionName, aligner1, filter1),
			},
			{
				ResourceName:      "google_monitoring_alert_policy.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMonitoringAlertPolicy_basicCfg(alertName, conditionName, aligner2, filter2),
			},
			{
				ResourceName:      "google_monitoring_alert_policy.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMonitoringAlertPolicy_full(t *testing.T) {

	alertName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	conditionName1 := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	conditionName2 := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlertPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringAlertPolicy_fullCfg(alertName, conditionName1, conditionName2),
			},
			{
				ResourceName:      "google_monitoring_alert_policy.full",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMonitoringAlertPolicy_mql(t *testing.T) {

	alertName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	conditionName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlertPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringAlertPolicy_mqlCfg(alertName, conditionName),
			},
			{
				ResourceName:      "google_monitoring_alert_policy.mql",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMonitoringAlertPolicy_log(t *testing.T) {

	alertName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	conditionName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlertPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringAlertPolicy_logCfg(alertName, conditionName),
			},
			{
				ResourceName:      "google_monitoring_alert_policy.log",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAlertPolicyDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_monitoring_alert_policy" {
				continue
			}

			name := rs.Primary.Attributes["name"]

			url := fmt.Sprintf("https://monitoring.googleapis.com/v3/%s", name)
			_, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				RawURL:    url,
				UserAgent: config.UserAgent,
			})

			if err == nil {
				return fmt.Errorf("Error, alert policy %s still exists", name)
			}
		}

		return nil
	}
}

func testAccMonitoringAlertPolicy_forecast(t *testing.T) {

	alertName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	conditionName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	filter := `metric.type=\"compute.googleapis.com/instance/disk/write_bytes_count\" AND resource.type=\"gce_instance\"`

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlertPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringAlertPolicy_forecastCfg(alertName, conditionName, "ALIGN_RATE", filter),
			},
			{
				ResourceName:      "google_monitoring_alert_policy.forecast",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMonitoringAlertPolicy_promql(t *testing.T) {

	alertName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	conditionName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlertPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringAlertPolicy_promqlCfg(alertName, conditionName),
			},
			{
				ResourceName:      "google_monitoring_alert_policy.promql",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMonitoringAlertPolicy_basicCfg(alertName, conditionName, aligner, filter string) string {
	return fmt.Sprintf(`
resource "google_monitoring_alert_policy" "basic" {
  display_name = "%s"
  enabled      = true
  combiner     = "OR"

  conditions {
    display_name = "%s"

    condition_threshold {
      aggregations {
        alignment_period   = "60s"
        per_series_aligner = "%s"
      }

      duration        = "60s"
      comparison      = "COMPARISON_GT"
      filter          = "%s"
      threshold_value = "0.5"
    }
  }
}
`, alertName, conditionName, aligner, filter)
}

func testAccMonitoringAlertPolicy_fullCfg(alertName, conditionName1, conditionName2 string) string {
	return fmt.Sprintf(`
resource "google_monitoring_alert_policy" "full" {
  display_name = "%s"
  combiner     = "OR"
  enabled      = true

  conditions {
    display_name = "%s"

    condition_threshold {
      threshold_value = 50
      filter          = "metric.type=\"compute.googleapis.com/instance/disk/write_bytes_count\" AND resource.type=\"gce_instance\""
      duration        = "60s"
      comparison      = "COMPARISON_GT"

      aggregations {
        alignment_period     = "60s"
        per_series_aligner   = "ALIGN_RATE"
        cross_series_reducer = "REDUCE_MEAN"

        group_by_fields = [
          "metric.label.device_name",
          "project",
          "resource.label.instance_id",
          "resource.label.zone",
        ]
      }

      trigger {
        percent = 10
      }
    }
  }

  conditions {
    display_name = "%s"

    condition_absent {
      duration = "3600s"
      filter   = "metric.type=\"compute.googleapis.com/instance/cpu/utilization\" AND resource.type=\"gce_instance\""

      aggregations {
        alignment_period     = "60s"
        cross_series_reducer = "REDUCE_NONE"
        per_series_aligner   = "ALIGN_MEAN"

        group_by_fields = [
          "project",
          "resource.label.instance_id",
          "resource.label.zone",
        ]
      }

      trigger {
        count = 1
      }
    }
  }

  documentation {
    content   = "test content"
    mime_type = "text/markdown"
  }
}
`, alertName, conditionName1, conditionName2)
}

func testAccMonitoringAlertPolicy_mqlCfg(alertName, conditionName string) string {
	return fmt.Sprintf(`
resource "google_monitoring_alert_policy" "mql" {
  display_name = "%s"
  combiner     = "OR"
  enabled      = true

  conditions {
    display_name = "%s"

    condition_monitoring_query_language {
      query           = "fetch gce_instance::compute.googleapis.com/instance/cpu/utilization | align mean_aligner() | window 5m | condition value.utilization > .15 '10^2.%%'"
      duration        = "60s"

      trigger {
        count = 2
      }
    }
  }

  documentation {
    content   = "test content"
    mime_type = "text/markdown"
  }
}
`, alertName, conditionName)
}

func testAccMonitoringAlertPolicy_logCfg(alertName, conditionName string) string {
	return fmt.Sprintf(`
resource "google_monitoring_alert_policy" "log" {
  display_name = "%s"
  combiner     = "OR"
  enabled      = true

  conditions {
    display_name = "%s"

    condition_matched_log {
      filter = "protoPayload.methodName=\"google.cloud.bigquery.v2.TableService.DeleteTable\""
      label_extractors = {
        "test" = "EXTRACT(protoPayload.request)"
      }
    }
  }

  alert_strategy {
    notification_rate_limit {
      period = "300s"
    }
    auto_close = "2000s"
  }

  documentation {
    content   = "test content"
    mime_type = "text/markdown"
  }
}
`, alertName, conditionName)
}

func testAccMonitoringAlertPolicy_forecastCfg(alertName, conditionName, aligner, filter string) string {
	return fmt.Sprintf(`
resource "google_monitoring_alert_policy" "forecast" {
  display_name = "%s"
  enabled      = true
  combiner     = "OR"

  conditions {
    display_name = "%s"

    condition_threshold {
      aggregations {
        alignment_period   = "60s"
        per_series_aligner = "%s"
      }

      duration        = "60s"
      forecast_options {
        forecast_horizon = "3600s"
      }
      comparison      = "COMPARISON_GT"
      filter          = "%s"
      threshold_value = "0.5"
    }
  }
}
`, alertName, conditionName, aligner, filter)
}

func testAccMonitoringAlertPolicy_promqlCfg(alertName, conditionName string) string {
	return fmt.Sprintf(`
resource "google_monitoring_alert_policy" "promql" {
  display_name = "%s"
  combiner     = "OR"
  enabled      = true

  conditions {
    display_name = "%s"
    
    condition_prometheus_query_language {
      query           = "vector(1)"
      duration        = "60s"
      labels = {
        "severity" = "page"
      }
      alert_rule      = "AlwaysOn"
      rule_group      = "abc"
    }
  }

  documentation {
    content   = "test content"
    mime_type = "text/markdown"
  }
}
`, alertName, conditionName)
}
