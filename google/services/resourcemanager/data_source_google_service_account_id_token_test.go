// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager_test

import (
	"context"
	"testing"

	"fmt"

	"google.golang.org/api/idtoken"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

const targetAudience = "https://foo.bar/"

func testAccCheckServiceAccountIdTokenValue(name, audience string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ms := s.RootModule()

		rs, ok := ms.Resources[name]
		if !ok {
			return fmt.Errorf("can't find %s in state", name)
		}

		v, ok := rs.Primary.Attributes["id_token"]
		if !ok {
			return fmt.Errorf("id_token not found")
		}

		_, err := idtoken.Validate(context.Background(), v, audience)
		if err != nil {
			return fmt.Errorf("token validation failed: %v", err)
		}

		return nil
	}
}

func TestAccDataSourceGoogleServiceAccountIdToken_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.google_service_account_id_token.default"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleServiceAccountIdToken_basic(targetAudience),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "target_audience", targetAudience),
					testAccCheckServiceAccountIdTokenValue(resourceName, targetAudience),
				),
			},
		},
	})
}

func testAccCheckGoogleServiceAccountIdToken_basic(targetAudience string) string {

	return fmt.Sprintf(`
data "google_service_account_id_token" "default" {
  target_audience = "%s"
}
`, targetAudience)
}

func TestAccDataSourceGoogleServiceAccountIdToken_impersonation(t *testing.T) {
	t.Parallel()

	resourceName := "data.google_service_account_id_token.default"
	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, envvar.GetTestProjectFromEnv(), serviceAccount)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleServiceAccountIdToken_impersonation_datasource(targetAudience, targetServiceAccountEmail),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "target_audience", targetAudience),
					testAccCheckServiceAccountIdTokenValue(resourceName, targetAudience),
				),
			},
		},
	})
}

func testAccCheckGoogleServiceAccountIdToken_impersonation_datasource(targetAudience string, targetServiceAccount string) string {

	return fmt.Sprintf(`
data "google_service_account_access_token" "default" {
	target_service_account = "%s"
	scopes                 = ["userinfo-email", "https://www.googleapis.com/auth/cloud-platform"]
	lifetime               = "30s"
}

provider google {
	alias  = "impersonated"
	access_token = data.google_service_account_access_token.default.access_token
}

data "google_service_account_id_token" "default" {
	provider = google.impersonated
	target_service_account = "%s"
	target_audience = "%s"
}
`, targetServiceAccount, targetServiceAccount, targetAudience)
}
