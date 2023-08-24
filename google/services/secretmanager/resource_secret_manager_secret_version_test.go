// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package secretmanager_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccSecretManagerSecretVersion_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerSecretVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecretManagerSecretVersion_basic(context),
			},
			{
				ResourceName:      "google_secret_manager_secret_version.secret-version-basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSecretManagerSecretVersion_disable(context),
			},
			{
				ResourceName:      "google_secret_manager_secret_version.secret-version-basic",
				ImportState:       true,
				ImportStateVerify: true,
				// at this point the secret data is disabled and so reading the data on import will
				// give an empty string
				ImportStateVerifyIgnore: []string{"secret_data"},
			},
			{
				Config: testAccSecretManagerSecretVersion_basic(context),
			},
			{
				ResourceName:      "google_secret_manager_secret_version.secret-version-basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSecretManagerSecretVersion_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "tf-test-secret-version-%{random_suffix}"
  
  labels = {
    label = "my-label"
  }

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "secret-version-basic" {
  secret = google_secret_manager_secret.secret-basic.name

  secret_data = "my-tf-test-secret%{random_suffix}"
  enabled = true
}
`, context)
}

func testAccSecretManagerSecretVersion_disable(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "tf-test-secret-version-%{random_suffix}"

  labels = {
    label = "my-label"
  }

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "secret-version-basic" {
  secret = google_secret_manager_secret.secret-basic.name

  secret_data = "my-tf-test-secret%{random_suffix}"
  enabled = false
}
`, context)
}
