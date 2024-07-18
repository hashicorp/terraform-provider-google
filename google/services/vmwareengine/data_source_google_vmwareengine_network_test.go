// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package vmwareengine_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceVmwareEngineNetwork_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":   acctest.RandString(t, 10),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckVmwareengineNetworkDestroyProducer(t),
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
resource "google_project" "project" {
  project_id      = "tf-test%{random_suffix}"
  name            = "tf-test%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}

resource "google_project_service" "vmwareengine" {
  project = google_project.project.project_id
  service = "vmwareengine.googleapis.com"
}

resource "time_sleep" "sleep" {
  create_duration = "1m"

  depends_on = [
    google_project_service.vmwareengine,
  ]
}

resource "google_vmwareengine_network" "nw" {
  project           = google_project.project.project_id
  name              = "tf-test-sample-network%{random_suffix}"
  location          = "global" # Standard network needs to be global
  type              = "STANDARD"
  description       = "VMwareEngine standard network sample"

  depends_on = [
    time_sleep.sleep # Sleep allows permissions in the new project to propagate
  ]
}

data "google_vmwareengine_network" "ds" {
  name     = google_vmwareengine_network.nw.name
  project  = google_project.project.project_id
  location = "global"
}
`, context)
}
