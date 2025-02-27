// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceComputeRegionSslPolicy(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComputeRegionSslPolicyConfig(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState(
						"data.google_compute_region_ssl_policy.policy",
						"google_compute_region_ssl_policy.foobar",
					),
				),
			},
		},
	})
}

func testAccDataSourceComputeRegionSslPolicyConfig(policyName string) string {
	return fmt.Sprintf(`
resource "google_compute_region_ssl_policy" "foobar" {
  name            = "tf-test-policyds-%s"
  region          = "us-central1"
  profile         = "MODERN"
  min_tls_version = "TLS_1_2"
}

data "google_compute_region_ssl_policy" "policy" {
  name   = google_compute_region_ssl_policy.foobar.name
  region = google_compute_region_ssl_policy.foobar.region
}
`, policyName)
}
