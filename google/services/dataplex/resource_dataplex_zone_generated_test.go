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

func TestAccDataplexZone_BasicZone(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"region":        envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataplexZoneDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataplexZone_BasicZone(context),
			},
			{
				ResourceName:      "google_dataplex_zone.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataplexZone_BasicZoneUpdate0(context),
			},
			{
				ResourceName:      "google_dataplex_zone.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataplexZone_BasicZone(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataplex_zone" "primary" {
  discovery_spec {
    enabled = false
  }

  lake     = google_dataplex_lake.basic.name
  location = "%{region}"
  name     = "tf-test-zone%{random_suffix}"

  resource_spec {
    location_type = "MULTI_REGION"
  }

  type         = "RAW"
  description  = "Zone for DCL"
  display_name = "Zone for DCL"
  labels       = {}
  project      = "%{project_name}"
}

resource "google_dataplex_lake" "basic" {
  location     = "%{region}"
  name         = "tf-test-lake%{random_suffix}"
  description  = "Lake for DCL"
  display_name = "Lake for DCL"

  labels = {
    my-lake = "exists"
  }

  project = "%{project_name}"
}


`, context)
}

func testAccDataplexZone_BasicZoneUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataplex_zone" "primary" {
  discovery_spec {
    enabled = false
  }

  lake     = google_dataplex_lake.basic.name
  location = "%{region}"
  name     = "tf-test-zone%{random_suffix}"

  resource_spec {
    location_type = "MULTI_REGION"
  }

  type         = "RAW"
  description  = "Zone for DCL Updated"
  display_name = "Zone for DCL"

  labels = {
    updated_label = "exists"
  }

  project = "%{project_name}"
}

resource "google_dataplex_lake" "basic" {
  location     = "%{region}"
  name         = "tf-test-lake%{random_suffix}"
  description  = "Lake for DCL"
  display_name = "Lake for DCL"

  labels = {
    my-lake = "exists"
  }

  project = "%{project_name}"
}


`, context)
}

func testAccCheckDataplexZoneDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "rs.google_dataplex_zone" {
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

			obj := &dataplex.Zone{
				Lake:        dcl.String(rs.Primary.Attributes["lake"]),
				Location:    dcl.String(rs.Primary.Attributes["location"]),
				Name:        dcl.String(rs.Primary.Attributes["name"]),
				Type:        dataplex.ZoneTypeEnumRef(rs.Primary.Attributes["type"]),
				Description: dcl.String(rs.Primary.Attributes["description"]),
				DisplayName: dcl.String(rs.Primary.Attributes["display_name"]),
				Project:     dcl.StringOrNil(rs.Primary.Attributes["project"]),
				CreateTime:  dcl.StringOrNil(rs.Primary.Attributes["create_time"]),
				State:       dataplex.ZoneStateEnumRef(rs.Primary.Attributes["state"]),
				Uid:         dcl.StringOrNil(rs.Primary.Attributes["uid"]),
				UpdateTime:  dcl.StringOrNil(rs.Primary.Attributes["update_time"]),
			}

			client := transport_tpg.NewDCLDataplexClient(config, config.UserAgent, billingProject, 0)
			_, err := client.GetZone(context.Background(), obj)
			if err == nil {
				return fmt.Errorf("google_dataplex_zone still exists %v", obj)
			}
		}
		return nil
	}
}
