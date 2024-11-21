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

func TestAccEphemeralServiceAccountJwt_basic(t *testing.T) {
	t.Parallel()

	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "jwt-basic", serviceAccount)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountJwt_basic(targetServiceAccountEmail),
			},
		},
	})
}

func TestAccEphemeralServiceAccountJwt_withDelegates(t *testing.T) {
	t.Parallel()

	initialServiceAccount := envvar.GetTestServiceAccountFromEnv(t)
	delegateServiceAccountEmailOne := acctest.BootstrapServiceAccount(t, "jwt-delegate1", initialServiceAccount)          // SA_2
	delegateServiceAccountEmailTwo := acctest.BootstrapServiceAccount(t, "jwt-delegate2", delegateServiceAccountEmailOne) // SA_3
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "jwt-target", delegateServiceAccountEmailTwo)         // SA_4

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountJwt_withDelegates(delegateServiceAccountEmailOne, delegateServiceAccountEmailTwo, targetServiceAccountEmail),
			},
		},
	})
}

func TestAccEphemeralServiceAccountJwt_withExpiresIn(t *testing.T) {
	t.Parallel()

	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "expiry", serviceAccount)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountJwt_withExpiresIn(targetServiceAccountEmail),
			},
		},
	})
}

func testAccEphemeralServiceAccountJwt_basic(serviceAccountEmail string) string {
	return fmt.Sprintf(`
ephemeral "google_service_account_jwt" "jwt" {
  target_service_account = "%s"
  payload               = jsonencode({
    "sub": "%[1]s",
    "aud": "https://example.com"
  })
}
`, serviceAccountEmail)
}

func testAccEphemeralServiceAccountJwt_withDelegates(delegateServiceAccountEmailOne, delegateServiceAccountEmailTwo, targetServiceAccountEmail string) string {
	return fmt.Sprintf(`
ephemeral "google_service_account_jwt" "jwt" {
  target_service_account = "%s"
  delegates = [
    "%s",
    "%s",
  ]
  payload               = jsonencode({
    "sub": "%[1]s",
    "aud": "https://example.com"
  })
}
# The delegation chain is:
# SA_1 (initialServiceAccountEmail) -> SA_2 (delegateServiceAccountEmailOne) -> SA_3 (delegateServiceAccountEmailTwo) -> SA_4 (targetServiceAccountEmail)
`, targetServiceAccountEmail, delegateServiceAccountEmailOne, delegateServiceAccountEmailTwo)
}

func testAccEphemeralServiceAccountJwt_withExpiresIn(serviceAccountEmail string) string {
	return fmt.Sprintf(`
ephemeral "google_service_account_jwt" "jwt" {
  target_service_account = "%s"
  expires_in            = 3600
  payload               = jsonencode({
    "sub": "%[1]s",
    "aud": "https://example.com"
  })
}
`, serviceAccountEmail)
}
