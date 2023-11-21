// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package vmwareengine_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccVmwareengineNetwork_vmwareEngineNetworkUpdate(t *testing.T) {
	t.Parallel()
	context := map[string]interface{}{
		"region":          envvar.GetTestRegionFromEnv(),
		"random_suffix":   acctest.RandString(t, 10),
		"organization":    envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
	}

	configTemplate := vmwareEngineNetworkConfigTemplate(context)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVmwareengineNetworkDestroyProducer(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
			"time":   {},
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(configTemplate, "description1"),
			},
			{
				ResourceName:            "google_vmwareengine_network.default-nw",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name"},
			},
			{
				Config: fmt.Sprintf(configTemplate, "description2"),
			},
			{
				ResourceName:            "google_vmwareengine_network.default-nw",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name"},
			},
		},
	})
}

func vmwareEngineNetworkConfigTemplate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vmwareengine_network" "default-nw" {
  project     = google_project_service.acceptance.project
  name        = "%{region}-default"
  location    = "%{region}"
  type        = "LEGACY"
  description = "%s"
}

# there can be only 1 Legacy network per region for a given project, so creating new project to isolate tests.
resource "google_project" "acceptance" {
  name            = "tf-test-%{random_suffix}"
  project_id      = "tf-test-%{random_suffix}"
  org_id          = "%{organization}"
  billing_account = "%{billing_account}"
}

resource "google_project_service" "acceptance" {
  project  = google_project.acceptance.project_id
  service  = "vmwareengine.googleapis.com"

  # Needed for CI tests for permissions to propagate, should not be needed for actual usage
  depends_on = [time_sleep.wait_60_seconds]
}

resource "time_sleep" "wait_60_seconds" {
  depends_on = [google_project.acceptance]

  create_duration = "60s"
}
`, context)
}
