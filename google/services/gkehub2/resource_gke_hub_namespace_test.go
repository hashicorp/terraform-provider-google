// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package gkehub2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccGKEHub2Namespace_gkehubNamespaceBasicExample_update(t *testing.T) {
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
				Config: testAccGKEHub2Namespace_gkehubNamespaceBasicExample_basic(context),
			},
			{
				ResourceName:            "google_gke_hub_namespace.namespace",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"scope_namespace_id", "scope", "scope_id", "scope", "labels", "terraform_labels"},
			},
			{
				Config: testAccGKEHub2Namespace_gkehubNamespaceBasicExample_update(context),
			},
			{
				ResourceName:            "google_gke_hub_namespace.namespace",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"scope_namespace_id", "scope", "scope_id", "scope", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccGKEHub2Namespace_gkehubNamespaceBasicExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gke_hub_scope" "namespace" {
  scope_id = "tf-test-scope%{random_suffix}"
}


resource "google_gke_hub_namespace" "namespace" { 
  scope_namespace_id = "tf-test-namespace%{random_suffix}"
  scope_id = "tf-test-scope%{random_suffix}"
  scope = "${google_gke_hub_scope.namespace.name}"
  namespace_labels = {
      keyb = "valueb"
      keya = "valuea"
      keyc = "valuec"
  }
  labels = {
      keyb = "valueb"
      keya = "valuea"
      keyc = "valuec" 
  }
  depends_on = [google_gke_hub_scope.namespace]
}
`, context)
}

func testAccGKEHub2Namespace_gkehubNamespaceBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gke_hub_scope" "namespace" {
  scope_id = "tf-test-scope%{random_suffix}"
}


resource "google_gke_hub_namespace" "namespace" { 
  scope_namespace_id = "tf-test-namespace%{random_suffix}"
  scope_id = "tf-test-scope%{random_suffix}"
  scope = "${google_gke_hub_scope.namespace.name}"
  namespace_labels = {
      updated_keyb = "updated_valueb"
      updated_keya = "updated_valuea"
      updated_keyc = "updated_valuec"
  }
  labels = {
      updated_keyb = "updated_valueb"
      updated_keya = "updated_valuea"
      updated_keyc = "updated_valuec" 
  }
  depends_on = [google_gke_hub_scope.namespace]
}
`, context)
}
