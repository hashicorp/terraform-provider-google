// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package accesscontextmanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Since each test here is acting on the same organization and only one AccessPolicy
// can exist, they need to be run serially. See AccessPolicy for the test runner.

func testAccAccessContextManagerServicePerimeterIngressPolicy_basicTest(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	org := envvar.GetTestOrgFromEnv(t)
	//projects := acctest.BootstrapServicePerimeterProjects(t, 1)
	policyTitle := acctest.RandString(t, 10)
	perimeterTitle := "perimeter"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAccessContextManagerServicePerimeterIngressPolicy_basic(org, policyTitle, perimeterTitle),
			},
			{
				ResourceName:      "google_access_context_manager_service_perimeter.test-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAccessContextManagerServicePerimeterIngressPolicy_destroy(org, policyTitle, perimeterTitle),
				Check:  testAccCheckAccessContextManagerServicePerimeterIngressPolicyDestroyProducer(t),
			},
		},
	})
}

func testAccCheckAccessContextManagerServicePerimeterIngressPolicyDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_access_context_manager_service_perimeter_ingress_policy" {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{AccessContextManagerBasePath}}{{perimeter}}")
			if err != nil {
				return err
			}

			res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err != nil {
				return err
			}

			v, ok := res["status"]
			if !ok || v == nil {
				return nil
			}

			res = v.(map[string]interface{})
			v, ok = res["perimeter"]
			if !ok || v == nil {
				return nil
			}

			resources := v.([]interface{})
			if len(resources) == 0 {
				return nil
			}

			return fmt.Errorf("expected 0 resources in perimeter, found %d: %v", len(resources), resources)
		}

		return nil
	}
}

func testAccAccessContextManagerServicePerimeterIngressPolicy_basic(org, policyTitle, perimeterTitleName string) string {
	return fmt.Sprintf(`
%s

resource "google_access_context_manager_service_perimeter_ingress_policy" "test-access1" {
  perimeter = google_access_context_manager_service_perimeter.test-access.name
	ingress_from {
		identity_type = "ANY_IDENTITY"
	}
	ingress_to {
		resources = [ "*" ]
		operations {
			service_name = "bigquery.googleapis.com"
	
			method_selectors {
				method = "BigQueryStorage.ReadRows"
			}
	
			method_selectors {
				method = "TableService.ListTables"
			}
	
			method_selectors {
				permission = "bigquery.jobs.get"
			}
		}
	
		operations {
			service_name = "storage.googleapis.com"
	
			method_selectors {
				method = "google.storage.objects.create"
			}
		}
  	}
}

resource "google_access_context_manager_service_perimeter_ingress_policy" "test-access2" {
	perimeter = google_access_context_manager_service_perimeter.test-access.name
	ingress_from {
		identity_type = "ANY_IDENTITY"
	}
}

`, testAccAccessContextManagerServicePerimeterIngressPolicy_destroy(org, policyTitle, perimeterTitleName))
}

func testAccAccessContextManagerServicePerimeterIngressPolicy_destroy(org, policyTitle, perimeterTitleName string) string {
	return fmt.Sprintf(`
resource "google_access_context_manager_access_policy" "test-access" {
  parent = "organizations/%s"
  title  = "%s"
}

resource "google_access_context_manager_service_perimeter" "test-access" {
  parent         = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}"
  name           = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}/servicePerimeters/%s"
  title          = "%s"
  status {
    restricted_services = ["storage.googleapis.com"]
  }

  lifecycle {
  	ignore_changes = [status[0].ingress_policies]
  }
}
`, org, policyTitle, perimeterTitleName, perimeterTitleName)
}
