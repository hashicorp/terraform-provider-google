// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package gkehub2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccGKEHub2Scope_gkehubScopeBasicExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGKEHub2Scope_gkehubScopeBasicExample_basic(context),
			},
			{
				ResourceName:            "google_gke_hub_scope.scope",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"scope_id", "labels", "terraform_labels"},
			},
			{
				Config: testAccGKEHub2Scope_gkehubScopeBasicExample_update(context),
			},
			{
				ResourceName:            "google_gke_hub_scope.scope",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"scope_id", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccGKEHub2Scope_gkehubScopeBasicExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gke_hub_scope" "scope" {
  scope_id = "tf-test-scope%{random_suffix}"
  labels = {
    keyb = "valueb"
    keya = "valuea"
    keyc = "valuec" 
  }
}
`, context)
}

func testAccGKEHub2Scope_gkehubScopeBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gke_hub_scope" "scope" {
  scope_id = "tf-test-scope%{random_suffix}"
  labels = {
    updated_keyb = "updated_valueb"
    updated_keya = "updated_valuea"
    updated_keyc = "updated_valuec" 
  }
}
`, context)
}
