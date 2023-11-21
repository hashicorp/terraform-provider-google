// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package vmwareengine_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceVmwareEngineNetwork_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVmwareengineNetworkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVmwareEngineNetworkConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores("data.google_vmwareengine_network.ds", "google_vmwareengine_network.nw", map[string]struct{}{}),
				),
			},
		},
	})
}

func testAccDataSourceVmwareEngineNetworkConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vmwareengine_network" "nw" {
    name              = "tf-test-sample-network%{random_suffix}"
    location          = "global" # Standard network needs to be global
    type              = "STANDARD"
    description       = "VMwareEngine standard network sample"
}

data "google_vmwareengine_network" "ds" {
  name     = google_vmwareengine_network.nw.name
  location = "global"
  depends_on = [
    google_vmwareengine_network.nw,
  ]
}
`, context)
}
