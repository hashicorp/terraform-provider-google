// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager_test

import (
	"testing"

	"fmt"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testAccCheckServiceAccountAccessTokenValue(name, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ms := s.RootModule()
		rs, ok := ms.Outputs[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		// TODO: validate the token belongs to the service account
		if rs.Value == "" {
			return fmt.Errorf("%s Cannot be empty", name)
		}

		return nil
	}
}

func TestAccDataSourceGoogleServiceAccountAccessToken_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.google_service_account_access_token.default"
	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, envvar.GetTestProjectFromEnv(), serviceAccount)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:  testAccCheckGoogleServiceAccountAccessToken_datasource(targetServiceAccountEmail),
				Destroy: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "target_service_account", targetServiceAccountEmail),
					testAccCheckServiceAccountAccessTokenValue("access_token", targetServiceAccountEmail),
				),
			},
		},
	})
}

func testAccCheckGoogleServiceAccountAccessToken_datasource(targetServiceAccountID string) string {

	return fmt.Sprintf(`
data "google_service_account_access_token" "default" {
  target_service_account = "%s"
  scopes                 = ["userinfo-email", "https://www.googleapis.com/auth/cloud-platform"]
  lifetime               = "30s"
}

output "access_token" {
  value = data.google_service_account_access_token.default.access_token
  sensitive = true
}
`, targetServiceAccountID)
}
