// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package dataplex_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccDataplexDatascan_dataplexDatascanBasicProfileExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataplexDatascanDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataplexDatascan_dataplexDatascanBasicProfileExample(context),
			},
			{
				ResourceName:            "google_dataplex_datascan.basic_profile",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "data_scan_id"},
			},
		},
	})
}

func testAccDataplexDatascan_dataplexDatascanBasicProfileExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataplex_datascan" "basic_profile" {
  location     = "us-central1"
  data_scan_id = "tf-test-datascan%{random_suffix}"

  data {
	  resource = "//bigquery.googleapis.com/projects/bigquery-public-data/datasets/samples/tables/shakespeare"
  }

  execution_spec {
    trigger {
      on_demand {}
    }
  }

data_profile_spec {}

  project = "%{project_name}"
}
`, context)
}

func TestAccDataplexDatascan_dataplexDatascanFullProfileExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataplexDatascanDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataplexDatascan_dataplexDatascanFullProfileExample(context),
			},
			{
				ResourceName:            "google_dataplex_datascan.full_profile",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "data_scan_id"},
			},
		},
	})
}

func testAccDataplexDatascan_dataplexDatascanFullProfileExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataplex_datascan" "full_profile" {
  location     = "us-central1"
  display_name = "Full Datascan Profile"
  data_scan_id = "tf-test-datascan%{random_suffix}"
  description  = "Example resource - Full Datascan Profile"
  labels = {
    author = "billing"
  }

  data {
    resource = "//bigquery.googleapis.com/projects/bigquery-public-data/datasets/samples/tables/shakespeare"
  }

  execution_spec {
    trigger {
      schedule {
        cron = "TZ=America/New_York 1 1 * * *"
      }
    }
  }

  data_profile_spec {
    sampling_percent = 80
    row_filter = "word_count > 10"
    include_fields {
      field_names = ["word_count"]
    }
    exclude_fields {
      field_names = ["property_type"]
    }
  }

  project = "%{project_name}"
}
`, context)
}

func TestAccDataplexDatascan_dataplexDatascanBasicQualityExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataplexDatascanDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataplexDatascan_dataplexDatascanBasicQualityExample(context),
			},
			{
				ResourceName:            "google_dataplex_datascan.basic_quality",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "data_scan_id"},
			},
		},
	})
}

func testAccDataplexDatascan_dataplexDatascanBasicQualityExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataplex_datascan" "basic_quality" {
  location     = "us-central1"
  data_scan_id = "tf-test-datascan%{random_suffix}"

  data {
    resource = "//bigquery.googleapis.com/projects/bigquery-public-data/datasets/samples/tables/shakespeare"
  }

  execution_spec {
    trigger {
      on_demand {}
    }
  }

  data_quality_spec {
    rules {
      dimension = "VALIDITY"
      name = "rule1"
      description = "rule 1 for validity dimension"
      table_condition_expectation {
        sql_expression = "COUNT(*) > 0"
      }
    }
  }

  project = "%{project_name}"
}
`, context)
}

func TestAccDataplexDatascan_dataplexDatascanFullQualityExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataplexDatascanDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataplexDatascan_dataplexDatascanFullQualityExample(context),
			},
			{
				ResourceName:            "google_dataplex_datascan.full_quality",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "data_scan_id"},
			},
		},
	})
}

func testAccDataplexDatascan_dataplexDatascanFullQualityExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataplex_datascan" "full_quality" {
  location = "us-central1"
  display_name = "Full Datascan Quality"
  data_scan_id = "tf-test-datascan%{random_suffix}"
  description = "Example resource - Full Datascan Quality"
  labels = {
    author = "billing"
  }

  data {
    resource = "//bigquery.googleapis.com/projects/bigquery-public-data/datasets/austin_bikeshare/tables/bikeshare_stations"
  }

  execution_spec {
    trigger {
      schedule {
        cron = "TZ=America/New_York 1 1 * * *"
      }
    }
    field = "modified_date"
  }

  data_quality_spec {
    sampling_percent = 5
    row_filter = "station_id > 1000"
    rules {
      column = "address"
      dimension = "VALIDITY"
      threshold = 0.99
      non_null_expectation {}
    }

    rules {
      column = "council_district"
      dimension = "VALIDITY"
      ignore_null = true
      threshold = 0.9
      range_expectation {
        min_value = 1
        max_value = 10
        strict_min_enabled = true
        strict_max_enabled = false
      }
    }

    rules {
      column = "power_type"
      dimension = "VALIDITY"
      ignore_null = false
      regex_expectation {
        regex = ".*solar.*"
      }
    }

    rules {
      column = "property_type"
      dimension = "VALIDITY"
      ignore_null = false
      set_expectation {
        values = ["sidewalk", "parkland"]
      }
    }


    rules {
      column = "address"
      dimension = "UNIQUENESS"
      uniqueness_expectation {}
    }

    rules {
      column = "number_of_docks"
      dimension = "VALIDITY"
      statistic_range_expectation {
        statistic = "MEAN"
        min_value = 5
        max_value = 15
        strict_min_enabled = true
        strict_max_enabled = true
      }
    }

    rules {
      column = "footprint_length"
      dimension = "VALIDITY"
      row_condition_expectation {
        sql_expression = "footprint_length > 0 AND footprint_length <= 10"
      }
    }

    rules {
      dimension = "VALIDITY"
      table_condition_expectation {
        sql_expression = "COUNT(*) > 0"
      }
    }
  }


  project = "%{project_name}"
}
`, context)
}

func testAccCheckDataplexDatascanDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_dataplex_datascan" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{DataplexBasePath}}projects/{{project}}/locations/{{location}}/dataScans/{{data_scan_id}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("DataplexDatascan still exists at %s", url)
			}
		}

		return nil
	}
}
