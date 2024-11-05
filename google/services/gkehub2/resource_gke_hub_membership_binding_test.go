// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package gkehub2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccGKEHub2MembershipBinding_gkehubMembershipBindingBasicExample_update(t *testing.T) {
	// Currently failing
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"project":         envvar.GetTestProjectFromEnv(),
		"location":        envvar.GetTestRegionFromEnv(),
		"random_suffix":   acctest.RandString(t, 10),
		"network_name":    acctest.BootstrapSharedTestNetwork(t, "gke-cluster"),
		"subnetwork_name": acctest.BootstrapSubnet(t, "gke-cluster", acctest.BootstrapSharedTestNetwork(t, "gke-cluster")),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGKEHub2MembershipBinding_gkehubMembershipBindingBasicExample_basic(context),
			},
			{
				ResourceName:            "google_gke_hub_membership_binding.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"membership_binding_id", "scope", "membership_id", "location", "labels", "terraform_labels"},
			},
			{
				Config: testAccGKEHub2MembershipBinding_gkehubMembershipBindingBasicExample_update(context),
			},
			{
				ResourceName:            "google_gke_hub_membership_binding.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"membership_binding_id", "scope", "membership_id", "location", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccGKEHub2MembershipBinding_gkehubMembershipBindingBasicExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_container_cluster" "primary" {
  name               = "tf-test-basic-cluster%{random_suffix}"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection = false
  network       = "%{network_name}"
  subnetwork    = "%{subnetwork_name}"
}

resource "google_gke_hub_membership" "example" {
  membership_id = "tf-test-membership%{random_suffix}"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${google_container_cluster.primary.id}"
    }
  }
  
  depends_on = [google_container_cluster.primary]
}

resource "google_gke_hub_scope" "example" {
  scope_id = "tf-test-scope%{random_suffix}"
}

resource "google_gke_hub_membership_binding" "example" {
  membership_binding_id = "tf-test-membership-binding%{random_suffix}"
  scope = google_gke_hub_scope.example.name
  membership_id = "tf-test-membership%{random_suffix}"
  location = "global"
  labels = {
      keyb = "valueb"
      keya = "valuea"
      keyc = "valuec" 
  }
  depends_on = [
    google_gke_hub_membership.example,
    google_gke_hub_scope.example
  ]
}
`, context)
}

func testAccGKEHub2MembershipBinding_gkehubMembershipBindingBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_container_cluster" "primary" {
  name               = "tf-test-basic-cluster%{random_suffix}"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection = false
  network       = "%{network_name}"
  subnetwork    = "%{subnetwork_name}"
}

resource "google_gke_hub_membership" "example" {
  membership_id = "tf-test-membership%{random_suffix}"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${google_container_cluster.primary.id}"
    }
  }
  
  depends_on = [google_container_cluster.primary]
}

resource "google_gke_hub_scope" "example" {
  scope_id = "tf-test-scope%{random_suffix}"
}

resource "google_gke_hub_scope" "example2" {
  scope_id = "tf-test-scope2%{random_suffix}"
}

resource "google_gke_hub_membership_binding" "example" {
  membership_binding_id = "tf-test-membership-binding%{random_suffix}"
  scope = google_gke_hub_scope.example2.name
  membership_id = "tf-test-membership%{random_suffix}"
  location = "global"
  labels = {
      updated_keyb = "updated_valueb"
      updated_keya = "updated_valuea"
      updated_keyc = "updated_valuec" 
  }
  depends_on = [
    google_gke_hub_membership.example,
    google_gke_hub_scope.example2
  ]
}
`, context)
}
