// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceComputeSecurityPolicy_basic(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSecurityPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComputeSecurityPolicy_basic(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_compute_security_policy.sp1", "google_compute_security_policy.policy"),
					acctest.CheckDataSourceStateMatchesResourceState("data.google_compute_security_policy.sp2", "google_compute_security_policy.policy"),
				),
			},
		},
	})
}

func testAccDataSourceComputeSecurityPolicy_basic(suffix string) string {
	return fmt.Sprintf(`
resource "google_compute_security_policy" "policy" {
  name = "my-policy-%s"

  rule {
    action      = "deny(403)"
    priority    = "1000"
    description = "Deny access to IPs in 9.9.9.0/24"

    match {
      versioned_expr = "SRC_IPS_V1"

      config {
        src_ip_ranges = ["9.9.9.0/24"]
      }
    }
  }

  rule {
    action      = "allow"
    priority    = "2147483647"
    description = "default rule"

    match {
      versioned_expr = "SRC_IPS_V1"

      config {
        src_ip_ranges = ["*"]
      }
    }
  }
}

data "google_compute_security_policy" "sp1" {
  name    = google_compute_security_policy.policy.name
  project = google_compute_security_policy.policy.project
}

data "google_compute_security_policy" "sp2" {
  self_link = google_compute_security_policy.policy.self_link
}
`, suffix)
}
