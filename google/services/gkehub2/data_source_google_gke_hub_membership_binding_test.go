// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package gkehub2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceGoogleGKEHub2MembershipBinding_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGKEHub2MembershipBindingDestroyProducer(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleGKEHub2MembershipBinding_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_gke_hub_membership_binding.example", "google_gke_hub_membership_binding.example"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleGKEHub2MembershipBinding_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gke_hub_membership" "example" {
  membership_id = "tf-test-membership%{random_suffix}"
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

data "google_gke_hub_membership_binding" "example" {
  location = google_gke_hub_membership_binding.example.location
  membership_id = google_gke_hub_membership_binding.example.membership_id
  membership_binding_id = google_gke_hub_membership_binding.example.membership_binding_id
}
`, context)
}
