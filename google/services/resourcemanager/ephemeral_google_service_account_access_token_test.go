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

func TestAccEphemeralServiceAccountToken_basic(t *testing.T) {
	t.Parallel()

	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "basic", serviceAccount)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.AccTestPreCheck(t) },
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountToken_basic(targetServiceAccountEmail),
			},
		},
	})
}

func TestAccEphemeralServiceAccountToken_withDelegates(t *testing.T) {
	t.Parallel()

	project := envvar.GetTestProjectFromEnv()
	initialServiceAccount := envvar.GetTestServiceAccountFromEnv(t)
	delegateServiceAccountEmailOne := acctest.BootstrapServiceAccount(t, "delegate1", initialServiceAccount)          // SA_2
	delegateServiceAccountEmailTwo := acctest.BootstrapServiceAccount(t, "delegate2", delegateServiceAccountEmailOne) // SA_3
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "target", delegateServiceAccountEmailTwo)         // SA_4

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.AccTestPreCheck(t) },
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountToken_withDelegates(initialServiceAccount, delegateServiceAccountEmailOne, delegateServiceAccountEmailTwo, targetServiceAccountEmail, project),
			},
		},
	})
}

func TestAccEphemeralServiceAccountToken_withCustomLifetime(t *testing.T) {
	t.Parallel()

	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "lifetime", serviceAccount)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.AccTestPreCheck(t) },
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountToken_withCustomLifetime(targetServiceAccountEmail),
			},
		},
	})
}

func testAccEphemeralServiceAccountToken_basic(serviceAccountEmail string) string {
	return fmt.Sprintf(`
ephemeral "google_service_account_access_token" "token" {
  target_service_account = "%s"
  scopes                = ["https://www.googleapis.com/auth/cloud-platform"]
}
`, serviceAccountEmail)
}

func testAccEphemeralServiceAccountToken_withDelegates(initialServiceAccountEmail, delegateServiceAccountEmailOne, delegateServiceAccountEmailTwo, targetServiceAccountEmail, project string) string {
	return fmt.Sprintf(`
ephemeral "google_service_account_access_token" "test" {
  target_service_account = "%s"
  delegates = [
    "%s",
    "%s",
  ]
  scopes = ["https://www.googleapis.com/auth/cloud-platform"]
  lifetime = "3600s"
}

# The delegation chain is:
# SA_1 (initialServiceAccountEmail) -> SA_2 (delegateServiceAccountEmailOne) -> SA_3 (delegateServiceAccountEmailTwo) -> SA_4 (targetServiceAccountEmail)
`, targetServiceAccountEmail, delegateServiceAccountEmailOne, delegateServiceAccountEmailTwo)
}

func testAccEphemeralServiceAccountToken_withCustomLifetime(serviceAccountEmail string) string {
	return fmt.Sprintf(`
ephemeral "google_service_account_access_token" "token" {
  target_service_account = "%s"
  scopes                = ["https://www.googleapis.com/auth/cloud-platform"]
  lifetime              = "3600s"
}
`, serviceAccountEmail)
}
