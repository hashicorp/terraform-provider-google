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

package monitoring_test

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

func TestAccMonitoringGroup_monitoringGroupBasicExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringGroupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringGroup_monitoringGroupBasicExample(context),
			},
			{
				ResourceName:      "google_monitoring_group.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMonitoringGroup_monitoringGroupBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_monitoring_group" "basic" {
  display_name = "tf-test MonitoringGroup%{random_suffix}"

  filter = "resource.metadata.region=\"europe-west2\""
}
`, context)
}

func TestAccMonitoringGroup_monitoringGroupSubgroupExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringGroupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringGroup_monitoringGroupSubgroupExample(context),
			},
			{
				ResourceName:      "google_monitoring_group.subgroup",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMonitoringGroup_monitoringGroupSubgroupExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_monitoring_group" "parent" {
  display_name = "tf-test MonitoringParentGroup%{random_suffix}"
  filter       = "resource.metadata.region=\"europe-west2\""
}

resource "google_monitoring_group" "subgroup" {
  display_name = "tf-test MonitoringSubGroup%{random_suffix}"
  filter       = "resource.metadata.region=\"europe-west2\""
  parent_name  =  google_monitoring_group.parent.name
}
`, context)
}

func testAccCheckMonitoringGroupDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_monitoring_group" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{MonitoringBasePath}}v3/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:               config,
				Method:               "GET",
				Project:              billingProject,
				RawURL:               url,
				UserAgent:            config.UserAgent,
				ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.IsMonitoringConcurrentEditError},
			})
			if err == nil {
				return fmt.Errorf("MonitoringGroup still exists at %s", url)
			}
		}

		return nil
	}
}
