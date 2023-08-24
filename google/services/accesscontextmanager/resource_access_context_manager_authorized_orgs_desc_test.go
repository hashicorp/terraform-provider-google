// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package accesscontextmanager_test

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

func testAccAccessContextManagerAuthorizedOrgsDesc_basicTest(t *testing.T) {
	context := map[string]interface{}{
		"org_id": envvar.GetTestOrgFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAccessContextManagerAuthorizedOrgsDescDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAccessContextManagerAuthorizedOrgsDesc_accessContextManagerAuthorizedOrgsDescBasicExample(context),
			},
			{
				ResourceName:            "google_access_context_manager_authorized_orgs_desc.authorized-orgs-desc",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent"},
			},
		},
	})
}

func testAccAccessContextManagerAuthorizedOrgsDesc_accessContextManagerAuthorizedOrgsDescBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_access_context_manager_authorized_orgs_desc" "authorized-orgs-desc" {
  parent = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}"
  name   = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}/authorizedOrgsDescs/fakeDescName"
  authorization_type = "AUTHORIZATION_TYPE_TRUST"
  asset_type = "ASSET_TYPE_CREDENTIAL_STRENGTH"
  authorization_direction = "AUTHORIZATION_DIRECTION_TO"
  orgs = ["organizations/12345", "organizations/98765"]
}

resource "google_access_context_manager_access_policy" "test-access" {
  parent = "organizations/%{org_id}"
  title  = "my policy"
}
`, context)
}

func testAccCheckAccessContextManagerAuthorizedOrgsDescDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_access_context_manager_authorized_orgs_desc" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{AccessContextManagerBasePath}}{{name}}")
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
				return fmt.Errorf("AccessContextManagerAuthorizedOrgsDesc still exists at %s", url)
			}
		}

		return nil
	}
}
