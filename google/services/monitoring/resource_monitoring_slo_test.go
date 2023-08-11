// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package monitoring_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func setTestCheckMonitoringSloId(res string, sloId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		updateId, err := getTestResourceMonitoringSloId(res, s)
		if err != nil {
			return err
		}
		*sloId = updateId
		return nil
	}
}

func testCheckMonitoringSloIdAfterUpdate(res string, sloId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		updateId, err := getTestResourceMonitoringSloId(res, s)
		if err != nil {
			return err
		}

		if sloId == nil {
			return fmt.Errorf("unexpected error, slo ID was not set")
		}

		if *sloId != updateId {
			return fmt.Errorf("unexpected mismatch in slo ID after update, resource was recreated. Initial %q, Updated %q",
				*sloId, updateId)
		}
		return nil
	}
}

func getTestResourceMonitoringSloId(res string, s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources[res]
	if !ok {
		return "", fmt.Errorf("not found: %s", res)
	}

	if rs.Primary.ID == "" {
		return "", fmt.Errorf("no ID is set for %s", res)
	}

	if v, ok := rs.Primary.Attributes["slo_id"]; ok {
		return v, nil
	}

	return "", fmt.Errorf("slo_id not set on resource %s", res)
}

func TestAccMonitoringSlo_basic(t *testing.T) {
	t.Parallel()

	var generatedId string
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringSloDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringSlo_basic(),
				Check:  setTestCheckMonitoringSloId("google_monitoring_slo.primary", &generatedId),
			},
			{
				ResourceName:      "google_monitoring_slo.primary",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"service"},
			},
			{
				Config: testAccMonitoringSlo_basicUpdate(),
				Check:  testCheckMonitoringSloIdAfterUpdate("google_monitoring_slo.primary", &generatedId),
			},
			{
				ResourceName:      "google_monitoring_slo.primary",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"service"},
			},
		},
	})
}

func TestAccMonitoringSlo_availabilitySli(t *testing.T) {
	t.Parallel()

	var generatedId string
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringSloDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringSlo_availabilitySli(),
				Check:  setTestCheckMonitoringSloId("google_monitoring_slo.primary", &generatedId),
			},
			{
				ResourceName:      "google_monitoring_slo.primary",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"service"},
			},
			{
				Config: testAccMonitoringSlo_basicUpdate(),
				Check:  testCheckMonitoringSloIdAfterUpdate("google_monitoring_slo.primary", &generatedId),
			},
			{
				ResourceName:      "google_monitoring_slo.primary",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"service"},
			},
		},
	})
}
func TestAccMonitoringSlo_requestBased(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringSloDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringSloForSli(
					randomSuffix,
					testAccMonitoringSloSli_requestBasedDistributionMaxOnly()),
			},
			{
				ResourceName:      "google_monitoring_slo.test_slo",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"service"},
			},
			{
				Config: testAccMonitoringSloForSli(
					randomSuffix,
					testAccMonitoringSloSli_requestBasedGoodTotalRatio_goodAndTotal()),
			},
			{
				ResourceName:      "google_monitoring_slo.test_slo",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"service"},
			},
			{
				Config: testAccMonitoringSloForSli(
					randomSuffix,
					testAccMonitoringSloSli_requestBasedGoodTotalRatio_goodAndBad()),
			},
			{
				ResourceName:      "google_monitoring_slo.test_slo",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"service"},
			},
		},
	})
}

func TestAccMonitoringSlo_windowBased_updateSlis(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringSloDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringSloForSli(
					randomSuffix,
					testAccMonitoringSloSli_windowBasedSliGoodBadMetricFilter(),
				),
			},
			{
				ResourceName:      "google_monitoring_slo.test_slo",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"service"},
			},
			{
				Config: testAccMonitoringSloForSli(
					randomSuffix,
					testAccMonitoringSloSli_windowBasedSliGoodTotalRatioThreshold_goodBad(),
				),
			},
			{
				ResourceName:      "google_monitoring_slo.test_slo",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"service"},
			},
			{
				Config: testAccMonitoringSloForSli(
					randomSuffix,
					testAccMonitoringSloSli_windowBasedSliMetricMeanRange(),
				),
			},
			{
				ResourceName:      "google_monitoring_slo.test_slo",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"service"},
			},
			{
				Config: testAccMonitoringSloForSli(
					randomSuffix,
					testAccMonitoringSloSli_windowBasedSliMetricSumRange(),
				),
			},
			{
				ResourceName:      "google_monitoring_slo.test_slo",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"service"},
			},
		},
	})
}

func TestAccMonitoringSlo_windowBasedGoodTotalRatioThresholdSlis(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringSloDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringSloForSli(
					randomSuffix,
					testAccMonitoringSloSli_windowBasedSliGoodTotalRatioThreshold_distributionCut(),
				),
			},
			{
				ResourceName:      "google_monitoring_slo.test_slo",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"service"},
			},
			{
				Config: testAccMonitoringSloForSli(
					randomSuffix,
					testAccMonitoringSloSli_windowBasedSliGoodTotalRatioThreshold_distributionCutMaxOnly(),
				),
			},
			{
				ResourceName:      "google_monitoring_slo.test_slo",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"service"},
			},
			{
				Config: testAccMonitoringSloForSli(
					randomSuffix,
					testAccMonitoringSloSli_windowBasedSliGoodTotalRatioThreshold_goodTotal(),
				),
			},
			{
				ResourceName:      "google_monitoring_slo.test_slo",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"service"},
			},
			{
				Config: testAccMonitoringSloForSli(
					randomSuffix,
					testAccMonitoringSloSli_windowBasedSliGoodBadMetricFilter(),
				),
			},
			{
				ResourceName:      "google_monitoring_slo.test_slo",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"service"},
			},
			{
				Config: testAccMonitoringSloSli_windowBasedSliGoodTotalRatioThreshold_basicSli(),
			},
			{
				ResourceName:      "google_monitoring_slo.test_slo",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"service"},
			},
		},
	})
}

func TestAccMonitoringSlo_windowBasedMetricMeanRangeSlis(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringSloDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringSloForSli(
					randomSuffix,
					testAccMonitoringSloSli_windowBasedSliMetricMeanRange(),
				),
			},
			{
				ResourceName:      "google_monitoring_slo.test_slo",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"service"},
			},
			{
				Config: testAccMonitoringSloForSli(
					randomSuffix,
					testAccMonitoringSloSli_windowBasedSliMetricMeanRangeUpdate(),
				),
			},
			{
				ResourceName:      "google_monitoring_slo.test_slo",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"service"},
			},
		},
	})
}

func TestAccMonitoringSlo_windowBasedMetricSumRangeSlis(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringSloDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringSloForSli(
					randomSuffix,
					testAccMonitoringSloSli_windowBasedSliMetricSumRange(),
				),
			},
			{
				ResourceName:      "google_monitoring_slo.test_slo",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"service"},
			},
			{
				Config: testAccMonitoringSloForSli(
					randomSuffix,
					testAccMonitoringSloSli_windowBasedSliMetricSumRangeUpdate(),
				),
			},
			{
				ResourceName:      "google_monitoring_slo.test_slo",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"service"},
			},
		},
	})
}

func TestAccMonitoringSlo_genericService(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringSloDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringSlo_generic(randomSuffix),
			},
			{
				ResourceName:      "google_monitoring_slo.primary",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"service"},
			},
		},
	})
}

func testAccMonitoringSlo_basic() string {
	return `
data "google_monitoring_app_engine_service" "ae" {
  module_id = "default"
}

resource "google_monitoring_slo" "primary" {
  service = data.google_monitoring_app_engine_service.ae.service_id

  goal = 0.9
  rolling_period_days = 1

  basic_sli {
    latency {
      threshold = "1s"
    }
  }
  user_labels = {
    my_key       = "my_value"
    my_other_key = "my_other_value"
  }
}
`
}

func testAccMonitoringSlo_basicUpdate() string {
	return `
data "google_monitoring_app_engine_service" "ae" {
  module_id = "default"
}

resource "google_monitoring_slo" "primary" {
  service = data.google_monitoring_app_engine_service.ae.service_id

  goal = 0.8
  display_name = "Terraform Test updated SLO"
  calendar_period = "WEEK"

  basic_sli {
    latency {
      threshold = "2s"
    }
  }
}
`
}

func testAccMonitoringSlo_generic(randSuffix string) string {
	return fmt.Sprintf(`
resource "google_monitoring_service" "srv" {
	service_id = "tf-test-srv-%s"
	display_name = "My Basic CloudEnpoints Service"
	basic_service {
		service_type  = "CLOUD_ENDPOINTS"
		service_labels = {
			service = "another-endpoint"
		}
	}
}
	  
	
resource "google_monitoring_slo" "primary" {
	service = google_monitoring_service.srv.service_id
	
	goal = 0.9
	rolling_period_days = 1
	
	basic_sli {
		availability {
		}
	}
}
`, randSuffix)
}

func testAccMonitoringSlo_availabilitySli() string {
	return `
data "google_monitoring_app_engine_service" "ae" {
  module_id = "default"
}

resource "google_monitoring_slo" "primary" {
  service = data.google_monitoring_app_engine_service.ae.service_id

  goal = 0.9
  rolling_period_days = 1

  basic_sli {
	availability {
	}
  }
}
`
}

func testAccMonitoringSloForSli(randSuffix, sliConfig string) string {
	return fmt.Sprintf(`
resource "google_monitoring_custom_service" "srv" {
  service_id = "tf-test-custom-srv-%s"
  display_name = "My Custom Service"
}

resource "google_monitoring_slo" "test_slo" {
  service = google_monitoring_custom_service.srv.service_id
  display_name = "Terraform Test SLO"

  goal = 0.9
  rolling_period_days = 30

  %s


}
`, randSuffix, sliConfig)
}

func testAccMonitoringSloSli_requestBasedDistributionMaxOnly() string {
	return `
request_based_sli {
	distribution_cut {
		distribution_filter = join(" AND ", [
			"metric.type=\"serviceruntime.googleapis.com/api/request_latencies\"",
			"resource.type=\"consumed_api\"",
		])
		range {
			max = 10
		}
	}
}
`
}

func testAccMonitoringSloSli_requestBasedGoodTotalRatio_goodAndTotal() string {
	return `
request_based_sli {
	good_total_ratio {
		good_service_filter = join(" AND ", [
			"metric.type=\"serviceruntime.googleapis.com/api/request_count\"",
			"resource.type=\"consumed_api\"",
			"metric.label.\"response_code\"=\"200\"",
		])
		total_service_filter = join(" AND ", [
			"metric.type=\"serviceruntime.googleapis.com/api/request_count\"",
			"resource.type=\"consumed_api\"",
		])
	}
}
`
}

func testAccMonitoringSloSli_requestBasedGoodTotalRatio_goodAndBad() string {
	return `
request_based_sli {
	good_total_ratio {
		good_service_filter = join(" AND ", [
			"metric.type=\"serviceruntime.googleapis.com/api/request_count\"",
			"resource.type=\"consumed_api\"",
			"metric.label.\"response_code\"=\"200\"",
		])
		bad_service_filter = join(" AND ", [
			"metric.type=\"serviceruntime.googleapis.com/api/request_count\"",
			"resource.type=\"consumed_api\"",
			"metric.label.\"response_code\"=\"400\"",
		])
	}
}
`
}

func testAccMonitoringSloSli_windowBasedSliGoodBadMetricFilter() string {
	return fmt.Sprintf(`
windows_based_sli {
  window_period = "1200s"
  good_bad_metric_filter =  join(" AND ", [
    "metric.type=\"monitoring.googleapis.com/uptime_check/check_passed\"",
    "resource.type=\"uptime_url\"",
  ])
}
`)
}

func testAccMonitoringSloSli_windowBasedSliGoodTotalRatioThreshold_distributionCut() string {
	return fmt.Sprintf(`
windows_based_sli {
  window_period = "400s"
	good_total_ratio_threshold {
		threshold = 0.1
		performance {
			distribution_cut {
				distribution_filter = join(" AND ", [
					"metric.type=\"serviceruntime.googleapis.com/api/request_latencies\"",
					"resource.type=\"consumed_api\"",
				])
	
				range {
					min = 1
					max = 9
				}
			}
		}
	}
}
`)
}

func testAccMonitoringSloSli_windowBasedSliGoodTotalRatioThreshold_distributionCutMaxOnly() string {
	return fmt.Sprintf(`
windows_based_sli {
  window_period = "2400s"
	good_total_ratio_threshold {
		threshold = 0.1
		performance {
			distribution_cut {
				distribution_filter = join(" AND ", [
					"metric.type=\"serviceruntime.googleapis.com/api/request_latencies\"",
					"resource.type=\"consumed_api\"",
				])
	
				range {
					max = 9
				}
			}
		}
	}
}
`)
}

func testAccMonitoringSloSli_windowBasedSliGoodTotalRatioThreshold_goodTotal() string {
	return fmt.Sprintf(`
windows_based_sli {
  window_period = "2400s"
	good_total_ratio_threshold {
		threshold = 0.1
		performance {
			good_total_ratio {
				good_service_filter = join(" AND ", [
					"metric.type=\"serviceruntime.googleapis.com/api/request_count\"",
					"resource.type=\"consumed_api\"",
					"metric.label.\"response_code\"=\"200\"",
				])
				total_service_filter = join(" AND ", [
					"metric.type=\"serviceruntime.googleapis.com/api/request_count\"",
					"resource.type=\"consumed_api\"",
				])
			}
		}
	}
}
`)
}

func testAccMonitoringSloSli_windowBasedSliGoodTotalRatioThreshold_goodBad() string {
	return fmt.Sprintf(`
windows_based_sli {
  window_period = "2400s"
	good_total_ratio_threshold {
		threshold = 0.1
		performance {
			good_total_ratio {
				good_service_filter = join(" AND ", [
					"metric.type=\"serviceruntime.googleapis.com/api/request_count\"",
					"resource.type=\"consumed_api\"",
					"metric.label.\"response_code\"=\"200\"",
				])
				bad_service_filter = join(" AND ", [
				"metric.type=\"serviceruntime.googleapis.com/api/request_count\"",
				"resource.type=\"consumed_api\"",
				"metric.label.\"response_code\"=\"400\"",
			])
			}
		}
	}
}
`)
}

func testAccMonitoringSloSli_windowBasedSliMetricMeanRange() string {
	return fmt.Sprintf(`
windows_based_sli {
	window_period = "600s"
	metric_mean_in_range {
		time_series = join(" AND ", [
			"metric.type=\"agent.googleapis.com/cassandra/client_request/latency/95p\"",
			"resource.type=\"gce_instance\"",
			])
		
		range {
			max = 50000000
		}
	}
}
`)
}

func testAccMonitoringSloSli_windowBasedSliMetricMeanRangeUpdate() string {
	return `
windows_based_sli {
	window_period = "600s"
	metric_mean_in_range {
		time_series = join(" AND ", [
			"metric.type=\"agent.googleapis.com/cassandra/client_request/latency/99p\"",
			"resource.type=\"gce_instance\"",
			])
		
		range {
			min = 1
			max = 70000000
		}
	}
}
`
}

func testAccMonitoringSloSli_windowBasedSliMetricSumRange() string {
	return fmt.Sprintf(`
windows_based_sli {
	window_period = "600s"
	metric_sum_in_range {
		time_series = join(" AND ", [
			"metric.type=\"monitoring.googleapis.com/uptime_check/request_latency\"",
			"resource.type=\"uptime_url\"",
		])

		range {
			max = 5000
		}
	}
}
`)
}

func testAccMonitoringSloSli_windowBasedSliMetricSumRangeUpdate() string {
	return `
windows_based_sli {
	window_period = "600s"
	metric_sum_in_range {
		time_series = join(" AND ", [
			"metric.type=\"monitoring.googleapis.com/uptime_check/request_latency\"",
			"resource.type=\"gce_instance\"",
		])

		range {
			min = 10
			max = 6000
		}
	}
}
`
}

func testAccMonitoringSloSli_windowBasedSliGoodTotalRatioThreshold_basicSli() string {
	return fmt.Sprintf(`
data "google_monitoring_app_engine_service" "ae" {
	module_id = "default"
}
	  
resource "google_monitoring_slo" "test_slo" {
	service = data.google_monitoring_app_engine_service.ae.service_id
	goal = 0.9
	rolling_period_days = 30
	windows_based_sli {
		window_period = "400s"
		good_total_ratio_threshold {
			threshold = 0.1
			basic_sli_performance {
				availability {
				}
			}
		}
    }
}`)
}
