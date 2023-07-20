// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: DCL     ***
//
// ----------------------------------------------------------------------------
//
//     This file is managed by Magic Modules (https://github.com/GoogleCloudPlatform/magic-modules)
//     and is based on the DCL (https://github.com/GoogleCloudPlatform/declarative-resource-client-library).
//     Changes will need to be made to the DCL or Magic Modules instead of here.
//
//     We are not currently able to accept contributions to this file. If changes
//     are required, please file an issue at https://github.com/hashicorp/terraform-provider-google/issues/new/choose
//
// ----------------------------------------------------------------------------

package dataplex_test

import (
	"context"
	"fmt"
	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	dataplex "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/dataplex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccDataplexAsset_BasicAssetHandWritten(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"region":        envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataplexAssetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataplexAsset_BasicAssetHandWritten(context),
			},
			{
				ResourceName:            "google_dataplex_asset.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resource_spec.0.name"},
			},
			{
				Config: testAccDataplexAsset_BasicAssetHandWrittenUpdate0(context),
			},
			{
				ResourceName:            "google_dataplex_asset.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resource_spec.0.name"},
			},
		},
	})
}

func testAccDataplexAsset_BasicAssetHandWritten(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "basic_bucket" {
  name          = "tf-test-bucket%{random_suffix}"
  location      = "%{region}"
  uniform_bucket_level_access = true
  lifecycle {
    ignore_changes = [
      labels
    ]
  }
 
  project = "%{project_name}"
}
 
resource "google_dataplex_lake" "basic_lake" {
  name         = "tf-test-lake%{random_suffix}"
  location     = "%{region}"
  project = "%{project_name}"
}
 
 
resource "google_dataplex_zone" "basic_zone" {
  name         = "tf-test-zone%{random_suffix}"
  location     = "%{region}"
  lake = google_dataplex_lake.basic_lake.name
  type = "RAW"
 
  discovery_spec {
    enabled = false
  }
 
 
  resource_spec {
    location_type = "SINGLE_REGION"
  }
 
  project = "%{project_name}"
}
 
 
resource "google_dataplex_asset" "primary" {
  name          = "tf-test-asset%{random_suffix}"
  location      = "%{region}"
 
  lake = google_dataplex_lake.basic_lake.name
  dataplex_zone = google_dataplex_zone.basic_zone.name
 
  discovery_spec {
    enabled = false
  }
 
  resource_spec {
    name = "projects/%{project_name}/buckets/tf-test-bucket%{random_suffix}"
    type = "STORAGE_BUCKET"
  }
 
  project = "%{project_name}"
  depends_on = [
    google_storage_bucket.basic_bucket
  ]
}
`, context)
}

func testAccDataplexAsset_BasicAssetHandWrittenUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "basic_bucket" {
  name          = "tf-test-bucket%{random_suffix}"
  location      = "%{region}"
  uniform_bucket_level_access = true
  lifecycle {
    ignore_changes = [
      labels
    ]
  }
 
  project = "%{project_name}"
}
 
resource "google_dataplex_lake" "basic_lake" {
  name         = "tf-test-lake%{random_suffix}"
  location     = "%{region}"
  project = "%{project_name}"
}
 
 
resource "google_dataplex_zone" "basic_zone" {
  name         = "tf-test-zone%{random_suffix}"
  location     = "%{region}"
  lake = google_dataplex_lake.basic_lake.name
  type = "RAW"
 
  discovery_spec {
    enabled = false
  }
 
 
  resource_spec {
    location_type = "SINGLE_REGION"
  }
 
  project = "%{project_name}"
}
 
 
resource "google_dataplex_asset" "primary" {
  name          = "tf-test-asset%{random_suffix}"
  location      = "%{region}"
 
  lake = google_dataplex_lake.basic_lake.name
  dataplex_zone = google_dataplex_zone.basic_zone.name
 
  discovery_spec {
    enabled = false
  }
 
  resource_spec {
    name = "projects/%{project_name}/buckets/tf-test-bucket%{random_suffix}"
    type = "STORAGE_BUCKET"
  }
 
  project = "%{project_name}"
  depends_on = [
    google_storage_bucket.basic_bucket
  ]
}
`, context)
}

func testAccCheckDataplexAssetDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "rs.google_dataplex_asset" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			billingProject := ""
			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			obj := &dataplex.Asset{
				DataplexZone: dcl.String(rs.Primary.Attributes["dataplex_zone"]),
				Lake:         dcl.String(rs.Primary.Attributes["lake"]),
				Location:     dcl.String(rs.Primary.Attributes["location"]),
				Name:         dcl.String(rs.Primary.Attributes["name"]),
				Description:  dcl.String(rs.Primary.Attributes["description"]),
				DisplayName:  dcl.String(rs.Primary.Attributes["display_name"]),
				Project:      dcl.StringOrNil(rs.Primary.Attributes["project"]),
				CreateTime:   dcl.StringOrNil(rs.Primary.Attributes["create_time"]),
				State:        dataplex.AssetStateEnumRef(rs.Primary.Attributes["state"]),
				Uid:          dcl.StringOrNil(rs.Primary.Attributes["uid"]),
				UpdateTime:   dcl.StringOrNil(rs.Primary.Attributes["update_time"]),
			}

			client := transport_tpg.NewDCLDataplexClient(config, config.UserAgent, billingProject, 0)
			_, err := client.GetAsset(context.Background(), obj)
			if err == nil {
				return fmt.Errorf("google_dataplex_asset still exists %v", obj)
			}
		}
		return nil
	}
}
