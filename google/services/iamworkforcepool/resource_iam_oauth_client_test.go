// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package iamworkforcepool_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccIAMWorkforcePoolOauthClient_full(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIAMWorkforcePoolOauthClientDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIAMWorkforcePoolOauthClient_full(context),
			},
			{
				ResourceName:            "google_iam_oauth_client.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "oauth_client_id"},
			},
			{
				Config: testAccIAMWorkforcePoolOauthClient_full_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_iam_oauth_client.example", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_iam_oauth_client.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "oauth_client_id"},
			},
			{
				Config: testAccIAMWorkforcePoolOauthClient_full_cleanOptionalFields(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_iam_oauth_client.example", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_iam_oauth_client.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "oauth_client_id"},
			},
		},
	})
}

func testAccIAMWorkforcePoolOauthClient_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_oauth_client" "example" {
  oauth_client_id = "tf-test-example-client-id%{random_suffix}"
  display_name              = "Display Name of OAuth client"
  description               = "A sample OAuth client"
  location                  = "global"
  disabled                  = false
  allowed_grant_types       = ["AUTHORIZATION_CODE_GRANT"]
  allowed_redirect_uris     = ["https://www.example.com"]
  allowed_scopes            = ["https://www.googleapis.com/auth/cloud-platform"]
  client_type               = "CONFIDENTIAL_CLIENT"
}
`, context)
}

func testAccIAMWorkforcePoolOauthClient_full_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_oauth_client" "example" {
  oauth_client_id = "tf-test-example-client-id%{random_suffix}"
  display_name              = "Updated displayName"
  description               = "Updated description"
  location                  = "global"
  disabled                  = true
  allowed_grant_types       = ["AUTHORIZATION_CODE_GRANT"]
  allowed_redirect_uris     = ["https://www.update.com"]
  allowed_scopes            = ["https://www.googleapis.com/auth/cloud-platform", "openid"]
  client_type               = "CONFIDENTIAL_CLIENT"
}
`, context)
}

func testAccIAMWorkforcePoolOauthClient_full_cleanOptionalFields(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_oauth_client" "example" {
  oauth_client_id = "tf-test-example-client-id%{random_suffix}"
  location                  = "global"
  disabled                  = true
  allowed_grant_types       = ["AUTHORIZATION_CODE_GRANT"]
  allowed_redirect_uris     = ["https://www.update.com"]
  allowed_scopes            = ["https://www.googleapis.com/auth/cloud-platform", "openid"]
  client_type               = "CONFIDENTIAL_CLIENT"
}
`, context)
}
