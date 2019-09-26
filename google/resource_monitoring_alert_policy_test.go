package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Stackdriver tests cannot be run in parallel otherwise they will error out with:
// Error 503: Too many concurrent edits to the project configuration. Please try again.

func TestAccMonitoringAlertPolicy(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"basic":  testAccMonitoringAlertPolicy_basic,
		"full":   testAccMonitoringAlertPolicy_full,
		"update": testAccMonitoringAlertPolicy_update,
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

	alertName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	conditionName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	filter := `metric.type=\"compute.googleapis.com/instance/disk/write_bytes_count\" AND resource.type=\"gce_instance\"`

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlertPolicyDestroy,
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

	alertName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	conditionName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	filter1 := `metric.type=\"compute.googleapis.com/instance/disk/write_bytes_count\" AND resource.type=\"gce_instance\"`
	aligner1 := "ALIGN_RATE"
	filter2 := `metric.type=\"compute.googleapis.com/instance/cpu/utilization\" AND resource.type=\"gce_instance\"`
	aligner2 := "ALIGN_MAX"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlertPolicyDestroy,
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

	alertName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	conditionName1 := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	conditionName2 := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlertPolicyDestroy,
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

func testAccCheckAlertPolicyDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_monitoring_alert_policy" {
			continue
		}

		name := rs.Primary.Attributes["name"]

		url := fmt.Sprintf("https://monitoring.googleapis.com/v3/%s", name)
		_, err := sendRequest(config, "GET", "", url, nil)

		if err == nil {
			return fmt.Errorf("Error, alert policy %s still exists", name)
		}
	}

	return nil
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
        cross_series_reducer = "REDUCE_MEAN"
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
