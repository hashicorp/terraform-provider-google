// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package accesscontextmanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Since each test here is acting on the same organization and only one AccessPolicy
// can exist, they need to be run serially. See AccessPolicy for the test runner.

func testAccAccessContextManagerAccessLevels_basicTest(t *testing.T) {
	org := envvar.GetTestOrgFromEnv(t)
	vpcName := fmt.Sprintf("test-vpc-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAccessContextManagerAccessLevelsDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAccessContextManagerAccessLevels_basic(org, "my policy", "corpnet_access", "prodnet_access"),
			},
			{
				ResourceName:      "google_access_context_manager_access_levels.test-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAccessContextManagerAccessLevels_basicUpdated(org, "my new policy", "corpnet_access", "prodnet_access", vpcName),
			},
			{
				ResourceName:      "google_access_context_manager_access_levels.test-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAccessContextManagerAccessLevel_empty(org, "my new policy"),
			},
			{
				ResourceName:      "google_access_context_manager_access_levels.test-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAccessContextManagerAccessLevelsDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_access_context_manager_access_levels" {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{AccessContextManagerBasePath}}{{parent}}/accessLevels")
			if err != nil {
				return err
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("AccessLevels still exists at %s", url)
			}
		}

		return nil
	}
}

func testAccAccessContextManagerAccessLevels_basic(org, policyTitle, levelTitleName1, levelTitleName2 string) string {
	return fmt.Sprintf(`
resource "google_access_context_manager_access_policy" "test-access" {
  parent = "organizations/%s"
  title  = "%s"
}

resource "google_access_context_manager_access_levels" "test-access" {
  parent      = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}"

  access_levels {
	name        = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}/accessLevels/%s"
	title       = "%s"
	description = "hello"
	basic {
	  combining_function = "AND"
	  conditions {
	    ip_subnetworks = ["192.0.4.0/24"]
	  }
    }
  }

  access_levels {
	name        = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}/accessLevels/%s"
	title       = "%s"
	description = "hello again"
	basic {
	  conditions {
	    ip_subnetworks = ["176.0.2.0/24"]
	  }
    }
  }
}
`, org, policyTitle, levelTitleName1, levelTitleName1, levelTitleName2, levelTitleName2)
}

func testAccAccessContextManagerAccessLevels_basicUpdated(org, policyTitle, levelTitleName1, levelTitleName2, vpcName string) string {
	return fmt.Sprintf(`
resource "google_access_context_manager_access_policy" "test-access" {
  parent = "organizations/%s"
  title  = "%s"
}

resource "google_compute_network" "vpc_network" {
	name = "%s"
}

resource "google_access_context_manager_access_levels" "test-access" {
  parent      = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}"

  access_levels {
	name        = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}/accessLevels/%s"
	title       = "%s"
	description = "hello"
	basic {
	  combining_function = "AND"
	  conditions {
	    ip_subnetworks = ["192.0.2.0/24"]
	  }
    }
  }

  access_levels {
    name        = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}/accessLevels/%s"
    title       = "%s"
    description = "hello again"
    basic {
      conditions {
        vpc_network_sources {
          vpc_subnetwork {
            network = "//compute.googleapis.com/${google_compute_network.vpc_network.id}"
            vpc_ip_subnetworks = ["20.0.5.0/24"]
          }
        }
      }
    }
  }
}
`, org, policyTitle, vpcName, levelTitleName1, levelTitleName1, levelTitleName2, levelTitleName2)
}

func testAccAccessContextManagerAccessLevel_empty(org, policyTitle string) string {
	return fmt.Sprintf(`
resource "google_access_context_manager_access_policy" "test-access" {
  parent = "organizations/%s"
  title  = "%s"
}

resource "google_access_context_manager_access_levels" "test-access" {
  parent      = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}"
}
`, org, policyTitle)
}
