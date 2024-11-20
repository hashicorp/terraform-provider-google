// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccEphemeralServiceAccountIdToken_basic(t *testing.T) {
	t.Parallel()

	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "idtoken", serviceAccount)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountIdToken_basic(targetServiceAccountEmail),
			},
		},
	})
}

func TestAccEphemeralServiceAccountIdToken_withDelegates(t *testing.T) {
	t.Parallel()

	initialServiceAccount := envvar.GetTestServiceAccountFromEnv(t)
	delegateServiceAccountEmailOne := acctest.BootstrapServiceAccount(t, "id-delegate1", initialServiceAccount)          // SA_2
	delegateServiceAccountEmailTwo := acctest.BootstrapServiceAccount(t, "id-delegate2", delegateServiceAccountEmailOne) // SA_3
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "id-target", delegateServiceAccountEmailTwo)         // SA_4

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountIdToken_withDelegates(delegateServiceAccountEmailOne, delegateServiceAccountEmailTwo, targetServiceAccountEmail),
			},
		},
	})
}

func TestAccEphemeralServiceAccountIdToken_withIncludeEmail(t *testing.T) {
	t.Parallel()

	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "idtoken-email", serviceAccount)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountIdToken_withIncludeEmail(targetServiceAccountEmail),
			},
		},
	})
}

func testAccEphemeralServiceAccountIdToken_basic(serviceAccountEmail string) string {
	return fmt.Sprintf(`
ephemeral "google_service_account_id_token" "token" {
  target_service_account = "%s"
  target_audience       = "https://example.com"
}
`, serviceAccountEmail)
}

func testAccEphemeralServiceAccountIdToken_withDelegates(delegateServiceAccountEmailOne, delegateServiceAccountEmailTwo, targetServiceAccountEmail string) string {
	return fmt.Sprintf(`
ephemeral "google_service_account_id_token" "token" {
  target_service_account = "%s"
  delegates = [
    "%s",
    "%s",
  ]
  target_audience       = "https://example.com"
}

# The delegation chain is:
# SA_1 (initialServiceAccountEmail) -> SA_2 (delegateServiceAccountEmailOne) -> SA_3 (delegateServiceAccountEmailTwo) -> SA_4 (targetServiceAccountEmail)
`, targetServiceAccountEmail, delegateServiceAccountEmailOne, delegateServiceAccountEmailTwo)
}

func testAccEphemeralServiceAccountIdToken_withIncludeEmail(serviceAccountEmail string) string {
	return fmt.Sprintf(`
ephemeral "google_service_account_id_token" "token" {
  target_service_account = "%s"
  target_audience       = "https://example.com"
  include_email        = true
}
`, serviceAccountEmail)
}
