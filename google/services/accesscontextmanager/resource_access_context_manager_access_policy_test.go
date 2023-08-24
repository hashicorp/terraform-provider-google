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
// can exist, they need to be run serially
func TestAccAccessContextManager(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"access_policy":                    testAccAccessContextManagerAccessPolicy_basicTest,
		"access_policy_scoped":             testAccAccessContextManagerAccessPolicy_scopedTest,
		"service_perimeter":                testAccAccessContextManagerServicePerimeter_basicTest,
		"service_perimeter_update":         testAccAccessContextManagerServicePerimeter_updateTest,
		"service_perimeter_resource":       testAccAccessContextManagerServicePerimeterResource_basicTest,
		"access_level":                     testAccAccessContextManagerAccessLevel_basicTest,
		"access_level_full":                testAccAccessContextManagerAccessLevel_fullTest,
		"access_level_custom":              testAccAccessContextManagerAccessLevel_customTest,
		"access_levels":                    testAccAccessContextManagerAccessLevels_basicTest,
		"access_level_condition":           testAccAccessContextManagerAccessLevelCondition_basicTest,
		"service_perimeter_egress_policy":  testAccAccessContextManagerServicePerimeterEgressPolicy_basicTest,
		"service_perimeter_ingress_policy": testAccAccessContextManagerServicePerimeterIngressPolicy_basicTest,
		"service_perimeters":               testAccAccessContextManagerServicePerimeters_basicTest,
		"gcp_user_access_binding":          testAccAccessContextManagerGcpUserAccessBinding_basicTest,
		"authorized_orgs_desc":             testAccAccessContextManagerAuthorizedOrgsDesc_basicTest,
	}

	for name, tc := range testCases {
		// shadow the tc variable into scope so that when
		// the loop continues, if t.Run hasn't executed tc(t)
		// yet, we don't have a race condition
		// see https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc(t)
		})
	}
}

func testAccAccessContextManagerAccessPolicy_basicTest(t *testing.T) {
	org := envvar.GetTestOrgFromEnv(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAccessContextManagerAccessPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAccessContextManagerAccessPolicy_basic(org, "my policy"),
			},
			{
				ResourceName:      "google_access_context_manager_access_policy.test-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAccessContextManagerAccessPolicy_basic(org, "my new policy"),
			},
			{
				ResourceName:      "google_access_context_manager_access_policy.test-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAccessContextManagerAccessPolicyDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_access_context_manager_access_policy" {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{AccessContextManagerBasePath}}accessPolicies/{{name}}")
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
				return fmt.Errorf("AccessPolicy still exists at %s", url)
			}
		}

		return nil
	}
}

func testAccAccessContextManagerAccessPolicy_basic(org, title string) string {
	return fmt.Sprintf(`
resource "google_access_context_manager_access_policy" "test-access" {
  parent = "organizations/%s"
  title  = "%s"
}
`, org, title)
}

func testAccAccessContextManagerAccessPolicy_scopedTest(t *testing.T) {
	org := envvar.GetTestOrgFromEnv(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAccessContextManagerAccessPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAccessContextManagerAccessPolicy_scoped(org, "scoped policy"),
			},
			{
				ResourceName:      "google_access_context_manager_access_policy.test-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAccessContextManagerAccessPolicy_scoped(org, title string) string {
	return fmt.Sprintf(`
data "google_project" "project" {
}
resource "google_access_context_manager_access_policy" "test-access" {
  parent = "organizations/%s"
  title  = "%s"
  scopes = ["projects/${data.google_project.project.number}"]
}
`, org, title)
}
