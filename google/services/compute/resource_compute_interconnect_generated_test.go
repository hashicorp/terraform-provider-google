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

package compute_test

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

func TestAccComputeInterconnect_computeInterconnectBasicTestExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInterconnectDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInterconnect_computeInterconnectBasicTestExample(context),
			},
			{
				ResourceName:            "google_compute_interconnect.example-interconnect",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "terraform_labels"},
			},
		},
	})
}

func testAccComputeInterconnect_computeInterconnectBasicTestExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_compute_interconnect" "example-interconnect" {
  name                 = "tf-test-example-interconnect%{random_suffix}"
  customer_name        = "internal_customer" # Special customer only available for Google testing.
  interconnect_type    = "DEDICATED"
  link_type            = "LINK_TYPE_ETHERNET_10G_LR"
  location             = "https://www.googleapis.com/compute/v1/${data.google_project.project.id}/global/interconnectLocations/z2z-us-east4-zone1-lciadl-a" # Special location only available for Google testing.
  requested_link_count = 1
  admin_enabled        = true
  description          = "example description"
  macsec_enabled       = false
  noc_contact_email    = "user@example.com"
  requested_features   = []
  labels = {
    mykey = "myvalue"
  }
}
`, context)
}

func testAccCheckComputeInterconnectDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_interconnect" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{ComputeBasePath}}projects/{{project}}/global/interconnects/{{name}}")
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
				return fmt.Errorf("ComputeInterconnect still exists at %s", url)
			}
		}

		return nil
	}
}
