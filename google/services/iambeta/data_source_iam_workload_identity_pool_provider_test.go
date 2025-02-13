// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package iambeta_test

import (
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceIAMBetaWorkloadIdentityPoolProvider_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIAMBetaWorkloadIdentityPoolProviderDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIAMBetaWorkloadIdentityPoolProviderBasic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_iam_workload_identity_pool_provider.foo", "google_iam_workload_identity_pool_provider.bar"),
				),
			},
		},
	})
}

func testAccDataSourceIAMBetaWorkloadIdentityPoolProviderBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_workload_identity_pool" "pool" {
	workload_identity_pool_id = "pool-%{random_suffix}"
}

resource "google_iam_workload_identity_pool_provider" "bar" {
	workload_identity_pool_id          = google_iam_workload_identity_pool.pool.workload_identity_pool_id
	workload_identity_pool_provider_id = "bar-provider-%{random_suffix}"
	display_name                       = "Name of provider"
	description                        = "OIDC identity pool provider for automated test"
	disabled                           = true
	attribute_condition                = "\"e968c2ef-047c-498d-8d79-16ca1b61e77e\" in assertion.groups"
	attribute_mapping                  = {
		"google.subject" = "assertion.sub"
	}
	oidc {
		allowed_audiences = ["https://example.com/gcp-oidc-federation"]
		issuer_uri        = "https://sts.windows.net/azure-tenant-id"
	}
  }

data "google_iam_workload_identity_pool_provider" "foo" {
	workload_identity_pool_id          = google_iam_workload_identity_pool.pool.workload_identity_pool_id
	workload_identity_pool_provider_id = google_iam_workload_identity_pool_provider.bar.workload_identity_pool_provider_id
}
`, context)
}
