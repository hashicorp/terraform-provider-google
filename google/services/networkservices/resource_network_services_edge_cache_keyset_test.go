// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package networkservices_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccNetworkServicesEdgeCacheKeyset_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkServicesEdgeCacheKeysetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkServicesEdgeCacheKeyset_networkServicesEdgeCacheKeysetBasicExample(context),
			},
			{
				ResourceName:            "google_network_services_edge_cache_keyset.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
			{
				Config: testAccNetworkServicesEdgeCacheKeyset_update(context),
			},
			{
				ResourceName:            "google_network_services_edge_cache_keyset.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
		},
	})
}

func testAccNetworkServicesEdgeCacheKeyset_update(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_network_services_edge_cache_keyset" "default" {
  name                 = "default%{random_suffix}"
  description          = "T2"
  public_key {
    id = "my-public-key-2"
    value = "hzd03llxB1u5FOLKFkZ6_wCJqC7jtN0bg7xlBqS6WVM"
  }
	labels = {
		a = "a"
	}
}
`, context)
}
