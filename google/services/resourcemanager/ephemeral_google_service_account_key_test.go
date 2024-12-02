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

func TestAccEphemeralServiceAccountKey_basic(t *testing.T) {
	t.Parallel()

	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "key-basic", serviceAccount)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountKey_setup(targetServiceAccountEmail),
			},
			{
				Config: testAccEphemeralServiceAccountKey_basic(targetServiceAccountEmail),
			},
		},
	})
}

func testAccEphemeralServiceAccountKey_setup(serviceAccount string) string {
	return fmt.Sprintf(`
resource "google_service_account_key" "key" {
  service_account_id = "%s"
  public_key_type = "TYPE_X509_PEM_FILE"
}
`, serviceAccount)
}

func testAccEphemeralServiceAccountKey_basic(serviceAccount string) string {
	return fmt.Sprintf(`
resource "google_service_account_key" "key" {
  service_account_id = "%s"
  public_key_type = "TYPE_X509_PEM_FILE"
}

ephemeral "google_service_account_key" "key" {
  name            = google_service_account_key.key.name
  public_key_type = "TYPE_X509_PEM_FILE"
}
`, serviceAccount)
}
