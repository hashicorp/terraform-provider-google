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

package networkconnectivity_test

import (
	"context"
	"fmt"
	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	networkconnectivity "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/networkconnectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccNetworkConnectivityHub_BasicHub(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkConnectivityHubDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConnectivityHub_BasicHub(context),
			},
			{
				ResourceName:      "google_network_connectivity_hub.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNetworkConnectivityHub_BasicHubUpdate0(context),
			},
			{
				ResourceName:      "google_network_connectivity_hub.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNetworkConnectivityHub_BasicHub(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_network_connectivity_hub" "primary" {
  name        = "tf-test-hub%{random_suffix}"
  description = "A sample hub"

  labels = {
    label-one = "value-one"
  }

  project = "%{project_name}"
}


`, context)
}

func testAccNetworkConnectivityHub_BasicHubUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_network_connectivity_hub" "primary" {
  name        = "tf-test-hub%{random_suffix}"
  description = "An updated sample hub"

  labels = {
    label-two = "value-one"
  }

  project = "%{project_name}"
}


`, context)
}

func testAccCheckNetworkConnectivityHubDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "rs.google_network_connectivity_hub" {
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

			obj := &networkconnectivity.Hub{
				Name:        dcl.String(rs.Primary.Attributes["name"]),
				Description: dcl.String(rs.Primary.Attributes["description"]),
				Project:     dcl.StringOrNil(rs.Primary.Attributes["project"]),
				CreateTime:  dcl.StringOrNil(rs.Primary.Attributes["create_time"]),
				State:       networkconnectivity.HubStateEnumRef(rs.Primary.Attributes["state"]),
				UniqueId:    dcl.StringOrNil(rs.Primary.Attributes["unique_id"]),
				UpdateTime:  dcl.StringOrNil(rs.Primary.Attributes["update_time"]),
			}

			client := transport_tpg.NewDCLNetworkConnectivityClient(config, config.UserAgent, billingProject, 0)
			_, err := client.GetHub(context.Background(), obj)
			if err == nil {
				return fmt.Errorf("google_network_connectivity_hub still exists %v", obj)
			}
		}
		return nil
	}
}
