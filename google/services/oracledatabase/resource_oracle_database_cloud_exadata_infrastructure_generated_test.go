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

package oracledatabase_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccOracleDatabaseCloudExadataInfrastructure_oracledatabaseCloudExadataInfrastructureBasicExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       "oci-terraform-testing",
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckOracleDatabaseCloudExadataInfrastructureDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOracleDatabaseCloudExadataInfrastructure_oracledatabaseCloudExadataInfrastructureBasicExample(context),
			},
			{
				ResourceName:            "google_oracle_database_cloud_exadata_infrastructure.my-cloud-exadata",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cloud_exadata_infrastructure_id", "labels", "location", "terraform_labels"},
			},
		},
	})
}

func testAccOracleDatabaseCloudExadataInfrastructure_oracledatabaseCloudExadataInfrastructureBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_oracle_database_cloud_exadata_infrastructure" "my-cloud-exadata"{
  display_name = "OFake exadata displayname"
  cloud_exadata_infrastructure_id = "ofake-exadata"
  location = "us-east4"
  project = "%{project}"
  properties {
    shape = "Exadata.X9M"
    compute_count= "2"
    storage_count= "3"
  }
}
`, context)
}

func TestAccOracleDatabaseCloudExadataInfrastructure_oracledatabaseCloudExadataInfrastructureFullExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       "oci-terraform-testing",
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckOracleDatabaseCloudExadataInfrastructureDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOracleDatabaseCloudExadataInfrastructure_oracledatabaseCloudExadataInfrastructureFullExample(context),
			},
			{
				ResourceName:            "google_oracle_database_cloud_exadata_infrastructure.my-cloud-exadata",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cloud_exadata_infrastructure_id", "labels", "location", "terraform_labels"},
			},
		},
	})
}

func testAccOracleDatabaseCloudExadataInfrastructure_oracledatabaseCloudExadataInfrastructureFullExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_oracle_database_cloud_exadata_infrastructure" "my-cloud-exadata"{
  display_name = "OFake exadata displayname"
  cloud_exadata_infrastructure_id = "ofake-exadata-id"
  location = "us-east4"
  project = "%{project}"
  gcp_oracle_zone = "us-east4-b-r1"
  properties {
    shape = "Exadata.X9M"
    compute_count= "2"
    storage_count= "3"
    customer_contacts {
      email = "xyz@example.com"
    }
    maintenance_window {
      custom_action_timeout_mins       = "20"
      days_of_week                     = ["SUNDAY"]
      hours_of_day                     = [4]
      is_custom_action_timeout_enabled = "0"
      lead_time_week                   = "1"
      months                           = ["JANUARY","APRIL","MAY","OCTOBER"]
      patching_mode                    = "ROLLING"
      preference                       = "CUSTOM_PREFERENCE"
      weeks_of_month                   = [4]
    }
    total_storage_size_gb = "196608"
  }
  labels = {
    "label-one" = "value-one"
  }
}
`, context)
}

func testAccCheckOracleDatabaseCloudExadataInfrastructureDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_oracle_database_cloud_exadata_infrastructure" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{OracleDatabaseBasePath}}projects/{{project}}/locations/{{location}}/cloudExadataInfrastructures/{{cloud_exadata_infrastructure_id}}")
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
				return fmt.Errorf("OracleDatabaseCloudExadataInfrastructure still exists at %s", url)
			}
		}

		return nil
	}
}
