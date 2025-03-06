// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package iamworkforcepool_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccIAMWorkforcePoolOauthClientCredential_full(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIAMWorkforcePoolOauthClientCredentialDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIAMWorkforcePoolOauthClientCredential_full(context),
			},
			{
				ResourceName:            "google_iam_oauth_client_credential.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "oauth_client_credential_id", "oauthclient"},
			},
			{
				Config: testAccIAMWorkforcePoolOauthClientCredential_full_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_iam_oauth_client_credential.example", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_iam_oauth_client_credential.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "oauth_client_credential_id", "oauthclient"},
			},
			{
				Config: testAccIAMWorkforcePoolOauthClientCredential_full_cleanOptionalFields(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_iam_oauth_client_credential.example", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_iam_oauth_client_credential.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "oauth_client_credential_id", "oauthclient"},
			},
			// Set disabled to `true` so the client credential can be deleted
			{
				Config: testAccIAMWorkforcePoolOauthClientCredential_full_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_iam_oauth_client_credential.example", plancheck.ResourceActionUpdate),
					},
				},
			},
		},
	})
}

func testAccIAMWorkforcePoolOauthClientCredential_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_oauth_client" "oauth_client" {
  oauth_client_id           = "tf-test-example-client-id%{random_suffix}"
  location                  = "global"
  allowed_grant_types       = ["AUTHORIZATION_CODE_GRANT"]
  allowed_redirect_uris     = ["https://www.example.com"]
  allowed_scopes            = ["https://www.googleapis.com/auth/cloud-platform"]
  client_type               = "CONFIDENTIAL_CLIENT"
}

resource "google_iam_oauth_client_credential" "example" {
  oauthclient	                = google_iam_oauth_client.oauth_client.oauth_client_id
  location                      = google_iam_oauth_client.oauth_client.location
  oauth_client_credential_id    = "tf-test-cred-id%{random_suffix}"
  disabled                      = true
  display_name                  = "Display Name of credential"
}
`, context)
}

func testAccIAMWorkforcePoolOauthClientCredential_full_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_oauth_client" "oauth_client" {
  oauth_client_id           = "tf-test-example-client-id%{random_suffix}"
  location                  = "global"
  allowed_grant_types       = ["AUTHORIZATION_CODE_GRANT"]
  allowed_redirect_uris     = ["https://www.example.com"]
  allowed_scopes            = ["https://www.googleapis.com/auth/cloud-platform"]
  client_type               = "CONFIDENTIAL_CLIENT"
}

resource "google_iam_oauth_client_credential" "example" {
  oauthclient	                = google_iam_oauth_client.oauth_client.oauth_client_id
  location                      = google_iam_oauth_client.oauth_client.location
  oauth_client_credential_id    = "tf-test-cred-id%{random_suffix}"
  disabled                      = true
  display_name                  = "Updated displayName"
}
`, context)
}

func testAccIAMWorkforcePoolOauthClientCredential_full_cleanOptionalFields(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_oauth_client" "oauth_client" {
  oauth_client_id           = "tf-test-example-client-id%{random_suffix}"
  location                  = "global"
  allowed_grant_types       = ["AUTHORIZATION_CODE_GRANT"]
  allowed_redirect_uris     = ["https://www.example.com"]
  allowed_scopes            = ["https://www.googleapis.com/auth/cloud-platform"]
  client_type               = "CONFIDENTIAL_CLIENT"
}

resource "google_iam_oauth_client_credential" "example" {
  oauthclient	                = google_iam_oauth_client.oauth_client.oauth_client_id
  location                      = google_iam_oauth_client.oauth_client.location
  oauth_client_credential_id    = "tf-test-cred-id%{random_suffix}"
}
`, context)
}
